package modv2

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/afero"
)

/* FsRetriever is an interface that returns file bytes and if the object should be cached in later steps */
type Retriever interface {
	Retrieve(resource string) (content []byte, cached bool, err error)
}

/* FsRetriever is a filesystem retriever */
type FsRetriever struct {
	fs afero.Fs
}

/* NewFs returns a FsRetriever from a afero filesystem */
func NewFs(fs afero.Fs) FsRetriever {
	return FsRetriever{
		fs: fs,
	}
}

/* Retrieve implements the retriever interface for FsRetriever */
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

/* ProcessRequest resolves a resourve into its components */
func ProcessRequest(resource string) (string, string, string, error) {
	locationVersion := strings.Split(resource, "@")
	if len(locationVersion) != 2 {
		return "", resource, "", nil
	}
	repoResource := locationVersion[0]
	version := locationVersion[1]
	parts := strings.Split(repoResource, "/")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("resource must be in form gitx.com/user/repo/resource.ext@hash")
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource, version, nil
}
