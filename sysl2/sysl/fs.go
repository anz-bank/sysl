package main

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

func fileOrStdout(filename string, stdout io.Writer) (io.Writer, error) {
	if filename == "-" {
		return stdout, nil
	}
	w, err := os.Create(filename)
	return w, errors.Wrapf(err, "Error creating output file %q", filename)
}

func fileExists(filename string, fs http.FileSystem) bool {
	f, err := fs.Open(filename)
	if err != nil {
		return false
	}
	_, err = f.Stat()
	return err == nil
}

func dirExists(dirName string) bool {
	info, err := os.Stat(dirName)
	return err == nil && info.IsDir()
}

type fsFileStream struct {
	*antlr.InputStream
	filename string
}

func newFSFileStream(filename string, fs http.FileSystem) (*fsFileStream, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, f); err != nil {
		return nil, err
	}

	s := string(buf.Bytes())
	return &fsFileStream{antlr.NewInputStream(s), filename}, nil
}
