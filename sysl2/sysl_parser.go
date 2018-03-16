package main

import (
	"fmt"
	"math"

	"anz-bank/sysl/sysl2/proto"
)

func getMinMaxCount(t *sysl.Term) (int, int) {
	if t.Quantifier == nil {
		return 1, 1
	}

	switch t.Quantifier.Union.(type) {
	case *sysl.Quantifier_Optional: // 0 or 1
		return 0, 1
	case *sysl.Quantifier_ZeroPlus: // *
		return 0, math.MaxInt32
	case *sysl.Quantifier_OnePlus: // +
		return 1, math.MaxInt32
		// case *sysl.Quantifier_Separator:
	}
	return 1, 1
}

func parse(g *sysl.Grammar, input []int, val interface{}) (bool, []int, []interface{}) {
	if len(input) == 0 {
		fmt.Println("got input length zero!!!")
	}
	result := false
	tree := make([]interface{}, 0)

	switch val.(type) {
	case *sysl.Sequence:
		tok := val.(*sysl.Sequence)

		for _, t := range tok.Term {
			if t == nil {
				// nil == epsilon
				fmt.Println("matched nil")
				result = true
				continue
			}
			minCount, maxCount := getMinMaxCount(t)
			matchCount := 0
			res := false
			subTree := make([]interface{}, 0)

			for matchCount < maxCount {
				var remaining []int
				var matchedTerm interface{}
				switch t.Atom.Union.(type) {
				case *sysl.Atom_Choices:
					choice := t.GetAtom().GetChoices()
					res, remaining, matchedTerm = parse(g, input, choice)
				case *sysl.Atom_Rulename:
					nt := t.GetAtom().GetRulename()
					fmt.Printf("checking %s\n", nt.Name)
					res, remaining, matchedTerm = parse(g, input, g.Rules[nt.Name])
				default: //Atom_String_ and Atom_Regexp
					if len(input) == 0 {
						fmt.Printf("input is empty\n")
						res = false
					} else {
						term := t.GetAtom().Id
						in := int32(input[0])
						res = term == in
						remaining = input[1:]
						matchedTerm = t
					}
				}
				if res {
					matchCount++
					input = remaining
					subTree = append(subTree, matchedTerm)
				} else {
					if matchCount < minCount {
						return false, nil, nil
					} else if minCount == 0 && matchCount == 0 {
						//TODO: check maxCount?
						subTree = append(subTree, matchedTerm)
						break
					} else if matchCount >= minCount {
						break
					}
				}
			}
			if len(subTree) > 1 {
				tree = append(tree, subTree)
			} else if len(subTree) == 1 {
				tree = append(tree, subTree[0])
			}
		}
		fmt.Println("out of loop")
		result = true
		fmt.Println(tree)
		return result, input, tree
	case *sysl.Rule:
		r := val.(*sysl.Rule)
		fmt.Println("got " + r.GetName().Name)
		res, remaining, subTree := parse(g, input, r.Choices)
		if res {
			fmt.Printf("matched rulename (%s)\n", r.GetName().Name)
			rule := make(map[string]interface{})
			rule[r.GetName().Name] = subTree[0]
			tree = append(tree, rule)
			fmt.Println(tree)
			return true, remaining, tree
		}
		tree = append(tree, nil)
	case *sysl.Choice:
		p := val.(*sysl.Choice)
		fmt.Printf("choices count : (%d)\n", len(p.Sequence))
		for i, alt := range p.Sequence {
			fmt.Printf("trying choice (%d)\n", i)
			res, remaining, subTree := parse(g, input, alt)
			if res {
				fmt.Printf("matched choice :(%d)\n", i)
				choice := make(map[int][]interface{})
				choice[i] = subTree
				tree = append(tree, choice)
				fmt.Println(tree)
				return true, remaining, tree
			}
		}
		result = false
		tree = append(tree, nil)
	}
	return result, input, tree
}

// CheckGrammar returns true if the grammar consumes the whole string
// g - grammar to use
// text - to parse
// Start rule
func checkGrammar(g *sysl.Grammar, tokens []int, start string) bool {
	result, out, tree := parse(g, tokens, g.Rules[start])
	fmt.Println(tree)
	return (result && len(out) == 0)
}
