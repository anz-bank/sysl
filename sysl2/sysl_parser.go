package main

import (
	"fmt"

	sysl "github.com/anz-bank/sysl/sysl2/proto"
)

var termMap map[string]*sysl.Term

func symbolTerm(item interface{}) symbol {
	return item.(symbol)
}

func ruleSeq(item interface{}, rulename string, choice int) []interface{} {
	// r := item.([]interface{})
	fmt.Printf("%T\n", item)
	rule, ok := item.(map[string]map[int][]interface{})
	if ok {
		Choice := rule[rulename]
		seq := Choice[choice]
		fmt.Printf("%d\n", len(seq))
		return seq
	}
	return nil
}

// rule := lhs ':' rhs ';'
// lhs := lowercaseName
// rhs := choice
// choice := seq ( '|' seq)*
// seq := term+
// term := atom quantifier?
// atom := STRING | ruleName | '(' choice  ')'
func buildGrammar(ast []interface{}) *sysl.Grammar {
	rule := ruleSeq(ast[0], "rule", 0)

	lhs := ruleSeq(rule[0], "lhs", 0)
	ruleName := symbolTerm(lhs[0]).tok.text
	fmt.Printf("%s\n", ruleName)

	rhs := ruleSeq(rule[2], "rhs", 0)
	choice := ruleSeq(rhs[0], "choice", 0)
	for _, seq := range choice {
		s0 := ruleSeq(seq, "seq", 0)
		if s0 != nil {
			s := s0[0].([]interface{})
			fmt.Printf("%d\n", len(s))

			for _, term := range s {
				t0 := ruleSeq(term, "term", 0)
				fmt.Printf("%d\n", len(t0))
				atom := ruleSeq(t0[0], "atom", 0)
				fmt.Printf("%T\n", atom[0])
				fmt.Printf("%+v\n", symbolTerm(atom[0]))
				quantifier := ruleSeq(t0[1], "quantifier", 0)
				fmt.Printf("%+v\n", quantifier)
			}
		}
	}
	return nil
}
