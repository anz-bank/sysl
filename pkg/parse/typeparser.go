//nolint:golint,stylecheck,interfacer
package parse

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/anz-bank/sysl/pkg/grammar"
	"github.com/anz-bank/sysl/pkg/sysl"
)

// typeParser is a simplified SyslParserListener which only handles enough rules to return a *sysl.Type
// This should be used by any rule which expects a type
type typeParser struct {
	parser.BaseSyslParserListener
	appName  []string
	typePath []string

	namestrTarget *[]string //which of the above slices should the Name_str be appended to
}

func (t *typeParser) run(syslContext *sysl.Scope, ctx parser.ITypesContext) *sysl.Type {
	t.appName = []string{}
	t.typePath = []string{}

	if ndt := ctx.(*parser.TypesContext).NativeDataTypes(); ndt != nil {
		result := &sysl.Type{}
		primType, constraints := primitiveFromNativeDataType(ndt)
		if primType != nil {
			result.Type = primType
			if constraints != nil {
				result.Constraint = []*sysl.Type_Constraint{constraints}
			}
		}
		return result
	}

	antlr.NewParseTreeWalker().Walk(t, ctx)

	ref := &sysl.Scope{
		Appname: &sysl.AppName{Part: t.appName},
		Path:    t.typePath,
	}
	if len(t.appName) == 0 {
		ref.Appname = nil
	}

	return &sysl.Type{
		Type: &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Context: syslContext,
				Ref:     ref,
			},
		},
	}
}

func (t *typeParser) EnterName_str(ctx *parser.Name_strContext) {
	*t.namestrTarget = append(*t.namestrTarget, ctx.GetText())
}

func (t *typeParser) EnterUser_defined_type(ctx *parser.User_defined_typeContext) {
	t.namestrTarget = &t.typePath
}

func (t *typeParser) EnterReference(ctx *parser.ReferenceContext) {
	t.namestrTarget = &t.appName
}

func (t *typeParser) ExitApp_name(ctx *parser.App_nameContext) {
	// This can only come from a reference, which means the remaining name_str's are the type path
	t.namestrTarget = &t.typePath
}

func (t *typeParser) doSet(syslContext *sysl.Scope, sc sourceCtxHelper, ctx *parser.Set_typeContext) *sysl.Type {
	child := parseAllowedTypedContexts(syslContext, sc, ctx.Types())
	parent := &sysl.Type{Type: &sysl.Type_Set{Set: child}}
	t.handleCommonCollection(parent, child, ctx.Size_spec())
	return parent
}

func (t *typeParser) doSequence(syslContext *sysl.Scope,
	sc sourceCtxHelper,
	ctx *parser.Sequence_typeContext,
) *sysl.Type {
	child := parseAllowedTypedContexts(syslContext, sc, ctx.Types())
	parent := &sysl.Type{Type: &sysl.Type_Sequence{Sequence: child}}
	t.handleCommonCollection(parent, child, ctx.Size_spec())
	return parent
}

func (t *typeParser) handleCommonCollection(parent, child *sysl.Type, sizeSpec parser.ISize_specContext) {
	parent.Attrs = child.Attrs
	parent.Opt = child.Opt
	parent.SourceContext = child.SourceContext

	child.Attrs = nil
	child.Opt = false
	child.SourceContext = nil

	if sizeSpec != nil {
		if child.GetPrimitive() != sysl.Type_NO_Primitive {
			spec := sizeSpec.(*parser.Size_specContext)
			child.Constraint = makeTypeConstraint(child.GetPrimitive(), spec)
		}
	}
}

func parseAllowedTypedContexts(syslContext *sysl.Scope,
	sc sourceCtxHelper,
	contexts ...antlr.ParserRuleContext,
) *sysl.Type {
	tb := typeParser{}
	for _, ctx := range contexts {
		switch x := ctx.(type) {
		case *parser.View_type_specContext:
			return parseAllowedTypedContexts(syslContext, sc, x.Types(), x.Collection_type())
		case *parser.Collection_typeContext:
			return parseAllowedTypedContexts(syslContext, sc, x.Set_type(), x.Sequence_type())
		case *parser.TypesContext:
			return tb.run(syslContext, x)
		case *parser.Set_typeContext:
			return tb.doSet(syslContext, sc, x)
		case *parser.Sequence_typeContext:
			return tb.doSequence(syslContext, sc, x)
		}
	}
	panic("Should never get here")
}
