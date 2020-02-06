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
func (service *WebService) SwaggerUI(contents []byte) http.Handler {
	if service.SwaggerUILink == "" {
		service.SwaggerUILink = "/"
	}
	specDoc, _ := loads.Analyzed(contents, "2.0")
	b, _ := json.MarshalIndent(specDoc.Spec(), "", "  ")
	handler := http.NotFoundHandler()
	handler = middleware.Redoc(middleware.RedocOpts{
		BasePath: service.SwaggerUILink,
		SpecURL:  path.Join(service.SwaggerUILink, "swagger.json"), //"https://raw.githubusercontent.com/chimauwah/services-api-tech-challenge/master/swagger.json",
		Path:     "/",
	}, handler)
	handler = handlers.CORS()(middleware.Spec(service.SwaggerUILink, b, handler))
	return handler
}
