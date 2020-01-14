package database

import (
	"os"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func CreateTableDepthMap(tableMap map[string]*sysl.Type) map[int][]string {
	var completedTableDepthMap = map[int][]string{}
	var incompleteTableDepthMap = map[string]int{}
	var completeTableDepthMap = map[string]int{}
	var visitedTableAttrDepth = map[string]string{}
	for tableName := range tableMap {
		incompleteTableDepthMap[tableName] = 0
	}
	processTableDepth(tableMap, completedTableDepthMap, completeTableDepthMap, incompleteTableDepthMap,
		visitedTableAttrDepth)
	return completedTableDepthMap
}

func processTableDepth(
	tableMap map[string]*sysl.Type,
	completedTableDepthMap map[int][]string,
	completeTableDepthMap map[string]int,
	incompleteTableDepthMap map[string]int,
	visitedTableAttrs map[string]string,
) {
	for tableName := range incompleteTableDepthMap {
		processComplete, size, tempVisitedAttrs := findTableDepth(tableName, tableMap[tableName],
			visitedTableAttrs, completeTableDepthMap)
		if processComplete {
			processedTablesSlice := completedTableDepthMap[size]
			if processedTablesSlice == nil {
				processedTablesSlice = nil
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
		processTableDepth(tableMap, completedTableDepthMap, completeTableDepthMap, incompleteTableDepthMap,
			visitedTableAttrs)
	}
}

func findTableDepth(
	tableName string,
	table *sysl.Type,
	visitedTableAttrs map[string]string,
	completeTableDepthMap map[string]int,
) (bool, int, map[string]string) {
	var allAttrProcessed bool = true
	var tableDepth int = 0
	var tempVisitedAttrs = map[string]string{}
	if relEntity := table.GetRelation(); relEntity != nil {
		var attrNames []string
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

func GenerateFromSQLMap(m []ScriptOutput, fs afero.Fs, logger *logrus.Logger) error { //nolint:interfacer
	for _, e := range m {
		err := errors.Wrapf(afero.WriteFile(fs, e.filename, []byte(e.content),
			os.ModePerm), "writing %q", e.filename)
		if err != nil {
			logger.Errorf("error received while writing the file %s. The error message is - %s", e.filename, err.Error())
			return err
		}
	}
	return nil
}

func isAutoIncrementAndPrimaryKey(attrType *sysl.Type) (bool, bool) {
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
func getDataTypeAndSize(attrType *sysl.Type) (string, int64) {
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
	return syslDataType, attributeSize
}

func sortColumnNamesIntoList(attrMap map[string]*sysl.Type) []string {
	var sortedColumnNames []string
	for columnName := range attrMap {
		sortedColumnNames = append(sortedColumnNames, columnName)
	}
	sort.Strings(sortedColumnNames)
	return sortedColumnNames
}
