// File adapted from https://github.com/go-swagger/go-swagger/blob/master/cmd/swagger/commands/serve.go
// to be used within go and served an arbitrary endpoints (not just /docs)

package catalog

import (
	"net/http"
	"path"

	"github.com/gorilla/handlers"

	"github.com/go-openapi/runtime/middleware"
)

// SwaggerUI takes the contents of a swagger file and creates a handler for the interactive redoc
func (s *Server) SwaggerUI(contents []byte) http.Handler {
	if s.Flavor == "" {
		s.Flavor = "redoc"
	}
	if s.BasePath == "" {
		s.BasePath = "/"
	}
	if s.Path == "" {
		s.Path = "/"
	}

	handler := middleware.Redoc(middleware.RedocOpts{
		BasePath: s.BasePath,
		SpecURL:  path.Join("/", "swagger.json"),
		Path:     s.Path,
	}, nil)

	handler = handlers.CORS()(middleware.Spec(s.BasePath, contents, handler))
	return handler
}
