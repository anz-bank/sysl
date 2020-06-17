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

func LoadOpenAPIText(args OutputData, text string, logger *logrus.Logger) (out string, err error) {
	if strings.Contains(text, "swagger") {
		return LoadSwaggerText(args, text, logger)
	}
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData([]byte(text))
	if err != nil {
		return "", err
	}
	return importOpenAPI(args, swagger, logger, "")
}

func importOpenAPI(args OutputData,
	swagger *openapi3.Swagger,
	logger *logrus.Logger, basepath string) (out string, err error) {
	l := &loader{
		logger:            logger,
		externalSpecs:     make(map[string]*loader),
		spec:              swagger,
		types:             TypeList{},
		intermediateTypes: TypeList{},
		mode:              args.Mode,
		swaggerRoot:       args.SwaggerRoot,
	}
	l.convertTypes()
	endpoints := l.convertEndpoints()

	result := &bytes.Buffer{}
	w := newWriter(result, logger)
	if err := w.Write(l.convertInfo(args, basepath), l.types, endpoints...); err != nil {
		return "", err
	}
	return result.String(), nil
}

type loader struct {
	logger        *logrus.Logger
	externalSpecs map[string]*loader
	spec          *openapi3.Swagger
	types         TypeList
	// intermediateTypes is a temporary list which places the type is in parsing process still.
	// It can help to support circular dependency, like type A has an array contains type A itself.
	intermediateTypes TypeList
	swaggerRoot       string
	mode              string
	globalParams      Parameters
}

func (l *loader) newLoaderWithExternalSpec(path string, swagger *openapi3.Swagger) {
	l.externalSpecs[path] = &loader{
		logger:        l.logger,
		externalSpecs: make(map[string]*loader),
		spec:          swagger,
		types:         TypeList{},
		mode:          l.mode,
		swaggerRoot:   filepath.Dir(path),
	}
	l.externalSpecs[path].convertTypes()
	// external refs are usually found during initEndpoints, this is to find all external refs
	l.externalSpecs[path].convertEndpoints()
}

