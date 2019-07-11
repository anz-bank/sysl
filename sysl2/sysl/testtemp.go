package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testTempFile struct {
	f *os.File
}

func newTestTempFile(t *testing.T, dir, pattern string) *testTempFile {
	f, err := ioutil.TempFile("", "sysl_test.go-TestJSONPB-*.json")
	require.NoError(t, err, "newTestTempFile(%#v, %#v)", dir, pattern)
	return &testTempFile{f}
}

func testTempFilename(t *testing.T, dir, pattern string) string {
	tf := newTestTempFile(t, dir, pattern)
	defer tf.CloseAndRemove()
	return tf.Name()
}

func (tf *testTempFile) Name() string {
	return tf.f.Name()
}

func (tf *testTempFile) CloseAndRemove() {
	tf.f.Close()
	os.Remove(tf.Name())
}
