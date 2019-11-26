package parser

import (
	"os"
	s "strings"
	"unsafe"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/cornelk/hashmap"
	"github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var (
	syslLexerLog = os.Getenv("SYSL_LEXER_LOG") != ""

	keywords = [...]string{
		"sequence of",
		"set of",
		"return",
		"for",
		"one of",
		"else",
		"if",
		"loop",
		"until",
		"alt",
		"while",
	}

	// Antlr doesn't support reentrant Go lexer state, so we work around it with
	// a fast lock-free hash map.
	lexerStates = &hashmap.HashMap{}
)

const importKeyword = "import"

type lexerState struct {
	prevToken []antlr.Token
	level     stack

	spaces        int
	linenum       int
	inSqBrackets  int
	parens        int
	blockTextLine int
	gotNewLine    bool
	gotHTTPVerb   bool
	gotView       bool
	noMoreImports bool // Used to allow the import keyword after the application definition has started
}

func ls(l *SyslLexer) *lexerState {
	key := uintptr(unsafe.Pointer(l))
	if state, has := lexerStates.Get(key); has {
		return state.(*lexerState)
	}
	state := &lexerState{}
	lexerStates.Set(key, state)
	return state
}

func DeleteLexerState(l *SyslLexer) {
	key := uintptr(unsafe.Pointer(l))
	lexerStates.Del(key)
}

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

func startsWithKeyword(l *lexerState, text string) bool {
	var lower = s.ToLower(text)

	for _, kw := range keywords {
		if s.HasPrefix(lower, kw) {
			return true
		}
	}

	// import is only a keyword before the application starts
	isImport := s.HasPrefix(text, importKeyword) && !l.noMoreImports

	return isImport
}

func createDedentToken(source *antlr.TokenSourceCharStreamPair) *antlr.CommonToken {
	return antlr.NewCommonToken(source, SyslLexerDEDENT, 0, 0, 0)
}

func createIndentToken(source *antlr.TokenSourceCharStreamPair) *antlr.CommonToken {
	return antlr.NewCommonToken(source, SyslLexerINDENT, 0, 0, 0)
}

type stack []int

func (s *stack) Push(o int) {
	*s = append(*s, o)
}

func (s *stack) Pop() int {
	l := len(*s)
	ret := (*s)[l-1]
	*s = (*s)[:l-1]
	return ret
}

func (s *stack) Size() int {
	return len(*s)
}

func (s *stack) Peek() int {
	return (*s)[len(*s)-1]
}

func getPreviousIndent(s stack) int {
	if s.Size() == 0 {
		return 0
	}
	// peek, read but not remove HEAD
	return s.Peek()
}

// trimText Token Text
func trimText(l *SyslLexer) string {
	return s.TrimSpace(l.GetText())
}

func getNextToken(l *SyslLexer) antlr.Token {
	ls := ls(l)
	if len(ls.prevToken) > 0 {
		// poll, retrieve head
		nextTok := ls.prevToken[0]
		ls.prevToken = ls.prevToken[1:]
		return nextTok
	}

	next := l.BaseLexer.NextToken()
	if syslLexerLog {
		logrus.Info(next)
	}
	// return NEWLINE
	if ls.gotNewLine {
		switch next.GetTokenType() {
		case SyslLexerNEWLINE, SyslLexerNEWLINE_2, SyslLexerEMPTY_LINE, SyslLexerE_NL, SyslLexerE_EMPTY_LINE:
			fallthrough
		case SyslLexerINDENTED_COMMENT, SyslLexerEMPTY_COMMENT, SyslLexerE_INDENTED_COMMENT:
			fallthrough
		case SyslLexerE_DOT_NAME_NL:
			return next
		}
	}
	// regular whitespace, return as is.
	// return from here only when we encounter HIDDEN after INDENT has been generated
	// after processing NL.
	if !ls.gotNewLine && next.GetChannel() == antlr.TokenHiddenChannel {
		ls.spaces = 0
		return next
	} else if next.GetTokenType() == SyslLexerSYSL_COMMENT {
		ls.spaces = 0
		return next
	}

	if next.GetTokenType() == antlr.TokenEOF {
		ls.spaces = 0 // done with the file
	} else if !ls.gotNewLine {
		return next
	}

	for ls.spaces != getPreviousIndent(ls.level) {
		if ls.spaces > getPreviousIndent(ls.level) {
			ls.level.Push(ls.spaces)
			ls.prevToken = append(ls.prevToken, createIndentToken(next.GetSource()))
		} else {
			ls.level.Pop()
			ls.prevToken = append(ls.prevToken, createDedentToken(next.GetSource()))
		}
	}

	ls.gotNewLine = false
	ls.prevToken = append(ls.prevToken, next)
	// poll, retrieve head
	temp := ls.prevToken[0]
	ls.prevToken = ls.prevToken[1:]
	return temp
}
