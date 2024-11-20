package pbutil

import (
	"fmt"
	"io"
	"path"
	"regexp"

	"github.com/spf13/afero"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/anz-bank/sysl/pkg/sysl"
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

/*
	 Creates a directory path based on the Application name. Returns an open file descriptor.
	 E.g. if the Application is:
	 		'ANZx :: Foo :: Bar'
		  it creates the following file and returns the descriptor:
			<basePath>/ANZx/Foo/Bar/<fileName>
*/
func CreatePathForApplication(
	appName string,
	basePath string,
	appdata *sysl.Application,
	fileName string,
	fs afero.Fs) (afero.File, error) {
	dirName := path.Join(append([]string{basePath}, appdata.Name.GetPart()...)...)
	err := fs.MkdirAll(dirName, 0755)
	if err != nil {
		return nil, err
	}

	destFileName := path.Join(dirName, fileName)
	fd, err := fs.Create(destFileName)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

/*
Iterates through Applications and does the following for each one:
  - Creates an output path structure based on the Application name/parts
  - Identifies the correct output method to run based on `outputMode`
  - Executes the output method
*/
func OutputSplitApplications(
	module *sysl.Module,
	outputMode string,
	opt OutputOptions,
	basePath string,
	fileName string,
	fs afero.Fs) error {
	var err error = nil
	for appName, app := range module.Apps {
		fd, err := CreatePathForApplication(appName, basePath, app, fileName, fs)
		if err != nil {
			return err
		}

		outputFilePath := fd.Name()

		writer := func() error {
			switch outputMode {
			case "json":
				return FJSONPBWithOpt(fd, app, opt)
			case "textpb":
				fd.Close() // This method takes a filename, not a descriptor so close it.
				return TextPBWithOpt(app, outputFilePath, fs, opt)
			default:
				fd.Close() // This method takes a filename, not a descriptor so close it.
				return GeneratePBBinaryMessageFile(app, outputFilePath, fs)
			}
		}
		err = writer()
		fd.Close()
		if err != nil {
			break
		}
	}
	return err
}
