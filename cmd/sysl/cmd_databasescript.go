package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func GenerateDatabaseScripts(scriptParams *CmdDatabaseScriptParams, model *sysl.Module,
	logger *logrus.Logger) ([]database.ScriptOutput, error) {
	logger.Debugf("Application names: %v\n", scriptParams.appNames)
	logger.Debugf("title: %s\n", scriptParams.title)
	logger.Debugf("outputDir: %s\n", scriptParams.outputDir)
	logger.Debugf("db type: %s\n", scriptParams.dbType)
	appNamesStr := strings.TrimSpace(scriptParams.appNames)
	if appNamesStr == "" {
		logger.Error("no application name specified")
		return nil, fmt.Errorf("no application names specified")
	}
	appNames := strings.Split(appNamesStr, database.Delimiter)
	outputSlice := processSysl(model, appNames, scriptParams.outputDir, scriptParams.title,
		scriptParams.dbType, logger)
	return outputSlice, nil
}

func processSysl(mod *sysl.Module,
	appNames []string, outputDir, title, dbType string, logger *logrus.Logger) []database.ScriptOutput {
	var outputSlice []database.ScriptOutput
	apps := mod.GetApps()
	for _, appName := range appNames {
		app := apps[appName]
		if app != nil {
			v := database.MakeDatabaseScriptView(title, logger)
			outStr := v.GenerateDatabaseScriptCreate(app.GetTypes(), dbType, appName)
			outputFile := filepath.Join(outputDir, appName+database.SQLExtension)
			outputStruct := database.MakeScriptOutput(outputFile, outStr)
			outputSlice = append(outputSlice, *outputStruct)
		}
	}
	return outputSlice
}

type databaseScriptCmd struct {
	CmdDatabaseScriptParams
}

func (p *databaseScriptCmd) Name() string       { return "generate-db-scripts" }
func (p *databaseScriptCmd) MaxSyslModule() int { return 1 }

func (p *databaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate db script").Alias("generatedbscripts")

	cmd.Flag("title", "file title").Short('t').StringVar(&p.title)
	cmd.Flag("output-dir", "output directory for generated file").Short('o').StringVar(&p.outputDir)
	cmd.Flag("app-names", "application names to parse").Short('a').StringVar(&p.appNames)
	cmd.Flag("db-type", "database type e.g postgres").Short('d').StringVar(&p.dbType)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *databaseScriptCmd) Execute(args ExecuteArgs) error {
	outputSlice, err := GenerateDatabaseScripts(&p.CmdDatabaseScriptParams, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return database.GenerateFromSQLMap(outputSlice, args.Filesystem, args.Logger)
}
