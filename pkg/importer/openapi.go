package importer

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type openapiv3 struct {
	appName  string
	pkg      string
	basePath string // Used for Swagger2.0 files which have a basePath field
	logger   *logrus.Logger
	fs       afero.Fs
	types    TypeList

	nameStack []string
}

func NewOpenAPILoader(logger *logrus.Logger, fs afero.Fs) *openapi3.SwaggerLoader {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.LoadSwaggerFromURIFunc = func(
		loader *openapi3.SwaggerLoader, url *url.URL) (swagger *openapi3.Swagger, err error) {
		if url.Host == "" && url.Scheme == "" {
			logger.Debugf("Loading openapi ref: %s", url.String())
			data, err := afero.ReadFile(fs, pathFromURL(url))
			if err != nil {
				return nil, err
			}
			if strings.Contains(string(data), "swagger:") {
				swagger, _, err = convertToOpenAPI3(data, url, loader)
			} else {
				swagger, err = loader.LoadSwaggerFromDataWithPath(data, url)
			}

			if err != nil {
				return nil, err
			}
			return swagger, loader.ResolveRefsIn(swagger, url)
		}
		return nil, fmt.Errorf("unable to load openapi URI: %s", url.String())
	}
	return loader
}

func NewOpenAPIV3Importer(logger *logrus.Logger, fs afero.Fs) Importer {
	return &openapiv3{
		logger: logger,
		fs:     fs,
	}
}

func (o *openapiv3) WithAppName(appName string) Importer {
	o.appName = appName
	return o
}

func (o *openapiv3) WithPackage(packageName string) Importer {
	o.pkg = packageName
	return o
}

func (o *openapiv3) Load(file string) (string, error) {
	loader := NewOpenAPILoader(o.logger, o.fs)
	spec, err := loader.LoadSwaggerFromFile(file)
	if err != nil {
		spec, _ = loader.LoadSwaggerFromData([]byte(file)) // nolint: errcheck
	}
	return o.convertSpec(spec)
}

func (o *openapiv3) pushName(n string) func() {
	o.nameStack = append(o.nameStack, n)
	return func() {
		o.popName()
	}
}
func (o *openapiv3) popName() { o.nameStack = o.nameStack[:len(o.nameStack)-1] }

func orderedKeys(mapObj interface{}) []string {
	var typeNames []string
	for _, k := range reflect.ValueOf(mapObj).MapKeys() {
		typeNames = append(typeNames, k.String())
	}
	sort.Strings(typeNames)
	return typeNames
}

func (o *openapiv3) convertSpec(spec *openapi3.Swagger) (string, error) {
	o.types = TypeList{}
	for _, name := range orderedKeys(spec.Components.Schemas) {
		ref := spec.Components.Schemas[name]
		if _, found := o.types.Find(name); !found {
			if ref.Value == nil {
				o.types.Add(NewStringAlias(name))
			} else {
				o.types.Add(o.loadTypeSchema(name, ref.Value))
			}
		}
	}

	endpoints := make(map[string][]Endpoint, len(methodDisplayOrder))
	for _, k := range methodDisplayOrder {
		endpoints[k] = nil
	}

	for path, ep := range spec.Paths {
		meps := o.buildEndpoint(path, ep)
		for _, mep := range meps {
			endpoints[mep.Method] = append(endpoints[mep.Method], mep.Endpoints...)
		}
	}
	o.types.Sort()
	meps := make([]MethodEndpoints, 0, len(methodDisplayOrder))
	for _, k := range methodDisplayOrder {
		me := MethodEndpoints{
			Method:    k,
			Endpoints: endpoints[k],
		}
		sort.Slice(me.Endpoints, func(i, j int) bool {
			return strings.Compare(me.Endpoints[i].Path, me.Endpoints[j].Path) < 0
		})
		meps = append(meps, me)
	}
	result := &bytes.Buffer{}
	err := newWriter(result, o.logger).Write(o.buildSyslInfo(spec, o.basePath), o.types, meps...)

	return result.String(), err
}

