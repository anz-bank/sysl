package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubMgrGet(t *testing.T) {
	githubmod := &githubMgr{}
	githubmod.Init()
	testMods := Modules{}

	mod, err := githubmod.Get(RemoteDepsFile, "", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = githubmod.Get(RemoteDepsFile, MasterBranch, &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = githubmod.Get(RemoteDepsFile, "v0.0.1", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)

	mod, err = githubmod.Get("github.com/anz-bank/wrong/path", "", &testMods)
	assert.Error(t, err)
	assert.Nil(t, mod)
}

func TestGitHubMgrFind(t *testing.T) {
	githubmod := &githubMgr{}
	testMods := Modules{}
	master := &Module{Name: MasterBranch, Version: "v0.0.1"}
	mod1 := &Module{Name: "modulepath", Version: "v1.0.0-41f04d3bba15"}
	mod2 := &Module{Name: "modulepath", Version: "v0.2.0"}
	testMods.Add(master)
	testMods.Add(mod1)
	testMods.Add(mod2)

	assert.Equal(t, master, githubmod.Find("master/filename", "v0.0.1", &testMods))
	assert.Equal(t, master, githubmod.Find("master/filename2", "v0.0.1", &testMods))
	assert.Equal(t, master, githubmod.Find(".//master/filename", "v0.0.1", &testMods))
	assert.Equal(t, master, githubmod.Find(MasterBranch, "v0.0.1", &testMods))
	assert.Nil(t, githubmod.Find("master2/filename", "v0.0.1", &testMods))

	assert.Nil(t, githubmod.Find("master/filename", MasterBranch, &testMods))

	assert.Equal(t, mod2, githubmod.Find("modulepath/filename", "v0.2.0", &testMods))
	assert.Nil(t, githubmod.Find("modulepath/filename", "v1.0.0", &testMods))
}
