package parser

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNestedGrammar(t *testing.T) {
	text := `{ 123, {EXPR: 1 + 2 * 3 :} }`

	nested := makeNestedGrammarParser(text, makeRepeatSeq(makeQuantifierOptional()), makeEXPR())
	actual := []token{}

	nested.stack.push(nested.grammars["array"])

	for {
		tok := nested.getLexer().nextToken()
		if tok.id == -1 {
			p := nested.stack.top()
			tok = p.l.nextToken()
			if tok.id == -1 {
				break
			}
			nested.getLexer().currentIndex = p.l.currentIndex
		} else {
			switch tok.id {
			case 0:
				logrus.Println(tok.text)
				nested.pushGrammar(tok.text)
				tok.id = PUSH_GRAMMAR
			case 1:
				logrus.Println(tok.text)
				nested.popGrammar()
				tok.id = POP_GRAMMAR
			}
			p := nested.stack.top()
			p.l.currentIndex = nested.getLexer().currentIndex
		}
		actual = append(actual, tok)
	}
	logrus.Println(actual)
}
