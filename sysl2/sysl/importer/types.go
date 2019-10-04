package importer

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-openapi/spec"
)

type Type interface {
	Name() string
}

type StandardType struct {
	name       string
	Properties []Field
}

func (s *StandardType) Name() string { return s.name }

type SyslBuiltIn struct {
	name string
}

func (s *SyslBuiltIn) Name() string { return s.name }

type Alias struct {
	name   string
	Target Type
}

func NewStringAlias(name string) Type {
	return &Alias{
		name:   name,
		Target: &SyslBuiltIn{name: "string"},
	}
}

func (s *Alias) Name() string { return s.name }

type Array struct {
	name  string
	Items Type
}

func (s *Array) Name() string { return s.name }

type Enum struct {
	name string
}

func (s *Enum) Name() string { return s.name }

type Field struct {
	Name     string
	Type     Type
	Optional bool
}

type TypeList []Type
type FieldList []Field

// nolint:gochecknoglobals
var builtIns = []string{"int32", "int64", "int", "float", "string", "date", "bool", "decimal", "datetime", "xml"}

// nolint:gochecknoglobals
var swaggerToSyslMappings = map[string]string{
	"boolean": "bool",
	"date":    "date",
}

// nolint:gochecknoglobals
var swaggerFormats = []string{"int32", "int64", "float", "double", "date", "date-time", "byte"}

func IsKeyword(name string) bool {
	for _, kw := range builtIns {
		if name == kw {
			return true
		}
	}
	return false
}

func (s *StandardType) AddProperties(props map[string]spec.Schema, requiredProps []string, data *typeData) {
	keys := []string{}
	fields := map[string]Field{}
	for pname, prop := range props {
		propType, found := data.knownTypes.FindFromSchema(prop, data)
		if !found {
			refType := findReferencedType(prop, data)
			if ref, ok := data.doc.Definitions[refType]; ok {
				propType = createTypeFromSchema(refType, &ref, data)
			}
			if refType == "object" {
				propType = NewStringAlias(fmt.Sprintf("%s_%s_obj", s.Name(), pname))
				data.knownTypes = append(data.knownTypes, propType)
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

func mapSwaggerTypeAndFormatToType(typeName, format string, logger *logrus.Logger) string {
	typeName = strings.ToLower(typeName)
	format = strings.ToLower(format)
	if format != "" && !contains(format, swaggerFormats) {
		logger.Errorf("unknown format '%s' being used, ignoring...\n", format)
		format = ""
	}

	conversions := map[string]map[string]string{
		"string": {
			"":          "string",
			"date":      "date",
			"date-time": "datetime",
			"byte":      "string",
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

func InitTypes(doc *spec.Swagger, logger *logrus.Logger) TypeList {
	data := &typeData{
		doc:        doc,
		knownTypes: TypeList{},
		logger:     logger,
	}
	for name, definition := range doc.Definitions {
		def := definition
		if _, found := data.knownTypes.Find(name); !found {
			_ = createTypeFromSchema(name, &def, data)
		}
	}

	sort.SliceStable(data.knownTypes, func(i, j int) bool {
		return strings.Compare(data.knownTypes[i].Name(), data.knownTypes[j].Name()) < 0
	})
	return data.knownTypes
}

func (t TypeList) Find(name string) (Type, bool) {
	if builtin, ok := checkBuiltInTypes(name); ok {
		return builtin, ok
	}

	for _, n := range t {
		if n.Name() == name {
			return n, true
		}
	}
	return &StandardType{}, false
}

func (t TypeList) FindFromSchema(schema spec.Schema, data *typeData) (Type, bool) {
	if isArrayType(schema) {
		items, found := t.FindFromSchema(*schema.Items.Schema, data)
		if found {
			return &Array{Items: items}, true
		}
	}
	return t.Find(findReferencedType(schema, data))
}

func checkBuiltInTypes(name string) (Type, bool) {
	if syslType, ok := swaggerToSyslMappings[name]; ok {
		return &SyslBuiltIn{name: syslType}, true
	}

	if contains(name, builtIns) {
		return &SyslBuiltIn{name: name}, true
	}
	return &StandardType{}, false
}

func isArrayType(schema spec.Schema) bool {
	if len(schema.Type) == 1 {
		typeName := schema.Type[0]
		if typeName == "array" && schema.Items != nil {
			return true
		}
	}

	return false
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

func contains(needle string, haystack []string) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

type typeData struct {
	doc        *spec.Swagger
	knownTypes TypeList
	logger     *logrus.Logger
}

func createTypeFromSchema(name string, schema *spec.Schema, data *typeData) Type {
	var item Type
	if len(schema.Properties) == 0 {
		if isArrayType(*schema) {
			nested := NewStringAlias(fmt.Sprintf("%s_obj", name))
			data.knownTypes = append(data.knownTypes, nested)
			item = &Array{name: name, Items: nested}
		} else if len(schema.Enum) > 0 {
			item = &Enum{name: name}
		} else if refType := findReferencedType(*schema, data); refType != "" {
			log.Printf("WARNING: swagger type '%s' is malformed\n", name)
			t, found := data.knownTypes.Find(refType)
			if !found {
				if ref, ok := data.doc.Definitions[refType]; ok {
					t = createTypeFromSchema(refType, &ref, data)
				} else {
					t, _ = data.knownTypes.Find("string")
				}
			}
			item = &Alias{
				name:   name,
				Target: t,
			}
		}
	} else {
		st := &StandardType{
			name: name,
		}
		st.AddProperties(schema.Properties, schema.Required, data)
		item = st
	}

	data.knownTypes = append(data.knownTypes, item)

	return item
}
