package mod

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Equal(t, 1, len(testMods))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods[0])
}

func TestGetByFilename(t *testing.T) {
	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename("modulepath/filename", ""))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename("modulepath/filename2", ""))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename(".//modulepath/filename", ""))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods.GetByFilename("modulepath", ""))
	assert.Nil(t, testMods.GetByFilename("modulepath2/filename", ""))
}

func TestGetByFilepathWithoutValidMod(t *testing.T) {
	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Nil(t, testMods.GetByFilename("another_modulepath/filename", ""))
}

func TestGetByFilepathWithNilMods(t *testing.T) {
	var testMods Modules
	assert.Nil(t, testMods.GetByFilename("modulepath/filename", ""))
}

func TestFind(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	filename := "github.com/anz-bank/sysl/tests/deps.sysl"
	mod, err := Find(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, "github.com/anz-bank/sysl", mod.Name)
}

func TestFindWithWrongPath(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	wrongpath := "wrong_file_path/deps.sysl"
	mod, err := Find(wrongpath, "")
	assert.Equal(t, fmt.Sprintf("%s not found", wrongpath), err.Error())
	assert.Nil(t, mod)
}
