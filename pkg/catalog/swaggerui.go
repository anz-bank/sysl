// File adapted from https://github.com/go-swagger/go-swagger/blob/master/cmd/swagger/commands/serve.go
// to be used within go and served an arbitrary endpoints (not just /docs)

package catalog

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/go-openapi/spec"
	"github.com/gorilla/handlers"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// SwaggerUI takes the contents of a swagger file and creates a handler for the interactive redoc
func (s *Server) SwaggerUI(contents []byte) http.Handler {
	if s.Flavor == "" {
		s.Flavor = "redoc"
	}
	specDoc, err := loads.Analyzed(contents, "2.0")

	if err != nil {
		panic(err)
	}

	if s.Flatten {
		var err error
		specDoc, err = specDoc.Expanded(&spec.ExpandOptions{
			SkipSchemas:         false,
			ContinueOnError:     true,
			AbsoluteCircularRef: true,
		})

		if err != nil {
			panic(err)
		}
	}

	b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		panic(err)
	}

	basePath := s.BasePath
	if basePath == "" {
		basePath = "/"
	}

	listener, err := net.Listen("tcp4", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		panic(err)
	}
	sh, sp, err := swag.SplitHostPort(listener.Addr().String())
	if err != nil {
		panic(err)
	}
	if sh == "0.0.0.0" {
		sh = "localhost"
	}
	visit := s.DocURL
	handler := http.NotFoundHandler()
	if !s.NoUI {
		if s.Flavor == "redoc" {
			handler = middleware.Redoc(middleware.RedocOpts{
				BasePath: basePath,
				SpecURL:  path.Join(s.Resource, "swagger.json"),
				Path:     s.Path,
			}, handler)
			visit = fmt.Sprintf("http://%s:%d%s", sh, sp, path.Join(basePath, s.Path))
		} else if visit != "" || s.Flavor == "swagger" {
			if visit == "" {
				visit = "http://petstore.swagger.io/"
			}
			u, err := url.Parse(visit)
			if err != nil {
				panic(err)
			}
			q := u.Query()
			q.Add("url", fmt.Sprintf("http://%s:%d%s", sh, sp, path.Join(s.Resource, "swagger.json")))
			u.RawQuery = q.Encode()
			visit = u.String()
		}
	}
	fmt.Println(visit)
	handler = handlers.CORS()(middleware.Spec(basePath, b, handler))
	return handler
}
