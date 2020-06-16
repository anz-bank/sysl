package importer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var formatsToTest = []Format{
	XSD,
	OpenAPI3,
	SYSL,
	Swagger,
}

var formatTests = []struct {
	testName     string
	fileName     string
	fileContents []byte
	out          *Format
	err          error
}{
	{"Parses XML ext", ".xml", []byte{}, &XSD, nil},
	{"Parses XSD ext", ".xsd", []byte{}, &XSD, nil},
	{"Parses sysl ext", ".sysl", []byte{}, &SYSL, nil},
	{"Parses openapi3 yaml files", ".yaml", []byte(`openapi: "3.0.0"`), &OpenAPI3, nil},
	{"Parses openapi3 yaml files", ".yaml", []byte(`openapi: "3.0.3"`), &OpenAPI3, nil},
	{"Parses openapi3 json files", ".json", []byte(`{"openapi": "3.0.0"}`), &OpenAPI3, nil},
	{"Parses swagger2 yaml files", ".yaml", []byte(`swagger: "2.0"`), &Swagger, nil},
	{"Parses swagger2 json files", ".json", []byte(`{"swagger": "2.0"}`), &Swagger, nil},
	{"Parses xsd files", ".xsd", []byte(`xml`), &XSD, nil},
	{"Parses xml files", ".xml", []byte(`xml`), &XSD, nil},
	{"Fails for unknown extension", ".abcde", []byte{}, nil, errors.New("error detecting input file format for: ")},
	{"Fails for invalid json", ".json", []byte{}, nil, errors.New("error converting json to yaml for: ")},
	{"Fails for empty string", "", []byte{}, nil, errors.New("error detecting input file format for: ")},
}

//nolint:scopelint
func TestGuessFileType(t *testing.T) {
	for _, tt := range formatTests {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()
			guessedType, err := GuessFileType(tt.fileName, tt.fileContents, formatsToTest)
			assert.Equal(t, tt.err, err)
			if tt.out != nil {
				assert.Equal(t, tt.out.Name, guessedType.Name)
			}
		})
	}
}
