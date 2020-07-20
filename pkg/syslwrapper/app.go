//nolint:lll
package syslwrapper

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

// App is a simplified representation of an application in sysl
type App struct {
	Name       string
	Attributes map[string]string    // Contains app level attributes. All attributes assumed to be strings.
	Endpoints  map[string]*Endpoint // Contains all application endpoints
	Types      map[string]*Type     // Contains all type definitions in the application, excluding types defined in Params and Responses
}

// Endpoint is a simplified representation of a Sysl endpoint
type Endpoint struct {
	Summary     string                // Short human-readable description of what the endpoint does
	Description string                // Longer description of what the endpoint does
	Path        string                // Path
	Params      map[string]*Parameter // Request parameters
	Response    map[string]*Parameter // Response parameters
	Downstream  []string              // TODO: Work out the dependency graph of each application. Store downstreams in this field.
}

type Parameter struct {
	In          string // Valid values include {body, header, path, query}. Cookie not supported
	Description string
	Name        string
	Type        *Type
}

// Type represents a simplified Sysl Type.
type Type struct {
	Description string
	PrimaryKey  string           // Used to represent the primary key for type relation and key for map types.
	Optional    bool             // Used to represent if the type is optional
	Reference   string           // Used to represent type references. In the format of app:typename.
	Type        string           // Used to indicate the type. Can be one of {"bool", "int", "float", "decimal", "string", "string_8", "bytes", "date", "datetime", "xml", "uuid", "ref", "list", "map", "enum", "tuple", relation}
	Items       []*Type          // Used to represent map types, where the 0 index is the key type, and 1 index is the value type.
	Enum        map[int64]string // Used to represent enums
	Properties  map[string]*Type // Used to represent tuple and relation types.
}

type AppMapper struct {
	Module      *sysl.Module
	Types       map[string]*sysl.Type // A map of all sysl types. Key is in format "appname:typename".
	SimpleTypes map[string]*Type      // A map of all simplified types. Key is in format "appname:typename".
}

// MakeAppMapper creates an appmapper
func MakeAppMapper(m *sysl.Module) *AppMapper {
	return &AppMapper{
		Module: m,
	}
}

// ImportModule takes a sysl module and maps them into an array of simplified App structs
// It resolves any type references and cross application calls
func (am *AppMapper) Map() (map[string]*App, error) {
	var simpleApps = make(map[string]*App, 15)
	for _, app := range am.Module.Apps {
		simpleApp := am.BuildApplication(app)
		simpleApps[simpleApp.Name] = simpleApp
	}
	return simpleApps, nil
}

// BuildApplication returns a clean representation of a Sysl Application
// which hides the complexities of the protobuf generated type.
func (am *AppMapper) BuildApplication(a *sysl.Application) *App {
	cleanApp := &App{
		Name:       strings.Join(a.GetName().GetPart(), " "),
		Attributes: am.mapAttributes(a.GetAttrs()),
		Endpoints:  am.mapEndpoints(strings.Join(a.GetName().GetPart(), " "), a.GetEndpoints()),
		Types:      am.mapTypes(strings.Join(a.GetName().GetPart(), " "), a.GetTypes()),
	}
	return cleanApp
}

// Creates a map of all types
// TODO Check if colon is valid in typename
func (am *AppMapper) IndexTypes() map[string]*sysl.Type {
	var typeIndex map[string]*sysl.Type = make(map[string]*sysl.Type, 10)
	for appName, app := range am.Module.Apps {
		for typeName, typeVal := range app.Types {
			typeIndex[appName+":"+typeName] = typeVal
		}
	}
	am.Types = typeIndex
	return typeIndex
}

func (am *AppMapper) ConvertTypes() map[string]*Type {
	simpleTypes := make(map[string]*Type)
	for typeName, syslType := range am.Types {
		if simpleType := am.MapType(syslType); simpleType != nil {
			simpleTypes[typeName] = simpleType
		}
	}
	am.SimpleTypes = simpleTypes
	return simpleTypes
}

// Iterates over types and resolves typerefs
func (am *AppMapper) ResolveTypes() {
	for key, value := range am.Types {
		am.Types[key] = am.resolveType(value)
	}
}

// Only handles string attributes, support for int64, float64 and array attributes not implemented
func (am *AppMapper) mapAttributes(attributes map[string]*sysl.Attribute) map[string]string {
	var attr = make(map[string]string, 15)
	for key, value := range attributes {
		attr[key] = value.GetS()
	}
	return attr
}

