package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

const syslModuleName = "syslmodules"

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
		err := syslModInit(args)
		if err != nil {
			return err
		}
	}
	return nil
}

func syslModInit(args ExecuteArgs) error {
	// ignore folder creation error
	_err := args.Filesystem.Mkdir(syslRootMarker, 0755)
	if _err != nil {
		args.Logger.Warn(_err.Error())
	}

	err := os.Chdir(syslRootMarker)
	if err != nil {
		return err
	}

	out, err := exec.Command("go", "mod", "init", syslModuleName).CombinedOutput()
	if err != nil {
		return err
	}

	args.Logger.Debug(out)
	return nil
}
