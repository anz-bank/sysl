package database

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func (v *ScriptView) writeCreateSQLForATable(
	tableName string,
	table *sysl.Type_Relation,
	visitedAttributes map[string]string,
) {
	v.stringBuilder.WriteString(fmt.Sprintf("CREATE TABLE %s(\n", tableName))
	var foreignKeyConstraints, primaryKeys, attrNames []string
	var lineNumbers []int32
	lineNumberMap := map[int32]string{}
	for columnName := range table.AttrDefs {
		column := table.AttrDefs[columnName]
		lineNumber := column.GetSourceContext().GetStart().GetLine() //nolint:staticcheck
		lineNumberMap[lineNumber] = columnName
		lineNumbers = append(lineNumbers, lineNumber)
	}
	sort.Slice(lineNumbers, func(i, j int) bool { return lineNumbers[i] < lineNumbers[j] })
	for _, lineNo := range lineNumbers {
		attrName := lineNumberMap[lineNo]
		attrNames = append(attrNames, attrName)
	}
	var tableData string
	for _, attrName := range attrNames {
		attrType := table.AttrDefs[attrName]
		s, _ := v.writeCreateSQLForAColumn(attrType, tableName, attrName, &primaryKeys,
			&foreignKeyConstraints, visitedAttributes)
		tableData += s
	}
	tableData = v.addConstraints(tableData, tableName, foreignKeyConstraints, primaryKeys)
	if strings.HasSuffix(tableData, ",") {
		tableData = tableData[:len(tableData)-1]
	}
	v.stringBuilder.WriteString(tableData)
	v.stringBuilder.WriteString("\n);\n")
}

func (v *ScriptView) writeModifySQLForATable(
	tableName string,
	entityNew *sysl.Type_Relation,
	entityOld *sysl.Type_Relation,
	visitedAttributes map[string]string,
) {
	var primaryKeys []string
	dropColumnQueries := ""
	attrDefsNew := entityNew.AttrDefs
	attrDefsOld := entityOld.AttrDefs
	attrNamesListOld := sortColumnNamesIntoList(attrDefsOld)
	attrNamesListNew := sortColumnNamesIntoList(attrDefsNew)
	primaryKeyChanged := false
	primaryKeyExisted := false

	for _, attrNameOld := range attrNamesListOld {
		//column dropped
		attrTypeOld := attrDefsOld[attrNameOld]
		attrTypeNew := attrDefsNew[attrNameOld]
		if attrTypeNew == nil {
			_, wasDeletedAttrAPrimaryKey := isAutoIncrementAndPrimaryKey(attrTypeOld)
			if wasDeletedAttrAPrimaryKey {
				primaryKeyChanged = true
				primaryKeyExisted = true
			}
			dropColumnQueries += fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;\n", tableName, attrNameOld)
		}
	}
	// Create or modify all the columns

	for _, attrNameNew := range attrNamesListNew {
		attrTypeOld := attrDefsOld[attrNameNew]
		attrTypeNew := attrDefsNew[attrNameNew]
		if attrTypeOld == nil {
			//attribute added
			var foreignKeyConstraints []string
			str, isNewColumnPK := v.writeCreateSQLForAColumn(attrTypeNew, tableName, attrNameNew,
				&primaryKeys, &foreignKeyConstraints, visitedAttributes)
			str = strings.TrimSpace(str)
			str = str[:len(str)-1]
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;\n", tableName, str))
			if len(foreignKeyConstraints) > 0 {
				constraint := foreignKeyConstraints[0]
				constraint = constraint[:len(constraint)-1]
				v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD %s;\n", tableName, strings.TrimSpace(constraint)))
			}
			if isNewColumnPK {
				primaryKeyChanged = true
			}
		}
		if attrTypeOld != nil {
			//column retained. Find out it anything changed about the column. And then write alter queries for those columns
			primaryKeyChangedByColumn, wasOldPrimaryKey := v.writeModifySQLForAColumn(attrTypeOld, attrTypeNew,
				tableName, attrNameNew, &primaryKeys, visitedAttributes)
			if primaryKeyChangedByColumn {
				primaryKeyChanged = true
			}
			if wasOldPrimaryKey {
				primaryKeyExisted = true
			}
		}
	}
	pkConstraintName := strings.ToUpper(tableName + "_PK")

	//DROP PK IF it existed and has changed
	if primaryKeyExisted && primaryKeyChanged {
		v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;\n", tableName, pkConstraintName))
	}
	//DELETE COLUMNS
	v.stringBuilder.WriteString(dropColumnQueries)
	//ADD A PRIMARY KEY
	if primaryKeyChanged {
		pk := v.getPrimaryKeyString(primaryKeys)
		v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY(%s);\n",
			tableName, pkConstraintName, pk))
	}
}
func (v *ScriptView) writeCreateSQLForAColumn(attrType *sysl.Type, tableName, attrName string,
	primaryKeys, foreignKeyConstraints *[]string, visitedAttributes map[string]string) (string, bool) {
	var s string
	isAutoIncrement, isPrimaryKey := isAutoIncrementAndPrimaryKey(attrType)
	if isPrimaryKey {
		*primaryKeys = append(*primaryKeys, attrName)
	}
	if typeRef := attrType.GetTypeRef(); typeRef != nil {
		path0 := typeRef.GetRef().Path[0]
		path1 := typeRef.GetRef().Path[1]
		datatype := visitedAttributes[path0+"."+path1]
		s = fmt.Sprintf("  %s %s,\n",
			attrName, datatype)
		fkName := strings.ToUpper(tableName + "_" + attrName + "_FK")
		*foreignKeyConstraints = append(
			*foreignKeyConstraints,
			"  CONSTRAINT "+fkName+" FOREIGN KEY("+attrName+") REFERENCES "+path0+" ("+path1+"),")
		visitedAttributes[tableName+"."+attrName] = datatype
	} else {
		if isAutoIncrement {
			s = fmt.Sprintf("  %s %s,\n", attrName, "bigserial")
			visitedAttributes[tableName+"."+attrName] = bigIntConst
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
			visitedAttributes[tableName+"."+attrName] = datatype
		}
	}
	return s, isPrimaryKey
}

