package datamodeldiagram

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/sysl"
)

const classString = `class`
const relationArrow = `}--`
const tupleArrow = `*--`
const entityLessThanArrow = `<< `
const entityGreaterThanArrow = ` >>`

type DataModelParam struct {
	cmdutils.ClassLabeler
	Mod     *sysl.Module
	App     *sysl.Application
	Project string
	Title   string
	Epname  bool // If %(epname) is specified
}

type DataModelView struct {
	cmdutils.ClassLabeler
	Mod           *sysl.Module
	StringBuilder *strings.Builder
	Symbols       map[string]*cmdutils.Var
	Project       string
	Title         string
}

type RelationshipParam struct {
	Entity       string
	Relationship string
	Count        uint32
}

type EntityViewParam struct {
	EntityColor  string
	EntityHeader string
	EntityName   string
	EntityAlias  string
	IgnoredTypes map[string]struct{}
	Types        map[string]*sysl.Type
}

func MakeDataModelView(
	p cmdutils.ClassLabeler, mod *sysl.Module, stringBuilder *strings.Builder,
	title, project string,
) *DataModelView {
	return &DataModelView{
		ClassLabeler:  p,
		Mod:           mod,
		StringBuilder: stringBuilder,
		Project:       project,
		Title:         title,
		Symbols:       make(map[string]*cmdutils.Var),
	}
}

func (v *DataModelView) UniqueVarForAppName(nameParts ...string) string {
	// TODO: when DrawTuple actually separates Appname and TypeName fix this
	withoutEmptyStrings := []string{}
	for _, s := range nameParts {
		if s != "" {
			withoutEmptyStrings = append(withoutEmptyStrings, s)
		}
	}
	appName := strings.Join(withoutEmptyStrings, ".")
	if s, ok := v.Symbols[appName]; ok {
		return s.Alias
	}

	i := len(v.Symbols)
	alias := fmt.Sprintf("_%d", i)
	label := v.LabelClass(appName)
	s := &cmdutils.Var{
		Agent: cmdutils.MakeAgent(map[string]*sysl.Attribute{}),
		Order: i,
		Label: label,
		Alias: alias,
	}
	v.Symbols[appName] = s

	return s.Alias
}

func (v *DataModelView) DrawRelationship(relationshipMap map[string]map[string]RelationshipParam, viewType string) {
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
				v.StringBuilder.WriteString(fmt.Sprintf("%s %s \"%s\" %s\n", relName, viewType,
					relationshipMap[relName][childName].Relationship, relationshipMap[relName][childName].Entity))
			}
		}
	}
}

func (v *DataModelView) DrawRelation(
	viewParam EntityViewParam,
	entity *sysl.Type_Relation,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	entityTokens := strings.Split(viewParam.EntityName, ".")
	encEntity := v.UniqueVarForAppName(entityTokens[len(entityTokens)-1])
	v.StringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s(%s,%s)%s {\n", classString, viewParam.EntityName,
		encEntity, entityLessThanArrow, viewParam.EntityHeader, viewParam.EntityColor, entityGreaterThanArrow))

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
		v.StringBuilder.WriteString(s)
	}
	v.StringBuilder.WriteString("}\n")
}

func (v *DataModelView) DrawPrimitive(
	viewParam EntityViewParam,
	entity string,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	entityTokens := strings.Split(viewParam.EntityName, ".")
	encEntity := v.UniqueVarForAppName(entityTokens[len(entityTokens)-1])
	v.StringBuilder.WriteString(fmt.Sprintf("%s \"%s\" as %s %s(%s,%s) %s%s {\n",
		classString, viewParam.EntityName, encEntity,
		entityLessThanArrow, viewParam.EntityHeader, viewParam.EntityColor, strings.ToLower(entity), entityGreaterThanArrow))

	if _, exists := relationshipMap[encEntity]; !exists {
		relationshipMap[encEntity] = map[string]RelationshipParam{}
	}
	v.StringBuilder.WriteString("}\n")
}

