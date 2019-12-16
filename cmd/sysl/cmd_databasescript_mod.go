package main

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func GenerateModDatabaseScripts(datagenParams *CmdContextParamDatagen,
	modelOld *sysl.Module, modelNew *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("project: %v\n", datagenParams.project)
	logger.Debugf("title: %s\n", datagenParams.title)
	logger.Debugf("filter: %s\n", datagenParams.filter)
	logger.Debugf("output: %s\n", datagenParams.output)

	spclass := constructFormatParser("", datagenParams.classFormat)

	// The "project" app that specifies the data models to be built
	var appOld, appNew *sysl.Application
	var existsOld, existsNew bool
	appOld, existsOld = modelOld.GetApps()[datagenParams.project]
	_, existsNew = modelOld.GetApps()[datagenParams.project]
	if !(existsOld || existsNew) {
		return nil, fmt.Errorf("project not found in sysl for either new or old script")
	}

	epname := "Relational-Model"
	endptOld := appOld.GetEndpoints()[epname]
	endptNew := appNew.GetEndpoints()[epname]
	outputDir := datagenParams.output
	if strings.Contains(outputDir, "%(epname)") {
		of := MakeFormatParser(datagenParams.output)
		outputDir = of.FmtOutput(datagenParams.project, epname, "", nil)
	}
	generateModDatabaseScripts(spclass, outmap, modelOld, modelNew, endptOld.GetStmt(), endptNew.GetStmt(),
		datagenParams.title, datagenParams.project, outputDir)
	return outmap, nil
}

func generateModDatabaseScripts(pclass ClassLabeler, outmap map[string]string, modOld *sysl.Module, modNew *sysl.Module,
	stmtsOld []*sysl.Statement, stmtsNew []*sysl.Statement, title, project, outDir string) {
	appsOld := modOld.GetApps()
	appsNew := modNew.GetApps()
	var stringBuilder strings.Builder
	v := MakeEmptyDatabaseScriptView()

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
				processTablesForModifiedApps(tablesWithActions, title, project,
					outDir, pclass, stringBuilder, outmap)
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
				v := MakeDatabaseScriptView(pclass, &stringBuilder, dataParam.title, dataParam.project)
				outputStr := v.GenerateDatabaseScriptCreate(dataParam)
				outmap[outDir] = outputStr
			}
		}
	}
}

func processTablesForModifiedApps(tableDetails []TableDetails, title,
	project, outDir string, pclass ClassLabeler,
	stringBuilder strings.Builder, outmap map[string]string) {
	dataParam := &DatabaseScriptModifyParam{
		tableDetails: tableDetails,
		title:        title,
		project:      project,
	}
	v := MakeDatabaseScriptView(pclass, &stringBuilder, dataParam.title, dataParam.project)
	outputStr := v.GenerateDatabaseScriptModify(dataParam)
	outmap[outDir] = outputStr
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
	CmdContextParamDatagen
}

func (p *modDatabaseScriptCmd) Name() string            { return "datamodel" }
func (p *modDatabaseScriptCmd) RequireSyslModule() bool { return true }

func (p *modDatabaseScriptCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("data")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.classFormat)

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)

	p.AddFlag(cmd)

	cmd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)
	cmd.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').StringVar(&p.filter)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *modDatabaseScriptCmd) Execute(args ExecuteArgs) error {
	outmap, err := GenerateModDatabaseScripts(&p.CmdContextParamDatagen, args.Module, args.ModuleNew, args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(outmap, args.Filesystem)
}
