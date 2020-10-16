package arrai

import (
	"context"
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

// ExecutionError encapsulates detailed error msgs from arr.ai runtime.
type ExecutionError struct {
	Context  string
	Err      error
	ShortMsg string
}

func (e ExecutionError) Error() string { return e.Context + ": " + e.Err.Error() }

// EvaluateScript evaluates script with passed parameters.
// It help to pass Go's type parameters to arrai script explicitly.
// TODO: will move it to arrai when it is ready.
func EvaluateScript(arraiScript string, scriptParams ...interface{}) (rel.Value, error) {
	finalScript := fmt.Sprintf("%s%s", arraiScript, toScriptParams(scriptParams...))
	return syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), "", finalScript)
}

// RunBundle runs an arr.ai bundle with the passed parameters set as //os.args[1:].
// It help to pass Go's type parameters to arrai script explicitly.
func EvaluateBundle(bundle []byte, args ...string) (rel.Value, error) {
	args = append([]string{""}, args...)
	return syntax.EvaluateBundle(bundle, args...)
}

func toScriptParams(scriptParams ...interface{}) string {
	var result strings.Builder
	result.WriteString("(")
	for i, param := range scriptParams {
		switch t := param.(type) {
		case string:
			result.WriteString(fmt.Sprintf("`%s`", t))
		case bool, int, int8, int16, int32, int64, uint, uint8, uint16,
			uint32, uint64, uintptr, float32, float64, complex64, complex128:
			result.WriteString(fmt.Sprintf("%v", t))
		default:
			panic(fmt.Sprintf("invalid Go's basic types: %T", param))
		}
		if i+1 < len(scriptParams) {
			result.WriteString(", ")
		}
	}
	result.WriteString(")")

	return result.String()
}
