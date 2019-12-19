package database

import (
	"fmt"

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

type DatabaseScriptView struct {
	stringBuilder *strings.Builder
	title         string
	appName       string
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

func MakeDatabaseScriptView(title, appName string,
) *DatabaseScriptView {
	var stringBuilder strings.Builder
	return &DatabaseScriptView{
		stringBuilder: &stringBuilder,
		title:         title,
		appName:       appName,
	}
}
func MakeEmptyDatabaseScriptView() *DatabaseScriptView {
	return &DatabaseScriptView{}
}

func (v *DatabaseScriptView) GenerateDatabaseScriptModify(tableDetails []TableDetails,
	dbType string) string {
	v.stringBuilder.WriteString(fmt.Sprintf("/*TITLE : %s*/\n", v.title))
	v.stringBuilder.WriteString(fmt.Sprintf(databaseScriptHeader))
	appHeader := "\n\n/*-----------------------Relation Model : " +
		v.appName + "-----------------------------------------------*/\n"
	v.stringBuilder.WriteString(fmt.Sprintf(appHeader))

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

func (v *DatabaseScriptView) GenerateDatabaseScriptCreate(tableMap map[string]*proto.Type,
	dbType string) string {
	v.stringBuilder.WriteString(fmt.Sprintf("/*TITLE : %s*/\n", v.title))
	v.stringBuilder.WriteString(fmt.Sprintf(databaseScriptHeader))
	appHeader := "\n\n/*-----------------------Relation Model : " +
		v.appName + "-----------------------------------------------*/\n"
	v.stringBuilder.WriteString(fmt.Sprintf(appHeader))
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