func (o *openapiv3) buildSyslInfo(spec *openapi3.Swagger, basepath string) SyslInfo {
	info := SyslInfo{
		OutputData: OutputData{
			AppName: o.appName,
			Package: o.pkg,
		},
		Title:       spec.Info.Title,
		Description: spec.Info.Description,
		OtherFields: []string{},
	}
	values := []string{
		"version", spec.Info.Version,
		"termsOfService", spec.Info.TermsOfService,
		"basePath", basepath,
	}
	if spec.Info.License != nil {
		values = append(values, "license", spec.Info.License.Name)
	}
	if len(spec.Servers) > 0 {
		u, err := url.Parse(spec.Servers[0].URL)
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

const openapiv3DefinitionPrefix = "#/components/schemas/"

func sortProperties(props FieldList) {
	sort.SliceStable(props, func(i, j int) bool {
		return strings.Compare(props[i].Name, props[j].Name) < 0
	})
}

func attrsForArray(schema *openapi3.Schema) []string {
	var attrs []string
	for name, fn := range map[string]func() string{
		"min": func() string { return fmt.Sprint(schema.MinItems) },
		"max": func() string {
			if schema.MaxItems != nil {
				return fmt.Sprint(*schema.MaxItems)
			}
			return ""
		},
	} {
		if val := fn(); val != "" {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, name, val))
		}
	}
	return attrs
}

func attrsForString(schema *openapi3.Schema) []string {
	var attrs []string
	for name, fn := range map[string]func() string{
		"min": func() string { return fmt.Sprint(schema.MinLength) },
		"max": func() string {
			if schema.MaxLength != nil {
				return fmt.Sprint(*schema.MaxLength)
			}
			return ""
		},
		"regex": func() string { return schema.Pattern },
	} {
		if val := fn(); val != "" {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, name, val))
		}
	}
	return attrs
}

func getAttrs(schema *openapi3.Schema) []string {
	switch schema.Type {
	case ArrayTypeName:
		return attrsForArray(schema)
	case ObjectTypeName:
	case StringTypeName:
		return attrsForString(schema)
	case "integer", "int":
	}
	return nil
}

func typeNameFromSchemaRef(ref *openapi3.SchemaRef) string {
	if idx := strings.Index(ref.Ref, openapiv3DefinitionPrefix); idx >= 0 {
		return getSyslSafeName(strings.TrimPrefix(ref.Ref[idx:], openapiv3DefinitionPrefix))
	}
	switch ref.Value.Type {
	case ArrayTypeName:
		return typeNameFromSchemaRef(ref.Value.Items)
	case ObjectTypeName, "":
		return ObjectTypeName
	case "boolean":
		return "bool"
	case StringTypeName, "integer", "number":
		return mapSwaggerTypeAndFormatToType(ref.Value.Type, ref.Value.Format, logrus.StandardLogger())
	default:
		return getSyslSafeName(ref.Value.Type)
	}
}

func nameOnlyType(name string) Type {
	return &SyslBuiltIn{name: name}
}

func (o *openapiv3) typeAliasForSchema(ref *openapi3.SchemaRef) Type {
	name := typeNameFromSchemaRef(ref)
	t, found := o.types.Find(name)
	if !found {
		t = nameOnlyType(name)
	}
	if name == ObjectTypeName {
		t = nameOnlyType(strings.Join(o.nameStack, "_"))
	}

	if _, ok := t.(*Array); !ok && ref.Value.Type == "array" {
		return &Array{Items: t}
	}
	return t
}

func makeSizeSpec(min uint64, max *uint64) *sizeSpec {
	switch {
	case min > 0 && max != nil:
		return &sizeSpec{
			Min:     int(min),
			Max:     int(*max),
			MaxType: MaxSpecified,
		}
	case min > 0:
		return &sizeSpec{
			Min:     int(min),
			MaxType: OpenEnded,
		}
	}
	return nil
}

