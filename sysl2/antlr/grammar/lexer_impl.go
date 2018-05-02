package parser

import (
	"fmt"
	s "strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func calcSpaces(text string) int {
	s := 0

	for i := 0; i < len(text); i++ {
		if text[i] == ' ' {
			s++
		}
		if text[i] == '\t' {
			s += 4
		}
	}
	return s
}

func startsWithKeyword(text string) bool {
	var lower = s.ToLower(text)
	// var first = s.Split(lower, " ")[0]

	for k := range keywords {
		if s.HasPrefix(lower, keywords[k]) {
			return true
		}
	}
	return false
}

func createDedentToken(source *antlr.TokenSourceCharStreamPair) *antlr.CommonToken {
	t := antlr.NewCommonToken(source, SyslLexerDEDENT, 0, 0, 0)

	return t
}

func createIndentToken(source *antlr.TokenSourceCharStreamPair) *antlr.CommonToken {
	t := antlr.NewCommonToken(source, SyslLexerINDENT, 0, 0, 0)

	return t
}

func GettokenName(t antlr.Token, l *SyslLexer) string {
	if t.GetTokenType() > 0 && t.GetTokenType() < len(l.SymbolicNames) {
		return l.SymbolicNames[t.GetTokenType()]
	}
	if t.GetTokenType() == antlr.TokenEOF {
		return "EOF"
	}

	return fmt.Sprintf("%d", t.GetTokenType())
}

var keywords = [...]string{"set of", "return", "for", "one of", "else", "if", "loop", "until"}
var prevToken []antlr.Token

type Stack struct {
	stack []int
	index int
}

func NewStack() *Stack {
	s := new(Stack)
	s.index = 0
	return s
}

func (s *Stack) Push(o int) {
	s.stack = append(s.stack, o)
	s.index++
}

func (s *Stack) Pop() int {
	if s.index < 0 {
		panic("empty stack")
	}
	l := len(s.stack)
	ret := s.stack[l-1]
	s.stack = s.stack[:l-1]
	s.index--
	return ret
}

func (s *Stack) Size() int {
	return s.index
}

func (s *Stack) Peek() int {
	return s.stack[s.index-1]
}

var level = NewStack()

func getPreviousIndent(l *Stack) int {
	if level.Size() == 0 {
		return 0
	}
	// peek, read but not remove HEAD
	return l.Peek()
}

// GetNextToken ...
func GetNextToken(l *SyslLexer) antlr.Token {
	if len(prevToken) > 0 {
		// poll, retrieve head
		nextTok := prevToken[0]
		prevToken = prevToken[1:]
		// tokenName := GettokenName(nextTok, l)
		return nextTok
	}

	next := l.BaseLexer.NextToken()
	// tokenName := GettokenName(next, l)
	// return NEWLINE
	if gotNewLine && (next.GetTokenType() == SyslLexerNEWLINE || next.GetTokenType() == SyslLexerNEWLINE_2 || next.GetTokenType() == SyslLexerEMPTY_LINE || next.GetTokenType() == SyslLexerINDENTED_COMMENT || next.GetTokenType() == SyslLexerEMPTY_COMMENT) {
		return next
	}
	// regular whitespace, return as is.
	// return from here only when we encounter HIDDEN after INDENT has been generated
	// after processing NL.
	if !gotNewLine && next.GetChannel() == antlr.TokenHiddenChannel {
		spaces = 0
		return next
	} else if next.GetTokenType() == SyslLexerSYSL_COMMENT {
		spaces = 0
		return next
	}

	if next.GetTokenType() == antlr.TokenEOF {
		spaces = 0 // done with the file
	} else if !gotNewLine {
		return next
	}

	// if spaces != getPreviousIndent(level) {
	// 	// System.out.println(linenum + " :will_que " + tokenName)
	// 	fmt.Println(fmt.Sprintf("%d :will_que - %s ", linenum, tokenName))
	// }

	for spaces != getPreviousIndent(level) {
		if spaces > getPreviousIndent(level) {
			level.Push(spaces)
			prevToken = append(prevToken, createIndentToken(next.GetSource()))
		} else {
			level.Pop()
			prevToken = append(prevToken, createDedentToken(next.GetSource()))
		}
	}

	gotNewLine = false
	prevToken = append(prevToken, next)
	// poll, retrieve head
	temp := prevToken[0]
	prevToken = prevToken[1:]
	return temp
}
