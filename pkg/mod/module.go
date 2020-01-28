package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var SyslModules = os.Getenv("SYSL_MODULES") != "off"

type Module struct {
	Name string
	Dir  string
}

type Modules []*Module

func (m *Modules) Add(v *Module) {
	*m = append(*m, v)
}

func (m *Modules) GetByFilename(p string) *Module {
	p = filepath.Clean(p)

	for _, v := range *m {
		if hasPathPrefix(filepath.Clean(v.Name), p) {
			return v
		}
	}

	return nil
}

func Find(name string) (*Module, error) {
	err := goGetByFilepath(name)
	if err != nil {
		return nil, fmt.Errorf("%s not found", name)
	}

	var modules Modules
	if err = modules.Load(); err != nil {
		return nil, fmt.Errorf("error loading modules: %s", err.Error())
	}

	if modules == nil {
		return nil, fmt.Errorf("modules list is empty")
	}

	mod := modules.GetByFilename(name)
	if mod == nil {
		return nil, fmt.Errorf("error finding module of file %s", name)
	}
	return mod, nil
}

func hasPathPrefix(prefix, s string) bool {
	switch {
	default:
		return false
	case len(s) == len(prefix):
		return s == prefix
	case len(s) > len(prefix):
		if prefix != "" && prefix[len(prefix)-1] == filepath.Separator {
			return strings.HasPrefix(s, prefix)
		}
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}
}
