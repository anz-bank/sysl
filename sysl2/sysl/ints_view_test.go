package main

import (
	"strings"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestVarManagerForComponent(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		symbols:       map[string]*_var{},
	}

	//When
	result := v.VarManagerForComponent("test", map[string]string{})

	//Then
	assert.Equal(t, "_0", result)
}

func TestVarManagerForComponentWithNameMap(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		symbols: map[string]*_var{
			"appName": {
				alias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForComponent("test", map[string]string{
		"test": "appName",
	})

	//Then
	assert.Equal(t, "_1", result)
}

func TestVarManagerForComponentWithExistingName(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		symbols: map[string]*_var{
			"test": {
				alias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForComponent("test", map[string]string{})

	//Then
	assert.Equal(t, "_1", result)
}

func TestVarManagerForEPA(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"b": nil,
					},
				},
			},
		},
		drawableApps: map[string]struct{}{},
		symbols:      map[string]*_var{},
	}

	//When
	result := v.VarManagerForEPA("a : b")

	//Then
	assert.Equal(t, "_0", result)
}

func TestVarManagerForEPAWithExistingName(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"b": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
			},
		},
		drawableApps: map[string]struct{}{},
		symbols: map[string]*_var{
			"a : b": {
				alias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForEPA("a : b")

	//Then
	assert.Equal(t, "_1", result)
}

func TestVarManagerForTopState(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		topSymbols:    map[string]*_topVar{},
	}

	//When
	result := v.VarManagerForTopState("a : b")

	//Then
	assert.Equal(t, "_0", result)
}

func TestVarManagerForTopStateWithExistingName(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		topSymbols: map[string]*_topVar{
			"a : b": {
				topAlias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForTopState("a : b")

	//Then
	assert.Equal(t, "_1", result)
}

func TestBuildClusterForIntsView(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"epa": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
				"b": {
					Endpoints: map[string]*sysl.Endpoint{
						"epb": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
			},
		},
		drawableApps: map[string]struct{}{},
		topSymbols:   map[string]*_topVar{},
		symbols:      map[string]*_var{},
	}
	deps := []AppDependency{
		{
			Self: AppElement{
				Name:     "a",
				Endpoint: "epa",
			},
			Target: AppElement{
				Name:     "b",
				Endpoint: "epb",
			},
		},
	}

	//When
	v.buildClusterForEPAView(deps, "")

	//Then
	assert.Equal(t, `state "" as X_0 {
  state "" as _0
  state "" as _1
}
state "" as X_1 {
  state "" as _2
}
`, v.stringBuilder.String())
}

func TestBuildClusterForComponentView(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		topSymbols:    map[string]*_topVar{},
		symbols:       map[string]*_var{},
	}
	apps := []string{"a :: A", "a :: A", "b :: B", "c :: C"}

	//When
	v.buildClusterForIntsView(apps)

	//Then
	assert.Equal(t, `package "a" {
[] as _0
}
`, v.stringBuilder.String())
}

func TestGenerateIntsView(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	viewParams := &viewParams{}
	deps := []AppDependency{
		{
			Self: AppElement{
				Name:     "a",
				Endpoint: "epa",
			},
			Target: AppElement{
				Name:     "b",
				Endpoint: "epb",
			},
		},
	}
	params := &IntsParam{
		integrations: deps,
	}
	args := &Args{}
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"epa": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
				"b": {
					Endpoints: map[string]*sysl.Endpoint{
						"epb": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
			},
		},
		drawableApps: map[string]struct{}{},
		topSymbols:   map[string]*_topVar{},
		symbols:      map[string]*_var{},
	}

	//When
	v.generateIntsView(args, *viewParams, params)

	//Then
	assert.Equal(t, `@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}
[] as _0
[] as _1
_0 --> _1 <<indirect>>
@enduml`, v.stringBuilder.String())
}

func TestGenerateEPAView(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Call{
				Call: &sysl.Call{
					Target: &sysl.AppName{
						Part: []string{"b"},
					},
					Endpoint: "epb",
				},
			},
		},
	}
	viewParams := &viewParams{
		diagramTitle:       "test",
		highLightColor:     "blue",
		arrowColor:         "red",
		indirectArrowColor: "grey",
	}
	deps := []AppDependency{
		{
			Self: AppElement{
				Name:     "a",
				Endpoint: "epa",
			},
			Target: AppElement{
				Name:     "b",
				Endpoint: "epb",
			},
			Statement: &sysl.Statement{
				Stmt: &sysl.Statement_Call{},
			},
		},
	}
	params := &IntsParam{
		integrations: deps,
		app:          &sysl.Application{},
	}
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"epa": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
							Stmt: stmts,
						},
					},
				},
				"b": {
					Endpoints: map[string]*sysl.Endpoint{
						"epb": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
				"test": {
					Attrs: map[string]*sysl.Attribute{
						"appfmt": {
							Attribute: &sysl.Attribute_S{
								S: "test",
							},
						},
					},
				},
			},
		},
		project:      "test",
		drawableApps: map[string]struct{}{},
		topSymbols:   map[string]*_topVar{},
		symbols:      map[string]*_var{},
	}

	//When
	v.generateEPAView(*viewParams, params)

	//Then
	assert.Equal(t, `@startuml
title test
left to right direction
scale max 16384 height
hide empty description
skinparam state {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
  BackgroundColor<<highlight>> blue
  ArrowColor red
  ArrowColor<<indirect>> grey
}
state "test" as X_0 {
  state "test" as _0
  state "test" as _1
}
state "test" as X_1 {
  state "test" as _2
}
_0 -[#grey]-> _1
_1 -[#black]> _2
@enduml`, v.stringBuilder.String())
}

