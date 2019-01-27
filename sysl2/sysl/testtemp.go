package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTempFile struct {
	f *os.File
}

func newTestTempFile(t *testing.T, dir, pattern string) *testTempFile {
	f, err := ioutil.TempFile("", "github.com-sysl-sysl2-sysl-sysl_test.go-TestJSONPB-*.json")
	if assert.NoError(t, err, "newTestTempFile(%#v, %#v)", dir, pattern) {
		return &testTempFile{f}
	}
	return nil
}

func testTempFilename(t *testing.T, dir, pattern string) string {
	if tf := newTestTempFile(t, dir, pattern); tf != nil {
		defer tf.CloseAndRemove()
		return tf.Name()
	}
	return ""
}

func (tf *testTempFile) File() *os.File {
	return tf.f
}

func (tf *testTempFile) Name() string {
	return tf.f.Name()
}

func (tf *testTempFile) CloseAndRemove() {
	tf.f.Close()
	os.Remove(tf.Name())
}
