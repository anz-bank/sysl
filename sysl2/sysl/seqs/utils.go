package seqs

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
)

var (
	isoCtrlRE                = regexp.MustCompile("^iso_ctrl_(.*)_txt$")
	returnParamSpliterRE     = regexp.MustCompile(`,?(![^{]*\})`)
	returnTypeValueSpliterRE = regexp.MustCompile(`\s*<:\s*`)
	typeSetOfRE              = regexp.MustCompile(`set\s+of\s+(.+)$`)
	typeOneOfRE              = regexp.MustCompile(`one\s+of\s*{(.+)}$`)
	typeListSpliterRE        = regexp.MustCompile(`\s*,\s*`)
	endpointLabelReplaceRE   = regexp.MustCompile(`^.*? -> `)
	endpointParserRE         = regexp.MustCompile(`(?P<appname>.*?)\s*<-\s*(?P<epname>.*?)(?:\s*\[upto\s+(?P<upto>.*)\])*$`)
)

func TransformBlackBoxes(blackboxes []*sysl.Attribute) [][]string {
	bbs := [][]string{}
	for _, vals := range blackboxes {
		sub_bbs := []string{}
		for _, val := range vals.GetA().Elt {
			sub_bbs = append(sub_bbs, val.GetS())
		}
		bbs = append(bbs, sub_bbs)
	}

	return bbs
}

func ParseBlackBoxesFromArgument(blackboxFlags []string) [][]string {
	bbs := [][]string{}
	for _, blackboxFlag := range blackboxFlags {
		sub_bbs := []string{}
		sub_bbs = append(sub_bbs, strings.Split(blackboxFlag, ",")...)
		bbs = append(bbs, sub_bbs)
	}

	return bbs
}

func FindMatchItems(origin string) []string {
	re := regexp.MustCompile(`(%\(\w+\))`)
	return re.FindAllString(origin, -1)
}

func RemoveWrapper(origin string) string {
	replaced := strings.Replace(origin, "%(", "", 1)
	replaced = strings.Replace(replaced, ")", "", 1)
	return replaced
}

func RemovePercentSymbol(origin string) string {
	return strings.Replace(origin, "%", "", -1)
}

func MergeAttributes(app, edpnt map[string]*sysl.Attribute) map[string]*sysl.Attribute {
	result := make(map[string]*sysl.Attribute)
	for k, v := range app {
		result[k] = v
	}
	for k, v := range edpnt {
		result[k] = v
	}

	return result
}

func getAppName(appName *sysl.AppName) string {
	return strings.Join(appName.Part, " :: ")
}

func getApplicationAttrs(m *sysl.Module, appName string) map[string]*sysl.Attribute {
	if app, ok := m.Apps[appName]; ok {
		return app.Attrs
	}
	return nil
}

func getSortedISOCtrlSlice(attrs map[string]*sysl.Attribute) []string {
	s := make([]string, 0)

	for k := range attrs {
		match := isoCtrlRE.FindStringSubmatch(k)
		if len(match) > 1 {
			s = append(s, match[1])
		}
	}
	sort.Strings(s)
	return s
}

func getSortedISOCtrlStr(attrs map[string]*sysl.Attribute) string {
	return strings.Join(getSortedISOCtrlSlice(attrs), ", ")
}

func formatArgs(s *sysl.Module, appName, parameterTypeName string) string {
	if len(appName) == 0 || len(parameterTypeName) == 0 {
		return ""
	}

	val := func(a *sysl.Attribute) string {
		if s := a.GetS(); len(s) > 0 {
			return strings.ToUpper(s[:1])
		}
		return "?"
	}

	conf := ""
	integ := ""
	if app, ok := s.Apps[appName]; ok {
		if t, ok := app.Types[parameterTypeName]; ok {
			if v, ok := t.Attrs["iso_conf"]; ok {
				conf = val(v)
			}
			if v, ok := t.Attrs["iso_integ"]; ok {
				integ = val(v)
			}
		}
	}

	c := "green"
	if conf == "R" {
		c = "red"
	}

	return fmt.Sprintf("<color blue>%s.%s</color> <<color %s>%s, %s</color>>", appName, parameterTypeName, c, conf, integ)
}

func formatReturnParam(s *sysl.Module, payload string) []string {
	types := make([]string, 0)
	if len(payload) > 0 {
		paramSlice := returnParamSpliterRE.Split(payload, -1)
		for _, param := range paramSlice {
			ptype := param

			valueTypeSlice := returnTypeValueSpliterRE.Split(param, -1)
			if len(valueTypeSlice) == 2 {
				ptype = valueTypeSlice[1]
			}

			if _, ok := sysl.Type_Primitive_value[strings.ToUpper(ptype)]; !ok {
				if m := typeSetOfRE.FindStringSubmatch(ptype); len(m) > 0 {
					ptype = m[1]
				}
				if m := typeOneOfRE.FindStringSubmatch(ptype); len(m) > 0 {
					types = append(types, typeListSpliterRE.Split(m[1], -1)...)
				} else {
					types = append(types, ptype)
				}
			}
		}
	}

	rargs := make([]string, 0, len(types))
	for _, t := range types {
		if !strings.Contains(t, "...") && strings.Contains(t, ".") {
			aps := strings.Split(t, ".")
			if len(aps) > 1 {
				rarg := formatArgs(s, aps[0], aps[1])
				if len(rarg) > 0 {
					rargs = append(rargs, rarg)
				}
			}
		} else {
			rargs = append(rargs, t)
		}
	}
	return rargs
}

func getReturnPayload(s *sysl.Module, stmts []*sysl.Statement) string {
	for _, v := range stmts {
		var subStmts []*sysl.Statement
		switch stmt := v.Stmt.(type) {
		case *sysl.Statement_Call, *sysl.Statement_Action:
			continue
		case *sysl.Statement_Ret:
			return stmt.Ret.GetPayload()
		case *sysl.Statement_Alt:
			for _, c := range stmt.Alt.Choice {
				if p := getReturnPayload(s, c.Stmt); len(p) > 0 {
					return p
				}
			}
		case *sysl.Statement_Cond:
			subStmts = stmt.Cond.Stmt
		case *sysl.Statement_Loop:
			subStmts = stmt.Loop.Stmt
		case *sysl.Statement_LoopN:
			subStmts = stmt.LoopN.Stmt
		case *sysl.Statement_Foreach:
			subStmts = stmt.Foreach.Stmt
		case *sysl.Statement_Group:
			subStmts = stmt.Group.Stmt
		}

		if p := getReturnPayload(s, subStmts); len(p) > 0 {
			return p
		}
	}
	return ""
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

func normalizeEndpointName(endpointName string) string {
	return endpointLabelReplaceRE.ReplaceAllLiteralString(endpointName, " ⬄ ")
}
