package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSFileSystem(t *testing.T) {
	fs := &osFileSystem{root: "."}

	f, err := fs.Open("sysl.go")
	assert.NoError(t, err)
	stat, err := f.Stat()
	assert.NoError(t, err)
	assert.False(t, stat.IsDir())
	assert.True(t, stat.Size() > 0)

	_, err = fs.Open("doesn't.exist")
	assert.Error(t, err)
}

func TestFSFileStream(t *testing.T) {
	fs := &osFileSystem{root: "."}

	s, err := newFSFileStream("sysl.go", fs)
	assert.NoError(t, err)
	assert.Equal(t, "package main", s.GetText(0, 11))
}
