package parser

import (
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

func LoadAndGetDefaultApp(filename string, fs afero.Fs) (*sysl.Module, string, error) {
	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return &sysl.Module{}, "", err
	}
	tree, err := ParseString(string(data))
	if err != nil {
		return &sysl.Module{}, "", err
	}
	var defaultApp string
	mod := &sysl.Module{}
	mod.Apps = map[string]*sysl.Application{}
	for _, app := range tree.AllApplication() {
		name, app := buildApplication(app)
		mod.Apps[name] = app
		if defaultApp == "" {
			defaultApp = name
		}
	}

	return mod, defaultApp, nil
}
