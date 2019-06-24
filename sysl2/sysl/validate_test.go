package main

import (
	"fmt"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestValidatorLogMsg(t *testing.T) {

	cases := map[string]struct {
		input []ValidationMsg
	}{
		"Error":     {input: []ValidationMsg{{Message: "Error msg", MsgType: ERROR}}},
		"Warn":      {input: []ValidationMsg{{Message: "Warning msg", MsgType: WARN}}},
		"Info":      {input: []ValidationMsg{{Message: "Info msg", MsgType: INFO}}},
		"Undefined": {input: []ValidationMsg{{Message: "Undefined msg", MsgType: UNDEF}}},
	}

	for name, test := range cases {
		input := test.input
		t.Run(name, func(t *testing.T) {
			logMsg(input...)
		})
	}
}

func TestValidatorViewOutput(t *testing.T) {
	cases := map[string]struct {
		input     *sysl.Type
		expected1 string
		expected2 bool
	}{
		"Primitive string": {
			input:     &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}},
			expected1: "STRING", expected2: false},
		"Primitive bool": {
			input:     &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_BOOL}},
			expected1: "BOOL", expected2: false},
		"Primitive int": {
			input:     &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_INT}},
			expected1: "INT", expected2: false},
		"Primitive float": {
			input:     &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_FLOAT}},
			expected1: "FLOAT", expected2: false},
		"Primitive decimal": {
			input:     &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_DECIMAL}},
			expected1: "DECIMAL", expected2: false},
		"Primitive no type": {
			input:     &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_EMPTY}},
			expected1: "EMPTY", expected2: false},
		"Sequence of primitives": {
			input: &sysl.Type{Type: &sysl.Type_Sequence{
				Sequence: &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_INT}}}},
			expected1: "INT", expected2: true},
		"Sequence of ref": {
			input: &sysl.Type{Type: &sysl.Type_Sequence{
				Sequence: &sysl.Type{Type: &sysl.Type_TypeRef{TypeRef: &sysl.ScopedRef{
					Ref: &sysl.Scope{Path: []string{"RefType"}}}}}}},
			expected1: "RefType", expected2: true},
		"Ref": {
			input: &sysl.Type{Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{Ref: &sysl.Scope{Path: []string{"RefType"}}}}},
			expected1: "RefType", expected2: false},
		"Unknown": {
			input:     &sysl.Type{Type: &sysl.Type_Map_{}},
			expected1: "Unknown", expected2: false},
	}

	for name, test := range cases {
		input := test.input
		expected1 := test.expected1
		expected2 := test.expected2
		t.Run(name, func(t *testing.T) {
			out, isColl := viewOutput(input)
			assert.Equal(t, expected1, out, "Unexpected result")
			assert.Equal(t, expected2, isColl, "Unexpected result")
		})
	}
}

