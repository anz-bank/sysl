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
	}

	return result
}

func TestSimpleEP(t *testing.T) {
	filename := "tests/test1.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestAttribs(t *testing.T) {
	filename := "tests/attribs.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestIfElse(t *testing.T) {
	filename := "tests/if_else.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestArgs(t *testing.T) {
	filename := "tests/args.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestSimpleEPWithSpaces(t *testing.T) {
	filename := "tests/with_spaces.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestSimpleEP2(t *testing.T) {
	filename := "tests/test4.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestSimpleEndpointParams(t *testing.T) {
	filename := "tests/ep_params.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestOneOfStatements(t *testing.T) {
	filename := "tests/oneof.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestDuplicateEndpoints(t *testing.T) {
	filename := "tests/duplicate.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestEventing(t *testing.T) {
	filename := "tests/eventing.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestCollector(t *testing.T) {
	filename := "tests/collector.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestPubSubCollector(t *testing.T) {
	filename := "tests/pubsub_collector.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestDocstrings(t *testing.T) {
	filename := "tests/docstrings.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestMixins(t *testing.T) {
	filename := "tests/mixin.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}
func TestForLoops(t *testing.T) {
	filename := "tests/for_loop.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestGroupStmt(t *testing.T) {
	filename := "tests/group_stmt.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestUntilLoop(t *testing.T) {
	filename := "tests/until_loop.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestTuple(t *testing.T) {
	filename := "tests/test2.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestInplaceTuple(t *testing.T) {
	filename := "tests/inplace_tuple.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestRelational(t *testing.T) {
	filename := "tests/school.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestImports(t *testing.T) {
	filename := "tests/library.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestRootArg(t *testing.T) {
	filename := "school.sysl"
	root := "tests"
	if loadAndCompare(Parse(filename, root), filename, root) == false {
		t.Error("failed")
	}
}

func TestRestApi(t *testing.T) {
	filename := "tests/test_rest_api.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestRestApiQueryParams(t *testing.T) {
	filename := "tests/rest_api_query_params.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestSimpleProject(t *testing.T) {
	filename := "tests/project.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}

func TestUrlParamOrder(t *testing.T) {
	filename := "tests/rest_url_params.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == true {
		t.Error("failed")
	}
}

func TestRestApi_WrongOrder(t *testing.T) {
	filename := "tests/bad_order.sysl"
	if loadAndCompare(Parse(filename, ""), filename, "") == false {
		t.Error("failed")
	}
}
