package main

import (
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/test_rig"
	"gopkg.in/alecthomas/kingpin.v2"
)

type testRigCmd struct {
	cmdutils.CmdContextParamTestRig
}

func (p *testRigCmd) Name() string {
	return "test-rig"
}

func (p *testRigCmd) MaxSyslModule() int {
	return 1
}

func (p *testRigCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate test rig")
	cmd.Flag("template", "variables file name (json)").Required().StringVar(&p.TemplateFileName)
	cmd.Flag("output-dir", "directory to put generated files").StringVar(&p.OutputDir)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *testRigCmd) Execute(args cmdutils.ExecuteArgs) error {
	var err error
	err = refineCmd(p)
	if err != nil {
		return err
	}
	err = test_rig.GenerateRig(p.TemplateFileName, p.OutputDir, args.Modules)
	if err != nil {
		return err
	}
	return nil
}

func refineCmd(p *testRigCmd) error {
	// vars file should exist
	// app-name should match what we have in sysl modules

	return nil
}
