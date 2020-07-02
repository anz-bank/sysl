package importer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
)

const openapiv3DefinitionPrefix = "#/components/schemas/"

func MakeOpenAPI3Importer(logger *logrus.Logger, basePath string, filePath string) *OpenAPI3Importer {
	return &OpenAPI3Importer{
		logger:            logger,
		externalSpecs:     make(map[string]*OpenAPI3Importer),
		types:             TypeList{},
		intermediateTypes: TypeList{},
		basePath:          basePath,
		swaggerRoot:       filePath,
	}
}

type OpenAPI3Importer struct {
	appName       string
	pkg           string
	basePath      string // Used for Swagger2.0 files which have a basePath field
	logger        *logrus.Logger
	externalSpecs map[string]*OpenAPI3Importer
	spec          *openapi3.Swagger
	types         TypeList
	// intermediateTypes is a temporary list which places the type is in parsing process still.
	// It can help to support circular dependency, like type A has an array contains type A itself.
	intermediateTypes TypeList
	swaggerRoot       string
	globalParams      Parameters
}

func (l *OpenAPI3Importer) Load(input string) (string, error) {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	url, err := url.Parse(l.swaggerRoot)
	if err != nil {
		return "", err
	}
	swagger, err := loader.LoadSwaggerFromDataWithPath([]byte(input), url)
	if err != nil {
		return "", fmt.Errorf("error loading openapi3 file:%w", err)
	}
	l.spec = swagger
	return l.Parse()
}

func (l *OpenAPI3Importer) Parse() (string, error) {
	l.convertTypes()
	endpoints := l.convertEndpoints()

	result := &bytes.Buffer{}
	w := newWriter(result, l.logger)
	if err := w.Write(l.convertInfo(OutputData{
		AppName: l.appName,
		Package: l.pkg,
	}, l.basePath), l.types, endpoints...); err != nil {
		return "", err
	}
	return result.String(), nil
}

// Set the AppName of the imported app
func (l *OpenAPI3Importer) WithAppName(appName string) Importer {
	l.appName = appName
	return l
}

// Set the package attribute of the imported app
func (l *OpenAPI3Importer) WithPackage(pkg string) Importer {
	l.pkg = pkg
	return l
}

// basepath represents the Swagger basepath value.
// This is a swagger only field that isn't relevant to openapi3
func (l *OpenAPI3Importer) convertInfo(args OutputData, basepath string) SyslInfo {
	info := SyslInfo{
		OutputData:  args,
		Title:       l.spec.Info.Title,
		Description: l.spec.Info.Description,
		OtherFields: []string{},
	}
	values := []string{
		"version", l.spec.Info.Version,
		"termsOfService", l.spec.Info.TermsOfService,
		"basePath", basepath,
	}
	if l.spec.Info.License != nil {
		values = append(values, "license", l.spec.Info.License.Name)
	}
	if len(l.spec.Servers) > 0 {
		u, err := url.Parse(l.spec.Servers[0].URL)
		if err == nil {
			values = append(values, "host", u.Hostname())
		}
	}
	for i := 0; i < len(values); i += 2 {
		key := values[i]
		val := values[i+1]
		if val != "" {
			info.OtherFields = append(info.OtherFields, key, val)
		}
	}
	return info
}

func (l *OpenAPI3Importer) convertTypes() {
	// First init the swagger -> sysl mappings
	var swaggerToSyslMappings = map[string]string{
		"boolean": "bool",
		"date":    "date",
	}
	for swaggerName, syslName := range swaggerToSyslMappings {
		l.types.Add(&ImportedBuiltInAlias{
			name:   swaggerName,
			Target: &SyslBuiltIn{syslName},
		})
	}
	for name, schema := range l.spec.Components.Schemas {
		if _, has := l.types.Find(name); !has {
			if v := schema.Value; v != nil {
				if v.Type == ObjectTypeName && len(v.Properties) == 0 {
					continue // skip
				}
			}
			_ = l.typeFromSchema(name, schema.Value)
		}
	}
	l.types.Sort()
}

