package integrationdiagram

import (
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestVarManagerForComponent(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		Symbols:       map[string]*cmdutils.Var{},
	}

	//When
	result := v.VarManagerForComponent("test", map[string]string{})

	//Then
	assert.Equal(t, "_0", result)
}

func TestVarManagerForComponentWithNameMap(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		Symbols: map[string]*cmdutils.Var{
			"appName": {
				Alias: "_1",
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
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		Symbols: map[string]*cmdutils.Var{
			"test": {
				Alias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForComponent("test", map[string]string{})

	//Then
	assert.Equal(t, "_1", result)
}

func TestVarManagerForEPA(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Endpoints: map[string]*sysl.Endpoint{
						"b": nil,
					},
				},
			},
		},
		DrawableApps: map[string]struct{}{},
		Symbols:      map[string]*cmdutils.Var{},
	}

	//When
	result := v.VarManagerForEPA("a : b")

	//Then
	assert.Equal(t, "_0", result)
}

func TestVarManagerForEPAWithExistingName(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod: &sysl.Module{
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
		DrawableApps: map[string]struct{}{},
		Symbols: map[string]*cmdutils.Var{
			"a : b": {
				Alias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForEPA("a : b")

	//Then
	assert.Equal(t, "_1", result)
}

func TestVarManagerForTopState(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		TopSymbols:    map[string]*_topVar{},
	}

	//When
	result := v.VarManagerForTopState("a : b")

	//Then
	assert.Equal(t, "_0", result)
}

func TestVarManagerForTopStateWithExistingName(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		TopSymbols: map[string]*_topVar{
			"a : b": {
				TopAlias: "_1",
			},
		},
	}

	//When
	result := v.VarManagerForTopState("a : b")

	//Then
	assert.Equal(t, "_1", result)
}

func TestBuildClusterForIntsView(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod: &sysl.Module{
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
		DrawableApps: map[string]struct{}{},
		TopSymbols:   map[string]*_topVar{},
		Symbols:      map[string]*cmdutils.Var{},
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
	v.BuildClusterForEPAView(deps, "")

	//Then
	assert.Equal(t, `state "a" as X_0 {
  state "epa" as _0
  state "epb client" as _1
}
state "b" as X_1 {
  state "epb" as _2
}
`, v.StringBuilder.String())
}

func TestBuildClusterForComponentView(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		TopSymbols:    map[string]*_topVar{},
		Symbols:       map[string]*cmdutils.Var{},
	}
	apps := []string{"a :: A", "a :: A", "b :: B", "c :: C"}

	//When
	v.BuildClusterForIntsView(apps)

	//Then
	assert.Equal(t, `package "a" {
[A] as _0
}
`, v.StringBuilder.String())
}

func TestGenerateIntsView(t *testing.T) {
	t.Parallel()

	//Given
	v := makeIntDiagramVisitor()

	//When
	v.GenerateIntsView(
		&Args{},
		ViewParams{},
		&IntsParam{
			Integrations: []AppDependency{
				{
					Self:   AppElement{Name: "a", Endpoint: "epa"},
					Target: AppElement{Name: "b", Endpoint: "epb"},
				},
			},
			Apps: []string{"a", "b"},
		},
	)

	//Then
	assert.Equal(t, `@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}
[a] as _0
[b] as _1
_0 --> _1 <<indirect>>
[c] as _2
_2 <|.. _1
@enduml`, v.StringBuilder.String())
}

func TestGenerateIntsViewWithCustomAppfmt(t *testing.T) {
	t.Parallel()

	//Given
	v := makeIntDiagramVisitorWithAppfmt("**%(appname)**")

	//When
	v.GenerateIntsView(
		&Args{},
		ViewParams{},
		&IntsParam{
			Integrations: []AppDependency{
				{
					Self:   AppElement{Name: "a", Endpoint: "epa"},
					Target: AppElement{Name: "b", Endpoint: "epb"},
				},
			},
			Apps: []string{"a", "b"},
		},
	)

	//Then
	assert.Equal(t, `@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}
[**a**] as _0
[**b**] as _1
_0 --> _1 <<indirect>>
[**c**] as _2
_2 <|.. _1
@enduml`, v.StringBuilder.String())
}

func testEPAModule() *sysl.Module {
	return &sysl.Module{
		Apps: map[string]*sysl.Application{
			"a": {Endpoints: map[string]*sysl.Endpoint{
				"epa": {
					Attrs: map[string]*sysl.Attribute{"test": nil},
					Stmt: []*sysl.Statement{
						{
							Stmt: &sysl.Statement_Call{Call: &sysl.Call{
								Target:   &sysl.AppName{Part: []string{"b"}},
								Endpoint: "epb",
							}},
						},
					},
				}},
			},
			"b": {
				Endpoints: map[string]*sysl.Endpoint{
					"epb": {Attrs: map[string]*sysl.Attribute{"test": nil}},
				},
			},
			"test": {
				Attrs: map[string]*sysl.Attribute{
					"appfmt": {Attribute: &sysl.Attribute_S{S: "test"}},
				},
			},
		},
	}
}

func TestGenerateEPAView(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           testEPAModule(),
		Project:       "test",
		DrawableApps:  map[string]struct{}{},
		TopSymbols:    map[string]*_topVar{},
		Symbols:       map[string]*cmdutils.Var{},
	}

	//When
	v.GenerateEPAView(
		ViewParams{
			DiagramTitle:       "test",
			HighLightColor:     "blue",
			ArrowColor:         "red",
			IndirectArrowColor: "grey",
		},
		&IntsParam{
			Integrations: []AppDependency{
				{
					Self:   AppElement{Name: "a", Endpoint: "epa"},
					Target: AppElement{Name: "b", Endpoint: "epb"},
					Statement: &sysl.Statement{
						Stmt: &sysl.Statement_Call{},
					},
				},
			},
			App: &sysl.Application{},
		},
	)

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
@enduml`, v.StringBuilder.String())
}

func TestGenerateEPAViewEndpointPattern(t *testing.T) {
	t.Parallel()

	mod := testEPAModule()
	mod.Apps["b"].Endpoints["epb"].Attrs["patterns"] = &sysl.Attribute{
		Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{
			{Attribute: &sysl.Attribute_S{S: "tls"}},
		}}},
	}

	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           mod,
		Project:       "test",
		DrawableApps:  map[string]struct{}{},
		TopSymbols:    map[string]*_topVar{},
		Symbols:       map[string]*cmdutils.Var{},
	}

	//When
	v.GenerateEPAView(
		ViewParams{
			DiagramTitle:       "test",
			HighLightColor:     "blue",
			ArrowColor:         "red",
			IndirectArrowColor: "grey",
		},
		&IntsParam{
			Integrations: []AppDependency{
				{
					Self:   AppElement{Name: "a", Endpoint: "epa"},
					Target: AppElement{Name: "b", Endpoint: "epb"},
					Statement: &sysl.Statement{
						Stmt: &sysl.Statement_Call{},
						Attrs: map[string]*sysl.Attribute{
							"x": {Attribute: &sysl.Attribute_I{I: 42}},
						},
					},
				},
			},
			App: &sysl.Application{},
		},
	)

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
@enduml`, v.StringBuilder.String())
}

