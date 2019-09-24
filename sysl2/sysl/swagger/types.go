package swagger

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-openapi/spec"
)

type Type struct {
	Name       string
	Properties []Field
	isAlias    bool
	isArray    bool
	isEnum     bool
}

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

func (t *Type) AddProperties(props map[string]spec.Schema, requiredProps []string, data *typeData) {
	keys := []string{}
	fields := map[string]Field{}
	for pname, prop := range props {
		isArray, refType := checkArrayType(prop, data.logger)
		if !isArray {
			refType = findReferencedType(prop, data.logger)
		}
		propType, found := data.knownTypes.Find(refType)
		if !found {
			if ref, ok := data.doc.Definitions[refType]; ok {
				propType = createTypeFromSchema(refType, &ref, data)
			}
		}
		if refType == "object" {
			propType = Type{
				Name:    fmt.Sprintf("%s_%s_obj", t.Name, pname),
				isAlias: true,
			}
			data.knownTypes = append(data.knownTypes, propType)
		}
		if isArray {
			propType = Type{
				isArray:    true,
				Properties: FieldList{Field{Type: propType}},
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
		t.Properties = append(t.Properties, fields[k])
	}
}

func mapSwaggerTypeAndFormatToType(typeName, format string, logger *logrus.Logger) string {
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
		return strings.Compare(data.knownTypes[i].Name, data.knownTypes[j].Name) < 0
	})
	return data.knownTypes
}

func (t TypeList) Find(name string) (Type, bool) {

	if builtin, ok := checkBuiltInTypes(name); ok {
		return builtin, ok
	}

	for _, n := range t {
		if n.Name == name {
			return n, true
		}
	}
	return Type{}, false
}

func checkBuiltInTypes(name string) (Type, bool) {

	if syslType, ok := swaggerToSyslMappings[name]; ok {
		return Type{Name: syslType}, true
	}

	if contains(name, builtIns) {
		return Type{Name: name}, true
	}
	return Type{}, false
}

func checkArrayType(schema spec.Schema, logger *logrus.Logger) (bool, string) {
	if len(schema.Type) == 1 {
		typeName := schema.Type[0]
		if typeName == "array" && schema.Items != nil {
			return true, findReferencedType(*schema.Items.Schema, logger)
		}
	}

	return false, ""
}

func findReferencedType(schema spec.Schema, logger *logrus.Logger) string {
	if len(schema.Type) == 1 {
		if isArray, items := checkArrayType(schema, logger); isArray {
			return "sequence of " + items
		}
		return mapSwaggerTypeAndFormatToType(schema.Type[0], schema.Format, logger)
	} else if len(schema.Type) == 0 && schema.Items != nil {
		return findReferencedType(*schema.Items.Schema, logger)
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
	item := Type{
		Name: name,
	}
	if len(schema.Properties) == 0 {
		item.isAlias = true
		if isArray, _ := checkArrayType(*schema, data.logger); isArray {
			item.isArray = true
			nested := Type{
				Name:    name + "_obj",
				isAlias: true,
			}
			data.knownTypes = append(data.knownTypes, nested)
			item.Properties = []Field{
				{
					Type: nested,
				},
			}
		} else if len(schema.Enum) > 0 {
			item.isEnum = true
		} else if refType := findReferencedType(*schema, data.logger); refType != "" {
			log.Printf("WARNING: swagger type '%s' is malformed\n", name)
			t, found := data.knownTypes.Find(refType)
			if !found {
				if ref, ok := data.doc.Definitions[refType]; ok {
					t = createTypeFromSchema(refType, &ref, data)
				} else {
					t, _ = data.knownTypes.Find("string")
				}
			}
			item.Properties = []Field{
				{
					Type: t,
				},
			}
		}
	} else {
		item.AddProperties(schema.Properties, schema.Required, data)
	}

	data.knownTypes = append(data.knownTypes, item)

	return item
}
