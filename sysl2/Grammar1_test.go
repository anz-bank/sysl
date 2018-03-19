package main

import (
	"testing"

	sysl "github.com/anz-bank/sysl/sysl2/proto"
)

func testParser(g *sysl.Grammar, numTerms int, tokens []int, text string, expected bool, t *testing.T) (bool, []interface{}) {
	p := makeParser(g, text)

	if len(p.terminals) != numTerms {
		t.Error("got incorrect number of terms", p.terminals)
	}

	actual := make([]token, 0)

	for i, expected := range tokens {
		tok := p.l.nextToken()
		actual = append(actual, tok)
		if tok.id != expected {
			t.Errorf("got the wrong token at %d expected: %+v, got %+v", i, expected, tok)
		}
	}
	actual = actual[:len(actual)-1]
	result, tree := p.parseGrammar(&actual)

	if result != expected {
		t.Error("failed to parse " + text)
	}
	return result, tree
}

// S –> bab | bA
// A –> d | cA
func TestGrammar(t *testing.T) {
	tokens := []int{2, 3, 2, -1}
	text := "bab"

	testParser(makeGrammar1(), 4, tokens, text, true, t)
}

func TestGrammar1(t *testing.T) {
	tokens := []int{2, 1, 1, 1, 0, -1}
	text := "bcccd"

	testParser(makeGrammar1(), 4, tokens, text, true, t)
}

func TestGrammarIncompleteInput(t *testing.T) {
	tokens := []int{2, 3, -1}
	text := "ba"

	testParser(makeGrammar1(), 4, tokens, text, false, t)
}

func TestEXPR1(t *testing.T) {
	text := "1 + 3 * 7"
	tokens := []int{6, 0, 6, 2, 6, -1}

	testParser(makeEXPR(), 7, tokens, text, true, t)
}

func TestOBJ(t *testing.T) {
	text := "{}"
	tokens := []int{0, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierZeroPlus()), 4, tokens, text, true, t)
}

func TestOBJ2(t *testing.T) {
	text := "{123}"
	tokens := []int{0, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierZeroPlus()), 4, tokens, text, true, t)
}

func TestOBJ3(t *testing.T) {
	text := "{123, 246, 567}"
	tokens := []int{0, 1, 3, 1, 3, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierZeroPlus()), 4, tokens, text, true, t)
}

func TestOBJPlusNegative(t *testing.T) {
	text := "{123}"
	// NOTE THE +
	// obj
	//    : '{' pair (',' pair)+ '}'
	//    | '{' '}'
	//    ;
	tokens := []int{0, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierOnePlus()), 4, tokens, text, false, t)
}

func TestOBJPlusPositive(t *testing.T) {
	text := "{123, 234}"
	// NOTE THE +
	// obj
	//    : '{' pair (',' pair)+ '}'
	//    | '{' '}'
	//    ;
	tokens := []int{0, 1, 3, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierOnePlus()), 4, tokens, text, true, t)
}

func TestOBJOptional(t *testing.T) {
	text := "{}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	tokens := []int{0, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierOptional()), 4, tokens, text, true, t)
}

func TestOBJOptional1(t *testing.T) {
	text := "{123}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	tokens := []int{0, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierOptional()), 4, tokens, text, true, t)
}

func TestOBJOptional2(t *testing.T) {
	text := "{123, 234}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	tokens := []int{0, 1, 3, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierOptional()), 4, tokens, text, true, t)
}

func TestOBJOptionalNegative(t *testing.T) {
	text := "{123, 234, 567}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	tokens := []int{0, 1, 3, 1, 3, 1, 2, -1}

	testParser(makeRepeatSeq(makeQuantifierOptional()), 4, tokens, text, false, t)
}

func TestJSON_1(t *testing.T) {
	text := "{}"
	tokens := []int{3, 4, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_2(t *testing.T) {
	text := `{ "abc" : 123 }`
	tokens := []int{3, 5, 6, 7, 4, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_3(t *testing.T) {
	text := `{
        "abc" : 123 ,
        "def" : 4563456
    }`
	tokens := []int{3, 5, 6, 7, 1, 5, 6, 7, 4, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_Array1(t *testing.T) {
	text := `[]`
	tokens := []int{0, 2, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_Array2(t *testing.T) {
	text := `[ "abc" ]`
	tokens := []int{0, 5, 2, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_Array3(t *testing.T) {
	text := `[ "abc" , 123 ]`
	tokens := []int{0, 5, 1, 7, 2, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_Array4(t *testing.T) {
	text := `[ {"abc" : 123} ]`
	tokens := []int{0, 3, 5, 6, 7, 4, 2, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestJSON_Array5(t *testing.T) {
	text := `{
        "array": [
            {
                "abc": 123
            }
        ]
    }`
	tokens := []int{3, 5, 6, 0, 3, 5, 6, 7, 4, 2, 4, -1}

	testParser(makeJSON(makeQuantifierZeroPlus()), 8, tokens, text, true, t)
}

func TestEBNF1(t *testing.T) {
	text := `expr : INT | ID | expr;`
	tokens := []int{2, 9, 1, 5, 1, 5, 2, 10, -1}

	testParser(makeEBNF(), 11, tokens, text, true, t)
}

func TestEBNF2(t *testing.T) {
	text := `expr : INT*;`
	tokens := []int{2, 9, 1, 6, 10, -1}

	testParser(makeEBNF(), 11, tokens, text, true, t)
}

func TestEBNF3(t *testing.T) {
	text := `obj : "{" number* "}" ;`
	p := makeParser(makeEBNF(), text)

	actual := make([]token, 0)

	for {
		tok := p.l.nextToken()
		if tok.id == -1 {
			break
		}
		actual = append(actual, tok)
	}
	// actual = actual[:len(actual)-1]
	result, tree := p.parseGrammar(&actual)
	if !result {
		t.Errorf("unable to parse text=(%s)", text)
	}
	buildGrammar(tree)
}

func TestFirstSet1(t *testing.T) {
	g := makeEXPR()
	terms := makeBuilder().buildTerminalsList(g.Rules)

	if len(terms) != 7 {
		t.Error("got incorrect number of terms", terms)
	}

	first, follow := buildFirstFollowSet(g)

	if len(first) < 0 || len(follow) < 0 {
		t.Error("failed to calculate first set of E\n")
	}
}

func TestFirstSet2(t *testing.T) {
	g := makeG2()
	terms := makeBuilder().buildTerminalsList(g.Rules)

	if len(terms) != 3 {
		t.Error("got incorrect number of terms", terms)
	}

	first, follow := buildFirstFollowSet(g)

	if len(first) < 0 || len(follow) < 0 {
		t.Error("failed to calculate first set of E\n")
	}
}
