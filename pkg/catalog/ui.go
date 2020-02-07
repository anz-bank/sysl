// File adapted from https://github.com/go-swagger/go-swagger/blob/master/cmd/swagger/commands/serve.go
// to be used within go and served an arbitrary endpoints (not just /docs)

package catalog

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
)

// SwaggerUI takes the contents of a swagger file and creates a handler for the interactive redoc
func (service *WebService) SwaggerUI(contents []byte) (http.Handler, error) {
	if service.SwaggerUILink == "" {
		service.SwaggerUILink = "/"
	}
	specDoc, err := loads.Analyzed(contents, "2.0")
	if err != nil {
		return nil, err
	}
	b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return nil, err
	}
	handler := http.NotFoundHandler()
	handler = middleware.Redoc(middleware.RedocOpts{
		BasePath: service.SwaggerUILink,
		SpecURL:  path.Join(service.SwaggerUILink, "swagger.json"),
		Path:     "/",
	}, handler)
	handler = handlers.CORS()(middleware.Spec(service.SwaggerUILink, b, handler))
	return handler, nil
}
