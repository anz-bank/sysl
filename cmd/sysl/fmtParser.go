package main

import (
	"regexp"
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/proto_old"
	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var (
	itemReDefault    = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|%\()`)
	itemReStatement  = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|[|)]|%\()`)
	itemReEnd        = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|\)|%\()`)
	itemReVarStart   = regexp.MustCompile(`^%\(`)
	itemReVar        = regexp.MustCompile(`^(@?\w+)`)
	itemReCondOper   = regexp.MustCompile(`^[!=]=`)
	itemReSearch     = regexp.MustCompile(`^~/([^/]+)/`)
	itemReCondVal    = regexp.MustCompile(`^\'([\w ]+)\'`)
	itemReStmtOper   = regexp.MustCompile(`^[=?]`)
	itemReNoStmtOper = regexp.MustCompile(`^\|`)
	itemReStmtEnd    = regexp.MustCompile(`^\)`)
)

const (
	MatchSymbol = iota + 1
	MatchWord
	MatchLookahead
)

type FormatParser struct {
	self   string
	curPos int
	stk    []string
	result string
	oper   string
}

func MakeFormatParser(fmtStr string) *FormatParser {
	return &FormatParser{
		self: fmtStr,
	}
}

func (fp *FormatParser) LabelEndpoint(p *EndpointLabelerParam) string {
	attrs := map[string]string{
		"epname":       p.EndpointName,
		"human":        p.Human,
		"human_sender": p.HumanSender,
		"args":         p.Args,
		"patterns":     p.Patterns,
		"needs_int":    p.NeedsInt,
		"controls":     p.Controls,
	}
	mergeAttributesMap(attrs, p.Attrs)

	return fp.Parse(attrs)
}

func (fp *FormatParser) LabelApp(appname, controls string, attrs map[string]*sysl.Attribute) string {
	valMap := map[string]string{
		"appname":  appname,
		"controls": controls,
	}
	mergeAttributesMap(valMap, attrs)

	return fp.Parse(valMap)
}

func (fp *FormatParser) LabelClass(classname string) string {
	valMap := map[string]string{
		"classname": classname,
	}
	return fp.Parse(valMap)
}

func (fp *FormatParser) FmtSeq(epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	valMap := map[string]string{
		"epname":     epname,
		"eplongname": eplongname,
	}
	mergeAttributesMap(valMap, attrs)

	return fp.Parse(valMap)
}

func (fp *FormatParser) FmtOutput(appname, epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	valMap := map[string]string{
		"appname":    appname,
		"epname":     epname,
		"eplongname": eplongname,
	}
	mergeAttributesMap(valMap, attrs)

	return fp.Parse(valMap)
}

func (fp *FormatParser) Parse(attrs map[string]string) string {
	log.Debugf("self: %s", fp.self)
	log.Debugf("attrs: %v", attrs)
	fp.expansions(itemReDefault, attrs)
	formatted := strings.Replace(fp.result, "\n", "\\n", -1)
	fp.clear()
	log.Debugf("format string: %s", formatted)

	return formatted
}

func mergeAttributesMap(val map[string]string, attrs map[string]*sysl.Attribute) {
	for k, v := range attrs {
		val["@"+k] = v.GetS()
	}
}

func (fp *FormatParser) expansions(re *regexp.Regexp, attrs map[string]string) {
	var result string
	for fp.eat(re) {
		prefix := fp.pop()
		prefix = removePercentSymbol(prefix)
		result += prefix

		if fp.eat(itemReVarStart) {
			var yesStmt, noStmt, varName, value string
			isYesStmt, isUseVal, isEqualOper, isSearched := false, true, false, false
			if !fp.eat(itemReVar) {
				panic("missing variable reference")
			}
			varName = fp.pop()
			value = attrs[varName]

			// conditionals
			if fp.eat(itemReCondOper) {
				isEqualOper = true
				isUseVal = false
				conOper := fp.oper
				fp.oper = ""
				if !fp.eat(itemReCondVal) {
					panic("missing conditional value")
				}
				conVal := fp.pop()
				if conOper == "==" && value == conVal {
					isYesStmt = true
				}
				if conOper == "!=" && value != conVal {
					isYesStmt = true
				}
			}

			if fp.eat(itemReSearch) {
				isSearched = true
				isUseVal = false
				reWordBoundary := regexp.MustCompile(fp.pop())
				if reWordBoundary.MatchString(value) {
					isYesStmt = true
				}
			}

			have := fp.eat(itemReStmtOper)
			if have {
				isUseVal = false
				fp.expansions(itemReStatement, attrs)
				yesStmt = fp.result
			}
			haveNot := fp.eat(itemReNoStmtOper)
			if haveNot {
				fp.expansions(itemReEnd, attrs)
				noStmt = fp.result
			}
			if !fp.eat(itemReStmtEnd) {
				panic("unclosed expansion")
			}

			if isUseVal {
				result += value
				fp.result = result
				continue
			}
			if !isSearched && !isEqualOper && value != "" {
				isYesStmt = true
			}
			if isYesStmt {
				result += yesStmt
			} else {
				result += noStmt
			}
			fp.result = result
		} else {
			fp.result = result
			return
		}
	}
}

func (fp *FormatParser) eat(re *regexp.Regexp) bool {
	matchStr := fp.self[fp.curPos:]
	if matchStr == "" {
		return false
	}
	subSelfArr := re.FindAllStringSubmatch(matchStr, 1)
	if subSelfArr == nil {
		return false
	}
	subSelf := subSelfArr[0]
	subSelfLen := len(subSelf)

	switch subSelfLen {
	case MatchSymbol:
		fp.curPos += len(subSelf[0])
		fp.oper = subSelf[0]
	case MatchWord:
		insertion := subSelf[0]
		fp.stk = append(fp.stk, subSelf[1])
		fp.curPos += len(insertion)
	case MatchLookahead:
		insertion := subSelf[1]
		fp.stk = append(fp.stk, insertion)
		fp.curPos += len(insertion)
	}

	return subSelf != nil
}

func (fp *FormatParser) pop() string {
	if fp.stk == nil {
		return ""
	}

	n := len(fp.stk)
	if n == 0 {
		return ""
	}
	popped := fp.stk[n-1]
	fp.stk = fp.stk[:n-1]

	return popped
}

func (fp *FormatParser) clear() {
	fp.result = ""
	fp.curPos = 0
	fp.stk = []string{}
	fp.oper = ""
}

func removePercentSymbol(src string) string {
	substitute := string([]byte{1})
	src = strings.Replace(src, "%%", substitute, -1)
	src = strings.Replace(src, "%", "", -1)
	src = strings.Replace(src, substitute, "%", -1)

	return src
}
