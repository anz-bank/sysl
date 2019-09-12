package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

const localPlantuml = "http://localhost:8080/plantuml"

func GenerateDataModels(datagenParams *CmdContextParamDatagen) (map[string]string, error) {
	outmap := make(map[string]string)

	log.Debugf("root: %s\n", *datagenParams.root)
	log.Debugf("project: %v\n", *datagenParams.project)
	log.Debugf("title: %s\n", *datagenParams.title)
	log.Debugf("filter: %s\n", *datagenParams.filter)
	log.Debugf("modules: %s\n", *datagenParams.modules)
	log.Debugf("output: %s\n", *datagenParams.output)

	if *datagenParams.plantuml == "" {
		plantuml := os.Getenv("SYSL_PLANTUML")
		datagenParams.plantuml = &plantuml
		if *datagenParams.plantuml == "" {
			*datagenParams.plantuml = localPlantuml
		}
	}
	log.Debugf("plantuml: %s\n", *datagenParams.plantuml)

	spclass := constructFormatParser("", *datagenParams.classFormat)
	mod, err := loadApp(*datagenParams.modules, syslutil.NewChrootFs(afero.NewOsFs(), *datagenParams.root))

	if err != nil {
		return nil, err
	}

	// The "project" app that specifies the data models to be built
	var app *sysl.Application
	var exists bool
	if app, exists = mod.GetApps()[*datagenParams.project]; !exists {
		return nil, fmt.Errorf("project not found in sysl")
	}

	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := *datagenParams.output
		if strings.Contains(outputDir, "%(epname)") {
			of := MakeFormatParser(*datagenParams.output)
			outputDir = of.FmtOutput(*datagenParams.project, epname, endpt.GetLongName(), endpt.GetAttrs())
		}
		if *datagenParams.filter != "" {
			re := regexp.MustCompile(*datagenParams.filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		generateDataModel(spclass, outmap, mod, endpt.GetStmt(), *datagenParams.title, *datagenParams.project, outputDir)
	}
	return outmap, nil
}

func configureCmdlineForDatagen(sysl *kingpin.Application) *CmdContextParamDatagen {
	data := sysl.Command("data", "Generate data models")
	returnValues := &CmdContextParamDatagen{}

	returnValues.root = data.Flag("root", "sysl root directory for input model file (default: .)").Default(".").String()
	returnValues.classFormat = data.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").String()
	returnValues.title = data.Flag("title", "diagram title").Short('t').String()
	returnValues.plantuml = data.Flag("plantuml", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n")).Short('p').String()
	returnValues.output = data.Flag("output",
		"output file(default: %(epname).png)").Default("%(epname).png").Short('o').String()
	returnValues.project = data.Flag("project", "project pseudo-app to render").Short('j').String()
	returnValues.filter = data.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').String()
	returnValues.modules = data.Arg("modules",
		strings.Join([]string{"input files without .sysl extension and with leading /",
			"eg: /project_dir/my_models",
			"combine with --root if needed"}, "\n")).String()

	return returnValues
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
