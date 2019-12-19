package main

import (
	"fmt"
	"path/filepath"
	"strings"

	db "github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

func GenerateDatabaseScripts(scriptParams *CmdDatabaseScript,
	logger *logrus.Logger) ([]db.ScriptOutput, error) {
	logger.Debugf("Application names: %v\n", scriptParams.appNames)
	logger.Debugf("title: %s\n", scriptParams.title)
	logger.Debugf("outputDir: %s\n", scriptParams.outputDir)
	logger.Debugf("inputDir: %s\n", scriptParams.inputDir)
	logger.Debugf("source: %s\n", scriptParams.source)
	logger.Debugf("db type: %s\n", scriptParams.dbType)
	appNamesStr := strings.TrimSpace(scriptParams.appNames)
	if appNamesStr == "" {
		return nil, fmt.Errorf("No application names specified")
	}
	model, _, err := LoadSyslModule(scriptParams.inputDir, scriptParams.source, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	appNames := strings.Split(appNamesStr, db.Delimiter)
	outputSlice := processSysl(model, appNames, scriptParams.outputDir, scriptParams.title, scriptParams.dbType)
	return outputSlice, nil
}
func processSysl(mod *sysl.Module,
	appNames []string, outputDir, title, dbType string) []db.ScriptOutput {
	outputSlice := []db.ScriptOutput{}
	apps := mod.GetApps()
	for _, appName := range appNames {
		app := apps[appName]
		if app != nil {
			v := db.MakeDatabaseScriptView(title, appName)
			outStr := v.GenerateDatabaseScriptCreate(app.GetTypes(), dbType)
			outputFile := filepath.Join(outputDir, appName+db.SQLExtension)
			outputStruct := db.MakeScriptOutput(outputFile, outStr)
			outputSlice = append(outputSlice, *outputStruct)
		}
	}
	return outputSlice
}

type databaseScriptCmd struct {
	CmdDatabaseScript
}

func (p *databaseScriptCmd) Name() string            { return "generate-db-scripts" }
func (p *databaseScriptCmd) RequireSyslModule() bool { return false }

func (p *databaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate postgres sql script").Alias("generatedbscripts")

	cmd.Flag("title", "file title").Short('t').StringVar(&p.title)
	cmd.Flag("input-dir", "input dir").Short('i').StringVar(&p.inputDir)
	cmd.Flag("source", "source sysl").Short('s').StringVar(&p.source)
	cmd.Flag("output-dir", "output directory for generated file").Short('o').StringVar(&p.outputDir)
	cmd.Flag("app-names", "application names to parse").Short('a').StringVar(&p.appNames)
	cmd.Flag("db-type", "database type e.g postgres").Short('d').StringVar(&p.dbType)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *databaseScriptCmd) Execute(args ExecuteArgs) error {
	outputSlice, err := GenerateDatabaseScripts(&p.CmdDatabaseScript, args.Logger)
	if err != nil {
		return err
	}
	return db.GenerateFromSQLMap(outputSlice, args.Filesystem)
}
