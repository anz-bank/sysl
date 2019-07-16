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

func TestValidatorGetTypeName(t *testing.T) {
	cases := map[string]struct {
		input    *sysl.Type
		expected string
	}{
		"Primitive string": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}},
			expected: "STRING"},
		"Primitive bool": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_BOOL}},
			expected: "BOOL"},
		"Primitive int": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_INT}},
			expected: "INT"},
		"Primitive float": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_FLOAT}},
			expected: "FLOAT"},
		"Primitive decimal": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_DECIMAL}},
			expected: "DECIMAL"},
		"Primitive no type": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_EMPTY}},
			expected: "EMPTY"},
		"Sequence of primitives": {
			input: &sysl.Type{Type: &sysl.Type_Sequence{
				Sequence: &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_INT}}}},
			expected: "INT"},
		"Sequence of ref": {
			input: &sysl.Type{Type: &sysl.Type_Sequence{
				Sequence: &sysl.Type{Type: &sysl.Type_TypeRef{TypeRef: &sysl.ScopedRef{
					Ref: &sysl.Scope{Path: []string{"RefType"}}}}}}},
			expected: "RefType"},
		"Ref": {
			input: &sysl.Type{Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{Ref: &sysl.Scope{Path: []string{"RefType"}}}}},
			expected: "RefType"},
		"Unknown": {
			input:    &sysl.Type{Type: &sysl.Type_Map_{}},
			expected: "Unknown"},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			typeName := getTypeName(input)
			assert.Equal(t, expected, typeName, "Unexpected result")
		})
	}
}

func TestValidatorIsCollectionType(t *testing.T) {
	cases := map[string]struct {
		input    *sysl.Type
		expected bool
	}{
		"Sequence": {
			input:    &sysl.Type{Type: &sysl.Type_Sequence{}},
			expected: true},
		"Map": {
			input:    &sysl.Type{Type: &sysl.Type_Map_{}},
			expected: true},
		"Set": {
			input:    &sysl.Type{Type: &sysl.Type_Set{}},
			expected: true},
		"List": {
			input:    &sysl.Type{Type: &sysl.Type_List_{}},
			expected: true},
		"Primitive string": {
			input:    &sysl.Type{Type: &sysl.Type_Primitive_{}},
			expected: false},
		"Ref": {
			input: &sysl.Type{Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{Ref: &sysl.Scope{Path: []string{"RefType"}}}}},
			expected: false},
		"Unknown": {
			input:    &sysl.Type{Type: &sysl.Type_NoType_{}},
			expected: false},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			typeName := isCollectionType(input)
			assert.Equal(t, expected, typeName, "Unexpected result")
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
	cases := map[string]struct {
		input    inputData
		expected bool
	}{
		"Same primitive types": {
			input:    inputData{type1: stringType, type2: stringType},
			expected: true},
		"Different primitive types1": {
			input:    inputData{type1: stringType, type2: intType},
			expected: false},
		"Different primitive types2": {
			input:    inputData{type1: intType, type2: stringType},
			expected: false},
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
			expected: true},
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
			expected: false},
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
			expected: false},
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
			expected: true},
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
			expected: false},
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
			expected: false},
		"Different types1": {
			input:    inputData{type1: noType(), type2: stringType},
			expected: false},
		"Different types2": {
			input:    inputData{type1: stringType, type2: noType()},
			expected: false},
		"Different types3": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}, type2: stringType},
			expected: false},
		"Different types3.5": {
			input: inputData{type2: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}},
					},
				},
			}, type1: stringType},
			expected: false},
		"Different types4": {
			input: inputData{type1: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: []string{"StatementList"}},
					},
				},
			}, type2: stringType},
			expected: false},
		"Nil types": {
			input:    inputData{type1: nil, type2: nil},
			expected: false},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			isSame := hasSameType(input.type1, input.type2)
			assert.True(t, expected == isSame, "Unexpected result")
		})
	}
}

func TestValidatorResolveExprType(t *testing.T) {
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
			syslType, validationMsgs := resolveExprType(input.expr, input.viewName)
			assert.True(t, hasSameType(expected.syslType, syslType), "Unexpected result")
			assert.Equal(t, expected.validationMsgs, validationMsgs, "Unexpected result")
		})
	}
}

