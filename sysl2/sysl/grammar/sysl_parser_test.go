package parser

import (
	"fmt"
	"testing"

	parser "github.com/anz-bank/sysl/sysl2/naive"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func TestSyslLexer(t *testing.T) {
	input := `
# Example template file demonstrating how to use the template views
# -> Will dump the input model to a simplified text file


DependencyTemplate:

    !view findEpDeps(ep <: sysl.Endpoint) -> set of string:
        ep -> (s:
            let targets = ep.value.stmts where(.type == "call") -> <set of string> (t:
                arg = '"' + t.target + '"'
            )
            $  {$ targets $}
            out = targets flatten(.arg)
        )

    !view dump(app <: sysl.App) -> string:
        app -> (:
            let targets = app.endpoints -> <set of string> (ep:
                out = findEpDeps(ep)
            )

            $ {{ "{$ app.name $}": {$ targets flatten(.out) flatten(.out) flatten(.) $} }}\
        )

    !view start(mod <: sysl.TemplateInput) -> sysl.TemplateResult:
        mod -> (:
            Filename = "foo.json"
            let apps = mod.Apps -> <sequence of string> (a:
                out = dump(a)
            )
            Data = "[\n" + apps + "\n]"
        )


`

	lexer := NewSyslLexer(antlr.NewInputStream(input))
	defer DeleteLexerState(lexer)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	p := NewSyslParser(stream)
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	//	p.AddErrorListener(&errorListener)

	p.BuildParseTrees = true

	p.Sysl_file()
	for _, token := range stream.GetAllTokens() {
		if token.GetTokenType() == parser.EOF {
			fmt.Printf("<<EOF>>\n")
			continue
		}

		text := lexerSymbolicNames[token.GetTokenType()]
		if text == "" {
			text = lexerLiteralNames[token.GetTokenType()]
		}
		fmt.Printf("% 3d  %-20s '%s'\n", token.GetTokenType(), text, token.GetText())
	}
}
