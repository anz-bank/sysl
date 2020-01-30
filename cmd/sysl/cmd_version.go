package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"
)

type versionCmd struct{}

func (c *versionCmd) Name() string       { return "version" }
func (c *versionCmd) MaxSyslModule() int { return 0 }

func (c *versionCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(c.Name(), "Print the build information for Sysl executables.")
	return cmd
}

func (c *versionCmd) Execute(args ExecuteArgs) error {
	tag, err := latestGitVersionTag()
	if err != nil {
		return err
	}

	fmt.Printf("sysl version %s %s/%s\n", tag, runtime.GOOS, runtime.GOARCH)
	return nil
}

func latestGitVersionTag() (string, error) {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0", "--match=v[0-9]*").Output()
	if err != nil {
		return "", err
	}
	return string(out[:len(out)-1]), nil
}
