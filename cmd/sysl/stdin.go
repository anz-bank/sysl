package main

import (
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
)

// stdinFile represents a file to parse provided via stdin (JSON encoded). The content is the source
// to parse, and the path indicates the location from which imports should be resolved.
//
// For example, without stdin, if there is some file foo.sysl on disk, `sysl cmd path/to/foo.sysl`
// would parse that file and execute command cmd on the resulting module. To do the same thing via
// stdin:
//
// `echo '[{"path": "path/to/foo.sysl", "content": "$(cat path/to/foo.sysl)"}]' | sysl cmd`
//
// The advantage is that the content of foo.sysl need not match what is on disk. For example, it
// could contain some buffered content that should be validated before being saved.
type stdinFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// loadFromStdin reads in a list of stdinFiles from the Reader.
func loadStdinFiles(stdin io.Reader) ([]stdinFile, error) {
	var files []stdinFile

	src, err := io.ReadAll(stdin)
	if err != nil {
		return nil, err
	}

	if len(src) > 0 {
		if err := json.Unmarshal(src, &files); err != nil || len(files) == 0 {
			return nil, errors.New("stdin must be a JSON array of {path, content} file objects")
		}
	}

	for _, file := range files {
		file.Path = filepath.Clean(file.Path)
	}

	return files, nil
}
