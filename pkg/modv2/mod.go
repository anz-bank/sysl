package modv2

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/afero"
)

type FsRetriever struct {
	fs afero.Fs
}

func NewFs(fs afero.Fs) FsRetriever {
	return FsRetriever{
		fs: fs,
	}
}

/* FsRetriever is an interface that returns a Object and if the object should be cached in later steps */
type Retriever interface {
	Retrieve(resource string) (content []byte, cached bool, err error)
}

func (r FsRetriever) Retrieve(resource string) ([]byte, bool, error) {
	file, err := r.fs.Open(resource)
	if file == nil {
		return nil, false, fmt.Errorf("error opening file %w", err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false, fmt.Errorf("error opening file %w", err)
	}
	return b, true, nil
}
