package main

import (
	"fmt"

	"anz-bank/sysl/sysl2/proto"
)

func makeQuantifierOptional() *sysl.Quantifier {
	return &sysl.Quantifier{Union: &sysl.Quantifier_Optional{}}
}

func makeQuantifierZeroPlus() *sysl.Quantifier {
	return &sysl.Quantifier{Union: &sysl.Quantifier_ZeroPlus{}}
}

func makeQuantifierOnePlus() *sysl.Quantifier {
	return &sysl.Quantifier{Union: &sysl.Quantifier_OnePlus{}}
}

func makeStringTerm(str string) *sysl.Term {
	return &sysl.Term{Atom: &sysl.Atom{Union: &sysl.Atom_String_{String_: str}}, Quantifier: nil}
}

func makeRegexpTerm(str string) *sysl.Term {
	return &sysl.Term{Atom: &sysl.Atom{Union: &sysl.Atom_Regexp{Regexp: str}}, Quantifier: nil}
}

func makeSequence(terms ...*sysl.Term) *sysl.Sequence {
	seq := sysl.Sequence{Term: terms}
	return &seq
}

func makeRule(name string) (*sysl.RuleName, *sysl.Term) {
	ruleName := sysl.RuleName{Name: name}
	ruleTerm := sysl.Term{Atom: &sysl.Atom{Union: &sysl.Atom_Rulename{Rulename: &ruleName}}, Quantifier: nil}
	return &ruleName, &ruleTerm
}

// S –> bab | bA
// A –> d | cA

func makeGrammar1() sysl.Grammar {
	a := makeStringTerm("a")
	b := makeStringTerm("b")
	c := makeStringTerm("c")
	d := makeStringTerm("d")

	ruleNameA, A := makeRule("A")

	a1 := makeSequence(d)
	a2 := makeSequence(c, A)

	choiceA := sysl.Choice{Sequence: []*sysl.Sequence{a1, a2}}
	ruleA := sysl.Rule{Name: ruleNameA, Choices: &choiceA}

	ruleNameS, _ := makeRule("S")
	s1 := makeSequence(b, a, b)
	s2 := makeSequence(b, A)
	choiceS := sysl.Choice{Sequence: []*sysl.Sequence{s1, s2}}
	ruleS := sysl.Rule{Name: ruleNameS, Choices: &choiceS}

	prods := map[string]*sysl.Rule{}
	prods["S"] = &ruleS
	prods["A"] = &ruleA

	return sysl.Grammar{Name: "test", Rules: prods}
}

func makeEBNF() sysl.Grammar {
	star := makeStringTerm("*")
	plus := makeStringTerm("+")
	qn := makeStringTerm("?")

	alt := makeStringTerm("|")
	colon := makeStringTerm(":")
	semiColon := makeStringTerm(";")
	openParen := makeStringTerm("(")
	closeParen := makeStringTerm(")")

	tokenName := makeRegexpTerm("[A-Z][0-9A-Z_]+")
	lowercaseName := makeRegexpTerm("[a-z][0-9a-z_]+")

	lhsName, lhsTerm := makeRule("lhs")
	rhsName, rhsTerm := makeRule("rhs")
	ruleName, _ := makeRule("rule")
	choiceName, choiceTerm := makeRule("choice")
	seqName, seqTerm := makeRule("seq")
	termName, termTerm := makeRule("term")
	atomName, atomTerm := makeRule("atom")
	quantifierName, quantifierTerm := makeRule("quantifier")

	zeroPlusChoiceSeq := makeSequence(alt, seqTerm)
	zeroPlusChoice := sysl.Choice{Sequence: []*sysl.Sequence{zeroPlusChoiceSeq}}

	zeroPlusChoiceTerm := sysl.Term{Atom: &sysl.Atom{Union: &sysl.Atom_Choices{Choices: &zeroPlusChoice}}}
	zeroPlusChoiceTerm.Quantifier = makeQuantifierZeroPlus()

	q1 := makeSequence(star)
	q2 := makeSequence(plus)
	q3 := makeSequence(qn)
	quantifierChoice := sysl.Choice{Sequence: []*sysl.Sequence{q1, q2, q3}}
	quantifier := sysl.Rule{Name: quantifierName, Choices: &quantifierChoice}

	atom1 := makeSequence(tokenName)
	atom2 := makeSequence(lowercaseName)
	atom3 := makeSequence(openParen, choiceTerm, closeParen)
	atomChoice := sysl.Choice{Sequence: []*sysl.Sequence{atom1, atom2, atom3}}
	atom := sysl.Rule{Name: atomName, Choices: &atomChoice}

	quantifierTerm.Quantifier = makeQuantifierOptional()
	term1 := makeSequence(atomTerm, quantifierTerm)
	termChoice := sysl.Choice{Sequence: []*sysl.Sequence{term1}}
	term := sysl.Rule{Name: termName, Choices: &termChoice}

	termTerm.Quantifier = makeQuantifierOnePlus()
	seq1 := makeSequence(termTerm)
	seqChoice := sysl.Choice{Sequence: []*sysl.Sequence{seq1}}
	seq := sysl.Rule{Name: seqName, Choices: &seqChoice}

	choice1 := makeSequence(seqTerm, &zeroPlusChoiceTerm)
	choiceChoice := sysl.Choice{Sequence: []*sysl.Sequence{choice1}}
	choice := sysl.Rule{Name: choiceName, Choices: &choiceChoice}

	lhs1 := makeSequence(lowercaseName)
	lhsChoice := sysl.Choice{Sequence: []*sysl.Sequence{lhs1}}
	lhs := sysl.Rule{Name: lhsName, Choices: &lhsChoice}

	rhs1 := makeSequence(choiceTerm)
	rhsChoice := sysl.Choice{Sequence: []*sysl.Sequence{rhs1}}
	rhs := sysl.Rule{Name: rhsName, Choices: &rhsChoice}

	rule1 := makeSequence(lhsTerm, colon, rhsTerm, semiColon)
	ruleChoice := sysl.Choice{Sequence: []*sysl.Sequence{rule1}}
	rule := sysl.Rule{Name: ruleName, Choices: &ruleChoice}

	prods := map[string]*sysl.Rule{}
	prods["rule"] = &rule
	prods["lhs"] = &lhs
	prods["rhs"] = &rhs
	prods["choice"] = &choice
	prods["seq"] = &seq
	prods["term"] = &term
	prods["atom"] = &atom
	prods["quantifier"] = &quantifier

	return sysl.Grammar{Name: "EBNF", Rules: prods}

}

