package importer

import (
	"bytes"
	"fmt"
	"net/url"
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
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.LoadSwaggerFromURIFunc = func(loader *openapi3.SwaggerLoader, url *url.URL) (swagger *openapi3.Swagger, err error) {
		if url.Host == "" && url.Scheme == "" {
			o.logger.Infof("Loading openapi ref: %s", url.String())
			data, err := afero.ReadFile(o.fs, url.Path)
			if err != nil {
				return nil, err
			}
			return loader.LoadSwaggerFromDataWithPath(data, url)
		}
		return nil, fmt.Errorf("unable to load openapi URI: %s", url.String())
	}

	spec, err := loader.LoadSwaggerFromFile(file)
	if err != nil {
		spec, err = loader.LoadSwaggerFromData([]byte(file))
	}
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

	var endpoints []MethodEndpoints
	for path, ep := range spec.Paths {
		endpoints = append(endpoints, o.buildEndpoint(path, ep)...)
	}
	o.types.Sort()

	result := &bytes.Buffer{}
	err = newWriter(result, o.logger).Write(o.buildSyslInfo(spec), o.types, endpoints...)

	return result.String(), err
}

func (o *openapiv3) buildSyslInfo(spec *openapi3.Swagger) SyslInfo {
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
	case "array":
		return attrsForArray(schema)
	case "object":
	case "string":
		return attrsForString(schema)
	case "integer", "int":
	}
	return nil
}

func typeNameFromSchemaRef(ref *openapi3.SchemaRef) string {
	if strings.HasPrefix(ref.Ref, openapiv3DefinitionPrefix) {
		return convertToSyslSafe(strings.TrimPrefix(ref.Ref, openapiv3DefinitionPrefix))
	}
	switch ref.Value.Type {
	case "array":
		return typeNameFromSchemaRef(ref.Value.Items)
	case "object":
		return "/* FIXME */"
	case "boolean":
		return "bool"
	case "string", "integer", "number":
		return mapSwaggerTypeAndFormatToType(ref.Value.Type, ref.Value.Format, logrus.StandardLogger())
	default:
		return convertToSyslSafe(ref.Value.Type)
	}
}

func nameOnlyType(name string) Type {
	return &SyslBuiltIn{name: name}
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

func buildField(name string, prop *openapi3.SchemaRef) Field {
	f := Field{Name: name}
	typeName := typeNameFromSchemaRef(prop)
	f.Type = nameOnlyType(convertToSyslSafe(typeName))
	if prop.Value.Type == "array" {
		f.SizeSpec = makeSizeSpec(prop.Value.MinItems, prop.Value.MaxItems)
		f.Type = &Array{
			name:  "",
			Items: nameOnlyType(typeName),
			Attrs: nil,
		}
	} else if typeName == "string" {
		f.SizeSpec = makeSizeSpec(prop.Value.MinLength, prop.Value.MaxLength)
		if r := prop.Value.Pattern; r != "" {
			f.Attributes = append(f.Attributes, fmt.Sprintf(`regex="%s"`, r))
		}
		if e := prop.Value.Enum; len(e) != 0 {
			var vals []string
			for _, opt := range e {
				vals = append(vals, fmt.Sprintf(`"%s"`, opt))
			}
			sort.Strings(vals)
			f.Attributes = append(f.Attributes, fmt.Sprintf(`enum_values=[%s]`, strings.Join(vals, ", ")))
		}
	}
	return f
}

func (o *openapiv3) loadTypeSchema(name string, schema *openapi3.Schema) Type {
	name = convertToSyslSafe(name)
	switch schema.Type {
	case "array":
		items := o.loadTypeSchema("", schema.Items.Value)
		return &Array{name: name, Items: items, Attrs: getAttrs(schema)}
	case "object":
		obj := &StandardType{
			name:       convertToSyslSafe(name),
			Properties: nil,
			Attributes: getAttrs(schema),
		}

		for name, prop := range schema.Properties {
			f := buildField(name, prop)
			f.Optional = !contains(name, schema.Required)
			obj.Properties = append(obj.Properties, f)
		}
		sortProperties(obj.Properties)
		return obj
	default:
		o.logger.Warnf("unknown scheme type: %s", schema.Type)
		return NewStringAlias(name)
	}
	return nil
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
			for mediatype, obj := range req.Value.Content {
				param := Param{In: "body"}
				param.Field = buildField(typeNameFromSchemaRef(obj.Schema)+"Request", obj.Schema)
				if mediatype != "" {
					param.Attributes = append(param.Attributes, fmt.Sprintf(`mediatype="%s"`, mediatype))
				}
				ep.Params.Add(param)
			}
		}
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
			for mediaType, val := range resp.Value.Content {
				tname := getSyslSafeName(typeNameFromSchemaRef(val.Schema))
				t := nameOnlyType(tname)
				if val.Schema.Value.Type == "array" {
					t = &Array{Items: t}
				}
				f := Field{
					Name:       tname,
					Type:       t,
					Attributes: []string{fmt.Sprintf("mediatype=\"%s\"", mediaType)},
				}
				respType.Properties = append(respType.Properties, f)
			}
			for name, _ := range resp.Value.Headers {
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
		p := Param{
			Field: buildField(item.Value.Name, item.Value.Schema),
			In:    item.Value.In,
		}
		p.Optional = !item.Value.Required
		out.Add(p)
	}
	return out
}
