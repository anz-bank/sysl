package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Parallel()

	content := `keywords:
        while | return |
        if | else | import`

	regexes := make([]string, 7)
	// Initialzing on separate lines helps in figuring out the the id of each token.
	// id is the index in the regexes array.
	regexes[0] = `[|]`
	regexes[1] = `[:]`
	regexes[2] = `[+]`
	regexes[3] = `[*]`
	regexes[4] = `[?]`
	regexes[5] = `import`
	regexes[6] = `[a-z][a-z0-9]*`

	tokens := []token{
		{6, "keywords"},
		{1, ":"},
		{6, "while"},
		{0, "|"},
		{6, "return"},
		{0, "|"},
		{6, "if"},
		{0, "|"},
		{6, "else"},
		{0, "|"},
		{5, "import"},
		{-1, ""},
	}

	lex := makeLexer(content, regexes)
	for i, expected := range tokens {
		tok := lex.nextToken()
		assert.Equalf(t, expected, tok, "wrong token @ %d: expected %+v but got %+v", i, expected, tok)
	}
}
