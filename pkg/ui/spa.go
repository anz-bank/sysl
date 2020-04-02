package ui

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
)

// Adapted from the Readme at https://github.com/gorilla/mux

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
	fileSystem fileSystem
}

// fileSystem abstracts away the filesystem implementation for better testing
type fileSystem interface {
	Stat(name string) (os.FileInfo, error)
	Dir(name string) http.FileSystem
}

// pkgerFS implements the fileSystem interface and wraps calls to pkger
type pkgerFS struct{}

func (f pkgerFS) Stat(name string) (os.FileInfo, error) {
	return pkger.Stat(name)
}

func (f pkgerFS) Dir(name string) http.FileSystem {
	return pkger.Dir(name)
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = h.fileSystem.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.FileServer(h.fileSystem.Dir(h.indexPath)).ServeHTTP(w, r)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 lspframework server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(h.fileSystem.Dir(h.staticPath)).ServeHTTP(w, r)
}
