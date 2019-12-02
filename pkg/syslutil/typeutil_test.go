package syslutil

import (
	"testing"

	sysl "github.com/anz-bank/sysl/pkg/proto_old"
	"github.com/stretchr/testify/assert"
)

func TestHasSameType(t *testing.T) {
	t.Parallel()

	type inputData struct {
		type1 *sysl.Type
		type2 *sysl.Type
	}
	cases := map[string]struct {
		input    inputData
		expected bool
	}{
		"Same primitive types": {
			input:    inputData{type1: TypeString(), type2: TypeString()},
			expected: true},
		"Different primitive types1": {
			input:    inputData{type1: TypeString(), type2: TypeInt()},
			expected: false},
		"Different primitive types2": {
			input:    inputData{type1: TypeInt(), type2: TypeString()},
			expected: false},
		"Same transform typerefs1": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"Statement"}}}}},
				type2: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"Statement"}}}}}},
			expected: true},
		"Different transform typerefs1-1": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"Statement"}}}}},
				type2: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"StatementList"}}}}}},
			expected: false},
		"Different transform typerefs1-2": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"StatementList"}}}}},
				type2: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"Statement"}}}}}},
			expected: false},
		"Same transform typerefs2": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}},
				type2: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}}},
			expected: true},
		"Different transform typerefs2-1": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}},
				type2: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"StatementList"}}}}}}},
			expected: false},
		"Different transform typerefs2-2": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"StatementList"}}}}}},
				type2: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}}},
			expected: false},
		"Different types1": {
			input:    inputData{type1: TypeNone(), type2: TypeString()},
			expected: false},
		"Different types2": {
			input:    inputData{type1: TypeString(), type2: TypeNone()},
			expected: false},
		"Different types3": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}},
				type2: TypeString()},
			expected: false},
		"Different types3.5": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}},
				type2: TypeString()},
			expected: false},
		"Different types4": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"StatementList"}}}}},
				type2: TypeString()},
			expected: false},
		"Tuples": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_Tuple_{Tuple: &sysl.Type_Tuple{}}},
				type2: &sysl.Type{
					Type: &sysl.Type_Tuple_{Tuple: &sysl.Type_Tuple{}}}},
			expected: true},
		"Nil types": {
			input:    inputData{type1: nil, type2: nil},
			expected: false},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			isSame := HasSameType(input.type1, input.type2)
			assert.True(t, expected == isSame, "Unexpected result")
		})
	}
}

func TestGetTypeDetail(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input              *sysl.Type
		expectedTypeName   string
		expectedTypeDetail string
	}{
		"primitive": {input: TypeInt(), expectedTypeName: "primitive", expectedTypeDetail: "INT"},
		"type_ref path": {input: &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Ref: &sysl.Scope{Path: []string{"foo"}}}}},
			expectedTypeName:   "type_ref",
			expectedTypeDetail: "foo"},
		"type_ref appname": {input: &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Ref: &sysl.Scope{
						Path:    []string{"foo", ""},
						Appname: &sysl.AppName{Part: []string{"bar"}}}}}},
			expectedTypeName:   "type_ref",
			expectedTypeDetail: "bar"},
		"sequence": {input: &sysl.Type{
			Type: &sysl.Type_Sequence{
				Sequence: TypeString()}},
			expectedTypeName:   "sequence",
			expectedTypeDetail: "STRING"},
		"set": {input: &sysl.Type{
			Type: &sysl.Type_Set{
				Set: TypeString()}},
			expectedTypeName:   "set",
			expectedTypeDetail: "STRING"},
		"list": {input: &sysl.Type{
			Type: &sysl.Type_List_{
				List: &sysl.Type_List{
					Type: TypeFloat()}}},
			expectedTypeName:   "list",
			expectedTypeDetail: "FLOAT"},
		"tuple": {input: &sysl.Type{
			Type: &sysl.Type_Tuple_{}},
			expectedTypeName:   "tuple",
			expectedTypeDetail: ""},
		"relation": {input: &sysl.Type{
			Type: &sysl.Type_Relation_{}},
			expectedTypeName:   "relation",
			expectedTypeDetail: ""},
		"union": {input: &sysl.Type{
			Type: &sysl.Type_OneOf_{}},
			expectedTypeName:   "union",
			expectedTypeDetail: ""},
	}

	for name, test := range cases {
		input := test.input
		expectedTypeName := test.expectedTypeName
		expectedTypeDetail := test.expectedTypeDetail
		t.Run(name, func(t *testing.T) {
			typeName, typeDetail := GetTypeDetail(input)
			assert.Equal(t, expectedTypeName, typeName)
			assert.Equal(t, expectedTypeDetail, typeDetail)
		})
	}
}

func TestTypeNone(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_NoType_{NoType: &sysl.Type_NoType{}}}, TypeNone())
}

func TestTypeEmpty(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_EMPTY}}, TypeEmpty())
}

func TestTypeString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}}, TypeString())
}

func TestTypeInt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_INT}}, TypeInt())
}

func TestTypeFloat(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_FLOAT}}, TypeFloat())
}

func TestTypeDecimal(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_DECIMAL}}, TypeDecimal())
}

func TestTypeBool(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_BOOL}}, TypeBool())
}
