package main

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/anz-bank/sysl/src/proto"
)

type (
	nothing struct{}

	setS struct {
		m map[string]nothing
	}
)

func makeSetS(initial ...string) *setS {
	s := &setS{make(map[string]nothing)}

	for _, v := range initial {
		if len(v) > 0 {
			s.insert(v)
		}
	}

	return s
}

func makeSetSFromAttrs(attrs map[string]*sysl.Attribute) *setS {
	s := &setS{make(map[string]nothing)}

	if patterns, has := attrs["patterns"]; has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				s.insert(y.GetS())
			}
		}
	}

	return s
}

func (s *setS) contains(elem string) bool {
	_, ok := s.m[elem]
	return ok
}

func (s *setS) insert(elem string) {
	s.m[elem] = nothing{}
}

func (s *setS) remove(elem string) {
	delete(s.m, elem)
}

func (s *setS) count() int {
	return len(s.m)
}

func (s *setS) union(other *setS) *setS {
	m := make(map[string]nothing)

	for k := range s.m {
		m[k] = nothing{}
	}

	for k := range other.m {
		m[k] = nothing{}
	}

	return &setS{m}
}

func (s *setS) intersection(other *setS) *setS {
	m := make(map[string]nothing)

	for k := range s.m {
		if _, ok := other.m[k]; ok {
			m[k] = nothing{}
		}
	}

	return &setS{m}
}

func (s *setS) difference(other *setS) *setS {
	m := make(map[string]nothing)

	for k := range s.m {
		if _, ok := other.m[k]; !ok {
			m[k] = nothing{}
		}
	}

	return &setS{m}
}

func (s *setS) getSortedElements() []string {
	o := make([]string, 0, len(s.m))

	for k := range s.m {
		o = append(o, k)
	}

	sort.Strings(o)

	return o
}

func (s *setS) copy() *setS {
	m := make(map[string]nothing)

	for k := range s.m {
		m[k] = nothing{}
	}

	return &setS{m}
}

type (
	agent struct {
		category int
		name     string
	}

	participant struct {
		agent
		order int
		label string
		alias string
	}

	onNewParticipantCreating func(appname, alias string, order int) *participant

	participantManager struct {
		sync.Mutex
		symbols                      map[string]*participant
		handleNewParticipantCreating onNewParticipantCreating
	}
)

var (
	agents = map[string]agent{
		"human":    {0, "actor"},
		"ui":       {1, "boundary"},
		"cron":     {2, "control"},
		"db":       {4, "database"},
		"external": {5, "control"},
	}
)

