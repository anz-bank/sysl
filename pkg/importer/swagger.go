package importer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/sirupsen/logrus"
)

func LoadSwaggerText(args OutputData, oas2spec string, logger *logrus.Logger) (out string, err error) {
	oas3spec, basePath, err := convertToOpenAPI3([]byte(oas2spec))
	if err != nil {
		return "", err
	}
	importer := MakeOpenAPI3Importer(logger, basePath, "")
	importer.WithAppName(args.AppName).WithPackage(args.Package)
	importer.spec = oas3spec
	return importer.Parse()
}

func MakeOpenAPI2Importer(logger *logrus.Logger, basePath string, filePath string) *OpenAPI2Importer {
	return &OpenAPI2Importer{OpenAPI3Importer: &OpenAPI3Importer{
		logger:            logger,
		externalSpecs:     make(map[string]*OpenAPI3Importer),
		types:             TypeList{},
		intermediateTypes: TypeList{},
		basePath:          basePath,
		swaggerRoot:       filePath,
	}}
}

type OpenAPI2Importer struct {
	openAPI2Spec string
	*OpenAPI3Importer
}

func (l *OpenAPI2Importer) Load(oas2spec string) (string, error) {
	oas3spec, basePath, err := convertToOpenAPI3([]byte(oas2spec))
	if err != nil {
		return "", fmt.Errorf("error converting openapi 2:%w", err)
	}
	l.openAPI2Spec = oas2spec
	l.basePath = basePath
	l.spec = oas3spec
	return l.Parse()
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
func convertToOpenAPI3(data []byte) (*openapi3.Swagger, string, error) {
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
