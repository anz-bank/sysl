package parser

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

type cases struct {
	text   string
	tokens []int
	result bool
}

func testParser(
	g *sysl.Grammar,
	numTerms int,
	tokens []int,
	text string,
	expectedResult bool,
	t *testing.T,
) (bool, []interface{}) {
	p := makeParser(g, text)

	assert.Len(t, p.terminals, numTerms)

	actual := make([]token, 0)

	for i, expected := range tokens {
		tok := p.l.nextToken()
		actual = append(actual, tok)
		assert.Equal(t, expected, tok.id, "%d: %#v", i, expected)
	}
	actual = actual[:len(actual)-1]
	result, tree := p.parseGrammar(&actual)
	assert.Equal(t, expectedResult, result)
	return result, tree
}

// S –> bab | bA
// A –> d | cA
func TestGrammar1(t *testing.T) {
	t.Parallel()

	cases := []cases{
		{"bab", []int{2, 3, 2, -1}, true},
		{"bcccd", []int{2, 1, 1, 1, 0, -1}, true},
		{"ba", []int{2, 3, -1}, false},
	}
	g := makeGrammar1()
	for i := range cases {
		testParser(g, 4, cases[i].tokens, cases[i].text, cases[i].result, t)
	}
}

func TestEXPR1(t *testing.T) {
	t.Parallel()

	text := "1 + 3 * 7"
	tokens := []int{6, 0, 6, 2, 6, -1}

	testParser(makeEXPR(), 7, tokens, text, true, t)
}

func TestOBJ(t *testing.T) {
	t.Parallel()

	g := makeRepeatSeq(makeQuantifierZeroPlus())
	cases := []cases{
		{"{}", []int{0, 2, -1}, true},
		{"{123}", []int{0, 1, 2, -1}, true},
		{"{123, 246, 567}", []int{0, 1, 3, 1, 3, 1, 2, -1}, true},
	}
	for i := range cases {
		testParser(g, 4, cases[i].tokens, cases[i].text, cases[i].result, t)
	}
}

func TestOBJPlus(t *testing.T) {
	t.Parallel()

	// NOTE THE +
	// obj
	//    : '{' pair (',' pair)+ '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOnePlus())
	cases := []cases{
		{"{123}", []int{0, 1, 2, -1}, false},
		{"{123, 234}", []int{0, 1, 3, 1, 2, -1}, true},
	}
	for i := range cases {
		testParser(g, 4, cases[i].tokens, cases[i].text, cases[i].result, t)
	}
}