//nolint: funlen
func (v *DataModelView) DrawTuple(
	viewParam EntityViewParam,
	entity *sysl.Type_Tuple,
	relationshipMap map[string]map[string]RelationshipParam,
) {
	var alias, typeName string
	if viewParam.EntityAlias == "" {
		alias = viewParam.EntityName
	} else {
		alias = viewParam.EntityAlias
		// add space for better formatting and allow empty space when alias == ""
		typeName = " " + viewParam.EntityName
	}
	encEntity := v.UniqueVarForAppName(strings.Split(viewParam.EntityName, ".")...)

	// this will create a class header, with alias it will look like the following:
	// class "AliasName" as _0 << (D,orchid) TypeName >> {
	// Without the alias:
	// class "TypeName" as _0 << (D,orchid) >> {
	v.StringBuilder.WriteString(
		fmt.Sprintf(
			"%s \"%s\" as %s %s(%s,%s)%s%s {\n",
			classString, alias,
			encEntity, entityLessThanArrow, viewParam.EntityHeader,
			viewParam.EntityColor, typeName, entityGreaterThanArrow,
		),
	)
	var appName string
	var relation string
	var collectionString string
	var isPrimitiveList bool

	// sort and iterate over attributes
	attrNames := []string{}
	for attrName := range entity.AttrDefs {
		attrNames = append(attrNames, attrName)
	}
	sort.Strings(attrNames)
	for _, attrName := range attrNames {
		var path = []string{}
		//If the first element before the dot isn't what we passed in, then it's referring to another app
		if arr := strings.Split(attrName, "."); len(arr) <= 1 {
			appName = strings.Split(viewParam.EntityName, ".")[0]
		} else if arr[0] != appName {
			appName = arr[0]
		}
		attrType := entity.AttrDefs[attrName]

		if _, exists := relationshipMap[encEntity]; !exists {
			relationshipMap[encEntity] = map[string]RelationshipParam{}
		}
		if attrType.GetPrimitive() == sysl.Type_NO_Primitive {
			switch {
			case attrType.GetList() != nil:
				if attrType.GetList().GetType().GetPrimitive() == sysl.Type_NO_Primitive {
					path = attrType.GetList().GetType().GetTypeRef().GetRef().Path
				} else {
					isPrimitiveList = true
					path = append(path, strings.ToLower(attrType.GetList().GetType().GetPrimitive().String()))
				}
				collectionString = fmt.Sprintf("+ %s : **List <%s>**\n", attrName, path[0])
				relation = `0..*` //nolint:goconst
			case attrType.GetSet() != nil:
				if attrType.GetSet().GetPrimitive() == sysl.Type_NO_Primitive {
					path = attrType.GetSet().GetTypeRef().GetRef().Path
				} else {
					isPrimitiveList = true
					path = append(path, strings.ToLower(attrType.GetSet().GetPrimitive().String()))
				}
				collectionString = fmt.Sprintf("+ %s : **Set <%s>**\n", attrName, path[0])
				relation = `0..*` //nolint:goconst
			case attrType.GetSequence() != nil:
				if attrType.GetSequence().GetPrimitive() == sysl.Type_NO_Primitive {
					path = attrType.GetSequence().GetTypeRef().GetRef().Path
					if len(path) > 1 {
						appName = path[0]
					}
				} else {
					isPrimitiveList = true
					path = append(path, strings.ToLower(attrType.GetSequence().GetPrimitive().String()))
				}
				fullName := strings.Join(path, ".")
				collectionString = fmt.Sprintf("+ %s : **Sequence <%s>**\n", attrName, strings.Join(path, "."))
				if _, ok := viewParam.IgnoredTypes[fullName]; ok {
					v.StringBuilder.WriteString(collectionString)
					continue
				}
				relation = `0..*` //nolint:goconst
			case attrType.GetTypeRef() != nil:
				arr := attrType.GetTypeRef().GetRef().Path
				// If the array is larger than 1 then we don't neeed appName
				// TODO: Fix this when appname and typename are parsed differently
				if len(arr) > 1 {
					appName = ""
				}
				path = append(path, strings.Join(arr, "."))
				collectionString = fmt.Sprintf("+ %s : **%s**\n", attrName, path[0])
				fullName := strings.Join(attrType.GetTypeRef().GetRef().Path, ".")
				if _, ok := viewParam.IgnoredTypes[fullName]; ok {
					v.StringBuilder.WriteString(collectionString)
					continue
				}
				relation = `1..1 `
			default:
				continue
			}
			if !isPrimitiveList {
				typeName := path[0]
				if len(path) > 1 {
					appName = path[0]
					typeName = path[1]
				}
				if viewParam.Types[appName+"."+typeName] == nil && viewParam.Types[typeName] == nil {
					v.StringBuilder.WriteString(collectionString)
					continue
				}
				v.StringBuilder.WriteString(collectionString)
				if _, mulRelation := relationshipMap[encEntity][v.UniqueVarForAppName(appName, typeName)]; mulRelation {
					relationshipMap[encEntity][v.UniqueVarForAppName(appName, typeName)] = RelationshipParam{
						Entity:       relationshipMap[encEntity][v.UniqueVarForAppName(appName, typeName)].Entity,
						Relationship: relationshipMap[encEntity][v.UniqueVarForAppName(appName, typeName)].Relationship,
						Count:        relationshipMap[encEntity][v.UniqueVarForAppName(appName, typeName)].Count + 1,
					}
				} else {
					relationshipMap[encEntity][v.UniqueVarForAppName(appName, typeName)] = RelationshipParam{
						Entity:       v.UniqueVarForAppName(appName, typeName),
						Relationship: relation,
						Count:        1,
					}
				}
			} else {
				v.StringBuilder.WriteString(collectionString)
			}
		} else {
			v.StringBuilder.WriteString(fmt.Sprintf("+ %s : %s\n", attrName, strings.ToLower(attrType.GetPrimitive().String())))
		}
	}
	v.StringBuilder.WriteString("}\n")
}

