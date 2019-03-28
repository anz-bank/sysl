package main

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetAppName(t *testing.T) {
	assert.Equal(t,
		"test :: name",
		getAppName(&sysl.AppName{Part: []string{"test", "name"}}),
	)
}

func TestGetApp(t *testing.T) {
	app := &sysl.Application{Attrs: map[string]*sysl.Attribute{}}
	assert.Equal(t,
		app,
		getApp(
			&sysl.AppName{Part: []string{"test", "name"}},
			&sysl.Module{Apps: map[string]*sysl.Application{"test :: name": app}},
		),
	)
}

func TestPattern(t *testing.T) {
	attrs := map[string]*sysl.Attribute{
		"patterns": {Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{
			{Attribute: &sysl.Attribute_S{S: "abstract"}},
			{Attribute: &sysl.Attribute_S{S: "human"}},
		}}}},
	}
	assert.True(t, HasPattern(attrs, "abstract"))
	assert.True(t, HasPattern(attrs, "human"))
	assert.False(t, HasPattern(attrs, "abstracthuman"))
}

func TestIsNotSameAppWithPartLength(t *testing.T) {
	assert.False(t, isSameApp(
		&sysl.AppName{Part: []string{"test", "name"}},
		&sysl.AppName{Part: []string{"name1"}},
	))
}

func TestIsNotSameAppWithPartContent(t *testing.T) {
	assert.False(t, isSameApp(
		&sysl.AppName{Part: []string{"test", "name"}},
		&sysl.AppName{Part: []string{"test", "name1"}},
	))
}

func TestIsSameApp(t *testing.T) {
	assert.True(t, isSameApp(
		&sysl.AppName{Part: []string{"test", "name"}},
		&sysl.AppName{Part: []string{"test", "name"}},
	))
}

func TestIsSameCall(t *testing.T) {
	assert.True(t, isSameCall(
		&sysl.Call{Target: &sysl.AppName{Part: []string{"test", "name"}}, Endpoint: "endpt"},
		&sysl.Call{Target: &sysl.AppName{Part: []string{"test", "name"}}, Endpoint: "endpt"},
	))
}
