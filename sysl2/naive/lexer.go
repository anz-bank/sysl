package parser

import (
	"regexp"
)

type token struct {
	id   int
	text string
}

type lexer struct {
	currentIndex     int
	regexs           []*regexp.Regexp
	content          string
	ws               *regexp.Regexp
	ignoreWhiteSpace bool
}

// Regular whitespace delimited
func makeLexer(content string, regexs []string) lexer {
	compiledRegexes := make([]*regexp.Regexp, len(regexs))

	for i, reg := range regexs {
		compiledRegex := regexp.MustCompile(reg)
		compiledRegexes[i] = compiledRegex
	}

	return lexer{
		ignoreWhiteSpace: true,
		regexs:           compiledRegexes,
		currentIndex:     0,
		content:          content,
		ws:               regexp.MustCompile(`[^ \t\r\n]+`),
	}
}

func (l *lexer) nextToken() token {
	longestMatchIndex := -1
	var tokenString string
	if l.currentIndex < len(l.content) {
		longestMatchLength := 0
		locWS := l.ws.FindStringIndex(l.content[l.currentIndex:])
		if locWS != nil {
			start := l.currentIndex + locWS[0]
			end := l.currentIndex + locWS[1]
			str := l.content[start:end]

			for i, r := range l.regexs {
				loc := r.FindStringIndex(str)
				if loc != nil {
					matchLen := loc[1] - loc[0]
					if loc[0] == 0 && matchLen > longestMatchLength {
						longestMatchLength = loc[1]
						longestMatchIndex = i
						tokenString = str[loc[0]:loc[1]]
					}
				}
			}
			if longestMatchIndex != -1 {
				l.currentIndex = (start + longestMatchLength)
			}
		}
	}
	return token{longestMatchIndex, tokenString}
}
