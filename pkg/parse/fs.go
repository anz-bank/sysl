package parse

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
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
