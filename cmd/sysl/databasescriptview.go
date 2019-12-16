package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	proto "github.com/anz-bank/sysl/pkg/sysl"
)

const databaseScriptHeader = `/* ---------------------------------------------
Autogenerated script from sysl
--------------------------------------------- */
`
const defaultTextSize = 50
const strConst = `string`
const intConst = `integer`

type DatabaseScriptParam struct {
	ClassLabeler
	types   map[string]*proto.Type
	project string
	title   string
}

type DatabaseScriptModifyParam struct {
	ClassLabeler
	tableDetails []TableDetails
	project      string
	title        string
}

type TableDetails struct {
	table    *proto.Type
	tableOld *proto.Type
	action   string
	name     string
}

type DatabaseScriptView struct {
	ClassLabeler
	stringBuilder *strings.Builder
	symbols       map[string]*_var
	project       string
	title         string
}

func MakeDatabaseScriptView(
	p ClassLabeler, stringBuilder *strings.Builder,
	title, project string,
) *DatabaseScriptView {
	return &DatabaseScriptView{
		ClassLabeler:  p,
		stringBuilder: stringBuilder,
		project:       project,
		title:         title,
		symbols:       make(map[string]*_var),
	}
}

func MakeEmptyDatabaseScriptView() *DatabaseScriptView {
	return &DatabaseScriptView{}
}

func (v *DatabaseScriptView) GenerateDatabaseScriptModify(dataParam *DatabaseScriptModifyParam) string {
	relationshipMap := map[string]map[string]RelationshipParam{}
	if dataParam.title != "" {
		fmt.Fprintf(v.stringBuilder, "/*TITLE : %s*/\n", dataParam.title)
	}
	v.stringBuilder.WriteString(databaseScriptHeader)
	visitedAttributes := make(map[string]string)
	// sort and iterate over each entity type the selected application
	// *Type_Tuple_ OR *Type_Relation_
	tableDetails := dataParam.tableDetails
	for _, tableDetail := range tableDetails {
		switch tableDetail.action {
		case "ADD":
			v.writeCreateSQLForRelation(tableDetail.name, tableDetail.table.GetRelation(), relationshipMap, visitedAttributes)
		case "RETAIN":
			v.writeModifySQLForRelation(tableDetail.name, tableDetail.table.GetRelation(),
				tableDetail.tableOld.GetRelation(), relationshipMap, visitedAttributes)
		case "DELETE":
			v.stringBuilder.WriteString(fmt.Sprintf("DROP TABLE %s\n", tableDetail.name))
		}
	}
	return v.stringBuilder.String()
}

func (v *DatabaseScriptView) GenerateDatabaseScriptCreate(dataParam *DatabaseScriptParam) string {
	var isRelation bool
	relationshipMap := map[string]map[string]RelationshipParam{}
	if dataParam.title != "" {
		fmt.Fprintf(v.stringBuilder, "/*TITLE : %s*/\n", dataParam.title)
	}
	v.stringBuilder.WriteString(databaseScriptHeader)
	visitedAttributes := make(map[string]string)
	// sort and iterate over each entity type the selected application
	// *Type_Tuple_ OR *Type_Relation_
	typeMap := dataParam.types
	completedTableDepthMap := v.createTableDepthMap(typeMap)
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
			table := typeMap[tableName]
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
			entityType := typeMap[entityName]
			if relEntity := entityType.GetRelation(); relEntity != nil {
				isRelation = true
				v.writeCreateSQLForRelation(entityName, relEntity, relationshipMap, visitedAttributes)
			}
		}
	}
	if isRelation {
		v.writeCreateSQLForRelationship(relationshipMap, relationArrow)
	}
	return v.stringBuilder.String()
}

func (v *DatabaseScriptView) writeCreateSQLForRelationship(
	relationshipMap map[string]map[string]RelationshipParam, viewType string) {
	relNames := []string{}
	for relName := range relationshipMap {
		relNames = append(relNames, relName)
	}
	sort.Strings(relNames)
	for _, relName := range relNames {
		childNames := []string{}
		for childName := range relationshipMap[relName] {
			childNames = append(childNames, childName)
		}
		sort.Strings(childNames)
		for _, childName := range childNames {
			for cnt := relationshipMap[relName][childName].Count; cnt > 0; cnt-- {
				v.stringBuilder.WriteString(fmt.Sprintf("%s %s \"%s\" %s\n", relName, viewType,
					relationshipMap[relName][childName].Relationship, relationshipMap[relName][childName].Entity))
			}
		}
	}
}

