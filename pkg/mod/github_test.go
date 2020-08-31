package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubMgrGet(t *testing.T) {
	githubmod := &githubMgr{}
	githubmod.Init(".sysl")
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
	local := &Module{Name: "local", Version: "v0.0.1"}
	mod1 := &Module{Name: "remote", Version: "v1.0.0-41f04d3bba15"}
	mod2 := &Module{Name: "remote", Version: "v0.2.0"}
	testMods.Add(local)
	testMods.Add(mod1)
	testMods.Add(mod2)

	assert.Equal(t, local, githubmod.Find("local/filename", "v0.0.1", &testMods))
	assert.Equal(t, local, githubmod.Find("local/filename2", "v0.0.1", &testMods))
	assert.Equal(t, local, githubmod.Find(".//local/filename", "v0.0.1", &testMods))
	assert.Equal(t, local, githubmod.Find("local", "v0.0.1", &testMods))
	assert.Nil(t, githubmod.Find("local2/filename", "v0.0.1", &testMods))

	assert.Nil(t, githubmod.Find("local/filename", MasterBranch, &testMods))

	assert.Equal(t, mod2, githubmod.Find("remote/filename", "v0.2.0", &testMods))
	assert.Nil(t, githubmod.Find("remote/filename", "v1.0.0", &testMods))
}

func TestGetGitHubRepoPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		filename string
		path     *githubRepoPath
	}{
		{"github.com/anz-bank/sysl", nil},
		{"github.com/anz-bank/sysl/", nil},
		{"github.com/anz-bank/sysl/deps.sysl", &githubRepoPath{"anz-bank", "sysl", "deps.sysl"}},
		{"github.com/anz-bank/sysl/nested/module/deps.sysl", &githubRepoPath{"anz-bank", "sysl", "nested/module/deps.sysl"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			p, err := getGitHubRepoPath(tt.filename)
			if tt.path == nil {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.path, p)
		})
	}
}
