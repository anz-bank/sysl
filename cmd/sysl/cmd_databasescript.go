package main

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func GenerateDatabaseScripts(scriptParams *CmdDatabaseScript,
	logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("project: %v\n", scriptParams.project)
	logger.Debugf("title: %s\n", scriptParams.title)
	logger.Debugf("output: %s\n", scriptParams.output)
	logger.Debugf("root: %s\n", scriptParams.root)
	logger.Debugf("source: %s\n", scriptParams.source)
	model, _, err := LoadSyslModule(scriptParams.root, scriptParams.source, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	// The "project" app that specifies the data models to be built
	var app *sysl.Application
	var exists bool
	if app, exists = model.GetApps()[scriptParams.project]; !exists {
		return nil, fmt.Errorf("project not found in sysl")
	}

	epname := "Relational-Model"
	endpt := app.GetEndpoints()[epname]
	outputDir := scriptParams.output
	if strings.Contains(outputDir, "%(epname)") {
		of := MakeFormatParser(scriptParams.output)
		outputDir = of.FmtOutput(scriptParams.project, epname, "", nil)
	}
	generateDatabaseScripts(outmap, model, endpt.GetStmt(), scriptParams.title, scriptParams.project, outputDir)
	return outmap, nil
}

func generateDatabaseScripts(outmap map[string]string, mod *sysl.Module,
	stmts []*sysl.Statement, title, project, outDir string) {
	apps := mod.GetApps()

	// Parse all the applications in the project
	outStr := ""
	if title != "" {
		outStr += "/*TITLE : " + title + "*/\n"
	}
	outStr += databaseScriptHeader
	for _, stmt := range stmts {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			var stringBuilder strings.Builder
			app := apps[a.Action.Action]
			if app != nil {
				dataParam := &DatabaseScriptParam{
					types:   app.GetTypes(),
					title:   title,
					project: project,
				}
				v := MakeDatabaseScriptView(&stringBuilder, dataParam.title, dataParam.project)
				outStr += "\n\n/*-----------------------Relation Model : " +
					a.Action.Action + "-----------------------------------------------*/\n"
				outStr += v.GenerateDatabaseScriptCreate(dataParam)
			}
		}
	}
	outmap[outDir] = outStr
}

type databaseScriptCmd struct {
	CmdDatabaseScript
}

func (p *databaseScriptCmd) Name() string            { return "generate-script" }
func (p *databaseScriptCmd) RequireSyslModule() bool { return false }

func (p *databaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate postgres sql script").Alias("generatescript")

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)
	cmd.Flag("rootDir", "root dir").Short('r').StringVar(&p.root)
	cmd.Flag("source", "source sysl").Short('s').StringVar(&p.source)
	cmd.Flag("output",
		"output file (default: %(epname).sql)",
	).Default("%(epname).sql").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *databaseScriptCmd) Execute(args ExecuteArgs) error {
	outmap, err := GenerateDatabaseScripts(&p.CmdDatabaseScript, args.Logger)
	if err != nil {
		return err
	}
	return GenerateFromSQLMap(outmap, args.Filesystem)
}
