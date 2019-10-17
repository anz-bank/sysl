package main

import (
	"fmt"
	"io/ioutil"

	"github.com/anz-bank/sysl/sysl2/sysl/exporter"
	"gopkg.in/alecthomas/kingpin.v2"
)

type swaggerExportCmd struct {
	appName    string
	outSwagger string
}

func (p *swaggerExportCmd) Name() string            { return "export-swagger" }
func (p *swaggerExportCmd) RequireSyslModule() bool { return true }

func (p *swaggerExportCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Generate swagger")
	cmd.Flag("app-name",
		"name of the sysl app defined in sysl model."+
			" if there are multiple apps defined in sysl model,"+
			" swagger will be generated only for the given app").Required().Short('a').StringVar(&p.appName)
	cmd.Flag("output", "output filename").Default("output.yaml").Short('o').StringVar(&p.outSwagger)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *swaggerExportCmd) Execute(args ExecuteArgs) error {
	syslApp, syslAppFound := args.Module.GetApps()[p.appName]
	if !syslAppFound {
		return fmt.Errorf("app not found in the Sysl file")
	}
	swaggerExporter := exporter.MakeSwaggerExporter(syslApp, args.Logger)
	err := swaggerExporter.GenerateSwagger()
	if err != nil {
		return err
	}
	outYaml, err := swaggerExporter.SerializeToYaml()
	if err != nil {
		return err
	}
	fileWriteError := ioutil.WriteFile(p.outSwagger, outYaml, 0644)
	if fileWriteError != nil {
		return fileWriteError
	}
	return nil
}
