package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	db "github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

func GenerateModDatabaseScripts(scriptParams *CmdDatabaseScriptMod,
	logger *logrus.Logger) ([]db.ScriptOutput, error) {

	logger.Debugf("Application names: %v\n", scriptParams.appNames)
	logger.Debugf("title: %s\n", scriptParams.title)
	logger.Debugf("outputDir: %s\n", scriptParams.outputDir)
	logger.Debugf("inputDir: %s\n", scriptParams.inputDir)
	logger.Debugf("source org: %s\n", scriptParams.orgSource)
	logger.Debugf("source new: %s\n", scriptParams.newSource)
	logger.Debugf("db type: %s\n", scriptParams.dbType)
	appNamesStr := strings.TrimSpace(scriptParams.appNames)
	if appNamesStr == "" {
		return nil, fmt.Errorf("No application names specified")
	}
	modelOld, _, err1 := LoadSyslModule(scriptParams.inputDir, scriptParams.orgSource, afero.NewOsFs(), logger)
	if err1 != nil {
		return nil, err1
	}
	modelNew, _, err2 := LoadSyslModule(scriptParams.inputDir, scriptParams.newSource, afero.NewOsFs(), logger)
	if err2 != nil {
		return nil, err2
	}
	appNames := strings.Split(appNamesStr, db.Delimiter)
	outputSlice := processModSysl(modelOld.GetApps(), modelNew.GetApps(), appNames,
		scriptParams.title, scriptParams.outputDir, scriptParams.dbType)
	return outputSlice, nil
}

func processModSysl(appsOld, appsNew map[string]*sysl.Application,
	appNames []string, title, outputDir, dbType string) []db.ScriptOutput {
	outputSlice := []db.ScriptOutput{}
	for _, appName := range appNames {
		appOld := appsOld[appName]
		appNew := appsNew[appName]
		if appOld != nil && appNew != nil {
			typeMapOld := appOld.GetTypes()
			typeMapNew := appNew.GetTypes()
			tableDepthMapOld := db.CreateTableDepthMap(typeMapOld)
			tableDepthMapNew := db.CreateTableDepthMap(typeMapNew)
			tablesWithActions := findAddedDeletedRetainedTables(typeMapOld, typeMapNew, tableDepthMapOld, tableDepthMapNew)
			outStr := processTablesForModifiedApps(tablesWithActions, title, appName, dbType)
			outputFile := filepath.Join(outputDir, appName+db.SQLExtension)
			outputStruct := db.MakeScriptOutput(outputFile, outStr)
			outputSlice = append(outputSlice, *outputStruct)
		}
	}
	for _, appName := range appNames {
		appNew := appsNew[appName]
		appOld := appsOld[appName]
		if appNew != nil && appOld == nil {
			v := db.MakeDatabaseScriptView(title, appName)
			outStr := v.GenerateDatabaseScriptCreate(appNew.GetTypes(), dbType)
			outputFile := filepath.Join(outputDir, appName+db.SQLExtension)
			outputStruct := db.MakeScriptOutput(outputFile, outStr)
			outputSlice = append(outputSlice, *outputStruct)
		}
	}
	return outputSlice
}

func processTablesForModifiedApps(tableDetails []db.TableDetails, title, appName, dbType string) string {
	v := db.MakeDatabaseScriptView(title, appName)
	outputStr := v.GenerateDatabaseScriptModify(tableDetails, dbType)
	return outputStr
}

func findAddedDeletedRetainedTables(
	tableMapOld map[string]*sysl.Type,
	tableMapNew map[string]*sysl.Type,
	tableDepthMapOld, tableDepthMapNew map[int][]string) []db.TableDetails {
	tableDepthsListNew := makeSortedListOfTableDepth(tableDepthMapNew)
	tableWithAction := []db.TableDetails{}
	//Add the retained and added tables
	for _, depth := range tableDepthsListNew {
		tableNames := tableDepthMapNew[depth]
		sort.Strings(tableNames)
		for _, tableName := range tableNames {
			found := ifTableExisted(tableName, tableMapOld)
			if found {
				tableDetails := db.MakeTableDetails(tableMapNew[tableName],
					tableMapOld[tableName], db.Retain, tableName)
				tableWithAction = append(tableWithAction, *tableDetails)
			} else {
				tableDetails := db.MakeTableDetails(tableMapNew[tableName],
					nil, db.Add, tableName)
				tableWithAction = append(tableWithAction, *tableDetails)
			}
		}
	}

	return tableWithAction
}

func ifTableExisted(tableName string, tableMap map[string]*sysl.Type) bool {
	item := tableMap[tableName]
	return item != nil
}

func makeSortedListOfTableDepth(tableDepthMap map[int][]string) []int {
	depthList := []int{}
	for depth := range tableDepthMap {
		depthList = append(depthList, depth)
	}
	sort.Ints(depthList)
	return depthList
}

type modDatabaseScriptCmd struct {
	CmdDatabaseScriptMod
}

func (p *modDatabaseScriptCmd) Name() string            { return "generate-db-scripts-delta" }
func (p *modDatabaseScriptCmd) RequireSyslModule() bool { return false }

func (p *modDatabaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("generatedbscriptsdelta")

	cmd.Flag("title", "file title").Short('t').StringVar(&p.title)
	cmd.Flag("input-dir", "input dir").Short('r').StringVar(&p.inputDir)
	cmd.Flag("org-source", "org source sysl").Short('s').StringVar(&p.orgSource)
	cmd.Flag("new-source", "new source sysl").Short('n').StringVar(&p.newSource)
	cmd.Flag("output-dir", "output directory").Short('o').StringVar(&p.outputDir)
	cmd.Flag("app-names", "application names to read").Short('a').StringVar(&p.appNames)
	cmd.Flag("db-type", "database type e.g postgres").Short('d').StringVar(&p.dbType)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *modDatabaseScriptCmd) Execute(args ExecuteArgs) error {
	outputSlice, err := GenerateModDatabaseScripts(&p.CmdDatabaseScriptMod, args.Logger)
	if err != nil {
		return err
	}
	return db.GenerateFromSQLMap(outputSlice, args.Filesystem)
}
