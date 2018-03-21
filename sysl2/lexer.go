package main

import (
	"fmt"
    "regexp"
	"sort"

	"github.com/anz-bank/sysl/sysl2/proto"
)

type token struct {
    id   int
    text string
}

type lexer struct {
    currentIndex        int
    regexs              []*regexp.Regexp
    content             string
	WS               *regexp.Regexp
    ignoreWhiteSpace    bool
}

var arr []string
var tokens map[string]int32
var index int32

func buildFromChoice(choice *sysl.Choice) {

	for _, s := range choice.Sequence {

		for _, t := range s.Term {
			if t == nil {
				continue
			}
			var str string
			a := t.Atom
			switch a.Union.(type) {
			case *sysl.Atom_String_:
				str = t.GetAtom().GetString_()
			case *sysl.Atom_Regexp:
				str = t.GetAtom().GetRegexp()
			case *sysl.Atom_Choices:
				buildFromChoice(t.GetAtom().GetChoices())
			}
			if str != "" {
				i, has := tokens[str]
				if !has {
					fmt.Println("token: [" + str + "]")
					tokens[str] = index
					if len(str) == 1 {
						str = "[" + str + "]"
					}
					arr = append(arr, str)
					a.Id = index
					index++
				} else {
					fmt.Printf("token: [%s] = %d (id=%d)\n", str, i, a.Id)

				}
			}
		}
	}
}

// assigns new value to Atom.Id
func getTerminals(rules map[string]*sysl.Rule) []string {
	arr = make([]string, 0)
	tokens = make(map[string]int32)
	index = 0

	var ks []string
	for key := range rules {
		ks = append(ks, key)
	}
	sort.Strings(ks)
	for _, key := range ks {
		fmt.Println("Key: " + key)
		b := rules[key]
		buildFromChoice(b.Choices)
	}

	return arr
}

// Regular whitespace delimited
func makeLexer(content string, regexs []string) lexer {
    compiledRegexes := make([]*regexp.Regexp, len(regexs))

    for i, reg := range regexs {
        compiledRegex := regexp.MustCompile(reg)
        compiledRegexes[i] = compiledRegex
    }

    return lexer{
        ignoreWhiteSpace:    true,
        regexs:              compiledRegexes,
        currentIndex:        0,
        content:             content,
		WS:               regexp.MustCompile(`[^ \t\r\n]+`),
    }
}

func (l *lexer) nextToken() token {
    longestMatchIndex := -1
    var tokenString string
    if l.currentIndex < len(l.content) {
        longestMatchLength := 0
		locWS := l.WS.FindStringIndex(l.content[l.currentIndex:])
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
