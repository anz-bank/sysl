package main

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/eval"
)

type replCmd struct{}

func (p *replCmd) Name() string       { return "repl" }
func (p *replCmd) MaxSyslModule() int { return 0 }

func (p *replCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Enter a sysl REPL")
	return cmd
}

func (p *replCmd) Execute(args cmdutils.ExecuteArgs) error {
	s := &eval.Scope{}
	repl := eval.NewREPL(args.Stdin, args.Stdout)
	for {
		if err := repl(s, nil, nil); err != nil {
			return nil // means EOF
		}
	}
}
