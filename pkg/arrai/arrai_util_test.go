package arrai

import (
	"testing"

	"github.com/arr-ai/arrai/rel"

	"github.com/stretchr/testify/assert"
)

func TestToScriptParams(t *testing.T) {
	assert.Equal(t, "(`param1`, `param2`, `param3`)", toScriptParams("param1", "param2", "param3"))
	assert.Equal(t, "(23, 45, true)", toScriptParams(23, 45, true))
	assert.Equal(t, "(`23`, 45, true)", toScriptParams("23", 45, true))
}

func TestEvaluateScript(t *testing.T) {
	val, err := EvaluateScript("let calc = \\number number + 1;calc", 2)
	assert.NoError(t, err)
	assert.Equal(t, val, rel.NewNumber(3))
}
