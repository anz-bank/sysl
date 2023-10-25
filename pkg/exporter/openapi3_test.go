//nolint:lll
package exporter

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/tidwall/gjson"
)

func readSyslModule(filename string) (*sysl.Module, error) {
	mod, err := parse.NewParser().ParseFromFs(filename,
		syslutil.NewChrootFs(afero.NewOsFs(), "."))
	return mod, err
}

func TestExportWithFile(t *testing.T) {
	t.Parallel()

	mod, err := readSyslModule("./test-data/openapi3/petstore.sysl")
	assert.NoError(t, err)
	mapper := syslwrapper.MakeAppMapper(mod)
	mapper.IndexTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)

	exporter := MakeOpenAPI3Exporter(simpleApps, &logrus.Logger{})
	err = exporter.Export()
	assert.NoError(t, err)
	outputSpecJSON, err := exporter.SerializeOutput("Swagger Petstore", "json")
	assert.NoError(t, err)

	errorType := gjson.Get(string(outputSpecJSON), "components.schemas.Error.type")
	errorRequired := gjson.Get(string(outputSpecJSON), "components.schemas.Error.required")
	for _, val := range errorRequired.Array() {
		assert.True(t, val.Str == "message" || val.Str == "code")
	}
	errorPropertiesCode := gjson.Get(string(outputSpecJSON), "components.schemas.Error.properties.code.type")
	errorPropertiesCodeFormat := gjson.Get(string(outputSpecJSON), "components.schemas.Error.properties.code.format")
	errorPropertiesMessage := gjson.Get(string(outputSpecJSON), "components.schemas.Error.properties.message.type")
	assert.Equal(t, "object", errorType.Str)
	assert.Equal(t, "integer", errorPropertiesCode.Str)
	assert.Equal(t, "int64", errorPropertiesCodeFormat.Str)
	assert.Equal(t, "string", errorPropertiesMessage.Str)

	newPetNameDescription := gjson.Get(string(outputSpecJSON), "components.schemas.NewPet.properties.name.description")
	assert.Equal(t, "The Pets Name", newPetNameDescription.Str)

	getPetsSuccessResponse := gjson.Get(string(outputSpecJSON), "paths./pets.get.responses.200.content.application/json.schema.$ref")
	getPetsErrorResponse := gjson.Get(string(outputSpecJSON), "paths./pets.get.responses.default.content.application/json.schema.$ref")
	postPetsRequired := gjson.Get(string(outputSpecJSON), "paths./pets.post.requestBody.required")
	assert.Equal(t, "#/components/schemas/NewPetResponse", getPetsSuccessResponse.Str)
	assert.Equal(t, "#/components/schemas/Error", getPetsErrorResponse.Str)
	assert.True(t, postPetsRequired.Bool())

	xVer := gjson.Get(string(outputSpecJSON), "info.x-version")
	assert.Equal(t, "1.0.1", xVer.Str)
	pathXParam := gjson.Get(string(outputSpecJSON), "paths./pets.get.x-operation-param")
	assert.Equal(t, "xOperationParam", pathXParam.Str)
}
func TestExport(t *testing.T) {
	t.Parallel()

	type1 := syslwrapper.MakeMap(
		syslwrapper.MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"}),
		syslwrapper.MakePrimitive("string"))
	type2 := syslwrapper.MakePrimitive("string")
	param1 := syslwrapper.MakeParam("Login", type1)
	var app1 = syslwrapper.MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = syslwrapper.MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": app1,
			"app2": app2,
		},
	}

	mapper := syslwrapper.MakeAppMapper(mod)
	mapper.IndexTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)

	exporter := MakeOpenAPI3Exporter(simpleApps, &logrus.Logger{})
	err = exporter.Export()
	assert.NoError(t, err)

	outputSpec, err := exporter.SerializeOutput("app1", "json")
	assert.NoError(t, err)

	paramName := gjson.Get(string(outputSpec), "paths.testEndpoint.get.parameters.0.name")
	assert.Equal(t, "Login", paramName.Str)
	paramSchema := gjson.Get(string(outputSpec), "paths.testEndpoint.get.parameters.0.schema.type")
	assert.Equal(t, "object", paramSchema.Str)
}