func TestValidatorValidateTransform(t *testing.T) {
	type inputData struct {
		transform *sysl.Expr_Transform
		viewName  string
		implViews map[string]*sysl.View
		out       string
	}

	transformModule, tfmAppName := loadAndGetDefaultApp("tests", "transform1.sysl")
	grammarModule, grammarAppName := loadAndGetDefaultApp("tests", "go.gen.sysl")

	grammar := grammarModule.GetApps()[grammarAppName]
	tfmViews := transformModule.GetApps()[tfmAppName].GetViews()

	cases := map[string]struct {
		input    inputData
		expected []ValidationMsg
	}{
		"Equal": {
			input: inputData{
				viewName:  "TfmValid",
				transform: tfmViews["TfmValid"].GetExpr().GetTransform(),
				implViews: tfmViews,
				out:       "MethodDecl"},
			expected: []ValidationMsg{}},
		"Not Equal": {
			input: inputData{
				viewName:  "TfmInvalid",
				transform: tfmViews["TfmInvalid"].GetExpr().GetTransform(),
				implViews: tfmViews,
				out:       "MethodDecl"},
			expected: []ValidationMsg{
				{Message: "In view 'TfmInvalid', type 'FunctionName' is missing", MsgType: ERROR}}},
		"Absent optional": {
			input: inputData{
				viewName:  "TfmNoOptional",
				transform: tfmViews["TfmNoOptional"].GetExpr().GetTransform(),
				implViews: tfmViews,
				out:       "MethodDecl"},
			expected: []ValidationMsg{}},
		"Excess attributes without optionals": {
			input: inputData{
				viewName:  "TfmExcessAttrs1",
				transform: tfmViews["TfmExcessAttrs1"].GetExpr().GetTransform(),
				implViews: tfmViews,
				out:       "MethodDecl"},
			expected: []ValidationMsg{
				{Message: "In view 'TfmExcessAttrs1', excess attribute is defined: 'ExcessAttr1'", MsgType: ERROR}}},
		"Excess attributes with optionals": {
			input: inputData{
				viewName:  "TfmExcessAttrs2",
				transform: tfmViews["TfmExcessAttrs2"].GetExpr().GetTransform(),
				implViews: tfmViews,
				out:       "MethodDecl"},
			expected: []ValidationMsg{
				{Message: "In view 'TfmExcessAttrs2', excess attribute is defined: 'ExcessAttr1'", MsgType: ERROR}}},
		"Valid choice": {
			input: inputData{
				viewName:  "ValidChoice",
				transform: tfmViews["ValidChoice"].GetExpr().GetTransform(),
				implViews: tfmViews,
				out:       "Statement"},
			expected: []ValidationMsg{}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateTransform(grammar, input.transform, input.viewName, input.implViews, input.out)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateTransformInnerTypes(t *testing.T) {
	type inputData struct {
		transform *sysl.Expr_Transform
		viewName  string
		implViews map[string]*sysl.View
		out       string
	}

	transformModule, tfmAppName := loadAndGetDefaultApp("tests", "transform1.sysl")
	grammarModule, grammarAppName := loadAndGetDefaultApp("tests", "go.gen.sysl")

	grammar := grammarModule.GetApps()[grammarAppName]
	tfmViews := transformModule.GetApps()[tfmAppName].GetViews()

	cases := map[string]struct {
		input    inputData
		expected []ValidationMsg
	}{
		"Valid inner type": {
			input: inputData{
				transform: tfmViews["ValidInnerAttrs"].GetExpr().GetTransform(),
				viewName:  "ValidInnerAttrs",
				implViews: tfmViews,
				out:       "goFile",
			},
			expected: []ValidationMsg{}},
		"Invalid inner type": {
			input: inputData{
				transform: tfmViews["InvalidInnerAttrs"].GetExpr().GetTransform(),
				viewName:  "InvalidInnerAttrs",
				implViews: tfmViews,
				out:       "goFile",
			},
			expected: []ValidationMsg{
				{Message: "In view 'InvalidInnerAttrs', type 'PackageName' is missing", MsgType: ERROR},
				{Message: "In view 'InvalidInnerAttrs', excess attribute is defined: 'Foo'", MsgType: ERROR}}},
	}
	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateTransform(grammar, input.transform, input.viewName, input.implViews, input.out)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateTransformChoiceTypes(t *testing.T) {
	type inputData struct {
		transform *sysl.Expr_Transform
		viewName  string
		implViews map[string]*sysl.View
		out       string
	}

	transformModule, tfmAppName := loadAndGetDefaultApp("tests", "transform1.sysl")
	grammarModule, grammarAppName := loadAndGetDefaultApp("tests", "go.gen.sysl")

	grammar := grammarModule.GetApps()[grammarAppName]
	tfmViews := transformModule.GetApps()[tfmAppName].GetViews()

	cases := map[string]struct {
		input    inputData
		expected []ValidationMsg
	}{
		"Valid choice": {
			input: inputData{
				transform: tfmViews["ValidChoice"].GetExpr().GetTransform(),
				viewName:  "ValidChoice",
				implViews: tfmViews,
				out:       "Statement"},
			expected: []ValidationMsg{}},
		"Invalid choice": {
			input: inputData{
				transform: tfmViews["InvalidChoice"].GetExpr().GetTransform(),
				viewName:  "InvalidChoice",
				implViews: tfmViews,
				out:       "Statement"},
			expected: []ValidationMsg{
				{Message: "In view 'InvalidChoice', invalid choice has been made as : 'Foo'", MsgType: ERROR},
				{Message: "In view 'InvalidChoice', excess attribute is defined: 'Foo'", MsgType: ERROR}}},
		"Valid choice combination": {
			input: inputData{
				transform: tfmViews["ValidChoiceCombination"].GetExpr().GetTransform(),
				viewName:  "ValidChoiceCombination",
				implViews: tfmViews,
				out:       "MethodSpec"},
			expected: []ValidationMsg{}},
		"Valid choice non-combination": {
			input: inputData{
				transform: tfmViews["ValidChoiceNonCombination"].GetExpr().GetTransform(),
				viewName:  "ValidChoiceNonCombination",
				implViews: tfmViews,
				out:       "MethodSpec"},
			expected: []ValidationMsg{}},
		"Invalid choice combination excess": {
			input: inputData{
				transform: tfmViews["InvalidChoiceCombinationExcess"].GetExpr().GetTransform(),
				viewName:  "InvalidChoiceCombinationExcess",
				implViews: tfmViews,
				out:       "MethodSpec"},
			expected: []ValidationMsg{{
				Message: "In view 'InvalidChoiceCombinationExcess', excess attribute is defined: 'Foo'",
				MsgType: ERROR}}},
		"Invalid choice combination missing": {
			input: inputData{
				transform: tfmViews["InvalidChoiceCombiMissing"].GetExpr().GetTransform(),
				viewName:  "InvalidChoiceCombiMissing",
				implViews: tfmViews,
				out:       "MethodSpec"},
			expected: []ValidationMsg{
				{Message: "In view 'InvalidChoiceCombiMissing', type 'Signature' is missing", MsgType: ERROR},
				{Message: "In view 'InvalidChoiceCombiMissing', excess attribute is defined: 'Foo'", MsgType: ERROR}}},
		"Invalid choice non-combination missing": {
			input: inputData{
				transform: tfmViews["InvalidChoiceNonCombination"].GetExpr().GetTransform(),
				viewName:  "InvalidChoiceNonCombination",
				implViews: tfmViews,
				out:       "MethodSpec"},
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
			actual := validateTransform(grammar, input.transform, input.viewName, input.implViews, input.out)
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidate(t *testing.T) {
	transformModule, tfmAppName := loadAndGetDefaultApp("tests", "transform2.sysl")
	grammarModule, grammarAppName := loadAndGetDefaultApp("tests", "go.gen.sysl")

	grammar := grammarModule.GetApps()[grammarAppName]
	transform := transformModule.GetApps()[tfmAppName]

	actual := validate(grammar, transform, "goFile")
	assert.Equal(t, []ValidationMsg{}, actual, "Unexpected result")
}

func TestValidatorLoadTransformSuccess(t *testing.T) {
	tfm, err := loadTransform("tests", "transform2.sysl")
	assert.NotNil(t, tfm, "Unexpected result")
	assert.Nil(t, err, "Unexpected result")
}

func TestValidatorLoadTransformError(t *testing.T) {
	tfm, err := loadTransform("foo", "bar.sysl")
	assert.Nil(t, tfm, "Unexpected result")
	assert.NotNil(t, err, "Unexpected result")
}

func TestValidatorLoadGrammarSuccess(t *testing.T) {
	grammar, err := loadGrammar("tests/go.gen.g")
	assert.NotNil(t, grammar, "Unexpected result")
	assert.Nil(t, err, "Unexpected result")
}

func TestValidatorLoadGrammarError(t *testing.T) {
	grammar, err := loadGrammar("foo/bar.g")
	assert.Nil(t, grammar, "Unexpected result")
	assert.NotNil(t, err, "Unexpected result")
}
