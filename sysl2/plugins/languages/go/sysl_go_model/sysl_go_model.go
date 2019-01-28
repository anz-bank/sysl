package main

import (
	"context"

	"github.com/anz-bank/sysl/src/proto"
	g "github.com/anz-bank/sysl/sysl2/codegen/golang"
	"github.com/anz-bank/sysl/sysl2/plugins"
	"github.com/anz-bank/sysl/sysl2/plugins/languages/go"
)

type modelPlugin struct{}

func (*modelPlugin) GenerateCode(
	ctx context.Context,
	request *plugins.GenerateCodeRequest,
) (*plugins.GenerateCodeResponse, error) {
	m := request.Module
	app := m.Apps[request.Parameter]

	decls := []g.Decl{}
	for _, tname := range sysl.TypeNamesInSourceOrder(app.Types) {
		t := app.Types[tname]
		goType := g.Public(tname)

		decls = append(decls,
			g.TypeDecl(
				g.NewTypeSpec(goType, g.GoTypeForSyslType(t, goType, app)),
			).WithDocf("represents the Sysl-defined %s type.", tname),
		)

		if oneof, ok := t.Type.(*sysl.Type_OneOf_); ok {
			for _, ftype := range sysl.TypesInSourceOrder(oneof.OneOf.Type) {
				ref := ftype.Type.(*sysl.Type_TypeRef).TypeRef.Ref.Path[0]
				tagMethodType := *g.UnionTagMethodType(goType)
				decls = append(decls,
					g.Method("", g.Star(g.I(ref)), "IsA", tagMethodType).
						WithDocf("tags %s as implementing %s.", ref, goType),
				)
			}
		}
	}

	return helpers.SingleGoFileGenerateCodeResponse(g.File{
		Name:    *g.I(g.GoPackageName(app)),
		Imports: append(g.StandardImports, g.ImportSpec{Path: *g.String("fmt")}),
		Decls:   decls,
	})
}

func main() {
	plugins.ServeCodeGenerator(&modelPlugin{})
}
