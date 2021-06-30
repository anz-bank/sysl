package arrai

import (
	"context"
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/arrai/translate/pb"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/arr-ai/frozen"

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

// SyslPbToValue loads a Sysl protobuf message from a path and serializes it to an arr.ai value.
func SyslPbToValue(pbPath string) (rel.Value, error) {
	m, err := pbutil.FromPB(pbPath, afero.NewOsFs())
	if err != nil {
		return nil, err
	}
	return SyslModuleToValue(m)
}

// SyslPbToValue serializes a Sysl protobuf message to an arr.ai value.
func SyslModuleToValue(module *sysl.Module) (rel.Value, error) {
	return pb.FromProtoValue(protoreflect.ValueOf(module.ProtoReflect()))
}

// EvaluateScript evaluates script with passed parameters.
// It help to pass Go's type parameters to arrai script explicitly.
// TODO: will move it to arrai when it is ready.
func EvaluateScript(arraiScript string, scriptParams ...interface{}) (rel.Value, error) {
	finalScript := fmt.Sprintf("(%s)(%s)", arraiScript, toScriptParams(scriptParams...))
	return syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), "", finalScript)
}

// RunBundle runs an arr.ai bundle with the passed parameters set as //os.args[1:].
// It help to pass Go's type parameters to arrai script explicitly.
func EvaluateBundle(bundle []byte, args ...string) (rel.Value, error) {
	args = append([]string{""}, args...)
	return syntax.EvaluateBundle(bundle, args...)
}

// EvaluateGrammar parses a wbnf grammar from source, uses the grammar (and a root rule) to
// parse a source string, and returns the resulting AST.
func EvaluateGrammar(wbnf, rule, source string) (rel.Value, error) {
	script := fmt.Sprintf(`//grammar.parse({://grammar.lang.wbnf:%s:})`, wbnf)
	return EvaluateScript(script, rule, source)
}

// EvaluateMacro parses a wbnf grammar from source, uses the grammar (and a root rule) to
// parse a source string, transforms the resulting AST with the tx function (arr.ai source) and
// returns the output.
func EvaluateMacro(wbnf, rule, tx, source string) (rel.Value, error) {
	ast, err := EvaluateGrammar(wbnf, rule, source)
	if err != nil {
		return nil, err
	}
	return EvaluateScript(tx, ast)
}

// EvaluateMacroSimple parses a wbnf grammar from source, uses the grammar (and a root rule) to
// parse a source string, transforms the resulting AST with the simpleTransform function from
// arrai/contrib/util, and returns the output.
func EvaluateMacroSimple(wbnf, rule, source string) (rel.Value, error) {
	return EvaluateMacro(wbnf, rule, `//{github.com/arr-ai/arrai/contrib/util}.simpleTransform`, source)
}

// toScriptParams aggregates a list of values into arr.ai source for the arguments of a function
// call, serializing and escaping the parameter values appropriately as needed.
func toScriptParams(scriptParams ...interface{}) string {
	var result strings.Builder
	for i, param := range scriptParams {
		switch t := param.(type) {
		case string:
			result.WriteString(fmt.Sprintf("`%s`", t))
		case bool, int, int8, int16, int32, int64, uint, uint8, uint16,
			uint32, uint64, uintptr, float32, float64, complex64, complex128:
			result.WriteString(fmt.Sprintf("%v", t))
		case rel.Value:
			s, err := syntax.PrettifyString(t, 2)
			if err != nil {
				panic(fmt.Sprintf("failed to serialize arr.ai value: %T", param))
			}
			result.WriteString(s)
		default:
			panic(fmt.Sprintf("invalid type for script param: %T", param))
		}
		if i+1 < len(scriptParams) {
			result.WriteString(", ")
		}
	}
	return result.String()
}

func ToStrings(x interface{}) []string {
	switch xs := x.(type) {
	case nil:
		return nil
	case []string:
		return xs
	case []interface{}:
		ss := make([]string, len(xs))
		for i, x := range xs {
			ss[i] = x.(string)
		}
		return ss
	default:
		panic(fmt.Errorf("not a []string: %T", xs))
	}
}

func ToStringInterfaceMap(x interface{}) map[string]interface{} {
	switch t := x.(type) {
	case nil:
		return nil
	case map[string]interface{}:
		return t
	case frozen.Map:
		m := make(map[string]interface{}, t.Count())
		ctx := context.Background()
		for i := t.Range(); i.Next(); {
			k, v := i.Entry()
			m[k.(rel.Value).String()] = v.(rel.Value).Export(ctx)
		}
		return m
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(t))
		for k, v := range t {
			m[fmt.Sprintf("%s", k)] = v
		}
		return m
	default:
		panic(fmt.Errorf("not a map[interface{}]interface{}: %T", t))
	}
}
