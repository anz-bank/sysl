package exporter

import (
	proto "github.com/anz-bank/sysl/src/proto"
	yaml "github.com/ghodss/yaml"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

type SwaggerExporter struct {
	app          *proto.Application
	buildSwagger *spec.Swagger
	log          *logrus.Logger
}

func MakeSwaggerExporter(app *proto.Application, logger *logrus.Logger) *SwaggerExporter {
	return &SwaggerExporter{
		app:          app,
		buildSwagger: &spec.Swagger{},
		log:          logger,
	}
}

func (v *SwaggerExporter) GenerateSwagger() error {
	v.buildSwagger.Swagger = "2.0"
	v.buildSwagger.Host = v.app.GetAttrs()["host"].GetS()

	v.buildSwagger.SwaggerProps.Info = &spec.Info{}
	v.buildSwagger.Paths = &spec.Paths{}
	v.buildSwagger.Paths.Paths = map[string]spec.PathItem{}
	v.buildSwagger.Definitions = spec.Definitions{}

	v.buildSwagger.SwaggerProps.Info.Title = v.app.LongName
	v.buildSwagger.SwaggerProps.Info.Description = v.app.GetAttrs()["description"].GetS()
	v.buildSwagger.SwaggerProps.Info.Version = v.app.GetAttrs()["version"].GetS()
	if v.buildSwagger.SwaggerProps.Info.Version == "" {
		v.buildSwagger.SwaggerProps.Info.Version = "0.0.0"
	}

	// parse type defs
	typeExporter := makeTypeExporter(v.log)
	typeExportError := typeExporter.exportTypes(v.app.GetTypes(), v.buildSwagger.Definitions)
	if typeExportError != nil {
		return typeExportError
	}

	endpointExporter := makeEndpointExporter(typeExporter, v.log)

	// iterate over each endpoint in the selected application
	for endpointName, endpoint := range v.app.Endpoints {
		err := endpointExporter.exportEndpoint(endpointName, endpoint, v.buildSwagger.Paths.Paths)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *SwaggerExporter) SerializeToYaml() ([]byte, error) {
	jsonSpec, err := v.buildSwagger.MarshalJSON()
	if err != nil {
		return nil, err
	}
	yamlSpec, err := yaml.JSONToYAML(jsonSpec)
	if err != nil {
		return nil, err
	}
	return yamlSpec, nil
}
