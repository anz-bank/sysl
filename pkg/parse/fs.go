package parse

import (
	"bytes"
	"io"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/sysl/pkg/mod"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func fileExists(filename string, fs afero.Fs) bool {
	f, err := fs.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()
	logrus.Debugf("opened file %s", f.Name())

	_, err = f.Stat()
	return err == nil
}

type fsFileStream struct {
	*antlr.InputStream
	filename string
}

func newFSFileStream(filename string, fs afero.Fs) (*fsFileStream, *mod.Module, error) {
	f, mod, err := fs.OpenWithModule(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, f); err != nil {
		return nil, err
	}

	return &fsFileStream{antlr.NewInputStream(buf.String()), filename}, mod, nil
}
