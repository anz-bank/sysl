package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"testing"

	"anz-bank/sysl/src/proto"
	"github.com/golang/protobuf/proto"
)

func loadAndCompare(filename string, m2 *sysl.Module) bool {
	cmd := exec.Command("sysl", "pb", "-o", filename+".pb", filename)
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
	// TextPB(m2)
	// TextPB(&mod)

	return proto.Equal(&mod, m2)
}

func TestSimpleEP(t *testing.T) {
	filename := "tests/test1.sysl"
	if loadAndCompare(filename, Parse(filename, "")) == false {
		t.Error("failed")
	}
}
func TestSimpleEP2(t *testing.T) {
	filename := "tests/test4.sysl"
	if loadAndCompare(filename, Parse(filename, "")) == false {
		t.Error("failed")
	}
}
func TestTuple(t *testing.T) {
	filename := "tests/test2.sysl"
	if loadAndCompare(filename, Parse(filename, "")) == false {
		t.Error("failed")
	}
}
func TestRelational(t *testing.T) {
	filename := "tests/test3.sysl"
	if loadAndCompare(filename, Parse(filename, "")) == false {
		t.Error("failed")
	}
}
