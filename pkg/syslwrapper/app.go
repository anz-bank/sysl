package syslutil

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

// App is a simplified representation of an application in sysl
type App struct {
	Name       string
	Attributes map[string]string
	Endpoints  map[string]*Endpoint
	Types      map[string]*Type
}

type Endpoint struct {
	Path       string
	Params     map[string]*Parameter
	Response   map[string]*Parameter
	Downstream []string
}

type Parameter struct {
	Name string
	Type *Type
}

type Type struct {
	Reference  string // The full name of the app where the type is defined
	Type       string
	Items      []*Type
	Enum       map[string]int64
	Properties map[string]*Type
}

type AppMapper struct {
	Module *sysl.Module
	Types  map[string]*sysl.Type // A map of all non reference sysl types
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
		simpleApp, err := am.BuildApplication(app)
		if err != nil {
			return nil, err
		}
		simpleApps[simpleApp.Name] = simpleApp
	}
	return simpleApps, nil
}

// BuildApplication returns a clean representation of a Sysl Application
// which hides the complexities of the protobuf generated type.
func (am *AppMapper) BuildApplication(a *sysl.Application) (*App, error) {
	cleanApp := &App{
		Name:       strings.Join(a.GetName().GetPart(), " "),
		Attributes: am.mapAttributes(a.GetAttrs()),
		Endpoints:  am.mapEndpoints(strings.Join(a.GetName().GetPart(), " "), a.GetEndpoints()),
		Types:      am.mapTypes(strings.Join(a.GetName().GetPart(), " "), a.GetTypes()),
	}
	return cleanApp, nil
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

// TODO: Resolve Parameters

// Iterates over types and resolves typerefs
func (am *AppMapper) resolveTypes() {
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
	for typeName, _ := range syslTypes {
		simpleTypes[typeName] = am.convertTypeToString(am.Types[appName+":"+typeName])
	}
	return simpleTypes
}

func (am *AppMapper) mapEndpoints(appName string, ep map[string]*sysl.Endpoint) map[string]*Endpoint {
	endpoints := make(map[string]*Endpoint, 15)
	for key, value := range ep {
		endpoints[key] = &Endpoint{
			Path:     key,
			Params:   am.mapAllParams(value),
			Response: am.mapResponse(value.GetStmt()),
		}
	}
	return endpoints
}

func (am *AppMapper) mapResponse(stmt []*sysl.Statement) map[string]*Parameter {
	responseTypes := make(map[string]*Parameter, 15)
	for i := range stmt {
		var returnType *Type
		var returnName string
		if stmt[i].GetRet() == nil {
			continue
		}
		if strings.Contains(stmt[i].GetRet().Payload, "<:") {
			returnStatement := strings.Split(stmt[i].GetRet().Payload, " <: ")
			returnRef := strings.Replace(returnStatement[1], ".", ":", 1)
			returnType = am.convertTypeToString(am.Types[returnRef])
			returnName = returnStatement[0]

		} else {
			// If it's a primitive, we return it
			returnName = "default"
			returnType = &Type{
				Type: stmt[i].GetRet().Payload,
			}

		}

		responseTypes[returnName] = &Parameter{
			Name: returnName,
			Type: returnType,
		}
	}

	return responseTypes
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
		fmt.Println(urlParam.GetName())
		params[urlParam.GetName()] = &Parameter{
			Name: urlParam.GetName(),
			Type: am.resolveParamType(urlParam.GetType()),
		}
	}
	return params
}

func (am *AppMapper) mapQueryParams(queryParams []*sysl.Endpoint_RestParams_QueryParam, params map[string]*Parameter) map[string]*Parameter {
	if queryParams == nil {
		return params
	}
	for _, queryParam := range queryParams {
		fmt.Println(queryParam.GetName())
		params[queryParam.GetName()] = &Parameter{
			Name: queryParam.GetName(),
			Type: am.resolveParamType(queryParam.GetType()),
		}
	}
	return params
}

func (am *AppMapper) mapParams(p []*sysl.Param) map[string]*Parameter {
	params := make(map[string]*Parameter, 15)
	for _, param := range p {
		params[param.GetName()] = &Parameter{
			Name: param.GetName(),
			Type: am.resolveParamType(param.GetType()),
		}
	}
	return params
}

func (am *AppMapper) resolveParamType(syslType *sysl.Type) *Type {
	// Lookup type in Types
	typeResolved := am.resolveType(syslType)
	return am.convertTypeToString(typeResolved)
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
		resolved, err = am.TypeFromRef(t)
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

// Resolves type references
// Type references within the same application don't have the appName, so it must be passed in
func (am *AppMapper) TypeFromRef(t *sysl.Type) (*sysl.Type, error) {
	// TypeRefs can have various formats.
	// When a type defined in the same app is referenced
	// 	- no context is provided
	// 	- the ref.path[0] element is the type name
	// When a type from another app is referenced in a parameter
	// 	- context is NOT provided
	//  - ref.appName is provided
	// 	- the ref.path[0] element is the type name
	// When a type from another app is referenced
	// 	- context is provided
	// 	- the ref.path[0] element is the application name
	var appName string
	if t == nil {
		return nil, fmt.Errorf("invalid arguments")
	}

	ref := t.GetTypeRef().GetRef()
	if ref.GetAppname() != nil {
		appName = strings.Join(ref.Appname.Part, "")
	}
	typeName := strings.Join(ref.GetPath(), ".")
	if len(ref.GetPath()) > 1 {
		appName = ref.Path[0]
		typeName = ref.Path[1]
	}
	resolvedType, ok := am.Types[appName+":"+typeName]
	if !ok {
		return nil, fmt.Errorf("unable to find type ref for %s", typeName)
	}
	return resolvedType, nil
}

// Converts sysl type to a string representatino of the type
func (am *AppMapper) convertTypeToString(t *sysl.Type) *Type {
	var simpleType string
	var items []*Type
	var properties map[string]*Type
	var enum map[string]int64
	var ref string

	if t == nil {
		fmt.Printf("empty type")
		return &Type{}
	}

	switch t.Type.(type) {
	case *sysl.Type_NoType_:
		simpleType = "notype"
	case *sysl.Type_Primitive_:
		simpleType = am.convertPrimitive(t.String())
	case *sysl.Type_Enum_:
		simpleType = "enum"
		enum = t.GetEnum().GetItems()
	case *sysl.Type_List_:
		simpleType = "list"
		items = append(items, am.convertTypeToString(t.GetList().Type))
	case *sysl.Type_Map_:
		simpleType = "map"
		items = append(items, am.convertTypeToString(t.GetMap().Key))
		items = append(items, am.convertTypeToString(t.GetMap().Value))
	case *sysl.Type_TypeRef:
		simpleType = "ref"
		ref = t.GetRef().GetAppname().String() + ":" + t.GetTypeRef().GetRef().GetPath()
	case *sysl.Type_Tuple_:
		simpleType = "tuple"
		properties = make(map[string]*Type, 15)
		for k, v := range t.GetTuple().AttrDefs {
			properties[k] = am.convertTypeToString(v)
		}
	}

	return &Type{
		Type:       simpleType,
		Items:      items,
		Properties: properties,
		Enum:       enum,
	}
}

func (am *AppMapper) convertPrimitive(typeStr string) string {
	primTypeFirstLine := strings.Split(typeStr, " ")[0]
	primType := strings.Split(primTypeFirstLine, ":")[1]
	primTypeLower := strings.ToLower(primType)
	primTypeNoSpace := strings.TrimSpace(primTypeLower)
	return primTypeNoSpace
}
