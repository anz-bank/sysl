package relmod

import (
	"context"
	"fmt"
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/syntax"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/require"
)

var testAttrs = map[string]*sysl.Attribute{
	"patterns": {Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{
		Elt: []*sysl.Attribute{
			{Attribute: &sysl.Attribute_S{S: "tag1"}},
			{Attribute: &sysl.Attribute_S{S: "tag2"}},
		},
	}}},
	"string": {Attribute: &sysl.Attribute_S{S: "s"}},
	"array": {Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{
		Elt: []*sysl.Attribute{
			{Attribute: &sysl.Attribute_S{S: "a1"}},
			{Attribute: &sysl.Attribute_S{S: "a2"}},
		},
	}}},
}

func TestTags(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"tag1", "tag2"}, tags(testAttrs))
}

func TestAnnos(t *testing.T) {
	t.Parallel()

	expected := map[string]interface{}{
		"array":  rel.NewArray(rel.NewString([]rune("a1")), rel.NewString([]rune("a2"))),
		"string": rel.NewString([]rune("s")),
	}

	assert.Equal(t, expected, annos(testAttrs))
}

func TestParseReturnPayload(t *testing.T) {
	t.Parallel()

	ctx, err := withPayloadParser(context.Background())
	require.NoError(t, err)

	mustParseReturnPayload := func(p string) StatementReturn {
		r, err := parseReturnPayload(ctx, p)
		require.NoError(t, err)
		return r
	}

	assert.Equal(t, StatementReturn{Status: "ok"}, mustParseReturnPayload("ok"))

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypePrimitive{"any"}},
		mustParseReturnPayload("ok <: list of unmatched"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypePrimitive{"int"}},
		mustParseReturnPayload("ok <: int"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypeSet{TypePrimitive{"int"}}},
		mustParseReturnPayload("ok <: set of int"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypeRef{TypePath: []string{"Type"}}},
		mustParseReturnPayload("ok <: Type"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypeSequence{TypeRef{TypePath: []string{"Type"}}}},
		mustParseReturnPayload("ok <: sequence of Type"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypeRef{
			AppName:  []string{"App"},
			TypePath: []string{"Type"},
		}},
		mustParseReturnPayload("ok <: App.Type"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypeRef{
			AppName:  []string{"Namespace", "App"},
			TypePath: []string{"Type"},
		}},
		mustParseReturnPayload("ok <: Namespace::App.Type"),
	)

	assert.Equal(t,
		StatementReturn{Status: "ok", Type: TypeRef{
			AppName:  []string{"A"},
			TypePath: []string{"B"},
		}},
		mustParseReturnPayload("ok <: A.B"),
	)

	assert.Equal(t,
		StatementReturn{
			Status: "ok",
			Type:   TypePrimitive{"int"},
			Attr: StatementReturnAttrs{
				Modifier: []string{"tag"},
				Nvp: map[interface{}]interface{}{
					"k":  "v",
					"ak": map[string]interface{}{"a": []interface{}{"a1", "a2"}},
				},
			},
		},
		mustParseReturnPayload(`ok <: int [~tag, k="v", ak=["a1", "a2"]]`),
	)

	assert.Equal(t,
		StatementReturn{
			Status: "ok",
			Type:   TypePrimitive{"string"},
			Attr: StatementReturnAttrs{
				Nvp: map[interface{}]interface{}{
					"annotation": map[string]interface{}{
						"a": []interface{}{"as", "an", "array"},
					},
				},
			},
		},
		mustParseReturnPayload(`ok <: string [annotation=["as", "an", "array"]]`),
	)

	assert.Equal(t,
		StatementReturn{
			Status: "ok",
			Type:   TypePrimitive{"string"},
			Attr: StatementReturnAttrs{
				Nvp: map[interface{}]interface{}{
					"annotation": map[string]interface{}{
						"a": []interface{}{
							map[string]interface{}{"a": []interface{}{"or", "as", "an"}},
							map[string]interface{}{"a": []interface{}{"array", "of", "arrays"}},
						},
					},
				},
			},
		},
		mustParseReturnPayload(`ok <: string [annotation=[["or", "as", "an"], ["array", "of", "arrays"]]]`),
	)
}

func TestReturn(t *testing.T) {
	grammar := `
		type		-> PRIMITIVE | ref;
		ref			-> (app=(\w+):"::" ".")? type=\w+;
		PRIMITIVE	-> 'int' | 'int32' | 'int64' | 'float' | 'float32' | 'float64' | 'decimal'
					 | 'bool' | 'bytes' | 'string' | 'date' | 'datetime' | 'any';
	`
	tx := `\ast cond ast {
		(:PRIMITIVE, ...): (primitive: PRIMITIVE.'' rank (:.@)),
		(:ref, ...): (
			appName: ref.app?:[] >> (.'' rank (:.@)),
			typePath: [ref.type.'' rank (:.@)],
		),
	}`

	actual, err := arrai.EvaluateMacro(grammar, "type", tx, "int")
	require.NoError(t, err)
	rel.AssertEqualValues(t, eval(`(primitive: "int")`), actual)

	actual, err = arrai.EvaluateMacro(grammar, "type", tx, "Type")
	require.NoError(t, err)
	rel.AssertEqualValues(t, eval(`(appName: {}, typePath: ["Type"])`), actual)

	actual, err = arrai.EvaluateMacro(grammar, "type", tx, "App.Type")
	require.NoError(t, err)
	rel.AssertEqualValues(t, eval(`(appName: ["App"], typePath: ["Type"])`), actual)

	actual, err = arrai.EvaluateMacro(grammar, "type", tx, "Namespace::App.Type")
	require.NoError(t, err)
	rel.AssertEqualValues(t, eval(`(appName: ["Namespace", "App"], typePath: ["Type"])`), actual)
}

func eval(src string) rel.Value {
	v, err := syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), "", src)
	if err != nil {
		panic(fmt.Errorf("invalid arr.ai source: %s (%s)", src, err))
	}
	return v
}