func (v *DatabaseScriptView) writeCreateSQLForRelation(
	entityName string,
	entity *proto.Type_Relation,
	relationshipMap map[string]map[string]RelationshipParam,
	visitedAttributes map[string]string,
) {
	v.stringBuilder.WriteString(fmt.Sprintf("CREATE TABLE %s(\n", entityName))
	foreignKeyConstraints := []string{}
	primaryKeys := []string{}

	lineNumbers := []int32{}
	attrNames := []string{}
	lineNumberMap := make(map[int32]string)
	for columnName := range entity.AttrDefs {
		column := entity.AttrDefs[columnName]
		lineNumber := column.GetSourceContext().GetStart().GetLine()
		lineNumberMap[lineNumber] = columnName
		lineNumbers = append(lineNumbers, lineNumber)
	}
	sort.Slice(lineNumbers, func(i, j int) bool { return lineNumbers[i] < lineNumbers[j] })
	for _, lineNo := range lineNumbers {
		attrName := lineNumberMap[lineNo]
		attrNames = append(attrNames, attrName)
	}
	// sort and iterate over attributes
	/*attrNames := []string{}
	for attrName := range entity.AttrDefs {
		attrNames = append(attrNames, attrName)
	}
	sort.Strings(attrNames)*/
	var tableData string
	for _, attrName := range attrNames {
		attrType := entity.AttrDefs[attrName]
		s, _ := v.writeCreateSQLForAColumn(attrType, entityName, attrName, &primaryKeys,
			&foreignKeyConstraints, visitedAttributes)
		tableData += s
	}
	tableData = v.addConstraints(tableData, entityName, foreignKeyConstraints, primaryKeys)
	if strings.HasSuffix(tableData, ",") {
		tableData = tableData[:len(tableData)-1]
	}
	v.stringBuilder.WriteString(tableData)
	v.stringBuilder.WriteString("\n);\n")
}