func (am *AppMapper) mapTypes(appName string, syslTypes map[string]*sysl.Type) map[string]*Type {
	simpleTypes := make(map[string]*Type, 15)
	for typeName := range syslTypes {
		if simpleTypeFromLookup := am.MapType(am.Types[appName+":"+typeName]); simpleTypeFromLookup != nil {
			simpleTypes[typeName] = simpleTypeFromLookup
		}
	}
	return simpleTypes
}

func (am *AppMapper) mapEndpoints(appName string, ep map[string]*sysl.Endpoint) map[string]*Endpoint {
	endpoints := make(map[string]*Endpoint, 15)
	for key, value := range ep {
		endpoints[key] = &Endpoint{
			Summary:     am.GetAttribute(value.GetAttrs(), "summary"),
			Path:        key,
			Params:      am.mapAllParams(value),
			Response:    am.mapResponse(value.GetStmt(), appName),
			Description: value.Docstring,
		}
	}
	return endpoints
}

// Parses return statements and builds response parameters
func (am *AppMapper) mapResponse(stmt []*sysl.Statement, appName string) map[string]*Parameter {
	responseTypes := make(map[string]*Parameter, 15)
	for i := range stmt {
		var returnType *Type
		var returnName string
		if stmt[i].GetRet() == nil {
			continue
		}

		if strings.Contains(stmt[i].GetRet().Payload, "<:") {
			returnStatement := strings.Split(stmt[i].GetRet().Payload, " <: ")
			returnName = returnStatement[0]
			returnType = am.mapReturnType(returnStatement[1], appName)
		} else {
			returnType = am.mapReturnType(stmt[i].GetRet().Payload, appName)
			// Default return name of 200
			returnName = "200"
		}

		responseTypes[returnName] = &Parameter{
			Name: returnName,
			Type: returnType,
		}
	}
	return responseTypes
}

// Checks if the return value is a complex type such as sequence of string
func (am *AppMapper) mapReturnType(retValue string, appName string) *Type {
	var returnType *Type
	switch {
	case strings.Contains(retValue, "sequence of "):
		listType := strings.Replace(retValue, "sequence of ", "", 1)
		returnType = &Type{
			Type: "list",
			Items: []*Type{
				am.mapSimpleReturnType(listType, appName),
			},
		}
	case strings.Contains(retValue, "set of "):
		listType := strings.Replace(retValue, "set of ", "", 1)
		returnType = &Type{
			Type: "set",
			Items: []*Type{
				am.mapSimpleReturnType(listType, appName),
			},
		}
	default:
		returnType = am.mapSimpleReturnType(retValue, appName)
	}
	return returnType
}

// Parses return values such as Foobar.Error and Error or string
func (am *AppMapper) mapSimpleReturnType(retName string, appName string) *Type {
	var returnType *Type
	if strings.Contains(retName, ".") {
		returnRef := strings.Replace(retName, ".", ":", 1)
		returnType = &Type{
			Type:      "ref",
			Reference: returnRef,
		}
	} else {
		// Handle references to type in same app, e.g 200 <: Error
		_, ok := am.Types[appName+":"+retName]
		if ok {
			returnType = &Type{
				Type:      "ref",
				Reference: appName + ":" + retName,
			}
		}
		// Type reference not found. Trying to convert primitive
		if IsPrimitive(retName) {
			returnType = &Type{
				Type: retName,
			}
		}
	}
	return returnType
}

// IsPrimitive takes an input type string and returns true if the string is a builtin sysl primitive type.
// This includes double, int64, float64, string, bool, date, datetime
// Mostly used for parsing return statements
func IsPrimitive(typeName string) bool {
	switch typeName {
	case "double", "int64", "float64", "string", "bool", "date", "datetime": //nolint:goconst
		return true
	default:
		return false
	}
}

func (am *AppMapper) mapAllParams(endpoint *sysl.Endpoint) map[string]*Parameter {
	params := am.mapParams(endpoint.Param)
	params = am.mapRestParams(endpoint.RestParams, params)
	return params
}

func (am *AppMapper) mapRestParams(p *sysl.Endpoint_RestParams, params map[string]*Parameter) map[string]*Parameter {
	if p.GetQueryParam() != nil {
		params = am.mapQueryParams(p.QueryParam, params)
	}

	if p.GetUrlParam() != nil {
		params = am.mapURLParams(p.UrlParam, params)
	}

	return params
}