func makeAgent(attrs map[string]*sysl.Attribute) agent {
	if patterns, has := attrs["patterns"]; has {
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

func (p *participant) getAlias() string {
	return p.alias
}

func (p participant) String() string {
	return fmt.Sprintf(`%s "%s" as %s`, p.name, p.label, p.alias)
}

func makeParticipantManager(handler onNewParticipantCreating) *participantManager {
	if handler == nil {
		panic("handleNewParticipantCreating can not be nil")
	}

	return &participantManager{
		symbols:                      make(map[string]*participant),
		handleNewParticipantCreating: handler,
	}
}

func (pm *participantManager) getOrCreateNewParticipant(appname string) *participant {
	pm.Lock()
	defer pm.Unlock()
	if v, ok := pm.symbols[appname]; ok {
		return v
	}

	i := len(pm.symbols)
	v := pm.handleNewParticipantCreating(appname, fmt.Sprintf("_%d", i), i)
	if v != nil {
		pm.symbols[appname] = v
	}

	return v
}

func (pm *participantManager) getSortedParticipants() []*participant {
	s := make([]*participant, 0, len(pm.symbols))

	pm.Lock()
	for _, v := range pm.symbols {
		s = append(s, v)
	}
	pm.Unlock()

	sort.Slice(s, func(i, j int) bool {
		if s[i].category != s[j].category {
			return s[i].category < s[j].category
		}
		return s[i].order < s[j].order
	})

	return s
}

type (
	textWriter struct {
		ind            int
		complete       bool
		autogenWarning bool
		body           bytes.Buffer
		head           bytes.Buffer
		active         map[string]int
	}
)

func (w *textWriter) WriteString(s string) (n int, err error) {
	if !strings.Contains(s, "\n") {
		if w.complete {
			w.writeIndent()
		}
		w.complete = false
		return w.body.WriteString(s)
	}
	return w.Write([]byte(s))
}

func (w *textWriter) Write(p []byte) (n int, err error) {
	newline := []byte("\n")
	newlines := bytes.Count(p, newline)
	if newlines == 0 {
		if w.complete {
			w.writeIndent()
		}
		n, err = w.body.Write(p)
		w.complete = false
		return n, err
	}

	frags := bytes.SplitN(p, newline, newlines+1)

	for i, frag := range frags {
		if w.complete {
			w.writeIndent()
		}
		nn, err := w.body.Write(frag)
		n += nn
		if err != nil {
			return n, err
		}
		if i+1 < len(frags) {
			if err := w.body.WriteByte('\n'); err != nil {
				return n, err
			}
			n++
		}
	}
	w.complete = len(frags[len(frags)-1]) == 0
	return n, nil
}

func (w *textWriter) WriteByte(c byte) error {
	if w.complete {
		w.writeIndent()
	}
	err := w.body.WriteByte(c)
	w.complete = c == '\n'
	return err
}

func (w *textWriter) writeIndent() {
	if !w.complete {
		return
	}

	spaces := make([]byte, 0, w.ind)
	for i := 0; i < w.ind; i++ {
		spaces = append(spaces, ' ')
	}
	w.body.Write(spaces)
	w.complete = false
}

func (w *textWriter) writeHead(s string) {
	w.head.WriteString(s)
	if s[len(s)-1] != '\n' {
		w.head.WriteByte('\n')
	}
}

func (w *textWriter) indent() {
	w.ind++
}

func (w *textWriter) unindent() {
	if w.ind == 0 {
		return
	}
	w.ind--
}

func (w *textWriter) activate(s string) {
	w.active[s]++
	fmt.Fprintf(w, "activate %s\n", s)
}

func (w *textWriter) deactivate(s string) {
	if _, ok := w.active[s]; !ok {
		return
	}

	w.active[s]--
	fmt.Fprintf(w, "deactivate %s\n", s)
	if w.active[s] == 0 {
		delete(w.active, s)
	}
}

func (w textWriter) String() string {
	if w.body.Len() == 0 || w.head.Len() == 0 {
		return ""
	}

	var sb strings.Builder
	if w.autogenWarning {
		fmt.Fprintln(&sb, "''''''''''''''''''''''''''''''''''''''''''")
		fmt.Fprintln(&sb, "''                                      ''")
		fmt.Fprintln(&sb, "''  AUTOGENERATED CODE -- DO NOT EDIT!  ''")
		fmt.Fprintln(&sb, "''                                      ''")
		fmt.Fprintln(&sb, "''''''''''''''''''''''''''''''''''''''''''")
		fmt.Fprintln(&sb)
	}
	fmt.Fprintln(&sb, "@startuml")
	sb.WriteString(w.head.String())
	sb.WriteString(w.body.String())
	fmt.Fprintln(&sb, "@enduml")

	return sb.String()
}

func getSortedISOCtrls(attrs map[string]*sysl.Attribute) []string {
	s := make([]string, 0)

	reg := regexp.MustCompile("iso_ctrl_(.*)_txt")
	for k := range attrs {
		if !strings.Contains(k, "iso_ctrl") {
			continue
		}
		match := reg.FindStringSubmatch(k)
		if len(match) > 1 {
			s = append(s, match[1])
		}
	}
	sort.Strings(s)
	return s
}

func formatArgs(s *sysl.Module, appName, parameterTypeName string) string {
	if len(appName) == 0 || len(parameterTypeName) == 0 {
		return ""
	}

	conf := "?"
	integ := "?"
	if app, ok := s.Apps[appName]; ok {
		if t, exist := app.Types[parameterTypeName]; exist {
			for k, v := range t.Attrs {
				if s := v.GetS(); len(s) > 0 {
					switch k {
					case "iso_conf":
						conf = strings.ToUpper(s[:1])
					case "iso_integ":
						integ = strings.ToUpper(s[:1])
					default:
						continue
					}
				}
			}
		}
	}

	isocolor := "green"
	if conf == "R" {
		isocolor = "red"
	}

	return fmt.Sprintf("<color blue>%s.%s</color> <<color %s>%s, %s</color>>", appName, parameterTypeName, isocolor, conf, integ)
}

func formatReturnParam(s *sysl.Module, payload string) []string {
	ptns := make([]string, 0)
	if len(payload) > 0 {
		re := regexp.MustCompile(`,?(![^{]*\})`)
		rns := re.Split(payload, -1)
		for _, rn := range rns {
			ptn := rn
			if strings.Count(rn, "<:") == 1 {
				rex := regexp.MustCompile(`\s*<:\s*`)
				ps := rex.Split(rn, -1)
				if len(ps) == 2 {
					ptn = ps[1]
				}
			}

			if _, ok := sysl.Type_Primitive_value[strings.ToUpper(ptn)]; !ok {
				rex := regexp.MustCompile(`set\s+of\s+(.+)$`)
				if m := rex.FindStringSubmatch(ptn); len(m) > 0 {
					ptn = m[1]
				}
				rex = regexp.MustCompile(`one\s+of\s*{(.+)}$`)
				if m := rex.FindStringSubmatch(ptn); len(m) > 0 {
					rex = regexp.MustCompile(`\s*,\s*`)
					ptns = append(ptns, rex.Split(m[1], -1)...)
				} else {
					ptns = append(ptns, ptn)
				}
			}
		}
	}

	rargs := make([]string, 0, len(ptns))
	for _, ptn := range ptns {
		if !strings.Contains(ptn, "...") && strings.Contains(ptn, ".") {
			aps := strings.Split(ptn, ".")
			if len(aps) > 1 {
				rarg := formatArgs(s, aps[0], aps[1])
				if len(rarg) > 0 {
					rargs = append(rargs, rarg)
				}
			}
		} else {
			rargs = append(rargs, ptn)
		}
	}
	return rargs
}

func getStatementsReturnPayload(s *sysl.Module, stmts []*sysl.Statement) string {
	for _, v := range stmts {
		switch stmt := v.Stmt.(type) {
		case *sysl.Statement_Call, *sysl.Statement_Action:
			continue
		case *sysl.Statement_Cond:
			p := getStatementsReturnPayload(s, stmt.Cond.Stmt)
			if len(p) > 0 {
				return p
			}
		case *sysl.Statement_Loop:
			p := getStatementsReturnPayload(s, stmt.Loop.Stmt)
			if len(p) > 0 {
				return p
			}
		case *sysl.Statement_LoopN:
			p := getStatementsReturnPayload(s, stmt.LoopN.Stmt)
			if len(p) > 0 {
				return p
			}
		case *sysl.Statement_Foreach:
			p := getStatementsReturnPayload(s, stmt.Foreach.Stmt)
			if len(p) > 0 {
				return p
			}
		case *sysl.Statement_Group:
			p := getStatementsReturnPayload(s, stmt.Group.Stmt)
			if len(p) > 0 {
				return p
			}
		case *sysl.Statement_Alt:
			for _, c := range stmt.Alt.Choice {
				p := getStatementsReturnPayload(s, c.Stmt)
				if len(p) > 0 {
					return p
				}
			}
		case *sysl.Statement_Ret:
			return stmt.Ret.GetPayload()
		}
	}
	return ""
}

func getAndFmtReturnPayload(s *sysl.Module, stmts []*sysl.Statement) string {
	p := getStatementsReturnPayload(s, stmts)
	pf := formatReturnParam(s, p)

	return strings.Join(pf, " | ")
}

func getAndFmtParam(s *sysl.Module, params []*sysl.Param) []string {
	r := make([]string, 0, len(params))
	for _, v := range params {
		if ref_type := v.GetType().GetTypeRef(); ref_type != nil {
			if ref := ref_type.GetRef(); ref != nil {
				an := getAppName(ref.GetAppname())
				pn := strings.Join(ref.GetPath(), ".")
				eparg := formatArgs(s, an, pn)
				if len(eparg) > 0 {
					r = append(r, eparg)
				}
			}
		}
	}
	return r
}
