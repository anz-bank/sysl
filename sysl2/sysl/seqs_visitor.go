package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
)

type (
	visitor interface {
		visit(s element)
	}

	element interface {
		accept(v visitor)
	}

	endpointElement struct {
		// read only
		fromApp                *sysl.AppName
		appName                string
		endpointName           string
		senderPatterns         *setS
		senderEndpointPatterns *setS
		stmt                   *sysl.Statement
		deactivate             bool
		// belows are read write
		application      *sysl.Application
		endpoint         *sysl.Endpoint
		appPatterns      *setS
		endpointPatterns *setS
		sender           string
		agent            string
		payload          string
	}

	statementsElement struct {
		endpointElement
		stmts          []*sysl.Statement
		deactivate     bool
		lastParentStmt bool
	}

	sdVisitor struct {
		w               *textWriter
		pm              *participantManager
		m               *sysl.Module
		already_visited map[string]int
		uptos           map[string]string
		fmtEndpoint     SfmtEP
	}
)

func (s *endpointElement) accept(v visitor) {
	v.visit(s)
}

func (s *statementsElement) accept(v visitor) {
	v.visit(s)
}

func (s *sdVisitor) visit(elem element) {
	switch t := elem.(type) {
	case *endpointElement:
		s.visitEndpoint(t)
	case *statementsElement:
		s.visitStmts(t)
	}
}

func (s *sdVisitor) visitEndpoint(e *endpointElement) {
	if !s.normalizeApplication(e) {
		logrus.Errorf("Application with name %s does not found.", e.appName)
		return
	}
	if !s.normalizeEndpoint(e) {
		logrus.Errorf("Endpoint with name %s does not found in Application with name %s", e.endpointName, e.appName)
		return
	}

	s.normalizeSender(e)
	s.normalizeAgent(e)
	s.normalizeAppPatterns(e)
	s.normalizeEndpointPatterns(e)
	s.normalizePayload(e)

	s.testCallingTimer(e)
	s.testCallingSelf(e)

	if len(e.endpoint.Stmt) > 0 {
		hitBlackbox := s.testHitBlackbox(e)

		if !hitBlackbox {
			active := !e.appPatterns.contains("human") && !e.appPatterns.contains("cron")
			if active {
				s.w.activate(e.agent)
			}
			visiting := fmt.Sprintf("%s <- %s", e.appName, e.endpointName)
			s.already_visited[visiting] += 1
			p := &statementsElement{
				endpointElement: *e,
				stmts:           e.endpoint.Stmt,
				deactivate:      active,
				lastParentStmt:  true,
			}
			s.visit(p)
			s.already_visited[visiting] -= 1
			if s.already_visited[visiting] == 0 {
				delete(s.already_visited, visiting)
			}
			if active {
				s.w.deactivate(e.agent)
			}
		}
	}
}

func (s *sdVisitor) normalizeSender(e *endpointElement) {
	if e.fromApp != nil {
		e.sender = s.pm.getOrCreateNewParticipant(getAppName(e.fromApp)).getAlias()
	} else {
		e.sender = "["
	}
}

func (s *sdVisitor) normalizeAgent(e *endpointElement) {
	e.agent = s.pm.getOrCreateNewParticipant(e.appName).getAlias()
}

func (s *sdVisitor) normalizeApplication(e *endpointElement) bool {
	app, ok := s.m.Apps[e.appName]
	e.application = app

	return ok
}

func (s *sdVisitor) normalizeEndpoint(e *endpointElement) bool {
	endpoint, ok := e.application.Endpoints[e.endpointName]
	e.endpoint = endpoint

	return ok
}

func (s *sdVisitor) normalizeAppPatterns(e *endpointElement) {
	e.appPatterns = makeSetSFromAttrs(e.application.Attrs)
}

func (s *sdVisitor) normalizeEndpointPatterns(e *endpointElement) {
	e.endpointPatterns = makeSetSFromAttrs(e.endpoint.Attrs)
}

func (s *sdVisitor) normalizePayload(e *endpointElement) {
	e.payload = getAndFmtReturnPayload(s.m, e.endpoint.Stmt)
}

