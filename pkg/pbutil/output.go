package pbutil

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// JSONPB ...
func JSONPB(m protoreflect.ProtoMessage, filename string, fs afero.Fs) error {
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
func FJSONPB(w io.Writer, m protoreflect.ProtoMessage) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	ma := protojson.MarshalOptions{Multiline: true, Indent: " ", EmitUnpopulated: false}
	mb, err := ma.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(mb)
	return err
}

// TextPB ...
func TextPB(m protoreflect.ProtoMessage, filename string, fs afero.Fs) error {
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
func FTextPB(w io.Writer, m protoreflect.ProtoMessage) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	pt := prototext.MarshalOptions{Multiline: true, Indent: "  ", EmitUnknown: false}
	mt, err := pt.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(mt)
	return err
}

// GeneratePBBinaryMessageFile generates binary message to the file specified by `filename`.
func GeneratePBBinaryMessageFile(m protoreflect.ProtoMessage, filename string, fs afero.Fs) error {
	if m == nil {
		return fmt.Errorf("module is nil: %#v", filename)
	}

	f, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return GeneratePBBinaryMessage(f, m)
}

// GeneratePBBinaryMessage generates binary message to IO writer specified by `w`.
func GeneratePBBinaryMessage(w io.Writer, m protoreflect.ProtoMessage) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	bytes, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}
