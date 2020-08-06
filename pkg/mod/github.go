package mod

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
)

type githubMgr struct {
	client *github.Client
}

var syslModulesCacheDir string

func (d *githubMgr) Init() {
	if d.client == nil {
		d.client = github.NewClient(nil)

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		syslModulesCacheDir = filepath.Join(usr.HomeDir, ".sysl")
	}
}

func (d *githubMgr) Get(filename, ver string, m *Modules) (*Module, error) {
	repoPath, err := getGitHubRepoPath(filename)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	var ref *github.RepositoryContentGetOptions
	if ver != "" {
		ref = &github.RepositoryContentGetOptions{Ref: ver}
	}

	var fileContent *github.RepositoryContent
	for i := 0; i < 10; i++ {
		fileContent, _, _, err = d.client.Repositories.GetContents(ctx, repoPath.owner, repoPath.repo, repoPath.path, ref)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); !ok {
				return nil, err
			}

			log.Println("hit GitHub rate limit: ", i)
			if i == 9 {
				return nil, err
			}
			continue
		}
		break
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}
	if ver == "" {
		ver = "v0.0.0-" + fileContent.GetSHA()[:12]
	}

	name := strings.Join([]string{"github.com", repoPath.owner, repoPath.repo}, "/")
	dir := filepath.Join(syslModulesCacheDir, "github.com", repoPath.owner, repoPath.repo)
	dir = AppendVersion(dir, ver)
	new := &Module{
		Name:    name,
		Dir:     dir,
		Version: ver,
	}

	fname := filepath.Join(dir, repoPath.path)
	if !fileExists(fname) {
		err = writeFile(fname, []byte(content))
		if err != nil {
			return nil, err
		}
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

type githubRepoPath struct {
	owner string
	repo  string
	path  string
}

func getGitHubRepoPath(filename string) (*githubRepoPath, error) {
	names := strings.Split(filename, "/")
	if len(names) < 4 {
		return nil, fmt.Errorf("the imported module path %s is invalid", filename)
	}
	if names[0] != "github.com" {
		return nil, errors.New("non-github.com repository is not supported under GitHub mode")
	}

	owner := names[1]
	repo := names[2]
	path := path.Join(names[3:]...)

	return &githubRepoPath{
		owner: owner,
		repo:  repo,
		path:  path,
	}, nil
}

func writeFile(filename string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}
