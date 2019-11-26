package main

import (
	"fmt"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"
)

type infoCmd struct{}

func (p *infoCmd) Name() string            { return "info" }
func (p *infoCmd) RequireSyslModule() bool { return false }

func (p *infoCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Show binary information")
	return cmd
}

func (p *infoCmd) Execute(args ExecuteArgs) error {
	fmt.Printf("Version    : %s\n", Version)
	fmt.Printf("Commit ID  : %s\n", GitCommit)
	fmt.Printf("Build date : %s\n", BuildDate)
	fmt.Printf("GOOS       : %s\n", runtime.GOOS)
	fmt.Printf("GOARCH     : %s\n", runtime.GOARCH)
	fmt.Printf("Go Version : %s\n", runtime.Version())
	fmt.Printf("Build OS   : %s\n", BuildOS)
	return nil
}
