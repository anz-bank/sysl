package datamodeldiagram

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func GenerateDataModelsWithProjectMannerModule(datagenParams *cmdutils.CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("project: %v\n", datagenParams.Project)
	logger.Debugf("title: %s\n", datagenParams.Title)
	logger.Debugf("filter: %s\n", datagenParams.Filter)
	logger.Debugf("output: %s\n", datagenParams.Output)

	spclass := sequencediagram.ConstructFormatParser("", datagenParams.ClassFormat)

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
		// TODO: Fix the need to pass `epNameMode` into function
		epNameMode := strings.Contains(datagenParams.Output, "%(epname)")
		GenerateDataModel(spclass, outmap, model, endpt.GetStmt(), datagenParams.Title,
			datagenParams.Project, outputDir, epNameMode)
	}
	return outmap, nil
}

/**
 * It is added to help reviewing generated data model with sysl
 * file produced by command import. Generate data model diagrams using the following command:
 * sysl data      -d --root=/Users/guest/data -t Test -o Test.png Test
 * sysl datamodel -d --root=/Users/guest/data -t Test -o Test.png Test.sysl
 */
func GenerateDataModelsWithPureModule(datagenParams *cmdutils.CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("title: %s\n", datagenParams.Title)
	logger.Debugf("output: %s\n", datagenParams.Output)

	spclass := sequencediagram.ConstructFormatParser("", datagenParams.ClassFormat)

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
				Mod:    model,
				App:    app,
				Title:  datagenParams.Title,
				Epname: strings.Contains(datagenParams.Output, "%(epname)"),
			}
			v := MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
			outmap[outputDir] = v.GenerateDataView(dataParam)
		}
	}

	return outmap, nil
}

func GenerateDataModel(pclass cmdutils.ClassLabeler, outmap map[string]string, mod *sysl.Module,
	stmts []*sysl.Statement, title, project, outDir string, epName bool) {
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
					Epname:  epName,
				}
				v := MakeDataModelView(pclass, dataParam.Mod, &stringBuilder, dataParam.Title, dataParam.Project)
				outmap[outDir] = v.GenerateDataView(dataParam)
			}
		}
	}
}

func GenerateDataModels(datagenParams *cmdutils.CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	if datagenParams.Direct {
		// The sysl file is not project manner
		return GenerateDataModelsWithPureModule(datagenParams, model, logger)
	}
	// Sysl file is project manner
	return GenerateDataModelsWithProjectMannerModule(datagenParams, model, logger)
}
