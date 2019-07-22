package main

import (
	"flag"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

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
		"Nil": {
			input:    nil,
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
	p := NewParser()
	transform, _ := loadTransform("tests", "transform1.sysl", p)

	var entryPointView *sysl.View
	var nonEntryPointView *sysl.View
	var invalidEntryPointView *sysl.View

	entryPointView = transform.Views[start]
	nonEntryPointView = transform.Views["TfmDefaultEmpty"]
	invalidEntryPointView = transform.Views["EntryPointInvalid"]

	cases := map[string]struct {
		input    *sysl.Application
		expected map[string][]Msg
	}{
		"Exists": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{start: entryPointView, "nonEntryPoint": nonEntryPointView}},
			expected: map[string][]Msg{}},
		"Not exists": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{"nonEntryPoint": nonEntryPointView}},
			expected: map[string][]Msg{
				"EntryPoint": {{MessageID: ErrEntryPointUndefined, MessageData: []string{start}}}}},
		"Incorrect output": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{start: invalidEntryPointView, "nonEntryPoint": nonEntryPointView}},
			expected: map[string][]Msg{
				"EntryPoint": {{MessageID: ErrInvalidEntryPointReturn, MessageData: []string{start, start}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(nil, input, &Parser{messages: map[string][]Msg{}})
			validator.validateEntryPoint(start)
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateFileName(t *testing.T) {
	viewName := "filename"
	p := NewParser()
	transform, _ := loadTransform("tests", "transform1.sysl", p)

	var fileNameView *sysl.View
	var nonFileNameView *sysl.View
	var invalidFileNameView1 *sysl.View
	var invalidFileNameView2 *sysl.View
	var invalidFileNameView3 *sysl.View

	fileNameView = transform.Views[viewName]
	nonFileNameView = transform.Views["TfmDefaultEmpty"]
	invalidFileNameView1 = transform.Views["TfmFilenameInvalid1"]
	invalidFileNameView2 = transform.Views["TfmFilenameInvalid2"]
	invalidFileNameView3 = transform.Views["TfmFilenameInvalid3"]

	cases := map[string]struct {
		input    *sysl.Application
		expected map[string][]Msg
	}{
		"Exists": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{viewName: fileNameView, "nonEntryPoint": nonFileNameView}},
			expected: map[string][]Msg{}},
		"Not exists": {
			input: &sysl.Application{Views: map[string]*sysl.View{"tfmDefaultEmpty": nonFileNameView}},
			expected: map[string][]Msg{"filename": {
				{MessageID: ErrUndefinedView, MessageData: []string{viewName}}}}},
		"Incorrect output": {
			input: &sysl.Application{Views: map[string]*sysl.View{viewName: invalidFileNameView1}},
			expected: map[string][]Msg{"filename": {
				{MessageID: ErrInvalidReturn, MessageData: []string{viewName, "string"}}}}},
		"Incorrect assignment": {
			input: &sysl.Application{Views: map[string]*sysl.View{viewName: invalidFileNameView2}},
			expected: map[string][]Msg{"filename": {
				{MessageID: ErrMissingReqField, MessageData: []string{viewName, viewName, "string"}}}}},
		"Excess assignment": {
			input: &sysl.Application{Views: map[string]*sysl.View{viewName: invalidFileNameView3}},
			expected: map[string][]Msg{"filename": {
				{MessageID: ErrExcessAttr, MessageData: []string{"foo", viewName, "string"}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(nil, input, &Parser{messages: map[string][]Msg{}})
			validator.validateFileName()
			actual := validator.GetMessages()
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
			input:    inputData{type1: typeString(), type2: typeString()},
			expected: true},
		"Different primitive types1": {
			input:    inputData{type1: typeString(), type2: typeInt()},
			expected: false},
		"Different primitive types2": {
			input:    inputData{type1: typeInt(), type2: typeString()},
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
			input:    inputData{type1: typeNone(), type2: typeString()},
			expected: false},
		"Different types2": {
			input:    inputData{type1: typeString(), type2: typeNone()},
			expected: false},
		"Different types3": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}},
				type2: typeString()},
			expected: false},
		"Different types3.5": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{Part: []string{"Statement"}}}}}},
				type2: typeString()},
			expected: false},
		"Different types4": {
			input: inputData{
				type1: &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Path: []string{"StatementList"}}}}},
				type2: typeString()},
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
			isSame := hasSameType(input.type1, input.type2)
			assert.True(t, expected == isSame, "Unexpected result")
		})
	}
}

