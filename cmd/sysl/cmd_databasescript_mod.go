package main

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func GenerateModDatabaseScripts(scriptParams *cmdutils.CmdDatabaseScriptParams, modelOld, modelNew *sysl.Module,
	logger *logrus.Logger) ([]database.ScriptOutput, error) {
	logger.Debugf("Application names: %v\n", scriptParams.AppNames)
	logger.Debugf("title: %s\n", scriptParams.Title)
	logger.Debugf("outputDir: %s\n", scriptParams.OutputDir)
	logger.Debugf("db type: %s\n", scriptParams.DBType)
	appNamesStr := strings.TrimSpace(scriptParams.AppNames)
	if appNamesStr == "" {
		logger.Error("no application name specified")
		return nil, fmt.Errorf("no application names specified")
	}
	appNames := strings.Split(appNamesStr, database.Delimiter)
	v := database.MakeDatabaseScriptView(scriptParams.Title, logger)
	outputSlice := v.ProcessModSysls(modelOld.GetApps(), modelNew.GetApps(), appNames,
		scriptParams.OutputDir, scriptParams.DBType)
	return outputSlice, nil
}

type modDatabaseScriptCmd struct {
	cmdutils.CmdDatabaseScriptParams
}

func (p *modDatabaseScriptCmd) Name() string       { return "generate-db-scripts-delta" }
func (p *modDatabaseScriptCmd) MaxSyslModule() int { return 2 }

func (p *modDatabaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate delta db scripts").Alias("generatedbscriptsdelta")

	cmd.Flag("title", "file title").Short('t').StringVar(&p.Title)
	cmd.Flag("output-dir", "output directory").Short('o').StringVar(&p.OutputDir)
	cmd.Flag("app-names", "application names to read").Short('a').StringVar(&p.AppNames)
	cmd.Flag("db-type", "database type e.g postgres").Short('d').StringVar(&p.DBType)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *modDatabaseScriptCmd) Execute(args cmdutils.ExecuteArgs) error {
	if len(args.Modules) < 2 {
		return fmt.Errorf("this command needs min 2 module(s)")
	}
	outputSlice, err := GenerateModDatabaseScripts(&p.CmdDatabaseScriptParams,
		args.Modules[0], args.Modules[1], args.Logger)
	if err != nil {
		return err
	}
	return database.GenerateFromSQLMap(outputSlice, args.Filesystem, args.Logger)
}