func TestGenerateEPAViewSameApp(t *testing.T) {
	t.Parallel()

	mod := testEPAModule()
	mod.Apps["b"].Endpoints["epa"] = proto.Clone(mod.Apps["a"].Endpoints["epa"]).(*sysl.Endpoint)

	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           mod,
		Project:       "test",
		DrawableApps:  map[string]struct{}{},
		TopSymbols:    map[string]*_topVar{},
		Symbols:       map[string]*cmdutils.Var{},
	}

	//When
	v.GenerateEPAView(
		ViewParams{
			DiagramTitle:       "test",
			HighLightColor:     "blue",
			ArrowColor:         "red",
			IndirectArrowColor: "grey",
		},
		&IntsParam{
			Integrations: []AppDependency{
				{
					Self:   AppElement{Name: "b", Endpoint: "epa"},
					Target: AppElement{Name: "b", Endpoint: "epb"},
					Statement: &sysl.Statement{
						Stmt: &sysl.Statement_Call{},
						Attrs: map[string]*sysl.Attribute{
							"x": {Attribute: &sysl.Attribute_I{I: 42}},
							"patterns": {
								Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{
									{Attribute: &sysl.Attribute_S{S: "tls"}},
								}}},
							},
						},
					},
				},
			},
			App: &sysl.Application{},
		},
	)

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
_0 -[#grey]-> _1
@enduml`, v.StringBuilder.String())
}

func TestGenerateView(t *testing.T) {
	t.Parallel()

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
		Integrations: deps,
		App:          &sysl.Application{},
		Endpt: &sysl.Endpoint{
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
[a] as _0
[b] as _1
_0 --> _1 <<indirect>>
@enduml`, result)
}

