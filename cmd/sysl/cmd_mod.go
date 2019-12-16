package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type modCmd struct {}

func (p *modCmd) Name() string            { return "mod" }
func (p *modCmd) RequireSyslModule() bool { return false }

func (p *modCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Provide access to operations on modules")
	return cmd
}

func (p *modCmd) Execute(args ExecuteArgs) error {
	return nil
}
