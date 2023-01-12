package pbutil

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

var ErrUnknownExtension = errors.New("unknown extension")

// FromPB unmarshals a Sysl module from an encoded protobuf message file.
func FromPB(pbPath string, fs afero.Fs) (*sysl.Module, error) {
	f, err := fs.Open(pbPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	in, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	m, err := FromPBByteContents(pbPath, in)

	// If unknown just try pb
	if errors.Is(err, ErrUnknownExtension) {
		m = &sysl.Module{}
		err = proto.Unmarshal(in, m)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse sysl module: %w", err)
	}

	return m, nil
}

// FromPBStringContents unmarshals a Sysl module from an encoded protobuf message.
func FromPBStringContents(pbPath, contents string) (*sysl.Module, error) {
	return fromPBContents(pbPath, contents, func(cont interface{}) []byte {
		return []byte(cont.(string))
	})
}

// FromPBByteContents unmarshals a Sysl module from an encoded protobuf message.
func FromPBByteContents(pbPath string, contents []byte) (*sysl.Module, error) {
	return fromPBContents(pbPath, contents, func(cont interface{}) []byte {
		return cont.([]byte)
	})
}

func fromPBContents(pbPath string, contents interface{}, toBytes func(interface{}) []byte) (*sysl.Module, error) {
	m := &sysl.Module{}

	switch {
	case strings.HasSuffix(pbPath, ".pb"):
		err := proto.Unmarshal(toBytes(contents), m)

		return m, err

	case strings.HasSuffix(pbPath, ".pb.json"):
		err := protojson.Unmarshal(toBytes(contents), m)

		return m, err

	case strings.HasSuffix(pbPath, ".textpb"):
		err := prototext.Unmarshal(toBytes(contents), m)

		return m, err
	}

	return nil, ErrUnknownExtension
}