func (o *openapiv3) buildField(name string, prop *openapi3.SchemaRef) Field {
	f := Field{Name: name}
	typeName := typeNameFromSchemaRef(prop)

	if prop.Ref != "" {
		f.Type = nameOnlyType(typeName)
		return f
	}

	defer o.pushName(name)()

	if prop.Value.Type == ArrayTypeName && prop.Value.Items.Ref != "" {
		f.Type = &Array{Items: nameOnlyType(typeNameFromSchemaRef(prop.Value.Items))}
		f.SizeSpec = makeSizeSpec(prop.Value.MinItems, prop.Value.MaxItems)
		return f
	}

	f.Type = o.typeAliasForSchema(prop)
	isArray := prop.Value.Type == ArrayTypeName
	if typeName == ObjectTypeName {
		if isArray {
			prop = prop.Value.Items
		}
		ns := o.nameStack
		o.nameStack = nil
		defer func() { o.nameStack = ns }()
		t := o.loadTypeSchema(strings.Join(ns, "_"), prop.Value)
		o.types.Add(t)
		if isArray && prop.Ref == "" {
			f.Type = &Array{Items: t}
		} else {
			f.Type = t
		}
	} else if typeName == StringTypeName {
		f.SizeSpec = makeSizeSpec(prop.Value.MinLength, prop.Value.MaxLength)
		f.Attributes = attrsForStrings(prop.Value)
	}
	return f
}

func attrsForStrings(schema *openapi3.Schema) []string {
	var attrs []string
	if r := schema.Pattern; r != "" {
		attrs = append(attrs, fmt.Sprintf(`regex="%s"`, getSyslSafeName(r)))
	}
	if e := schema.Enum; len(e) != 0 && false { // remove the `&& false` when enum_values are added
		var vals []string
		for _, opt := range e {
			vals = append(vals, fmt.Sprintf(`"%s"`, opt))
		}
		sort.Strings(vals)
		attrs = append(attrs, fmt.Sprintf(`enum_values=[%s]`, strings.Join(vals, ", ")))
	}
	return attrs
}

func (o *openapiv3) loadTypeSchema(name string, schema *openapi3.Schema) Type {
	name = getSyslSafeName(name)
	defer o.pushName(name)()
	switch schema.Type {
	case "array":
		var items Type
		if childName := typeNameFromSchemaRef(schema.Items); childName == ObjectTypeName {
			defer o.pushName("obj")()
			items = o.loadTypeSchema(name+"_obj", schema.Items.Value)
			o.types.Add(items)
		} else {
			items = o.typeAliasForSchema(schema.Items)
		}
		return &Array{name: name, Items: items, Attrs: getAttrs(schema)}
	case ObjectTypeName, "":
		obj := &StandardType{
			name:       name,
			Properties: nil,
			Attributes: getAttrs(schema),
		}

		// AllOf means this object is composed of all of the sub-schemas (and potentially addiitonal properties)
		for _, subschema := range schema.AllOf {
			subType := o.loadTypeSchema("", subschema.Value)
			if subObj, ok := subType.(*StandardType); ok {
				obj.Properties = append(obj.Properties, subObj.Properties...)
			}
		}

		for fname, prop := range schema.Properties {
			f := o.buildField(fname, prop)
			f.Optional = !contains(fname, schema.Required)
			obj.Properties = append(obj.Properties, f)
		}
		if len(obj.Properties) == 0 {
			return NewStringAlias(name)
		}
		sortProperties(obj.Properties)
		return obj
	default:
		if schema.Type == "string" && schema.Enum != nil {
			return &Enum{name: name, Attrs: attrsForStrings(schema)}
		}
		if t, ok := checkBuiltInTypes(mapSwaggerTypeAndFormatToType(schema.Type, schema.Format, o.logger)); ok {
			return &Alias{name: name, Target: t}
		}
		o.logger.Warnf("unknown scheme type: %s", schema.Type)
		return NewStringAlias(name)
	}
}