func TestOBJOptional(t *testing.T) {
	t.Parallel()

	// NOTE THE ?
	// obj
	//    : '{' number (',' number)? '}'
	//    | '{' '}'
	//    ;
	g := makeRepeatSeq(makeQuantifierOptional())
	cases := []cases{
		{"{}", []int{0, 2, -1}, true},
		{"{123}", []int{0, 1, 2, -1}, true},
		{"{123, 234}", []int{0, 1, 3, 1, 2, -1}, true},
		{"{123, 234, 567}", []int{0, 1, 3, 1, 3, 1, 2, -1}, false},
	}
	for i := range cases {
		testParser(g, 4, cases[i].tokens, cases[i].text, cases[i].result, t)
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	g := makeJSON(makeQuantifierZeroPlus())
	cases := []cases{
		{"{}", []int{3, 4, -1}, true},
		{`{ "abc" : 123 }`, []int{3, 5, 6, 7, 4, -1}, true},
		{
			`{
      "abc" : 123 ,
      "def" : 4563456
      }
      `,
			[]int{3, 5, 6, 7, 1, 5, 6, 7, 4, -1}, true,
		},
		{`[]`, []int{0, 2, -1}, true},
		{`[ "abc" ]`, []int{0, 5, 2, -1}, true},
		{`[ "abc" , 123 ]`, []int{0, 5, 1, 7, 2, -1}, true},
		{`[ {"abc" : 123} ]`, []int{0, 3, 5, 6, 7, 4, 2, -1}, true},
		{`{
      "array": [
          {
              "abc": 123
          }
      ]
      }`, []int{3, 5, 6, 0, 3, 5, 6, 7, 4, 2, 4, -1}, true},
	}

	for i := range cases {
		testParser(g, 8, cases[i].tokens, cases[i].text, cases[i].result, t)
	}
}

func TestFirstSet1(t *testing.T) {
	t.Parallel()

	g := makeEXPR()
	terms := makeBuilder().buildTerminalsList(g.Rules)

	assert.Len(t, terms, 7)

	first, follow := buildFirstFollowSet(g)

	//TODO: revisit this test
	assert.NotEmpty(t, first, 0)
	assert.NotEmpty(t, follow, 0)
}

func TestFirstSet2(t *testing.T) {
	t.Parallel()

	g := makeG2()
	terms := makeBuilder().buildTerminalsList(g.Rules)

	if len(terms) != 3 {
		t.Error("got incorrect number of terms", terms)
	}

	first, follow := buildFirstFollowSet(g)

	//TODO: revisit this test
	assert.NotEmpty(t, first, 0)
	assert.NotEmpty(t, follow, 0)
}

func TestEBNF1(t *testing.T) {
	t.Parallel()

	g := makeEBNF()
	cases := []cases{
		{`expr : INT | ID | expr;`, []int{1, 8, 1, 4, 1, 4, 1, 9, -1}, true},
		{`expr : INT*;`, []int{1, 8, 1, 5, 9, -1}, true},
	}

	for i := range cases {
		testParser(g, 10, cases[i].tokens, cases[i].text, cases[i].result, t)
	}
}

func TestBuildEBNFGrammar(t *testing.T) {
	t.Parallel()

	// Both grammars are equivalent
	grammars := []string{
		`s : 'd' | 'c' s ; `,
		`s : 'c'* 'd' ; `,
	}
	text := "ccd"
	tokens := [][]int{
		{1, 1, 0, -1},
		{0, 0, 1, -1},
	}
	choiceBranch := []int{1, 0}

	for i := range grammars {
		g := ParseEBNF(grammars[i], "obj", "s")

		result, ast := testParser(g, 2, tokens[i], text, true, t)
		assert.True(t, result, "grammars[%d]", i)
		choiceActual, s := ruleSeq(ast[0], "s")
		assert.Equal(t, choiceBranch[i], choiceActual)
		assert.Len(t, s, 2)
	}
}

func TestBuildEBNFGrammar2(t *testing.T) {
	t.Parallel()

	text := `
        s : 'b' 'a' 'b' | 'b' a;
        a : 'd' | 'c' a;
        `
	g := ParseEBNF(text, "obj", "s")
	if len(g.Rules) != 2 {
		t.Errorf("incorrect number of rules")
	}
	text = "bab"
	tokens := []int{2, 3, 2, -1}

	result, ast := testParser(g, 4, tokens, text, true, t)
	assert.True(t, result)

	choiceActual, s := ruleSeq(ast[0], "s")
	assert.Zero(t, choiceActual)
	assert.Len(t, s, 3)
}

func TestBuildEBNFGrammar3(t *testing.T) {
	t.Parallel()

	text := `
        s : a (',' a)*;
        a : 'd';
        `
	g := ParseEBNF(text, "obj", "s")
	if len(g.Rules) != 2 {
		t.Errorf("incorrect number of rules")
	}
	text = "d,d"
	tokens := []int{0, 1, 0, -1}

	result, ast := testParser(g, 2, tokens, text, true, t)
	assert.True(t, result)

	choiceActual, s := ruleSeq(ast[0], "s")
	assert.Zero(t, choiceActual)
	assert.Len(t, s, 2)
}

func TestBuildEBNFGrammar4(t *testing.T) {
	t.Parallel()

	text := `
          s: a a a;
          a: 'a' | 'b' | 'c' | 'd';
        `
	g := ParseEBNF(text, "obj", "s")
	if len(g.Rules) != 2 {
		t.Errorf("incorrect number of rules")
	}
	text = "bad"
	tokens := []int{1, 0, 3, -1}

	result, ast := testParser(g, 4, tokens, text, true, t)
	assert.True(t, result)
	choiceActual, s := ruleSeq(ast[0], "s")
	assert.Zero(t, choiceActual)
	assert.Len(t, s, 3)
}
