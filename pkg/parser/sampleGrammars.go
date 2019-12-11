package parser

func makeQuantifierOptional() *Quantifier {
	return &Quantifier{Union: &Quantifier_Optional{}}
}

func makeQuantifierZeroPlus() *Quantifier {
	return &Quantifier{Union: &Quantifier_ZeroPlus{}}
}

func makeQuantifierOnePlus() *Quantifier {
	return &Quantifier{Union: &Quantifier_OnePlus{}}
}

func makeStringTerm(str string) *Term {
	return &Term{Atom: &Atom{Union: &Atom_String_{String_: str}}, Quantifier: nil}
}

func makeRegexpTerm(str string) *Term {
	return &Term{Atom: &Atom{Union: &Atom_Regexp{Regexp: str}}, Quantifier: nil}
}

func makeSequence(terms ...*Term) *Sequence {
	seq := Sequence{Term: terms}
	return &seq
}

func makeRule(name string) (*RuleName, *Term) {
	ruleName := RuleName{Name: name}
	ruleTerm := Term{Atom: &Atom{Union: &Atom_Rulename{Rulename: &ruleName}}, Quantifier: nil}
	return &ruleName, &ruleTerm
}

// S –> bab | bA
// A –> d | cA
func makeGrammar1() *Grammar {
	a := makeStringTerm("a")
	b := makeStringTerm("b")
	c := makeStringTerm("c")
	d := makeStringTerm("d")

	ruleNameA, A := makeRule("A")
	ruleNameS, _ := makeRule("S")

	return &Grammar{
		Name:  "test",
		Start: "S",
		Rules: map[string]*Rule{
			"S": {
				Name: ruleNameS,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(b, a, b),
						makeSequence(b, A)},
				},
			},
			"A": {
				Name: ruleNameA,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(d),
						makeSequence(c, A),
					},
				},
			},
		},
	}
}

// grammar := rule+
// rule := lhs ':' rhs ';'
// lhs := lowercaseName
// rhs := choice
// choice := seq ( '|' seq)*
// seq := term+
// term := atom quantifier?
// atom := STRING | ruleName | '(' choice  ')'
func makeEBNF() *Grammar {
	star := makeRegexpTerm("[*]")
	plus := makeRegexpTerm("[+]")
	qn := makeRegexpTerm("[?]")
	alt := makeRegexpTerm("[|]")
	colon := makeRegexpTerm("[:]")
	semiColon := makeRegexpTerm("[;]")
	openParen := makeRegexpTerm("[(]")
	closeParen := makeRegexpTerm("[)]")
	STRING := makeRegexpTerm(`['][^']*[']`)

	ruleNameRef := makeRegexpTerm("[a-zA-Z][0-9a-zA-Z_]*")

	lhsName, lhsTerm := makeRule("lhs")
	rhsName, rhsTerm := makeRule("rhs")
	ruleName, ruleTerm := makeRule("rule")
	grammarName, _ := makeRule("grammar")
	choiceName, choiceTerm := makeRule("choice")
	seqName, seqTerm := makeRule("seq")
	atomName, atomTerm := makeRule("atom")
	termName, termTerm := makeRule("term")
	termTerm.Quantifier = makeQuantifierOnePlus()
	quantifierName, quantifierTerm := makeRule("quantifier")
	quantifierTerm.Quantifier = makeQuantifierOptional()

	zeroPlusChoiceTerm := Term{
		Atom: &Atom{
			Union: &Atom_Choices{
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(alt, seqTerm),
					},
				},
			},
		},
		Quantifier: makeQuantifierZeroPlus(),
	}

	ruleTerm.Quantifier = makeQuantifierOnePlus()

	return &Grammar{
		Name:  "EBNF",
		Start: "grammar",
		Rules: map[string]*Rule{
			"grammar": {
				Name: grammarName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(ruleTerm),
					},
				},
			},
			"rule": {
				Name: ruleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(lhsTerm, colon, rhsTerm, semiColon),
					},
				},
			},
			"lhs": {
				Name: lhsName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(ruleNameRef),
					},
				},
			},
			"rhs": {
				Name: rhsName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(choiceTerm),
					},
				},
			},
			"choice": {
				Name: choiceName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(seqTerm, &zeroPlusChoiceTerm),
					},
				},
			},
			"seq": {
				Name: seqName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(termTerm),
					},
				},
			},
			"term": {
				Name: termName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(atomTerm, quantifierTerm),
					},
				},
			},
			"atom": {
				Name: atomName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(STRING),
						makeSequence(ruleNameRef),
						makeSequence(openParen, choiceTerm, closeParen),
					},
				},
			},
			"quantifier": {
				Name: quantifierName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(star),
						makeSequence(plus),
						makeSequence(qn),
					},
				},
			},
		}}
}

