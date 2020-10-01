package syslutil

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

// GetTypeDetail returns name of the type and details in string format
func GetTypeDetail(t *sysl.Type) (typeName string, typeDetail string) {
	switch x := t.Type.(type) {
	case *sysl.Type_Primitive_:
		typeName = "primitive"
		typeDetail = sysl.Type_Primitive_name[int32(x.Primitive)]
	case *sysl.Type_TypeRef:
		typeName = "type_ref"
		if x.TypeRef.Ref != nil {
			if len(x.TypeRef.Ref.Path) == 1 {
				typeDetail = x.TypeRef.Ref.Path[0]
				if x.TypeRef.Ref.Appname != nil {
					typeDetail = GetAppName(x.TypeRef.Ref.Appname) + "." + typeDetail
				}
			} else if x.TypeRef.Ref.Appname != nil && len(x.TypeRef.Ref.Appname.Part) > 0 {
				typeDetail = GetAppName(x.TypeRef.Ref.Appname)
			}
		}
	case *sysl.Type_Sequence:
		typeName = "sequence"
		_, d := GetTypeDetail(x.Sequence)
		typeDetail = d
	case *sysl.Type_Set:
		typeName = "set"
		_, d := GetTypeDetail(x.Set)
		typeDetail = d
	case *sysl.Type_List_:
		typeName = "list"
		_, d := GetTypeDetail(x.List.Type)
		typeDetail = d
	case *sysl.Type_Tuple_:
		typeName = "tuple"
	case *sysl.Type_Relation_:
		typeName = "relation"
	case *sysl.Type_OneOf_:
		typeName = "union"
	}
	return typeName, typeDetail
}

// TypeNone returns none-type
func TypeNone() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_NoType_{NoType: &sysl.Type_NoType{}}}
}

// TypeEmpty returns empty primitive type
func TypeEmpty() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_EMPTY}}
}

// TypeString returns string type
func TypeString() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}}
}

// TypeInt returns integer type
func TypeInt() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_INT}}
}

// TypeFloat returns float type
func TypeFloat() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_FLOAT}}
}

// TypeDecimal returns decimal type
func TypeDecimal() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_DECIMAL}}
}

// TypeBool returns boolean type
func TypeBool() *sysl.Type {
	return &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_BOOL}}
}

// HasSameType returns true if type 2 matches with type 1
func HasSameType(type1 *sysl.Type, type2 *sysl.Type) bool {
	if type1 == nil || type2 == nil {
		return false
	}

	switch type1.GetType().(type) {
	case *sysl.Type_Primitive_:
		return type1.GetPrimitive() == type2.GetPrimitive()
	case *sysl.Type_TypeRef:
		if type2.GetTypeRef() != nil {
			ref1 := type1.GetTypeRef().GetRef()
			ref2 := type2.GetTypeRef().GetRef()

			if ref1.GetAppname() != nil && ref2.GetAppname() != nil {
				return ref1.GetAppname().GetPart()[0] == ref2.GetAppname().GetPart()[0]
			} else if ref1.GetPath() != nil && ref2.GetPath() != nil {
				return ref1.GetPath()[0] == ref2.GetPath()[0]
			}
		}
	case *sysl.Type_Tuple_:
		return type2.GetTuple() != nil
	}

	return false
}

func JoinAppName(name *sysl.AppName) string {
	return strings.Join(name.GetPart(), " :: ")
}

func JoinTypePath(path []string) string {
	return strings.Join(path, ".")
}
