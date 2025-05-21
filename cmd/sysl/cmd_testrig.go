package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/testrig"
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
	templateFileHelp := `Variables file name (json format). Example content:
	{
		"dbfront": {
			"name": "dbfront",
			"import": "< golang import path to generated code >",
			"port": "8080",
			"impl": {
				"name": "dbfront_impl",
				"import": "< golang import path to your implementation >",
				"interface_factory": "GetServiceInterfaceImplementation",
				"callback_factory": "GetCallback"
			}
		}
	}`
	p.TemplateFileName = cmd.Flag("template", templateFileHelp).Required().ExistingFile()
	p.OutputDir = cmd.Flag("output-dir", "Directory name to put generated files in").Default(".").String()
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *testRigCmd) Execute(args cmdutils.ExecuteArgs) error {
	err := testrig.GenerateRig(args.Filesystem, *p.TemplateFileName, *p.OutputDir, args.Modules)
	if err != nil {
		return err
	}
	return nil
}
