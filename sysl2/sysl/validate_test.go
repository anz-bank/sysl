package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTypeRef(t *testing.T) {
	transform, _ := loadAndGetDefaultApp("tests", "validation_transform.sysl")

	typeRefs := []string{"list", "string", "BinExpr", "Conditional", "string", "View", "Transform", "variable"}

	for _, tfm := range transform.GetApps() {
		tfmView := tfm.Views["TfmTypeRefTest"]
		for i, stmt := range tfmView.GetExpr().GetTransform().GetStmt() {
			if stmt.GetAssign() != nil {
				assert.Equal(t, typeRefs[i], typeRef(stmt.GetAssign().GetExpr()))
			}
		}
	}
}

func TestReadView(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected map[string]content
	}{
		"Basic": {input: "TfmReadViewBasicTest",
			expected: map[string]content{"TransformVar": {Name: "TransformVar", TypeRef: "string"}}},
		"With children": {input: "TfmReadViewChildrenTest",
			expected: map[string]content{"TransformVar": {Name: "TransformVar", TypeRef: "Transform",
				Children: map[string]content{"Foo": {Name: "Foo", TypeRef: "string"}}}}},
		"Empty": {input: "TfmReadViewEmptyTest", expected: map[string]content{}},
	}

	transform, _ := loadAndGetDefaultApp("tests", "validation_transform.sysl")

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := map[string]content{}
			for _, tfm := range transform.GetApps() {
				tfmView := tfm.Views[input]

				for _, stmt := range tfmView.GetExpr().GetTransform().GetStmt() {
					if stmt.GetAssign() != nil {
						actual[stmt.GetAssign().GetName()] = readView(stmt)
					}
				}
			}
			assert.Equalf(t, expected, actual, "Unexpected result", input)
		})
	}
}

func TestOutputType(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected string
	}{
		"Primitive": {input: "TfmOutputTypePrimitiveTest",
			expected: "STRING"},
		"Primitive Collection": {input: "TfmOutputTypePrimitiveCollectionTest",
			expected: "sequence of STRING"},
		"Transform": {input: "TfmOutputTypeTransformTest",
			expected: "Transform"},
		"Transform Collection": {input: "TfmOutputTypeTfmCollectionTest",
			expected: "sequence of Transform"},
	}

	transform, _ := loadAndGetDefaultApp("tests", "validation_transform.sysl")

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := ""
			for _, tfm := range transform.GetApps() {
				tfmView := tfm.Views[input]
				actual = outputType(tfmView.GetRetType())
			}
			assert.Equalf(t, expected, actual, "Unexpected result", input)
		})
	}
}

func TestValidateEntryPoint(t *testing.T) {
	start := "EntryPoint"

	cases := map[string]struct {
		input    Views
		expected ValidationMsg
	}{
		"Exists": {input: map[string]View{start: {Name: start, OutputType: start}},
			expected: ValidationMsg{message: "", msgType: UNDEF}},
		"Not exists": {input: map[string]View{"view": {Name: "view"}},
			expected: ValidationMsg{message: "[Validator]: Entry point view: " + start + " is undefined", msgType: ERROR}},
		"Incorrect output": {input: map[string]View{start: {Name: start, OutputType: "string"}},
			expected: ValidationMsg{message: "[Validator]: Output type of entry point view: " + start + " should be " + start, msgType: ERROR}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateEntryPoint(input, start)
			assert.Equalf(t, expected, actual, "Unexpected result", input)
		})
	}
}

func TestValidateFileName(t *testing.T) {
	viewName := "filename"

	cases := map[string]struct {
		input    Views
		expected ValidationMsg
	}{
		"Valid": {input: map[string]View{viewName: {Name: viewName, OutputType: "STRING", Body: map[string]content{viewName: {Name: viewName}}}},
			expected: ValidationMsg{message: "", msgType: UNDEF}},
		"Invalid": {input: map[string]View{viewName: {Name: viewName, OutputType: "STRING", Body: map[string]content{"Bar": {Name: "Bar"}}}},
			expected: ValidationMsg{message: "[Validator]: In view filename : Missing type: filename", msgType: ERROR}},
		"Not exists": {input: map[string]View{},
			expected: ValidationMsg{message: "[Validator]: view: filename is undefined", msgType: ERROR}},
		"Incorrect output": {input: map[string]View{viewName: {Name: viewName, OutputType: "Foo", Body: map[string]content{viewName: {Name: viewName}}}},
			expected: ValidationMsg{message: "[Validator]: Output type of view: filename should be string", msgType: ERROR}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			actual := validateFileName(input)
			assert.Equalf(t, expected, actual, "Unexpected result", input)
		})
	}
}

func TestValidateRetType(t *testing.T) {
	cases := map[string]struct {
		input    string
		isSeq    bool
		expected []ValidationMsg
	}{
		"Missing type": {input: "TfmRetTypeMissingTest", isSeq: false,
			expected: []ValidationMsg{{message: "[Validator]: In view TfmRetTypeMissingTest : Missing type: ParameterDecl", msgType: ERROR}}},
		"Valid type": {input: "TfmRetTypeValidTest", isSeq: false},
		"Sequence invalid": {input: "TfmRetTypeSeqInvalidTest", isSeq: true,
			expected: []ValidationMsg{{message: "[Validator]: In view TfmRetTypeSeqInvalidTest : Expects sequence type as return value", msgType: ERROR}}},
		"Sequence valid": {input: "TfmRetTypeSeqValidTest", isSeq: true},
		"Multiple validation failures": {input: "TfmRetTypeMultiInvalidTest", isSeq: true,
			expected: []ValidationMsg{
				{message: "[Validator]: In view TfmRetTypeMultiInvalidTest : Missing type: identifier", msgType: ERROR},
				{message: "[Validator]: In view TfmRetTypeMultiInvalidTest : Expects sequence type as return value", msgType: ERROR}}},
	}

	transform, _ := loadAndGetDefaultApp("tests", "validation_transform.sysl")
	grm, _ := loadAndGetDefaultApp("tests", "go.gen.sysl")

	for name, test := range cases {
		input := test.input
		isSeq := test.isSeq
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			var view View
			for _, tfm := range transform.GetApps() {
				tfmView := tfm.Views[input]
				view = View{
					Name:       input,
					Body:       map[string]content{},
					OutputType: outputType(tfmView.GetRetType()),
					RawData:    tfmView,
				}

				for _, stmt := range tfmView.GetExpr().GetTransform().GetStmt() {
					if stmt.GetAssign() != nil {
						view.Body[stmt.GetAssign().GetName()] = readView(stmt)
					}
				}
			}

			for _, grammarSysl := range grm.GetApps() {
				outputType := view.OutputType
				typeRef := strings.Replace(outputType, "sequence of", "", -1)
				typeRef = strings.TrimSpace(typeRef)

				if grammarType, exists := grammarSysl.Types[typeRef]; exists {
					assert.Equalf(t, expected, validateRetType(view, grammarType, isSeq), "Unexpected result", input)
				}
			}

		})
	}
}

func TestValidate(t *testing.T) {
	transform, _ := loadAndGetDefaultApp("tests", "valid_transform.sysl")
	grm, _ := loadAndGetDefaultApp("tests", "go.gen.sysl")

	validate("goFile", transform, grm)
}

func TestLogMsg(t *testing.T) {
	logMsg(ValidationMsg{message: "Info msg", msgType: INFO})
	logMsg(ValidationMsg{message: "Info msg", msgType: WARN})
	logMsg(ValidationMsg{message: "Info msg", msgType: ERROR})
}
