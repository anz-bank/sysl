package importer

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_OpenAPI3Importer_loadTypeSchema(t *testing.T) {
	t.Parallel()
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	schema, err := loader.LoadSwaggerFromFile("tests/openapi3/one-of.yaml")
	assert.NoError(t, err)
	o := &OpenAPI3Importer{logger: logrus.New()}

	catType := &SyslBuiltIn{
		name: "Cat",
	}
	dogType := &SyslBuiltIn{
		name: "Dog",
	}
	expected := &Union{name: "Pet", Options: []Field{{Name: "Cat", Type: catType}, {Name: "Dog", Type: dogType}}}
	if got := o.loadTypeSchema("Pet", schema.Components.Schemas["Pet"].Value); !reflect.DeepEqual(got, expected) {
		t.Errorf("OpenAPI3Importer.loadTypeSchema() = %v, want %v", got, expected)
	}
}
