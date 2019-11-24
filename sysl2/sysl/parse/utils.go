package parse

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	sysl "github.com/anz-bank/sysl/src/proto"
)

func makeExpr(ctx *sysl.SourceContext) *sysl.Expr {
	return &sysl.Expr{
		SourceContext: ctx,
	}
}

type sourceCtxHelper struct {
	filename string
}

func (s sourceCtxHelper) Get(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	return s.get(ctx, false)
}

func (s sourceCtxHelper) GetWithText(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	return s.get(ctx, true)
}

func (s sourceCtxHelper) get(ctx *antlr.BaseParserRuleContext, withText bool) *sysl.SourceContext {
	start := ctx.GetStart()
	end := ctx.GetStop()

	var text string
	if withText {
		is := start.GetInputStream()
		endChar := end.GetStop()
		if endChar == 0 {
			endChar = is.Size() - 1
		}
		text = is.GetText(start.GetStart(), endChar)
	}

	return &sysl.SourceContext{
		File: s.filename,
		Start: &sysl.SourceContext_Location{
			Line: int32(start.GetLine()),
			Col:  int32(start.GetColumn()),
		},
		End: &sysl.SourceContext_Location{
			Line: int32(end.GetLine()),
			Col:  int32(end.GetColumn()),
		},
		Text: text,
	}
}
