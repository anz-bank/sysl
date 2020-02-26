package main

import (
	"errors"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/mod"
	"gopkg.in/alecthomas/kingpin.v2"
)

type modCmd struct {
	modName string
}

var (
	initCmd *kingpin.CmdClause
)

func (m *modCmd) Name() string       { return "mod" }
func (m *modCmd) MaxSyslModule() int { return 0 }

func (m *modCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(m.Name(), "provides access to operations on Sysl modules")
	initCmd = cmd.Command("init", "initializes and writes a new go.mod to the current directory")
	initCmd.Arg("name", "name of the sysl module").StringVar(&m.modName)

	return cmd
}

func (m *modCmd) Execute(args cmdutils.ExecuteArgs) error {
	// subcommands are somewhat funky using the cmd_runner
	switch args.Command {
	case initCmd.FullCommand():
		return mod.SyslModInit(m.modName, args.Logger)
	default:
		return errors.New("command not recognized")
	}
}