// E  -> T E'
// E' -> + T E' | -TE' |epsilon
// T  -> F T'
// T' -> * F T' | /FT' |epsilon
// F  -> (E) | int
func makeEXPR() sysl.Grammar {
	plus := makeStringTerm("+")
	minus := makeStringTerm("-")
	star := makeStringTerm("*")
	divide := makeStringTerm("/")
	openParen := makeStringTerm("(")
	closeParen := makeStringTerm(")")
	integer := makeRegexpTerm("[0-9]+")

	ERuleName, ETerm := makeRule("E")
	ETailRuleName, ETailTerm := makeRule("ETail")
	TRuleName, TTerm := makeRule("T")
	TTailRuleName, TTailTerm := makeRule("TTail")
	factorRuleName, factorTerm := makeRule("factor")

	E1 := makeSequence(TTerm, ETailTerm)
	EChoice := sysl.Choice{Sequence: []*sysl.Sequence{E1}}
	ERule := sysl.Rule{Name: ERuleName, Choices: &EChoice}

	ETail1 := makeSequence(plus, TTerm, ETailTerm)
	ETail2 := makeSequence(minus, TTerm, ETailTerm)
	ETail3 := makeSequence(nil)
	ETailChoice := sysl.Choice{Sequence: []*sysl.Sequence{ETail1, ETail2, ETail3}}
	ETailRule := sysl.Rule{Name: ETailRuleName, Choices: &ETailChoice}

	T1 := makeSequence(factorTerm, TTailTerm)
	TChoice := sysl.Choice{Sequence: []*sysl.Sequence{T1}}
	TRule := sysl.Rule{Name: TRuleName, Choices: &TChoice}

	TTail1 := makeSequence(star, factorTerm, TTailTerm)
	TTail2 := makeSequence(divide, factorTerm, TTailTerm)
	TTail3 := makeSequence(nil)
	TTailChoice := sysl.Choice{Sequence: []*sysl.Sequence{TTail1, TTail2, TTail3}}
	TTailRule := sysl.Rule{Name: TTailRuleName, Choices: &TTailChoice}

	factor1 := makeSequence(openParen, ETerm, closeParen)
	factor2 := makeSequence(integer)
	factorChoice := sysl.Choice{Sequence: []*sysl.Sequence{factor1, factor2}}
	factorRule := sysl.Rule{Name: factorRuleName, Choices: &factorChoice}

	prods := map[string]*sysl.Rule{}
	prods["E"] = &ERule
	prods["ETail"] = &ETailRule
	prods["T"] = &TRule
	prods["TTail"] = &TTailRule
	prods["factor"] = &factorRule

	return sysl.Grammar{Name: "EXPR", Rules: prods}

}

