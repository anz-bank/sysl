package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		entityColor:  "orchid",
		entityHeader: "D",
		entityName:   "uuid",
	}
	relationshipMap := map[string]map[string]RelationshipParam{}
	v.drawPrimitive(viewParam, "INT", relationshipMap)
	actual := v.stringBuilder.String()

	expected := "class \"uuid\" as _0 << (D,orchid) >> {\n" +
		"+ id : int\n}\n"
	assert.EqualValues(t, expected, actual, nil)
}
