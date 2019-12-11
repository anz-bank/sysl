package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anz-bank/sysl/pkg/parser"
)

func Test_SimpleChoiceConversion(t *testing.T) {
	text := `
          s: a;
        `
	grammar := parser.ParseEBNF(text, "obj", "s")
	require.Equal(t, 1, len(grammar.Rules))

	rp := makeRule(grammar.Rules["s"])
	require.Equal(t, grammar.Rules["s"], rp.rule)

	require.Len(t, rp.types.Items(), 1)
	require.EqualValues(t, []string{"a"}, rp.deps.ToSlice())
}

func Test_SimpleChoiceConversion2(t *testing.T) {
	text := `
          s: a | b;
        `
	grammar := parser.ParseEBNF(text, "obj", "s")
	require.Equal(t, 1, len(grammar.Rules))

	rp := makeRule(grammar.Rules["s"])
	require.Equal(t, grammar.Rules["s"], rp.rule)

	require.Len(t, rp.types.Items(), 1)
	require.EqualValues(t, []string{"a", "b"}, rp.deps.ToSortedSlice())
}

func Test_SimpleChoiceConversion3(t *testing.T) {
	text := `
          s: a+ | b;
        `
	grammar := parser.ParseEBNF(text, "obj", "s")
	require.Equal(t, 1, len(grammar.Rules))

	rp := makeRule(grammar.Rules["s"])
	require.Equal(t, grammar.Rules["s"], rp.rule)

	require.Len(t, rp.types.Items(), 2)
	require.EqualValues(t, []string{"__s_01", "a", "b"}, rp.deps.ToSortedSlice())
}

func Test_SimpleChoiceConversionComplicated(t *testing.T) {
	text := `
          s: (a|b)*|c;
        `
	grammar := parser.ParseEBNF(text, "obj", "s")
	require.Equal(t, 1, len(grammar.Rules))

	rp := makeRule(grammar.Rules["s"])
	require.Equal(t, grammar.Rules["s"], rp.rule)

	require.Len(t, rp.types.Items(), 3)
	require.EqualValues(t, []string{"__s_01", "__s_02", "a", "b", "c"}, rp.deps.ToSortedSlice())
}
