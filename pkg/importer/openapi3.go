package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type OpenAPI3Importer struct {
	appName  string
	pkg      string
	basePath string // Used for Swagger2.0 files which have a basePath field
	logger   *logrus.Logger
	fs       afero.Fs
	types    TypeList

	nameStack []string
}

func NewOpenAPIV3Importer(logger *logrus.Logger, fs afero.Fs) Importer {
	return &OpenAPI3Importer{
		logger: logger,
		fs:     fs,
	}
}

func (o *OpenAPI3Importer) LoadFile(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return o.Load(string(bs))
}

func (o *OpenAPI3Importer) Load(file string) (string, error) {
	loader := NewOpenAPI3Loader(o.logger, o.fs)
	spec, err := loader.LoadSwaggerFromFile(file)
	if err != nil {
		spec, _ = loader.LoadSwaggerFromData([]byte(file)) // nolint: errcheck
	}
	return o.convertSpec(spec)
}

func (o *OpenAPI3Importer) WithAppName(appName string) Importer {
	o.appName = appName
	return o
}

func (o *OpenAPI3Importer) WithPackage(packageName string) Importer {
	o.pkg = packageName
	return o
}

func NewOpenAPI3Loader(logger *logrus.Logger, fs afero.Fs) *openapi3.SwaggerLoader {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.LoadSwaggerFromURIFunc = func(
		loader *openapi3.SwaggerLoader, url *url.URL) (swagger *openapi3.Swagger, err error) {
		if url.Host == "" && url.Scheme == "" {
			logger.Debugf("Loading OpenAPI3 ref: %s", url)

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
		return nil, fmt.Errorf("unable to load OpenAPI3 URI: %s", url.String())
	}
	return loader
}

func (o *OpenAPI3Importer) pushName(n string) func() {
	o.nameStack = append(o.nameStack, n)
	return func() {
		o.popName()
	}
}

func (o *OpenAPI3Importer) popName() { o.nameStack = o.nameStack[:len(o.nameStack)-1] }

func (o *OpenAPI3Importer) convertSpec(spec *openapi3.Swagger) (string, error) {
	// Convert types
	o.types = TypeList{}
	for name, ref := range spec.Components.Schemas {
		if _, found := o.types.Find(name); !found {
			if ref.Value == nil {
				o.types.Add(NewStringAlias(name))
			} else {
				o.types.Add(o.loadTypeSchema(name, ref.Value))
			}
		}
	}

	// Convert endpoints
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
		me.Sort()
		meps = append(meps, me)
	}

	result := &bytes.Buffer{}
	err := newWriter(result, o.logger).Write(o.buildSyslInfo(spec, o.basePath), o.types, meps...)

	return result.String(), err
}

func (o *OpenAPI3Importer) buildSyslInfo(spec *openapi3.Swagger, basepath string) SyslInfo {
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
	case OpenAPI_ARRAY:
		return attrsForArray(schema)
	case OpenAPI_OBJECT:
	case OpenAPI_STRING:
		return attrsForString(schema)
	case OpenAPI_INTEGER, "int":
	}
	return nil
}

func typeNameFromSchemaRef(ref *openapi3.SchemaRef) string {
	if idx := strings.Index(ref.Ref, openapiv3DefinitionPrefix); idx >= 0 {
		return getSyslSafeName(strings.TrimPrefix(ref.Ref[idx:], openapiv3DefinitionPrefix))
	}
	switch ref.Value.Type {
	case OpenAPI_ARRAY:
		return typeNameFromSchemaRef(ref.Value.Items)
	case OpenAPI_OBJECT, OpenAPI_EMPTY:
		if ref.Value.Title != "" {
			return ref.Value.Title
		}
		return OpenAPI_OBJECT
	case OpenAPI_BOOLEAN:
		return syslutil.Type_BOOL
	case OpenAPI_STRING, OpenAPI_INTEGER, OpenAPI_NUMBER:
		return mapOpenAPITypeAndFormatToType(ref.Value.Type, ref.Value.Format, logrus.StandardLogger())
	default:
		return getSyslSafeName(ref.Value.Type)
	}
}

func nameOnlyType(name string) Type {
	return &SyslBuiltIn{name: name}
}

