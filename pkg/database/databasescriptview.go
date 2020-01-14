package database

import (
	"fmt"
	"path/filepath"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"sort"
	"strings"
)

type TableDetails struct {
	table    *sysl.Type
	tableOld *sysl.Type
	action   string
	name     string
}

type ScriptView struct {
	stringBuilder *strings.Builder
	logger        *logrus.Logger
	title         string
}

type ScriptOutput struct {
	filename string
	content  string
}

func MakeScriptOutput(filename, content string) *ScriptOutput {
	return &ScriptOutput{
		filename: filename,
		content:  content,
	}
}

func MakeTableDetails(table, tableOld *sysl.Type,
	action, name string,
) *TableDetails {
	return &TableDetails{
		table:    table,
		tableOld: tableOld,
		action:   action,
		name:     name,
	}
}

func MakeDatabaseScriptView(title string, logger *logrus.Logger,
) *ScriptView {
	var stringBuilder strings.Builder
	return &ScriptView{
		stringBuilder: &stringBuilder,
		title:         title,
		logger:        logger,
	}
}

func (v *ScriptView) GenerateDatabaseScriptCreate(tableMap map[string]*sysl.Type,
	dbType, appName string) string {
	v.stringBuilder.WriteString(fmt.Sprintf("/*TITLE : %s*/\n", v.title))
	v.stringBuilder.WriteString(databaseScriptHeader)
	appHeader := "\n\n/*-----------------------Relation Model : " +
		appName + "-----------------------------------------------*/\n"
	v.stringBuilder.WriteString(appHeader)
	visitedAttributes := map[string]string{}
	completedTableDepthMap := CreateTableDepthMap(tableMap)
	var depthsFound []int
	for depth := range completedTableDepthMap {
		depthsFound = append(depthsFound, depth)
	}
	sort.Ints(depthsFound)
	for _, depth := range depthsFound {
		tableNames := completedTableDepthMap[depth]
		var lineNumbers []int32
		var entityNames []string
		lineNumberMap := map[int32]string{}
		for _, tableName := range tableNames {
			table := tableMap[tableName]
			lineNumber := table.GetSourceContext().GetStart().GetLine()
			lineNumberMap[lineNumber] = tableName
			lineNumbers = append(lineNumbers, lineNumber)
		}
		sort.Slice(lineNumbers, func(i, j int) bool { return lineNumbers[i] < lineNumbers[j] })
		for _, lineNo := range lineNumbers {
			entityName := lineNumberMap[lineNo]
			entityNames = append(entityNames, entityName)
		}
		for _, entityName := range entityNames {
			entityType := tableMap[entityName]
			if relEntity := entityType.GetRelation(); relEntity != nil {
				v.writeCreateSQLForATable(entityName, relEntity, visitedAttributes)
			}
		}
	}
	return v.stringBuilder.String()
}

func (v *ScriptView) ProcessModSysls(appsOld, appsNew map[string]*sysl.Application,
	appNames []string, outputDir, dbType string) []ScriptOutput {
	var outputSlice []ScriptOutput
	for _, appName := range appNames {
		appOld := appsOld[appName]
		appNew := appsNew[appName]
		if appOld != nil && appNew != nil {
			v.stringBuilder.Reset()
			typeMapOld := appOld.GetTypes()
			typeMapNew := appNew.GetTypes()
			tableDepthMapOld := CreateTableDepthMap(typeMapOld)
			tableDepthMapNew := CreateTableDepthMap(typeMapNew)
			tablesWithActions := findAddedDeletedRetainedTables(typeMapOld, typeMapNew, tableDepthMapOld, tableDepthMapNew)
			outStr := v.processTablesForModifiedApps(tablesWithActions, v.title, appName, dbType)
			outputFile := filepath.Join(outputDir, appName+SQLExtension)
			outputStruct := MakeScriptOutput(outputFile, outStr)
			outputSlice = append(outputSlice, *outputStruct)
		} else if appNew != nil && appOld == nil {
			v.stringBuilder.Reset()
			outStr := v.GenerateDatabaseScriptCreate(appNew.GetTypes(), dbType, appName)
			outputFile := filepath.Join(outputDir, appName+SQLExtension)
			outputStruct := MakeScriptOutput(outputFile, outStr)
			outputSlice = append(outputSlice, *outputStruct)
		}
	}
	return outputSlice
}

func (v *ScriptView) processTablesForModifiedApps(tableDetails []TableDetails,
	title, appName, dbType string) string {
	outputStr := v.generateDatabaseScriptModify(tableDetails, dbType, appName)
	return outputStr
}

func findAddedDeletedRetainedTables(
	tableMapOld map[string]*sysl.Type,
	tableMapNew map[string]*sysl.Type,
	tableDepthMapOld, tableDepthMapNew map[int][]string) []TableDetails {
	tableDepthsListNew := makeSortedListOfTableDepth(tableDepthMapNew)
	var tableWithAction []TableDetails
	//Add the retained and added tables
	for _, depth := range tableDepthsListNew {
		tableNames := tableDepthMapNew[depth]
		sort.Strings(tableNames)
		for _, tableName := range tableNames {
			_, ok := tableMapOld[tableName]
			if ok {
				tableDetails := MakeTableDetails(tableMapNew[tableName],
					tableMapOld[tableName], Retain, tableName)
				tableWithAction = append(tableWithAction, *tableDetails)
			} else {
				tableDetails := MakeTableDetails(tableMapNew[tableName],
					nil, Add, tableName)
				tableWithAction = append(tableWithAction, *tableDetails)
			}
		}
	}

	return tableWithAction
}

func makeSortedListOfTableDepth(tableDepthMap map[int][]string) []int {
	var depthList []int
	for depth := range tableDepthMap {
		depthList = append(depthList, depth)
	}
	sort.Ints(depthList)
	return depthList
}

func (v *ScriptView) generateDatabaseScriptModify(tableDetails []TableDetails,
	dbType, appName string) string {
	v.stringBuilder.WriteString(fmt.Sprintf("/*TITLE : %s*/\n", v.title))
	v.stringBuilder.WriteString(databaseScriptHeader)
	appHeader := "\n\n/*-----------------------Relation Model : " +
		appName + "-----------------------------------------------*/\n"
	v.stringBuilder.WriteString(appHeader)

	visitedAttributes := map[string]string{}
	for _, tableDetail := range tableDetails {
		switch tableDetail.action {
		case "ADD":
			v.writeCreateSQLForATable(tableDetail.name, tableDetail.table.GetRelation(), visitedAttributes)
		case "RETAIN":
			v.writeModifySQLForATable(tableDetail.name, tableDetail.table.GetRelation(),
				tableDetail.tableOld.GetRelation(), visitedAttributes)
		default:
			v.logger.Warnf("the table action is spcified as %s, which is not valid\n", tableDetail.action)
		}
	}
	return v.stringBuilder.String()
}
