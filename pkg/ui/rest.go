package ui

import (
	"net/http"

	"encoding/json"

	"github.com/anz-bank/sysl/pkg/exporter"
	"github.com/sirupsen/logrus"

	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
)

// buildRestHandler creates a http handler to serve up SwaggerUI docs
func (b *APIDocBuilder) buildRestHandler(basepath string) error {
	swaggerExporter := exporter.MakeSwaggerExporter(b.app, logrus.New())
	err := swaggerExporter.GenerateSwagger()
	if err != nil {
		return err
	}
	output, err := swaggerExporter.SerializeOutput("json")
	if err != nil {
		return err
	}
	h, err := b.SwaggerUI(output, basepath)
	if err != nil {
		return err
	}
	b.doc.handler = h
	return nil
}

// SwaggerUI takes the contents of a swagger file and creates a handler to display the the endpoints
func (b *APIDocBuilder) SwaggerUI(contents []byte, basepath string) (http.Handler, error) {
	specDoc, err := loads.Analyzed(contents, "2.0")
	if err != nil {
		return nil, err
	}
	bytes, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return nil, err
	}
	// Save json to service object
	b.doc.spec = bytes

	handler := http.NotFoundHandler()
	handler = middleware.Redoc(middleware.RedocOpts{
		BasePath: basepath,
		SpecURL:  path.Join("/rest/spec", b.app.GetName().GetPart()[0]),
		Path:     "/",
	}, handler)
	handler = handlers.CORS()(middleware.Spec(basepath, bytes, handler))
	return handler, nil
}
