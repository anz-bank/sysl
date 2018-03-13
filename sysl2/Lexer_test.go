package main

import (
	"testing"
)

func TestLexer(t *testing.T) {
	content := `keywords:
        while | return |
        if | else | import`

	regexes := make([]string, 7)
	regexes[0] = `[|]`
	regexes[1] = `[:]`
	regexes[2] = `[+]`
	regexes[3] = `[*]`
	regexes[4] = `[?]`
	regexes[5] = `import`
	regexes[6] = `[a-z][a-z0-9]*`

	tokens := []int{6, 1, 6, 0, 6, 0, 6, 0, 6, 0, 5, -1}

	lex := makeLexer(content, regexes)
	for i, expected := range tokens {
		tok := lex.nextToken()
		if tok != expected {
			t.Errorf("wrong token @ %d: expected %d but got %d", i, expected, tok)
		}
	}
}
