package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/mod"
	"gopkg.in/alecthomas/kingpin.v2"
)

type modCmd struct {
	subcommand string
}

func (m *modCmd) Name() string       { return "mod" }
func (m *modCmd) MaxSyslModule() int { return 0 }

func (m *modCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	opts := []string{"init"}
	cmd := app.Command(m.Name(), "Configure sysl module system")
	cmd.Arg("command", fmt.Sprintf("commands: [%s]", strings.Join(opts, ","))).Required().EnumVar(&m.subcommand, opts...)

	return cmd
}

func (m *modCmd) Execute(args ExecuteArgs) error {
	if m.subcommand == "init" {
		return mod.SyslModInit(args.Logger)
	}

	return errors.New("command not recognized")
}
