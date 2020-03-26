package parser

import (
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

func buildSourceContext(node ast.Node) *sysl.SourceContext {
	leafs := Leafs(node)
	if len(leafs) < 1 {
		return &sysl.SourceContext{}
	}

	// FIXME
	return &sysl.SourceContext{}

	first := leafs[0].Scanner()
	end := leafs[len(leafs)-1].Scanner()

Loop:
	for _, l := range leafs[1:] {
		val := l.Scanner().String()
		switch val {
		case ":", "\n", "\r\n":
			end = l.Scanner()
			break Loop
		}
	}

	return &sysl.SourceContext{
		File:  first.Filename(),
		Start: buildSourceContextLocation(first),
		End:   buildSourceContextLocation(end),
	}
}

func buildSourceContextLocation(from parser.Scanner) *sysl.SourceContext_Location {
	line, col := from.Position()
	return &sysl.SourceContext_Location{
		Line: int32(line),
		Col:  int32(col),
	}
}