func TestValidatorValidateEntryPoint(t *testing.T) {
	start := "EntryPoint"
	transform, _ := loadAndGetDefaultApp("tests", "transform1.sysl")

	var entryPointView *sysl.View
	var nonEntryPointView *sysl.View
	var invalidEntryPointView *sysl.View

	for _, tfm := range transform.GetApps() {
		entryPointView = tfm.Views[start]
		nonEntryPointView = tfm.Views["TfmDefaultEmpty"]
		invalidEntryPointView = tfm.Views["EntryPointInvalid"]
	}

	cases := map[string]struct {
		input    map[string]*sysl.View
		expected []ValidationMsg
	}{
		"Exists": {input: map[string]*sysl.View{start: entryPointView, "nonEntryPoint": nonEntryPointView},
			expected: nil},
		"Not exists": {input: map[string]*sysl.View{"nonEntryPoint": nonEntryPointView},
			expected: []ValidationMsg{
				{Message: fmt.Sprintf("Entry point view: '%s' is undefined", start), MsgType: ERROR}}},
		"Incorrect output": {
			input: map[string]*sysl.View{start: invalidEntryPointView, "nonEntryPoint": nonEntryPointView},
			expected: []ValidationMsg{
				{Message: fmt.Sprintf("Output type of entry point view: '%s' should be '%s'", start, start),
					MsgType: ERROR}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateEntryPoint(input, start)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateFileName(t *testing.T) {
	viewName := "filename"
	transform, _ := loadAndGetDefaultApp("tests", "transform1.sysl")

	var fileNameView *sysl.View
	var nonFileNameView *sysl.View
	var invalidFileNameView1 *sysl.View
	var invalidFileNameView2 *sysl.View
	var invalidFileNameView3 *sysl.View

	for _, tfm := range transform.GetApps() {
		fileNameView = tfm.Views[viewName]
		nonFileNameView = tfm.Views["TfmDefaultEmpty"]
		invalidFileNameView1 = tfm.Views["TfmFilenameInvalid1"]
		invalidFileNameView2 = tfm.Views["TfmFilenameInvalid2"]
		invalidFileNameView3 = tfm.Views["TfmFilenameInvalid3"]
	}

	cases := map[string]struct {
		input    map[string]*sysl.View
		expected []ValidationMsg
	}{
		"Exists": {input: map[string]*sysl.View{viewName: fileNameView, "nonEntryPoint": nonFileNameView},
			expected: []ValidationMsg{}},
		"Not exists": {input: map[string]*sysl.View{"tfmDefaultEmpty": nonFileNameView},
			expected: []ValidationMsg{{Message: "View 'filename' is undefined", MsgType: ERROR}}},
		"Incorrect output": {input: map[string]*sysl.View{viewName: invalidFileNameView1},
			expected: []ValidationMsg{{Message: "In view 'filename', output type should be 'string'", MsgType: ERROR}}},
		"Incorrect assignment": {input: map[string]*sysl.View{viewName: invalidFileNameView2},
			expected: []ValidationMsg{{Message: "In view 'filename' : missing type: 'filename'", MsgType: ERROR}}},
		"Excess assignment": {input: map[string]*sysl.View{viewName: invalidFileNameView3},
			expected: []ValidationMsg{
				{Message: fmt.Sprintf("In view 'filename' : Excess assignments: 'foo'"), MsgType: ERROR}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateFileName(input)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorHasSameType(t *testing.T) {
	type inputData struct {
		type1 *sysl.Type
		type2 *sysl.Type
	}
	type expectedData struct {
		isSame       bool
		expectedType string
	}
	cases := map[string]struct {
		input    inputData
		expected expectedData
	}{
		"Same primitive types": {
			input:    inputData{type1: stringType, type2: stringType},
			expected: expectedData{isSame: true, expectedType: "STRING"}},
		"Different primitive types1": {
			input:    inputData{type1: stringType, type2: intType},
			expected: expectedData{isSame: false, expectedType: "STRING"}},
		"Different primitive types2": {
			input:    inputData{type1: intType, type2: stringType},
			expected: expectedData{isSame: false, expectedType: "INT"}},
		"Same transform typerefs1": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"Statement"}},
					},
				},
			}, type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"Statement"}},
					},
				},
			}},
			expected: expectedData{isSame: true, expectedType: "Statement"}},
		"Different transform typerefs1-1": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"Statement"}},
					},
				},
			}, type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"StatementList"}},
					},
				},
			}},
			expected: expectedData{isSame: false, expectedType: "Statement"}},
		"Different transform typerefs1-2": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"StatementList"}},
					},
				},
			}, type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"Statement"}},
					},
				},
			}},
			expected: expectedData{isSame: false, expectedType: "StatementList"}},
		"Same transform typerefs2": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}, type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}},
			expected: expectedData{isSame: true, expectedType: "Statement"}},
		"Different transform typerefs2-1": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}, type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"StatementList"}}},
					},
				},
			}},
			expected: expectedData{isSame: false, expectedType: "Statement"}},
		"Different transform typerefs2-2": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"StatementList"}}},
					},
				},
			}, type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}},
			expected: expectedData{isSame: false, expectedType: "StatementList"}},
		"Different types1": {
			input:    inputData{type1: noType(), type2: stringType},
			expected: expectedData{isSame: false, expectedType: "Unresolved"}},
		"Different types2": {
			input:    inputData{type1: stringType, type2: noType()},
			expected: expectedData{isSame: false, expectedType: "STRING"}},
		"Different types3": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}, type2: stringType},
			expected: expectedData{isSame: false, expectedType: "Statement"}},
		"Different types4": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"StatementList"}},
					},
				},
			}, type2: stringType},
			expected: expectedData{isSame: false, expectedType: "StatementList"}},
		"Nil types": {
			input:    inputData{type1: nil, type2: nil},
			expected: expectedData{isSame: false, expectedType: "nil"}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			isSame, actual := hasSameType(input.type1, input.type2)
			assert.Equal(t, expected.expectedType, actual, "Unexpected result")
			assert.True(t, expected.isSame == isSame, "Unexpected result")
		})
	}
}

