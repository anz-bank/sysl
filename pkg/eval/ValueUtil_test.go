package eval

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeValueBool(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_B{B: true}},
		MakeValueBool(true),
	)
	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_B{B: false}},
		MakeValueBool(false),
	)
}

func TestMakeValueI64(t *testing.T) {
	t.Parallel()

	for i, value := range []int64{0, 1, 2, -1, 1<<63 - 1, -(1 << 63)} {
		assert.Equal(t,
			&sysl.Value{Value: &sysl.Value_I{I: value}},
			MakeValueI64(value),
			"%d: %v", i, value,
		)
	}
}

func TestMakeValueString(t *testing.T) {
	t.Parallel()

	for i, value := range []string{"", "x", "😀", "foo-bar"} {
		assert.Equal(t,
			&sysl.Value{Value: &sysl.Value_S{S: value}},
			MakeValueString(value),
			"%d: %v", i, value,
		)
	}
}

func TestMakeValueMap(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{},
		}}},
		MakeValueMap(),
	)
}

func TestMakeValueList(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_List_{List: &sysl.Value_List{
			Value: []*sysl.Value{},
		}}},
		MakeValueList(),
	)
}

func TestMakeValueSet(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Set{Set: &sysl.Value_List{
			Value: []*sysl.Value{},
		}}},
		MakeValueSet(),
	)
}

func TestIsCollectionType(t *testing.T) {
	t.Parallel()

	assert.True(t, IsCollectionType(MakeValueList()))
	assert.True(t, IsCollectionType(MakeValueSet()))
	assert.False(t, IsCollectionType(MakeValueString("1, 2, 3")))
}

func TestAddItemToValueMap(t *testing.T) {
	t.Parallel()

	m := MakeValueMap()
	AddItemToValueMap(m, "key", MakeValueString("value"))
	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"key": {Value: &sysl.Value_S{S: "value"}},
			},
		}}},
		m,
	)
	AddItemToValueMap(m, "key2", MakeValueString("value2"))
	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"key":  {Value: &sysl.Value_S{S: "value"}},
				"key2": {Value: &sysl.Value_S{S: "value2"}},
			},
		}}},
		m,
	)
}

func TestAppendItemToValueList(t *testing.T) {
	t.Parallel()

	m := MakeValueList()
	AppendItemToValueList(m.GetList(), MakeValueI64(42))
	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_List_{List: &sysl.Value_List{
			Value: []*sysl.Value{
				{Value: &sysl.Value_I{I: 42}},
			},
		}}},
		m,
	)
	AppendItemToValueList(m.GetList(), MakeValueString("value"))
	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_List_{List: &sysl.Value_List{
			Value: []*sysl.Value{
				{Value: &sysl.Value_I{I: 42}},
				{Value: &sysl.Value_S{S: "value"}},
			},
		}}},
		m,
	)
}

func TestAttributeToValueString(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		MakeValueString("hello"),
		attributeToValue(&sysl.Attribute{Attribute: &sysl.Attribute_S{S: "hello"}}),
	)
}

func TestAttributeToValueList(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		MakeValueList(),
		attributeToValue(&sysl.Attribute{Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{}},
		}}),
	)
	assert.Equal(t,
		MakeValueList(
			MakeValueString("hello"),
		),
		attributeToValue(&sysl.Attribute{Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{
				{Attribute: &sysl.Attribute_S{S: "hello"}},
			}},
		}}),
	)
}

func TestAttributeToValueOther(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		(*sysl.Value)(nil),
		attributeToValue(&sysl.Attribute{Attribute: &sysl.Attribute_I{I: 42}}),
	)
}

func TestAttributeToValueComnposite(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		MakeValueList(
			MakeValueList(),
			MakeValueString("hello"),
			MakeValueList(
				MakeValueString("goodbye"),
				MakeValueString("thanks for all the fish"),
			),
		),
		attributeToValue(&sysl.Attribute{Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{
				{Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{}}}},
				{Attribute: &sysl.Attribute_S{S: "hello"}},
				{Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{
					{Attribute: &sysl.Attribute_S{S: "goodbye"}},
					{Attribute: &sysl.Attribute_S{S: "thanks for all the fish"}},
				}}}},
			}},
		}}),
	)
}

func assertTypeDetail(t *testing.T, expectedTypeName, expectedTypeDetail string, typ *sysl.Type) {
	typeName, typeDetail := syslutil.GetTypeDetail(typ)
	assert.Equal(t, expectedTypeName, typeName)
	assert.Equal(t, expectedTypeDetail, typeDetail)
}

func TestGetTypeDetailPrimitive(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"primitive", "STRING",
		&sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}},
	)
	assertTypeDetail(t,
		"primitive", "BOOL",
		&sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_BOOL}},
	)
}

func TestGetTypeDetailTypeRef(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"type_ref", "foo",
		&sysl.Type{Type: &sysl.Type_TypeRef{TypeRef: &sysl.ScopedRef{Ref: &sysl.Scope{
			Path: []string{"foo"},
		}}}},
	)
	assertTypeDetail(t,
		"type_ref", "foo",
		&sysl.Type{Type: &sysl.Type_TypeRef{TypeRef: &sysl.ScopedRef{Ref: &sysl.Scope{
			Appname: &sysl.AppName{Part: []string{"foo"}},
		}}}},
	)
}

func TestGetTypeDetailSequence(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"sequence", "STRING",
		&sysl.Type{Type: &sysl.Type_Sequence{
			Sequence: &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}},
		}},
	)
}

func TestGetTypeDetailSet(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"set", "STRING",
		&sysl.Type{Type: &sysl.Type_Set{
			Set: &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}},
		}},
	)
}

