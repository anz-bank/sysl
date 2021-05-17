package arrai

import (
	"context"
	"fmt"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/syntax"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"

	"github.com/stretchr/testify/assert"
)

const dateGrammar = `date -> y=\d{4} "-" m=\d{2} "-" d=\d{2};`

func TestToScriptParams(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "`param1`, `param2`, `param3`", toScriptParams("param1", "param2", "param3"))
	assert.Equal(t, "23, 45, true", toScriptParams(23, 45, true))
	assert.Equal(t, "`23`, 45, true", toScriptParams("23", 45, true))
}

func TestEvaluateScript(t *testing.T) {
	t.Parallel()

	val, err := EvaluateScript(`let increment = \n n + 1; increment`, 2)
	assert.NoError(t, err)
	assert.Equal(t, val, rel.NewNumber(3))
}

func TestEvaluateGrammar(t *testing.T) {
	t.Parallel()

	actual, err := EvaluateGrammar(dateGrammar, "date", "2020-06-09")
	require.NoError(t, err)

	expected := arrai(`('': [4\'-', 7\'-'], @rule: 'date', d: ('': 8\'09'), m: ('': 5\'06'), y: ('': '2020'))`)

	rel.AssertEqualValues(t, expected, actual)
}

func TestEvaluateMacro(t *testing.T) {
	t.Parallel()

	tx := `\ast ast -> (year: .y, month: .m, day: .d) :> //eval.value(.'')`

	actual, err := EvaluateMacro(dateGrammar, "date", tx, "2020-06-09")
	require.NoError(t, err)

	rel.AssertEqualValues(t, arrai(`(year: 2020, month: 6, day: 9)`), actual)
}

func TestEvaluateMacroSimple(t *testing.T) {
	t.Parallel()

	actual, err := EvaluateMacroSimple(dateGrammar, "date", "2020-06-09")
	require.NoError(t, err)

	rel.AssertEqualValues(t, arrai(`(y: '2020', m: 5\'06', d: 8\'09')`), actual)
}

func arrai(src string) rel.Value {
	v, err := syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), "", src)
	if err != nil {
		panic(fmt.Errorf("invalid arr.ai source: %s (%s)", src, err))
	}
	return v
}
