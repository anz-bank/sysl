package seqs

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
)

var out io.Writer = os.Stdout

var (
	itemRE                   = regexp.MustCompile(`(%\(\w+\))`)
	isoCtrlRE                = regexp.MustCompile("^iso_ctrl_(.*)_txt$")
	returnTypeValueSpliterRE = regexp.MustCompile(`\s*<:\s*`)
	typeSetOfRE              = regexp.MustCompile(`set\s+of\s+(.+)$`)
	typeOneOfRE              = regexp.MustCompile(`one\s+of\s*{(.+)}$`)
	typeListSpliterRE        = regexp.MustCompile(`\s*,\s*`)
	endpointLabelReplaceRE   = regexp.MustCompile(`^.*? -> `)
	endpointParserRE         = regexp.MustCompile(`(?P<appname>.*?)\s*<-\s*(?P<epname>.*?)(?:\s*\[upto\s+(?P<upto>.*)\])*$`)
)

func TransformBlackBoxes(blackboxes []*sysl.Attribute) [][]string {
	bbs := make([][]string, 0, len(blackboxes))
	for _, vals := range blackboxes {
		subBbs := make([]string, 0)
		for _, val := range vals.GetA().Elt {
			subBbs = append(subBbs, val.GetS())
		}
		if len(subBbs) > 0 {
			bbs = append(bbs, subBbs)
		}
	}

	return bbs
}

func ParseBlackBoxesFromArgument(blackboxFlags []string) [][]string {
	bbs := make([][]string, 0, len(blackboxFlags))
	for _, blackboxFlag := range blackboxFlags {
		subBbs := make([]string, 0)
		subBbs = append(subBbs, strings.Split(blackboxFlag, ",")...)
		if len(subBbs) > 0 {
			bbs = append(bbs, subBbs)
		}
	}

	return bbs
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

func ParseAttributesFormat(init, formatted string, attrs map[string]string) string {
	formatted = replaceAttributes(init, attrs)
	formatted = replaceAttributesWithRules(formatted, attrs)

	return formatted
}

func replaceAttributesWithRules(formatted string, attrs map[string]string) string {
	re := regexp.MustCompile(`%\(@?\w+([!=]=\'.*\')?\?[^\)\?]*\)`)
	result := re.FindString(formatted)
	if result != "" {
		var subFormatted, conditionExp, conditionVal, yesStat, noStat string
		re = regexp.MustCompile(`%\(@?\w+`)
		key := removeWrapper(re.FindString(result))
		val := attrs[key]
		fmt.Fprintln(out, fmt.Sprintf("key: %s, val: %s", key, val))
		re = regexp.MustCompile(`[!=]='\w+'`)
		conditionStat := re.FindString(result)
		if conditionStat != "" {
			conditionExp = string([]rune(conditionStat)[0:2])
			fmt.Fprintln(out, fmt.Sprintf("conditionExp: %s", conditionExp))
			conditionVal = removeConditionWrapper(conditionStat)
		}
		re = regexp.MustCompile(`\?[^\)\?]*`)
		stat := strings.Replace(re.FindString(result), "?", "", 1)
		stats := strings.Split(stat, "|")
		yesStat = stats[0]
		if len(stats) > 1 {
			noStat = stats[1]
		}
		isYesStat := false
		if conditionExp == "!=" && val != "" && val != conditionVal {
			isYesStat = true
		}
		if conditionExp == "==" && val != "" && val == conditionVal {
			isYesStat = true
		}
		if conditionExp == "" && val != "" {
			isYesStat = true
		}
		subFormatted = noStat
		if isYesStat {
			subFormatted = yesStat
		}
		fmt.Fprintln(out, fmt.Sprintf("isYesStat: %v", isYesStat))

		formatted = strings.Replace(formatted, result, subFormatted, 1)
		formatted = replaceAttributesWithRules(formatted, attrs)
	}

	return formatted
}

func replaceAttributes(formatted string, attrs map[string]string) string {
	re := regexp.MustCompile(`%\(@?\w+\)`)
	result := re.FindString(formatted)
	if result != "" {
		formatted = strings.Replace(formatted, result, string(attrs[removeWrapper(result)]), 1)
		formatted = replaceAttributes(formatted, attrs)
	}

	return formatted
}

func removeWrapper(origin string) string {
	replaced := strings.Replace(origin, "%(", "", 1)
	replaced = strings.Replace(replaced, "@", "", 1)
	replaced = strings.Replace(replaced, ")", "", 1)

	return replaced
}

func removeConditionWrapper(origin string) string {
	replaced := strings.Replace(origin, "!", "", 1)
	replaced = strings.Replace(replaced, "=", "", 2)
	replaced = strings.Replace(replaced, "'", "", 2)

	return replaced
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
			conf = val(t.Attrs["iso_conf"])
			integ = val(t.Attrs["iso_integ"])
		}
	}

	c := "green"
	if conf == "R" {
		c = "red"
	}

	return fmt.Sprintf("<color blue>%s.%s</color> <<color %s>%s, %s</color>>", appName, parameterTypeName, c, conf, integ)
}

func formatReturnParam(s *sysl.Module, payload string) []string {
	// the original regex from python is `,(?![^{]*\})`
	// however golang regex does not support negative look ahead
	// that is the reason why I write this split function
	split := func(s string) []string {
		slice := make([]string, 0)
		inBrace := 0
		startIndex := 0
		for i, v := range s {
			switch v {
			case ',':
				if inBrace == 0 {
					slice = append(slice, s[startIndex:i])
					startIndex = i + 1
				}
			case '{':
				inBrace++
			case '}':
				inBrace--
			}
		}
		if startIndex < len(s) {
			slice = append(slice, s[startIndex:])
		}
		return slice
	}
	types := make([]string, 0)
	if len(payload) > 0 {
		paramSlice := split(payload)
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

func getReturnPayload(stmts []*sysl.Statement) string {
	for _, v := range stmts {
		var subStmts []*sysl.Statement
		switch stmt := v.Stmt.(type) {
		case *sysl.Statement_Call, *sysl.Statement_Action:
			continue
		case *sysl.Statement_Ret:
			return stmt.Ret.GetPayload()
		case *sysl.Statement_Alt:
			for _, c := range stmt.Alt.Choice {
				if p := getReturnPayload(c.Stmt); len(p) > 0 {
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

		if p := getReturnPayload(subStmts); len(p) > 0 {
			return p
		}
	}
	return ""
}

func getAndFmtParam(s *sysl.Module, params []*sysl.Param) []string {
	r := make([]string, 0, len(params))
	for _, v := range params {
		if refType := v.GetType().GetTypeRef(); refType != nil {
			if ref := refType.GetRef(); ref != nil {
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