func TestValidatorValidateViews(t *testing.T) {
	p := NewParser()
	transform, _ := loadTransform("tests", "transform1.sysl", p)
	grammar, _ := loadGrammar("tests/grammar.sysl")

	cases := map[string]struct {
		input    string
		expected map[string][]Msg
	}{
		"Equal": {
			input:    "TfmValid",
			expected: map[string][]Msg{"TfmValid": nil}},
		"Not Equal": {
			input: "TfmInvalid",
			expected: map[string][]Msg{"TfmInvalid": {
				{MessageID: ErrMissingReqField, MessageData: []string{"FunctionName", "TfmInvalid", "MethodDecl"}}}}},
		"Absent optional": {
			input:    "TfmNoOptional",
			expected: map[string][]Msg{"TfmNoOptional": nil}},
		"Excess attributes without optionals": {
			input: "TfmExcessAttrs1",
			expected: map[string][]Msg{"TfmExcessAttrs1": {
				{MessageID: ErrExcessAttr, MessageData: []string{"ExcessAttr1", "TfmExcessAttrs1", "MethodDecl"}}}}},
		"Excess attributes with optionals": {
			input: "TfmExcessAttrs2",
			expected: map[string][]Msg{"TfmExcessAttrs2": {
				{MessageID: ErrExcessAttr, MessageData: []string{"ExcessAttr1", "TfmExcessAttrs2", "MethodDecl"}}}}},
		"Valid choice": {
			input:    "ValidChoice",
			expected: map[string][]Msg{"ValidChoice": nil}},
		"Relational Type": {
			input:    "Relational",
			expected: map[string][]Msg{"Relational": nil}},
		"Inner relational Type": {
			input:    "InnerRelational",
			expected: map[string][]Msg{"InnerRelational": nil}},
		"Transform variable valid": {
			input:    "TransformVarValid",
			expected: map[string][]Msg{"TransformVarValid": nil}},
		"Transform variable redefined": {
			input: "TransformVarRedefined",
			expected: map[string][]Msg{"TransformVarRedefined": {
				{MessageID: 409, MessageData: []string{"TransformVarRedefined", "varDeclaration"}}}}},
		"Transform inner-variable redefined": {
			input: "TransformInnerVarRedefined",
			expected: map[string][]Msg{"TransformInnerVarRedefined": {
				{MessageID: 409, MessageData: []string{"TransformInnerVarRedefined:varDeclaration", "foo"}}}}},
		"Transform assign redefined": {
			input: "TransformAssignRedefined",
			expected: map[string][]Msg{"TransformAssignRedefined": {
				{MessageID: 409, MessageData: []string{"TransformAssignRedefined", "VarDecl"}}}}},
		"Transform inner-assign redefined": {
			input: "TransformInnerAssignRedefined",
			expected: map[string][]Msg{"TransformInnerAssignRedefined": {
				{MessageID: 409, MessageData: []string{"TransformInnerAssignRedefined:VarDecl", "TypeName"}}}}},
		"Transform variable invalid": {
			input: "TransformVarInvalid:varDeclaration",
			expected: map[string][]Msg{"TransformVarInvalid:varDeclaration": {
				{MessageID: 405, MessageData: []string{"identifier", "TransformVarInvalid:varDeclaration", "VarDecl"}},
				{MessageID: 406, MessageData: []string{"foo", "TransformVarInvalid:varDeclaration", "VarDecl"}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(
				grammar, transform,
				&Parser{
					assignTypes: map[string]TypeData{input: p.assignTypes[input]},
					letTypes:    map[string]TypeData{input: p.letTypes[input]},
					messages:    map[string][]Msg{input: p.messages[input]}})
			validator.validateViews()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateViewsInnerTypes(t *testing.T) {
	p := NewParser()
	transform, _ := loadTransform("tests", "transform1.sysl", p)
	grammar, _ := loadGrammar("tests/grammar.sysl")

	cases := map[string]struct {
		inputAssign map[string]TypeData
		inputLet    map[string]TypeData
		expected    map[string][]Msg
	}{
		"Valid inner type": {
			inputAssign: map[string]TypeData{"ValidInnerAttrs": p.GetAssigns()["ValidInnerAttrs"]},
			inputLet:    map[string]TypeData{"ValidInnerAttrs": p.GetLets()["ValidInnerAttrs"]},
			expected:    map[string][]Msg{}},
		"Invalid inner type": {
			inputAssign: map[string]TypeData{"InvalidInnerAttrs": p.GetAssigns()["InvalidInnerAttrs"]},
			inputLet:    map[string]TypeData{"InvalidInnerAttrs": p.GetLets()["InvalidInnerAttrs"]},
			expected: map[string][]Msg{"InvalidInnerAttrs": {
				{MessageID: ErrMissingReqField, MessageData: []string{"PackageName", "InvalidInnerAttrs", "PackageClause"}},
				{MessageID: ErrExcessAttr, MessageData: []string{"Foo", "InvalidInnerAttrs", "PackageClause"}}}}},
	}
	for name, test := range cases {
		inputAssign := test.inputAssign
		inputLet := test.inputLet
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(grammar, transform,
				&Parser{assignTypes: inputAssign, letTypes: inputLet, messages: map[string][]Msg{}})
			validator.validateViews()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateViewsChoiceTypes(t *testing.T) {
	p := NewParser()
	transform, _ := loadTransform("tests", "transform1.sysl", p)
	grammar, _ := loadGrammar("tests/grammar.sysl")

	cases := map[string]struct {
		inputAssign map[string]TypeData
		inputLet    map[string]TypeData
		expected    map[string][]Msg
	}{
		"Valid choice": {
			inputAssign: map[string]TypeData{"ValidChoice": p.GetAssigns()["ValidChoice"]},
			inputLet:    map[string]TypeData{"ValidChoice": p.GetLets()["ValidChoice"]},
			expected:    map[string][]Msg{}},
		"Invalid choice": {
			inputAssign: map[string]TypeData{"InvalidChoice": p.GetAssigns()["InvalidChoice"]},
			inputLet:    map[string]TypeData{"InvalidChoice": p.GetLets()["InvalidChoice"]},
			expected: map[string][]Msg{"InvalidChoice": {
				{MessageID: ErrInvalidOption, MessageData: []string{"InvalidChoice", "Foo", "Statement"}},
				{MessageID: ErrExcessAttr, MessageData: []string{"Foo", "InvalidChoice", "Statement"}}}}},
		"Valid choice combination": {
			inputAssign: map[string]TypeData{"ValidChoiceCombination": p.GetAssigns()["ValidChoiceCombination"]},
			inputLet:    map[string]TypeData{"ValidChoiceCombination": p.GetLets()["ValidChoiceCombination"]},
			expected:    map[string][]Msg{}},
		"Valid choice non-combination": {
			inputAssign: map[string]TypeData{"ValidChoiceNonCombination": p.GetAssigns()["ValidChoiceNonCombination"]},
			inputLet:    map[string]TypeData{"ValidChoiceNonCombination": p.GetLets()["ValidChoiceNonCombination"]},
			expected:    map[string][]Msg{}},
		"Invalid choice combination excess": {
			inputAssign: map[string]TypeData{
				"InvalidChoiceCombinationExcess": p.GetAssigns()["InvalidChoiceCombinationExcess"]},
			inputLet: map[string]TypeData{
				"InvalidChoiceCombinationExcess": p.GetLets()["InvalidChoiceCombinationExcess"]},
			expected: map[string][]Msg{"InvalidChoiceCombinationExcess": {{
				MessageID:   ErrExcessAttr,
				MessageData: []string{"Foo", "InvalidChoiceCombinationExcess", "MethodSpec"}}}}},
		"Invalid choice combination missing": {
			inputAssign: map[string]TypeData{
				"InvalidChoiceCombiMissing": p.GetAssigns()["InvalidChoiceCombiMissing"]},
			inputLet: map[string]TypeData{
				"InvalidChoiceCombiMissing": p.GetLets()["InvalidChoiceCombiMissing"]},
			expected: map[string][]Msg{"InvalidChoiceCombiMissing": {
				{MessageID: ErrMissingReqField,
					MessageData: []string{"Signature", "InvalidChoiceCombiMissing", "MethodSpec"}},
				{MessageID: ErrExcessAttr,
					MessageData: []string{"Foo", "InvalidChoiceCombiMissing", "MethodSpec"}}}}},
		"Invalid choice non-combination missing": {
			inputAssign: map[string]TypeData{
				"InvalidChoiceNonCombination": p.GetAssigns()["InvalidChoiceNonCombination"]},
			inputLet: map[string]TypeData{
				"InvalidChoiceNonCombination": p.GetLets()["InvalidChoiceNonCombination"]},
			expected: map[string][]Msg{"InvalidChoiceNonCombination": {
				{
					MessageID:   ErrInvalidOption,
					MessageData: []string{"InvalidChoiceNonCombination", "Interface", "MethodSpec"}},
				{
					MessageID:   ErrExcessAttr,
					MessageData: []string{"Interface", "InvalidChoiceNonCombination", "MethodSpec"}}}}},
	}
	for name, test := range cases {
		inputAssign := test.inputAssign
		inputLet := test.inputLet
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(grammar, transform, &Parser{
				assignTypes: inputAssign, letTypes: inputLet, messages: map[string][]Msg{}})
			validator.validateViews()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidate(t *testing.T) {
	p := NewParser()
	transform, _ := loadTransform("tests", "transform2.sysl", p)
	grammar, _ := loadGrammar("tests/grammar.sysl")

	validator := NewValidator(grammar, transform, p)
	validator.Validate("goFile")
	actual := validator.GetMessages()
	assert.Equal(t, map[string][]Msg{}, actual, "Unexpected result")
}

func TestValidatorLoadTransformSuccess(t *testing.T) {
	p := NewParser()
	tfm, err := loadTransform("tests", "transform2.sysl", p)
	assert.NotNil(t, tfm, "Unexpected result")
	assert.Nil(t, err, "Unexpected result")
}

func TestValidatorLoadTransformError(t *testing.T) {
	p := NewParser()
	tfm, err := loadTransform("foo", "bar.sysl", p)
	assert.Nil(t, tfm, "Unexpected result")
	assert.NotNil(t, err, "Unexpected result")
}

func TestValidatorLoadGrammarSuccess(t *testing.T) {
	grammar, err := loadGrammar("tests/grammar.sysl")
	assert.NotNil(t, grammar, "Unexpected result")
	assert.Nil(t, err, "Unexpected result")
}

func TestValidatorLoadGrammarError(t *testing.T) {
	grammar, err := loadGrammar("foo/bar.g")
	assert.Nil(t, grammar, "Unexpected result")
	assert.NotNil(t, err, "Unexpected result")
}

func TestValidatorDoValidate(t *testing.T) {
	cases := map[string]struct {
		args     []string
		flags    *flag.FlagSet
		isErrNil bool
	}{
		"Success": {
			args: []string{
				"sysl2", "validate", "-root-transform", "tests", "-transform", "transform2.sysl", "-grammar",
				"tests/grammar.sysl", "-start", "goFile"},
			flags: flag.NewFlagSet("validate", flag.PanicOnError), isErrNil: true},
		"Flag parse fail": {
			args: []string{
				"sysl2", "validate", "-root-transforms", "tests", "-transform", "transform2.sysl", "-grammar",
				"tests/grammar.sysl", "-start", "goFile"},
			flags: flag.NewFlagSet("validate", flag.ContinueOnError), isErrNil: false},
		"Grammar loading fail": {
			args: []string{
				"sysl2", "validate", "-root-transform", "tests", "-transform", "transform2.sysl", "-grammar",
				"tests/go.sysl", "-start", "goFile"},
			flags: flag.NewFlagSet("validate", flag.PanicOnError), isErrNil: false},
		"Transform loading fail": {
			args: []string{
				"sysl2", "validate", "-root-transform", "tests", "-transform", "tfm.sysl", "-grammar",
				"tests/grammar.sysl", "-start", "goFile"},
			flags: flag.NewFlagSet("validate", flag.PanicOnError), isErrNil: false},
		"Has validation messages": {
			args: []string{
				"sysl2", "validate", "-root-transform", "tests", "-transform", "transform1.sysl", "-grammar",
				"tests/grammar.sysl", "-start", "goFile"},
			flags: flag.NewFlagSet("validate", flag.PanicOnError), isErrNil: false},
	}

	for name, test := range cases {
		args := test.args
		flags := test.flags
		isErrNil := test.isErrNil
		t.Run(name, func(t *testing.T) {
			err := DoValidate(flags, args)
			if isErrNil {
				assert.Nil(t, err, "Unexpected result")
			} else {
				assert.NotNil(t, err, "Unexpected result")
			}
		})
	}
}
