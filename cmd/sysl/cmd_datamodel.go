package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"gopkg.in/alecthomas/kingpin.v2"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func generateDataModelsWithProjectMannerModule(datagenParams *cmdutils.CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("project: %v\n", datagenParams.Project)
	logger.Debugf("title: %s\n", datagenParams.Title)
	logger.Debugf("filter: %s\n", datagenParams.Filter)
	logger.Debugf("output: %s\n", datagenParams.Output)

	spclass := ConstructFormatParser("", datagenParams.ClassFormat)

	// The "project" app that specifies the data models to be built
	var app *sysl.Application
	var exists bool
	if app, exists = model.GetApps()[datagenParams.Project]; !exists {
		return nil, fmt.Errorf("project not found in sysl")
	}

	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := datagenParams.Output
		if strings.Contains(outputDir, "%(epname)") {
			of := cmdutils.MakeFormatParser(datagenParams.Output)
			outputDir = of.FmtOutput(datagenParams.Project, epname, endpt.GetLongName(), endpt.GetAttrs())
		}
		if datagenParams.Filter != "" {
			re := regexp.MustCompile(datagenParams.Filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		generateDataModel(spclass, outmap, model, endpt.GetStmt(), datagenParams.Title, datagenParams.Project, outputDir)
	}
	return outmap, nil
}

/**
 * It is added to help reviewing generated data model with sysl
 * file produced by command import. Generate data model diagrams using the following command:
 * sysl data      -d --root=/Users/guest/data -t Test -o Test.png Test
 * sysl datamodel -d --root=/Users/guest/data -t Test -o Test.png Test.sysl
 */
func generateDataModelsWithPureModule(datagenParams *cmdutils.CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("title: %s\n", datagenParams.Title)
	logger.Debugf("output: %s\n", datagenParams.Output)

	spclass := ConstructFormatParser("", datagenParams.ClassFormat)

	apps := model.GetApps()
	for appName := range apps {
		app := apps[appName]
		outputDir := datagenParams.Output
		if strings.Contains(outputDir, "%(epname)") {
			of := cmdutils.MakeFormatParser(datagenParams.Output)
			outputDir = of.FmtOutput(appName, appName, app.GetLongName(), app.GetAttrs())
		}
		var stringBuilder strings.Builder
		if app != nil {
			dataParam := &DataModelParam{
				Mod:   model,
				App:   app,
				Title: datagenParams.Title,
			}
			v := MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
			outmap[outputDir] = v.GenerateDataView(dataParam)
		}
	}

	return outmap, nil
}

func generateDataModel(pclass cmdutils.ClassLabeler, outmap map[string]string, mod *sysl.Module,
	stmts []*sysl.Statement, title, project, outDir string) {
	apps := mod.GetApps()

	// Parse all the applications in the project
	for _, stmt := range stmts {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			var stringBuilder strings.Builder
			app := apps[a.Action.Action]
			if app != nil {
				dataParam := &DataModelParam{
					Mod:     mod,
					App:     app,
					Title:   title,
					Project: project,
				}
				v := MakeDataModelView(pclass, dataParam.Mod, &stringBuilder, dataParam.Title, dataParam.Project)
				outmap[outDir] = v.GenerateDataView(dataParam)
			}
		}
	}
}

func generateDataModels(datagenParams *cmdutils.CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	if datagenParams.Direct {
		// The sysl file is not project manner
		return generateDataModelsWithPureModule(datagenParams, model, logger)
	}
	// Sysl file is project manner
	return generateDataModelsWithProjectMannerModule(datagenParams, model, logger)
}

type datamodelCmd struct {
	plantumlmixin
	cmdutils.CmdContextParamDatagen
}

func (p *datamodelCmd) Name() string       { return "datamodel" }
func (p *datamodelCmd) MaxSyslModule() int { return 1 }

func (p *datamodelCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("data")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.ClassFormat)

	cmd.Flag("title", "Diagram title").Short('t').StringVar(&p.Title)

	p.AddFlag(cmd)

	cmd.Flag("output",
		"Output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.Output)
	cmd.Flag("project", "Project pseudo-app to render").Short('j').StringVar(&p.Project)
	cmd.Flag("direct", "Process data model directly without project manner").Short('d').BoolVar(&p.Direct)
	cmd.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').StringVar(&p.Filter)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *datamodelCmd) Execute(args cmdutils.ExecuteArgs) error {
	outmap, err := generateDataModels(&p.CmdContextParamDatagen, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(outmap, args.Filesystem)
}
