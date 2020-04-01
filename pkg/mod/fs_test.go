package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
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
	t.Parallel()

	filename := "deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, "../../tests/")
	f, err := fs.Open(filename)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenLocalFileFailed(t *testing.T) {
	t.Parallel()

	filename := "wrong.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, "../../tests/")
	f, err := fs.Open(filename)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found", filepath.Join("../../tests/", filename)), err.Error())
}

func TestOpenRemoteFile(t *testing.T) {
	t.Parallel()

	filename := "github.com/anz-bank/sysl/demo/examples/Modules/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, "")
	f, err := fs.Open(filename)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenRemoteFileFailed(t *testing.T) {
	t.Parallel()

	filename := "github.com/wrong/repo/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, "")
	f, err := fs.Open(filename)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found", filepath.FromSlash(filename)), err.Error())
}

func TestOpenRemoteFileWithRoot(t *testing.T) {
	t.Parallel()

	root := "github.com/anz-bank/sysl"
	path := "demo/examples/Modules/deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, root)
	f, err := fs.Open(path)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenFile(t *testing.T) {
	t.Parallel()

	filename := "deps.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, "../../tests/")
	f, err := fs.OpenFile(filename, os.O_RDWR, 0600)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenFileFailed(t *testing.T) {
	t.Parallel()

	filename := "wrong.sysl"
	_, memfs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(memfs, "../../tests/")
	f, err := fs.OpenFile(filename, os.O_RDWR, 0600)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found", filepath.Join("../../tests/", filename)), err.Error())
}