func (o *openapiv3) buildEndpoint(path string, item *openapi3.PathItem) []MethodEndpoints {
	var res []MethodEndpoints
	ops := map[string]*openapi3.Operation{
		"GET":    item.Get,
		"PUT":    item.Put,
		"POST":   item.Post,
		"DELETE": item.Delete,
		"PATCH":  item.Patch,
	}

	commonParams := o.buildParams(item.Parameters)
	for method, op := range ops {
		if op == nil {
			continue
		}
		me := MethodEndpoints{Method: method}
		ep := Endpoint{
			Path:        getSyslSafeName(path),
			Description: op.Description,
			Params:      commonParams.Extend(o.buildParams(op.Parameters)),
			Responses:   nil,
		}
		if req := op.RequestBody; req != nil {
			for mediaType, obj := range req.Value.Content {
				param := Param{In: "body"}
				param.Field = o.fieldForMediaType(mediaType, obj.Schema, "Request")
				ep.Params.Add(param)
			}
		}
		typePrefix := getSyslSafeName(convertToSyslSafe(cleanEndpointPath(path))) + "_"
		for statusCode, resp := range op.Responses {
			text := "error"
			if statusCode[0] == '2' {
				text = "ok"
			}
			respType := &StandardType{
				name:       typePrefix + text,
				Properties: FieldList{},
			}
			for mediaType, val := range resp.Value.Content {
				f := o.fieldForMediaType(mediaType, val.Schema, "")
				respType.Properties = append(respType.Properties, f)
			}
			for name := range resp.Value.Headers {
				f := Field{
					Name:       name,
					Attributes: []string{"~header"},
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
					o.types.Add(respType)
					r.Type = respType
				}
			}
			ep.Responses = append(ep.Responses, r)
		}

		me.Endpoints = append(me.Endpoints, ep)
		res = append(res, me)
	}
	return res
}

func (o *openapiv3) buildParams(params openapi3.Parameters) Parameters {
	var out Parameters
	for _, item := range params {
		name := item.Value.Name
		if item.Value.In == "query" {
			name = convertToSyslSafe(name)
		}
		p := Param{
			Field: o.buildField(name, item.Value.Schema),
			In:    item.Value.In,
		}
		// Avoid putting sequences into the params
		if a, ok := p.Field.Type.(*Array); ok {
			p.Field.Type = o.types.AddAndRet(&Alias{name: item.Value.Name, Target: a})
		}
		p.Optional = !item.Value.Required
		out.Add(p)
	}
	return out
}

func (o *openapiv3) fieldForMediaType(mediatype string, schema *openapi3.SchemaRef, typeSuffix string) Field {
	tname := typeNameFromSchemaRef(schema)
	field := o.buildField(tname+typeSuffix, schema)
	if _, found := o.types.Find(tname); !found {
		o.types.Add(o.loadTypeSchema(tname, schema.Value))
	}
	if a, ok := field.Type.(*Array); ok && typeSuffix == "Request" {
		field.Type = o.types.AddAndRet(&Alias{name: field.Name, Target: a})
	}
	if mediatype != "" {
		field.Attributes = append(field.Attributes, fmt.Sprintf(`mediatype="%s"`, mediatype))
	}
	return field
}

func pathToURL(filename string) (*url.URL, error) {
	if runtime.GOOS == "windows" {
		// Windows pathing doesnt work well with the openapi3 package, so we need to fudge the URL for refs to work
		u, err := url.Parse("file:///" + filepath.ToSlash(filename))
		if err != nil {
			return nil, err
		}
		u.Scheme = ""
		return u, nil
	}
	return url.Parse(filename)
}

func pathFromURL(u *url.URL) string {
	if runtime.GOOS == "windows" {
		return filepath.FromSlash(u.Path[1:])
	}
	return u.Path
}
