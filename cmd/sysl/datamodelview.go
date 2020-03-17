package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	proto "github.com/anz-bank/sysl/pkg/sysl"
)

const relationArrow = `}--`
const tupleArrow = `*--`
const entityLessThanArrow = `<< (`
const entityGreaterThanArrow = `) >>`
const classString = `class`

type DataModelParam struct {
	cmdutils.ClassLabeler
	Mod     *proto.Module
	App     *proto.Application
	Project string
	Title   string
}

type DataModelView struct {
	cmdutils.ClassLabeler
	mod           *proto.Module
	stringBuilder *strings.Builder
	symbols       map[string]*cmdutils.Var
	project       string
	title         string
}

type RelationshipParam struct {
	Entity       string
	Relationship string
	Count        uint32
}

type EntityViewParam struct {
	entityColor  string
	entityHeader string
	entityName   string
}

func MakeDataModelView(
	p cmdutils.ClassLabeler, mod *proto.Module, stringBuilder *strings.Builder,
	title, project string,
) *DataModelView {
	return &DataModelView{
		ClassLabeler:  p,
		mod:           mod,
		stringBuilder: stringBuilder,
		project:       project,
		title:         title,
		symbols:       make(map[string]*cmdutils.Var),
	}
}

func (v *DataModelView) UniqueVarForAppName(appName string) string {
	if s, ok := v.symbols[appName]; ok {
		return s.Alias
	}

	i := len(v.symbols)
	alias := fmt.Sprintf("_%d", i)
	label := v.LabelClass(appName)
	s := &cmdutils.Var{
		Agent: cmdutils.MakeAgent(map[string]*proto.Attribute{}),
		Order: i,
		Label: label,
		Alias: alias,
	}
	v.symbols[appName] = s

	return s.Alias
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
			for cnt := relationshipMap[relName][childName].Count; cnt > 0; cnt-- {
				v.stringBuilder.WriteString(fmt.Sprintf("%s %s \"%s\" %s\n", relName, viewType,
					relationshipMap[relName][childName].Relationship, relationshipMap[relName][childName].Entity))
			}
		}
	}
}

func (v *DataModelView) drawRelation(
	viewParam EntityViewParam,
	entity *proto.Type_Relation,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	entityTokens := strings.Split(viewParam.entityName, ".")
	encEntity := v.UniqueVarForAppName(entityTokens[len(entityTokens)-1])
	v.stringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s%s,%s%s {\n", classString, viewParam.entityName,
		encEntity, entityLessThanArrow, viewParam.entityHeader, viewParam.entityColor, entityGreaterThanArrow))

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
			targetEntity := v.UniqueVarForAppName(typeRef.GetRef().Path[0])
			s = fmt.Sprintf("+ %s : **%s.%s** <<FK>>\n",
				attrName,
				typeRef.GetRef().Path[0],
				typeRef.GetRef().Path[1])
			if _, exists := relationshipMap[encEntity]; !exists {
				relationshipMap[encEntity] = map[string]RelationshipParam{}
			}
			if _, mulRelation := relationshipMap[encEntity][targetEntity]; mulRelation {
				relationshipMap[encEntity][targetEntity] = RelationshipParam{
					Entity:       relationshipMap[encEntity][targetEntity].Entity,
					Relationship: relationshipMap[encEntity][targetEntity].Relationship,
					Count:        relationshipMap[encEntity][targetEntity].Count + 1,
				}
			} else {
				relationshipMap[encEntity][targetEntity] = RelationshipParam{
					Entity:       targetEntity,
					Relationship: " ",
					Count:        1,
				}
			}
		} else {
			s = fmt.Sprintf("+ %s : %s\n", attrName, strings.ToLower(attrType.GetPrimitive().String()))
		}
		v.stringBuilder.WriteString(s)
	}
	v.stringBuilder.WriteString("}\n")
}

