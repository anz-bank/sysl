package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anz-bank/sysl/pkg/sysl"
)

type scope struct {
	appname []string
	path    []string
}
type typereftestdata struct {
	scope    scope
	appname  []string
	path     []string
	expected string // will be compared against the Type.Docstring or empty
}

func (t typereftestdata) testName() string {
	return fmt.Sprintf("%s.%s", strings.Join(t.appname, "::"), strings.Join(t.path, "."))
}

func (t typereftestdata) makeRef() *sysl.ScopedRef {
	return &sysl.ScopedRef{
		Context: &sysl.Scope{
			Appname: &sysl.AppName{
				Part: t.scope.appname,
			},
			Path: t.scope.path,
		},
		Ref: &sysl.Scope{
			Appname: &sysl.AppName{
				Part: t.appname,
			},
			Path: t.path,
		},
	}
}

func (t typereftestdata) Run(tt *testing.T, mod *sysl.Module) {
	tt.Run(t.testName(), func(tt *testing.T) {
		ref := t.makeRef()
		untouched := ref.String()

		res := resolveTypeRef(mod, ref)
		if t.expected == "" {
			assert.Nil(tt, res)
		} else {
			assert.NotNil(tt, res)
			if res != nil {
				assert.Equal(tt, t.expected, res.Docstring)
			}
		}
		// Ensure the ref hasn't changed
		assert.Equal(tt, untouched, ref.String())
	})
}

func Test_ResolveTypeRef_NoScope(t *testing.T) {
	mod := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"A": {
				Types: map[string]*sysl.Type{
					"Foo":     {Docstring: "A::Foo"},
					"Foo.Bar": {Docstring: "A::Foo.Bar"},
					"Bar":     {Docstring: "A::Bar"},
				},
			},
		},
	}

	for _, test := range []typereftestdata{
		{appname: []string{"A"}, path: []string{"Foo"}, expected: "A::Foo"},
		{appname: []string{"A"}, path: []string{"Foo.Bar"}, expected: "A::Foo.Bar"},
		{appname: []string{"A"}, path: []string{"Bar"}, expected: "A::Bar"},
		{appname: []string{"A"}, path: []string{"Bingo"}},
		{appname: []string{""}, path: []string{"Foo"}},
	} {
		test.Run(t, mod)
	}
}

func Test_ResolveTypeRef_SimpleScope(t *testing.T) {
	mod := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"A": {
				Types: map[string]*sysl.Type{
					"Foo":     {Docstring: "A::Foo"},
					"Foo.Bar": {Docstring: "A::Foo.Bar"},
					"Bar":     {Docstring: "A::Bar"},
				},
			},
		},
	}

	for _, test := range []typereftestdata{
		{scope: scope{appname: []string{"A"}, path: []string{}},
			appname: []string{"A"}, path: []string{"Foo"}, expected: "A::Foo"},
		{scope: scope{appname: []string{"A"}, path: []string{}},
			appname: []string{}, path: []string{"Foo"}, expected: "A::Foo"},
		{scope: scope{appname: []string{"B"}, path: []string{}},
			appname: []string{"A"}, path: []string{"Foo"}, expected: "A::Foo"},
		{scope: scope{appname: []string{"A"}, path: []string{"Foo"}},
			appname: nil, path: []string{"Bar"}, expected: "A::Foo.Bar"},
	} {
		test.Run(t, mod)
	}
}

func Test_ResolveTypeRef_ComplexScope(t *testing.T) {
	mod := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"A": {
				Types: map[string]*sysl.Type{
					"Foo":     {Docstring: "A::Foo"},
					"Foo.Bar": {Docstring: "A::Foo.Bar"},
					"Bar":     {Docstring: "A::Bar"},
				},
			},
			"A :: B": {
				Types: map[string]*sysl.Type{
					"Foo":     {Docstring: "A::B::Foo"},
					"Foo.Bar": {Docstring: "A::B::Foo.Bar"},
					"Bar":     {Docstring: "A::B::Bar"},
				},
			},
			"A :: B :: C": {
				Types: map[string]*sysl.Type{
					"Foo":     {Docstring: "A::B::C::Foo"},
					"Foo.Bar": {Docstring: "A::B::C::Foo.Bar"},
					"Bar":     {Docstring: "A::B::C::Bar"},
				},
			},
		},
	}

	for _, test := range []typereftestdata{
		{scope: scope{appname: []string{"A", "B", "C"}, path: []string{}},
			appname: []string{}, path: []string{"Foo"}, expected: "A::B::C::Foo"},
		{scope: scope{appname: []string{"A", "B", "C"}, path: []string{"X", "Y"}},
			appname: []string{}, path: []string{"Foo"}, expected: "A::B::C::Foo"},
		{scope: scope{appname: []string{"A", "B"}, path: []string{"Foo"}},
			appname: []string{}, path: []string{"Bar"}, expected: "A::B::Foo.Bar"},
		{scope: scope{appname: []string{"A", "B"}, path: []string{"Foo"}},
			appname: []string{"B"}, path: []string{"Bar"}, expected: "A::B::Bar"},
	} {
		test.Run(t, mod)
	}
}

func Test_ResolveTypeRef_DoesntMutateRef(t *testing.T) {

	mod := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"Model": {
				Types: map[string]*sysl.Type{
					"Error":    {Docstring: "1"},
					"Response": {Docstring: "2"},
					"Bar":      {Docstring: "A::Bar"},
				},
			}},
	}
	for _, test := range []typereftestdata{
		{scope: scope{appname: []string{"Model"}, path: []string{"Response"}},
			appname: nil, path: []string{"Error"}, expected: "1"},
	} {
		test.Run(t, mod)
	}
}