func (v *DatabaseScriptView) writeModifySQLForRelation(
	entityName string,
	entityNew *proto.Type_Relation,
	entityOld *proto.Type_Relation,
	relationshipMap map[string]map[string]RelationshipParam,
	visitedAttributes map[string]string,
) {
	primaryKeys := []string{}
	dropColumnQueries := ""
	attrDefsNew := entityNew.AttrDefs
	attrDefsOld := entityOld.AttrDefs
	attrNamesListOld := v.sortColumnNamesIntoList(attrDefsOld)
	attrNamesListNew := v.sortColumnNamesIntoList(attrDefsNew)
	primaryKeyChanged := false
	primaryKeyExisted := false

	for _, attrNameOld := range attrNamesListOld {
		//column dropped
		attrTypeOld := attrDefsOld[attrNameOld]
		attrTypeNew := attrDefsNew[attrNameOld]
		if attrTypeNew == nil {
			_, wasDeletedAttrAPrimaryKey := v.isAutoIncrementAndPrimaryKey(attrTypeOld)
			if wasDeletedAttrAPrimaryKey {
				primaryKeyChanged = true
				primaryKeyExisted = true
			}
			dropColumnQueries += fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;\n", entityName, attrNameOld)
		}
	}
	// Create or modify all the columns

	for _, attrNameNew := range attrNamesListNew {
		attrTypeOld := attrDefsOld[attrNameNew]
		attrTypeNew := attrDefsNew[attrNameNew]
		if attrTypeOld == nil {
			//attribute added
			var foreignKeyConstraints = []string{}
			str, isNewColumnPK := v.writeCreateSQLForAColumn(attrTypeNew, entityName, attrNameNew,
				&primaryKeys, &foreignKeyConstraints, visitedAttributes)
			str = strings.TrimSpace(str)
			str = str[:len(str)-1]
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;\n", entityName, str))
			if len(foreignKeyConstraints) > 0 {
				constraint := foreignKeyConstraints[0]
				constraint = constraint[:len(constraint)-1]
				v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD %s;\n", entityName, strings.TrimSpace(constraint)))
			}
			if isNewColumnPK {
				primaryKeyChanged = true
			}
		}
		if attrTypeOld != nil {
			//column retained. Find out it anything changed about the column. And then write alter queries for those columns
			primaryKeyChangedByColumn, wasOldPrimaryKey := v.writeModifySQLForAColumn(attrTypeOld, attrTypeNew,
				entityName, attrNameNew, &primaryKeys, visitedAttributes)
			if primaryKeyChangedByColumn {
				primaryKeyChanged = true
			}
			if wasOldPrimaryKey {
				primaryKeyExisted = true
			}
		}
	}
	pkConstraintName := strings.ToUpper(entityName + "_PK")

	//DROP PK IF it existed and has changed
	if primaryKeyExisted && primaryKeyChanged {
		v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;\n", entityName, pkConstraintName))
	}
	//DELETE COLUMNS
	v.stringBuilder.WriteString(dropColumnQueries)
	//ADD A PRIMARY KEY
	if primaryKeyChanged {
		pk := v.getPrimaryKeyString(primaryKeys)
		v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY(%s);\n",
			entityName, pkConstraintName, pk))
	}
}
func (v *DatabaseScriptView) writeCreateSQLForAColumn(attrType *proto.Type, entityName, attrName string,
	primaryKeys, foreignKeyConstraints *[]string, visitedAttributes map[string]string) (string, bool) {
	var s string
	isAutoIncrement, isPrimaryKey := v.isAutoIncrementAndPrimaryKey(attrType)
	if isPrimaryKey {
		*primaryKeys = append(*primaryKeys, attrName)
	}
	if typeRef := attrType.GetTypeRef(); typeRef != nil {
		s = fmt.Sprintf("  %s %s,\n",
			attrName, visitedAttributes[typeRef.GetRef().Path[0]+"."+typeRef.GetRef().Path[1]])
		fkName := strings.ToUpper(entityName + "_" + attrName + "_FK")
		*foreignKeyConstraints = append(*foreignKeyConstraints, "  CONSTRAINT "+fkName+" FOREIGN KEY("+attrName+") REFERENCES ")
		*foreignKeyConstraints = append(*foreignKeyConstraints, typeRef.GetRef().Path[0]+" ("+typeRef.GetRef().Path[1]+"),")
	} else {
		if isAutoIncrement {
			s = fmt.Sprintf("  %s %s,\n", attrName, "bigserial")
			visitedAttributes[entityName+"."+attrName] = intConst
		} else {
			syslDataType := strings.ToLower(attrType.GetPrimitive().String())
			var attributeSize int64
			attributeSize = defaultTextSize
			if syslDataType == strConst {
				constraint := attrType.GetConstraint()
				if len(constraint) > 0 {
					length := constraint[0].GetLength()
					if length != nil {
						max := length.GetMax()
						if max > 0 {
							attributeSize = max
						}
					}
				}
			}
			var datatype = v.getPostgresDataTypes(syslDataType, attributeSize)
			s = fmt.Sprintf("  %s %s,\n", attrName, datatype)
			visitedAttributes[entityName+"."+attrName] = datatype
		}
	}
	return s, isPrimaryKey
}