func TestValidatorResolveVariableType(t *testing.T) {
	type inputData struct {
		viewName string
		expr     *sysl.Expr
	}
	type expectedData struct {
		syslType       *sysl.Type
		validationMsgs []ValidationMsg
	}

	expressions := map[string]*sysl.Expr{}

	transform, _ := loadAndGetDefaultApp("tests", "transform1.sysl")

	for _, tfm := range transform.GetApps() {
		for _, stmt := range tfm.Views["varTypeResolve"].GetExpr().GetTransform().GetStmt() {
			expressions[stmt.GetAssign().GetName()] = stmt.GetAssign().GetExpr()
		}
	}

	cases := map[string]struct {
		input    inputData
		expected expectedData
	}{
		"String": {
			input:    inputData{viewName: "varTypeResolve", expr: expressions["stringType"]},
			expected: expectedData{syslType: stringType, validationMsgs: []ValidationMsg{}}},
		"Int": {
			input:    inputData{viewName: "varTypeResolve", expr: expressions["intType"]},
			expected: expectedData{syslType: intType, validationMsgs: []ValidationMsg{}}},
		"Bool": {
			input:    inputData{viewName: "varTypeResolve", expr: expressions["boolType"]},
			expected: expectedData{syslType: boolType, validationMsgs: []ValidationMsg{}}},
		"Decimal": {
			input:    inputData{viewName: "varTypeResolve", expr: expressions["decimalType"]},
			expected: expectedData{syslType: decimalType, validationMsgs: []ValidationMsg{}}},
		"Transform type primitive": {
			input:    inputData{viewName: "varTypeResolve", expr: expressions["transformTypePrimitive"]},
			expected: expectedData{syslType: stringType, validationMsgs: []ValidationMsg{}}},
		"Transform type ref": {
			input: inputData{viewName: "varTypeResolve", expr: expressions["transformTypeRef"]},
			expected: expectedData{syslType: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"Statement"}},
					},
				},
			}, validationMsgs: []ValidationMsg{}}},
		"Valid unary result": {
			input:    inputData{viewName: "varTypeResolve", expr: expressions["unaryResultValid"]},
			expected: expectedData{syslType: boolType, validationMsgs: []ValidationMsg{}}},
		"Invalid unary result": {
			input: inputData{viewName: "varTypeResolve", expr: expressions["unaryResultInvalid"]},
			expected: expectedData{
				syslType: boolType,
				validationMsgs: []ValidationMsg{{
					Message: "In view 'varTypeResolve', unary operator used with non boolean type: 'STRING'",
					MsgType: 100},
				}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			syslType, validationMsgs := resolveVariableType(input.expr, input.viewName)
			hasSameType, _ := hasSameType(expected.syslType, syslType)
			assert.True(t, hasSameType, "Unexpected result")
			assert.Equal(t, expected.validationMsgs, validationMsgs, "Unexpected result")
		})
	}
}

