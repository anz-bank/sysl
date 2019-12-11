package parser

import (
	"math"
	"sort"

	"github.com/sirupsen/logrus"
)

type parser struct {
	g         *Grammar
	l         *lexer
	terminals []string
	tokens    *[]token
}

type builder struct {
	arr    []string
	tokens map[string]int32
}

type symbol struct {
	tok  token
	term *Term
}

const FIRST = 0
const EPSILON = -2
const EOF = -1
const PushGrammar = -3
const PopGrammar = -4

func (b *builder) handleChoice(choice *Choice) {
	for _, s := range choice.Sequence {
		for _, t := range s.Term {
			if t == nil {
				continue
			}
			var str string
			switch t.Atom.Union.(type) {
			case *Atom_String_:
				str = t.GetAtom().GetString_()
			case *Atom_Regexp:
				str = t.GetAtom().GetRegexp()
				str = `\A(?:` + str + `)`
			case *Atom_Choices:
				b.handleChoice(t.GetAtom().GetChoices())
			}
			if str != "" {
				if id, has := b.tokens[str]; !has {
					b.tokens[str] = int32(len(b.arr))
					t.Atom.Id = int32(len(b.arr))
					b.arr = append(b.arr, str)
				} else {
					t.Atom.Id = id
				}
			}
		}
	}
}

// assigns new value to Atom.Id
func (b *builder) buildTerminalsList(rules map[string]*Rule) []string {
	ks := make([]string, 0, len(rules))
	for key := range rules {
		ks = append(ks, key)
	}
	sort.Strings(ks)
	for _, key := range ks {
		b.handleChoice(rules[key].Choices)
	}
	return b.arr
}

func makeBuilder() *builder {
	return &builder{
		arr:    make([]string, 0),
		tokens: make(map[string]int32),
	}
}

func makeParser(g *Grammar, text string) *parser {
	terms := makeBuilder().buildTerminalsList(g.Rules)
	lexer := makeLexer(text, terms)
	return &parser{
		g:         g,
		terminals: terms,
		l:         &lexer,
	}
}

// GetTermMinMaxCount returns acceptable min and max counts
// by looking at term's quantifier
func GetTermMinMaxCount(t *Term) (int, int) {
	if t.Quantifier == nil {
		return 1, 1
	}

	switch t.Quantifier.Union.(type) {
	case *Quantifier_Optional: // 0 or 1
		return 0, 1
	case *Quantifier_ZeroPlus: // *
		return 0, math.MaxInt32
	case *Quantifier_OnePlus: // +
		return 1, math.MaxInt32
		// case *sysl.Quantifier_Separator:
	}
	return 1, 1
}

func (p *parser) parse(g *Grammar, input int, val interface{}) (bool, int, []interface{}) {
	if input == len(*(p.tokens)) {
		logrus.Debug("got input length zero!!!")
	}
	result := false
	tree := []interface{}{}

	switch val := val.(type) {
	case *Sequence:
		for index, t := range val.Term {
			if t == nil {
				// nil => epsilon, see makeEXPR()
				logrus.Debug("matched nil")
				continue
			}
			minCount, maxCount := GetTermMinMaxCount(t)
			matchCount := 0
			var res bool
			subTree := make([]interface{}, 0)
			singleTerm := minCount == maxCount && minCount == 1

			for matchCount < maxCount {
				var remaining int
				var matchedTerm interface{}
				switch x := t.Atom.Union.(type) {
				case *Atom_Choices:
					var parseResult []interface{}
					res, remaining, parseResult = p.parse(g, input, x.Choices)
					matchedTerm = parseResult[0]
				case *Atom_Rulename:
					logrus.Debugf("checking Rule %s (singleTerm: %v)", x.Rulename.Name, singleTerm)
					var parseResult []interface{}
					res, remaining, parseResult = p.parse(g, input, g.Rules[x.Rulename.Name])
					matchedTerm = parseResult[0]
				default: //Atom_String_ and Atom_Regexp
					if input == len(*p.tokens) {
						logrus.Debugf("input is empty")
						res = false
					} else {
						in := (*p.tokens)[input]
						res = int(t.GetAtom().Id) == in.id
						remaining = input + 1
						matchedTerm = symbol{in, t}
					}
				}
				if res {
					matchCount++
					input = remaining
					subTree = append(subTree, matchedTerm)
					logrus.Debugf("%d >> matched (%v): %d  -- len(subtree=%d)", index, matchedTerm, matchCount, len(subTree))
				} else {
					if matchCount < minCount {
						return false, EOF, nil
					}
					if matchCount == 0 {
						//TODO: check maxCount?
						subTree = append(subTree, matchedTerm)
					}
					break
				}
			}
			if singleTerm {
				tree = append(tree, subTree[0])
			} else {
				tree = append(tree, subTree)
			}
		}
		logrus.Debug("out of loop")
		return true, input, tree
	case *Rule:
		logrus.Debug("Entering Rule " + val.GetName().Name)
		res, remaining, subTree := p.parse(g, input, val.Choices)
		if res {
			logrus.Debugf("matched rulename (%s)", val.GetName().Name)
			logrus.Debugf("got choice (%T)", subTree[0])
			rule := map[string]map[int][]interface{}{
				val.GetName().Name: subTree[0].(map[int][]interface{}),
			}
			return true, remaining, append(tree, rule)
		}
		logrus.Debug("did not match " + val.GetName().Name)
		tree = append(tree, nil)
	case *Choice:
		logrus.Debugf("choices count : (%d)", len(val.Sequence))
		for i, alt := range val.Sequence {
			logrus.Debugf("trying choice (%d)", i)
			res, remaining, subTree := p.parse(g, input, alt)
			if res {
				logrus.Debugf("matched choice :(%d)", i)
				choice := map[int][]interface{}{i: subTree}
				return true, remaining, append(tree, choice)
			}
		}
		result = false
		tree = append(tree, nil)
	}
	return result, input, tree
}

