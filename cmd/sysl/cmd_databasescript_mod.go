package main

import (
	"fmt"
	"strings"

	db "github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func GenerateModDatabaseScripts(scriptParams *CmdDatabaseScript, modelOld, modelNew *sysl.Module,
	logger *logrus.Logger) ([]db.ScriptOutput, error) {
	logger.Debugf("Application names: %v\n", scriptParams.appNames)
	logger.Debugf("title: %s\n", scriptParams.title)
	logger.Debugf("outputDir: %s\n", scriptParams.outputDir)
	logger.Debugf("db type: %s\n", scriptParams.dbType)
	appNamesStr := strings.TrimSpace(scriptParams.appNames)
	if appNamesStr == "" {
		return nil, fmt.Errorf("no application names specified")
	}
	appNames := strings.Split(appNamesStr, db.Delimiter)
	v := db.MakeDatabaseScriptView(scriptParams.title)
	outputSlice := v.ProcessModSysls(modelOld.GetApps(), modelNew.GetApps(), appNames,
		scriptParams.outputDir, scriptParams.dbType)
	return outputSlice, nil
}

type modDatabaseScriptCmd struct {
	CmdDatabaseScript
}

func (p *modDatabaseScriptCmd) Name() string       { return "generate-db-scripts-delta" }
func (p *modDatabaseScriptCmd) MaxSyslModule() int { return 2 }

func (p *modDatabaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("generatedbscriptsdelta")

	cmd.Flag("title", "file title").Short('t').StringVar(&p.title)
	cmd.Flag("output-dir", "output directory").Short('o').StringVar(&p.outputDir)
	cmd.Flag("app-names", "application names to read").Short('a').StringVar(&p.appNames)
	cmd.Flag("db-type", "database type e.g postgres").Short('d').StringVar(&p.dbType)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *modDatabaseScriptCmd) Execute(args ExecuteArgs) error {
	if len(args.Modules) < 2 {
		return fmt.Errorf("this command needs min 2 module(s)")
	}
	outputSlice, err := GenerateModDatabaseScripts(&p.CmdDatabaseScript, args.Modules[0], args.Modules[1], args.Logger)
	if err != nil {
		return err
	}
	return db.GenerateFromSQLMap(outputSlice, args.Filesystem)
}
