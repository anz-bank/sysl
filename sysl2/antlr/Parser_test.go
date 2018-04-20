package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"anz-bank/sysl/src/proto"
	"github.com/golang/protobuf/proto"
)

func loadAndCompare(m2 *sysl.Module, filename string, root string) bool {
	output := filename + ".pb"

	args := []string{"pb", "-o", output, filename}
	if len(root) > 0 {
		root := []string{"--root", root}
		args[2] = root[1] + args[2]
		output = args[2]
		args = append(root, args...)
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
	// if !result {
	ioutil.WriteFile("generated.txt", []byte(proto.MarshalTextString(m2)), os.ModePerm)
	ioutil.WriteFile("golden.txt", []byte(proto.MarshalTextString(&mod)), os.ModePerm)
	// }

	return result
}

func TestSimpleEP(t *testing.T) {
	filename := "tests/test1.sysl"
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
func TestTuple(t *testing.T) {
	filename := "tests/test2.sysl"
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

func TestRestApi(t *testing.T) {
	filename := "tests/test_rest_api.sysl"
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
