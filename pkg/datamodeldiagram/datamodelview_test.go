package datamodeldiagram

import (
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type ClassLabelerMock struct {
	mock.Mock
}

func (mock *ClassLabelerMock) LabelClass(className string) string {
	args := mock.Called(className)
	return args.String(0)
}

func assertDraw(t *testing.T,
	expected, className string,
	drawFunc func(*DataModelView, map[string]map[string]RelationshipParam)) {
	clMock := new(ClassLabelerMock)
	clMock.On("LabelClass", className).Return("test")

	var stringBuilder strings.Builder
	v := MakeDataModelView(clMock, nil, &stringBuilder, "title", "project")
	relationshipMap := map[string]map[string]RelationshipParam{}
	drawFunc(v, relationshipMap)
	actual := v.StringBuilder.String()

	assert.EqualValues(t, expected, actual, nil)
	clMock.AssertExpectations(t)
}

func TestDrawPrimitive(t *testing.T) {
	t.Parallel()
	assertDraw(t,
		"class \"uuid\" as _0 << (D,orchid) int >> {\n}\n",
		"uuid",
		func(v *DataModelView, relationshipMap map[string]map[string]RelationshipParam) {
			v.DrawPrimitive(
				EntityViewParam{
					EntityColor:  "orchid",
					EntityHeader: "D",
					EntityName:   "uuid",
				},
				"INT",
				relationshipMap,
			)
		},
	)
}

func TestDrawTuple(t *testing.T) {
	t.Parallel()
	tuple := syslwrapper.MakeTuple(
		map[string]*sysl.Type{"attr1": syslwrapper.MakePrimitive("string")},
	).Type.(*sysl.Type_Tuple_).Tuple
	assertDraw(t,
		`class "aliasName" as _0 << (D,orchid) typeName >> {
+ attr1 : string
}
`,
		"typeName",
		func(v *DataModelView, relationshipMap map[string]map[string]RelationshipParam) {
			v.DrawTuple(
				EntityViewParam{
					EntityColor:  "orchid",
					EntityHeader: "D",
					EntityName:   "typeName",
					EntityAlias:  "aliasName",
				},
				tuple, relationshipMap,
			)
		},
	)

	assertDraw(t,
		`class "typeName" as _0 << (D,orchid) >> {
+ attr1 : string
}
`,
		"typeName",
		func(v *DataModelView, relationshipMap map[string]map[string]RelationshipParam) {
			v.DrawTuple(
				EntityViewParam{
					EntityColor:  "orchid",
					EntityHeader: "D",
					EntityName:   "typeName",
					EntityAlias:  "",
				},
				tuple, relationshipMap,
			)
		},
	)

	assertDraw(t,
		`class "typeName" as _0 << (D,orchid) >> {
+ attr1 : **Sequence <string>**
+ attr2 : **Set <int>**
+ attr3 : **List <bool>**
}
`,
		"typeName",
		func(v *DataModelView, relationshipMap map[string]map[string]RelationshipParam) {
			v.DrawTuple(
				EntityViewParam{
					EntityColor:  "orchid",
					EntityHeader: "D",
					EntityName:   "typeName",
					EntityAlias:  "",
				},
				syslwrapper.MakeTuple(
					map[string]*sysl.Type{
						"attr1": syslwrapper.MakeSequence(syslwrapper.MakePrimitive("string")),
						"attr2": syslwrapper.MakeSet(syslwrapper.MakePrimitive("int")),
						"attr3": syslwrapper.MakeList(syslwrapper.MakePrimitive("bool")),
					},
				).Type.(*sysl.Type_Tuple_).Tuple,
				relationshipMap,
			)
		},
	)
}
