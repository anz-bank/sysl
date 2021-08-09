package transform

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/test"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/wbnf/parser"
	"github.com/spf13/afero"
)

// EvalFileWithParam returns the result of evaluating the function in the specified file, passing in the specified
// parameter.
func EvalFileWithParam(fs afero.Fs, scriptPath string, param rel.Value) (rel.Value, error) {
	expr, err := ExprFileWithParam(fs, scriptPath, param)
	if err != nil {
		return nil, err
	}

	return expr.Eval(arraictx.InitRunCtx(context.Background()), rel.EmptyScope)
}

// EvalWithParam returns the result of evaluating the function in the specified string, passing in the specified
// parameter.
func EvalWithParam(script string, param rel.Value) (rel.Value, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	expr, err := ExprWithParam(script, cwd, param)
	if err != nil {
		return nil, err
	}

	return expr.Eval(arraictx.InitRunCtx(context.Background()), rel.EmptyScope)
}

// ExprFileWithParam returns the unevaluated expression of the function in the specified file after the specified
// parameter was applied to it.
func ExprFileWithParam(fs afero.Fs, scriptPath string, param rel.Value) (rel.Expr, error) {
	scriptBytes, err := afero.ReadFile(fs, scriptPath)
	if err != nil {
		return nil, err
	}

	return ExprWithParam(string(scriptBytes), scriptPath, param)
}

// ExprWithParam returns the unevaluated expression of the function in the specified string after the specified
// parameter was applied to it.
func ExprWithParam(script string, scriptPath string, param rel.Value) (rel.Expr, error) {
	closure, err := syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), scriptPath, script)
	if err != nil {
		return nil, err
	}

	if _, is := closure.(rel.Closure); !is {
		return nil, fmt.Errorf("supplied transform script is not a function")
	}

	return rel.NewCallExpr(*parser.NewScanner(""), closure, param), nil
}

// RunTests runs the test function in the specified string, passing in the specified param, and returns a populated
// test.File with the results.
func RunTests(testScript string, testScriptPath string, param rel.Value) (test.File, error) {
	testFile := test.File{Path: testScriptPath}

	start := time.Now()
	testExpr, err := ExprWithParam(testScript, testScriptPath, param)
	testFile.WallTime = time.Since(start)
	if err == nil {
		testFile.Results, err = test.RunExpr(arraictx.InitRunCtx(context.Background()), testExpr)
	}

	return testFile, err
}
