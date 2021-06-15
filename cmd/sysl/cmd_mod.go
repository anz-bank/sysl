package main

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/mod"
	"gopkg.in/alecthomas/kingpin.v2"
)

type modCmd struct {
	repo string
}

func (c *modCmd) Name() string       { return "mod" }
func (c *modCmd) MaxSyslModule() int { return 0 }

func (c *modCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(c.Name(), "Provides access to operations on modules.")
	cmd.Command("init", "initialize new module in current directory")
	getCmd := cmd.Command("get", "initialize new module in current directory")
	getCmd.Arg("repo", "get and add module").StringVar(&c.repo)
	updateCmd := cmd.Command("update", "add missing and remove unused modules")
	updateCmd.Arg("repo", "repo to update. update all if repo is not specified.").StringVar(&c.repo)
	return cmd
}

func (c *modCmd) Execute(args cmdutils.ExecuteArgs) error {
	subCmd := strings.TrimPrefix(args.Command, "mod ")
	retr, err := mod.Retriever(args.Filesystem)
	if err != nil {
		return err
	}
	return retr.Command(subCmd, c.repo)
}
