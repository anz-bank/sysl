package main

import (
	"fmt"
	"io/ioutil"

	"github.com/anz-bank/sysl/sysl2/sysl/exporter"
	"gopkg.in/alecthomas/kingpin.v2"
)

type exportCmd struct {
	appName string
	out     string
	target  string
	format  string
}

func (p *exportCmd) Name() string            { return "export" }
func (p *exportCmd) RequireSyslModule() bool { return true }

func (p *exportCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Export sysl to external types. Supported types: Swagger")
	cmd.Flag("app-name",
		"name of the sysl app defined in sysl model."+
			" if there are multiple apps defined in sysl model,"+
			" swagger will be generated only for the given app").Required().Short('a').StringVar(&p.appName)
	cmd.Flag("output", "output filename").Default("output.yaml").Short('o').StringVar(&p.out)
	cmd.Flag("target", "export target").Default("swagger").Short('t').StringVar(&p.target)
	cmd.Flag("format", "export format").Default("yaml").Short('f').StringVar(&p.format)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *exportCmd) Execute(args ExecuteArgs) error {
	syslApp, syslAppFound := args.Module.GetApps()[p.appName]
	var output []byte
	var err error
	if !syslAppFound {
		return fmt.Errorf("app not found in the Sysl file")
	}
	if p.target == "swagger" {
		swaggerExporter := exporter.MakeSwaggerExporter(syslApp, args.Logger)
		err = swaggerExporter.GenerateSwagger()
		if err != nil {
			return err
		}
		output, err = swaggerExporter.SerializeOutput(p.format)
		if err != nil {
			return err
		}
	}

	fileWriteError := ioutil.WriteFile(p.out, output, 0644)
	if fileWriteError != nil {
		return fileWriteError
	}
	return nil
}
