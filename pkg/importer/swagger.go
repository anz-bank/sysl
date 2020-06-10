package importer

import (
	"encoding/json"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/sirupsen/logrus"
)

func LoadSwaggerText(args OutputData, text string, logger *logrus.Logger) (out string, err error) {
	openapiv3, basePath, err := convertToOpenapiv3([]byte(text))
	if err != nil {
		return "", err
	}
	return importOpenAPI(args, openapiv3, logger, basePath)
}

func convertToOpenapiv3(data []byte) (*openapi3.Swagger, string, error) {
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