func TestMapOptional(t *testing.T) {
	t.Parallel()

	simpleApps := map[string]*syslwrapper.App{
		"TestApp": {
			Name: "TestApp",
			Endpoints: map[string]*syslwrapper.Endpoint{
				"TestEndpoint": {
					Summary: "GetPets",
					Path:    "GET /pets",
					Params: map[string]*syslwrapper.Parameter{
						"limit": {
							In:   "query",
							Name: "limit",
							Type: &syslwrapper.Type{
								Type:     "int",
								Optional: true,
							},
						},
					},
				},
			},
		},
	}

	exporter := MakeOpenAPI3Exporter(simpleApps, &logrus.Logger{})
	err := exporter.Export()
	assert.NoError(t, err)
	limitRequired := exporter.openapi3["TestApp"].Paths["/pets"].Get.Parameters.GetByInAndName("query", "limit").Required
	assert.False(t, limitRequired)
}

func TestMapEnums(t *testing.T) {
	t.Parallel()

	simpleApps := map[string]*syslwrapper.App{
		"TestApp": {
			Name: "TestApp",
			Endpoints: map[string]*syslwrapper.Endpoint{
				"TestEndpoint": {
					Summary: "GetPets",
					Path:    "GET /pets",
					Params: map[string]*syslwrapper.Parameter{
						"limit": {
							In:   "query",
							Name: "limit",
							Type: &syslwrapper.Type{
								Type: "enum",
								Enum: map[int64]string{
									1: "apple",
									2: "orange",
								},
							},
						},
					},
				},
			},
		},
	}

	exporter := MakeOpenAPI3Exporter(simpleApps, &logrus.Logger{})
	err := exporter.Export()
	assert.NoError(t, err)
	outputSpec, err := exporter.SerializeOutput("TestApp", "json")
	assert.NoError(t, err)
	enum1 := gjson.Get(string(outputSpec), "paths./pets.get.parameters.0.schema.enum")
	for _, val := range enum1.Array() {
		assert.True(t, val.Str == "apple" || val.Str == "orange")
	}
}

func TestExtensions(t *testing.T) {
	t.Parallel()

	simpleApps := map[string]*syslwrapper.App{
		"TestApp": {
			Name: "TestApp",
			Attributes: map[string]string{
				"version":   "1.0.0",
				"x-version": "1.0.1",
			},
			Endpoints: map[string]*syslwrapper.Endpoint{
				"TestEndpoint": {
					Summary: "GetPets",
					Path:    "GET /pets",
					Params: map[string]*syslwrapper.Parameter{
						"limit": {
							In:   "query",
							Name: "limit",
							Type: &syslwrapper.Type{
								Type: "int",
							},
						},
					},
					Extensions: map[string]interface{}{
						"x-operation-param": "xOperationParam",
					},
				},
			},
		},
	}

	exporter := MakeOpenAPI3Exporter(simpleApps, &logrus.Logger{})
	err := exporter.Export()
	assert.NoError(t, err)
	outputSpecJSON, err := exporter.SerializeOutput("TestApp", "json")
	assert.NoError(t, err)
	xVer := gjson.Get(string(outputSpecJSON), "info.x-version")
	assert.Equal(t, "1.0.1", xVer.Str)
	pathXParam := gjson.Get(string(outputSpecJSON), "paths./pets.get.x-operation-param")
	assert.Equal(t, "xOperationParam", pathXParam.Str)
}

func TestMakeOpenAPI3Exporter(t *testing.T) {
	t.Parallel()

	exporter := MakeOpenAPI3Exporter(map[string]*syslwrapper.App{}, &logrus.Logger{})
	expected := &OpenAPI3Exporter{
		apps:     map[string]*syslwrapper.App{},
		openapi3: map[string]*openapi3.T{},
		log:      &logrus.Logger{},
	}
	assert.Equal(t, exporter, expected)
}

func TestExportHandlesNamespaces(t *testing.T) {
	t.Parallel()
	appName := "Namespace :: testapp"
	mod, err := readSyslModule("./test-data/openapi3/namespace.sysl")
	assert.NoError(t, err)
	mapper := syslwrapper.MakeAppMapper(mod)
	mapper.IndexTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)

	exporter := MakeOpenAPI3Exporter(simpleApps, &logrus.Logger{})
	err = exporter.Export()
	assert.NoError(t, err)
	outputSpecJSON, err := exporter.SerializeOutput(appName, "json")
	assert.NoError(t, err)
	assert.Contains(t, string(outputSpecJSON), appName)
}
