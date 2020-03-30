package mod

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/stretchr/testify/assert"
)

func TestNewFs(t *testing.T) {
	t.Parallel()

	_, backendFs := syslutil.WriteToMemOverlayFs("/")
	fs := NewFs(backendFs, "")
	assert.Equal(t, "ChrootFS", fs.source.Name())
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
	assert.Equal(t, fmt.Sprintf("%s not found", filename), err.Error())
}
