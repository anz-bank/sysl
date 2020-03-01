package parser

import "github.com/antlr/antlr4/runtime/Go/antlr"

func NewThreadSafeSyslParser(input antlr.TokenStream) *SyslParser {
	p := NewSyslParser(input)

	// The generated parser code has the deserializer, ATN and DFA as
	// globals, and they are mutated. This makes the code non-thread-safe.
	// Rebuild these structures as local to the SyslParser instance.
	// This code is copied from the generated sysl_parser.go code.
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	p.Interpreter = antlr.NewParserATNSimulator(p, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	return p
}
