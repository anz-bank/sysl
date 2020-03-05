package main

import (
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/testrig"
	"gopkg.in/alecthomas/kingpin.v2"
)

type testRigCmd struct {
	TemplateFileName *string
	OutputDir        *string
}

func (p *testRigCmd) Name() string {
	return "test-rig"
}

func (p *testRigCmd) MaxSyslModule() int {
	return 1
}

func (p *testRigCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate test rig")
	p.TemplateFileName = cmd.Flag("template", "variables file name (json)").Required().ExistingFile()
	p.OutputDir = cmd.Flag("output-dir", "directory to put generated files").Default(".").String()
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *testRigCmd) Execute(args cmdutils.ExecuteArgs) error {
	var err error
	err = testrig.GenerateRig(*p.TemplateFileName, *p.OutputDir, args.Modules)
	if err != nil {
		return err
	}
	return nil
}
