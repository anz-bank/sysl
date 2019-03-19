package seqs

import (
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	log "github.com/sirupsen/logrus"
)

var (
	itemReDefault    = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|%\()`)
	itemReStatement  = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|[|)]|%\()`)
	itemReEnd        = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|\)|%\()`)
	itemReVarStart   = regexp.MustCompile(`^%\(`)
	itemReVar        = regexp.MustCompile(`(@?\w+)`)
	itemReCondOper   = regexp.MustCompile(`^[!=]=`)
	itemReSearch     = regexp.MustCompile(`~/([^/]+)/`)
	itemReCondVal    = regexp.MustCompile(`\'([\w ]+)\'`)
	itemReStmtOper   = regexp.MustCompile(`^[=?]`)
	itemReNoStmtOper = regexp.MustCompile(`^\|`)
	itemReStmtEnd    = regexp.MustCompile(`^\)`)
)

const (
	MATCH_SYMBOL = iota + 1
	MATCH_WORD
	MATCH_LOOK_AHEAD
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

func (f *FormatParser) Parse(attrs map[string]string) string {
	log.Debugf("self: %s", f.self)
	log.Debugf("attrs: %v", attrs)
	f.expansions(itemReDefault, attrs)
	formatted := f.result
	f.clear()
	log.Debugf("format string: %s", formatted)

	return formatted
}

func mergeAttributesMap(val map[string]string, attrs map[string]*sysl.Attribute) {
	for k, v := range attrs {
		val[k] = v.GetS()
	}
}

func (f *FormatParser) expansions(re *regexp.Regexp, attrs map[string]string) {
	var result string
	for f.eat(re) {
		prefix := f.pop()
		prefix = removePercentSymbol(prefix)
		result += prefix

		if f.eat(itemReVarStart) {
			var yesStmt, noStmt, varName, value string
			isYesStmt, isUseVal := false, true
			if !f.eat(itemReVar) {
				panic("missing variable reference")
			}
			varName = removeAtSymbol(f.pop())
			value = attrs[varName]

			// conditionals
			if f.eat(itemReCondOper) {
				isUseVal = false
				conOper := f.oper
				f.oper = ""
				if !f.eat(itemReCondVal) {
					panic("missing conditional value")
				}
				conVal := f.pop()
				if conOper == "==" && value == conVal {
					isYesStmt = true
				}
				if conOper == "!=" && value != conVal {
					isYesStmt = true
				}
			}

			if f.eat(itemReSearch) {
				isUseVal = false
				searchStr := f.pop()
				if strings.Contains(value, searchStr) {
					isYesStmt = true
				}
			}

			have := f.eat(itemReStmtOper)
			if have {
				isUseVal = false
				f.expansions(itemReStatement, attrs)
				yesStmt = f.result
			}
			have_not := f.eat(itemReNoStmtOper)
			if have_not {
				f.expansions(itemReEnd, attrs)
				noStmt = f.result
			}
			if !f.eat(itemReStmtEnd) {
				panic("unclosed expansion")
			}

			if isUseVal {
				result += value
				f.result = result
				continue
			}
			if !isYesStmt && value != "" {
				isYesStmt = true
			}
			if isYesStmt {
				result += yesStmt
			} else {
				result += noStmt
			}
			f.result = result
		} else {
			f.result = result
			return
		}
	}
}

func (f *FormatParser) eat(re *regexp.Regexp) bool {
	matchStr := string(f.self[f.curPos:])
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
	case MATCH_SYMBOL:
		f.curPos += len(subSelf[0])
		f.oper = subSelf[0]
	case MATCH_WORD:
		insertion := subSelf[0]
		f.stk = append(f.stk, subSelf[1])
		f.curPos += len(insertion)
	case MATCH_LOOK_AHEAD:
		insertion := subSelf[1]
		f.stk = append(f.stk, insertion)
		f.curPos += len(insertion)
	}

	return subSelf != nil
}

func (f *FormatParser) pop() string {
	if f.stk == nil {
		return ""
	}

	l := len(f.stk)
	if l == 0 {
		return ""
	}
	popped := f.stk[l-1]
	f.stk = f.stk[:l-1]

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

func removeAtSymbol(key string) string {
	return strings.Replace(key, "@", "", 1)
}