func (am *AppMapper) mapURLParams(urlParams []*sysl.Endpoint_RestParams_QueryParam, params map[string]*Parameter) map[string]*Parameter {
	if urlParams == nil {
		return params
	}
	for _, urlParam := range urlParams {
		params[urlParam.GetName()] = &Parameter{
			Name: urlParam.GetName(),
			Type: am.MapType(urlParam.GetType()),
			In:   "path",
		}
	}
	return params
}

func (am *AppMapper) mapQueryParams(queryParams []*sysl.Endpoint_RestParams_QueryParam, params map[string]*Parameter) map[string]*Parameter {
	if queryParams == nil {
		return params
	}
	for _, queryParam := range queryParams {
		params[queryParam.GetName()] = &Parameter{
			Name: queryParam.GetName(),
			Type: am.MapType(queryParam.GetType()),
			In:   "query",
		}
	}
	return params
}

func (am *AppMapper) mapParams(p []*sysl.Param) map[string]*Parameter {
	params := make(map[string]*Parameter, 15)
	for _, param := range p {
		params[param.GetName()] = &Parameter{
			Name: param.GetName(),
			Type: am.MapType(param.GetType()),
			In:   ParamIn(param.Type.GetAttrs()),
		}
	}
	return params
}

func (am *AppMapper) GetAttribute(attribute map[string]*sysl.Attribute, attributeName string) string {
	attr, ok := attribute[attributeName]
	if !ok {
		return ""
	}
	return attr.GetS()
}

func ParamIn(attrs map[string]*sysl.Attribute) string {
	if HasPattern(attrs, "body") {
		return "body"
	}
	return "header"
}

func HasPattern(attrs map[string]*sysl.Attribute, pattern string) bool {
	patterns, has := attrs["patterns"]
	if has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				if y.GetS() == pattern {
					return true
				}
			}
		}
	}
	return false
}

// Takes a sysl type and resolves all references recursively
// Resolves typerefs and collections of types into the base primitives
// TODO: Handle circular dependencies
func (am *AppMapper) resolveType(t *sysl.Type) *sysl.Type {
	if t == nil {
		return t
	}
	var resolved = t
	var err error
	switch t.Type.(type) {
	case *sysl.Type_OneOf_:
		for index, val := range t.GetOneOf().Type {
			resolved.GetOneOf().Type[index] = am.resolveType(val)
		}
	case *sysl.Type_Map_:
		resolved.GetMap().Key = am.resolveType(t.GetMap().Key)
		resolved.GetMap().Value = am.resolveType(t.GetMap().Value)
	case *sysl.Type_TypeRef:
		resolved, err = am.MapSyslType(t)
	case *sysl.Type_Tuple_:
		for key, value := range t.GetTuple().AttrDefs {
			resolved.GetTuple().AttrDefs[key] = am.resolveType(value)
		}
	case *sysl.Type_List_:
		resolved.GetList().Type = am.resolveType(t.GetList().Type)
	}

	if err != nil {
		panic(err)
	}
	return resolved
}

// MapSyslType converts types from sysl.Type to Type
func (am *AppMapper) MapSyslType(t *sysl.Type) (*sysl.Type, error) {
	var appName string
	if t == nil {
		return nil, fmt.Errorf("invalid arguments")
	}
	appName, typeName := am.GetRefDetails(t)
	resolvedType, ok := am.Types[appName+":"+typeName]
	if !ok {
		return nil, fmt.Errorf("unable to find type ref for %s", typeName)
	}
	return resolvedType, nil
}

// TypeRefs can have various formats.
// Case 1: When a type defined in the same app is referenced in a TYPE
// 	- no appname is provided in the path
// 	- the ref.path[0] element is the type name
// Case 2: When a type from another app is referenced in a TYPE
// 	- context is provided
// 	- the ref.path[0] element is the application name
// Case 3: When a type from another app is referenced in a parameter
// 	- context is NOT provided
//  - ref.appName is provided
// 	- the ref.path[0] element is the type name
// Case 4: When a type from the same app is referenced in a parameter
// 	- context is NOT provided
//  - ref.appName is provided AND is the type name (This is crazy and needs to be fixed)
func (am *AppMapper) GetRefDetails(t *sysl.Type) (appName string, typeName string) {
	ref := t.GetTypeRef().GetRef()
	if ref.GetPath() == nil {
		typeName = strings.Join(ref.Appname.Part, "")
		return appName, typeName
	}
	if len(ref.GetPath()) > 1 {
		appName = ref.Path[0]
		typeName = ref.Path[1]
	} else {
		if ref.GetAppname() != nil {
			appName = strings.Join(ref.Appname.Part, "")
		} else {
			appName = strings.Join(t.GetTypeRef().GetContext().Appname.Part, "")
		}
		typeName = strings.Join(ref.GetPath(), ".")
	}
	return appName, typeName
}

