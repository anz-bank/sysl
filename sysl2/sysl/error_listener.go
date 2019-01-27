package main

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// SyslParserErrorListener ...
type SyslParserErrorListener struct {
	*antlr.DefaultErrorListener
	hasErrors bool
}

// SyntaxError ...
func (d *SyslParserErrorListener) SyntaxError(
	recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	d.hasErrors = true
	fmt.Printf("SyntaxError: Token: %s\n", recognizer.GetSymbolicNames()[offendingSymbol.(*antlr.CommonToken).GetTokenType()])
}

// ReportAttemptingFullContext ...
func (d *SyslParserErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser,
	dfa *antlr.DFA, startIndex, stopIndex int,
	conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	fmt.Printf("ReportAttemptingFullContext: %d %d\n", startIndex, stopIndex)
}

// ReportAmbiguity ...
func (d *SyslParserErrorListener) ReportAmbiguity(recognizer antlr.Parser,
	dfa *antlr.DFA, startIndex, stopIndex int, exact bool,
	ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	fmt.Printf("ReportAmbiguity: %d %d\n", startIndex, stopIndex)
}

// ReportContextSensitivity ...
func (d *SyslParserErrorListener) ReportContextSensitivity(recognizer antlr.Parser,
	dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	fmt.Printf("ReportContextSensitivity: %d %d\n", startIndex, stopIndex)
}