func (l *loader) convertInfo(args OutputData, basepath string) SyslInfo {
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

func (l *loader) convertTypes() {
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

func (l *loader) typeFromRef(path string) Type {
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

func (l *loader) typeFromRemoteRef(path string) Type {
	cleaned := strings.Split(path, "#")
	if len(cleaned) != 2 {
		return nil
	}

	refPath, defPath := cleaned[0], openapiv3DefinitionPrefix+strings.TrimPrefix(cleaned[1], "/definitions/")
	if !filepath.IsAbs(refPath) || strings.HasPrefix(refPath, l.swaggerRoot) {
		var err error
		refPath, err = filepath.Abs(filepath.Join(l.swaggerRoot, refPath))
		if err != nil {
			panic(err)
		}
	}

	if externalLoader, fileLoaded := l.externalSpecs[refPath]; fileLoaded {
		return externalLoader.typeFromRef(defPath)
	}

	l.loadExternalSchema(refPath)
	return l.externalSpecs[refPath].typeFromRef(defPath)
}

func (l *loader) loadExternalSchema(path string) {
	l.newLoaderWithExternalSpec(path, l.getOpenapi3(path))
}

func guessYamlType(filename string, data []byte) string {
	for _, check := range []string{ModeSwagger, ModeOpenAPI} {
		if strings.Contains(string(data), check) {
			return check
		}
	}

	return "unknown"
}
func (l *loader) getOpenapi3(path string) *openapi3.Swagger {
	var swagger *openapi3.Swagger
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	mode := guessYamlType(path, data)
	switch mode {
	case ModeSwagger:
		swagger, _, err = convertToOpenapiv3(data)
	case ModeOpenAPI:
		swagger, err = openapi3.NewSwaggerLoader().LoadSwaggerFromData(data)
	default:
		panic("unknown mode: " + mode)
	}
	if err != nil {
		panic(err)
	}
	return swagger
}

func (l *loader) typeFromSchemaRef(name string, ref *openapi3.SchemaRef) Type {
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

func (l *loader) typeFromSchema(name string, schema *openapi3.Schema) Type {
	if t, found := l.types.Find(name); found {
		return t
	}
	switch schema.Type {
	case ObjectTypeName, "":
		t := &StandardType{
			name:       getSyslSafeEndpoint(name),
			Properties: FieldList{},
		}
		l.intermediateTypes.Add(t)
		for pname, pschema := range schema.Properties {
			var fieldType Type
			if pschema.Value != nil && pschema.Value.Type == ArrayTypeName {
				if atype := l.typeFromRef(pschema.Value.Items.Ref); atype != nil {
					fieldType = &Array{Items: atype}
				} else if atype := l.typeFromRef(pschema.Value.Items.Value.Type); atype != nil {
					fieldType = &Array{Items: atype}
				} else if pschema.Value.Items.Value.Type == ObjectTypeName {
					atype := l.typeFromSchema(name+"_"+getSyslSafeEndpoint(pname), pschema.Value.Items.Value)
					fieldType = &Array{Items: atype}
				}
			}
			if fieldType == nil {
				fieldType = l.typeFromSchemaRef(getSyslSafeEndpoint(name)+"_"+getSyslSafeEndpoint(pname), pschema)
			}
			f := Field{
				Name: getSyslSafeEndpoint(pname),
				Type: fieldType,
			}
			if !contains(pname, schema.Required) {
				f.Optional = true
			}
			t.Properties = append(t.Properties, f)
		}
		sortProperties(t.Properties)
		if len(t.Properties) == 0 {
			return l.types.AddAndRet(NewStringAlias(name))
		}
		return l.types.AddAndRet(t)
	case ArrayTypeName:
		t := &Array{
			name:  name,
			Items: l.typeFromSchemaRef(name+"_obj", schema.Items),
		}
		if name != "" {
			return l.types.AddAndRet(t)
		}
		return t
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

func (l *loader) convertEndpoints() []MethodEndpoints {
	epMap := map[string][]Endpoint{}

	l.initGlobalParams()

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
				epMap[method] = append(epMap[method], l.initEndpoint(path, op, params))
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
					Path:        getSyslSafeEndpoint(e.Path),
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

func (l *loader) initEndpoint(path string, op *openapi3.Operation, params Parameters) Endpoint {
	var responses []Response
	typePrefix := strings.NewReplacer(
		"/", "_",
		"{", "_",
		"}", "_",
		"-", "_").Replace(path) + "_"
	for statusCode, resp := range op.Responses {
		text := "error"
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

func (l *loader) initGlobalParams() {
	l.globalParams = Parameters{
		items:       map[string]Param{},
		insertOrder: []string{},
	}
	for name, param := range l.spec.Components.Parameters {
		l.globalParams.items[name] = l.buildParam(param.Value)
		l.globalParams.insertOrder = append(l.globalParams.insertOrder, name)
	}
}

func (l *loader) findResponse(ref *openapi3.ResponseRef) *openapi3.Response {
	if ref.Value != nil {
		return ref.Value
	}
	refName := ref.Ref[strings.LastIndex(ref.Ref, "/"):]
	if ref, ok := l.spec.Components.Responses[refName]; ok {
		return l.findResponse(ref)
	}
	return &openapi3.Response{}
}

func (l *loader) buildParams(params openapi3.Parameters) Parameters {
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

func (l *loader) buildParam(p *openapi3.Parameter) Param {
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

func hasToBeSyslSafe(in string) bool {
	return strings.ToLower(in) == "query"
}

func convertToSyslSafe(name string) string {
	if !strings.ContainsAny(name, "- ") {
		return name
	}

	syslSafe := strings.Builder{}
	toUppercase := false
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '-':
			toUppercase = true
		case ' ':
			continue
		default:
			if toUppercase {
				syslSafe.WriteString(strings.ToUpper(string(name[i])))
				toUppercase = false
			} else {
				syslSafe.WriteByte(name[i])
			}
		}
	}
	return syslSafe.String()
}
