package golang

import (
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
)

// GoPackageName derives the Go package name from an Application's go_package
// attribute.
func GoPackageName(app *sysl.Application) string {
	goPackage := app.Attrs["go_package"].Attribute.(*sysl.Attribute_S).S
	parts := strings.SplitN(goPackage, ";", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	importPathParts := strings.Split(parts[0], "/")
	return importPathParts[len(importPathParts)-1]
}

// GoImportPath derives the Go import path from an Application's go_package
// attribute.
func GoImportPath(app *sysl.Application) string {
	goPackage := app.Attrs["go_package"].Attribute.(*sysl.Attribute_S).S
	return strings.SplitN(goPackage, ";", 2)[0]
}

var primitiveTypes = map[sysl.Type_Primitive]Expr{
	sysl.Type_BOOL:     I("bool"),
	sysl.Type_INT:      I("int64"),
	sysl.Type_FLOAT:    I("float64"),
	sysl.Type_DECIMAL:  I("float64"),
	sysl.Type_STRING:   SliceType(I("rune")),
	sysl.Type_BYTES:    SliceType(I("byte")),
	sysl.Type_STRING_8: I("string"),
	sysl.Type_DATE:     Dot(I("time"), "Time"),
	sysl.Type_DATETIME: Dot(I("time"), "Time"),
	sysl.Type_XML:      I("string"),
	sysl.Type_UUID:     Dot(I("uuid"), "UUID"),
}

// StandardImports carries the imports required by internal code generation.
var StandardImports = []ImportSpec{
	{Path: *String("time")},
	{Path: *String("github.com/google/uuid")},
}

// GoTypeForSyslType returns the equivalent Go type for a Sysl type.
func GoTypeForSyslType(syslType *sysl.Type, tname string, app *sysl.Application) Expr {
	e, starable := func() (Expr, bool) {
		switch syslType := syslType.Type.(type) {
		case *sysl.Type_OneOf_:
			return Interface(
				*NewField(UnionTagMethodType(tname), "IsA"),
			), false
		case *sysl.Type_Primitive_:
			if goType, found := primitiveTypes[syslType.Primitive]; found {
				return goType, true
			}
		case *sysl.Type_Sequence:
			return SliceType(GoTypeForSyslType(syslType.Sequence, "", app)), true
		case *sysl.Type_Tuple_:
			fields := []Field{}
			for _, fname := range sysl.TypeNamesInSourceOrder(syslType.Tuple.AttrDefs) {
				ftype := syslType.Tuple.AttrDefs[fname]
				fields = append(fields, *NewField(GoTypeForSyslType(ftype, "", app), Public(fname)))
			}
			return Struct(fields...), true
		case *sysl.Type_TypeRef:
			ref := syslType.TypeRef.Ref.Path[0]
			referee, found := app.Types[ref]
			if !found {
				panic(errors.Errorf("Reference to unkown type: %s", ref))
			}
			_, starable := referee.Type.(*sysl.Type_Tuple_)
			return I(ref), starable
		}
		panic(errors.Errorf("Unhandled type: %#v", syslType))
	}()
	if starable && syslType.Opt {
		return Star(e)
	}
	return e

}

// UnionTagMethodType computes a FuncType for a union's tag method.
func UnionTagMethodType(tname string) *FuncType {
	return NewFuncType(ParenFields(*NewField(I(tname))))
}
