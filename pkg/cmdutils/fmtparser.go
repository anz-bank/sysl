package cmdutils

import (
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var (
	ItemReDefault    = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|%\()`)
	ItemReStatement  = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|[|)]|%\()`)
	ItemReEnd        = regexp.MustCompile(`((?:[^%]|%[^(\n]|\n)*?)($|\)|%\()`)
	ItemReVarStart   = regexp.MustCompile(`^%\(`)
	ItemReVar        = regexp.MustCompile(`^(@?\w+)`)
	ItemReCondOper   = regexp.MustCompile(`^[!=]=`)
	ItemReSearch     = regexp.MustCompile(`^~/([^/]+)/`)
	ItemReCondVal    = regexp.MustCompile(`^\'([\w ]+)\'`)
	ItemReStmtOper   = regexp.MustCompile(`^[=?]`)
	ItemReNoStmtOper = regexp.MustCompile(`^\|`)
	ItemReStmtEnd    = regexp.MustCompile(`^\)`)
)

const (
	MatchSymbol = iota + 1
	MatchWord
	MatchLookahead
)

type FormatParser struct {
	Self   string
	CurPos int
	Stk    []string
	Result string
	Oper   string
}

func MakeFormatParser(fmtStr string) *FormatParser {
	return &FormatParser{
		Self: fmtStr,
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
	MergeAttributesMap(attrs, p.Attrs)

	return fp.Parse(attrs)
}

func (fp *FormatParser) LabelApp(appname, controls string, attrs map[string]*sysl.Attribute) string {
	valMap := map[string]string{
		"appname":  appname,
		"controls": controls,
	}
	MergeAttributesMap(valMap, attrs)

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
	MergeAttributesMap(valMap, attrs)

	return fp.Parse(valMap)
}

func (fp *FormatParser) FmtOutput(appname, epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	valMap := map[string]string{
		"appname":    appname,
		"epname":     epname,
		"eplongname": eplongname,
	}
	MergeAttributesMap(valMap, attrs)

	return fp.Parse(valMap)
}

func (fp *FormatParser) Parse(attrs map[string]string) string {
	logrus.Debugf("self: %s", fp.Self)
	logrus.Debugf("attrs: %v", attrs)
	fp.Expansions(ItemReDefault, attrs)
	formatted := strings.ReplaceAll(fp.Result, "\n", "\\n")
	fp.Clear()
	logrus.Debugf("format string: %s", formatted)

	return formatted
}

func MergeAttributesMap(val map[string]string, attrs map[string]*sysl.Attribute) {
	for k, v := range attrs {
		val["@"+k] = v.GetS()
	}
}

func (fp *FormatParser) Expansions(re *regexp.Regexp, attrs map[string]string) {
	var result string
	for fp.Eat(re) {
		prefix := fp.Pop()
		prefix = RemovePercentSymbol(prefix)
		result += prefix

		if fp.Eat(ItemReVarStart) {
			var yesStmt, noStmt, varName, value string
			isYesStmt, isUseVal, isEqualOper, isSearched := false, true, false, false
			if !fp.Eat(ItemReVar) {
				panic("missing variable reference")
			}
			varName = fp.Pop()
			value = attrs[varName]

			// conditionals
			if fp.Eat(ItemReCondOper) {
				isEqualOper = true
				isUseVal = false
				conOper := fp.Oper
				fp.Oper = ""
				if !fp.Eat(ItemReCondVal) {
					panic("missing conditional value")
				}
				conVal := fp.Pop()
				if conOper == "==" && value == conVal {
					isYesStmt = true
				}
				if conOper == "!=" && value != conVal {
					isYesStmt = true
				}
			}

			if fp.Eat(ItemReSearch) {
				isSearched = true
				isUseVal = false
				reWordBoundary := regexp.MustCompile(fp.Pop())
				if reWordBoundary.MatchString(value) {
					isYesStmt = true
				}
			}

			have := fp.Eat(ItemReStmtOper)
			if have {
				isUseVal = false
				fp.Expansions(ItemReStatement, attrs)
				yesStmt = fp.Result
			}
			haveNot := fp.Eat(ItemReNoStmtOper)
			if haveNot {
				fp.Expansions(ItemReEnd, attrs)
				noStmt = fp.Result
			}
			if !fp.Eat(ItemReStmtEnd) {
				panic("unclosed expansion")
			}

			if isUseVal {
				result += value
				fp.Result = result
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
			fp.Result = result
		} else {
			fp.Result = result
			return
		}
	}
}

func (fp *FormatParser) Eat(re *regexp.Regexp) bool {
	matchStr := fp.Self[fp.CurPos:]
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
		fp.CurPos += len(subSelf[0])
		fp.Oper = subSelf[0]
	case MatchWord:
		insertion := subSelf[0]
		fp.Stk = append(fp.Stk, subSelf[1])
		fp.CurPos += len(insertion)
	case MatchLookahead:
		insertion := subSelf[1]
		fp.Stk = append(fp.Stk, insertion)
		fp.CurPos += len(insertion)
	}

	return subSelf != nil
}

func (fp *FormatParser) Pop() string {
	if fp.Stk == nil {
		return ""
	}

	n := len(fp.Stk)
	if n == 0 {
		return ""
	}
	popped := fp.Stk[n-1]
	fp.Stk = fp.Stk[:n-1]

	return popped
}

func (fp *FormatParser) Clear() {
	fp.Result = ""
	fp.CurPos = 0
	fp.Stk = []string{}
	fp.Oper = ""
}

func RemovePercentSymbol(src string) string {
	substitute := string([]byte{1})
	src = strings.ReplaceAll(src, "%%", substitute)
	src = strings.ReplaceAll(src, "%", "")
	src = strings.ReplaceAll(src, substitute, "%")

	return src
}
