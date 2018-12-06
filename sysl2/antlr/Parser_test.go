package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"anz-bank/sysl/src/proto"
	"github.com/golang/protobuf/proto"
)

func loadAndCompare(m2 *sysl.Module, filename string, root string) bool {

	// remove that does not match legacy.
	for _, app := range m2.Apps {
		app.SourceContext = nil
		// app.SourceContext.Start.Col = 0
		for _, ep := range app.Endpoints {
			ep.SourceContext = nil
		}
		for _, t := range app.Types {
			t.SourceContext = nil
		}
	}

	output := filename + ".pb"

	args := []string{"pb", "-o", output, filename}
	if len(root) > 0 {
		root_array := []string{"--root", root}
		args[2] = root_array[1] + "/" + args[2]
		output = args[2]
		args = append(root_array, args...)
	}

	cmd := exec.Command("sysl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return false
	}
	buf := bytes.NewBuffer(nil)

	f, _ := os.Open(output)
	io.Copy(buf, f)
	f.Close()

	mod := sysl.Module{}
	err = proto.Unmarshal(buf.Bytes(), &mod)
	if err != nil {
		fmt.Println(err)
		return false
	}
	result := proto.Equal(&mod, m2)
	// uncomment to compare
	if !result {
		TextPB(m2, "generated.txt")
		TextPB(&mod, "golden.txt")
		cmd = exec.Command("diff", "-y", "golden.txt", "generated.txt")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	return result
}

func testParse(t *testing.T, filename, root string) {
	if !loadAndCompare(Parse(filename, root), filename, root) {
		t.Error("failed")
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
	filename := "tests/rest_url_params.sysl"
	// output does not match with legacy
	// the order does not match
	// check the diff.
	if loadAndCompare(Parse(filename, ""), filename, "") == true {
		t.Error("failed")
	}
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
