package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByFilepath(t *testing.T) {
	t.Parallel()

	var testMods goModules
	testMods = append(testMods, &goModule{Path: "modulepath"})
	assert.Equal(t, &goModule{Path: "modulepath"}, testMods.GetByFilepath("modulepath/filename"))
	assert.Equal(t, &goModule{Path: "modulepath"}, testMods.GetByFilepath(".//modulepath/filename"))
}

func TestGetByFilepathWithoutValidMod(t *testing.T) {
	t.Parallel()

	var testMods goModules
	testMods = append(testMods, &goModule{Path: "modulepath"})
	assert.Nil(t, testMods.GetByFilepath("another_modulepath/filename"))
}

func TestGetByFilepathWithNilMods(t *testing.T) {
	t.Parallel()

	var testMods goModules
	assert.Nil(t, testMods.GetByFilepath("modulepath/filename"))
}