func (v *DataModelView) GenerateDataView(dataParam *DataModelParam) string {
	var isRelation bool
	appName := strings.Join(dataParam.App.Name.Part, "")
	relationshipMap := map[string]map[string]RelationshipParam{}
	v.StringBuilder.WriteString("@startuml\n")
	if dataParam.Title != "" {
		fmt.Fprintf(v.StringBuilder, "title %s\n", dataParam.Title)
	}
	v.StringBuilder.WriteString(integrationdiagram.PumlHeader)

	// sort and iterate over each entity type the selected application
	// *Type_Tuple_ OR *Type_Relation_
	typeMap := map[string]*sysl.Type{}
	ignoredTypes := map[string]struct{}{}
	//typeMap := dataParam.App.GetTypes()
	// TODO: Actually put The app/project name and the app in a struct so strings.split and join dont need to be used
	entityNames := []string{}
	for _, app := range dataParam.Mod.Apps {
		for entityName, entityValue := range app.GetTypes() {
			entityName = strings.Join(app.Name.GetPart(), "") + "." + entityName
			if entityValue.Type != nil {
				typeMap[entityName] = entityValue
				entityNames = append(entityNames, entityName)
			} else {
				ignoredTypes[entityName] = struct{}{}
			}
		}
	}

	sort.Strings(entityNames)
	for _, entityName := range entityNames {
		if dataParam.Epname && strings.Split(entityName, ".")[0] != appName {
			continue
		}
		entityType := typeMap[entityName]
		if relEntity := entityType.GetRelation(); relEntity != nil {
			isRelation = true
			viewParam := EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
				Types:        typeMap,
			}
			v.DrawRelation(viewParam, relEntity, relationshipMap)
		} else if tupEntity := entityType.GetTuple(); tupEntity != nil {
			isRelation = false
			viewParam := EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
				IgnoredTypes: ignoredTypes,
				Types:        typeMap,
			}
			v.DrawTuple(viewParam, tupEntity, relationshipMap)
		} else if pe := entityType.GetPrimitive(); pe != sysl.Type_NO_Primitive && len(strings.TrimSpace(pe.String())) > 0 {
			isRelation = false
			viewParam := EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
				Types:        typeMap,
			}
			v.DrawPrimitive(viewParam, pe.String(), relationshipMap)
		}
	}
	if isRelation {
		v.DrawRelationship(relationshipMap, relationArrow)
	} else {
		v.DrawRelationship(relationshipMap, tupleArrow)
	}
	v.StringBuilder.WriteString("@enduml\n")
	return v.StringBuilder.String()
}
