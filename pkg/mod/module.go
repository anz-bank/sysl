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

func (m *Modules) Add(item *Module) {
	*m = append(*m, item)
}

func (m *Modules) GetByFilename(p string) *Module {
	p = filepath.Clean(p)

	for _, m := range m {
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

	if mod == nil {
		return nil, fmt.Error("modules list is empty")
	}

	mod := modules.GetByFilename(name)
	if mod == nil {
		return nil, fmt.Errorf("error finding module of file %s", name)
	}
	return mod, nil
}
