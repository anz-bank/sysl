package main

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	sysl "github.com/anz-bank/sysl/src/proto_old"
	"github.com/sirupsen/logrus"
)

func GenerateDataModels(datagenParams *CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("project: %v\n", datagenParams.project)
	logger.Debugf("title: %s\n", datagenParams.title)
	logger.Debugf("filter: %s\n", datagenParams.filter)
	logger.Debugf("output: %s\n", datagenParams.output)

	spclass := constructFormatParser("", datagenParams.classFormat)

	// The "project" app that specifies the data models to be built
	var app *sysl.Application
	var exists bool
	if app, exists = model.GetApps()[datagenParams.project]; !exists {
		return nil, fmt.Errorf("project not found in sysl")
	}

	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := datagenParams.output
		if strings.Contains(outputDir, "%(epname)") {
			of := MakeFormatParser(datagenParams.output)
			outputDir = of.FmtOutput(datagenParams.project, epname, endpt.GetLongName(), endpt.GetAttrs())
		}
		if datagenParams.filter != "" {
			re := regexp.MustCompile(datagenParams.filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		generateDataModel(spclass, outmap, model, endpt.GetStmt(), datagenParams.title, datagenParams.project, outputDir)
	}
	return outmap, nil
}

func generateDataModel(pclass ClassLabeler, outmap map[string]string, mod *sysl.Module,
	stmts []*sysl.Statement, title, project, outDir string) {
	apps := mod.GetApps()

	// Parse all the applications in the project
	for _, stmt := range stmts {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			var stringBuilder strings.Builder
			app := apps[a.Action.Action]
			if app != nil {
				dataParam := &DataModelParam{
					mod:     mod,
					app:     app,
					title:   title,
					project: project,
				}
				v := MakeDataModelView(pclass, dataParam.mod, &stringBuilder, dataParam.title, dataParam.project)
				outmap[outDir] = v.GenerateDataView(dataParam)
			}
		}
	}
}

type datamodelCmd struct {
	plantumlmixin
	CmdContextParamDatagen
}

func (p *datamodelCmd) Name() string            { return "datamodel" }
func (p *datamodelCmd) RequireSyslModule() bool { return true }

func (p *datamodelCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("data")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.classFormat)

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)

	p.AddFlag(cmd)

	cmd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)
	cmd.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').StringVar(&p.filter)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *datamodelCmd) Execute(args ExecuteArgs) error {
	outmap, err := GenerateDataModels(&p.CmdContextParamDatagen, args.Module, args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(outmap, args.Filesystem)
}