// E  -> T E'
// E' -> + T E' | -TE' |epsilon
// T  -> F T'
// T' -> * F T' | /FT' |epsilon
// F  -> (E) | int
func makeEXPR() *Grammar {
	plus := makeRegexpTerm("[+]")
	minus := makeRegexpTerm("[-]")
	star := makeRegexpTerm("[*]")
	divide := makeRegexpTerm("[/]")
	openParen := makeRegexpTerm("[(]")
	closeParen := makeRegexpTerm("[)]")
	integer := makeRegexpTerm("[0-9]+")

	ERuleName, ETerm := makeRule("E")
	ETailRuleName, ETailTerm := makeRule("ETail")
	TRuleName, TTerm := makeRule("T")
	TTailRuleName, TTailTerm := makeRule("TTail")
	factorRuleName, factorTerm := makeRule("factor")

	return &Grammar{
		Name:  "EXPR",
		Start: "E",
		Rules: map[string]*Rule{
			"E": {
				Name: ERuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(TTerm, ETailTerm),
					},
				},
			},
			"ETail": {
				Name: ETailRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(plus, TTerm, ETailTerm),
						makeSequence(minus, TTerm, ETailTerm),
						makeSequence(nil),
					},
				},
			},
			"T": {
				Name: TRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(factorTerm, TTailTerm),
					},
				},
			},
			"TTail": {
				Name: TTailRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(star, factorTerm, TTailTerm),
						makeSequence(divide, factorTerm, TTailTerm),
						makeSequence(nil),
					},
				},
			},
			"factor": {
				Name: factorRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(openParen, ETerm, closeParen),
						makeSequence(integer),
					},
				},
			},
		},
	}
}

// obj
//    : '{' number (',' number)* '}'
//    | '{' '}'
//    ;
func makeRepeatSeq(quantifier *Quantifier) *Grammar {
	curlyOpen := makeRegexpTerm("[{]")
	curlyClosed := makeRegexpTerm("[}]")
	comma := makeRegexpTerm("[,]")
	number := makeRegexpTerm("[0-9]+")

	objRuleName, _ := makeRule("obj")
	obj2RuleName, obj2Term := makeRule("obj2")
	obj2Term.Quantifier = quantifier

	return &Grammar{
		Name:  "array",
		Start: "obj",
		Rules: map[string]*Rule{
			"obj": {
				Name: objRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(curlyOpen, number, obj2Term, curlyClosed),
						makeSequence(curlyOpen, curlyClosed),
					},
				},
			},
			"obj2": {
				Name: obj2RuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(comma, number),
					},
				},
			},
		},
	}
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
//    | array
func makeJSON(quantifier *Quantifier) *Grammar {
	// doubleQuote := makeStringTerm("\"")
	// singleQuote := makeStringTerm("'")
	curlyOpen := makeRegexpTerm("[{]")
	curlyClosed := makeRegexpTerm("[}]")
	comma := makeRegexpTerm("[,]")
	sqOpen := makeRegexpTerm("[[]")
	sqClose := makeRegexpTerm("[]]")
	colon := makeRegexpTerm("[:]")
	number := makeRegexpTerm("[0-9]+")
	STRING := makeRegexpTerm(`["][^"]*["]`)

	jsonRuleName, _ := makeRule("json")
	valueRuleName, valueTerm := makeRule("value")
	objRuleName, objTerm := makeRule("obj")
	pairRuleName, pairTerm := makeRule("pair")
	arrayRuleName, arrayTerm := makeRule("array")

	obj2Term := Term{
		Atom: &Atom{
			Union: &Atom_Choices{
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(comma, pairTerm),
					},
				},
			},
		},
		Quantifier: quantifier,
	}

	array2Term := Term{
		Atom: &Atom{
			Union: &Atom_Choices{
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(comma, valueTerm),
					},
				},
			},
		},
		Quantifier: quantifier,
	}

	return &Grammar{
		Name:  "json",
		Start: "json",
		Rules: map[string]*Rule{
			"obj": {
				Name: objRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(curlyOpen, pairTerm, &obj2Term, curlyClosed),
						makeSequence(curlyOpen, curlyClosed),
					},
				},
			},
			"value": {
				Name: valueRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(STRING),
						makeSequence(number),
						makeSequence(objTerm),
						makeSequence(arrayTerm),
					},
				},
			},
			"json": {
				Name: jsonRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(valueTerm),
					},
				},
			},
			"pair": {
				Name: pairRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(STRING, colon, valueTerm),
					},
				},
			},
			"array": {
				Name: arrayRuleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(sqOpen, valueTerm, &array2Term, sqClose),
						makeSequence(sqOpen, sqClose),
					},
				},
			},
		},
	}
}

func makeG2() *Grammar {
	a := makeStringTerm("a")
	b := makeStringTerm("b")
	d := makeStringTerm("d")

	SruleName, _ := makeRule("S")
	AruleName, ATerm := makeRule("A")
	BruleName, BTerm := makeRule("B")
	DruleName, DTerm := makeRule("D")

	return &Grammar{
		Name:  "G2",
		Start: "S",
		Rules: map[string]*Rule{
			"S": {
				Name: SruleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(ATerm, a),
					},
				},
			},
			"A": {
				Name: AruleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(BTerm, DTerm),
					},
				},
			},
			"B": {
				Name: BruleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(b), makeSequence(nil),
					},
				},
			},
			"D": {
				Name: DruleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(d), makeSequence(nil),
					},
				},
			},
		},
	}
}

func makeNestedGrammar() *Grammar {
	a := makeRegexpTerm("{[A-Za-z]+:")
	b := makeStringTerm(":}")
	SruleName, _ := makeRule("S")

	return &Grammar{
		Name:  "G2",
		Start: "S",
		Rules: map[string]*Rule{
			"S": {
				Name: SruleName,
				Choices: &Choice{
					Sequence: []*Sequence{
						makeSequence(a),
						makeSequence(b),
					},
				},
			},
		},
	}
}

// func main() {
// 	fmt.Println("parsing grammar")

// }
