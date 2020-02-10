package mod

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func removeFile(t *testing.T, fs afero.Fs, file string) {
	exists, err := afero.Exists(fs, file)
	assert.NoError(t, err)
	if exists {
		err = fs.Remove(file)
		assert.NoError(t, err)
	}
}

func TestSyslModInit(t *testing.T) {
	fs := afero.NewOsFs()
	logger := logrus.StandardLogger()

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := SyslModInit(logger)
	assert.NoError(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}

func TestSyslModInitAlreadyExists(t *testing.T) {
	fs := afero.NewOsFs()
	logger := logrus.StandardLogger()

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := SyslModInit(logger)
	assert.NoError(t, err)

	err = SyslModInit(logger)
	assert.Error(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}