func (s *sdVisitor) normalizeEndpointLabel(e *endpointElement) string {
	regEx := regexp.MustCompile(`^.*? -> `)
	label := regEx.ReplaceAllLiteralString(e.endpointName, " ⬄ ")

	if e.stmt != nil && e.stmt.GetCall() != nil {
		ptrns := ""
		if e.senderPatterns.count() > 0 || e.endpointPatterns.count() > 0 {
			ptrns = fmt.Sprintf("%s → %s",
				strings.Join(e.senderPatterns.getSortedElements(), ", "),
				strings.Join(e.endpointPatterns.getSortedElements(), ", "))
		}

		isoctrl := getSortedISOCtrls(e.endpoint.Attrs)

		epargs := getAndFmtParam(s.m, e.endpoint.Param)

		p := &epFmtParam{
			epname:   label,
			args:     strings.Join(epargs, " | "),
			patterns: ptrns,
			controls: strings.Join(isoctrl, ", "),
			attrs:    e.stmt.GetAttrs(),
		}

		human := e.appPatterns.contains("human")
		humanSender := e.senderPatterns.contains("human")
		cron := e.senderPatterns.contains("cron")
		needsInt := !(human || humanSender || cron) && e.sender != e.agent
		if human {
			p.human = "human"
		}
		if humanSender {
			p.human_sender = "human sender"
		}
		if needsInt {
			p.needs_int = "needs_int"
		}
		label = s.fmtEndpoint(p)
	}

	return label
}

func (s *sdVisitor) testCallingSelf(e *endpointElement) {
	callingSelf := e.fromApp != nil && getAppName(e.fromApp) == e.appName
	if !callingSelf && len(e.payload) == 0 && e.deactivate {
		s.w.deactivate(e.agent)
	}
}

func (s *sdVisitor) testCallingTimer(e *endpointElement) {
	human := e.appPatterns.contains("human")
	cron := e.appPatterns.contains("cron")

	if !((human && e.sender == "[") || cron) {
		icon := ""
		if e.endpointPatterns.contains("cron") {
			icon = "<&timer>"
		}
		label := s.normalizeEndpointLabel(e)
		fmt.Fprintf(s.w, "%s->%s : %s%s\n", e.sender, e.agent, icon, label)
	}
}

func (s *sdVisitor) testHitBlackbox(e *endpointElement) bool {
	visiting := fmt.Sprintf("%s <- %s", e.appName, e.endpointName)
	comment, hitUpto := s.uptos[visiting]
	_, hitVisited := s.already_visited[visiting]

	if hitUpto || hitVisited {
		if len(comment) == 0 {
			comment = "see below"
		}

		if len(e.payload) > 0 {
			s.w.activate(e.agent)
			s.w.indent()
			fmt.Fprintf(s.w, "note over %s: %s\n", e.agent, comment)
			fmt.Fprintf(s.w, "%s<--%s : %s\n", e.sender, e.agent, e.payload)
			s.w.unindent()
			s.w.deactivate(e.agent)
		} else {
			direct := "right"
			if e.sender > e.agent {
				direct = "left"
			}
			fmt.Fprintf(s.w, "note %s: %s\n", direct, comment)
		}
	}

	return hitUpto || hitVisited
}

func (s *sdVisitor) visitStmts(e *statementsElement) {
	for i, v := range e.stmts {
		switch c := v.Stmt.(type) {
		case *sysl.Statement_Call:
			s.visitCall(e, i, c.Call)
		case *sysl.Statement_Action:
			s.visitAction(e, i, c.Action)
		case *sysl.Statement_Cond:
			s.visitCond(e, i, c.Cond)
		case *sysl.Statement_Loop:
			s.visitLoop(e, i, c.Loop)
		case *sysl.Statement_LoopN:
			s.visitLoopN(e, i, c.LoopN)
		case *sysl.Statement_Foreach:
			s.visitForeach(e, i, c.Foreach)
		case *sysl.Statement_Group:
			s.visitGroup(e, i, c.Group)
		case *sysl.Statement_Alt:
			s.visitAlt(e, i, c.Alt)
		case *sysl.Statement_Ret:
			s.visitRet(e, i, c.Ret)
		default:
			logrus.Warn("No statement!")
		}
	}
}

