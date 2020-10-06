package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

// Format represents a format that can be imported into Sysl
type Format struct {
	Name      string   // Name of the format
	Signature string   // This is a string which can be used to uniquely identify the format
	FileExt   []string // The file extension of the format
}

var SYSL = Format{
	Name:      "sysl",
	Signature: "",
	FileExt:   []string{".sysl"},
}

var XSD = Format{
	Name:      "xsd",
	Signature: ``,
	FileExt:   []string{".xsd", ".xml"},
}

var Grammar = Format{
	Name:      "grammar",
	Signature: "",
	FileExt:   []string{".g"},
}

var Avro = Format{
	Name:      "avro",
	Signature: "",
	FileExt:   []string{".avsc"},
}

var SpannerSQL = Format{
	Name:      "spannerSQL",
	Signature: "",
	FileExt:   []string{".sql"},
}

var SpannerSQLDir = Format{
	Name:      "spannerSQLdir",
	Signature: "",
	FileExt:   []string{".up.sql"},
}

// OpenAPI3 is identified by the openapi header. -  The value MUST be "3.x.x".
// For more details refer to https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md#oasDocument
var OpenAPI3 = Format{
	Name:      "openapi3",
	Signature: `openapi:`,
	FileExt:   []string{".yaml", ".json", ".yml"},
}

// Swagger only has 2.0.0 as the single valid format -  The value MUST be "2.0".
// For more details refer to https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#swaggerObject
var Swagger = Format{
	Name:      "swagger",
	Signature: `swagger:`,
	FileExt:   []string{".yaml", ".json", ".yml"},
}

// GuessFileType detects the file based on the filename extension and the file itself.
// It returns the detected format if successful, or an error if not.
// It first tries to match the file extensions before checking the files for signatures such as
// `swagger: "2.0"`.
func GuessFileType(path string, isDir bool, content []byte, validFormats []Format) (*Format, error) {
	if isDir {
		if files, err := ioutil.ReadDir(path); err == nil {
			for _, info := range files {
				if strings.HasSuffix(info.Name(), ".up.sql") || strings.HasSuffix(info.Name(), ".up.ddl") {
					return &SpannerSQLDir, nil
				}
			}
		}
	}

	var matchesExt []Format
	ext := filepath.Ext(path)
	for _, format := range validFormats {
		for _, formatExt := range format.FileExt {
			if formatExt == ext {
				matchesExt = append(matchesExt, format)
				break
			}
		}
	}

	if len(matchesExt) == 1 {
		return &matchesExt[0], nil
	}

	var matchesSignature []Format
	// Convert to yaml so we only need to compare a single format
	if ext == ".json" {
		var err error
		content, err = yaml.JSONToYAML(content)
		if err != nil {
			return nil, fmt.Errorf("error converting spec to yaml for: %s", path)
		}
	}

	for _, format := range matchesExt {
		if strings.Contains(string(content), format.Signature) {
			matchesSignature = append(matchesSignature, format)
		}
	}

	if len(matchesSignature) == 1 {
		return &matchesSignature[0], nil
	}

	// We return an error if the number of matches is less than 0 or greater than 1
	return nil, fmt.Errorf("error detecting input file format for: %s", path)
}
