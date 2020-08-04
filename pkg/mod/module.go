package mod

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var SyslModules = os.Getenv("SYSL_MODULES") != "off"
var GitHubMode = os.Getenv("SYSL_GITHUB") == "on"

type Module struct {
	Name    string // "github.com/anz-bank/sysl"
	Dir     string // "/Users/foo/go/pkg/mod/github.com/anz-bank/sysl@v1.1.0"
	Version string // "v1.1.0"
}

type Modules []*Module

var modules Modules

type DependencyManager interface {
	Get(filename, ver string, m *Modules) (*Module, error)
	Find(filename, ver string, m *Modules) *Module
	Load(m *Modules) error
}

func (m *Modules) Add(v *Module) {
	*m = append(*m, v)
}

func (m *Modules) Len() int {
	return len(*m)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Find(name string, ver string) (*Module, error) {
	if !fileExists("go.mod") {
		return nil, errors.New("no go.mod file, run `go mod init` first")
	}

	var manager DependencyManager
	if GitHubMode {
		manager = &githubMgr{}
	} else {
		manager = &goModules{}
	}

	if modules.Len() == 0 {
		if err := manager.Load(&modules); err != nil {
			return nil, fmt.Errorf("error loading modules: %s", err.Error())
		}
	}

	mod := manager.Find(name, ver, &modules)
	if mod != nil {
		return mod, nil
	}

	return manager.Get(name, ver, &modules)
}

func hasPathPrefix(prefix, s string) bool {
	prefix = filepath.Clean(prefix)
	s = filepath.Clean(s)

	if len(s) > len(prefix) {
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}

	return s == prefix
}
