package exporter

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslwrapper"
	"github.com/getkin/kin-openapi/openapi3"
	yaml "github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

type OpenAPI3Exporter struct {
	apps     map[string]*syslwrapper.App
	openapi3 map[string]*openapi3.Swagger
	log      *logrus.Logger
}

func MakeOpenAPI3Exporter(apps map[string]*syslwrapper.App, logger *logrus.Logger) *OpenAPI3Exporter {
	return &OpenAPI3Exporter{
		apps:     apps,
		openapi3: make(map[string]*openapi3.Swagger),
		log:      logger,
	}
}

func (s *OpenAPI3Exporter) Export() error {
	for k, v := range s.apps {
		oas3spec, err := s.GenerateOpenAPI3(v)
		if err != nil {
			return err
		}
		s.openapi3[k] = oas3spec
	}
	return nil
}

func (s *OpenAPI3Exporter) SerializeOutput(appName string, mode string) ([]byte, error) {
	spec, ok := s.openapi3[appName]
	if !ok {
		return nil, errors.New("spec not found")
	}
	jsonSpec, err := json.MarshalIndent(spec, "", " ")
	if err != nil {
		return nil, err
	}
	if mode == "json" {
		return jsonSpec, nil
	}
	yamlSpec, err := yaml.JSONToYAML(jsonSpec)
	if err != nil {
		return nil, err
	}
	return yamlSpec, nil
}

func (s *OpenAPI3Exporter) GenerateOpenAPI3(app *syslwrapper.App) (*openapi3.Swagger, error) {
	spec := &openapi3.Swagger{}
	spec.OpenAPI = "3.0.0"
	spec.Info = &openapi3.Info{}
	spec.Info.Title = app.Name
	spec.Info.Version = app.Attributes["version"]
	spec.Info.Description = app.Attributes["description"]
	spec.Info.Contact = &openapi3.Contact{}
	spec.Info.Contact.Name = app.Attributes["contact.name"]
	spec.Info.Contact.Email = app.Attributes["contact.email"]
	spec.Info.Contact.URL = app.Attributes["contact.url"]

	for k, v := range app.Attributes {
		if strings.HasPrefix(k, "x-") {
			if spec.Info.Extensions == nil {
				spec.Info.Extensions = make(map[string]interface{})
			}
			spec.Info.Extensions[k] = v
		}
	}

	// TODO: Handle multiple environments in attributes
	// TODO: Handle server variables
	server := &openapi3.Server{
		URL:         app.Attributes["env.1.url"],
		Description: app.Attributes["env.1.description"],
		Variables:   map[string]*openapi3.ServerVariable{},
	}
	spec.AddServer(server)
	spec.Components = openapi3.NewComponents()
	spec.Components.Schemas = make(map[string]*openapi3.SchemaRef)
	for k, v := range app.Types {
		spec.Components.Schemas[k] = s.exportType(v)
	}
	for _, v := range app.Endpoints {
		var method, path string
		epPath := strings.Split(v.Path, " ")
		if len(epPath) > 1 {
			method = strings.Split(v.Path, " ")[0]
			path = strings.Split(v.Path, " ")[1]
		} else {
			method = "GET"
			path = v.Path
		}

		// Map Params
		operation := openapi3.NewOperation()
		// Convert to multiline string
		operation.Description = v.Description
		operation.Summary = v.Summary
		operation.Extensions = v.Extensions
		for paramName, paramItem := range v.Params {
			var param *openapi3.Parameter
			var payload *openapi3.SchemaRef
			switch paramItem.In {
			case "header":
				param = openapi3.NewHeaderParameter(paramName)
			case "path":
				param = openapi3.NewPathParameter(paramName)
			case "query":
				param = openapi3.NewQueryParameter(paramName)
			case "body":
				payload = s.exportType(paramItem.Type)
			}

			// TODO: Cookies https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md#parameterObject
			if param != nil {
				param.WithSchema(s.exportType(paramItem.Type).Value)
				param.Schema.Ref = s.exportType(paramItem.Type).Ref
				param.Required = !paramItem.Type.Optional
				operation.AddParameter(param)
			} else {
				operation.RequestBody = &openapi3.RequestBodyRef{
					Value: openapi3.NewRequestBody().WithJSONSchemaRef(payload).WithRequired(!paramItem.Type.Optional),
				}
			}
		}

		// Map Responses
		for _, value := range v.Response {
			response := openapi3.NewResponse()
			schemaRef := s.exportType(value.Type)
			response.WithDescription(value.Name)
			response.WithContent(openapi3.NewContentWithJSONSchemaRef(schemaRef))
			// Parse response code to int
			operation.AddResponse(parseResponseCode(value.Name), response)
		}

		spec.AddOperation(path, method, operation)
	}

	return spec, nil
}

func (s *OpenAPI3Exporter) exportType(t *syslwrapper.Type) *openapi3.SchemaRef {
	var ref string
	var value = openapi3.NewSchema()
	if t == nil {
		return nil
	}
	switch t.Type {
	case "bool":
		value = openapi3.NewBoolSchema()
	case "datetime":
		value = openapi3.NewDateTimeSchema()
	case "date":
		value = openapi3.NewStringSchema().WithFormat("date")
	case "string", "string_8":
		value = openapi3.NewStringSchema()
	case "float":
		value = openapi3.NewFloat64Schema()
		value.WithFormat("float")
	case "decimal":
		value = openapi3.NewFloat64Schema()
		value.WithFormat("double")
	case "int":
		value = openapi3.NewIntegerSchema()
		value.WithFormat("int64")
	case "uuid":
		value = openapi3.NewUUIDSchema()
	case "bytes":
		value = openapi3.NewBytesSchema()
	case "enum":
		// Only string schemas are currently supported
		value = openapi3.NewStringSchema().WithEnum(convertEnum(t.Enum).Data...)
	case "map":
		value = openapi3.NewObjectSchema()
		for k, v := range t.Properties {
			value.Properties[k] = s.exportType(v)
		}
	case "list":
		value = openapi3.NewArraySchema()
		value.Items = s.exportType(t.Items[0])
	case "tuple":
		var required []string
		value = openapi3.NewObjectSchema()
		for k, v := range t.Properties {
			value.Properties[k] = s.exportType(v)
			if !v.Optional {
				required = append(required, k)
			}
		}
		value.Required = required
	case "ref":
		ref = SyslRefToJSONSchema(t.Reference)
	}
	return openapi3.NewSchemaRef(ref, value)
}

type validInputs struct {
	Data []interface{}
}

func convertEnum(syslEnum map[int64]string) validInputs {
	enums := validInputs{}
	for _, str := range syslEnum {
		enums.Data = append(enums.Data, str)
	}
	return enums
}

func parseResponseCode(response string) int {
	// Handle sysl convention for using ok to mean 200
	if response == "ok" {
		return 200
	}
	responseCode, err := strconv.Atoi(response)
	if err != nil {
		return 0
	}
	return responseCode
}
func SyslRefToJSONSchema(syslRef string) string {
	reference := strings.Split(syslRef, ".")
	return "#/components/schemas/" + reference[1]
}