func TestGetTypeDetailList(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"list", "STRING",
		&sysl.Type{Type: &sysl.Type_List_{List: &sysl.Type_List{
			Type: &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}},
		}}},
	)
}

func TestGetTypeDetailTuple(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"tuple", "",
		&sysl.Type{Type: &sysl.Type_Tuple_{Tuple: &sysl.Type_Tuple{}}},
	)
}

func TestGetTypeDetailRelation(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"relation", "",
		&sysl.Type{Type: &sysl.Type_Relation_{Relation: &sysl.Type_Relation{}}},
	)
}

func TestGetTypeDetailUnion(t *testing.T) {
	t.Parallel()

	assertTypeDetail(t,
		"union", "",
		&sysl.Type{Type: &sysl.Type_OneOf_{OneOf: &sysl.Type_OneOf{}}},
	)
}

func TestFieldsToValueMapEmpty(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{},
		}}},
		fieldsToValueMap(map[string]*sysl.Type{}),
	)
}

func TestFieldsToValueMapPrimitive(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"a": {Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
					Items: map[string]*sysl.Value{
						"name":      MakeValueString("a"),
						"docstring": MakeValueString("this is a field"),
						"attrs":     MakeValueMap(),
						"type":      MakeValueString("primitive"),
						"primitive": MakeValueString("DATE"),
						"optional":  MakeValueBool(true),
						// "fields":    nil,
					},
				}}},
			},
		}}},
		fieldsToValueMap(map[string]*sysl.Type{
			"a": {
				Docstring: "this is a field",
				Type:      &sysl.Type_Primitive_{Primitive: sysl.Type_DATE},
				Opt:       true,
			},
		}),
	)
}

func TestFieldsToValueMapTuple(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"a": {Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
					Items: map[string]*sysl.Value{
						"name":      MakeValueString("a"),
						"docstring": MakeValueString("this is a field"),
						"attrs":     MakeValueMap(),
						"type":      MakeValueString("tuple"),
						"tuple":     MakeValueString(""),
						"optional":  MakeValueBool(false),
						"fields":    MakeValueMap(),
					},
				}}},
			},
		}}}.Value.(*sysl.Value_Map_).Map.Items["a"],
		fieldsToValueMap(map[string]*sysl.Type{
			"a": {
				Docstring: "this is a field",
				Type:      &sysl.Type_Tuple_{Tuple: &sysl.Type_Tuple{}},
			},
		}).Value.(*sysl.Value_Map_).Map.Items["a"],
	)
}

func TestFieldsToValueMapRelation(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"a": {Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
					Items: map[string]*sysl.Value{
						"name":      MakeValueString("a"),
						"docstring": MakeValueString("this is a field"),
						"attrs":     MakeValueMap(),
						"type":      MakeValueString("relation"),
						"relation":  MakeValueString(""),
						"optional":  MakeValueBool(false),
						"fields":    MakeValueMap(),
					},
				}}},
			},
		}}},
		fieldsToValueMap(map[string]*sysl.Type{
			"a": {
				Docstring: "this is a field",
				Type:      &sysl.Type_Relation_{Relation: &sysl.Type_Relation{}},
			},
		}),
	)
}

func TestStmtToValueAction(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"type":   MakeValueString("action"),
				"action": MakeValueString("doit"),
			},
		}}},
		stmtToValue(
			&sysl.Statement{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "doit"}}},
		),
	)
}

func TestStmtToValueCall(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		&sysl.Value{Value: &sysl.Value_Map_{Map: &sysl.Value_Map{
			Items: map[string]*sysl.Value{
				"type":     MakeValueString("call"),
				"endpoint": MakeValueString("doit"),
				"target":   MakeValueString("Foo :: Bar"),
			},
		}}},
		stmtToValue(
			&sysl.Statement{Stmt: &sysl.Statement_Call{Call: &sysl.Call{
				Target:   &sysl.AppName{Part: []string{"Foo", "Bar"}},
				Endpoint: "doit",
			}}},
		),
	)
}

func TestAddModule(t *testing.T) {
	mod, err := parse.NewParser().ParseFromFs("model_with_deps.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	s := Scope{}
	s.AddModule("module", mod)
	module := s["module"].GetMap().Items
	assert.Equal(t, "ModelWithDeps", module["apps"].GetMap().Items["ModelWithDeps"].GetMap().Items["name"].GetS())
	assert.Equal(t, "Dep", module["apps"].GetMap().Items["Dep"].GetMap().Items["name"].GetS())

	_, hasKey := module["types"]
	assert.Equal(t, true, hasKey)
}

func TestGetValueSlice_NullPanic(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
	}()
	_ = GetValueSlice(MakeValueBool(true))
}

func TestGetValueSlice_ListNoPanic(t *testing.T) {
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()
	_ = GetValueSlice(MakeValueList(MakeValueBool(true)))
}

func TestGetValueSlice_SetNoPanic(t *testing.T) {
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()
	_ = GetValueSlice(MakeValueSet())
}

func TestGetValueSlice_SameList(t *testing.T) {
	boolList := MakeValueList(MakeValueBool(true))
	result := GetValueSlice(boolList)
	assert.Equal(t, &result[0], &boolList.GetList().Value[0])
}

func TestGetValueSlice_SameSet(t *testing.T) {
	boolSet := MakeValueSet()
	boolSet.GetSet().Value = append(boolSet.GetSet().Value, MakeValueBool(true))
	result := GetValueSlice(boolSet)
	assert.Equal(t, &result[0], &boolSet.GetSet().Value[0])
}

func TestTypeToValue(t *testing.T) {
	dataType := sysl.Type{}
	val := TypeToValue(&dataType)
	// assert.Equal(t, &result[0], &boolSet.GetSet().Value[0])
	assert.Equal(t, "", val.GetMap().GetItems()["type"].GetS())
}
