package main

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/env"
	"gopkg.in/alecthomas/kingpin.v2"
)

type envCmd struct{}

func (c *envCmd) Name() string       { return "env" }
func (c *envCmd) MaxSyslModule() int { return 0 }

func (c *envCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	return app.Command(c.Name(), "Print sysl environment information.")
}

func (c *envCmd) Execute(args cmdutils.ExecuteArgs) error {
	for _, e := range env.VARS {
		fmt.Printf("%s=\"%s\"\n", e, e.Value())
	}
	return nil
}
