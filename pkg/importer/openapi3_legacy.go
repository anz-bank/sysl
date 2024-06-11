// TODO: This file contains the old importer for openapi3. It is kept here
// because openapi2 importer still uses this importer. The arrai importer for
// openapi3 still has some gaps to fully replace this importer (mainly handling
// references and reusable types for params, request bodies, etc). Remove this
// when the arrai importer can work with existing tests in openapi2 importer.

package importer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/utils"
)

type OpenAPI3Importer struct {
	appName  string
	pkg      string
	basePath string // Used for Swagger2.0 files which have a basePath field
	logger   *logrus.Logger
	fs       afero.Fs
	types    TypeList

	nameStack []string

	// to check for circular references when loading schemas.
	refMap map[string]bool
}

func NewLegacyOpenAPIV3Importer(logger *logrus.Logger, fs afero.Fs) Importer {
	return &OpenAPI3Importer{
		logger: logger,
		fs:     fs,
	}
}

func (o *OpenAPI3Importer) LoadFile(path string) (string, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return o.Load(string(bs))
}

func (o *OpenAPI3Importer) Load(file string) (string, error) {
	loader := NewOpenAPI3Loader(o.logger, o.fs)
	spec, err := loader.LoadFromFile(file)
	if err != nil {
		spec, _ = loader.LoadFromData([]byte(file)) // nolint: errcheck
	}
	return o.convertSpec(spec)
}

// Configure allows the imported Sysl application name, package and import directories to be specified.
func (o *OpenAPI3Importer) Configure(arg *ImporterArg) (Importer, error) {
	if arg.AppName == "" {
		return nil, errors.New("application name not provided")
	}
	o.appName = arg.AppName
	o.pkg = arg.PackageName
	return o, nil
}

