package seqs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	log "github.com/sirupsen/logrus"
)

type entry struct {
	appName      string
	endpointName string
	upto         string
}

func makeEntry(s string) *entry {
	match := endpointParserRE.FindStringSubmatch(s)

	out := &entry{}
	for i, name := range endpointParserRE.SubexpNames() {
		if i > 0 && i <= len(match) {
			switch name {
			case "appname":
				out.appName = match[i]
			case "epname":
				out.endpointName = match[i]
			case "upto":
				out.upto = match[i]
			}
		}
	}
	return out
}

type EndpointCollectionElement struct {
	title      string
	entries    []*entry
	uptos      *strSet
	blackboxes map[string]string
}

func MakeEndpointCollectionElement(title string, endpoints []string, blackboxes [][]string) *EndpointCollectionElement {
	entries := make([]*entry, 0, len(endpoints))
	uptos := make([]string, 0, len(endpoints))

	for _, v := range endpoints {
		entry := makeEntry(v)
		entries = append(entries, entry)

		uptos = append(uptos, fmt.Sprintf("%s <- %s", entry.appName, entry.endpointName))
	}

	bb := make(map[string]string)
	for _, b := range blackboxes {
		switch len(b) {
		case 0:
			continue
		case 1:
			bb[b[0]] = ""
		default:
			bb[b[0]] = b[1]
		}
	}

	return &EndpointCollectionElement{
		title:      title,
		entries:    entries,
		uptos:      makeStrSet(uptos...),
		blackboxes: bb,
	}
}

func (e *EndpointCollectionElement) Accept(v Visitor) {
	v.Visit(e)
}

type EndpointElement struct {
	fromApp                *sysl.AppName
	appName                string
	endpointName           string
	uptos                  map[string]string
	senderPatterns         *strSet
	senderEndpointPatterns *strSet
	stmt                   *sysl.Statement
	deactivate             func()
}

func (e *EndpointElement) Accept(v Visitor) {
	v.Visit(e)
}

func (e *EndpointElement) sender(v VarManager) string {
	if e.fromApp != nil {
		return v.UniqueVarForAppName(getAppName(e.fromApp))
	}

	return "["
}

func (e *EndpointElement) agent(v VarManager) string {
	return v.UniqueVarForAppName(e.appName)
}

func (e *EndpointElement) application(m *sysl.Module) *sysl.Application {
	if app, ok := m.Apps[e.appName]; ok {
		return app
	}
	panic(fmt.Sprintf("The application with name %s does not exists", e.appName))
}

func (e *EndpointElement) endpoint(a *sysl.Application) *sysl.Endpoint {
	if ep, ok := a.Endpoints[e.endpointName]; ok {
		return ep
	}
	panic(fmt.Sprintf("The endpoint with name %s does not exists in the Application with name %s",
		e.endpointName, e.appName))
}

func (e *EndpointElement) label(l EndpointLabeler, m *sysl.Module, ep *sysl.Endpoint, epp *strSet,
	isHuman, isHumanSender, needsInt bool) string {
	label := normalizeEndpointName(e.endpointName)

	if e.stmt != nil && e.stmt.GetCall() != nil {
		ptrns := func(a *strSet, b *strSet) string {
			if a.Len() > 0 || b.Len() > 0 {
				return fmt.Sprintf("%s -> %s", strings.Join(a.ToSlice(), ", "), strings.Join(b.ToSlice(), ", "))
			}
			return ""
		}(e.senderEndpointPatterns, epp)

		isoctrl := getSortedISOCtrlStr(ep.Attrs)
		epargs := getAndFmtParam(m, ep.Param)

		str := func(t bool, v string) string {
			if t {
				return v
			}
			return ""
		}

		param := &EndpointLabelerParam{
			EndpointName: label,
			Human:        str(isHuman, "human"),
			HumanSender:  str(isHumanSender, "human sender"),
			NeedsInt:     str(needsInt, "needs_int"),
			Args:         strings.Join(epargs, " | "),
			Patterns:     ptrns,
			Controls:     isoctrl,
			Attrs:        e.stmt.Attrs,
		}
		label = l.LabelEndpoint(param)
	}

	return label
}

