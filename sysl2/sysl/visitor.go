package main

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
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
	uptos      StrSet
	blackboxes map[string]*Upto
}

func MakeEndpointCollectionElement(
	title string,
	endpoints []string,
	blackboxes map[string]*Upto,
) *EndpointCollectionElement {
	entries := make([]*entry, 0, len(endpoints))
	uptos := make([]string, 0, len(endpoints))

	for _, v := range endpoints {
		entry := makeEntry(v)
		entries = append(entries, entry)

		uptos = append(uptos, fmt.Sprintf("%s <- %s", entry.appName, entry.endpointName))
	}

	bb := make(map[string]*Upto)
	for k, b := range blackboxes {
		if len(b.Comment) > 0 {
			bb[k] = b
			if len(b.Comment) == 1 {
				b.Comment = ""
			}
		}
	}

	return &EndpointCollectionElement{
		title:      title,
		entries:    entries,
		uptos:      MakeStrSet(uptos...),
		blackboxes: bb,
	}
}

func (e *EndpointCollectionElement) Accept(v Visitor) error {
	return v.Visit(e)
}

type UptoType int

const (
	UpTo                 = 0
	BBApplication        = 1
	BBEndpointCollection = 2
	BBCommandLine        = 3
)

type Upto struct {
	VisitCount int
	Comment    string
	ValueType  UptoType
}

type EndpointElement struct {
	fromApp                *sysl.AppName
	appName                string
	endpointName           string
	uptos                  map[string]*Upto
	senderPatterns         StrSet
	senderEndpointPatterns StrSet
	stmt                   *sysl.Statement
	deactivate             func()
}

func (e *EndpointElement) Accept(v Visitor) error {
	return v.Visit(e)
}

