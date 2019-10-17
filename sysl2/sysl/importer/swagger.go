package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/sirupsen/logrus"

	"github.com/go-openapi/loads"
)

func LoadSwaggerText(args OutputData, text string, logger *logrus.Logger) (out string, err error) {
	doc, err := loads.Analyzed(json.RawMessage(text), "2.0")
	if err != nil {
		logger.Errorf("Failed to load swagger spec: %s\n", err.Error())
		return "", err
	}

	result := &bytes.Buffer{}

	swagger := doc.Spec()
	types := InitSwaggerTypes(swagger, logger)
	globalParams := buildGlobalParams(swagger.Parameters, types, logger)
	endpoints := InitEndpoints(swagger, types, globalParams, logger)
	info := SyslInfo{
		OutputData:  args,
		Title:       "",
		Description: "",
		OtherFields: []string{},
	}
	if swagger.Info != nil {
		info.Title = swagger.Info.Title
		info.Description = swagger.Info.Description
		values := []string{
			"version", swagger.Info.Version,
			"host", swagger.Host,
			"license", "",
			"termsOfService", swagger.Info.TermsOfService}
		for i := 0; i < len(values); i += 2 {
			key := values[i]
			val := values[i+1]
			if val != "" {
				info.OtherFields = append(info.OtherFields, key, val)
			}
		}
	}

	w := newWriter(result, logger)
	if err := w.Write(info, types, swagger.BasePath, endpoints...); err != nil {
		return "", err
	}

	return result.String(), nil
}

func InitSwaggerTypes(doc *spec.Swagger, logger *logrus.Logger) TypeList {
	types := TypeList{}
	// First init the swagger -> sysl mappings
	var swaggerToSyslMappings = map[string]string{
		"boolean": "bool",
		"date":    "date",
	}
	for swaggerName, syslName := range swaggerToSyslMappings {
		types.Add(&ImportedBuiltInAlias{
			name:   swaggerName,
			Target: &SyslBuiltIn{syslName},
		})
	}

	data := &typeData{
		doc:        doc,
		knownTypes: types,
		logger:     logger,
	}
	for name, definition := range doc.Definitions {
		def := definition
		if _, found := data.knownTypes.Find(name); !found {
			_ = createTypeFromSchema(name, &def, data)
		}
	}

	data.knownTypes.Sort()
	return data.knownTypes
}

// nolint:gochecknoglobals
var swaggerFormats = []string{"int32", "int64", "float", "double", "date", "date-time", "byte"}

func mapSwaggerTypeAndFormatToType(typeName, format string, logger *logrus.Logger) string {
	typeName = strings.ToLower(typeName)
	format = strings.ToLower(format)
	if format != "" && !contains(format, swaggerFormats) {
		logger.Errorf("unknown format '%s' being used, ignoring...\n", format)
		format = ""
	}

	conversions := map[string]map[string]string{
		StringTypeName: {
			"":          StringTypeName,
			"date":      "date",
			"date-time": "datetime",
			"byte":      StringTypeName,
		},
		"integer": {
			"":      "int",
			"int32": "int32",
			"int64": "int64",
		},
		"number": {
			"":       "float",
			"double": "float",
			"float":  "float",
		},
	}

	if formatMap, ok := conversions[typeName]; ok {
		if result, ok := formatMap[format]; ok {
			return result
		}
		logger.Warnf("Unhandled (type, format) -> (%s, %s)\n", typeName, format)
		return mapSwaggerTypeAndFormatToType(typeName, "", logger)
	}

	return typeName
}

func findReferencedType(schema spec.Schema, data *typeData) string {
	if len(schema.Type) == 1 {
		return mapSwaggerTypeAndFormatToType(schema.Type[0], schema.Format, data.logger)
	} else if len(schema.Type) == 0 && schema.Items != nil {
		return findReferencedType(*schema.Items.Schema, data)
	}

	if refURL := schema.Ref.GetURL(); refURL != nil {
		return getReferenceFragment(refURL)
	}

	return ""
}

func getReferenceFragment(u *url.URL) string {
	parts := strings.Split(u.Fragment, "/")
	return parts[len(parts)-1]
}

type typeData struct {
	doc        *spec.Swagger
	knownTypes TypeList
	logger     *logrus.Logger
}

func createTypeFromSchema(name string, schema *spec.Schema, data *typeData) Type {
	var item Type
	if len(schema.Properties) == 0 {
		if isSwaggerArrayType(*schema) {
			nested := NewStringAlias(fmt.Sprintf("%s_obj", name))
			data.knownTypes.Add(nested)
			item = &Array{name: name, Items: nested}
		} else if len(schema.Enum) > 0 {
			item = &Enum{name: name}
		} else if refType := findReferencedType(*schema, data); refType != "" {
			data.logger.Warnf("WARNING: swagger type '%s' is malformed\n", name)
			t, found := data.knownTypes.Find(refType)
			if !found {
				if ref, ok := data.doc.Definitions[refType]; ok {
					t = createTypeFromSchema(refType, &ref, data)
				} else {
					t, _ = data.knownTypes.Find(StringTypeName)
				}
			}
			item = &ExternalAlias{
				name:   name,
				Target: t,
			}
		}
	} else {
		st := &StandardType{
			name: name,
		}
		AddSwaggerProperties(st, schema.Properties, schema.Required, data)
		item = st
	}

	data.knownTypes.Add(item)

	return item
}

func (t TypeList) FindFromSchema(schema spec.Schema, data *typeData) (Type, bool) {
	if isSwaggerArrayType(schema) {
		items, found := t.FindFromSchema(*schema.Items.Schema, data)
		if found {
			return &Array{Items: items}, true
		}
	}
	return t.Find(findReferencedType(schema, data))
}

func isSwaggerArrayType(schema spec.Schema) bool {
	if len(schema.Type) == 1 {
		typeName := schema.Type[0]
		if typeName == "array" && schema.Items != nil {
			return true
		}
	}

	return false
}

func AddSwaggerProperties(s *StandardType, props map[string]spec.Schema, requiredProps []string, data *typeData) {
	keys := []string{}
	fields := map[string]Field{}
	for pname, prop := range props {
		propType, found := data.knownTypes.FindFromSchema(prop, data)
		if !found {
			var refType string
			refType = findReferencedType(prop, data)
			if refType == "" || refType == "object" {
				p := prop
				propType = createTypeFromSchema(fmt.Sprintf("%s_%s_obj", s.Name(), pname), &p, data)
			}
			if ref, ok := data.doc.Definitions[refType]; ok {
				propType = createTypeFromSchema(refType, &ref, data)
			}
			if isSwaggerArrayType(prop) {
				refType = findReferencedType(*prop.Items.Schema, data)
				if refType == "object" {
					propType = createTypeFromSchema(fmt.Sprintf("%s_%s_obj", s.Name(), pname), prop.Items.Schema, data)
				} else {
					if ref, ok := data.doc.Definitions[refType]; ok {
						propType = createTypeFromSchema(refType, &ref, data)
					} else {
						data.logger.Errorf("Referenced type %s not found\n", refType)
					}
				}

				propType = &Array{Items: propType}
			}
		}
		f := Field{
			Name:     pname,
			Type:     propType,
			Optional: !contains(pname, requiredProps),
		}
		fields[pname] = f
		keys = append(keys, pname)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	for _, k := range keys {
		s.Properties = append(s.Properties, fields[k])
	}
}