func (o *OpenAPI3Importer) typeAliasForSchema(ref *openapi3.SchemaRef) Type {
	name := typeNameFromSchemaRef(ref)
	t, found := o.types.Find(name)
	if !found {
		t = nameOnlyType(name)
	}
	if name == OpenAPI_OBJECT {
		t = nameOnlyType(strings.Join(o.nameStack, "_"))
	}

	if _, ok := t.(*Array); !ok && ref.Value.Type == OpenAPI_ARRAY {
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

func (o *OpenAPI3Importer) buildField(name string, prop *openapi3.SchemaRef) Field {
	f := Field{Name: name}
	typeName := typeNameFromSchemaRef(prop)

	if prop.Ref != "" {
		f.Type = nameOnlyType(typeName)
		return f
	}

	defer o.pushName(name)()

	if prop.Value.Type == OpenAPI_ARRAY && prop.Value.Items.Ref != "" {
		f.Type = &Array{Items: nameOnlyType(typeNameFromSchemaRef(prop.Value.Items))}
		f.SizeSpec = makeSizeSpec(prop.Value.MinItems, prop.Value.MaxItems)
		return f
	}

	f.Type = o.typeAliasForSchema(prop)
	isArray := prop.Value.Type == OpenAPI_ARRAY
	switch typeName {
	// possible types: string, int, bool, object, date, float
	case OpenAPI_OBJECT:
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
	case OpenAPI_STRING:
		f.SizeSpec = makeSizeSpec(prop.Value.MinLength, prop.Value.MaxLength)
		f.Attrs = attrsForStrings(prop.Value)
	}

	e := exampleAttr(prop.Value.Example)
	if e != "" {
		f.Attrs = append(f.Attrs, e)
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

func (o *OpenAPI3Importer) loadTypeSchema(name string, schema *openapi3.Schema) Type {
	name = getSyslSafeName(name)
	defer o.pushName(name)()
	switch schema.Type {
	case OpenAPI_ARRAY:
		var items Type
		if childName := typeNameFromSchemaRef(schema.Items); childName == OpenAPI_OBJECT {
			defer o.pushName("obj")()
			items = o.loadTypeSchema(name+"_obj", schema.Items.Value)
			o.types.Add(items)
		} else {
			items = o.typeAliasForSchema(schema.Items)
		}

		return &Array{baseType: baseType{name: name, attrs: getAttrs(schema)}, Items: items}
	case OpenAPI_OBJECT, OpenAPI_EMPTY:
		obj := &StandardType{
			baseType: baseType{name: name, attrs: getAttrs(schema)},
		}

		if len(schema.OneOf) != 0 {
			var fields FieldList
			for _, subSchema := range schema.OneOf {
				field := o.buildField(typeNameFromSchemaRef(subSchema), subSchema)
				fields = append(fields, field)
			}
			return &Union{
				baseType: baseType{name: name},
				Options:  fields,
			}
		}

		// Removing this as it breaks an import file.
		// AllOf means this object is composed of all of the sub-schemas (and potentially additional properties)
		// for _, subschema := range schema.AllOf {
		// 	subType := o.loadTypeSchema("", subschema.Value)

		// 	if subObj, ok := subType.(*StandardType); ok {
		// 		obj.Properties = append(obj.Properties, subObj.Properties...)
		// 	}
		// }

		for fname, prop := range schema.Properties {
			f := o.buildField(fname, prop)
			f.Optional = !utils.Contains(fname, schema.Required)
			obj.Properties = append(obj.Properties, f)
		}
		if len(obj.Properties) == 0 {
			return NewStringAlias(name)
		}
		obj.Properties.Sort()
		return obj
	default:
		if schema.Type == OpenAPI_STRING && schema.Enum != nil {
			return &Enum{baseType{name: name, attrs: attrsForStrings(schema)}}
		}
		if t, ok := checkBuiltInTypes(mapOpenAPITypeAndFormatToType(schema.Type, schema.Format, o.logger)); ok {
			return &Alias{baseType: baseType{name: name}, Target: t}
		}
		o.logger.Warnf("unknown scheme type: %s", schema.Type)
		return NewStringAlias(name)
	}
}

func (o *OpenAPI3Importer) buildEndpoint(path string, item *openapi3.PathItem) []MethodEndpoints {
	var res []MethodEndpoints
	ops := map[string]*openapi3.Operation{
		syslutil.Method_GET:    item.Get,
		syslutil.Method_PUT:    item.Put,
		syslutil.Method_POST:   item.Post,
		syslutil.Method_DELETE: item.Delete,
		syslutil.Method_PATCH:  item.Patch,
	}

	commonParams := o.buildParams(item.Parameters)
	supportedCode := regexp.MustCompile("^ok|error|[1-5][0-9][0-9]$")
	errType := regexp.MustCompile("^Error|error$")

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
				param.Field = o.fieldForMediaType(mediaType, obj, "Request")
				ep.Params.Add(param)
			}
		}
		typePrefix := getSyslSafeName(convertToSyslSafe(cleanEndpointPath(path))) + "_"
		for statusCode, resp := range op.Responses {
			text := statusCode

			respType := &StandardType{
				baseType: baseType{
					name: typePrefix + text,
				},
				Properties: FieldList{},
			}

			for mediaType, val := range resp.Value.Content {
				f := o.fieldForMediaType(mediaType, val, "")
				respType.Properties = append(respType.Properties, f)
			}
			for name := range resp.Value.Headers {
				f := Field{
					Name:  name,
					Attrs: []string{"~header"},
				}
				if f.Type == nil {
					f.Type = StringAlias
				}
				respType.Properties = append(respType.Properties, f)
			}

			r := Response{}

			if len(respType.Properties) == 1 && respType.Properties[0].Attrs[0] != "~header" {
				r.Type = respType.Properties[0].Type
				r.Type.AddAttributes(respType.Properties[0].Attrs)
			} else if len(respType.Properties) > 0 {
				respType.Properties.Sort()
				o.types.Add(respType)
				r.Type = respType
			}

			match := supportedCode.MatchString(text)
			if !match {
				o.logger.Warnf("Custom response code %s is not supported", text)
				text = "ok"
				if r.Type != nil {
					if match = errType.MatchString(r.Type.Name()); match {
						text = "error"
					}
				}
			}
			r.Text = text

			ep.Responses = append(ep.Responses, r)
		}

		me.Endpoints = append(me.Endpoints, ep)
		res = append(res, me)
	}
	return res
}

func (o *OpenAPI3Importer) buildParams(params openapi3.Parameters) Parameters {
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
			p.Field.Type = o.types.AddAndRet(&Alias{baseType: baseType{name: item.Value.Name}, Target: a})
		}
		p.Optional = !item.Value.Required
		out.Add(p)
	}
	return out
}

