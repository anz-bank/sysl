package main

import (
	"errors"
	"fmt"
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
	switch m.subcommand {
	case "init":
		return syslModInit(args)
	}
	return errors.New("command not recognized")
}

func syslModInit(args ExecuteArgs) error {
	// makes the assumption that the CWD is not a go module since we hijack this command
	out, err := exec.Command("go", "mod", "init", syslModuleName).CombinedOutput()
	if err != nil {
		return err
	}

	args.Logger.Info(string(out))
	return nil
}
