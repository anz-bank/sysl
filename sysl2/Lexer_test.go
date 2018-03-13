package main

import (
	"testing"
)

func TestLexer(t *testing.T) {
	content := `keywords: 
				while | return | 
				if | else | import`

	regexs := make([]string, 7)
	regexs[0] = `[|]`
	regexs[1] = `[:]`
	regexs[2] = `[+]`
	regexs[3] = `[*]`
	regexs[4] = `[?]`
	regexs[5] = `import`
	regexs[6] = `[a-z][a-z0-9]*`

	tokens := []int{6, 1, 6, 0, 6, 0, 6, 0, 6, 0, 5, -1}

	lex := makeLexer(content, regexs)
	for i := range tokens {
		tok := lex.nextToken()
		if tok != tokens[i] {
			t.Errorf("got the wrong token: %d", i)
		}
	}
}
