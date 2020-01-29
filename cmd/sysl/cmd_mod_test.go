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

	fs.Remove("go.mod")
	fs.Remove("go.sum")

	err := syslModInit(execArgs)
	assert.NoError(t, err)

	fs.Remove("go.mod")
	fs.Remove("go.sum")
}

func TestSyslModInitAlreadyExists(t *testing.T) {
	fs := afero.NewOsFs()
	execArgs := ExecuteArgs{
		Modules:        nil,
		Filesystem:     fs,
		Logger:         logrus.StandardLogger(),
		DefaultAppName: "",
	}

	fs.Remove("go.mod")
	fs.Remove("go.sum")

	err := syslModInit(execArgs)
	assert.NoError(t, err)

	err = syslModInit(execArgs)
	assert.Error(t, err)

	fs.Remove("go.mod")
	fs.Remove("go.sum")
}
