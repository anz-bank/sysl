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

func buildTypeReference(syslContext *sysl.Scope, ctx parser.ITypesContext) *sysl.Type {
	tb := typeParser{}
	return tb.run(syslContext, ctx)
}
