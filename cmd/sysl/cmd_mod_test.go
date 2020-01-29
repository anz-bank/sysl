package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSyslModInit(t *testing.T) {
	fs := afero.NewOsFs()
	execArgs := ExecuteArgs{
		Modules:        nil,
		Filesystem:     fs,
		Logger:         logrus.StandardLogger(),
		DefaultAppName: "",
	}

	err := fs.RemoveAll(syslRootMarker)
	assert.NoError(t, err)

	err = syslModInit(execArgs)
	assert.NoError(t, err)

	// TODO: move this to Cleanup() in golang 1.14 (this test may pollute the cwd)
	err = fs.RemoveAll(syslRootMarker)
	assert.NoError(t, err)
}

func TestSyslModInitRootAlreadyExists(t *testing.T) {
	fs := afero.NewOsFs()
	execArgs := ExecuteArgs{
		Modules:        nil,
		Filesystem:     fs,
		Logger:         logrus.StandardLogger(),
		DefaultAppName: "",
	}

	err := fs.RemoveAll(syslRootMarker)
	assert.NoError(t, err)

	err = syslModInit(execArgs)
	assert.NoError(t, err)

	err = syslModInit(execArgs)
	assert.Error(t, err)

	// TODO: move this to Cleanup() in golang 1.14 (this test may pollute the cwd)
	err = fs.RemoveAll(syslRootMarker)
	assert.NoError(t, err)
}
