package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/anz-bank/sysl/pkg/cmdutils"
)

type validateCmd struct{}

func (p *validateCmd) Name() string       { return "validate" }
func (p *validateCmd) MaxSyslModule() int { return 1 }

func (p *validateCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Validate the sysl file")
	return cmd
}

func (p *validateCmd) Execute(args cmdutils.ExecuteArgs) error {
	// Nothing to do here, the runner loads the sysl file automatically. If we got here the file was successfully loaded
	return nil
}
