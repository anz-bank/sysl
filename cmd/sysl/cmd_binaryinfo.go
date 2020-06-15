package main

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"gopkg.in/alecthomas/kingpin.v2"
)

type infoCmd struct{}

func (p *infoCmd) Name() string       { return "info" }
func (p *infoCmd) MaxSyslModule() int { return 0 }

func (p *infoCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Show binary information")
	return cmd
}

func (p *infoCmd) Execute(args cmdutils.ExecuteArgs) error {
	fmt.Printf("Build:\n")
	fmt.Printf("  Version      : %s\n", Version)
	fmt.Printf("  Git Commit   : %s\n", GitFullCommit)
	fmt.Printf("  Date         : %s\n", BuildDate)
	fmt.Printf("  Go Version   : %s\n", GoVersion)
	fmt.Printf("  OS           : %s\n", BuildOS)

	return nil
}