// parseGrammar returns (true, ast) if the grammar consumes the whole string
func (p *parser) parseGrammar(arr *[]token) (bool, []interface{}) {
	p.tokens = arr
	result, out, tree := p.parse(p.g, 0, p.g.Rules[p.g.Start])
	return result && len(*p.tokens) == out, tree
}

func setFromTerm(first map[string]*intSet, t *Term) *intSet {
	var firstSet *intSet
	switch t.Atom.Union.(type) {
	case *Atom_Choices:
		panic("not handled yet")
	case *Atom_Rulename:
		nt := t.GetAtom().GetRulename()
		firstSet = first[nt.Name]
	default: //Atom_String_ and Atom_Regexp
		s := make(intSet)
		s.add(int(t.GetAtom().Id))
		firstSet = &s
	}
	return firstSet
}

func hasEpsilon(t *Term, first map[string]*intSet) bool {
	if _, ok := t.Atom.Union.(*Atom_Rulename); ok {
		return setFromTerm(first, t).has(EPSILON)
	}
	return false
}

func isNonTerminal(t *Term) bool {
	switch t.Atom.Union.(type) {
	case *Atom_Choices, *Atom_Rulename:
		return true
	}
	return false
}

func buildFirstFollowSet(g *Grammar) (map[string]*intSet, map[string]*intSet) {
	first := make(map[string]*intSet)
	follow := make(map[string]*intSet)

	for key := range g.Rules {
		first[key] = &intSet{}
		follow[key] = &intSet{}
	}

	// Follow Set Rule 1
	// 1. First put $ (the end of input marker) in Follow(S) (S is the start symbol)
	follow[g.Start].add(EOF)

	// check for epsilon
	for ruleName, rule := range g.Rules {
		for _, seq := range rule.Choices.Sequence {
			if seq.Term[0] == nil {
				first[ruleName].add(EPSILON)
			}
		}
	}

	updated := true
	for updated {
		updated = false
		for ruleName, rule := range g.Rules {
			for _, seq := range rule.Choices.Sequence {
				numEpsilon := 0
				for _, t := range seq.Term {
					if t == nil {
						continue
					}
					if !hasEpsilon(t, first) {
						updated = first[ruleName].union(setFromTerm(first, t)) || updated
						break
					} else {
						y1 := setFromTerm(first, t).clone()
						y1.remove(EPSILON)
						updated = first[ruleName].union(y1) || updated
						numEpsilon++
					}
					if len(seq.Term) == numEpsilon {
						first[ruleName].add(EPSILON)
					}
				}
			}
		}
	}
	updated = true
	for updated {
		updated = false
		for ruleName, rule := range g.Rules {
			for _, seq := range rule.Choices.Sequence {
				last := len(seq.Term) - 1
				for i := range seq.Term {
					// follow
					if i < last {
						A := seq.Term[i]
						B := seq.Term[i+1]
						Bfirst := setFromTerm(first, B).clone()
						// Rule 4
						// If there is a production X → aAB, where FIRST(B) contains ε,
						// then everything in FOLLOW(X) is in FOLLOW(A)
						if isNonTerminal(B) {
							if Bfirst.has(EPSILON) {
								Afollow := setFromTerm(follow, A)
								updated = Afollow.union(follow[ruleName]) || updated
							}
						}

						// Rule 2.
						// If there is a production X → aAB,
						// (where a can be a whole string)
						// then everything in FIRST(B) except for ε is placed in FOLLOW(A).
						if isNonTerminal(A) {
							Bfirst.remove(EPSILON)
							Afollow := setFromTerm(follow, A)
							updated = Afollow.union(Bfirst) || updated
						}
					} else if i == last {
						// Rule 3
						// If there is a production X → aB, then everything in FOLLOW(X) is in FOLLOW(B)
						B := seq.Term[i]
						if B != nil && isNonTerminal(B) {
							Bfollow := setFromTerm(follow, B)
							updated = Bfollow.union(follow[ruleName]) || updated
						}
					}
				}
			}
		}
	}
	return first, follow
}
