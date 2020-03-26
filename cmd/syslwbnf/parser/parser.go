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

	p := newParserScope()
	if err := p.Parse(tree); err != nil {
		return &sysl.Module{}, "", err
	}

	return p.mod, p.defaultApp, nil
}

type pscope struct {
	// NOTE: Limit what actually ends up in this type please! and use accessor methods!
	mod        *sysl.Module
	defaultApp string

	context *sysl.Scope
}

func newParserScope() *pscope {
	return &pscope{}
}

func (p *pscope) Parse(node SyslFileNode) error {
	var defaultApp string
	mod := &sysl.Module{}
	mod.Apps = map[string]*sysl.Application{}
	for _, app := range node.AllApplication() {
		name, app := p.buildApplication(app)
		mod.Apps[name] = app
		if defaultApp == "" {
			defaultApp = name
		}
	}
	p.mod = mod
	p.defaultApp = defaultApp
	return nil
}
