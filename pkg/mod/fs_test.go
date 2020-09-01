package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewFs(t *testing.T) {
	t.Parallel()

	_, backendFs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(backendFs, "")
	assert.Equal(t, "ChrootFS", fs.Fs.Name())
}

func TestFsName(t *testing.T) {
	t.Parallel()

	_, backendFs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(backendFs, "")
	assert.Equal(t, "ModSupportedFs", fs.Name())
}

func TestOpenLocalFile(t *testing.T) {
	SyslModules = false
	defer func() { SyslModules = true }()

	filename := "deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "../../tests/")
	f, err := mfs.Open(filename)
	assert.NoError(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenLocalFileFailed(t *testing.T) {
	SyslModules = false
	defer func() { SyslModules = true }()

	filename := "wrong.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "../../tests/")
	f, err := mfs.Open(filename)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found: no such file in current working directory", filename), err.Error())
}

func TestOpenRemoteFile(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	filename := "github.com/anz-bank/sysl/tests/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "")
	f, err := mfs.Open(filename)
	assert.NoError(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenRemoteFileFailed(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	filename := "github.com/wrong/repo/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "")
	f, err := mfs.Open(filename)
	assert.Nil(t, f)
	assert.Error(t, err)
}

func TestOpenRemoteFileWithRoot(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	root := "github.com/anz-bank/sysl"
	path := "demo/examples/Modules/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, root)
	f, err := mfs.Open(path)
	assert.NoError(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenFile(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	filename := "deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "../../tests/")
	f, err := mfs.OpenFile(filename, os.O_RDWR, 0600)
	assert.NoError(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenFileFailed(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	filename := "wrong.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "../../tests/")
	f, err := mfs.OpenFile(filename, os.O_RDWR, 0600)
	assert.Nil(t, f)
	assert.Error(t, err)
}

func removeGomodFile(t *testing.T, fs afero.Fs) {
	removeFile(t, fs, "go.mod")
	removeFile(t, fs, "go.sum")
}

func createGomodFile(t *testing.T, fs afero.Fs) {
	gomod, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer gomod.Close()
	_, err = gomod.WriteString("module github.com/anz-bank/sysl/pkg/mod")
	assert.NoError(t, err)
	err = gomod.Sync()
	assert.NoError(t, err)
}

func removeFile(t *testing.T, fs afero.Fs, file string) {
	exists, err := afero.Exists(fs, file)
	assert.NoError(t, err)
	if exists {
		err = fs.Remove(file)
		assert.NoError(t, err)
	}
}
