package parser

import "github.com/antlr/antlr4/runtime/Go/antlr"

func NewThreadSafeSyslLexer(input antlr.CharStream) *SyslLexer {
	l := NewSyslLexer(input)

	// The generated lexer code has the deserializer, ATN and DFA as
	// globals, and they are mutated. This makes the code non-thread-safe.
	// Rebuild these structures as local to the SyslLexer instance.
	// This code is copied from the generated sysl_lexer.go code.
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}

	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())
	return l
}
