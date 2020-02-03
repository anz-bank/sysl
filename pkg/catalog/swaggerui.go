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

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/go-openapi/spec"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// Server to serve a swagger spec with docs ui
type Server struct {
	Fs      afero.Fs
	Log     *logrus.Logger
	Modules []*sysl.Module

	BasePath string `long:"base-path" description:"the base path to serve the spec and UI at"`
	Path     string
	Resource string
	Flavor   string `short:"F" long:"flavor" description:"the flavor of docs, can be swagger or redoc" default:"redoc" choice:"redoc,swagger"` //nolint: lll
	DocURL   string `long:"doc-url" description:"override the url which takes a url query param to render the doc ui"`
	NoOpen   bool   `long:"no-open" description:"when present won't open the the browser to show the url"`
	NoUI     bool   `long:"no-ui" description:"when present, only the swagger spec will be served"`
	Flatten  bool   `long:"flatten" description:"when present, flatten the swagger spec before serving it"`
	Port     int    `long:"port" short:"p" description:"the port to serve this site" env:"PORT"`
	Host     string `long:"host" description:"the interface to serve this site, defaults to 0.0.0.0" env:"HOST"`
}

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
