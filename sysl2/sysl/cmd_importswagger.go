package main

import (
	"io/ioutil"

	"github.com/anz-bank/sysl/sysl2/sysl/swagger"
	"gopkg.in/alecthomas/kingpin.v2"
)

type importSwaggerCmd struct {
	swagger.OutputData
	filename string
	outfile  string
}

func (p *importSwaggerCmd) Name() string            { return "import-swagger" }
func (p *importSwaggerCmd) RequireSyslModule() bool { return false }

func (p *importSwaggerCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Convert swagger yaml/json -> sysl")
	cmd.Flag("input", "swagger input filename").Required().StringVar(&p.filename)
	cmd.Flag("app-name",
		"name of the sysl app to define in sysl model.").Required().StringVar(&p.AppName)
	cmd.Flag("package",
		"name of the sysl package to define in sysl model.").Required().StringVar(&p.Package)
	cmd.Flag("output", "output filename").Default("output.sysl").StringVar(&p.outfile)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *importSwaggerCmd) Execute(args ExecuteArgs) error {

	data, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}

	output, err := swagger.LoadSwaggerText(p.OutputData, string(data), args.Logger)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.outfile, []byte(output), 0644)
}
