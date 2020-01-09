package database

import (
	"fmt"
	"path/filepath"

	"github.com/anz-bank/sysl/pkg/sysl"
	proto "github.com/anz-bank/sysl/pkg/sysl"

	"sort"
	"strings"
)

type TableDetails struct {
	table    *proto.Type
	tableOld *proto.Type
	action   string
	name     string
}

type ScriptView struct {
	stringBuilder *strings.Builder
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

func MakeTableDetails(table, tableOld *proto.Type,
	action, name string,
) *TableDetails {
	return &TableDetails{
		table:    table,
		tableOld: tableOld,
		action:   action,
		name:     name,
	}
}

func MakeDatabaseScriptView(title string,
) *ScriptView {
	var stringBuilder strings.Builder
	return &ScriptView{
		stringBuilder: &stringBuilder,
		title:         title,
	}
}
func MakeEmptyDatabaseScriptView() *ScriptView {
	return &ScriptView{}
}

func (v *ScriptView) GenerateDatabaseScriptCreate(tableMap map[string]*proto.Type,
	dbType, appName string) string {
	v.stringBuilder.WriteString(fmt.Sprintf("/*TITLE : %s*/\n", v.title))
	v.stringBuilder.WriteString(fmt.Sprint(databaseScriptHeader))
	appHeader := "\n\n/*-----------------------Relation Model : " +
		appName + "-----------------------------------------------*/\n"
	v.stringBuilder.WriteString(fmt.Sprint(appHeader))
	visitedAttributes := make(map[string]string)
	completedTableDepthMap := CreateTableDepthMap(tableMap)
	depthsFound := []int{}
	for depth := range completedTableDepthMap {
		depthsFound = append(depthsFound, depth)
	}
	sort.Ints(depthsFound)
	for _, depth := range depthsFound {
		tableNames := completedTableDepthMap[depth]
		lineNumbers := []int32{}
		entityNames := []string{}
		lineNumberMap := make(map[int32]string)
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
	outputSlice := []ScriptOutput{}
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
	tableWithAction := []TableDetails{}
	//Add the retained and added tables
	for _, depth := range tableDepthsListNew {
		tableNames := tableDepthMapNew[depth]
		sort.Strings(tableNames)
		for _, tableName := range tableNames {
			found := ifTableExisted(tableName, tableMapOld)
			if found {
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

func (v *ScriptView) generateDatabaseScriptModify(tableDetails []TableDetails,
	dbType, appName string) string {
	v.stringBuilder.WriteString(fmt.Sprintf("/*TITLE : %s*/\n", v.title))
	v.stringBuilder.WriteString(fmt.Sprint(databaseScriptHeader))
	appHeader := "\n\n/*-----------------------Relation Model : " +
		appName + "-----------------------------------------------*/\n"
	v.stringBuilder.WriteString(fmt.Sprint(appHeader))

	visitedAttributes := make(map[string]string)
	for _, tableDetail := range tableDetails {
		switch tableDetail.action {
		case "ADD":
			v.writeCreateSQLForATable(tableDetail.name, tableDetail.table.GetRelation(), visitedAttributes)
		case "RETAIN":
			v.writeModifySQLForATable(tableDetail.name, tableDetail.table.GetRelation(),
				tableDetail.tableOld.GetRelation(), visitedAttributes)
		}
	}
	return v.stringBuilder.String()
}
