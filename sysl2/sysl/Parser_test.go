package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/golang/protobuf/proto"
)

func pyParse(filename, root string) (*sysl.Module, error) {
	output := filename + ".pb"

	args := []string{"pb", "-o", output, filename}
	if len(root) > 0 {
		rootArg := []string{"--root", root}
		// TODO: This looks dubious
		args[2] = root + "/" + args[2]
		output = args[2]
		args = append(rootArg, args...)
	}

	pySysl, found := os.LookupEnv("SYSL_PYTHON_BIN")
	if !found {
		pySysl = "sysl"
	}
	cmd := exec.Command(pySysl, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)

	f, _ := os.Open(output)
	io.Copy(buf, f)
	f.Close()

	module := &sysl.Module{}
	if err := proto.Unmarshal(buf.Bytes(), module); err != nil {
		return nil, err
	}
	return module, nil
}

func parseComparable(filename, root string) (*sysl.Module, error) {
	module, err := Parse(filename, root)
	if err != nil {
		return nil, err
	}

	// remove stuff that does not match legacy.
	for _, app := range module.Apps {
		app.SourceContext = nil
		// app.SourceContext.Start.Col = 0
		for _, ep := range app.Endpoints {
			ep.SourceContext = nil
		}
		for _, t := range app.Types {
			t.SourceContext = nil
		}
	}

	return module, nil
}

func parseAndCompare(filename, root string) (bool, error) {
	goModule, err := parseComparable(filename, root)
	if err != nil {
		return false, err
	}

	pyModule, err := pyParse(filename, root)
	if err != nil {
		return false, err
	}
	if proto.Equal(pyModule, goModule) {
		return true, nil
	}

	golden := "golden.txt"
	generated := "generated.txt"
	if err = TextPB(pyModule, golden); err != nil {
		return false, err
	}
	if err = TextPB(goModule, generated); err != nil {
		return false, err
	}
	cmd := exec.Command("diff", "-y", golden, generated)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return false, nil
}

func parseAndPrint(t *testing.T, filename, root string) error {
	goModule, err := parseComparable(filename, root)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Outputting %#v:\n", filename)
	FTextPB(os.Stderr, goModule)
	fmt.Fprintf(os.Stderr, "------------------------\n")
	return nil
}

func testParse(t *testing.T, filename, root string) {
	equal, err := parseAndCompare(filename, root)
	if err != nil {
		t.Errorf("Parsing error: %v", err)
	}
	if !equal {
		t.Error("Mismatch")
	}
}

func TestSimpleEP(t *testing.T) {
	testParse(t, "tests/test1.sysl", "")
}

func TestAttribs(t *testing.T) {
	testParse(t, "tests/attribs.sysl", "")
}

func TestIfElse(t *testing.T) {
	testParse(t, "tests/if_else.sysl", "")
}

func TestArgs(t *testing.T) {
	testParse(t, "tests/args.sysl", "")
}

func TestSimpleEPWithSpaces(t *testing.T) {
	testParse(t, "tests/with_spaces.sysl", "")
}

func TestSimpleEP2(t *testing.T) {
	testParse(t, "tests/test4.sysl", "")
}

func TestSimpleEndpointParams(t *testing.T) {
	testParse(t, "tests/ep_params.sysl", "")
}

func TestOneOfStatements(t *testing.T) {
	testParse(t, "tests/oneof.sysl", "")
}

func TestDuplicateEndpoints(t *testing.T) {
	testParse(t, "tests/duplicate.sysl", "")
}

func TestEventing(t *testing.T) {
	testParse(t, "tests/eventing.sysl", "")
}

func TestCollector(t *testing.T) {
	testParse(t, "tests/collector.sysl", "")
}

func TestPubSubCollector(t *testing.T) {
	testParse(t, "tests/pubsub_collector.sysl", "")
}

func TestDocstrings(t *testing.T) {
	testParse(t, "tests/docstrings.sysl", "")
}

func TestMixins(t *testing.T) {
	testParse(t, "tests/mixin.sysl", "")
}
func TestForLoops(t *testing.T) {
	testParse(t, "tests/for_loop.sysl", "")
}

func TestGroupStmt(t *testing.T) {
	testParse(t, "tests/group_stmt.sysl", "")
}

func TestUntilLoop(t *testing.T) {
	testParse(t, "tests/until_loop.sysl", "")
}

func TestTuple(t *testing.T) {
	testParse(t, "tests/test2.sysl", "")
}

func TestInplaceTuple(t *testing.T) {
	testParse(t, "tests/inplace_tuple.sysl", "")
}

func TestRelational(t *testing.T) {
	testParse(t, "tests/school.sysl", "")
}

func TestImports(t *testing.T) {
	testParse(t, "tests/library.sysl", "")
}

func TestRootArg(t *testing.T) {
	testParse(t, "school.sysl", "tests")
}

func TestRestApi(t *testing.T) {
	testParse(t, "tests/test_rest_api.sysl", "")
}

func TestRestApiQueryParams(t *testing.T) {
	testParse(t, "tests/rest_api_query_params.sysl", "")
}

func TestSimpleProject(t *testing.T) {
	testParse(t, "tests/project.sysl", "")
}

func TestUrlParamOrder(t *testing.T) {
	// Output won't match legacy; visually inspect the diff.
	parseAndCompare("tests/rest_url_params.sysl", "")
}

func TestRestApi_WrongOrder(t *testing.T) {
	testParse(t, "tests/bad_order.sysl", "")
}

func TestTransform(t *testing.T) {
	testParse(t, "tests/transform.sysl", "")
}

func TestImpliedDot(t *testing.T) {
	testParse(t, "tests/implied.sysl", "")
}

func TestStmts(t *testing.T) {
	testParse(t, "tests/stmts.sysl", "")
}

func TestMath(t *testing.T) {
	testParse(t, "tests/math.sysl", "")
}

func TestTableof(t *testing.T) {
	testParse(t, "tests/tableof.sysl", "")
}

func TestRank(t *testing.T) {
	testParse(t, "tests/rank.sysl", "")
}

func TestMatching(t *testing.T) {
	testParse(t, "tests/matching.sysl", "")
}

func TestNavigate(t *testing.T) {
	testParse(t, "tests/navigate.sysl", "")
}

func TestFuncs(t *testing.T) {
	testParse(t, "tests/funcs.sysl", "")
}

func TestPetshop(t *testing.T) {
	testParse(t, "../../demo/petshop/petshop.sysl", "")
}

func TestCrash(t *testing.T) {
	testParse(t, "tests/crash.sysl", "")
}
