package main

import (
	"testing"
)

// S –> bab | bA
// A –> d | cA

func TestGrammar(t *testing.T) {
	tokens := []int{2, 3, 2, -1}
	text := "bab"

	g := makeGrammar1()
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)

	for i := range tokens {
		tok := lex.nextToken()
		if tok != tokens[i] {
			t.Errorf("got the wrong token at %d expected: %d, got %d", i, tokens[i], tok)
		}
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "S")

	if !result {
		t.Error("failed to parse " + text)
	}
}

func TestGrammar1(t *testing.T) {
	tokens := []int{2, 1, 1, 1, 0, -1}
	text := "bcccd"

	g := makeGrammar1()
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)

	for i := range tokens {
		tok := lex.nextToken()
		if tok != tokens[i] {
			t.Errorf("got the wrong token at %d expected: %d, got %d", i, tokens[i], tok)
		}
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "S")

	if !result {
		t.Error("failed to parse " + text)
	}
}

func TestGrammarIncompleteInput(t *testing.T) {
	tokens := []int{2, 3, -1}
	text := "ba"

	g := makeGrammar1()
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)

	for i := range tokens {
		tok := lex.nextToken()
		if tok != tokens[i] {
			t.Errorf("got the wrong token at %d expected: %d, got %d", i, tokens[i], tok)
		}
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "S")

	if result != false {
		t.Error("failed to parse " + text)
	}
}

func TestEXPR1(t *testing.T) {
	text := "1 + 3 * 7"
	g := makeEXPR()

	terms := getTerminals(g.Rules)

	if len(terms) != 7 {
		t.Error("got incorrect number of terms", terms)
	}
	lex := makeLexer(text, terms)
	tokens := make([]int, 6)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "E")
	if !result {
		t.Error("failed to parse " + text)
	}
}

func TestOBJ(t *testing.T) {
	text := "{}"
	g := makeRepeatSeq(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 3)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if !result {
		t.Error("failed to parse " + text)
	}

}

func TestOBJ2(t *testing.T) {
	text := "{123}"
	g := makeRepeatSeq(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 4)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if !result {
		t.Error("failed to parse " + text)
	}

}

func TestOBJ3(t *testing.T) {
	text := "{123, 246, 567}"
	g := makeRepeatSeq(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 8)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if !result {
		t.Error("failed to parse " + text)
	}

}

func TestOBJPlusNegative(t *testing.T) {
	text := "{123}"
	// NOTE THE +
	// obj
	//    : '{' pair (',' pair)+ '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOnePlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 4)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if result != false {
		t.Error("failed to parse " + text)
	}

}

func TestOBJPlusPositive(t *testing.T) {
	text := "{123, 234}"
	// NOTE THE +
	// obj
	//    : '{' pair (',' pair)+ '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOnePlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 6)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if !result {
		t.Error("failed to parse " + text)
	}

}

func TestOBJOptional(t *testing.T) {
	text := "{}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOptional())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 3)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestOBJOptional1(t *testing.T) {
	text := "{123}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOptional())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 4)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestOBJOptional2(t *testing.T) {
	text := "{123, 234}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOptional())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 6)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestOBJOptionalNegative(t *testing.T) {
	text := "{123, 234, 567}"
	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOptional())
	terms := getTerminals(g.Rules)

	if len(terms) != 4 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 8)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "obj")
	if result != false {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_1(t *testing.T) {
	text := "{}"
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 3)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_2(t *testing.T) {
	text := `{ "abc" : 123 }`
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 6)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_3(t *testing.T) {
	text := `{ 
				"abc" : 123 ,
				"def" : 4563456
			}`
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 10)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_Array1(t *testing.T) {
	text := `[]`
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 3)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	if tokens[len(tokens)-1] != -1 {
		t.Error("did not get the complete set of tokens")
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_Array2(t *testing.T) {
	text := `[ "abc" ]`
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 4)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	if tokens[len(tokens)-1] != -1 {
		t.Error("did not get the complete set of tokens")
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_Array3(t *testing.T) {
	text := `[ "abc" , 123 ]`
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 6)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	if tokens[len(tokens)-1] != -1 {
		t.Error("did not get the complete set of tokens")
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_Array4(t *testing.T) {
	text := `[ {"abc" : 123} ]`
	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 8)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	if tokens[len(tokens)-1] != -1 {
		t.Error("did not get the complete set of tokens")
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestJSON_Array5(t *testing.T) {
	text := `{
		"array": [
			{
				"abc": 123
			}
		]
	}`

	g := makeJSON(makeQuantifierZeroPlus())
	terms := getTerminals(g.Rules)

	if len(terms) != 8 {
		t.Error("got incorrect number of terms", terms)
	}

	lex := makeLexer(text, terms)
	tokens := make([]int, 12)

	for i := range tokens {
		tokens[i] = lex.nextToken()
	}
	if tokens[len(tokens)-1] != -1 {
		t.Error("did not get the complete set of tokens")
	}
	result := checkGrammar(g, tokens[:len(tokens)-1], "json")
	if result != true {
		t.Error("failed to parse " + text)
	}

}

func TestEBNF1(t *testing.T) {
	g := makeEBNF()
	terms := getTerminals(g.Rules)

	if len(terms) != 10 {
		t.Error("got incorrect number of terms", terms)
	}
	content := `expr : INT | ID | expr;`

	tokens := []int{1, 8, 0, 4, 0, 4, 1, 9, -1}

	lex := makeLexer(content, terms)
	for i := range tokens {
		tok := lex.nextToken()
		if tok != tokens[i] {
			t.Errorf("got the wrong token at %d expected: %d, got %d", i, tokens[i], tok)
		}
	}

	result := checkGrammar(g, tokens[:len(tokens)-1], "rule")
	if !result {
		t.Error("failed to parse \n" + content)
	}

}