type StatementElement struct {
	EndpointElement
	stmts            []*sysl.Statement
	deactivate       func()
	isLastParentStmt bool
}

func (e *StatementElement) Accept(v Visitor) {
	v.Visit(e)
}

func (e *StatementElement) isLastStmt(i int) bool {
	return e.isLastParentStmt && i == len(e.stmts)-1
}

type SequenceDiagramVisitor struct {
	AppLabeler
	EndpointLabeler
	w       *SequenceDiagramWriter
	m       *sysl.Module
	visited map[string]int
	symbols map[string]*_var
}

func MakeSequenceDiagramVisitor(a AppLabeler, e EndpointLabeler, w *SequenceDiagramWriter, m *sysl.Module) *SequenceDiagramVisitor {
	return &SequenceDiagramVisitor{
		AppLabeler: a,
		EndpointLabeler: e,
		w:       w,
		m:       m,
		visited: make(map[string]int),
		symbols: make(map[string]*_var),
	}
}

func (v *SequenceDiagramVisitor) Visit(e Element) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
	}()

	switch t := e.(type) {
	case *EndpointCollectionElement:
		v.visitEndpointCollection(t)
	case *EndpointElement:
		v.visitEndpoint(t)
	case *StatementElement:
		v.visitStatment(t)
	}

	s := make([]*_var, 0, len(v.symbols))
	for _, item := range v.symbols {
		s = append(s, item)
	}
	sort.Slice(s, func(i, j int) bool {
		if s[i].category == s[j].category {
			return s[i].order < s[j].order
		}
		return s[i].category < s[j].category
	})
	for _, item := range s {
		v.w.WriteHead(item.String())
	}
}

func (v *SequenceDiagramVisitor) UniqueVarForAppName(appName string) string {
	if s, ok := v.symbols[appName]; ok {
		return s.alias
	}

	i := len(v.symbols)
	alias := fmt.Sprintf("_%d", i)
	attrs := getApplicationAttrs(v.m, appName)
	controls := getSortedISOCtrlStr(attrs)
	label := v.LabelApp(appName, controls, attrs)
	s := &_var{
		agent: makeAgent(attrs),
		order: i,
		label: label,
		alias: alias,
	}
	v.symbols[appName] = s

	return s.alias
}

func (v *SequenceDiagramVisitor) visitEndpointCollection(e *EndpointCollectionElement) {
	if len(e.title) > 0 {
		fmt.Fprintln(v.w, "title", e.title)
	}

	for _, entry := range e.entries {
		fmt.Fprintf(v.w, "== %s <- %s ==\n", entry.appName, entry.endpointName)

		visiting := fmt.Sprintf("%s <- %s", entry.appName, entry.endpointName)
		bbs := make(map[string]string)
		for k, v := range e.blackboxes {
			if k == entry.upto || (e.uptos.Contains(k) && k != visiting) {
				bbs[k] = "see below"
			} else {
				bbs[k] = v
			}
		}
		for k := range v.visited {
			delete(v.visited, k)
		}
		e := &EndpointElement{
			appName:      entry.appName,
			endpointName: entry.endpointName,
			uptos:        bbs,
			senderPatterns: makeStrSet(),
			senderEndpointPatterns: makeStrSet(),
		}

		v.visitEndpoint(e)
	}
}

