package main

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/sysl2/proto"
)

type createTerm func(string) *sysl.Term

var termMap map[string]*sysl.Term

func init() {
	termMap = make(map[string]*sysl.Term)
}

func returnSecond(arg string) *sysl.Term {
	_, t := makeRule(arg)
	return t
}

func getTerm(str string, m createTerm) *sysl.Term {
	t, has := termMap[str]
	if has {
		return t
	}
	termMap[str] = m(str)
	return termMap[str]
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
		s := s0[0].([]interface{})

		for _, term := range s {
			var t *sysl.Term
			_, t0 := ruleSeq(term, "term")
			atomType, atom := ruleSeq(t0[0], "atom")

			switch atomType {
			case 0:
				tokText := symbolTerm(atom[0]).tok.text
				tokText = strings.Trim(tokText, `"`)
				t = getTerm(tokText, makeTerm)
			case 2:
				t = getTerm(symbolTerm(atom[0]).tok.text, returnSecond)
			case 3:
				panic("not implemented yet.")
			default:
			}
			terms = append(terms, t)
			qs := t0[1].([]interface{})
			_, quantifier := ruleSeq(qs[0], "quantifier")
			if quantifier != nil {
				fmt.Printf("%+v\n", quantifier[0])
				// TODO: need to handle quantifers
			}
		}
	}
	return makeSequence(terms...)
}

func buildChoice(choice []interface{}) *sysl.Choice {
	choiceS := make([]*sysl.Sequence, 0)
	for option, seq := range choice {
		var s0 []interface{}
		if option > 0 {
			t := seq.([]interface{})[0]
			tt := t.(map[int][]interface{})[0]

			_, s0 = ruleSeq(tt[1], "seq")
		} else {
			_, s0 = ruleSeq(seq, "seq")
		}
		choiceS = append(choiceS, buildSequence(s0))
	}
	return &sysl.Choice{Sequence: choiceS}
}

func buildRule(ast []interface{}) *sysl.Rule {
	_, rule := ruleSeq(ast[0], "rule")
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

	for _, r := range grammar {
		rule := buildRule(r.([]interface{}))
		g.Rules[rule.GetName().Name] = rule
	}
	return &g
}

func parseEBNF(ebnfText string, name string, start string) *sysl.Grammar {
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
		fmt.Printf("unable to parse text=\n%s\n", ebnfText)
		return nil
	}
	return buildGrammar(name, start, tree)
}