func TestDrawSystemView(t *testing.T) {
	t.Parallel()

	//Given
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod:           &sysl.Module{},
		DrawableApps:  map[string]struct{}{},
		Symbols: map[string]*cmdutils.Var{
			"test": {
				Alias: "_1",
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
		Integrations: deps,
		App:          &sysl.Application{},
		Endpt: &sysl.Endpoint{
			Attrs: map[string]*sysl.Attribute{
				"epa": nil,
			},
		},
	}
	viewParams := &ViewParams{}
	nameMap := map[string]string{}

	//When
	v.DrawSystemView(*viewParams, params, nameMap)

	//Then
	assert.Equal(t, `[a] as _1
[b] as _2
_1 --> _2 <<indirect>>
`, v.StringBuilder.String())
}

func TestMakeIntsParam(t *testing.T) {
	t.Parallel()

	p := &IntsParam{[]string{"a"},
		map[string]struct{}{},
		[]AppDependency{},
		&sysl.Application{}, &sysl.Endpoint{}}

	assert.NotNil(t, p)
	assert.Equal(t, "a", p.Apps[0])
}

func TestMakeArgs(t *testing.T) {
	t.Parallel()

	a := &Args{"a", "p", true, true}

	assert.NotNil(t, a)
	assert.Equal(t, "a", a.Title)
}

func TestStringInSlice(t *testing.T) {
	t.Parallel()

	s := []string{"a", "b"}

	assert.True(t, StringInSlice("a", s))
}

// makeIntDiagramVisitor returns a populated IntsDiagramVisitor for testing diagram generation.
// The project does not have an appfmt attribute specified.
func makeIntDiagramVisitor() *IntsDiagramVisitor {
	return makeIntDiagramVisitorWithAppfmt("")
}

// makeIntDiagramVisitor returns a populated IntsDiagramVisitor for testing diagram generation.
// If provided, appfmt is set as the project's appfmt attribute.
func makeIntDiagramVisitorWithAppfmt(appfmt string) *IntsDiagramVisitor {
	var stringBuilder strings.Builder
	v := &IntsDiagramVisitor{
		StringBuilder: &stringBuilder,
		Mod: &sysl.Module{
			Apps: map[string]*sysl.Application{
				"a": {
					Name: &sysl.AppName{Part: []string{"a"}},
					Endpoints: map[string]*sysl.Endpoint{
						"epa": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
				},
				"b": {
					Name: &sysl.AppName{Part: []string{"b"}},
					Endpoints: map[string]*sysl.Endpoint{
						"epb": {
							Attrs: map[string]*sysl.Attribute{
								"test": nil,
							},
						},
					},
					Mixin2: []*sysl.Application{
						{Name: &sysl.AppName{Part: []string{"c"}}},
					},
				},
				"c": {
					Name: &sysl.AppName{Part: []string{"c"}},
				},
				"project": {},
			},
		},
		Project:      "project",
		DrawableApps: map[string]struct{}{},
		TopSymbols:   map[string]*_topVar{},
		Symbols:      map[string]*cmdutils.Var{},
	}

	if appfmt != "" {
		v.Mod.Apps["project"].Attrs = map[string]*sysl.Attribute{
			"appfmt": {Attribute: &sysl.Attribute_S{S: "**%(appname)**"}},
		}
	}
	return v
}
