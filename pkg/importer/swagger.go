package importer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/afero"

	"github.com/ghodss/yaml"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/sirupsen/logrus"
)

func MakeOpenAPI2Importer(logger *logrus.Logger, basePath string, filePath string) *OpenAPI2Importer {
	return &OpenAPI2Importer{openapiv3: &openapiv3{
		appName:  "",
		pkg:      "",
		basePath: basePath,
		logger:   logger,
		fs:       afero.NewOsFs(),
		types:    TypeList{},
	},
		filepath: filePath}
}

type OpenAPI2Importer struct {
	openAPI2Spec string
	filepath     string
	*openapiv3
}

func (l *OpenAPI2Importer) Load(oas2spec string) (string, error) {
	u, err := pathToURL(l.filepath)
	if err != nil {
		return "", err
	}
	oas3spec, basePath, err := convertToOpenAPI3([]byte(oas2spec), u, NewOpenAPILoader(l.logger, l.fs))
	if err != nil {
		return "", fmt.Errorf("error converting openapi 2:%w", err)
	}
	l.openAPI2Spec = oas2spec
	l.basePath = basePath
	return l.convertSpec(oas3spec, basePath)
}

// Set the AppName of the imported app
func (l *OpenAPI2Importer) WithAppName(appName string) Importer {
	l.appName = appName
	return l
}

// Set the package attribute of the imported app
func (l *OpenAPI2Importer) WithPackage(pkg string) Importer {
	l.pkg = pkg
	return l
}

// convertToOpenAPI3 takes a swagger spec and converts it to openapi3
func convertToOpenAPI3(data []byte, uri *url.URL, loader *openapi3.SwaggerLoader) (*openapi3.Swagger, string, error) {
	var swagger2 openapi2.Swagger
	jsondata, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal(jsondata, &swagger2)
	if err != nil {
		return nil, "", err
	}

	if len(swagger2.Schemes) == 0 {
		swagger2.Schemes = []string{"https"}
	}
	openapiv3, err := openapi2conv.ToV3Swagger(&swagger2)
	if err != nil {
		return nil, "", err
	}
	// openapi2conv doesnt handle external references correctly so to avoid that we need to serialise the converted
	// v3 doc back to text and manually replace #/components/* references (easier to do it via text instead of walking
	// the whole object tree again
	if oa3text, err := json.Marshal(openapiv3); err == nil {
		replacer := strings.NewReplacer(
			"#/definitions/", "#/components/schemas/",
			"#/responses/", "#/components/responses/",
			"#/parameters/", "#/components/parameters/",
		)
		oa3replaced := replacer.Replace(string(oa3text))
		_ = json.Unmarshal([]byte(oa3replaced), &openapiv3) //nolint: errcheck // we dont care about the error here
	}

	err = loader.ResolveRefsIn(openapiv3, uri)
	if err != nil {
		return nil, "", err
	}

	return openapiv3, swagger2.BasePath, nil
}

// nolint:gochecknoglobals
var swaggerFormats = []string{"int32", "int64", "float", "double", "date", "date-time", "byte", "binary"}

func mapSwaggerTypeAndFormatToType(typeName, format string, logger *logrus.Logger) string {
	typeName = strings.ToLower(typeName)
	format = strings.ToLower(format)
	if format != "" && !contains(format, swaggerFormats) {
		logger.Warnf("unknown format '%s' being used, ignoring...\n", format)
		format = ""
	}

	conversions := map[string]map[string]string{
		StringTypeName: {
			"":          StringTypeName,
			"date":      "date",
			"date-time": "datetime",
			"byte":      StringTypeName,
			"binary":    "bytes",
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
