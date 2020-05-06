package mod

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Equal(t, 1, len(testMods))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods[0])
}

func TestGetByFilename(t *testing.T) {
	t.Parallel()

	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename("modulepath/filename"))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename("modulepath/filename2"))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename(".//modulepath/filename"))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename("modulepath"))
	assert.Nil(t, testMods.GetByFilename("modulepath2/filename"))
}

func TestGetByFilepathWithoutValidMod(t *testing.T) {
	t.Parallel()

	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Nil(t, testMods.GetByFilename("another_modulepath/filename"))
}

func TestGetByFilepathWithNilMods(t *testing.T) {
	t.Parallel()

	var testMods Modules
	assert.Nil(t, testMods.GetByFilename("modulepath/filename"))
}

func TestFind(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	filename := "github.com/anz-bank/sysl/demo/examples/Modules/deps.sysl"
	mod, err := Find(filename)
	assert.Nil(t, err)
	assert.Equal(t, "github.com/anz-bank/sysl", mod.Name)
}

func TestFindWithWrongPath(t *testing.T) {
	fs := afero.NewOsFs()
	_, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer removeFile(t, fs, "go.mod")

	wrongpath := "wrong_file_path/deps.sysl"
	mod, err := Find(wrongpath)
	assert.Equal(t, fmt.Sprintf("%s not found", wrongpath), err.Error())
	assert.Nil(t, mod)
}
