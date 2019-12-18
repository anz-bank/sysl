package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func GenerateModDatabaseScripts(scriptParams *CmdDatabaseScriptMod,
	logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("project: %v\n", scriptParams.project)
	logger.Debugf("title: %s\n", scriptParams.title)
	logger.Debugf("output: %s\n", scriptParams.output)
	logger.Debugf("source: %s\n", scriptParams.orgSource)
	logger.Debugf("output: %s\n", scriptParams.newSource)

	modelOld, _, err1 := LoadSyslModule(scriptParams.root, scriptParams.orgSource, afero.NewOsFs(), logger)
	if err1 != nil {
		return nil, err1
	}
	modelNew, _, err2 := LoadSyslModule(scriptParams.root, scriptParams.newSource, afero.NewOsFs(), logger)
	if err2 != nil {
		return nil, err2
	}
	// The "project" app that specifies the data models to be built
	var appOld, appNew *sysl.Application
	var existsOld, existsNew bool
	appOld, existsOld = modelOld.GetApps()[scriptParams.project]
	appNew, existsNew = modelNew.GetApps()[scriptParams.project]
	if !(existsOld || existsNew) {
		return nil, fmt.Errorf("project not found in sysl for either new or old script")
	}

	epname := "Relational-Model"
	endptOld := appOld.GetEndpoints()[epname]
	endptNew := appNew.GetEndpoints()[epname]
	outputDir := scriptParams.output
	if strings.Contains(outputDir, "%(epname)") {
		of := MakeFormatParser(scriptParams.output)
		outputDir = of.FmtOutput(scriptParams.project, epname, "", nil)
	}
	generateModDatabaseScripts(outmap, modelOld, modelNew, endptOld.GetStmt(), endptNew.GetStmt(),
		scriptParams.title, scriptParams.project, outputDir)
	return outmap, nil
}

func generateModDatabaseScripts(outmap map[string]string, modOld *sysl.Module, modNew *sysl.Module,
	stmtsOld []*sysl.Statement, stmtsNew []*sysl.Statement, title, project, outDir string) {
	appsOld := modOld.GetApps()
	appsNew := modNew.GetApps()
	var stringBuilder strings.Builder
	v := MakeEmptyDatabaseScriptView()
	outStr := ""
	if title != "" {
		outStr += "/*TITLE : " + title + "*/\n"
	}
	outStr += databaseScriptHeader
	// Parse all the applications in the project
	for _, stmt := range stmtsOld {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			appOld := appsOld[a.Action.Action]
			appNew := appsNew[a.Action.Action]
			// add/retained in the new sysl file
			if appOld != nil && appNew != nil {
				typeMapOld := appOld.GetTypes()
				typeMapNew := appNew.GetTypes()
				tableDepthMapOld := v.createTableDepthMap(typeMapOld)
				tableDepthMapNew := v.createTableDepthMap(typeMapNew)
				tablesWithActions := findAddedDeletedRetainedTables(typeMapOld, typeMapNew, tableDepthMapOld, tableDepthMapNew)
				outStr += "\n\n/*-----------------------Relation Model : " +
					a.Action.Action + "-----------------------------------------------*/\n"
				outStr += processTablesForModifiedApps(tablesWithActions, title, project,
					outDir, stringBuilder)
			}
		}
	}
	for _, stmt := range stmtsNew {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			appNew := appsNew[a.Action.Action]
			appOld := appsOld[a.Action.Action]
			// app added in the new sysl file
			if appNew != nil && appOld == nil {
				dataParam := &DatabaseScriptParam{
					types:   appNew.GetTypes(),
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

func processTablesForModifiedApps(tableDetails []TableDetails, title,
	project, outDir string,
	stringBuilder strings.Builder) string {
	dataParam := &DatabaseScriptModifyParam{
		tableDetails: tableDetails,
		title:        title,
		project:      project,
	}
	v := MakeDatabaseScriptView(&stringBuilder, dataParam.title, dataParam.project)
	outputStr := v.GenerateDatabaseScriptModify(dataParam)
	return outputStr
}

func findAddedDeletedRetainedTables(
	tableMapOld map[string]*sysl.Type,
	tableMapNew map[string]*sysl.Type,
	tableDepthMapOld, tableDepthMapNew map[int][]string) []TableDetails {
	tableDepthsListNew := makeSortedListOfTableDepth(tableDepthMapNew)
	tableWithAction := []TableDetails{}
	//Add the retained and added tables
	for _, depth := range tableDepthsListNew {
		tableNames := tableDepthMapNew[depth]
		sort.Strings(tableNames)
		for _, tableName := range tableNames {
			found := ifTableExisted(tableName, tableMapOld)
			if found {
				tableDetails := TableDetails{
					table:    tableMapNew[tableName],
					tableOld: tableMapOld[tableName],
					action:   "RETAIN",
					name:     tableName,
				}
				tableWithAction = append(tableWithAction, tableDetails)
			} else {
				tableDetails := TableDetails{
					table:  tableMapNew[tableName],
					action: "ADD",
					name:   tableName,
				}
				tableWithAction = append(tableWithAction, tableDetails)
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
	plantumlmixin
	CmdDatabaseScriptMod
}

func (p *modDatabaseScriptCmd) Name() string            { return "generate-script-delta" }
func (p *modDatabaseScriptCmd) RequireSyslModule() bool { return false }

func (p *modDatabaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("generatescriptdelta")

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)
	cmd.Flag("rootdir", "root dir").Short('r').StringVar(&p.root)
	cmd.Flag("orgSource", "org source sysl").Short('s').StringVar(&p.orgSource)
	cmd.Flag("newSource", "new source sysl").Short('n').StringVar(&p.newSource)

	cmd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *modDatabaseScriptCmd) Execute(args ExecuteArgs) error {
	outmap, err := GenerateModDatabaseScripts(&p.CmdDatabaseScriptMod, args.Logger)
	if err != nil {
		return err
	}
	return GenerateFromSQLMap(outmap, args.Filesystem)
}
