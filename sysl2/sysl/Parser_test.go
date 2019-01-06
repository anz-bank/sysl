package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func readSyslModule(filename string) (*sysl.Module, error) {
	var buf bytes.Buffer

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	io.Copy(&buf, f)
	f.Close()

	module := &sysl.Module{}
	if err := proto.UnmarshalText(buf.String(), module); err != nil {
		return nil, err
	}
	return module, nil
}

func pySysl() string {
	if pySysl, found := os.LookupEnv("SYSL_PYTHON_BIN"); found {
		return pySysl
	}
	return "sysl"
}

func pyParse(filename, root, output string) (*sysl.Module, error) {
	args := []string{"textpb", "-o", output, filename}
	if len(root) > 0 {
		rootArg := []string{"--root", root}
		// TODO: This looks dubious
		args[2] = root + "/" + args[2]
		output = args[2]
		args = append(rootArg, args...)
	}

	cmd := exec.Command(pySysl(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return readSyslModule(output)
}

func parseComparable(filename, root string) (*sysl.Module, error) {
	fs := &osFileSystem{root}
	module, err := FSParse(filename, fs)
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

func parseAndCompare(filename, root, golden string, goldenModule *sysl.Module) (bool, error) {
	module, err := parseComparable(filename, root)
	if err != nil {
		return false, err
	}

	if proto.Equal(goldenModule, module) {
		return true, nil
	}

	generated, err := ioutil.TempFile("", "sysl-test-*.textpb")
	if err != nil {
		return false, err
	}
	defer generated.Close()
	defer os.Remove(generated.Name())

	if err = TextPB(goldenModule, golden); err != nil {
		return false, err
	}

	if err = FTextPB(generated, module); err != nil {
		return false, err
	}
	generated.Close()

	cmd := exec.Command("diff", "-y", golden, generated.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return false, err
	}

	return false, nil
}

func parseAndCompareWithPython(filename, root string) (bool, error) {
	fmt.Printf("%35s <=> $(%s %[1]s)\n", filename, pySysl())

	golden, err := ioutil.TempFile("", "sysl-test-golden-*.textpb")
	if err != nil {
		return false, err
	}
	defer golden.Close()
	defer os.Remove(golden.Name())

	pyModule, err := pyParse(filename, root, golden.Name())
	if err != nil {
		return false, err
	}

	return parseAndCompare(filename, root, golden.Name(), pyModule)
}

func parseAndCompareWithGolden(filename, root string) (bool, error) {
	golden := filename + ".golden.textpb"
	fmt.Printf("%35s <=> %s\n", filename, golden)

	goldenModule, err := readSyslModule(golden)
	if err != nil {
		return false, err
	}
	return parseAndCompare(filename, root, golden, goldenModule)
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

func testParseAgainstPython(t *testing.T, filename, root string) {
	equal, err := parseAndCompareWithPython(filename, root)
	assert.NoError(t, err)
	assert.True(t, equal, "Mismatch")
}

func testParseAgainstGolden(t *testing.T, filename, root string) {
	equal, err := parseAndCompareWithGolden(filename, root)
	assert.NoError(t, err)
	assert.True(t, equal, "Mismatch")
}

func TestParseMissingFile(t *testing.T) {
	_, err := parseAndCompareWithGolden("tests/doesn't.exist", "")
	assert.Error(t, err)
}

func TestParseBadFile(t *testing.T) {
	_, err := parseAndCompareWithGolden("sysl.go", "")
	assert.Error(t, err)
}

func TestSimpleEP(t *testing.T) {
	testParseAgainstPython(t, "tests/test1.sysl", "")
}

func TestAttribs(t *testing.T) {
	testParseAgainstPython(t, "tests/attribs.sysl", "")
}

func TestIfElse(t *testing.T) {
	testParseAgainstPython(t, "tests/if_else.sysl", "")
}

func TestArgs(t *testing.T) {
	testParseAgainstPython(t, "tests/args.sysl", "")
}

func TestSimpleEPWithSpaces(t *testing.T) {
	testParseAgainstPython(t, "tests/with_spaces.sysl", "")
}

func TestSimpleEP2(t *testing.T) {
	testParseAgainstPython(t, "tests/test4.sysl", "")
}

func TestSimpleEndpointParams(t *testing.T) {
	testParseAgainstPython(t, "tests/ep_params.sysl", "")
}

func TestOneOfStatements(t *testing.T) {
	testParseAgainstPython(t, "tests/oneof.sysl", "")
}

func TestDuplicateEndpoints(t *testing.T) {
	testParseAgainstPython(t, "tests/duplicate.sysl", "")
}

func TestEventing(t *testing.T) {
	testParseAgainstPython(t, "tests/eventing.sysl", "")
}

func TestCollector(t *testing.T) {
	testParseAgainstPython(t, "tests/collector.sysl", "")
}

func TestPubSubCollector(t *testing.T) {
	testParseAgainstPython(t, "tests/pubsub_collector.sysl", "")
}

func TestDocstrings(t *testing.T) {
	testParseAgainstPython(t, "tests/docstrings.sysl", "")
}

func TestMixins(t *testing.T) {
	testParseAgainstPython(t, "tests/mixin.sysl", "")
}
func TestForLoops(t *testing.T) {
	testParseAgainstPython(t, "tests/for_loop.sysl", "")
}

func TestGroupStmt(t *testing.T) {
	testParseAgainstPython(t, "tests/group_stmt.sysl", "")
}

func TestUntilLoop(t *testing.T) {
	testParseAgainstPython(t, "tests/until_loop.sysl", "")
}

func TestTuple(t *testing.T) {
	testParseAgainstPython(t, "tests/test2.sysl", "")
}

func TestInplaceTuple(t *testing.T) {
	testParseAgainstPython(t, "tests/inplace_tuple.sysl", "")
}

func TestRelational(t *testing.T) {
	testParseAgainstPython(t, "tests/school.sysl", "")
}

func TestImports(t *testing.T) {
	testParseAgainstPython(t, "tests/library.sysl", "")
}

func TestRootArg(t *testing.T) {
	testParseAgainstPython(t, "school.sysl", "tests")
}

func TestRestApi(t *testing.T) {
	testParseAgainstPython(t, "tests/test_rest_api.sysl", "")
}

func TestRestApiQueryParams(t *testing.T) {
	testParseAgainstPython(t, "tests/rest_api_query_params.sysl", "")
}

func TestSimpleProject(t *testing.T) {
	testParseAgainstPython(t, "tests/project.sysl", "")
}

func TestUrlParamOrder(t *testing.T) {
	filename := "tests/rest_url_params.sysl"
	parseAndCompareWithPython(filename, "")
	fmt.Printf("Output for %#v won't match legacy. Visually inspect the above diff.\n", filename)
}

func TestUrlParamOrderAgainstGolden(t *testing.T) {
	testParseAgainstGolden(t, "tests/rest_url_params.sysl", "")
}

func TestRestApi_WrongOrder(t *testing.T) {
	testParseAgainstPython(t, "tests/bad_order.sysl", "")
}

func TestTransform(t *testing.T) {
	testParseAgainstPython(t, "tests/transform.sysl", "")
}

func TestImpliedDot(t *testing.T) {
	testParseAgainstPython(t, "tests/implied.sysl", "")
}

func TestStmts(t *testing.T) {
	testParseAgainstPython(t, "tests/stmts.sysl", "")
}

func TestMath(t *testing.T) {
	testParseAgainstPython(t, "tests/math.sysl", "")
}

func TestTableof(t *testing.T) {
	testParseAgainstPython(t, "tests/tableof.sysl", "")
}

func TestRank(t *testing.T) {
	testParseAgainstPython(t, "tests/rank.sysl", "")
}

func TestMatching(t *testing.T) {
	testParseAgainstPython(t, "tests/matching.sysl", "")
}

func TestNavigate(t *testing.T) {
	testParseAgainstPython(t, "tests/navigate.sysl", "")
}

func TestFuncs(t *testing.T) {
	testParseAgainstPython(t, "tests/funcs.sysl", "")
}

func TestPetshop(t *testing.T) {
	testParseAgainstPython(t, "petshop.sysl", "../../demo/petshop")
}

func TestCrash(t *testing.T) {
	testParseAgainstPython(t, "tests/crash.sysl", "")
}
