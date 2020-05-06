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

func (m *Modules) GetByFilename(filename string) *Module {
	for _, mod := range *m {
		if hasPathPrefix(mod.Name, filename) {
			return mod
		}
	}

	return nil
}

func Find(name string) (*Module, error) {
	err := SyslModInit("")
	if err != nil && !strings.HasSuffix(err.Error(), "go.mod already exists\n") {
		return nil, err
	}

	err = goGetByFilepath(name)
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
	prefix = filepath.Clean(prefix)
	s = filepath.Clean(s)

	if len(s) > len(prefix) {
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}

	return s == prefix
}
