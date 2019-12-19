package database

import (
	"os"
	"sort"
	"strings"

	proto "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func CreateTableDepthMap(tableMap map[string]*proto.Type) map[int][]string {
	var completedTableDepthMap = make(map[int][]string)
	var incompleteTableDepthMap = make(map[string]int)
	var completeTableDepthMap = make(map[string]int)
	var visitedTableAttrDepth = make(map[string]string)
	for tableName := range tableMap {
		incompleteTableDepthMap[tableName] = 0
	}
	processTableDepth(tableMap, completedTableDepthMap, completeTableDepthMap, incompleteTableDepthMap,
		visitedTableAttrDepth)
	return completedTableDepthMap
}

func processTableDepth(
	tableMap map[string]*proto.Type,
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
		processTableDepth(tableMap, completedTableDepthMap, completeTableDepthMap, incompleteTableDepthMap,
			visitedTableAttrs)
	}
}

func findTableDepth(
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

func GenerateFromSQLMap(m []ScriptOutput, fs afero.Fs) error {
	for _, e := range m {
		err := errors.Wrapf(afero.WriteFile(fs, e.filename, []byte(e.content),
			os.ModePerm), "writing %q", e.filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func isAutoIncrementAndPrimaryKey(attrType *proto.Type) (bool, bool) {
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
func getDataTypeAndSize(attrType *proto.Type) (string, int64) {
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

func sortColumnNamesIntoList(attrMap map[string]*proto.Type) []string {
	sortedColumnNames := []string{}
	for columnName := range attrMap {
		sortedColumnNames = append(sortedColumnNames, columnName)
	}
	sort.Strings(sortedColumnNames)
	return sortedColumnNames
}
