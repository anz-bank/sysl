package main

import (
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func GenerateDataModelsView(datagenParams *CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("title: %s\n", datagenParams.title)
	logger.Debugf("filter: %s\n", datagenParams.filter)
	logger.Debugf("output: %s\n", datagenParams.output)

	spclass := constructFormatParser("", datagenParams.classFormat)

	var app *sysl.Application
	outputDir := datagenParams.output
	if strings.Contains(outputDir, "%(epname)") {
		of := MakeFormatParser(datagenParams.output)
		outputDir = of.FmtOutput(datagenParams.project, "Default", app.GetLongName(), app.GetAttrs())
	}

	generateDataModelView(spclass, outmap, model, datagenParams.title, datagenParams.project, outputDir)

	return outmap, nil
}

func generateDataModelView(pclass ClassLabeler, outmap map[string]string, mod *sysl.Module, title, project,
	outDir string) {
	apps := mod.GetApps()
	for key := range apps {
		var stringBuilder strings.Builder
		app := apps[key]
		if app != nil {
			dataParam := &DataModelParam{
				mod:   mod,
				app:   app,
				title: title,
			}
			v := MakeDataModelView(pclass, dataParam.mod, &stringBuilder, dataParam.title, "")
			outmap[outDir] = v.GenerateDataView(dataParam)
		}
	}
}

// Process pure Sysl datamodel file produced by importer cmd
type datamodelViewCmd struct {
	plantumlmixin
	CmdContextParamDatagen
}

func (p *datamodelViewCmd) Name() string       { return "datamodelview" }
func (p *datamodelViewCmd) MaxSyslModule() int { return 1 }

func (p *datamodelViewCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models view").Alias("dataview")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.classFormat)

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)

	p.AddFlag(cmd)

	cmd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').StringVar(&p.filter)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *datamodelViewCmd) Execute(args ExecuteArgs) error {
	outmap, err := GenerateDataModelsView(&p.CmdContextParamDatagen, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(outmap, args.Filesystem)
}
