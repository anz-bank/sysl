package main

import (
	"errors"
	"strings"

	"github.com/anz-bank/sysl/pkg/mod"
	"gopkg.in/alecthomas/kingpin.v2"
)

type modCmd struct {
	modName string
}

func (m *modCmd) Name() string       { return "mod" }
func (m *modCmd) MaxSyslModule() int { return 0 }

func (m *modCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(m.Name(), "Configure sysl module system")
	initCmd := cmd.Command("init", "sysl module init command")
	initCmd.Arg("mod name", "Name of sysl module").Required().StringVar(&m.modName)

	return cmd
}

func (m *modCmd) Execute(args ExecuteArgs) error {
	// subcommands are somewhat funky using the cmd_runner
	subcommand := strings.Split(args.Command, " ")[1]
	if subcommand == "init" {
		return mod.SyslModInit(m.modName, args.Logger)
	}

	return errors.New("command not recognized")
}
