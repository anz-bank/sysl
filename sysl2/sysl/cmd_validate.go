package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type validateCmd struct{}

func (p *validateCmd) Name() string            { return "validate" }
func (p *validateCmd) RequireSyslModule() bool { return true }

func (p *validateCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Validate the sysl file")
	return cmd
}

func (p *validateCmd) Execute(args ExecuteArgs) error {
	// Nothing to do here, the runner loads the sysl file automatically. If we got here the file was successfully loaded
	return nil
}
