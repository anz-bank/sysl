package main

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestAppDependency_String(t *testing.T) {
	// Given
	stmt := &sysl.Statement{}
	dep := &AppDependency{
		Self:      AppElement{Name: "AppA", Endpoint: "EndptA"},
		Target:    AppElement{Name: "AppB", Endpoint: "EndptB"},
		Statement: stmt,
	}
	expected := "AppA:EndptA:AppB:EndptB"

	// When
	actual := dep.String()

	// Then
	assert.Equal(t, expected, actual)
}
