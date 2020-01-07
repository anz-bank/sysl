package mod

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExternalFile(t *testing.T) {
	t.Parallel()

	filename := "github.com/anz-bank/sysl/cmd/sysl/deps.sysl"
	importFilename, err := GetExternalFile(filename)
	assert.Nil(t, err)

	mod := GoMods.GetByFilepath("github.com/anz-bank/sysl/cmd/sysl/deps.sysl")
	assert.Equal(t, filepath.Join(mod.Dir, "/cmd/sysl/deps.sysl"), importFilename)
}

func TestGetExternalFileWithWrongPath(t *testing.T) {
	t.Parallel()

	wrongpath := "wrong_file_path/deps.sysl"
	importFilename, err := GetExternalFile(wrongpath)
	assert.Equal(t, wrongpath, importFilename)
	assert.Equal(t, fmt.Sprintf("%s not found\n", wrongpath), err.Error())
}