func (v *SequenceDiagramVisitor) visitEndpoint(e *EndpointElement) {
	sender := e.sender(v)
	agent := e.agent(v)
	app := e.application(v.m)
	endpoint := e.endpoint(app)

	appPatterns := makeStrSetFromPatternsAttr(app.Attrs)
	endPointPatterns := makeStrSetFromPatternsAttr(endpoint.Attrs)

	isHuman := appPatterns.Contains("human")
	isHumanSender := e.senderPatterns.Contains("human")
	isCron := appPatterns.Contains("cron")
	isCronSender := e.senderPatterns.Contains("cron")
	needsInt := !(isHuman || isHumanSender || isCronSender) && sender != agent

	if !((isHuman && sender == "[") || isCron) {
		label := e.label(v, v.m, endpoint, endPointPatterns, isHuman, isHumanSender, needsInt)
		icon := func(a *strSet) string {
			if a.Contains("cron") {
				return "<&timer>"
			}
			return ""
		}(endPointPatterns)

		fmt.Fprintf(v.w, "%s->%s : %s%s\n", sender, agent, icon, label)
	}

	payload := strings.Join(formatReturnParam(v.m, getReturnPayload(v.m, endpoint.Stmt)), " | ")

	isCallingSelf := e.fromApp != nil && getAppName(e.fromApp) == e.appName

	if !isCallingSelf && len(payload) == 0 && e.deactivate != nil {
		e.deactivate()
	}

	if len(endpoint.Stmt) > 0 {
		visiting := fmt.Sprintf(" %s <- %s\n", e.appName, e.endpointName)
		comment, hitUpto := e.uptos[visiting]
		_, hitVisited := v.visited[visiting]

		if hitUpto || hitVisited {
			if len(payload) > 0 {
				v.w.Activate(agent)
				if len(comment) > 0 {
					fmt.Fprintf(v.w, "note over %s: %s\n", agent, comment)
				}
			} else {
				direct := "right"
				if sender > agent {
					direct = "left"
				}
				fmt.Fprintf(v.w, "note %s: %s\n", direct, comment)
			}
			if len(payload) > 0 {
				fmt.Fprintf(v.w, "%s<--%s : %s", sender, agent, payload)
				v.w.Deactivate(agent)
			}
		} else {
			deactivate := v.w.Activated(agent, isHuman || isCron)
			v.visited[visiting]++

			p := &StatementElement{
				EndpointElement:  *e,
				stmts:            endpoint.Stmt,
				deactivate:       deactivate,
				isLastParentStmt: true,
			}
			v.visitStatment(p)

			deactivate()
			v.visited[visiting]--
			if v.visited[visiting] == 0 {
				delete(v.visited, visiting)
			}
		}
	}
}

func (v *SequenceDiagramVisitor) visitStatment(e *StatementElement) {
	for i, s := range e.stmts {
		switch c := s.Stmt.(type) {
		case *sysl.Statement_Call:
			v.visitCall(e, i, c.Call)
		case *sysl.Statement_Action:
			v.visitAction(e, c.Action)
		case *sysl.Statement_Cond:
			v.visitCond(e, i, c.Cond)
		case *sysl.Statement_Loop:
			v.visitLoop(e, i, c.Loop)
		case *sysl.Statement_LoopN:
			v.visitLoopN(e, i, c.LoopN)
		case *sysl.Statement_Foreach:
			v.visitForeach(e, i, c.Foreach)
		case *sysl.Statement_Group:
			v.visitGroup(e, i, c.Group)
		case *sysl.Statement_Alt:
			v.visitAlt(e, i, c.Alt)
		case *sysl.Statement_Ret:
			v.visitRet(e, c.Ret)
		default:
			panic("No statement!")
		}
	}
}

func (v *SequenceDiagramVisitor) visitCall(e *StatementElement, i int, c *sysl.Call) {
	isLastStmt := e.isLastStmt(i)
	app := e.application(v.m)
	stmtPatterns := makeStrSetFromPatternsAttr(e.stmts[i].Attrs)
	senderPatterns := makeStrSetFromPatternsAttr(app.Attrs)
	endpointPatterns := makeStrSetFromPatternsAttr(e.endpoint(app).Attrs)

	p := &EndpointElement{
		fromApp:                app.GetName(),
		appName:                getAppName(c.GetTarget()),
		endpointName:           c.GetEndpoint(),
		uptos:                  e.uptos,
		senderPatterns:         senderPatterns,
		senderEndpointPatterns: endpointPatterns.Union(stmtPatterns),
		stmt:                   e.stmts[i],
		deactivate: func() {
			if isLastStmt {
				e.deactivate()
			}
		},
	}
	v.w.Indent()
	v.visitEndpoint(p)
	v.w.Unindent()
}

