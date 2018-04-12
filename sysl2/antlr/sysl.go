package main

import (
	"fmt"
	"os"

	"anz-bank/sysl/src/proto"
	"anz-bank/sysl/sysl2/antlr/grammar"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang/protobuf/proto"
)

// TextPB ...
func TextPB(m *sysl.Module) {
	fmt.Println(proto.MarshalTextString(m))
}

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
}

// Parse ...
func Parse(filename string, root string) *sysl.Module {
	input := antlr.NewFileStream(filename)
	lexer := parser.NewSyslLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewSyslParser(stream)

	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	errorListener := SyslParserErrorListener{}
	p.AddErrorListener(&errorListener)
	p.BuildParseTrees = true

	tree := p.Sysl_file()
	if errorListener.hasErrors {
		return nil
	}
	listener := NewTreeShapeListener(root)

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
	return listener.module
}

func main() {
	// fmt.Printf("Reading file %s\n", os.Args[1])
	Parse(os.Args[1], "")
	// TextPB(mod)
}