func TestValidatorValidateTypes(t *testing.T) {
	type inputData struct {
		expression *sysl.Expr
		viewName   string
		implViews  map[string]*sysl.View
		out        string
	}

	var tfmViews map[string]*sysl.View

	transform, _ := loadAndGetDefaultApp("tests", "transform1.sysl")
	grammar, _ := loadAndGetDefaultApp("tests", "go.gen.sysl")

	for _, tfm := range transform.GetApps() {
		tfmViews = tfm.Views
	}

	cases := map[string]struct {
		input    inputData
		expected []ValidationMsg
	}{
		"Equal": {
			input: inputData{
				viewName:   "TfmValid",
				expression: tfmViews["TfmValid"].GetExpr(),
				implViews:  tfmViews,
				out:        "MethodDecl"},
			expected: []ValidationMsg{}},
		"Not Equal": {
			input: inputData{
				viewName:   "TfmInvalid",
				expression: tfmViews["TfmInvalid"].GetExpr(),
				implViews:  tfmViews,
				out:        "MethodDecl"},
			expected: []ValidationMsg{
				{Message: "In view 'TfmInvalid', type 'FunctionName' is missing", MsgType: ERROR}}},
		"Absent optional": {
			input: inputData{
				viewName:   "TfmNoOptional",
				expression: tfmViews["TfmNoOptional"].GetExpr(),
				implViews:  tfmViews,
				out:        "MethodDecl"},
			expected: []ValidationMsg{}},
		"Excess attributes without optionals": {
			input: inputData{
				viewName:   "TfmExcessAttrs1",
				expression: tfmViews["TfmExcessAttrs1"].GetExpr(),
				implViews:  tfmViews,
				out:        "MethodDecl"},
			expected: []ValidationMsg{
				{Message: "In view 'TfmExcessAttrs1', excess attribute is defined: 'ExcessAttr1'", MsgType: ERROR}}},
		"Excess attributes with optionals": {
			input: inputData{
				viewName:   "TfmExcessAttrs2",
				expression: tfmViews["TfmExcessAttrs2"].GetExpr(),
				implViews:  tfmViews,
				out:        "MethodDecl"},
			expected: []ValidationMsg{
				{Message: "In view 'TfmExcessAttrs2', excess attribute is defined: 'ExcessAttr1'", MsgType: ERROR}}},
		"Valid choice": {
			input: inputData{
				viewName:   "ValidChoice",
				expression: tfmViews["ValidChoice"].GetExpr(),
				implViews:  tfmViews,
				out:        "Statement"},
			expected: []ValidationMsg{}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateTypes(grammar, input.expression, input.viewName, input.implViews, input.out)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateInnerTypes(t *testing.T) {
	type inputData struct {
		expression *sysl.Expr
		viewName   string
		implViews  map[string]*sysl.View
		out        string
	}

	transform, _ := loadAndGetDefaultApp("tests", "transform1.sysl")
	grammar, _ := loadAndGetDefaultApp("tests", "go.gen.sysl")

	var tfmViews map[string]*sysl.View
	for _, tfm := range transform.GetApps() {
		tfmViews = tfm.Views
	}

	cases := map[string]struct {
		input    inputData
		expected []ValidationMsg
	}{
		"Valid inner type": {
			input: inputData{
				expression: tfmViews["ValidInnerAttrs"].GetExpr(),
				viewName:   "ValidInnerAttrs",
				implViews:  tfmViews,
				out:        "goFile",
			},
			expected: []ValidationMsg{}},
		"Invalid inner type": {
			input: inputData{
				expression: tfmViews["InvalidInnerAttrs"].GetExpr(),
				viewName:   "InvalidInnerAttrs",
				implViews:  tfmViews,
				out:        "goFile",
			},
			expected: []ValidationMsg{
				{Message: "In view 'InvalidInnerAttrs', type 'PackageName' is missing", MsgType: ERROR},
				{Message: "In view 'InvalidInnerAttrs', excess attribute is defined: 'Foo'", MsgType: ERROR}}},
	}
	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateTypes(grammar, input.expression, input.viewName, input.implViews, input.out)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateChoiceTypes(t *testing.T) {
	type inputData struct {
		expression *sysl.Expr
		viewName   string
		implViews  map[string]*sysl.View
		out        string
	}

	transform, _ := loadAndGetDefaultApp("tests", "transform1.sysl")
	grammar, _ := loadAndGetDefaultApp("tests", "go.gen.sysl")

	var tfmViews map[string]*sysl.View
	for _, tfm := range transform.GetApps() {
		tfmViews = tfm.Views
	}

	cases := map[string]struct {
		input    inputData
		expected []ValidationMsg
	}{
		"Valid choice": {
			input: inputData{
				expression: tfmViews["ValidChoice"].GetExpr(),
				viewName:   "ValidChoice",
				implViews:  tfmViews,
				out:        "Statement"},
			expected: []ValidationMsg{}},
		"Invalid choice": {
			input: inputData{
				expression: tfmViews["InvalidChoice"].GetExpr(),
				viewName:   "InvalidChoice",
				implViews:  tfmViews,
				out:        "Statement"},
			expected: []ValidationMsg{
				{Message: "In view 'InvalidChoice', invalid choice has been made as : 'Foo'", MsgType: ERROR},
				{Message: "In view 'InvalidChoice', excess attribute is defined: 'Foo'", MsgType: ERROR}}},
		"Valid choice combination": {
			input: inputData{
				expression: tfmViews["ValidChoiceCombination"].GetExpr(),
				viewName:   "ValidChoiceCombination",
				implViews:  tfmViews,
				out:        "MethodSpec"},
			expected: []ValidationMsg{}},
		"Valid choice non-combination": {
			input: inputData{
				expression: tfmViews["ValidChoiceNonCombination"].GetExpr(),
				viewName:   "ValidChoiceNonCombination",
				implViews:  tfmViews,
				out:        "MethodSpec"},
			expected: []ValidationMsg{}},
		"Invalid choice combination excess": {
			input: inputData{
				expression: tfmViews["InvalidChoiceCombinationExcess"].GetExpr(),
				viewName:   "InvalidChoiceCombinationExcess",
				implViews:  tfmViews,
				out:        "MethodSpec"},
			expected: []ValidationMsg{{
				Message: "In view 'InvalidChoiceCombinationExcess', excess attribute is defined: 'Foo'",
				MsgType: ERROR}}},
		"Invalid choice combination missing": {
			input: inputData{
				expression: tfmViews["InvalidChoiceCombiMissing"].GetExpr(),
				viewName:   "InvalidChoiceCombiMissing",
				implViews:  tfmViews,
				out:        "MethodSpec"},
			expected: []ValidationMsg{
				{Message: "In view 'InvalidChoiceCombiMissing', type 'Signature' is missing", MsgType: ERROR},
				{Message: "In view 'InvalidChoiceCombiMissing', excess attribute is defined: 'Foo'", MsgType: ERROR}}},
		"Invalid choice non-combination missing": {
			input: inputData{
				expression: tfmViews["InvalidChoiceNonCombination"].GetExpr(),
				viewName:   "InvalidChoiceNonCombination",
				implViews:  tfmViews,
				out:        "MethodSpec"},
			expected: []ValidationMsg{
				{
					Message: "In view 'InvalidChoiceNonCombination', invalid choice has been made as : 'Interface'",
					MsgType: ERROR},
				{
					Message: "In view 'InvalidChoiceNonCombination', excess attribute is defined: 'Interface'",
					MsgType: 100}}},
	}
	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateTypes(grammar, input.expression, input.viewName, input.implViews, input.out)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateTransform(t *testing.T) {
	actual := validateTransform("tests", "transform2.sysl", "tests/go.gen.g", "goFile")
	assert.Equal(t, []ValidationMsg{}, actual, "Unexpected result")
}