func NewOpenAPI3Loader(logger *logrus.Logger, fs afero.Fs) *openapi3.Loader {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(
		loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		if url.Host == "" && url.Scheme == "" {
			logger.Debugf("Loading OpenAPI3 ref: %s", url)

			data, err := afero.ReadFile(fs, pathFromURL(url))
			if err != nil {
				return nil, err
			}

			return data, err
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

func (o *OpenAPI3Importer) convertSpec(spec *openapi3.T) (string, error) {
	// Convert types
	o.types = TypeList{}
	for name, ref := range spec.Components.Schemas {
		if _, found := o.types.Find(name); !found {
			if ref.Value == nil {
				o.types.Add(NewStringAlias(name))
			} else {
				t, err := o.loadTypeSchema(name, ref.Value)
				if err != nil {
					return "", err
				}
				o.types.Add(t)
			}
		}
	}

	// Convert endpoints
	endpoints := make(map[string][]Endpoint, len(methodDisplayOrder))
	for _, k := range methodDisplayOrder {
		endpoints[k] = nil
	}
	for path, ep := range spec.Paths.Map() {
		meps, err := o.buildEndpoint(path, ep)
		if err != nil {
			return "", err
		}
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

func (o *OpenAPI3Importer) buildSyslInfo(spec *openapi3.T, basepath string) SyslInfo {
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
const openapiv2DefinitionPrefix = "#/definitions/"

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
	switch {
	case schema.Type.Is(openapi3.TypeArray):
		return attrsForArray(schema)
	case schema.Type.Is(openapi3.TypeObject):
	case schema.Type.Is(openapi3.TypeString):
		return attrsForString(schema)
	case schema.Type.Is(openapi3.TypeInteger), schema.Type.Is("int"):
	}
	return nil
}

func typeNameFromSchemaRef(ref *openapi3.SchemaRef) string {
	if idx := strings.Index(ref.Ref, openapiv3DefinitionPrefix); idx >= 0 {
		return getSyslSafeName(strings.TrimPrefix(ref.Ref[idx:], openapiv3DefinitionPrefix))
	}
	if idx := strings.Index(ref.Ref, openapiv2DefinitionPrefix); idx >= 0 {
		return getSyslSafeName(strings.TrimPrefix(ref.Ref[idx:], openapiv2DefinitionPrefix))
	}
	switch {
	case ref.Value.Type.Is(openapi3.TypeArray):
		// if no Items then default to object
		if ref.Value.Items == nil {
			return OpenAPI_OBJECT
		}
		return typeNameFromSchemaRef(ref.Value.Items)
	case ref.Value.Type.Is(openapi3.TypeObject), ref.Value.Type.Is(OpenAPI_EMPTY), ref.Value.Type == nil:
		return OpenAPI_OBJECT
	case ref.Value.Type.Is(openapi3.TypeBoolean):
		return syslutil.Type_BOOL
	case ref.Value.Type.Is(openapi3.TypeString),
		ref.Value.Type.Is(openapi3.TypeInteger),
		ref.Value.Type.Is(openapi3.TypeNumber):
		return mapOpenAPITypeAndFormatToType(ref.Value.Type.Slice()[0], ref.Value.Format, logrus.StandardLogger())
	default:
		refValueType := OpenAPI_EMPTY
		if len(ref.Value.Type.Slice()) > 0 {
			refValueType = ref.Value.Type.Slice()[0]
		}
		return getSyslSafeName(refValueType)
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

	if _, ok := t.(*Array); !ok && ref.Value.Type.Is(openapi3.TypeArray) {
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

func (o *OpenAPI3Importer) buildField(name string, prop *openapi3.SchemaRef) (Field, error) {
	isArray := prop.Value.Type.Is(openapi3.TypeArray)

	// If type is array, but Items is nil then default to an array of object
	if isArray && prop.Value.Items == nil {
		prop.Value.Items = openapi3.NewSchemaRef("", openapi3.NewObjectSchema())
	}

	f := Field{Name: name}
	typeName := typeNameFromSchemaRef(prop)

	if prop.Ref != "" {
		f.Type = nameOnlyType(typeName)
		return f, nil
	}

	defer o.pushName(name)()

	if isArray && prop.Value.Items.Ref != "" {
		f.Type = &Array{Items: nameOnlyType(typeNameFromSchemaRef(prop.Value.Items))}
		// f.SizeSpec = makeSizeSpec(prop.Value.MinItems, prop.Value.MaxItems)
		return f, nil
	}

	f.Type = o.typeAliasForSchema(prop)
	switch typeName {
	// possible types: string, int, bool, object, date, float
	case OpenAPI_OBJECT:
		if isArray {
			prop = prop.Value.Items
		}
		ns := o.nameStack
		o.nameStack = nil
		defer func() { o.nameStack = ns }()
		t, err := o.loadTypeSchema(strings.Join(ns, "_"), prop.Value)
		if err != nil {
			return Field{}, err
		}
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

	return f, nil
}

func attrsForStrings(schema *openapi3.Schema) []string {
	var attrs []string
	if r := schema.Pattern; r != "" {
		attrs = append(attrs, fmt.Sprintf(`regex="%s"`, getSyslSafeURI(r)))
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

func (o *OpenAPI3Importer) loadTypeSchema(name string, schema *openapi3.Schema) (_ Type, err error) {
	if o.refMap == nil {
		o.refMap = make(map[string]bool)
	}
	name = getSyslSafeName(name)
	defer o.pushName(name)()
	setDefined := func(refName string) { o.refMap[refName] = true }
	switch {
	case schema.Type.Is(openapi3.TypeArray):
		var items Type
		if childName := typeNameFromSchemaRef(schema.Items); childName == OpenAPI_OBJECT {
			defer o.pushName("obj")()
			if o.isCircular(schema.Items) {
				return nil, errCircularType(o.nameStack)
			}
			if schema.Items.Ref != "" {
				o.refMap[schema.Items.Ref] = false
			}
			defer setDefined(schema.Items.Ref)
			items, err = o.loadTypeSchema(name+"_obj", schema.Items.Value)
			if err != nil {
				return nil, err
			}
			o.types.Add(items)
		} else {
			items = o.typeAliasForSchema(schema.Items)
		}

		return &Array{baseType: baseType{name: name, attrs: getAttrs(schema)}, Items: items}, nil
	case schema.Type.Is(openapi3.TypeObject), schema.Type.Is(OpenAPI_EMPTY), schema.Type == nil:
		obj := &StandardType{
			baseType: baseType{name: name, attrs: getAttrs(schema)},
		}

		if len(schema.OneOf) != 0 {
			var fields FieldList
			for _, subSchema := range schema.OneOf {
				field, err := o.buildField(typeNameFromSchemaRef(subSchema), subSchema)
				if err != nil {
					return nil, err
				}
				fields = append(fields, field)
			}
			return &Union{
				baseType: baseType{name: name},
				Options:  fields,
			}, nil
		}

		// Removing this as it breaks an import file.
		// AllOf means this object is composed of all of the sub-schemas (and potentially additional properties)
		for _, subschema := range schema.AllOf {
			if o.isCircular(subschema) {
				return nil, errCircularType(o.nameStack)
			}
			if subschema.Ref != "" {
				o.refMap[subschema.Ref] = false
			}
			defer setDefined(subschema.Ref)

			subType, err := o.loadTypeSchema("", subschema.Value)
			if err != nil {
				return nil, err
			}

			if subObj, ok := subType.(*StandardType); ok {
				if len(obj.Properties) == 0 {
					obj.Properties = subObj.Properties
				} else {
					for _, subProp := range subObj.Properties {
						idx := slices.IndexFunc(obj.Properties, func(f Field) bool { return f.Name == subProp.Name })
						switch {
						case idx < 0:
							obj.Properties = append(obj.Properties, subProp)
						case reflect.DeepEqual(obj.Properties[idx], subProp):
							// just ignore duplicates that are identical
						default:
							// warn, but ignore
							o.logger.Warnf("%s contains a duplicate field: '%s', ignoring the second definition", name, subProp.Name)
						}
					}
				}
			}
		}

		for fname, prop := range schema.Properties {
			f, err := o.buildField(fname, prop)
			if err != nil {
				return nil, err
			}
			f.Optional = !utils.Contains(fname, schema.Required)
			obj.Properties = append(obj.Properties, f)
		}
		if len(obj.Properties) == 0 {
			return NewStringAlias(name), nil
		}
		if err = obj.SortProperties(); err != nil {
			return nil, err
		}
		return obj, nil
	default:
		if schema.Type.Is(openapi3.TypeString) && schema.Enum != nil {
			return &Enum{baseType{name: name, attrs: attrsForStrings(schema)}}, nil
		}
		schemaType := OpenAPI_EMPTY
		if len(schema.Type.Slice()) > 0 {
			schemaType = schema.Type.Slice()[0]
		}
		if t, ok := checkBuiltInTypes(mapOpenAPITypeAndFormatToType(schemaType, schema.Format, o.logger)); ok {
			return &Alias{baseType: baseType{name: name}, Target: t}, nil
		}
		o.logger.Warnf("unknown scheme type: %s", schema.Type)
		return NewStringAlias(name), nil
	}
}

func (o *OpenAPI3Importer) isCircular(ref *openapi3.SchemaRef) bool {
	if o.refMap == nil || ref.Ref == "" {
		return false
	}
	t, visited := o.refMap[ref.Ref]
	// if it is visited but the type is false, type is still being defined and that means it is loading a circular type.
	return visited && !t
}

func errCircularType(nameStack []string) error {
	stack := make([]string, 0, len(nameStack))
	for _, n := range nameStack {
		if n != "" {
			stack = append(stack, n)
		}
	}
	return fmt.Errorf("circular reference detected for type: %s", strings.Join(stack, "."))
}

func (o *OpenAPI3Importer) buildEndpoint(path string, item *openapi3.PathItem) ([]MethodEndpoints, error) {
	var res []MethodEndpoints
	ops := map[string]*openapi3.Operation{
		syslutil.Method_GET:    item.Get,
		syslutil.Method_PUT:    item.Put,
		syslutil.Method_POST:   item.Post,
		syslutil.Method_DELETE: item.Delete,
		syslutil.Method_PATCH:  item.Patch,
	}

	commonParams, err := o.buildParams(item.Parameters)
	if err != nil {
		return nil, err
	}

	for method, op := range ops {
		if op == nil {
			continue
		}

		me := MethodEndpoints{Method: method}
		params, err := o.buildParams(op.Parameters)
		if err != nil {
			return nil, err
		}

		ep := &Endpoint{
			Path:        getSyslSafeURI(path),
			Description: op.Description,
			Params:      commonParams.Extend(params),
			Responses:   nil,
		}

		err = o.buildRequests(op.RequestBody, ep)
		if err != nil {
			return nil, err
		}

		for statusCode, resp := range op.Responses.Map() {
			err = o.buildResponses(statusCode, resp, method, path, op, ep)
			if err != nil {
				return nil, err
			}
		}

		me.Endpoints = append(me.Endpoints, *ep)
		res = append(res, me)
	}
	return res, nil
}

func (o *OpenAPI3Importer) buildRequests(req *openapi3.RequestBodyRef, ep *Endpoint) error {
	if req == nil {
		return nil
	}
	if ep == nil {
		return errors.New("nil endpoint")
	}

	fields := make(map[string]map[string]*openapi3.MediaType)
	for mediaType, obj := range req.Value.Content {
		schema := obj.Schema
		tname := typeNameFromSchemaRef(schema)
		if _, ok := fields[tname]; !ok {
			fields[tname] = make(map[string]*openapi3.MediaType)
		}
		fields[tname][mediaType] = obj
	}

	for _, content := range fields {
		mtType := mtReq
		if len(content) > 1 {
			mtType = mtMultiReq
		}
		for mediaType, obj := range content {
			field, err := o.fieldForMediaType(mediaType, obj, mtType)
			if err != nil {
				return err
			}
			ep.Params.Add(Param{In: "body", Field: field})
		}
	}
	return nil
}

func (o *OpenAPI3Importer) buildResponses(
	statusCode string,
	resp *openapi3.ResponseRef,
	method string,
	path string,
	op *openapi3.Operation,
	ep *Endpoint,
) error {
	supportedCode := regexp.MustCompile("^ok|error|[1-5][0-9][0-9]$")
	errType := regexp.MustCompile("^Error|error$")
	typePrefix := getSyslSafeURI(convertToSyslSafe(cleanEndpointPath(path))) + "_"
	text := statusCode

	respType := &StandardType{
		baseType: baseType{
			name: typePrefix + text,
		},
		Properties: FieldList{},
	}

	fields := make(map[string]map[string]*openapi3.MediaType)
	for mediaType, obj := range resp.Value.Content {
		schema := obj.Schema
		tname := typeNameFromSchemaRef(schema)
		if _, ok := fields[tname]; !ok {
			fields[tname] = make(map[string]*openapi3.MediaType)
		}
		fields[tname][mediaType] = obj
	}

	for _, content := range fields {
		mtType := mtResp
		if len(content) > 1 {
			mtType = mtMultiResp
		}
		for mediaType, obj := range content {
			f, err := o.fieldForMediaType(mediaType, obj, mtType)
			if f.Type.Name() == OpenAPI_OBJECT {
				validOperationID := regexp.MustCompile("^[a-zA-Z_]+$")
				if op.OperationID != "" && validOperationID.MatchString(op.OperationID) {
					f.Type.SetName(fmt.Sprintf("%s_%s_%s", method, op.OperationID, statusCode))
				} else {
					f.Type.SetName(fmt.Sprintf("%s_%s", method, respType.Name()))
				}
			}
			if err != nil {
				return err
			}
			respType.Properties = append(respType.Properties, f)
		}
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
		if err := respType.SortProperties(); err != nil {
			return err
		}
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
	return nil
}

func (o *OpenAPI3Importer) buildParams(params openapi3.Parameters) (Parameters, error) {
	var out Parameters
	for _, item := range params {
		name := item.Value.Name
		if item.Value.In == "query" {
			name = convertToSyslSafe(name)
		}
		field, err := o.buildField(name, item.Value.Schema)
		if err != nil {
			return out, err
		}
		p := Param{
			Field: field,
			In:    item.Value.In,
		}
		// Avoid putting sequences into the params
		if a, ok := p.Field.Type.(*Array); ok {
			p.Field.Type = o.types.AddAndRet(&Alias{baseType: baseType{name: item.Value.Name}, Target: a})
		}
		p.Optional = !item.Value.Required
		out.Add(p)
	}
	return out, nil
}

type mediaTypeFieldType uint8

const (
	mtReq mediaTypeFieldType = iota
	mtResp
	mtMultiReq
	mtMultiResp
)

func (o *OpenAPI3Importer) fieldForMediaType(
	mediatype string,
	mediaObj *openapi3.MediaType,
	mtType mediaTypeFieldType,
) (Field, error) {
	schema := mediaObj.Schema
	tname := typeNameFromSchemaRef(schema)

	medianame, typeSuffix := "", ""
	if mtType == mtMultiReq || mtType == mtMultiResp {
		medianame = utils.ToCamel(cleanMediaType(mediatype))
	}
	if mtType == mtReq || mtType == mtMultiReq {
		typeSuffix = "Request"
	}
	field, err := o.buildField(tname+medianame+typeSuffix, schema)
	if err != nil {
		return Field{}, err
	}

	if _, found := o.types.Find(tname); !found {
		t, err := o.loadTypeSchema(tname, schema.Value)
		if err != nil {
			return Field{}, err
		}
		o.types.Add(t)
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

	return field, nil
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
		// sysl calls MustUnescape on return statements, so encode '%' chars
		v = strings.ReplaceAll(v, "%", "%25")
		if b.Len() > len("examples=[") {
			b.Write([]byte{','})
		}
		b.Write([]byte(fmt.Sprintf("[\"%s\",%s]", k, v)))
	}

	b.Write([]byte{']'})

	return b.String()
}
