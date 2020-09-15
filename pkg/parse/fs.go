package parse

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/pkg/mod"
	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gop"
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

func newFSFileStream(filename string, retriever gop.Retriever) (s *fsFileStream, m *mod.Module, err error) {
	var res gop.Object
	repo, resource, version, err := app.ProcessRequest(filename)
	if err != nil {
		resource = filename
	}
	res, _, err = retriever.Retrieve(repo, resource, version)
	if err != nil {
		return nil, nil, err
	}
	return &fsFileStream{antlr.NewInputStream(string(res.Content)), filename}, m, nil
}