func (o *OpenAPI3Importer) fieldForMediaType(mediatype string, mediaObj *openapi3.MediaType, typeSuffix string) Field {
	schema := mediaObj.Schema
	tname := typeNameFromSchemaRef(schema)
	field := o.buildField(tname+typeSuffix, schema)
	if _, found := o.types.Find(tname); !found {
		o.types.Add(o.loadTypeSchema(tname, schema.Value))
	}
	if a, ok := field.Type.(*Array); ok && typeSuffix == "Request" {
		field.Type = o.types.AddAndRet(&Alias{baseType: baseType{name: field.Name}, Target: a})
	}

	if mediatype != "" {
		field.Attrs = append(field.Attrs, fmt.Sprintf(`mediatype="%s"`, mediatype))
	}

	e := exampleAttr(mediaObj.Example)
	if e != "" {
		field.Attrs = append(field.Attrs, e)
	}

	es := examplesAttr(mediaObj.Examples)
	if es != "" {
		field.Attrs = append(field.Attrs, es)
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

func exampleAttr(example interface{}) string {
	if e := exampleAttrStr(example); e != "" {
		return examplesAttrStr(map[string]string{"": e})
	}
	return ""
}

func examplesAttr(examples map[string]*openapi3.ExampleRef) string {
	if len(examples) == 0 {
		return ""
	}

	examplesStr := make(map[string]string)
	for k, e := range examples {
		examplesStr[k] = exampleAttrStr(e.Value.Value)
	}

	return examplesAttrStr(examplesStr)
}

func exampleAttrStr(example interface{}) string {
	switch t := example.(type) {
	case string:
		b, err := json.Marshal(t)
		if err != nil {
			fmt.Printf("JSON marshal example string %s error %s", t, err.Error())
			return ""
		}
		return string(b)
	case float64:
		return fmt.Sprintf(`"%s"`, strconv.FormatFloat(t, 'f', -1, 64))
	case bool:
		return fmt.Sprintf(`"%t"`, t)
	case map[string]interface{}, *openapi3.Example:
		b, err := json.Marshal(t)
		if err != nil {
			fmt.Printf("JSON marshal example string %s error %s", t, err.Error())
			return ""
		}
		return exampleAttrStr(string(b))
	case []interface{}:
		strs := []string{}
		for _, e := range t {
			strs = append(strs, exampleAttrStr(e))
		}
		return exampleAttrStr(fmt.Sprintf(`[%s]`, strings.Join(strs, ",")))
	case nil:
	default:
		fmt.Printf("Unhandled example type %T %s\n", t, t)
	}

	return ""
}

func examplesAttrStr(examples map[string]string) string {
	var b bytes.Buffer

	b.Write([]byte("examples=["))

	for _, k := range utils.OrderedKeys(examples) {
		v := examples[k]
		if b.Len() > len("examples=[") {
			b.Write([]byte{','})
		}
		b.Write([]byte(fmt.Sprintf("[\"%s\",%s]", k, v)))
	}

	b.Write([]byte{']'})

	return b.String()
}
