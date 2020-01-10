package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var SyslModules = os.Getenv("SYSL_MODULES") != ""

type Module struct {
	Name string
	Dir  string
}

type Modules []*Module

func (modules *Modules) Add(m *Module) {
	*modules = append(*modules, m)
}

func (modules Modules) GetByFilename(p string) *Module {
	if modules == nil {
		return nil
	}

	p = filepath.Clean(p)

	for _, m := range modules {
		if strings.HasPrefix(p, filepath.Clean(m.Name)) {
			return m
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

	mod := modules.GetByFilename(name)
	if mod == nil {
		return nil, fmt.Errorf("error finding module of file %s", name)
	}
	return mod, nil
}
