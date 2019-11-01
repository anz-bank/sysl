package main

import (
	"github.com/anz-bank/sysl/sysl2/sysl/eval"
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

	s := eval.Scope{}
	for name, app := range args.Module.Apps {
		s.AddApp(name, app)
	}
	return nil
}