func (s *sdVisitor) visitCall(e *statementsElement, i int, c *sysl.Call) {
	lastStmt := s.isLastStatement(e, i)
	stmtPatterns := makeSetSFromAttrs(e.stmts[i].GetAttrs())
	p := &endpointElement{
		fromApp:                e.application.GetName(),
		appName:                getAppName(c.GetTarget()),
		endpointName:           c.GetEndpoint(),
		senderPatterns:         e.appPatterns.copy(),
		senderEndpointPatterns: e.endpointPatterns.union(stmtPatterns),
		stmt:                   e.stmts[i],
		deactivate:             lastStmt && e.deactivate,
	}
	s.w.indent()
	s.visit(p)
	s.w.unindent()
}

func (s *sdVisitor) visitAction(e *statementsElement, i int, c *sysl.Action) {
	fmt.Fprintf(s.w, "%s -> %s : %s\n", e.agent, e.agent, c.GetAction())
}

func (s *sdVisitor) visitCond(e *statementsElement, i int, c *sysl.Cond) {
	s.visitGroupStmt(e, c.GetStmt(), s.isLastStatement(e, i), "opt %s\n", c.GetTest())
}

func (s *sdVisitor) visitLoop(e *statementsElement, i int, c *sysl.Loop) {
	s.visitGroupStmt(e, c.GetStmt(), s.isLastStatement(e, i), "loop %s %s\n", sysl.Loop_Mode_name[int32(c.GetMode())], c.GetCriterion())
}

func (s *sdVisitor) visitLoopN(e *statementsElement, i int, c *sysl.LoopN) {
	s.visitGroupStmt(e, c.GetStmt(), s.isLastStatement(e, i), "loop %d times\n", c.GetCount())
}

func (s *sdVisitor) visitForeach(e *statementsElement, i int, c *sysl.Foreach) {
	s.visitGroupStmt(e, c.GetStmt(), s.isLastStatement(e, i), "loop for each %s\n", c.GetCollection())
}

func (s *sdVisitor) visitGroup(e *statementsElement, i int, c *sysl.Group) {
	s.visitGroupStmt(e, c.GetStmt(), s.isLastStatement(e, i), "group %s\n", c.GetTitle())
}

func (s *sdVisitor) visitAlt(e *statementsElement, i int, c *sysl.Alt) {
	prefix := "alt"
	lastStmt := s.isLastStatement(e, i)
	for j, choice := range c.GetChoice() {
		lastAltStmt := lastStmt && j == len(c.GetChoice())-1
		s.visitBlockStmt(e, choice.GetStmt(), lastAltStmt, "%s %s", prefix, choice.GetCond())
		prefix = "else"
	}
	fmt.Fprintln(s.w, "end")
}

func (s *sdVisitor) visitRet(e *statementsElement, i int, c *sysl.Return) {
	rargs := formatReturnParam(s.m, c.GetPayload())
	fmt.Fprintf(s.w, "%s<--%s : %s\n", e.sender, e.agent, strings.Join(rargs, " | "))
}

func (s *sdVisitor) visitBlockStmt(e *statementsElement, stmts []*sysl.Statement, lastStmt bool, fmtStr string, args ...interface{}) {
	fmt.Fprintf(s.w, fmtStr, args...)

	s.w.indent()

	p := &statementsElement{
		endpointElement: e.endpointElement,
		deactivate:      e.deactivate,
		stmts:           stmts,
		lastParentStmt:  lastStmt,
	}
	s.visit(p)

	s.w.unindent()
}

func (s *sdVisitor) visitGroupStmt(e *statementsElement, stmts []*sysl.Statement, lastStmt bool, fmtStr string, args ...interface{}) {
	s.visitBlockStmt(e, stmts, lastStmt, fmtStr, args...)
	fmt.Fprintln(s.w, "end")
}

func (s *sdVisitor) isLastStatement(e *statementsElement, i int) bool {
	return e.lastParentStmt && i == len(e.stmts)-1
}
