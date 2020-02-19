package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type testRigCmd struct {
	CmdContextParamTestRig
}

func (p *testRigCmd) Name() string {
	return "test-rig"
}

func (p *testRigCmd) MaxSyslModule() int {
	return 1
}

func (p *testRigCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate test rig")
	cmd.Flag("vars", "variables file name (json)").Required().StringVar(&p.varFileName)
	cmd.Flag("app-names", "application names to parse").StringVar(&p.appNames)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *testRigCmd) Execute(args ExecuteArgs) error {
	return nil
}

func refineCmd(p *testRigCmd) error {
	// vars file should exist
	// app-names should match what we have in sysl modules

	return nil
}
