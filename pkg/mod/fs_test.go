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
	assert.Nil(t, err)
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
	assert.Equal(t, fmt.Sprintf("%s not found", filename), err.Error())
}

func TestOpenRemoteFile(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	filename := "github.com/anz-bank/sysl/demo/examples/Modules/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "")
	f, err := mfs.Open(filename)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenRemoteFileFailed(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	filename := "github.com/wrong/repo/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "")
	f, err := mfs.Open(filename)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found", filepath.FromSlash(filename)), err.Error())
}

func TestOpenRemoteFileWithRoot(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	root := "github.com/anz-bank/sysl"
	path := "demo/examples/Modules/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, root)
	f, err := mfs.Open(path)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenFile(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	filename := "deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "../../tests/")
	f, err := mfs.OpenFile(filename, os.O_RDWR, 0600)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenFileFailed(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	filename := "wrong.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	mfs := NewFs(memfs, "../../tests/")
	f, err := mfs.OpenFile(filename, os.O_RDWR, 0600)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found", filepath.Join("../../tests/", filename)), err.Error())
}