func (v *SequenceDiagramVisitor) visitAction(e *StatementElement, c *sysl.Action) {
	fmt.Fprintf(v.w, "%s -> %s : %s\n", e.agent(v), e.agent(v), c.GetAction())
}

func (v *SequenceDiagramVisitor) visitCond(e *StatementElement, i int, c *sysl.Cond) {
	v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "opt %s\n", c.GetTest())
}

func (v *SequenceDiagramVisitor) visitLoop(e *StatementElement, i int, c *sysl.Loop) {
	v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "loop %s %s\n",
		sysl.Loop_Mode_name[int32(c.GetMode())], c.GetCriterion())
}

func (v *SequenceDiagramVisitor) visitLoopN(e *StatementElement, i int, c *sysl.LoopN) {
	v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "loop %d times\n", c.GetCount())
}

func (v *SequenceDiagramVisitor) visitForeach(e *StatementElement, i int, c *sysl.Foreach) {
	v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "loop for each %s\n", c.GetCollection())
}

func (v *SequenceDiagramVisitor) visitGroup(e *StatementElement, i int, c *sysl.Group) {
	v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "group %s\n", c.GetTitle())
}

func (v *SequenceDiagramVisitor) visitAlt(e *StatementElement, i int, c *sysl.Alt) {
	prefix := "alt"
	lastStmt := e.isLastStmt(i)
	for j, choice := range c.GetChoice() {
		lastAltStmt := lastStmt && j == len(c.GetChoice())-1
		v.visitBlockStmt(e, choice.GetStmt(), lastAltStmt, "%s %s", prefix, choice.GetCond())
		prefix = "else"
	}
	fmt.Fprintln(v.w, "end")
}

func (v *SequenceDiagramVisitor) visitRet(e *StatementElement, c *sysl.Return) {
	rargs := formatReturnParam(v.m, c.GetPayload())
	fmt.Fprintf(v.w, "%s<--%s : %s\n", e.sender(v), e.agent(v), strings.Join(rargs, " | "))
}

func (v *SequenceDiagramVisitor) visitBlockStmt(e *StatementElement,
	stmts []*sysl.Statement, isLastStmt bool, fmtStr string, args ...interface{}) {
	fmt.Fprintf(v.w, fmtStr, args...)
	v.w.Indent()
	p := &StatementElement{
		EndpointElement:  e.EndpointElement,
		deactivate:       e.deactivate,
		stmts:            stmts,
		isLastParentStmt: isLastStmt,
	}
	v.visitStatment(p)
	v.w.Unindent()
}

func (v *SequenceDiagramVisitor) visitGroupStmt(e *StatementElement,
	stmts []*sysl.Statement, isLastStmt bool, fmtStr string, args ...interface{}) {
	v.visitBlockStmt(e, stmts, isLastStmt, fmtStr, args...)
	fmt.Fprintln(v.w, "end")
}

type agent struct {
	category int
	name     string
}

var agents = map[string]agent{
	"human":    {0, "actor"},
	"ui":       {1, "boundary"},
	"cron":     {2, "control"},
	"db":       {4, "database"},
	"external": {5, "control"},
}

func makeAgent(attrs map[string]*sysl.Attribute) agent {
	if patterns, ok := attrs["patterns"]; ok {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				if v, ok := agents[y.GetS()]; ok {
					return v
				}
			}
		}
	}

	return agent{3, "control"}
}

type _var struct {
	agent
	order int
	label string
	alias string
}

func (s _var) String() string {
	return fmt.Sprintf(`%s "%s" as %s`, s.name, s.label, s.alias)
}
