package parser

import (
	"encoding/json"
	"strings"

	"github.com/anz-bank/sysl/sysl2/proto"
	"github.com/sirupsen/logrus"
)

func makeStringAtom(str string) *sysl.Atom {
	return &sysl.Atom{
		Union: &sysl.Atom_String_{
			String_: str,
		},
	}
}

func makeRuleNameAtom(ruleName string) *sysl.Atom {
	return &sysl.Atom{
		Union: &sysl.Atom_Rulename{
			Rulename: &sysl.RuleName{
				Name: ruleName,
			},
		},
	}
}

func makeQuantifier(item interface{}) *sysl.Quantifier {
	var q *sysl.Quantifier
	qs := item.([]interface{})
	_, quantifier := ruleSeq(qs[0], "quantifier")

	if quantifier != nil {
		switch symbolTerm(quantifier[0]).tok.text {
		case "*":
			q = makeQuantifierZeroPlus()
		case "+":
			q = makeQuantifierOnePlus()
		case "?":
			q = makeQuantifierOptional()
		default:
			panic("not implemented yet.")
		}
	}
	return q
}

func makeTerm(a *sysl.Atom, q *sysl.Quantifier) *sysl.Term {
	return &sysl.Term{Atom: a, Quantifier: q}
}

func makeAtom(term interface{}) *sysl.Atom {
	var a *sysl.Atom

	atomType, atom := ruleSeq(term, "atom")

	switch atomType {
	case 0: // STRING
		tokText := symbolTerm(atom[0]).tok.text
		tokText = strings.Replace(tokText, `'`, `"`, 2)
		var val string
		if json.Unmarshal([]byte(tokText), &val) == nil {
			tokText = val
		}
		a = makeStringAtom(tokText)
	case 1: // lowercaseName
		a = makeRuleNameAtom(symbolTerm(atom[0]).tok.text)
	case 2: // '(' choice ')'
		choice := atom[1]
		c, r := ruleSeq(choice, "choice")
		if c != 0 {
			panic("unexpected index for rule choice")
		}
		a = &sysl.Atom{
			Union: &sysl.Atom_Choices{
				Choices: buildChoice(r),
			},
		}
	default:
		panic("not implemented yet.")
	}
	return a
}

func symbolTerm(item interface{}) symbol {
	return item.(symbol)
}

func getChoice(choice map[int][]interface{}) (int, []interface{}) {
	if len(choice) != 1 {
		panic("choice should only have 1 sequence")
	}
	for c := range choice {
		seq := choice[c]
		return c, seq
	}
	return -1, nil
}

// ruleSeq returns Rule.Choice.Sequence
func ruleSeq(item interface{}, rulename string) (int, []interface{}) {
	rule, ok := item.(map[string]map[int][]interface{})

	if ok {
		return getChoice(rule[rulename])
	}
	return -1, nil
}

func buildSequence(s0 []interface{}) *sysl.Sequence {
	terms := make([]*sysl.Term, 0)

	if s0 != nil {
		for _, term := range s0[0].([]interface{}) {
			_, t0 := ruleSeq(term, "term")
			atom := makeAtom(t0[0])
			quantifier := makeQuantifier(t0[1])
			terms = append(terms, makeTerm(atom, quantifier))
		}
	}
	return makeSequence(terms...)
}

func buildChoice(choice []interface{}) *sysl.Choice {
	choiceS := make([]*sysl.Sequence, 0)
	var s0 []interface{}
	_, s0 = ruleSeq(choice[0], "seq")

	choiceS = append(choiceS, buildSequence(s0))
	if len(choice) > 1 {
		x := choice[1].([]interface{})
		if x[0] != nil {
			for _, seq := range x {
				tt := seq.(map[int][]interface{})
				_, s0 = ruleSeq(tt[0][1], "seq")
				choiceS = append(choiceS, buildSequence(s0))
			}
		}
	}
	return &sysl.Choice{Sequence: choiceS}
}

func buildRule(ast interface{}) *sysl.Rule {
	_, rule := ruleSeq(ast, "rule")
	_, lhs := ruleSeq(rule[0], "lhs")
	ruleName, _ := makeRule(symbolTerm(lhs[0]).tok.text)
	_, rhs := ruleSeq(rule[2], "rhs")
	_, choice := ruleSeq(rhs[0], "choice")

	return &sysl.Rule{Name: ruleName, Choices: buildChoice(choice)}
}

// grammar := rule+
// rule := lhs ':' rhs ';'
// lhs := lowercaseName
// rhs := choice
// choice := seq ( '|' seq)*
// seq := term+
// term := atom quantifier?
// atom := STRING | ruleName | '(' choice  ')'
func buildGrammar(name string, start string, ast []interface{}) *sysl.Grammar {
	g := sysl.Grammar{
		Name:  name,
		Start: start,
		Rules: make(map[string]*sysl.Rule),
	}
	_, grammar := ruleSeq(ast[0], "grammar")

	for _, r := range grammar[0].([]interface{}) {
		rule := buildRule(r)
		g.Rules[rule.GetName().Name] = rule
	}
	return &g
}

// ParseEBNF Parses and build the EBNF grammar
func ParseEBNF(ebnfText string, name string, start string) *sysl.Grammar {
	p := makeParser(makeEBNF(), ebnfText)
	actual := make([]token, 0)

	for {
		tok := p.l.nextToken()
		if tok.id == -1 {
			break
		}
		actual = append(actual, tok)
	}

	result, tree := p.parseGrammar(&actual)
	if !result {
		logrus.Printf("unable to parse text=\n%s\n", ebnfText)
		return nil
	}
	return buildGrammar(name, start, tree)
}