func (v *DatabaseScriptView) writeModifySQLForAColumn(attrTypeOld, attrTypeNew *proto.Type, entityName,
	attrName string, primaryKeys *[]string, visitedAttributes map[string]string) (bool, bool) {
	typeRefNew := attrTypeNew.GetTypeRef()
	typeRefOld := attrTypeOld.GetTypeRef()
	primaryKeyChanged := false

	isAutoIncrementOld, isPrimaryKeyOld := v.isAutoIncrementAndPrimaryKey(attrTypeOld)
	isAutoIncrementNew, isPrimaryKeyNew := v.isAutoIncrementAndPrimaryKey(attrTypeNew)

	if isPrimaryKeyOld != isPrimaryKeyNew {
		primaryKeyChanged = true
		if isPrimaryKeyNew {
			*primaryKeys = append(*primaryKeys, attrName)
		}
	}
	datatype := ""
	fkName := strings.ToUpper(entityName + "_" + attrName + "_FK")
	if typeRefNew != nil {
		datatype = visitedAttributes[typeRefNew.GetRef().Path[0]+"."+typeRefNew.GetRef().Path[1]]
		if typeRefOld == nil {
			// typeref added. Add Foreign Key Constraint
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n",
				entityName, attrName, datatype))
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT "+fkName+" FOREIGN KEY(%s) REFERENCES %s(%s);\n",
				entityName, attrName, typeRefNew.GetRef().Path[0], typeRefNew.GetRef().Path[1]))
			visitedAttributes[entityName+"."+attrName] = datatype
		} else if !strings.EqualFold(typeRefNew.GetRef().Path[0]+"."+typeRefNew.GetRef().Path[1],
			typeRefOld.GetRef().Path[0]+"."+typeRefOld.GetRef().Path[1]) {
			//typeref changed. Drop previous foreign key constraint and add new.
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;\n", entityName, fkName))
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT "+fkName+" FOREIGN KEY(%s) REFERENCES %s(%s);\n",
				entityName, attrName, typeRefNew.GetRef().Path[0],
				typeRefNew.GetRef().Path[1]))
			visitedAttributes[entityName+"."+attrName] = datatype
		}
	} else {
		datatype = v.getDataTypeAndSize(attrTypeNew)
		datatypeOld := ""
		if typeRefOld != nil {
			// typeref removed and datatype has been added. Remove foreign key reference.
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;\n", entityName, fkName))
		} else {
			datatypeOld = v.getDataTypeAndSize(attrTypeOld)
		}
		if !strings.EqualFold(datatype, datatypeOld) {
			if isAutoIncrementNew != isAutoIncrementOld && isAutoIncrementNew {
				datatype = "bigserial"
				visitedAttributes[entityName+"."+attrName] = intConst
			} else {
				datatype = v.getDataTypeAndSize(attrTypeNew)
				visitedAttributes[entityName+"."+attrName] = datatype
			}
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n", entityName,
				attrName, datatype))
			//datatype has not changed. Check if the autoincrement has changed
		} else if isAutoIncrementNew != isAutoIncrementOld {
			if isAutoIncrementNew {
				//auto increment added for the attribute. Alter table and put datatype as bigserial
				datatype = "bigserial"
				visitedAttributes[entityName+"."+attrName] = intConst
			} else {
				//auto increment removed for the attribute. Alter table and put datatype as integer.
				datatype = v.getDataTypeAndSize(attrTypeNew)
				visitedAttributes[entityName+"."+attrName] = datatype
			}
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n",
				entityName, attrName, datatype))
		}
	}
	return primaryKeyChanged, isPrimaryKeyOld
}

func (v *DatabaseScriptView) isAutoIncrementAndPrimaryKey(attrType *proto.Type) (bool, bool) {
	isAutoIncrement := false
	isPrimaryKey := false
	if patterns := attrType.GetAttrs(); patterns != nil {
		if attributesArray := patterns["patterns"].GetA(); attributesArray != nil {
			nestedAttrs := attributesArray.GetElt()
			for _, nestedAttr := range nestedAttrs {
				if strings.EqualFold("autoinc", nestedAttr.GetS()) {
					isAutoIncrement = true
				}
				if strings.EqualFold("pk", nestedAttr.GetS()) {
					isPrimaryKey = true
				}
			}
		}
	}
	return isAutoIncrement, isPrimaryKey
}
func (v *DatabaseScriptView) getDataTypeAndSize(attrType *proto.Type) string {
	syslDataType := strings.ToLower(attrType.GetPrimitive().String())
	var attributeSize int64
	attributeSize = defaultTextSize
	if syslDataType == strConst {
		constraint := attrType.GetConstraint()
		if len(constraint) > 0 {
			length := constraint[0].GetLength()
			if length != nil {
				max := length.GetMax()
				if max > 0 {
					attributeSize = max
				}
			}
		}
	}
	return v.getPostgresDataTypes(syslDataType, attributeSize)
}

func (v *DatabaseScriptView) UniqueVarForAppName(appName string) string {
	if s, ok := v.symbols[appName]; ok {
		return s.alias
	}

	i := len(v.symbols)
	alias := fmt.Sprintf("_%d", i)
	label := v.LabelClass(appName)
	s := &_var{
		agent: makeAgent(map[string]*proto.Attribute{}),
		order: i,
		label: label,
		alias: alias,
	}
	v.symbols[appName] = s

	return s.alias
}

func (v *DatabaseScriptView) addConstraints(
	s string,
	entityName string,
	foreignKeyConstraints []string,
	primaryKeys []string,
) string {
	pk := v.getPrimaryKeyString(primaryKeys)
	if !strings.EqualFold(pk, "") {
		entityName = strings.ToUpper(entityName) + "_PK"
		s = s + "  CONSTRAINT " + entityName + " PRIMARY KEY(" + pk + "),"
	}
	for _, foreignKeyConstraint := range foreignKeyConstraints {
		s = s + "\n" + foreignKeyConstraint
	}
	return s
}

