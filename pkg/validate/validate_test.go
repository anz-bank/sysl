package validate

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/stretchr/testify/assert"
)

const testDir = "../../tests/"

func TestValidatorGetTypeName(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	start := "EntryPoint"
	p := parse.NewParser()
	transform, err := loadTransform("transform1.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

	var entryPointView *sysl.View
	var nonEntryPointView *sysl.View
	var invalidEntryPointView *sysl.View

	entryPointView = transform.Views[start]
	nonEntryPointView = transform.Views["TfmDefaultEmpty"]
	invalidEntryPointView = transform.Views["EntryPointInvalid"]

	cases := map[string]struct {
		input    *sysl.Application
		expected map[string][]msg.Msg
	}{
		"Exists": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{start: entryPointView, "nonEntryPoint": nonEntryPointView}},
			expected: map[string][]msg.Msg{}},
		"Not exists": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{"nonEntryPoint": nonEntryPointView}},
			expected: map[string][]msg.Msg{
				"EntryPoint": {{MessageID: msg.ErrEntryPointUndefined, MessageData: []string{start}}}}},
		"Incorrect output": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{start: invalidEntryPointView, "nonEntryPoint": nonEntryPointView}},
			expected: map[string][]msg.Msg{
				"EntryPoint": {{MessageID: msg.ErrInvalidEntryPointReturn, MessageData: []string{"STRING", start}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(nil, input, &parse.Parser{Messages: map[string][]msg.Msg{}})
			validator.validateEntryPoint(start)
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateFileName(t *testing.T) {
	t.Parallel()

	viewName := "filename"
	p := parse.NewParser()
	transform, err := loadTransform("transform1.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

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
		expected map[string][]msg.Msg
	}{
		"Exists": {
			input: &sysl.Application{
				Views: map[string]*sysl.View{viewName: fileNameView, "nonEntryPoint": nonFileNameView}},
			expected: map[string][]msg.Msg{}},
		"Not exists": {
			input: &sysl.Application{Views: map[string]*sysl.View{"tfmDefaultEmpty": nonFileNameView}},
			expected: map[string][]msg.Msg{"filename": {
				{MessageID: msg.ErrUndefinedView, MessageData: []string{viewName}}}}},
		"Incorrect output": {
			input: &sysl.Application{Views: map[string]*sysl.View{viewName: invalidFileNameView1}},
			expected: map[string][]msg.Msg{"filename": {
				{MessageID: msg.ErrInvalidReturn, MessageData: []string{viewName, "Expected type is string"}}}}},
		"Incorrect assignment": {
			input: &sysl.Application{Views: map[string]*sysl.View{viewName: invalidFileNameView2}},
			expected: map[string][]msg.Msg{"filename": {
				{MessageID: msg.ErrExcessAttr, MessageData: []string{"foo", viewName, "string"}},
				{MessageID: msg.ErrMissingReqField, MessageData: []string{viewName, viewName, "string"}}}}},
		"Excess assignment": {
			input: &sysl.Application{Views: map[string]*sysl.View{viewName: invalidFileNameView3}},
			expected: map[string][]msg.Msg{"filename": {
				{MessageID: msg.ErrExcessAttr, MessageData: []string{"foo", viewName, "string"}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(nil, input, &parse.Parser{Messages: map[string][]msg.Msg{}})
			validator.validateFileName()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateViews(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	transform, err := loadTransform("transform1.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

	grammar, err := LoadGrammar(testDir+"grammar.sysl", afero.NewOsFs())
	require.NoError(t, err)
	require.NotNil(t, grammar)

	cases := map[string]struct {
		input    string
		expected map[string][]msg.Msg
	}{
		"Equal": {
			input:    "TfmValid",
			expected: map[string][]msg.Msg{"TfmValid": nil}},
		"Not Equal": {
			input: "TfmInvalid",
			expected: map[string][]msg.Msg{"TfmInvalid": {
				{MessageID: msg.ErrMissingReqField, MessageData: []string{"FunctionName", "TfmInvalid", "MethodDecl"}}}}},
		"Absent optional": {
			input:    "TfmNoOptional",
			expected: map[string][]msg.Msg{"TfmNoOptional": nil}},
		"Excess attributes without optionals": {
			input: "TfmExcessAttrs1",
			expected: map[string][]msg.Msg{"TfmExcessAttrs1": {
				{MessageID: msg.ErrExcessAttr, MessageData: []string{"ExcessAttr1", "TfmExcessAttrs1", "MethodDecl"}}}}},
		"Excess attributes with optionals": {
			input: "TfmExcessAttrs2",
			expected: map[string][]msg.Msg{"TfmExcessAttrs2": {
				{MessageID: msg.ErrExcessAttr, MessageData: []string{"ExcessAttr1", "TfmExcessAttrs2", "MethodDecl"}}}}},
		"Valid choice": {
			input:    "ValidChoice",
			expected: map[string][]msg.Msg{"ValidChoice": nil}},
		"Relational Type": {
			input:    "Relational",
			expected: map[string][]msg.Msg{"Relational": nil}},
		"Inner relational Type": {
			input:    "InnerRelational",
			expected: map[string][]msg.Msg{"InnerRelational": nil}},
		"Transform variable valid": {
			input:    "TransformVarValid",
			expected: map[string][]msg.Msg{"TransformVarValid": nil}},
		"Transform variable redefined": {
			input: "TransformVarRedefined",
			expected: map[string][]msg.Msg{"TransformVarRedefined": {
				{MessageID: 409, MessageData: []string{"TransformVarRedefined", "varDeclaration"}}}}},
		"Transform inner-variable redefined": {
			input: "TransformInnerVarRedefined",
			expected: map[string][]msg.Msg{"TransformInnerVarRedefined": {
				{MessageID: 409, MessageData: []string{"TransformInnerVarRedefined:varDeclaration", "foo"}}}}},
		"Transform assign redefined": {
			input: "TransformAssignRedefined",
			expected: map[string][]msg.Msg{"TransformAssignRedefined": {
				{MessageID: 409, MessageData: []string{"TransformAssignRedefined", "VarDecl"}}}}},
		"Transform inner-assign redefined": {
			input: "TransformInnerAssignRedefined",
			expected: map[string][]msg.Msg{"TransformInnerAssignRedefined": {
				{MessageID: 409, MessageData: []string{"TransformInnerAssignRedefined:VarDecl", "TypeName"}}}}},
		"Transform variable invalid": {
			input: "TransformVarInvalid:varDeclaration",
			expected: map[string][]msg.Msg{"TransformVarInvalid:varDeclaration": {
				{MessageID: 405, MessageData: []string{"identifier", "TransformVarInvalid:varDeclaration", "VarDecl"}},
				{MessageID: 406, MessageData: []string{"foo", "TransformVarInvalid:varDeclaration", "VarDecl"}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(
				grammar, transform,
				&parse.Parser{
					AssignTypes: map[string]parse.TypeData{input: p.AssignTypes[input]},
					LetTypes:    map[string]parse.TypeData{input: p.LetTypes[input]},
					Messages:    map[string][]msg.Msg{input: p.Messages[input]}})
			validator.validateViews()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateViewsInnerTypes(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	transform, err := loadTransform("transform1.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

	grammar, err := LoadGrammar(testDir+"grammar.sysl", afero.NewOsFs())
	require.NoError(t, err)
	require.NotNil(t, grammar)

	cases := map[string]struct {
		inputAssign map[string]parse.TypeData
		inputLet    map[string]parse.TypeData
		expected    map[string][]msg.Msg
	}{
		"Valid inner type": {
			inputAssign: map[string]parse.TypeData{"ValidInnerAttrs": p.GetAssigns()["ValidInnerAttrs"]},
			inputLet:    map[string]parse.TypeData{"ValidInnerAttrs": p.GetLets()["ValidInnerAttrs"]},
			expected:    map[string][]msg.Msg{}},
		"Invalid inner type": {
			inputAssign: map[string]parse.TypeData{"InvalidInnerAttrs": p.GetAssigns()["InvalidInnerAttrs"]},
			inputLet:    map[string]parse.TypeData{"InvalidInnerAttrs": p.GetLets()["InvalidInnerAttrs"]},
			expected: map[string][]msg.Msg{"InvalidInnerAttrs": {
				{MessageID: msg.ErrMissingReqField, MessageData: []string{"PackageName", "InvalidInnerAttrs", "PackageClause"}},
				{MessageID: msg.ErrExcessAttr, MessageData: []string{"Foo", "InvalidInnerAttrs", "PackageClause"}}}}},
	}
	for name, test := range cases {
		inputAssign := test.inputAssign
		inputLet := test.inputLet
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(grammar, transform,
				&parse.Parser{AssignTypes: inputAssign, LetTypes: inputLet, Messages: map[string][]msg.Msg{}})
			validator.validateViews()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidateViewsChoiceTypes(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	transform, err := loadTransform("transform1.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

	grammar, err := LoadGrammar(testDir+"grammar.sysl", afero.NewOsFs())
	require.NoError(t, err)
	require.NotNil(t, grammar)

	cases := map[string]struct {
		inputAssign map[string]parse.TypeData
		inputLet    map[string]parse.TypeData
		expected    map[string][]msg.Msg
	}{
		"Valid choice": {
			inputAssign: map[string]parse.TypeData{"ValidChoice": p.GetAssigns()["ValidChoice"]},
			inputLet:    map[string]parse.TypeData{"ValidChoice": p.GetLets()["ValidChoice"]},
			expected:    map[string][]msg.Msg{}},
		"Invalid choice": {
			inputAssign: map[string]parse.TypeData{"InvalidChoice": p.GetAssigns()["InvalidChoice"]},
			inputLet:    map[string]parse.TypeData{"InvalidChoice": p.GetLets()["InvalidChoice"]},
			expected: map[string][]msg.Msg{"InvalidChoice": {
				{MessageID: msg.ErrInvalidOption, MessageData: []string{"InvalidChoice", "Foo", "Statement"}},
				{MessageID: msg.ErrExcessAttr, MessageData: []string{"Foo", "InvalidChoice", "Statement"}}}}},
		"Valid choice combination": {
			inputAssign: map[string]parse.TypeData{"ValidChoiceCombination": p.GetAssigns()["ValidChoiceCombination"]},
			inputLet:    map[string]parse.TypeData{"ValidChoiceCombination": p.GetLets()["ValidChoiceCombination"]},
			expected:    map[string][]msg.Msg{}},
		"Valid choice non-combination": {
			inputAssign: map[string]parse.TypeData{"ValidChoiceNonCombination": p.GetAssigns()["ValidChoiceNonCombination"]},
			inputLet:    map[string]parse.TypeData{"ValidChoiceNonCombination": p.GetLets()["ValidChoiceNonCombination"]},
			expected:    map[string][]msg.Msg{}},
		"Invalid choice combination excess": {
			inputAssign: map[string]parse.TypeData{
				"InvalidChoiceCombinationExcess": p.GetAssigns()["InvalidChoiceCombinationExcess"]},
			inputLet: map[string]parse.TypeData{
				"InvalidChoiceCombinationExcess": p.GetLets()["InvalidChoiceCombinationExcess"]},
			expected: map[string][]msg.Msg{"InvalidChoiceCombinationExcess": {{
				MessageID:   msg.ErrExcessAttr,
				MessageData: []string{"Foo", "InvalidChoiceCombinationExcess", "MethodSpec"}}}}},
		"Invalid choice combination missing": {
			inputAssign: map[string]parse.TypeData{
				"InvalidChoiceCombiMissing": p.GetAssigns()["InvalidChoiceCombiMissing"]},
			inputLet: map[string]parse.TypeData{
				"InvalidChoiceCombiMissing": p.GetLets()["InvalidChoiceCombiMissing"]},
			expected: map[string][]msg.Msg{"InvalidChoiceCombiMissing": {
				{MessageID: msg.ErrMissingReqField,
					MessageData: []string{"Signature", "InvalidChoiceCombiMissing", "MethodSpec"}},
				{MessageID: msg.ErrExcessAttr,
					MessageData: []string{"Foo", "InvalidChoiceCombiMissing", "MethodSpec"}}}}},
		"Invalid choice non-combination missing": {
			inputAssign: map[string]parse.TypeData{
				"InvalidChoiceNonCombination": p.GetAssigns()["InvalidChoiceNonCombination"]},
			inputLet: map[string]parse.TypeData{
				"InvalidChoiceNonCombination": p.GetLets()["InvalidChoiceNonCombination"]},
			expected: map[string][]msg.Msg{"InvalidChoiceNonCombination": {
				{
					MessageID:   msg.ErrInvalidOption,
					MessageData: []string{"InvalidChoiceNonCombination", "Interface", "MethodSpec"}},
				{
					MessageID:   msg.ErrExcessAttr,
					MessageData: []string{"Interface", "InvalidChoiceNonCombination", "MethodSpec"}}}}},
	}
	for name, test := range cases {
		inputAssign := test.inputAssign
		inputLet := test.inputLet
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(grammar, transform, &parse.Parser{
				AssignTypes: inputAssign, LetTypes: inputLet, Messages: map[string][]msg.Msg{}})
			validator.validateViews()
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}

func TestValidatorValidate(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	transform, err := loadTransform("transform2.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

	grammar, err := LoadGrammar(testDir+"grammar.sysl", afero.NewOsFs())
	require.NoError(t, err)
	require.NotNil(t, grammar)

	validator := NewValidator(grammar, transform, p)
	validator.Validate("goFile")
	actual := validator.GetMessages()
	assert.Equal(t, map[string][]msg.Msg{}, actual, "Unexpected result")
}

func TestValidatorLoadTransformSuccess(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	tfm, err := loadTransform("transform2.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	assert.NotNil(t, tfm, "Unexpected result")
	assert.Nil(t, err, "Unexpected result")
}

func TestValidatorLoadTransformError(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	tfm, err := loadTransform("bar.sysl", syslutil.NewChrootFs(afero.NewOsFs(), "foo"), p)
	assert.Nil(t, tfm, "Unexpected result")
	assert.NotNil(t, err, "Unexpected result")
}

func TestValidatorLoadGrammarSuccess(t *testing.T) {
	t.Parallel()

	grammar, err := LoadGrammar(testDir+"grammar.sysl", afero.NewOsFs())
	assert.NotNil(t, grammar, "Unexpected result")
	assert.Nil(t, err, "Unexpected result")
}

func TestValidatorLoadGrammarError(t *testing.T) {
	t.Parallel()

	grammar, err := LoadGrammar("foo/bar.g", afero.NewOsFs())
	assert.Nil(t, grammar, "Unexpected result")
	assert.NotNil(t, err, "Unexpected result")
}

func TestValidatorValidateTfmReturn(t *testing.T) {
	t.Parallel()

	p := parse.NewParser()
	transform, err := loadTransform("transform1.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir), p)
	require.NoError(t, err)
	require.NotNil(t, transform)

	grammar, err := LoadGrammar(testDir+"grammar.sysl", afero.NewOsFs())
	require.NoError(t, err)
	require.NotNil(t, grammar)

	cases := map[string]struct {
		input    string
		expected map[string][]msg.Msg
	}{
		"Valid view single obj": {
			input:    "validViewReturnSingleObj",
			expected: map[string][]msg.Msg{"validViewReturnSingleObj": nil}},
		"Valid view seq": {
			input:    "validViewReturnSequence",
			expected: map[string][]msg.Msg{"validViewReturnSequence": nil}},
		"Valid view set": {
			input:    "validViewReturnSet",
			expected: map[string][]msg.Msg{"validViewReturnSet": nil}},
		"Invalid view single obj": {
			input: "invalidViewReturnSingleObj",
			expected: map[string][]msg.Msg{"invalidViewReturnSingleObj": {
				{MessageID: 404, MessageData: []string{"invalidViewReturnSingleObj", "Expected a single VarDecl"}}}}},
		"Invalid view seq": {
			input: "invalidViewReturnSequence",
			expected: map[string][]msg.Msg{"invalidViewReturnSequence": {
				{MessageID: 404, MessageData: []string{"invalidViewReturnSequence", "Expected a sequence of VarDecl"}}}}},
		"Invalid view set": {
			input: "invalidViewReturnSet",
			expected: map[string][]msg.Msg{"invalidViewReturnSet": {
				{MessageID: 404, MessageData: []string{"invalidViewReturnSet", "Expected a set of VarDecl"}}}}},
		"Valid inner tfm single obj": {
			input:    "validInnerTfmReturnSingleObj",
			expected: map[string][]msg.Msg{"validInnerTfmReturnSingleObj": nil}},
		"Valid inner tfm collection": {
			input:    "validInnerTfmReturnCollection",
			expected: map[string][]msg.Msg{"validInnerTfmReturnCollection": nil}},
		"Invalid inner tfm single obj": {
			input: "InvalidInnerTfmReturnSingleObj",
			expected: map[string][]msg.Msg{"InvalidInnerTfmReturnSingleObj": {
				{MessageID: 404, MessageData: []string{"InvalidInnerTfmReturnSingleObj", "Expected a single Statement"}}}}},
		"Invalid inner tfm collection": {
			input: "InvalidInnerTfmReturnCollection",
			expected: map[string][]msg.Msg{"InvalidInnerTfmReturnCollection": {
				{MessageID: 404, MessageData: []string{"InvalidInnerTfmReturnCollection", "Expected a sequence of Statement"}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			validator := NewValidator(
				grammar, transform,
				&parse.Parser{Messages: map[string][]msg.Msg{input: p.Messages[input]}})

			view := transform.GetViews()[input]
			validator.validateTfmReturn(input, view.GetExpr(), view.GetRetType())
			actual := validator.GetMessages()
			assert.Equal(t, expected, actual, "Unexpected result")
		})
	}
}