// Converts sysl type to a string representation of the type
func (am *AppMapper) MapType(t *sysl.Type) *Type {
	var ref, primaryKey, simpleType string
	var items []*Type
	var properties map[string]*Type
	var enum map[int64]string

	if t == nil {
		return nil
	}

	switch t.Type.(type) {
	case *sysl.Type_NoType_:
		simpleType = "notype"
	case *sysl.Type_Primitive_:
		simpleType = am.convertPrimitive(t.String())
	case *sysl.Type_Enum_:
		simpleType = "enum"
		enum = make(map[int64]string)
		for str, index := range t.GetEnum().GetItems() {
			enum[index] = str
		}
	case *sysl.Type_Set:
		simpleType = "set"
		items = append(items, am.MapType(t.GetSet()))
	case *sysl.Type_Sequence:
		simpleType = "list" //nolint:goconst
		items = append(items, am.MapType(t.GetSequence()))
	case *sysl.Type_List_:
		simpleType = "list" //nolint:goconst
		items = append(items, am.MapType(t.GetList().Type))
	case *sysl.Type_Map_:
		// This type isn't currently used in Sysl. Tuples are used to represent Map Types instead, using a json_map_key attribute
		simpleType = "map"
		items = append(items, am.MapType(t.GetMap().Key))
		items = append(items, am.MapType(t.GetMap().Value))
	case *sysl.Type_TypeRef:
		simpleType = "ref"
		appName, typeName := am.GetRefDetails(t)
		ref = appName + ":" + typeName
	case *sysl.Type_Tuple_:
		// Currently maps in Sysl are represented as Tuples due to limitations in the grammar.
		// We need to check for the presence of a json_map_key attribute to distinguish between tuples and maps
		if mapKey := t.GetAttrs()["json_map_key"]; mapKey != nil {
			simpleType = "map"
			primaryKey = mapKey.GetS()
		} else {
			simpleType = "tuple"
		}
		properties = make(map[string]*Type, 15)
		for k, v := range t.GetTuple().AttrDefs {
			properties[k] = am.MapType(v)
		}
	case *sysl.Type_Relation_:
		simpleType = "relation"
		properties = make(map[string]*Type, 15)
		for k, v := range t.GetRelation().AttrDefs {
			switch v.Type.(type) {
			case *sysl.Type_TypeRef:
				appName, typeName := convertTableRef(v)
				properties[k] = &Type{
					Type:      "ref",
					Reference: appName + ":" + typeName,
				}
			default:
				properties[k] = am.MapType(v)
			}
		}
		if pk := t.GetRelation().GetPrimaryKey(); pk != nil {
			if pk.GetAttrName() != nil && len(pk.GetAttrName()) > 0 {
				primaryKey = pk.GetAttrName()[0]
			}
		}
	}

	return &Type{
		Reference:  ref,
		Optional:   t.GetOpt(),
		Type:       simpleType,
		Items:      items,
		Properties: properties,
		Enum:       enum,
		PrimaryKey: primaryKey,
	}
}

// Table refs must be handled differently as elements of the path are [TableName, Fieldname]
func convertTableRef(tableRef *sysl.Type) (appName string, typeName string) {
	appName = tableRef.GetTypeRef().GetContext().GetAppname().GetPart()[0]
	typeName = tableRef.GetTypeRef().GetRef().GetPath()[0]
	return appName, typeName
}

func (am *AppMapper) convertPrimitive(typeStr string) string {
	primTypeFirstLine := strings.Split(typeStr, " ")[0]
	primType := strings.Split(primTypeFirstLine, ":")[1]
	primTypeLower := strings.ToLower(primType)
	primTypeNoSpace := strings.TrimSpace(primTypeLower)
	return primTypeNoSpace
}
