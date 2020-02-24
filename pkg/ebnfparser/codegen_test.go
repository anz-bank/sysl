package ebnfparser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/alecthomas/assert"

	"github.com/anz-bank/sysl/pkg/eval"
)

func loadGrammarFromText(t *testing.T, text string, start string) *EbnfGrammar {
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "grammar.g", []byte(text), 0644))
	g, err := ReadGrammar(fs, "grammar.g", start)
	require.NoError(t, err)
	require.NotNil(t, g)
	return g
}

func TestGenerateOutput_Simple(t *testing.T) {
	g := loadGrammarFromText(t, `a: 'a' b 'x'; `, "a")
	out, err := GenerateOutput(g, eval.MakeValueString("hello"), nil)
	assert.NoError(t, err)
	assert.EqualValues(t, "a hello x", out)
}

func TestGenerateOutput_NoStartRule(t *testing.T) {
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "grammar.g", []byte(`a: 'a' b 'x'; `), 0644))
	g, err := ReadGrammar(fs, "grammar.g", "start")
	require.Error(t, err)
	require.Nil(t, g)
}

func TestGenerateOutput_Map(t *testing.T) {
	g := loadGrammarFromText(t, `a: 'a' b 'x'; `, "a")
	v := &sysl.Value_Map{
		Items: map[string]*sysl.Value{
			"b": eval.MakeValueString("hello"),
		},
	}
	out, err := GenerateOutput(g, &sysl.Value{
		Value: &sysl.Value_Map_{Map: v},
	}, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, "a hello x", out)
}

func TestGenerateOutput_Choices(t *testing.T) {
	g := loadGrammarFromText(t, `a: 'a' (b|c) 'x'; `, "a")
	v := &sysl.Value_Map{
		Items: map[string]*sysl.Value{
			"b": eval.MakeValueString("foo"),
			"c": eval.MakeValueString("hello"),
		},
	}
	out, err := GenerateOutput(g, &sysl.Value{
		Value: &sysl.Value_Map_{Map: v},
	}, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, "a foo x", out)
}

func makeList(items ...*sysl.Value) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_List_{List: &sysl.Value_List{
			Value: items,
		},
		},
	}
}

func TestGenerateOutput_Quants(t *testing.T) {
	type data struct {
		name      string
		g         string
		items     []*sysl.Value
		expected  string
		expectErr bool
	}
	for _, test := range []data{
		{"empty star", "a: b*;", []*sysl.Value{}, "", false},
		{"one star", "a: b*;", []*sysl.Value{eval.MakeValueString("foo")}, "foo", false},
		{"many star", "a: b*;", []*sysl.Value{eval.MakeValueString("a"),
			eval.MakeValueString("b")}, "a b", false},

		{"empty plus", "a: b+;", []*sysl.Value{}, "", true},
		{"one plus", "a: b+;", []*sysl.Value{eval.MakeValueString("foo")}, "foo", false},
		{"many plus", "a: b+;", []*sysl.Value{eval.MakeValueString("a"),
			eval.MakeValueString("b")}, "a b", false},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			v := &sysl.Value_Map{
				Items: map[string]*sysl.Value{
					"b": makeList(test.items...),
				},
			}
			g := loadGrammarFromText(t, test.g, "a")
			out, err := GenerateOutput(g, &sysl.Value{
				Value: &sysl.Value_Map_{Map: v},
			}, nil)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, test.expected, out)
			}
		})
	}
}

// the existing codegen uses the value 'true' to indicate that a type exists but is empty
func TestGenerateOutput_WantTypeGetBoolHack(t *testing.T) {
	g := loadGrammarFromText(t, `
Signature: Parameters Result?;
Parameters: '(' ParameterList? ')';
ParameterList     : ParameterDecl ParameterDeclC*;
ParameterDecl  : Identifier TypeName;
ParameterDeclC: ',' ParameterDecl;
`, "Signature")
	v := &sysl.Value_Map{
		Items: map[string]*sysl.Value{
			"Parameters": eval.MakeValueBool(true),
		},
	}
	out, err := GenerateOutput(g, &sysl.Value{
		Value: &sysl.Value_Map_{Map: v},
	}, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, "( )", out)
}
