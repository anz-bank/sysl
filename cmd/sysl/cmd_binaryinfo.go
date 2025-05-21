package main

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/alecthomas/kingpin/v2"
)

type infoCmd struct{}

func (p *infoCmd) Name() string       { return "info" }
func (p *infoCmd) MaxSyslModule() int { return 0 }

func (p *infoCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Show binary information")
	return cmd
}

func (p *infoCmd) Execute(args cmdutils.ExecuteArgs) error {
	fmt.Fprintf(args.Stdout, "Build:\n")
	fmt.Fprintf(args.Stdout, "  Version      : %s\n", Version)
	fmt.Fprintf(args.Stdout, "  Git Commit   : %s\n", GitFullCommit)
	fmt.Fprintf(args.Stdout, "  Date         : %s\n", BuildDate)
	fmt.Fprintf(args.Stdout, "  Go Version   : %s\n", GoVersion)
	fmt.Fprintf(args.Stdout, "  OS           : %s\n", BuildOS)

	return nil
}
