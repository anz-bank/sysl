package mod

import (
	"fmt"
	"os"
	"path"
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
	assert.Equal(t, fmt.Sprintf("%s not found", filename), err.Error())
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
	assert.Equal(t, fmt.Sprintf("%s not found", path.Clean(filename)), err.Error())
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
	assert.Equal(t, fmt.Sprintf("%s not found", path.Join("../../tests/", filename)), err.Error())
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

func TestExtractVersion(t *testing.T) {
	path, ver := ExtractVersion("github.com/anz-bank/sysl@v0.1")
	assert.Equal(t, "github.com/anz-bank/sysl", path)
	assert.Equal(t, "v0.1", ver)

	path, ver = ExtractVersion("github.com/anz-bank/sysl/pkg@v0.2")
	assert.Equal(t, "github.com/anz-bank/sysl/pkg", path)
	assert.Equal(t, "v0.2", ver)

	path, ver = ExtractVersion("github.com/anz-bank/sysl/pkg")
	assert.Equal(t, "github.com/anz-bank/sysl/pkg", path)
	assert.Equal(t, "", ver)
}
