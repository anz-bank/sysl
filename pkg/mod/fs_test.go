package mod

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewFs(t *testing.T) {
	t.Parallel()

	backendFs := afero.NewOsFs()
	fs := NewFs(backendFs)
	assert.Equal(t, backendFs, fs.source)
}

func TestOpenLocalFile(t *testing.T) {
	t.Parallel()

	filename := "github.com/anz-bank/sysl/demo/examples/Modules/deps.sysl"
	fs := NewFs(afero.NewOsFs())
	f, err := fs.Open(filename)
	assert.Nil(t, err)
	assert.Equal(t, "deps.sysl", filepath.Base(f.Name()))
}

func TestOpenRemoteFile(t *testing.T) {
	t.Parallel()

	filename := "github.com/ChloePlanet/sysltestpub/pubbanana.sysl"
	fs := NewFs(afero.NewOsFs())
	f, err := fs.Open(filename)
	assert.Nil(t, err)
	assert.Equal(t, "pubbanana.sysl", filepath.Base(f.Name()))
}

func TestOpenRemoteFileFailed(t *testing.T) {
	t.Parallel()

	filename := "github.com/wrong/sysltestpub/pubbanana.sysl"
	fs := NewFs(afero.NewOsFs())
	f, err := fs.Open(filename)
	assert.Nil(t, f)
	assert.Equal(t, fmt.Sprintf("%s not found", filename), err.Error())
}
