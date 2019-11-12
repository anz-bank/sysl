package parse

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	sysl "github.com/anz-bank/sysl/src/proto"
)

type sourceCtxHelper struct {
	filename string
}

func (s sourceCtxHelper) Get(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	start := ctx.GetStart()
	return &sysl.SourceContext{
		File: s.filename,
		Start: &sysl.SourceContext_Location{
			Line: int32(start.GetLine()),
			Col:  int32(start.GetColumn()),
		},
	}
}
