package pbutil

import (
	"fmt"
	"io/ioutil"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
)

// FromPB unmarshals a Sysl module from an encoded protobuf message.
func FromPB(pbPath string, fs afero.Fs) (*sysl.Module, error) {
	f, err := fs.Open(pbPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}
	in, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	m := &sysl.Module{}
	if err := proto.Unmarshal(in, m); err != nil {
		return nil, fmt.Errorf("failed to parse sysl module: %s", err)
	}
	return m, nil
}
