package pbutil

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/afero"
)

// JSONPB ...
func JSONPB(m proto.Message, filename string, fs afero.Fs) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	f, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return FJSONPB(f, m)
}

// FJSONPB ...
func FJSONPB(w io.Writer, m proto.Message) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	ma := jsonpb.Marshaler{Indent: " ", EmitDefaults: false}
	return ma.Marshal(w, m)
}

// TextPB ...
func TextPB(m proto.Message, filename string, fs afero.Fs) error {
	if m == nil {
		return fmt.Errorf("module is nil: %#v", filename)
	}

	f, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return FTextPB(f, m)
}

// FTextPB ...
func FTextPB(w io.Writer, m proto.Message) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	return proto.MarshalText(w, m)
}