// handles importing of reference types.
func (l *OpenAPI3Importer) typeFromRef(path string) Type {
	// matches with external file remote reference
	if isOpenAPIOrSwaggerExt(path) {
		if t := l.typeFromRemoteRef(path); t != nil {
			if _, recorded := l.types.Find(t.Name()); !recorded {
				l.types.Add(t)
				l.types.Sort()
			}
			return t
		}
	} else {
		path = strings.TrimPrefix(path, openapiv3DefinitionPrefix)
		if t, has := checkBuiltInTypes(path); has {
			return t
		}
		if t, ok := l.types.Find(path); ok {
			return t
		}
		// add following check in incompleted type to support circular dependency
		if t, ok := l.intermediateTypes.Find(path); ok {
			return t
		}
		if schema, has := l.spec.Components.Schemas[path]; has {
			return l.typeFromSchema(path, schema.Value)
		}
	}

	return nil
}

// OpenAPI specs can references type definitions in other files.
func (l *OpenAPI3Importer) typeFromRemoteRef(remoteRef string) Type {
	refPath, defPath := l.parseRef(remoteRef)
	if externalLoader, fileLoaded := l.externalSpecs[refPath]; fileLoaded {
		return externalLoader.typeFromRef(defPath)
	}

	l.loadExternalSchema(refPath)
	return l.externalSpecs[refPath].typeFromRef(defPath)
}

// parseRef breaks a reference string into a referencepath and a definitionpath
// It also converts swagger refs to openapi3 refs
// Remote refs are of the format:
//  #/components/schemas/Date
//  ../resources/users.yaml
// resources/users.yaml#/components/schemas/Date
func (l *OpenAPI3Importer) parseRef(ref string) (refPath string, defPath string) {
	refPath, defPath = splitRef(ref)
	defPath = toOpenAPI3Ref(defPath)
	refPath = filepath.Join(filepath.Dir(l.swaggerRoot), refPath)
	return refPath, defPath
}

func splitRef(ref string) (string, string) {
	cleaned := strings.Split(ref, "#")
	if len(cleaned) != 2 {
		return "", ""
	}
	return cleaned[0], "#" + cleaned[1]
}

func toOpenAPI3Ref(ref string) string {
	if strings.HasPrefix(ref, "#/definitions/") {
		ref = strings.Replace(ref, "#/definitions/", "#/components/schemas/", 1)
	}
	return ref
}

func (l *OpenAPI3Importer) loadExternalSchema(remoteRef string) {
	l.externalSpecs[remoteRef] = MakeOpenAPI3Importer(l.logger, "", remoteRef)
	l.externalSpecs[remoteRef].spec = l.externalSpecs[remoteRef].getExternalSpec(remoteRef)
	l.externalSpecs[remoteRef].convertTypes()
	// external refs are usually found during initEndpoints, this is to find all external refs
	l.externalSpecs[remoteRef].convertEndpoints()
}

// Grabs openapi or swagger spec from given path
func (l *OpenAPI3Importer) getExternalSpec(path string) *openapi3.Swagger {
	var swagger *openapi3.Swagger
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	formats := []Format{OpenAPI3, Swagger}
	format, err := GuessFileType(path, data, formats)
	if err != nil {
		panic(err)
	}

	switch format.Name {
	case Swagger.Name:
		swagger, _, err = convertToOpenAPI3(data)
	case OpenAPI3.Name:
		swagger, err = openapi3.NewSwaggerLoader().LoadSwaggerFromData(data)
	}
	if err != nil {
		panic(err)
	}
	return swagger
}

func (l *OpenAPI3Importer) typeFromSchemaRef(name string, ref *openapi3.SchemaRef) Type {
	if ref == nil {
		return nil
	}
	if ref.Ref != "" {
		if t := l.typeFromRef(ref.Ref); t != nil {
			return t
		}
	}
	return l.typeFromSchema(name, ref.Value)
}

func sortProperties(props FieldList) {
	sort.SliceStable(props, func(i, j int) bool {
		return strings.Compare(props[i].Name, props[j].Name) < 0
	})
}

func (l *OpenAPI3Importer) typeFromSchema(name string, schema *openapi3.Schema) Type {
	if t, found := l.types.Find(name); found {
		return t
	}
	switch schema.Type {
	case ObjectTypeName, "":
		return l.typeFromObject(name, schema)
	case ArrayTypeName:
		return l.typeFromArray(name, schema)
	default:
		if len(schema.Enum) > 0 {
			return l.types.AddAndRet(&Enum{name: name})
		}
		baseType := mapSwaggerTypeAndFormatToType(schema.Type, schema.Format, l.logger)
		if t, found := l.types.Find(baseType); found {
			return t
		}
		if s, has := l.spec.Components.Schemas[schema.Type]; has {
			return l.typeFromSchemaRef(schema.Type, s)
		}

		l.logger.Warnf("unknown schema.Type: %s", schema.Type)
		return l.types.AddAndRet(NewStringAlias(name))
	}
}

