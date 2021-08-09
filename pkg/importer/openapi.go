package importer

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

// OpenAPI Data Types: https://swagger.io/docs/specification/data-models/data-types/
// nolint:revive,stylecheck
const (
	OpenAPI_EMPTY   = ""
	OpenAPI_STRING  = "string"
	OpenAPI_OBJECT  = "object"
	OpenAPI_ARRAY   = "array"
	OpenAPI_BOOLEAN = "boolean"
	OpenAPI_INTEGER = "integer"
	OpenAPI_NUMBER  = "number"
)

// OpenAPI string formats: https://swagger.io/docs/specification/data-models/data-types/ -> String Formats
// nolint:gochecknoglobals,revive,stylecheck
var OpenAPIFormats = []OpenAPIFormat{
	OpenAPIFormat_INT32,
	OpenAPIFormat_INT64,
	OpenAPIFormat_FLOAT,
	OpenAPIFormat_DOUBLE,
	OpenAPIFormat_DATE,
	OpenAPIFormat_DATETIME,
	OpenAPIFormat_BYTE,
	OpenAPIFormat_BINARY,
	OpenAPIFormat_UUID,
}

type OpenAPIFormat = string

// nolint:revive,stylecheck
const (
	OpenAPIFormat_INT32    OpenAPIFormat = "int32"
	OpenAPIFormat_INT64    OpenAPIFormat = "int64"
	OpenAPIFormat_FLOAT    OpenAPIFormat = "float"
	OpenAPIFormat_DOUBLE   OpenAPIFormat = "double"
	OpenAPIFormat_DATE     OpenAPIFormat = "date"
	OpenAPIFormat_DATETIME OpenAPIFormat = "date-time"
	OpenAPIFormat_BYTE     OpenAPIFormat = "byte"
	OpenAPIFormat_BINARY   OpenAPIFormat = "binary"
	OpenAPIFormat_UUID     OpenAPIFormat = "uuid"
	OpenAPIFormat_URI      OpenAPIFormat = "uri"
)

func mapOpenAPITypeAndFormatToType(typeName, format string, logger *logrus.Logger) string {
	typeName = strings.ToLower(typeName)
	format = strings.ToLower(format)

	// {openapi_type: {openapi_format: sysl_type}}
	conversions := map[string]map[string]string{
		OpenAPI_STRING: {
			"":                     syslutil.Type_STRING,
			OpenAPIFormat_DATE:     syslutil.Type_DATE,
			OpenAPIFormat_DATETIME: syslutil.Type_DATETIME,
			OpenAPIFormat_BYTE:     syslutil.Type_BYTES,
			OpenAPIFormat_BINARY:   syslutil.Type_BYTES,
			OpenAPIFormat_UUID:     syslutil.Type_UUID,
			OpenAPIFormat_URI:      syslutil.Type_STRING,
		},
		OpenAPI_INTEGER: {
			"":                  syslutil.Type_INT,
			OpenAPIFormat_INT32: syslutil.Type_INT32,
			OpenAPIFormat_INT64: syslutil.Type_INT64,
		},
		OpenAPI_NUMBER: {
			"":                   syslutil.Type_FLOAT,
			OpenAPIFormat_DOUBLE: syslutil.Type_FLOAT,
			OpenAPIFormat_FLOAT:  syslutil.Type_FLOAT,
		},
	}

	if formatMap, ok := conversions[typeName]; ok {
		if result, ok := formatMap[format]; ok {
			return result
		}
		logger.Debugf("Unhandled (type, format) -> (%s, %s), ignoring...\n", typeName, format)
		return mapOpenAPITypeAndFormatToType(typeName, "", logger)
	}

	return typeName
}
