package pbutil

import (
	"fmt"
	"io"
	"regexp"

	"github.com/spf13/afero"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type OutputOptions struct {
	Compact bool
}

// JSONPB ...
func JSONPB(m protoreflect.ProtoMessage, filename string, fs afero.Fs) error {
	return JSONPBWithOpt(m, filename, fs, OutputOptions{})
}

func JSONPBWithOpt(m protoreflect.ProtoMessage, filename string, fs afero.Fs, o OutputOptions) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	f, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return FJSONPBWithOpt(f, m, o)
}

// Recognise extra whitespace after a JSON key.
var extraSpaceAfterKeyRE = regexp.MustCompile(`(?m)^(\s*"[^"]*": ) `)

// FJSONPB ...
func FJSONPB(w io.Writer, m protoreflect.ProtoMessage) error {
	return FJSONPBWithOpt(w, m, OutputOptions{})
}

func FJSONPBWithOpt(w io.Writer, m protoreflect.ProtoMessage, o OutputOptions) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	ma := protojson.MarshalOptions{Multiline: true, Indent: " ", EmitUnpopulated: false}
	if o.Compact {
		ma.Multiline = false
		ma.Indent = ""
	}
	mb, err := ma.Marshal(m)
	if err != nil {
		return err
	}

	// TODO: Remove after we get protobuf working locally and in docker builds.
	mb = extraSpaceAfterKeyRE.ReplaceAll(mb, []byte("$1"))

	_, err = w.Write(mb)
	return err
}

// TextPB ...
func TextPB(m protoreflect.ProtoMessage, filename string, fs afero.Fs) error {
	return TextPBWithOpt(m, filename, fs, OutputOptions{})
}

func TextPBWithOpt(m protoreflect.ProtoMessage, filename string, fs afero.Fs, o OutputOptions) error {
	if m == nil {
		return fmt.Errorf("module is nil: %#v", filename)
	}

	f, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return FTextPBWithOpt(f, m, o)
}

// FTextPB ...
func FTextPB(w io.Writer, m protoreflect.ProtoMessage) error {
	return FTextPBWithOpt(w, m, OutputOptions{})
}

func FTextPBWithOpt(w io.Writer, m protoreflect.ProtoMessage, o OutputOptions) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	pt := prototext.MarshalOptions{Multiline: true, Indent: " ", EmitUnknown: false}
	if o.Compact {
		pt.Multiline = false
		pt.Indent = ""
	}
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