func (l *OpenAPI3Importer) typeFromObject(name string, schema *openapi3.Schema) Type {
	t := &StandardType{
		name:       getSyslSafeName(name),
		Properties: FieldList{},
	}
	l.intermediateTypes.Add(t)
	for propName, propSchema := range schema.Properties {
		var fieldType Type
		if propSchema.Value != nil && propSchema.Value.Type == ArrayTypeName {
			if arrayRef := l.typeFromRef(propSchema.Value.Items.Ref); arrayRef != nil {
				fieldType = &Array{Items: arrayRef}
			} else if arrayItemType := l.typeFromRef(propSchema.Value.Items.Value.Type); arrayItemType != nil {
				fieldType = &Array{Items: arrayItemType}
			} else if propSchema.Value.Items.Value.Type == ObjectTypeName {
				arrayObj := l.typeFromSchema(name+"_"+getSyslSafeName(propName), propSchema.Value.Items.Value)
				fieldType = &Array{Items: arrayObj}
			}
		}
		if fieldType == nil {
			fieldType = l.typeFromSchemaRef(getSyslSafeName(name)+"_"+getSyslSafeName(propName), propSchema)
		}
		f := Field{
			Name: propName,
			Type: fieldType,
		}
		if !contains(propName, schema.Required) {
			f.Optional = true
		}
		t.Properties = append(t.Properties, f)
	}
	sortProperties(t.Properties)
	if len(t.Properties) == 0 {
		return l.types.AddAndRet(NewStringAlias(name))
	}
	return l.types.AddAndRet(t)
}

func (l *OpenAPI3Importer) typeFromArray(name string, schema *openapi3.Schema) Type {
	t := &Array{
		name:  name,
		Items: l.typeFromSchemaRef(name+"_obj", schema.Items),
	}
	if name != "" {
		return l.types.AddAndRet(t)
	}
	return t
}

func (l *OpenAPI3Importer) convertEndpoints() []MethodEndpoints {
	epMap := map[string][]Endpoint{}

	l.convertGlobalParams()

	for path, item := range l.spec.Paths {
		ops := map[string]*openapi3.Operation{
			"GET":    item.Get,
			"PUT":    item.Put,
			"POST":   item.Post,
			"DELETE": item.Delete,
			"PATCH":  item.Patch,
		}

		params := l.buildParams(item.Parameters)

		for method, op := range ops {
			if op != nil {
				epMap[method] = append(epMap[method], l.convertEndpoint(path, op, params))
			}
		}
	}

	for key := range epMap {
		key := key
		sort.SliceStable(epMap[key], func(i, j int) bool {
			return strings.Compare(epMap[key][i].Path, epMap[key][j].Path) < 0
		})
	}

	var result []MethodEndpoints
	for _, method := range methodDisplayOrder {
		if eps, ok := epMap[method]; ok {
			syslSafeEps := make([]Endpoint, 0, len(eps))
			for _, e := range eps {
				syslSafeEps = append(syslSafeEps, Endpoint{
					Path:        getSyslSafeName(e.Path),
					Description: e.Description,
					Params:      e.Params,
					Responses:   e.Responses,
				})
			}
			result = append(result, MethodEndpoints{
				Method:    method,
				Endpoints: syslSafeEps,
			})
		}
	}
	return result
}

func isSchemaDefinedObject(ref *openapi3.SchemaRef) bool {
	if ref == nil {
		return false
	}
	if val := ref.Value; val != nil && ref.Ref == "" {
		switch val.Type {
		case ObjectTypeName:
			return len(val.Properties) > 0
		case ArrayTypeName:
			return val.Items != nil
		case StringTypeName:
			return val.Format != ""
		}
	}
	return true
}

