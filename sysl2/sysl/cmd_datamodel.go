package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"regexp"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	log "github.com/sirupsen/logrus"
)

const localPlantuml = "http://localhost:8080/plantuml"

func GenerateDataModels(datagenParams *CmdContextParamDatagen) (map[string]string, error) {
	outmap := make(map[string]string)

	log.Debugf("project: %v\n", *datagenParams.project)
	log.Debugf("title: %s\n", *datagenParams.title)
	log.Debugf("filter: %s\n", *datagenParams.filter)
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
	mod := datagenParams.model

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


type dmCmd struct {
	title       string
	output      string
	project     string
	filter      string
	plantuml    string
	classFormat string
}

func (p *dmCmd) Name() string { return "data" }
func (p *dmCmd) RequireSyslModule() bool { return true }

func (p *dmCmd) Init(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Generate data models")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.classFormat)

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)

	cmd.Flag("plantuml",
		"base url of plantuml server (default: $SYSL_PLANTUML or "+
			"http://localhost:8080/plantuml see "+
			"http://plantuml.com/server.html#install for more info)",
	).Short('p').StringVar(&p.plantuml)

	cmd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)
	cmd.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').StringVar(&p.filter)

	return cmd
}

func (p *dmCmd) Execute(args ExecuteArgs) error {

	datagenParams := &CmdContextParamDatagen{
		model:        args.module,
		modelAppName: args.modAppName,
		title:        &p.title,
		output:       &p.output,
		project:      &p.project,
		filter:       &p.filter,
		plantuml:     &p.plantuml,
		classFormat:  &p.classFormat,
	}
	outmap, err := GenerateDataModels(datagenParams)
	if err != nil {
		return err
	}
	for k, v := range outmap {
		err := OutputPlantuml(k, *datagenParams.plantuml, v, args.fs)
		if err != nil {
			return fmt.Errorf("plantuml drawing error: %v", err)
		}
	}
	return nil
}
