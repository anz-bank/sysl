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
	local := &Module{Name: "local"}
	mod2 := &Module{Name: "remote", Version: "v0.2.0"}
	testMods.Add(local)
	testMods.Add(mod2)

	assert.Equal(t, local, gomod.Find("local/filename", "", &testMods))
	assert.Equal(t, local, gomod.Find("local/filename2", "", &testMods))
	assert.Equal(t, local, gomod.Find(".//local/filename", "", &testMods))
	assert.Equal(t, local, gomod.Find("local", "", &testMods))
	assert.Nil(t, gomod.Find("local2/filename", "", &testMods))

	assert.Equal(t, local, gomod.Find("local/filename", MasterBranch, &testMods))
	assert.Equal(t, local, gomod.Find("local/filename", "v0.0.1", &testMods))

	assert.Equal(t, mod2, gomod.Find("remote/filename", "v0.2.0", &testMods))
	assert.Nil(t, gomod.Find("remote/filename", "v1.0.0", &testMods))
}