func (v *DatabaseScriptView) getPrimaryKeyString(primaryKeys []string) string {
	pk := ""
	if len(primaryKeys) > 0 {
		for curIndex, primaryKey := range primaryKeys {
			if curIndex != 0 {
				pk = pk + "," + primaryKey
			} else {
				pk = primaryKey
			}
		}
	}
	return pk
}

func (v *DatabaseScriptView) getPostgresDataTypes(input string, size int64) string {
	switch input {
	case strConst:
		return "varchar (" + strconv.FormatInt(size, 10) + ")"
	case "int":
		return intConst
	case "date":
		return "date"
	case "timestamp":
		return "timestamp"
	default:
		return "varchar (50)"
	}
}

func (v *DatabaseScriptView) sortColumnNamesIntoList(attrMap map[string]*proto.Type) []string {
	sortedColumnNames := []string{}
	for columnName := range attrMap {
		sortedColumnNames = append(sortedColumnNames, columnName)
	}
	sort.Strings(sortedColumnNames)
	return sortedColumnNames
}

func (v *DatabaseScriptView) createTableDepthMap(typeMap map[string]*proto.Type) map[int][]string {
	var completedTableDepthMap = make(map[int][]string)
	var incompleteTableDepthMap = make(map[string]int)
	var completeTableDepthMap = make(map[string]int)
	var visitedTableAttrDepth = make(map[string]string)
	for tableName := range typeMap {
		incompleteTableDepthMap[tableName] = 0
	}
	v.processTableDepth(typeMap, completedTableDepthMap, completeTableDepthMap, incompleteTableDepthMap,
		visitedTableAttrDepth)
	return completedTableDepthMap
}

func (v *DatabaseScriptView) processTableDepth(
	typeMap map[string]*proto.Type,
	completedTableDepthMap map[int][]string,
	completeTableDepthMap map[string]int,
	incompleteTableDepthMap map[string]int,
	visitedTableAttrs map[string]string,
) {
	for tableName := range incompleteTableDepthMap {
		processComplete, size, tempVisitedAttrs := v.findTableDepth(tableName, typeMap[tableName],
			visitedTableAttrs, completeTableDepthMap)
		if processComplete {
			processedTablesSlice := completedTableDepthMap[size]
			if processedTablesSlice == nil {
				processedTablesSlice = []string{}
			}
			processedTablesSlice = append(processedTablesSlice, tableName)
			completedTableDepthMap[size] = processedTablesSlice
			completeTableDepthMap[tableName] = size
			delete(incompleteTableDepthMap, tableName)
			for tempAttr := range tempVisitedAttrs {
				visitedTableAttrs[tempAttr] = tempVisitedAttrs[tempAttr]
			}
		}
	}
	if len(incompleteTableDepthMap) != 0 {
		v.processTableDepth(typeMap, completedTableDepthMap, completeTableDepthMap, incompleteTableDepthMap,
			visitedTableAttrs)
	}
}

func (v *DatabaseScriptView) findTableDepth(
	tableName string,
	table *proto.Type,
	visitedTableAttrs map[string]string,
	completeTableDepthMap map[string]int,
) (bool, int, map[string]string) {
	var allAttrProcessed bool
	allAttrProcessed = true
	var tableDepth int
	tableDepth = 0
	var tempVisitedAttrs = make(map[string]string)
	if relEntity := table.GetRelation(); relEntity != nil {
		attrNames := []string{}
		for attrName := range relEntity.AttrDefs {
			attrNames = append(attrNames, attrName)
		}
		for _, attrName := range attrNames {
			attrType := relEntity.AttrDefs[attrName]
			if typeRef := attrType.GetTypeRef(); typeRef != nil {
				if val, ok := visitedTableAttrs[typeRef.GetRef().Path[0]+"."+typeRef.GetRef().Path[1]]; ok {
					newDepth := completeTableDepthMap[typeRef.GetRef().Path[0]] + 1
					tempVisitedAttrs[tableName+"."+attrName] = val
					if newDepth > tableDepth {
						tableDepth = newDepth
					}
				} else {
					allAttrProcessed = false
				}
			} else {
				tempVisitedAttrs[tableName+"."+attrName] = attrType.GetPrimitive().String()
			}
		}
	}
	return allAttrProcessed, tableDepth, tempVisitedAttrs
}
