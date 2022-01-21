package transform

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/anz-bank/sysl/pkg/arrai/relmod"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/test"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/wbnf/parser"
)

var errNotClosure = fmt.Errorf("supplied transform script is not a function")

// EvalWithParam returns the result of evaluating the function in the specified bytes, passing in the specified
// parameter.
func EvalWithParam(scriptBytes []byte, scriptPath string, param rel.Value) (rel.Value, error) {
	expr, err := exprWithParam(scriptBytes, scriptPath, param)
	if err != nil {
		return nil, err
	}

	return expr.Eval(arraictx.InitRunCtx(context.Background()), rel.EmptyScope)
}

// exprWithParam returns the unevaluated expression of the function in the specified file content, which may be a
// textual script (.arrai) or a bundle (.arraiz), after the specified parameter was applied to it.
func exprWithParam(scriptBytes []byte, scriptPath string, param rel.Value) (rel.Expr, error) {
	if http.DetectContentType(scriptBytes) == "application/zip" {
		closure, err := syntax.EvaluateBundle(scriptBytes)
		if err != nil {
			return nil, err
		}
		if _, is := closure.(rel.Closure); !is {
			return nil, errNotClosure
		}
		return rel.NewCallExpr(*parser.NewScanner(""), closure, param), nil
	}

	closure, err := syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), scriptPath, string(scriptBytes))
	if err != nil {
		return nil, err
	}

	if _, is := closure.(rel.Closure); !is {
		return nil, errNotClosure
	}

	return rel.NewCallExpr(*parser.NewScanner(""), closure, param), nil
}

// RunTests runs the test function in the specified string, passing in the specified param, and returns a populated
// test.File with the results.
func RunTests(scriptBytes []byte, testScriptPath string, param rel.Value) (test.File, error) {
	testFile := test.File{Path: testScriptPath}

	start := time.Now()
	testExpr, err := exprWithParam(scriptBytes, testScriptPath, param)
	testFile.WallTime = time.Since(start)
	if err == nil {
		testFile.Results, err = test.RunExpr(arraictx.InitRunCtx(context.Background()), testExpr)
	}

	return testFile, err
}

// buildTransformInput prepares the input tuple that is accepted as the only parameter for transforms.
func BuildTransformInput(modules []*sysl.Module, modulePaths []string) (rel.Tuple, error) {
	models := make([]syslModel, 0, len(modules))

	for i, module := range modules {
		modPath := "stdin"
		if len(modulePaths) > i {
			modPath = modulePaths[i]
		}
		mod, err := buildModel(module, modPath)
		if err != nil {
			return nil, err
		}
		models = append(models, mod)
	}

	input, err := rel.NewTupleFromMap(map[string]interface{}{"models": models})
	if err != nil {
		return nil, err
	}

	return input, nil
}

// buildModel create a syslModel struct by normalizing the supplied Sysl module and packaging it together with the
// original document model and the path. A collection of them is consumed by transforms, and allows them to choose
// the preferred type of model they want to work with.
func buildModel(module *sysl.Module, path string) (syslModel, error) {
	docMod, err := arrai.SyslModuleToValue(module)
	if err != nil {
		return syslModel{}, err
	}

	relMod, err := relmod.Normalize(context.Background(), module)
	if err != nil {
		return syslModel{}, err
	}

	return syslModel{path: path, doc: docMod, rel: *relMod}, nil
}

type syslModel struct {
	path string
	doc  rel.Value
	rel  relmod.Schema
}
