package syslutil

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/utils"
)

const namespaceSeparator = " :: "

// GetTypeDetail returns name of the type and details in string format
func GetTypeDetail(t *sysl.Type) (typeName string, typeDetail string) {
	switch x := t.Type.(type) {
	case *sysl.Type_Primitive_:
		typeName = "primitive"
		typeDetail = sysl.Type_Primitive_name[int32(x.Primitive)]
	case *sysl.Type_TypeRef:
		typeName = "type_ref"
		typeDetail = JoinTypeRef(x.TypeRef)
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

type BuiltInType = string

// BuiltInTypes are the builtin types of Sysl.
// There's no primitive int32 (and int64) type.
// The parser interprets the Sysl syntax int32 value
// as an int primitive with a bit width (sysl.Type_Constraint.BitWidth) of 32.
var BuiltInTypes = []BuiltInType{
	Type_NO_Primitive,
	Type_EMPTY,
	Type_ANY,
	Type_BOOL,
	Type_INT,
	Type_INT32,
	Type_INT64,
	Type_FLOAT,
	Type_DECIMAL,
	Type_STRING,
	Type_BYTES,
	Type_STRING_8,
	Type_DATE,
	Type_DATETIME,
	Type_XML,
	Type_UUID,
}

// nolint:golint,stylecheck
var (
	Type_NO_Primitive BuiltInType = GetTypePrimitiveName(sysl.Type_NO_Primitive)
	Type_EMPTY        BuiltInType = GetTypePrimitiveName(sysl.Type_EMPTY)
	Type_ANY          BuiltInType = GetTypePrimitiveName(sysl.Type_ANY)
	Type_BOOL         BuiltInType = GetTypePrimitiveName(sysl.Type_BOOL)
	Type_INT          BuiltInType = GetTypePrimitiveName(sysl.Type_INT)
	Type_INT32        BuiltInType = "int32"
	Type_INT64        BuiltInType = "int64"
	Type_FLOAT        BuiltInType = GetTypePrimitiveName(sysl.Type_FLOAT)
	Type_DECIMAL      BuiltInType = GetTypePrimitiveName(sysl.Type_DECIMAL)
	Type_STRING       BuiltInType = GetTypePrimitiveName(sysl.Type_STRING)
	Type_BYTES        BuiltInType = GetTypePrimitiveName(sysl.Type_BYTES)
	Type_STRING_8     BuiltInType = GetTypePrimitiveName(sysl.Type_STRING_8)
	Type_DATE         BuiltInType = GetTypePrimitiveName(sysl.Type_DATE)
	Type_DATETIME     BuiltInType = GetTypePrimitiveName(sysl.Type_DATETIME)
	Type_XML          BuiltInType = GetTypePrimitiveName(sysl.Type_XML)
	Type_UUID         BuiltInType = GetTypePrimitiveName(sysl.Type_UUID)
)

func IsBuiltIn(name string) bool {
	return utils.Contains(name, BuiltInTypes)
}

// nolint:interfacer
func GetTypePrimitiveName(t sysl.Type_Primitive) string {
	return strings.ToLower(t.String())
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

// JoinAppName returns a string with the parts of an app name joined on the namespace separator.
func JoinAppName(name *sysl.AppName) string {
	return JoinAppNameParts(name.GetPart()...)
}

// JoinAppNameParts returns a string with the parts of an app name joined on the namespace
// separator.
func JoinAppNameParts(parts ...string) string {
	return strings.Join(parts, namespaceSeparator)
}

// SplitAppNameParts returns the parts of an app name that was previously joined into a string.
func SplitAppNameParts(appName string) []string {
	parts := strings.Split(appName, strings.TrimSpace(namespaceSeparator))
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// JoinTypePath returns a string with the parts of a type reference path joined on ".".
func JoinTypePath(path []string) string {
	return strings.Join(path, ".")
}

// JoinTypeRef returns a string encoding the reference to a type.
//
// If the app name is present it will be joined on " :: " and prefixed to the path, joined by ".".
func JoinTypeRef(ref *sysl.ScopedRef) string {
	return JoinTypeRefScope(ref.Ref)
}

// JoinTypeRefScope returns a string encoding the reference to a type.
//
// If the app name is present it will be joined on " :: " and prefixed to the path, joined by ".".
func JoinTypeRefScope(ref *sysl.Scope) string {
	if ref.Appname != nil && len(ref.Appname.Part) > 0 {
		return JoinTypePath(append([]string{JoinAppName(ref.Appname)}, ref.Path...))
	}
	return JoinTypePath(ref.Path)
}
