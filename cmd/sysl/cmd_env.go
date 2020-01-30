package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

const KnownEnv = `
	SYSL_MODULES
	SYSL_PLANTUML
`

type envCmd struct{}

func (c *envCmd) Name() string       { return "env" }
func (c *envCmd) MaxSyslModule() int { return 0 }

func (c *envCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(c.Name(), "Print sysl environment information.")
	return cmd
}

func (c *envCmd) Execute(args ExecuteArgs) error {
	for _, e := range strings.Fields(KnownEnv) {
		v := os.Getenv(e)
		fmt.Printf("%s=\"%s\"\n", e, v)

	}

	return nil
}
