package parser

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

func makeStringAtom(str string) *Atom {
	return &Atom{
		Union: &Atom_String_{
			String_: str,
		},
	}
}

func makeRuleNameAtom(ruleName string) *Atom {
	return &Atom{
		Union: &Atom_Rulename{
			Rulename: &RuleName{
				Name: ruleName,
			},
		},
	}
}

func makeQuantifier(item interface{}) *Quantifier {
	qs := item.([]interface{})

	if _, quantifier := ruleSeq(qs[0], "quantifier"); quantifier != nil {
		switch symbolTerm(quantifier[0]).tok.text {
		case "*":
			return makeQuantifierZeroPlus()
		case "+":
			return makeQuantifierOnePlus()
		case "?":
			return makeQuantifierOptional()
		default:
			panic("not implemented yet.")
		}
	}
	return nil
}

func makeTerm(a *Atom, q *Quantifier) *Term {
	return &Term{Atom: a, Quantifier: q}
}

func makeAtom(term interface{}) *Atom {
	atomType, atom := ruleSeq(term, "atom")

	switch atomType {
	case 0: // STRING
		tokText := symbolTerm(atom[0]).tok.text
		tokText = strings.Replace(tokText, `'`, `"`, 2)
		var val string
		if json.Unmarshal([]byte(tokText), &val) == nil {
			tokText = val
		}
		return makeStringAtom(tokText)
	case 1: // lowercaseName
		return makeRuleNameAtom(symbolTerm(atom[0]).tok.text)
	case 2: // '(' choice ')'
		c, r := ruleSeq(atom[1], "choice")
		if c != 0 {
			logrus.Errorf("unexpected index for choice: %d", c)
			panic("unexpected index for rule choice")
		}
		return &Atom{
			Union: &Atom_Choices{
				Choices: buildChoice(r),
			},
		}
	default:
		panic("not implemented yet.")
	}
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
	if rule, ok := item.(map[string]map[int][]interface{}); ok {
		return getChoice(rule[rulename])
	}
	return -1, nil
}

func buildSequence(s0 []interface{}) *Sequence {
	terms := []*Term{}

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

func buildChoice(choice []interface{}) *Choice {
	_, s0 := ruleSeq(choice[0], "seq")
	choiceS := []*Sequence{buildSequence(s0)}
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
	return &Choice{Sequence: choiceS}
}

func buildRule(ast interface{}) *Rule {
	_, rule := ruleSeq(ast, "rule")
	_, lhs := ruleSeq(rule[0], "lhs")
	ruleName, _ := makeRule(symbolTerm(lhs[0]).tok.text)
	_, rhs := ruleSeq(rule[2], "rhs")
	_, choice := ruleSeq(rhs[0], "choice")

	return &Rule{Name: ruleName, Choices: buildChoice(choice)}
}

// grammar := rule+
// rule := lhs ':' rhs ';'
// lhs := lowercaseName
// rhs := choice
// choice := seq ( '|' seq)*
// seq := term+
// term := atom quantifier?
// atom := STRING | ruleName | '(' choice  ')'
func buildGrammar(name string, start string, ast []interface{}) *Grammar {
	_, grammar := ruleSeq(ast[0], "grammar")
	rules := map[string]*Rule{}
	for _, r := range grammar[0].([]interface{}) {
		rule := buildRule(r)
		rules[rule.GetName().Name] = rule
	}

	return &Grammar{
		Name:  name,
		Start: start,
		Rules: rules,
	}
}

// ParseEBNF Parses and build the EBNF grammar
func ParseEBNF(ebnfText string, name string, start string) *Grammar {
	p := makeParser(makeEBNF(), ebnfText)
	actual := []token{}

	for {
		tok := p.l.nextToken()
		if tok.id == -1 {
			break
		}
		actual = append(actual, tok)
	}

	result, tree := p.parseGrammar(&actual)
	if result {
		return buildGrammar(name, start, tree)
	}
	logrus.Printf("unable to parse text=\n%s\n", ebnfText)
	return nil
}