func (v *ScriptView) writeModifySQLForAColumn(attrTypeOld, attrTypeNew *sysl.Type, tableName,
	attrName string, primaryKeys *[]string, visitedAttributes map[string]string) (bool, bool) {
	typeRefNew := attrTypeNew.GetTypeRef()
	typeRefOld := attrTypeOld.GetTypeRef()
	primaryKeyChanged := false

	isAutoIncrementOld, isPrimaryKeyOld := isAutoIncrementAndPrimaryKey(attrTypeOld)
	isAutoIncrementNew, isPrimaryKeyNew := isAutoIncrementAndPrimaryKey(attrTypeNew)

	if isPrimaryKeyNew {
		*primaryKeys = append(*primaryKeys, attrName)
	}

	if isPrimaryKeyOld != isPrimaryKeyNew {
		primaryKeyChanged = true
	}
	datatype := ""
	fkName := strings.ToUpper(tableName + "_" + attrName + "_FK")
	if typeRefNew != nil {
		datatype = visitedAttributes[typeRefNew.GetRef().Path[0]+"."+typeRefNew.GetRef().Path[1]]
		if typeRefOld == nil {
			// typeref added. Add Foreign Key Constraint
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n",
				tableName, attrName, datatype))
			v.stringBuilder.WriteString(fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT "+fkName+" FOREIGN KEY(%s) REFERENCES %s(%s);\n",
				tableName, attrName, typeRefNew.GetRef().Path[0], typeRefNew.GetRef().Path[1]))
		}
	} else {
		syslDataType, attributeSize := getDataTypeAndSize(attrTypeNew)
		datatype = v.getPostgresDataTypes(syslDataType, attributeSize)
		datatypeOld := ""
		if typeRefOld != nil {
			// typeref removed and datatype has been added. Remove foreign key reference.
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;\n", tableName, fkName))
		} else {
			syslDataType, attributeSize := getDataTypeAndSize(attrTypeOld)
			datatypeOld = v.getPostgresDataTypes(syslDataType, attributeSize)
		}
		if !strings.EqualFold(datatype, datatypeOld) {
			syslDataType, attributeSize := getDataTypeAndSize(attrTypeNew)
			datatype = v.getPostgresDataTypes(syslDataType, attributeSize)
			v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n", tableName,
				attrName, datatype))
			//datatype has not changed. Check if the autoincrement has changed
		} else if isAutoIncrementNew != isAutoIncrementOld {
			if isAutoIncrementNew {
				//auto increment added for the attribute. Alter table and put datatype as bigserial
				sequenceName := tableName + "_" + attrName + "_seq"
				v.stringBuilder.WriteString(fmt.Sprintf("CREATE SEQUENCE %s;\n", sequenceName))
				v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n", tableName,
					attrName, datatype))
				v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT nextval('%s');\n", tableName,
					attrName, sequenceName))
				v.stringBuilder.WriteString(fmt.Sprintf("ALTER SEQUENCE %s OWNED BY %s.%s;\n", sequenceName,
					tableName, attrName))
				v.stringBuilder.WriteString(fmt.Sprintf("select setval('%s', coalesce(max(%s), 1)) from %s;\n",
					sequenceName, attrName, tableName))
				datatype = bigIntConst
			} else {
				v.stringBuilder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n",
					tableName, attrName, datatype))
			}
		}
	}
	visitedAttributes[tableName+"."+attrName] = datatype
	return primaryKeyChanged, isPrimaryKeyOld
}

func (v *ScriptView) addConstraints(
	s string,
	tableName string,
	foreignKeyConstraints []string,
	primaryKeys []string,
) string {
	pk := v.getPrimaryKeyString(primaryKeys)
	if !strings.EqualFold(pk, "") {
		tableName = strings.ToUpper(tableName) + "_PK"
		s = s + "  CONSTRAINT " + tableName + " PRIMARY KEY(" + pk + "),"
	}
	for _, foreignKeyConstraint := range foreignKeyConstraints {
		s = s + "\n" + foreignKeyConstraint
	}
	return s
}

func (v *ScriptView) getPrimaryKeyString(primaryKeys []string) string {
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

func (v *ScriptView) getPostgresDataTypes(input string, size int64) string {
	switch input {
	case strConst:
		return "varchar (" + strconv.FormatInt(size, 10) + ")"
	case "int":
		return "integer"
	case "date":
		return "date"
	default:
		return "varchar (50)"
	}
}