func (e *EndpointElement) sender(v VarManager) string {
	if e.fromApp != nil {
		return v.UniqueVarForAppName(syslutil.GetAppName(e.fromApp))
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

func (e *EndpointElement) label(
	l EndpointLabeler,
	m *sysl.Module,
	ep *sysl.Endpoint,
	epp StrSet,
	isHuman, isHumanSender, needsInt bool,
) string {
	label := normalizeEndpointName(e.endpointName)

	if e.stmt != nil && e.stmt.GetCall() != nil {
		ptrns := func(a StrSet, b StrSet) string {
			if len(a) > 0 || len(b) > 0 {
				return fmt.Sprintf("%s → %s", strings.Join(a.ToSortedSlice(), ", "), strings.Join(b.ToSortedSlice(), ", "))
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

func (e *StatementElement) Accept(v Visitor) error {
	return v.Visit(e)
}

func (e *StatementElement) isLastStmt(i int) bool {
	return e.isLastParentStmt && i == len(e.stmts)-1
}

type SequenceDiagramVisitor struct {
	AppLabeler
	EndpointLabeler
	w          *SequenceDiagramWriter
	m          *sysl.Module
	visited    map[string]int
	symbols    map[string]*_var
	currentApp string
	groupby    string
	groupboxes map[string]StrSet
}

func MakeSequenceDiagramVisitor(
	a AppLabeler,
	e EndpointLabeler,
	w *SequenceDiagramWriter,
	m *sysl.Module,
	appName string,
	group string,
) *SequenceDiagramVisitor {
	return &SequenceDiagramVisitor{
		AppLabeler:      a,
		EndpointLabeler: e,
		w:               w,
		m:               m,
		visited:         make(map[string]int),
		symbols:         make(map[string]*_var),
		currentApp:      appName,
		groupby:         group,
		groupboxes:      map[string]StrSet{},
	}
}

func (v *SequenceDiagramVisitor) Visit(e Element) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	var err error
	switch t := e.(type) {
	case *EndpointCollectionElement:
		err = v.visitEndpointCollection(t)
		for bbKey, bbVal := range t.blackboxes {
			if bbVal.ValueType == BBEndpointCollection && bbVal.VisitCount == 0 {
				log.Warnf("blackbox '%s' not hit in app %s\n", bbKey, v.currentApp)
			}
		}
	case *EndpointElement:
		err = v.visitEndpoint(t)
	case *StatementElement:
		err = v.visitStatment(t)
	}
	return err
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

func (v *SequenceDiagramVisitor) visitEndpointCollection(e *EndpointCollectionElement) error {
	if len(e.title) > 0 {
		fmt.Fprintln(v.w, "title", e.title)
	}

	for _, entry := range e.entries {
		allUptos := MakeStrSet()
		for _, entry := range e.entries {
			item := fmt.Sprintf("%s <- %s", entry.appName, entry.endpointName)
			allUptos.Insert(item)
		}
		allUptos.Insert(entry.upto)

		fmt.Fprintf(v.w, "== %s <- %s ==\n", entry.appName, entry.endpointName)

		visiting := fmt.Sprintf("%s <- %s", entry.appName, entry.endpointName)
		delete(allUptos, visiting)
		for k := range allUptos {
			e.blackboxes[k] = &Upto{
				ValueType: UpTo,
				Comment:   "see below",
			}
		}
		for k := range v.visited {
			delete(v.visited, k)
		}
		e := &EndpointElement{
			appName:                entry.appName,
			endpointName:           entry.endpointName,
			uptos:                  e.blackboxes,
			senderPatterns:         MakeStrSet(),
			senderEndpointPatterns: MakeStrSet(),
		}

		if err := e.Accept(v); err != nil {
			return err
		}
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
		if _, err := v.w.WriteHead(item.String()); err != nil {
			return err
		}
	}
	return nil
}

func (v *SequenceDiagramVisitor) visitEndpoint(e *EndpointElement) error {
	sender := e.sender(v)
	agent := e.agent(v)
	app := e.application(v.m)
	endpoint := e.endpoint(app)

	appPatterns := MakeStrSetFromAttr("patterns", app.Attrs)
	endPointPatterns := MakeStrSetFromAttr("patterns", endpoint.Attrs)

	isHuman := appPatterns.Contains("human")
	isHumanSender := e.senderPatterns.Contains("human")
	isCron := appPatterns.Contains("cron")
	isCronSender := e.senderPatterns.Contains("cron")
	needsInt := !(isHuman || isHumanSender || isCronSender) && sender != agent

	if len(v.groupby) > 0 {
		if attr, exists := app.GetAttrs()[v.groupby]; exists {
			if _, has := v.groupboxes[attr.GetS()]; !has {
				v.groupboxes[attr.GetS()] = MakeStrSet()
			}
			v.groupboxes[attr.GetS()].Insert(e.appName)
		}
	}

	if !((isHuman && sender == "[") || isCron) {
		label := e.label(v, v.m, endpoint, endPointPatterns, isHuman, isHumanSender, needsInt)
		icon := func(a StrSet) string {
			if a.Contains("cron") {
				return "<&timer>"
			}
			return ""
		}(endPointPatterns)

		fmt.Fprintf(v.w, "%s->%s : %s%s\n", sender, agent, icon, label)
	}

	payload := strings.Join(formatReturnParam(v.m, getReturnPayload(endpoint.Stmt)), " | ")

	isCallingSelf := e.fromApp != nil && syslutil.GetAppName(e.fromApp) == e.appName

	if !isCallingSelf && len(payload) == 0 && e.deactivate != nil {
		e.deactivate()
	}

	if len(endpoint.Stmt) > 0 {
		visiting := fmt.Sprintf("%s <- %s", e.appName, e.endpointName)
		upto, hitUpto := e.uptos[visiting]
		if hitUpto {
			e.uptos[visiting].VisitCount++
		}
		_, hitVisited := v.visited[visiting]

		if hitUpto || hitVisited {
			if len(payload) > 0 {
				v.w.Activate(agent)
				if len(upto.Comment) > 0 {
					fmt.Fprintf(v.w, "note over %s: %s\n", agent, upto.Comment)
				}
			} else {
				direct := "right"
				if sender > agent {
					direct = "left"
				}
				fmt.Fprintf(v.w, "note %s: %s\n", direct, upto.Comment)
			}
			if len(payload) > 0 {
				fmt.Fprintf(v.w, "%s<--%s : %s\n", sender, agent, payload)
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
			if err := p.Accept(v); err != nil {
				return err
			}

			deactivate()
			v.visited[visiting]--
			if v.visited[visiting] == 0 {
				delete(v.visited, visiting)
			}
		}
	}
	return nil
}

func (v *SequenceDiagramVisitor) visitStatment(e *StatementElement) error {
	for i, s := range e.stmts {
		var err error
		switch c := s.Stmt.(type) {
		case *sysl.Statement_Call:
			err = v.visitCall(e, i, c.Call)
		case *sysl.Statement_Action:
			err = v.visitAction(e, c.Action)
		case *sysl.Statement_Cond:
			err = v.visitCond(e, i, c.Cond)
		case *sysl.Statement_Loop:
			err = v.visitLoop(e, i, c.Loop)
		case *sysl.Statement_LoopN:
			err = v.visitLoopN(e, i, c.LoopN)
		case *sysl.Statement_Foreach:
			err = v.visitForeach(e, i, c.Foreach)
		case *sysl.Statement_Group:
			err = v.visitGroup(e, i, c.Group)
		case *sysl.Statement_Alt:
			err = v.visitAlt(e, i, c.Alt)
		case *sysl.Statement_Ret:
			err = v.visitRet(e, c.Ret)
		default:
			panic("Unrecognised statement type")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *SequenceDiagramVisitor) visitCall(e *StatementElement, i int, c *sysl.Call) error {
	isLastStmt := e.isLastStmt(i)
	app := e.application(v.m)
	stmtPatterns := MakeStrSetFromAttr("patterns", e.stmts[i].Attrs)
	senderPatterns := MakeStrSetFromAttr("patterns", app.Attrs)
	endpointPatterns := MakeStrSetFromAttr("patterns", e.endpoint(app).Attrs)

	p := &EndpointElement{
		fromApp:                app.GetName(),
		appName:                syslutil.GetAppName(c.GetTarget()),
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
	if err := p.Accept(v); err != nil {
		return err
	}
	v.w.Unindent()
	return nil
}

func (v *SequenceDiagramVisitor) visitAction(e *StatementElement, c *sysl.Action) error {
	_, err := fmt.Fprintf(v.w, "%s -> %s : %s\n", e.agent(v), e.agent(v), c.GetAction())
	return err
}

func (v *SequenceDiagramVisitor) visitCond(e *StatementElement, i int, c *sysl.Cond) error {
	return v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "opt %s\n", c.GetTest())
}

func (v *SequenceDiagramVisitor) visitLoop(e *StatementElement, i int, c *sysl.Loop) error {
	return v.visitGroupStmt(
		e, c.GetStmt(), e.isLastStmt(i), "loop %s %s\n",
		sysl.Loop_Mode_name[int32(c.GetMode())], c.GetCriterion(),
	)
}

func (v *SequenceDiagramVisitor) visitLoopN(e *StatementElement, i int, c *sysl.LoopN) error {
	return v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "loop %d times\n", c.GetCount())
}

func (v *SequenceDiagramVisitor) visitForeach(e *StatementElement, i int, c *sysl.Foreach) error {
	return v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "loop for each %s\n", c.GetCollection())
}

func (v *SequenceDiagramVisitor) visitGroup(e *StatementElement, i int, c *sysl.Group) error {
	return v.visitGroupStmt(e, c.GetStmt(), e.isLastStmt(i), "group %s\n", c.GetTitle())
}

func (v *SequenceDiagramVisitor) visitAlt(e *StatementElement, i int, c *sysl.Alt) error {
	prefix := "alt"
	lastStmt := e.isLastStmt(i)
	for j, choice := range c.GetChoice() {
		lastAltStmt := lastStmt && j == len(c.GetChoice())-1
		if err := v.visitBlockStmt(e, choice.GetStmt(), lastAltStmt, "%s %s\n", prefix, choice.GetCond()); err != nil {
			return err
		}
		prefix = "else"
	}
	_, err := fmt.Fprintln(v.w, "end")
	return err
}

func (v *SequenceDiagramVisitor) visitRet(e *StatementElement, c *sysl.Return) error {
	rargs := formatReturnParam(v.m, c.GetPayload())
	_, err := fmt.Fprintf(v.w, "%s<--%s : %s\n", e.sender(v), e.agent(v), strings.Join(rargs, " | "))
	return err
}

func (v *SequenceDiagramVisitor) visitBlockStmt(
	e *StatementElement,
	stmts []*sysl.Statement,
	isLastStmt bool,
	fmtStr string,
	args ...interface{},
) error {
	fmt.Fprintf(v.w, fmtStr, args...)
	v.w.Indent()
	p := &StatementElement{
		EndpointElement:  e.EndpointElement,
		deactivate:       e.deactivate,
		stmts:            stmts,
		isLastParentStmt: isLastStmt,
	}
	if err := v.visitStatment(p); err != nil {
		return err
	}
	v.w.Unindent()
	return nil
}

func (v *SequenceDiagramVisitor) visitGroupStmt(
	e *StatementElement,
	stmts []*sysl.Statement,
	isLastStmt bool,
	fmtStr string,
	args ...interface{},
) error {
	if err := v.visitBlockStmt(e, stmts, isLastStmt, fmtStr, args...); err != nil {
		return err
	}
	_, err := fmt.Fprintln(v.w, "end")
	return err
}

type agent struct {
	category int
	name     string
}

//nolint:gochecknoglobals
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