func (v *DataModelView) drawPrimitive(
	viewParam EntityViewParam,
	entity string,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	entityTokens := strings.Split(viewParam.entityName, ".")
	encEntity := v.UniqueVarForAppName(entityTokens[len(entityTokens)-1])
	v.stringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s%s,%s%s {\n", classString, viewParam.entityName,
		encEntity, entityLessThanArrow, viewParam.entityHeader, viewParam.entityColor, entityGreaterThanArrow))

	if _, exists := relationshipMap[encEntity]; !exists {
		relationshipMap[encEntity] = map[string]RelationshipParam{}
	}
	// Add default property id for primitive types
	v.stringBuilder.WriteString(fmt.Sprintf("+ %s : %s\n", "id", strings.ToLower(entity)))
	v.stringBuilder.WriteString("}\n")
}

func (v *DataModelView) drawTuple(
	viewParam EntityViewParam,
	entity *proto.Type_Tuple,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	entityTokens := strings.Split(viewParam.entityName, ".")
	encEntity := v.UniqueVarForAppName(entityTokens[len(entityTokens)-1])
	v.stringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s%s,%s%s {\n", classString, viewParam.entityName,
		encEntity, entityLessThanArrow, viewParam.entityHeader, viewParam.entityColor, entityGreaterThanArrow))
	var relation string
	var collectionString string
	var path []string
	var isPrimitiveList bool

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
				if attrType.GetList().GetType().GetPrimitive() == proto.Type_NO_Primitive {
					path = attrType.GetList().GetType().GetTypeRef().GetRef().Path
				} else {
					isPrimitiveList = true
					path = append(path, strings.ToLower(attrType.GetList().GetType().GetPrimitive().String()))
				}
				collectionString = fmt.Sprintf("+ %s : **List <%s>**\n", attrName, path[0])
				relation = `0..*`
			case attrType.GetSet() != nil:
				if attrType.GetSet().GetPrimitive() == proto.Type_NO_Primitive {
					path = attrType.GetSet().GetTypeRef().GetRef().Path
				} else {
					isPrimitiveList = true
					path = append(path, strings.ToLower(attrType.GetSet().GetPrimitive().String()))
				}
				collectionString = fmt.Sprintf("+ %s : **Set <%s>**\n", attrName, path[0])
				relation = `0..*`
			case attrType.GetTypeRef() != nil:
				path = attrType.GetTypeRef().GetRef().Path
				collectionString = fmt.Sprintf("+ %s : **%s**\n", attrName, path[0])
				relation = `1..1 `
			default:
				continue
			}
			v.stringBuilder.WriteString(collectionString)
			if !isPrimitiveList {
				if _, mulRelation := relationshipMap[encEntity][v.UniqueVarForAppName(path[0])]; mulRelation {
					relationshipMap[encEntity][v.UniqueVarForAppName(path[0])] = RelationshipParam{
						Entity:       relationshipMap[encEntity][v.UniqueVarForAppName(path[0])].Entity,
						Relationship: relationshipMap[encEntity][v.UniqueVarForAppName(path[0])].Relationship,
						Count:        relationshipMap[encEntity][v.UniqueVarForAppName(path[0])].Count + 1,
					}
				} else {
					relationshipMap[encEntity][v.UniqueVarForAppName(path[0])] = RelationshipParam{
						Entity:       v.UniqueVarForAppName(path[0]),
						Relationship: relation,
						Count:        1,
					}
				}
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
	if dataParam.Title != "" {
		fmt.Fprintf(v.stringBuilder, "title %s\n", dataParam.Title)
	}
	v.stringBuilder.WriteString(integrationdiagram.PumlHeader)

	// sort and iterate over each entity type the selected application
	// *Type_Tuple_ OR *Type_Relation_
	typeMap := dataParam.App.GetTypes()
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
		} else if pe := entityType.GetPrimitive(); pe != proto.Type_NO_Primitive && len(strings.TrimSpace(pe.String())) > 0 {
			isRelation = false
			viewParam := EntityViewParam{
				entityColor:  `orchid`,
				entityHeader: `D`,
				entityName:   entityName,
			}
			v.drawPrimitive(viewParam, pe.String(), relationshipMap)
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