// obj
//    : '{' number (',' number)* '}'
//    | '{' '}'
//    ;
func makeRepeatSeq(quantifier *sysl.Quantifier) sysl.Grammar {
	curlyOpen := makeStringTerm("{")
	curlyClosed := makeStringTerm("}")
	comma := makeStringTerm(",")
	number := makeRegexpTerm("[0-9]+")

	objRuleName, _ := makeRule("obj")
	obj2RuleName, obj2Term := makeRule("obj2")

	obj2 := makeSequence(comma, number)
	obj2Term.Quantifier = quantifier
	obj1 := makeSequence(curlyOpen, number, obj2Term, curlyClosed)
	obj3 := makeSequence(curlyOpen, curlyClosed)

	obj2Choice := sysl.Choice{Sequence: []*sysl.Sequence{obj2}}
	obj2Rule := sysl.Rule{Name: obj2RuleName, Choices: &obj2Choice}

	objChoice := sysl.Choice{Sequence: []*sysl.Sequence{obj1, obj3}}
	objRule := sysl.Rule{Name: objRuleName, Choices: &objChoice}

	prods := map[string]*sysl.Rule{}
	prods["obj"] = &objRule
	prods["obj2"] = &obj2Rule

	return sysl.Grammar{Name: "json", Rules: prods}

}

// json
//    : value
//    ;

// obj
//    : '{' pair (',' pair)* '}'
//    | '{' '}'
//    ;

// pair
//    : STRING ':' value
//    ;

// array
//    : '[' value (',' value)* ']'
//    | '[' ']'
//    ;

// value
//    : STRING
//    | NUMBER
//    | obj
//	  | array
func makeJSON(quantifier *sysl.Quantifier) sysl.Grammar {
	// doubleQuote := makeStringTerm("\"")
	// singleQuote := makeStringTerm("'")
	curlyOpen := makeStringTerm("{")
	curlyClosed := makeStringTerm("}")
	comma := makeStringTerm(",")
	sqOpen := makeStringTerm("[")
	sqClose := makeStringTerm("]")
	colon := makeStringTerm(":")
	number := makeRegexpTerm("[0-9]+")
	STRING := makeRegexpTerm(`["][^"]*["]`)

	jsonRuleName, _ := makeRule("json")
	valueRuleName, valueTerm := makeRule("value")
	objRuleName, objTerm := makeRule("obj")
	pairRuleName, pairTerm := makeRule("pair")
	arrayRuleName, arrayTerm := makeRule("array")

	subSequence := makeSequence(comma, pairTerm)
	obj2Choice := sysl.Choice{Sequence: []*sysl.Sequence{subSequence}}

	obj2Term := sysl.Term{Atom: &sysl.Atom{Union: &sysl.Atom_Choices{Choices: &obj2Choice}}, Quantifier: quantifier}
	obj1 := makeSequence(curlyOpen, pairTerm, &obj2Term, curlyClosed)
	obj2 := makeSequence(curlyOpen, curlyClosed)

	objChoice := sysl.Choice{Sequence: []*sysl.Sequence{obj1, obj2}}
	objRule := sysl.Rule{Name: objRuleName, Choices: &objChoice}

	value1 := makeSequence(STRING)
	value2 := makeSequence(number)
	value3 := makeSequence(objTerm)
	value4 := makeSequence(arrayTerm)
	valueChoice := sysl.Choice{Sequence: []*sysl.Sequence{value1, value2, value3, value4}}
	valueRule := sysl.Rule{Name: valueRuleName, Choices: &valueChoice}

	json1 := makeSequence(valueTerm)
	jsonChoice := sysl.Choice{Sequence: []*sysl.Sequence{json1}}
	jsonRule := sysl.Rule{Name: jsonRuleName, Choices: &jsonChoice}

	pair1 := makeSequence(STRING, colon, valueTerm)
	pairChoice := sysl.Choice{Sequence: []*sysl.Sequence{pair1}}
	pairRule := sysl.Rule{Name: pairRuleName, Choices: &pairChoice}

	arraySubSequence := makeSequence(comma, valueTerm)
	arraySubSequenceChoice := sysl.Choice{Sequence: []*sysl.Sequence{arraySubSequence}}

	array2Term := sysl.Term{Atom: &sysl.Atom{Union: &sysl.Atom_Choices{Choices: &arraySubSequenceChoice}}, Quantifier: quantifier}
	array1 := makeSequence(sqOpen, valueTerm, &array2Term, sqClose)
	array2 := makeSequence(sqOpen, sqClose)

	arrayChoice := sysl.Choice{Sequence: []*sysl.Sequence{array1, array2}}
	arrayRule := sysl.Rule{Name: arrayRuleName, Choices: &arrayChoice}

	prods := map[string]*sysl.Rule{}
	prods["obj"] = &objRule
	prods["value"] = &valueRule
	prods["json"] = &jsonRule
	prods["pair"] = &pairRule
	prods["array"] = &arrayRule

	return sysl.Grammar{Name: "json", Rules: prods}

}

func main() {
	fmt.Println("parsing grammar")

}
