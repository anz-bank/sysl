package main

import (
	"fmt"
	"io"
	"os"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// JSONPB ...
func JSONPB(m *sysl.Module, filename string) error {
	if m == nil {
		return fmt.Errorf("module is nil: %#v", filename)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return FJSONPB(f, m)
}

// FJSONPB ...
func FJSONPB(w io.Writer, m *sysl.Module) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	ma := jsonpb.Marshaler{Indent: " ", EmitDefaults: false}
	return ma.Marshal(w, m)
}

// TextPB ...
func TextPB(m *sysl.Module, filename string) error {
	if m == nil {
		return fmt.Errorf("module is nil: %#v", filename)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return FTextPB(f, m)
}

// FTextPB ...
func FTextPB(w io.Writer, m *sysl.Module) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	return proto.MarshalText(w, m)
}
