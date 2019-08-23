package main

import (
	"fmt"
	"sort"
	"strings"

	proto "github.com/anz-bank/sysl/src/proto"
)

const relationArrow = `}--`
const tupleArrow = `*--`
const entityLessThanArrow = `<< (`
const entityGreaterThanArrow = `) >>`
const classString = `class`

type DataModelParam struct {
	ClassLabeler
	mod     *proto.Module
	app     *proto.Application
	project string
	title   string
}

type DataModelView struct {
	ClassLabeler
	mod           *proto.Module
	stringBuilder *strings.Builder
	symbols       map[string]*_var
	project       string
	title         string
}

type RelationshipParam struct {
	Entity       string
	Relationship string
}

type EntityViewParam struct {
	entityColor  string
	entityHeader string
	entityName   string
}

func MakeDataModelView(
	p ClassLabeler, mod *proto.Module, stringBuilder *strings.Builder,
	title, project string,
) *DataModelView {
	return &DataModelView{
		ClassLabeler:  p,
		mod:           mod,
		stringBuilder: stringBuilder,
		project:       project,
		title:         title,
		symbols:       make(map[string]*_var),
	}
}

func (v *DataModelView) UniqueVarForAppName(appName string) string {
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

func (v *DataModelView) drawRelationship(relationshipMap map[string]map[string]RelationshipParam, viewType string) {
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
			v.stringBuilder.WriteString(fmt.Sprintf("%s %s \"%s\" %s\n", relName, viewType,
				relationshipMap[relName][childName].Relationship, relationshipMap[relName][childName].Entity))
		}
	}
}

func (v *DataModelView) drawRelation(
	viewParam EntityViewParam,
	entity *proto.Type_Relation,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	v.stringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s%s,%s%s {\n",
		classString, viewParam.entityName, v.UniqueVarForAppName(viewParam.entityName), entityLessThanArrow,
		viewParam.entityHeader, viewParam.entityColor, entityGreaterThanArrow))
	encEntity := v.UniqueVarForAppName(viewParam.entityName)

	// sort and iterate over attributes
	attrNames := []string{}
	for attrName := range entity.AttrDefs {
		attrNames = append(attrNames, attrName)
	}
	sort.Strings(attrNames)
	for _, attrName := range attrNames {
		attrType := entity.AttrDefs[attrName]
		var s string
		if typeRef := attrType.GetTypeRef(); typeRef != nil {
			s = fmt.Sprintf("+ %s : **%s.%s** <<FK>>\n",
				attrName,
				typeRef.GetRef().Path[0],
				typeRef.GetRef().Path[1])
			if _, exists := relationshipMap[encEntity]; !exists {
				relationshipMap[encEntity] = map[string]RelationshipParam{}
			}
			relationshipMap[encEntity][v.UniqueVarForAppName(typeRef.GetRef().Path[0])] = RelationshipParam{
				Entity:       v.UniqueVarForAppName(typeRef.GetRef().Path[0]),
				Relationship: " ",
			}
		} else {
			s = fmt.Sprintf("+ %s : %s\n", attrName, strings.ToLower(attrType.GetPrimitive().String()))
		}
		v.stringBuilder.WriteString(s)
	}
	v.stringBuilder.WriteString("}\n")
}

func (v *DataModelView) drawTuple(
	viewParam EntityViewParam,
	entity *proto.Type_Tuple,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	v.stringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s%s,%s%s {\n",
		classString, viewParam.entityName, v.UniqueVarForAppName(viewParam.entityName), entityLessThanArrow,
		viewParam.entityHeader, viewParam.entityColor, entityGreaterThanArrow))
	var relation string
	var collectionString string
	var path []string
	encEntity := v.UniqueVarForAppName(viewParam.entityName)

	// sort and iterate over attributes
	attrNames := []string{}
	for attrName := range entity.AttrDefs {
		attrNames = append(attrNames, attrName)
	}
	sort.Strings(attrNames)
	for _, attrName := range attrNames {
		attrType := entity.AttrDefs[attrName]
		if _, exists := relationshipMap[encEntity]; !exists {
			relationshipMap[encEntity] = map[string]RelationshipParam{}
		}
		if attrType.GetPrimitive() == proto.Type_NO_Primitive {
			switch {
			case attrType.GetList() != nil:
				path = attrType.GetSet().GetTypeRef().GetRef().Path
				collectionString = fmt.Sprintf("+ %s : **List <%s>**\n", attrName, path[0])
				relation = `0..*`
			case attrType.GetSet() != nil:
				path = attrType.GetSet().GetTypeRef().GetRef().Path
				collectionString = fmt.Sprintf("+ %s : **Set <%s>**\n", attrName, path[0])
				relation = `0..*`
			default:
				path = attrType.GetTypeRef().GetRef().Path
				collectionString = fmt.Sprintf("+ %s : **%s**\n", attrName, path[0])
				relation = `1..1 `
			}
			v.stringBuilder.WriteString(collectionString)
			relationshipMap[encEntity][v.UniqueVarForAppName(path[0])] = RelationshipParam{
				Entity:       v.UniqueVarForAppName(path[0]),
				Relationship: relation,
			}
		} else {
			v.stringBuilder.WriteString(fmt.Sprintf("+ %s : %s\n", attrName, strings.ToLower(attrType.GetPrimitive().String())))
		}
	}
	v.stringBuilder.WriteString("}\n")
}

func (v *DataModelView) GenerateDataView(dataParam *DataModelParam) string {
	var isRelation bool
	relationshipMap := map[string]map[string]RelationshipParam{}
	v.stringBuilder.WriteString("@startuml\n")
	if dataParam.title != "" {
		fmt.Fprintf(v.stringBuilder, "title %s\n", dataParam.title)
	}
	v.stringBuilder.WriteString(PumlHeader)

	// sort and iterate over each entity type the selected application
	// *Type_Tuple_ OR *Type_Relation_
	typeMap := dataParam.app.GetTypes()
	entityNames := []string{}
	for entityName := range typeMap {
		entityNames = append(entityNames, entityName)
	}
	sort.Strings(entityNames)
	for _, entityName := range entityNames {
		entityType := typeMap[entityName]
		if relEntity := entityType.GetRelation(); relEntity != nil {
			isRelation = true
			viewParam := EntityViewParam{
				entityColor:  `orchid`,
				entityHeader: `D`,
				entityName:   entityName,
			}
			v.drawRelation(viewParam, relEntity, relationshipMap)
		} else if tupEntity := entityType.GetTuple(); tupEntity != nil {
			isRelation = false
			viewParam := EntityViewParam{
				entityColor:  `orchid`,
				entityHeader: `D`,
				entityName:   entityName,
			}
			v.drawTuple(viewParam, tupEntity, relationshipMap)
		}
	}
	if isRelation {
		v.drawRelationship(relationshipMap, relationArrow)
	} else {
		v.drawRelationship(relationshipMap, tupleArrow)
	}
	v.stringBuilder.WriteString("@enduml\n")
	return v.stringBuilder.String()
}
