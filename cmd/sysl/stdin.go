package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strings"
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
func loadStdinFiles(src []byte) ([]stdinFile, error) {
	var files []stdinFile
	if len(src) > 0 {
		if err := json.Unmarshal(src, &files); err != nil || len(files) == 0 {
			errMsg := fmt.Sprintf("stdin must be a JSON array of {path, content} file objects. "+
				"JSON parsing error: %v", err)

			// From: https://adrianhesketh.com/2017/03/18/getting-line-and-character-positions-from-gos-json-unmarshal-errors/
			offset := int64(-1)
			if jsonError, ok := err.(*json.SyntaxError); ok {
				offset = jsonError.Offset - 1
			} else if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
				offset = jsonError.Offset - 1
			}

			// From: https://github.com/golang/go/issues/43513#issuecomment-755754498
			if offset >= 0 && offset < math.MaxInt32 {
				offset := int(offset)
				jsonUntilErr := string(src[:offset])
				line := 1 + strings.Count(jsonUntilErr, "\n")
				col := 1 + offset - (strings.LastIndex(jsonUntilErr, "\n") + len("\n"))
				errMsg += fmt.Sprintf(" (line %d, col %d)", line, col)
			}

			return nil, errors.New(errMsg)
		}
	}

	for _, file := range files {
		file.Path = filepath.Clean(file.Path)
	}

	return files, nil
}
