package mod

import (
	"context"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
)

var client *github.Client

type githubMgr struct{}

var syslModulesCacheDir string

func (*githubMgr) Init() {
	if client == nil {
		client = github.NewClient(nil)

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		syslModulesCacheDir = filepath.Join(usr.HomeDir, ".sysl")
	}
}

func (*githubMgr) Get(filename, ver string, m *Modules) (*Module, error) {
	names := strings.Split(filename, "/")
	if len(names) == 0 {
		return nil, errors.New("Empty path is not supported")
	}
	if names[0] != "github.com" {
		return nil, errors.New("Non-github.com repository is not supported under GitHub mode")
	}

	owner := names[1]
	repo := names[2]
	path := path.Join(names[3:]...)
	ctx := context.Background()
	var ref *github.RepositoryContentGetOptions
	if ver != "" {
		ref = &github.RepositoryContentGetOptions{Ref: ver}
	}

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, ref)
	if err != nil {
		return nil, err
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}
	if ver == "" {
		ver = "v0.0.0-" + fileContent.GetSHA()[:12]
	}

	name := filepath.Join("github.com", owner, repo)
	dir := filepath.Join(syslModulesCacheDir, name)
	dir = AppendVersion(dir, ver)
	fname := filepath.Join(dir, path)
	new := &Module{
		Name:    strings.Join([]string{"github.com", owner, repo}, "/"),
		Dir:     dir,
		Version: ver,
	}
	if !fileExists(fname) {
		if err = os.MkdirAll(filepath.Dir(fname), 0770); err != nil {
			return nil, err
		}
		file, err := os.Create(fname)
		if err != nil {
			return nil, err
		}
		if _, err = file.Write([]byte(content)); err != nil {
			return nil, err
		}
		defer file.Close()
		m.Add(new)
	}

	return new, nil
}

func (*githubMgr) Find(filename, ver string, m *Modules) *Module {
	if ver == "" || ver == MasterBranch {
		return nil
	}

	for _, mod := range *m {
		if hasPathPrefix(mod.Name, filename) {
			if mod.Version == ver {
				return mod
			}
		}
	}

	return nil
}

func (*githubMgr) Load(m *Modules) error {
	githubPath := filepath.Join(syslModulesCacheDir, "github.com")
	githubDir, err := os.Open(githubPath)
	if err != nil {
		return err
	}

	owners, err := githubDir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, owner := range owners {
		ownerDir, err := os.Open(filepath.Join(githubPath, owner))
		if err != nil {
			return err
		}
		repos, err := ownerDir.Readdirnames(-1)
		if err != nil {
			return err
		}
		for _, repo := range repos {
			p, ver := ExtractVersion(repo)
			name := filepath.Join("github.com", owner, p)
			m.Add(&Module{
				Name:    name,
				Dir:     filepath.Join(ownerDir.Name(), repo),
				Version: ver,
			})
		}
	}

	return nil
}
