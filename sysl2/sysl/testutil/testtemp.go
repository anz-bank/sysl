package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TempFile represents a temporary file.
type TempFile struct {
	f *os.File
}

// NewTempFile creates a TempFile.
func NewTempFile(t *testing.T, dir, pattern string) *TempFile {
	f, err := ioutil.TempFile("", "sysl_test.go-TestJSONPB-*.json")
	require.NoError(t, err, "NewTempFile(%#v, %#v)", dir, pattern)
	return &TempFile{f}
}

// TempFilename choose a name for a temporary file to be created by the caller.
func TempFilename(t *testing.T, dir, pattern string) string {
	tf := NewTempFile(t, dir, pattern)
	defer tf.CloseAndRemove()
	return tf.Name()
}

// Name returns the name of TempFile.
func (tf *TempFile) Name() string {
	return tf.f.Name()
}

// CloseAndRemove closes and removes a TempFile.
func (tf *TempFile) CloseAndRemove() {
	tf.f.Close()
	os.Remove(tf.Name())
}