func TestGenerateStateViewWhenLabelIsNotNil(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Call{
				Call: &sysl.Call{
					Target: &sysl.AppName{
						Part: []string{"b"},
					},
					Endpoint: "epb",
				},
			},
		},
	}
	viewParams := &viewParams{
		diagramTitle: "test",
	}
	deps := map[string]AppDependency{
		"a:epa:b:epb": {
			Self: AppElement{
				Name:     "a",
				Endpoint: "epa",
			},
			Target: AppElement{
				Name:     "b",
				Endpoint: "epb",
			},
			Statement: &sysl.Statement{
				Stmt: &sysl.Statement_Call{},
			},
		},
	}
	params := &IntsParam{
		integrations: deps,
		app: &sysl.Application{
			Attrs: map[string]*sysl.Attribute{
				"epfmt": {
					Attribute: &sysl.Attribute_S{
						S: "label",
					},
				},
			},
		},
	}
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"epa": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
							Stmt: stmts,
						},
					},
				},
				"b": {
					Endpoints: map[string]*sysl.Endpoint{
						"epb": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
				"test": {
					Attrs: map[string]*sysl.Attribute{
						"appfmt": {
							Attribute: &sysl.Attribute_S{
								S: "test",
							},
						},
					},
				},
			},
		},
		project:      "test",
		drawableApps: map[string]struct{}{},
		topSymbols:   map[string]*_topVar{},
		symbols:      map[string]*_var{},
	}

	//When
	v.generateStateView(*viewParams, params)

	//Then
	assert.Equal(t, `@startuml
title test
left to right direction
scale max 16384 height
hide empty description
skinparam state {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}
state "test" as X_0 {
  state "test" as _0
  state "test" as _1
}
state "test" as X_1 {
  state "test" as _2
}
_0 -[#silver]-> _1
_1 -[#black]> _2 : label
@enduml`, v.stringBuilder.String())
}

func TestGenerateView(t *testing.T) {
	//Given
	deps := []AppDependency{
		{
			Self: AppElement{
				Name:     "a",
				Endpoint: "epa",
			},
			Target: AppElement{
				Name:     "b",
				Endpoint: "epb",
			},
		},
	}
	params := &IntsParam{
		integrations: deps,
		app:          &sysl.Application{},
		endpt: &sysl.Endpoint{
			Attrs: map[string]*sysl.Attribute{
				"epa": nil,
			},
		},
	}
	args := &Args{}
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"a": {
				Endpoints: map[string]*sysl.Endpoint{
					"epa": {
						Attrs: map[string]*sysl.Attribute{
							"test": nil,
						},
					},
				},
			},
			"b": {
				Endpoints: map[string]*sysl.Endpoint{
					"epb": {
						Attrs: map[string]*sysl.Attribute{
							"test": nil,
						},
					},
				},
			},
		},
	}

	//When
	result := GenerateView(args, params, m)

	//Then
	assert.Equal(t, `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}
[] as _0
[] as _1
_0 --> _1 <<indirect>>
@enduml`, result)
}

func TestDrawSystemView(t *testing.T) {
	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		stringBuilder: &stringBuilder,
		mod:           &sysl.Module{},
		drawableApps:  map[string]struct{}{},
		symbols: map[string]*_var{
			"test": {
				alias: "_1",
			},
		},
	}
	deps := []AppDependency{
		{
			Self: AppElement{
				Name:     "a",
				Endpoint: "epa",
			},
			Target: AppElement{
				Name:     "b",
				Endpoint: "epb",
			},
		},
	}
	params := &IntsParam{
		integrations: deps,
		app:          &sysl.Application{},
		endpt: &sysl.Endpoint{
			Attrs: map[string]*sysl.Attribute{
				"epa": nil,
			},
		},
	}
	viewParams := &viewParams{}
	nameMap := map[string]string{}

	//When
	v.drawSystemView(*viewParams, params, nameMap)

	//Then
	assert.Equal(t, `[] as _1
[] as _2
_1 --> _2 <<indirect>>
`, v.stringBuilder.String())

}

func TestMakeIntsParam(t *testing.T) {
	p := &IntsParam{[]string{"a"},
		map[string]struct{}{},
		[]AppDependency{},
		&sysl.Application{}, &sysl.Endpoint{}}

	assert.NotNil(t, p)
	assert.Equal(t, "a", p.apps[0])
}

func TestMakeArgs(t *testing.T) {
	a := &Args{"a", "p", true, true}

	assert.NotNil(t, a)
	assert.Equal(t, "a", a.title)
}

func TestStringInSlice(t *testing.T) {
	s := []string{"a", "b"}

	assert.True(t, stringInSlice("a", s))
}
