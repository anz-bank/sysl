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
	first := leafs[0].Scanner()

	return &sysl.SourceContext{
		File:  first.Filename(),
		Start: buildSourceContextLocation(first),
	}
}

func buildSourceContextLocation(from parser.Scanner) *sysl.SourceContext_Location {
	line, col := from.Position()
	return &sysl.SourceContext_Location{
		Line: int32(line),
		Col:  int32(col),
	}
}