func (l *OpenAPI3Importer) convertEndpoint(path string, op *openapi3.Operation, params Parameters) Endpoint {
	var responses []Response
	typePrefix := getSyslSafeName(cleanEndpointPath(path)) + "_"
	for statusCode, resp := range op.Responses {
		text := statusCode
		if statusCode[0] == '2' {
			text = "ok"
		}
		respType := &StandardType{
			name:       typePrefix + text,
			Properties: FieldList{},
		}
		respVal := l.findResponse(resp)
		for mediaType, val := range respVal.Content {
			var t Type
			// try to not generate "EXTERNAL_" object types and use string instead
			if isSchemaDefinedObject(val.Schema) {
				t = l.typeFromSchemaRef("", val.Schema)
			} else {
				t = &Alias{name: typePrefix + statusCode, Target: StringAlias}
				l.types.Add(t)
			}
			f := Field{
				Name:       t.Name(),
				Attributes: []string{fmt.Sprintf("mediatype=\"%s\"", mediaType)},
				Type:       t,
			}
			respType.Properties = append(respType.Properties, f)
		}
		for name, header := range respVal.Headers {
			f := Field{
				Name:       name,
				Attributes: []string{"~header"},
				Type:       l.typeFromSchemaRef("", header.Value.Schema),
			}
			if f.Type == nil {
				f.Type = StringAlias
			}
			respType.Properties = append(respType.Properties, f)
		}

		r := Response{Text: text}
		if len(respType.Properties) > 0 {
			if len(respType.Properties) == 1 && respType.Properties[0].Attributes[0] != "~header" {
				r.Type = respType.Properties[0].Type
			} else {
				sortProperties(respType.Properties)
				l.types.Add(respType)
				r.Type = respType
			}
		}
		responses = append(responses, r)
	}

	res := Endpoint{
		Path:        path,
		Description: op.Description,
		Responses:   responses,
		Params:      params.Extend(l.buildParams(op.Parameters)),
	}

	if op.RequestBody != nil {
		for mediaType, content := range op.RequestBody.Value.Content {
			t := l.typeFromSchemaRef("", content.Schema)
			if _, ok := t.(*SyslBuiltIn); ok && content.Schema != nil && content.Schema.Ref != "" {
				parts := strings.Split(content.Schema.Ref, "/")
				t = l.types.AddAndRet(&Alias{name: parts[len(parts)-1], Target: t})
			} else if _, ok := t.(*Array); ok {
				// Cant have a sequence/set in the sysl params, so convert to a type alias
				// try to figure out the name of the param
				if val := content.Schema.Value; val != nil && val.Type == "array" && val.Items.Ref != "" {
					name := val.Items.Ref[strings.LastIndex(val.Items.Ref, "/")+1:]
					t = l.types.AddAndRet(&Alias{name: name, Target: t})
				}
			}
			p := Param{
				Field: Field{
					Name:       t.Name() + "Request",
					Type:       t,
					Optional:   !op.RequestBody.Value.Required,
					Attributes: []string{fmt.Sprintf("mediatype=\"%s\"", mediaType)},
					SizeSpec:   nil,
				},
				In: "body",
			}
			res.Params.Add(p)
		}
	}
	return res
}

func (l *OpenAPI3Importer) convertGlobalParams() {
	l.globalParams = Parameters{
		items:       map[string]Param{},
		insertOrder: []string{},
	}
	for name, param := range l.spec.Components.Parameters {
		l.globalParams.items[name] = l.buildParam(param.Value)
		l.globalParams.insertOrder = append(l.globalParams.insertOrder, name)
	}
}

func (l *OpenAPI3Importer) findResponse(ref *openapi3.ResponseRef) *openapi3.Response {
	if ref.Value != nil {
		return ref.Value
	}
	refName := ref.Ref[strings.LastIndex(ref.Ref, "/"):]
	if ref, ok := l.spec.Components.Responses[refName]; ok {
		return l.findResponse(ref)
	}
	return &openapi3.Response{}
}

func (l *OpenAPI3Importer) buildParams(params openapi3.Parameters) Parameters {
	out := Parameters{}
	for _, param := range params {
		var paramType Param
		if param.Ref != "" {
			paramType = l.globalParams.items[strings.TrimPrefix(param.Ref, "#/components/parameters/")]
		} else {
			paramType = l.buildParam(param.Value)
		}

		// Cant have a sequence/set in the sysl params, so convert to a type alias
		// try to figure out the name of the param
		if a, ok := paramType.Type.(*Array); ok {
			paramType.Type = &Alias{Target: a, name: a.Name()}
		}

		out.Add(paramType)
	}
	return out
}

func (l *OpenAPI3Importer) buildParam(p *openapi3.Parameter) Param {
	name := p.Name
	if hasToBeSyslSafe(p.In) {
		name = convertToSyslSafe(name)
	}
	t := l.typeFromSchemaRef(fmt.Sprintf("_param_%s", convertToSyslSafe(p.Name)), p.Schema)
	return Param{
		Field: Field{
			Name:       name,
			Type:       t,
			Optional:   !p.Required,
			Attributes: nil,
			SizeSpec:   nil,
		},
		In: p.In,
	}
}
