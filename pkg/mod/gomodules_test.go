package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoModulesGet(t *testing.T) {
	gomod := &goModules{}
	testMods := Modules{}

	mod, err := gomod.Get(RemoteDepsFile, "", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = gomod.Get(RemoteDepsFile, MasterBranch, &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = gomod.Get(RemoteDepsFile, "v0.0.1", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)

	mod, err = gomod.Get("github.com/anz-bank/wrongpath", "", &testMods)
	assert.Error(t, err)
	assert.Nil(t, mod)
}

func TestGoModulesFind(t *testing.T) {
	gomod := &goModules{}
	testMods := Modules{}
	master := &Module{Name: MasterBranch}
	testMods.Add(master)
	testMods.Add(&Module{Name: "modulepath", Version: "v1.0.0-20190902080502-41f04d3bba15"})
	testMods.Add(&Module{Name: "modulepath", Version: "v0.2.0"})

	assert.Equal(t, master, gomod.Find("master/filename", "", &testMods))
	assert.Equal(t, master, gomod.Find("master/filename2", "", &testMods))
	assert.Equal(t, master, gomod.Find(".//master/filename", "", &testMods))
	assert.Equal(t, master, gomod.Find(MasterBranch, "", &testMods))
	assert.Nil(t, gomod.Find("master2/filename", "", &testMods))

	assert.Equal(t, master, gomod.Find("master/filename", MasterBranch, &testMods))
	assert.Nil(t, gomod.Find("master/filename", "v0.0.1", &testMods))

	assert.Equal(t, &Module{Name: "modulepath", Version: "v0.2.0"}, gomod.Find("modulepath/filename", "v0.2.0", &testMods))
	assert.Nil(t, gomod.Find("modulepath/filename", "v1.0.0", &testMods))
}
