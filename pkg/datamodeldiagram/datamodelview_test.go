package datamodeldiagram

import (
	"strings"
	"testing"

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

func TestDrawPrimitive(t *testing.T) {
	clMock := new(ClassLabelerMock)
	clMock.On("LabelClass", "uuid").Return("test")

	var stringBuilder strings.Builder
	v := MakeDataModelView(clMock, nil, &stringBuilder, "title", "project")
	viewParam := EntityViewParam{
		EntityColor:  "orchid",
		EntityHeader: "D",
		EntityName:   "uuid",
	}
	relationshipMap := map[string]map[string]RelationshipParam{}
	v.DrawPrimitive(viewParam, "INT", relationshipMap)
	actual := v.StringBuilder.String()

	expected := "class \"uuid\" as _0 << (D,orchid) int >> {\n" + "}\n"
	assert.EqualValues(t, expected, actual, nil)
	clMock.AssertExpectations(t)
}
