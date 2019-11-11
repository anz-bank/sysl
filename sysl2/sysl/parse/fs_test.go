package parse

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	t.Parallel()

	// I think, therefore I am.
	assert.True(t, fileExists("fs_test.go", syslutil.NewChrootFs(afero.NewOsFs(), ".")))
}

func TestFileExistsBadFile(t *testing.T) {
	t.Parallel()

	assert.False(t, fileExists("x", syslutil.NewChrootFs(afero.NewOsFs(), "/non-existent.dir")))
	assert.False(t, fileExists("non-existent.file", syslutil.NewChrootFs(afero.NewOsFs(), ".")))
	assert.False(t, fileExists("non-existent.file", syslutil.NewChrootFs(afero.NewOsFs(), ".")))
}

func TestNewFSFileStream(t *testing.T) {
	t.Parallel()

	fs, err := newFSFileStream("fs_test.go", syslutil.NewChrootFs(afero.NewOsFs(), "."))
	if assert.NoError(t, err) {
		assert.Equal(t, "package parse\n", fs.GetText(0, 13))
	}
}

func TestNewFSFileStreamNotFound(t *testing.T) {
	t.Parallel()

	_, err := newFSFileStream("x", syslutil.NewChrootFs(afero.NewOsFs(), "/non-existent.dir"))
	assert.Error(t, err)
	_, err = newFSFileStream("non-existent.file", syslutil.NewChrootFs(afero.NewOsFs(), "."))
	assert.Error(t, err)
}

type flappyFile struct {
	data    []byte
	succeed bool
}

func (ff *flappyFile) Close() error {
	return nil
}

func (ff *flappyFile) Read(p []byte) (n int, err error) {
	N := len(ff.data)
	if len(p) >= N {
		if ff.succeed {
			copy(p, ff.data)
			return N, io.EOF
		}
		return 0, fmt.Errorf("Read failed")
	}
	copy(p, ff.data)
	ff.data = (ff.data)[:len(p)]
	logrus.Infof("ff.data = %#v", ff.data)
	return len(p), nil
}

func (*flappyFile) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("Seek() not implemented")
}

func (*flappyFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("Not a directory")
}

func (*flappyFile) Stat() (os.FileInfo, error) {
	return nil, fmt.Errorf("Stat() not implemented")
}

// An http.FileSystem in which every Open returns a file containing the given
// bytes, but returns an error for any attempt to read more data.
type flappyFileSystem struct {
	data    []byte
	succeed bool
}

func (ffs flappyFileSystem) Open(name string) (http.File, error) {
	ff := flappyFile(ffs)
	return &ff, nil
}

func TestFlappyFileSystem(t *testing.T) {
	t.Parallel()

	f, err := flappyFileSystem{[]byte("package ma"), false}.Open("won't.go")
	assert.NoError(t, err)

	_, err = f.Seek(0, io.SeekStart)
	assert.Error(t, err)

	_, err = f.Readdir(-1)
	assert.Error(t, err)

	_, err = f.Stat()
	assert.Error(t, err)

	var p [7]byte
	n, err := f.Read(p[:])
	if assert.NoError(t, err) {
		assert.Len(t, p, n)
		assert.Equal(t, "package", string(p[:]))
	}

	_, err = f.Read(p[:])
	assert.Error(t, err)

	content := "package main\n"
	f, err = flappyFileSystem{[]byte(content), true}.Open("will.go")
	assert.NoError(t, err)

	var q [20]byte
	n, err = f.Read(q[:])
	if assert.Equal(t, io.EOF, err) {
		if assert.Len(t, content, n) {
			assert.Equal(t, content, string(q[:n]))
		}
	}
}
