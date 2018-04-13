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
	args := []string{"pb", "-o", filename + ".pb", filename}
	if len(root) > 0 {
		root := []string{"--root", root}
		args = append(root, args...)
	}

	cmd := exec.Command("sysl", args...)
	err := cmd.Run()
	if err != nil {
		return false
	}
	buf := bytes.NewBuffer(nil)

	f, _ := os.Open(filename + ".pb")
	io.Copy(buf, f)
	f.Close()

	mod := sysl.Module{}
	err = proto.Unmarshal(buf.Bytes(), &mod)
	if err != nil {
		return false
	}
	// uncomment to compare
	fmt.Println("generated")
	TextPB(m2)
	fmt.Println("golden")
	TextPB(&mod)

	return proto.Equal(&mod, m2)
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

func TestRootArg(t *testing.T) {
	filename := "school.sysl"
	root := "tests"
	if loadAndCompare(Parse(filename, root), filename, root) == false {
		t.Error("failed")
	}
}
