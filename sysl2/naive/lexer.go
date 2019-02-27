package parser

import (
	"regexp"
	"sort"

	"github.com/anz-bank/sysl/sysl2/proto"
	"github.com/sirupsen/logrus"
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

type terminalBuilder struct {
	arr    []string
	tokens map[string]int32
	index  int32
}

func (b *terminalBuilder) buildFromChoice(choice *sysl.Choice) {
	for _, s := range choice.Sequence {
		for _, t := range s.Term {
			if t == nil {
				continue
			}
			var str string
			a := t.Atom
			switch x := a.Union.(type) {
			case *sysl.Atom_String_:
				str = x.String_
			case *sysl.Atom_Regexp:
				str = x.Regexp
			case *sysl.Atom_Choices:
				b.buildFromChoice(x.Choices)
			}
			if str != "" {
				if _, has := b.tokens[str]; !has {
					b.tokens[str] = b.index
					if len(str) == 1 {
						str = "[" + str + "]"
					}
					b.arr = append(b.arr, str)
					a.Id = b.index
					b.index++
				}
				logrus.Printf("token: [%s] (id=%d)\n", str, a.Id)
			}
		}
	}
}

// assigns new value to Atom.Id
func getTerminals(rules map[string]*sysl.Rule) []string {
	builder := terminalBuilder{
		arr:    []string{},
		tokens: map[string]int32{},
		index:  0,
	}

	var ks []string
	for key := range rules {
		ks = append(ks, key)
	}
	sort.Strings(ks)
	for _, key := range ks {
		logrus.Println("Key: " + key)
		b := rules[key]
		builder.buildFromChoice(b.Choices)
	}

	return builder.arr
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
