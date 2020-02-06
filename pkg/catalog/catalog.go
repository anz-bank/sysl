// Package catalog takes a sysl module with attributes defined (catalogFields) and serves a webserver listing the applications and endpoints
// It also uses GRPCUI and Redoc in order to generate an interactive page to interact with all the endpoints
// GRPC currently uses server reflection TODO: Support gpcui directly from swagger files
package catalog

import (
	"context"
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"google.golang.org/grpc"

	"net/http"

	"github.com/anz-bank/sysl/pkg/exporter"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/fullstorydev/grpcui/standalone"
)

// Server to set context of catalog
// Todo: Simplify this
type Server struct {
	Fs       afero.Fs       // Required
	Log      *logrus.Logger // Required
	Modules  []*sysl.Module // Required
	Fields   []string       // Required
	Host     string         // Required
	services []*WebService
	router   *mux.Router
	BasePath string
}

func (s *Server) Setup() {
	if s.BasePath == "" {
		s.BasePath = "/"
	}
	s.router = mux.NewRouter()
	s.routes()
}

// WebService is the type which will be rendered on the home page of the html/json as a row
type WebService struct {
	App           *sysl.Application `json:"-"`
	Fields        []string          `json:"-"`
	Attrs         map[string]string
	AppName       string
	SwaggerUILink string
	Log           *logrus.Logger // Required
	handler       http.Handler
}

func (c *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	c.router.ServeHTTP(w, r)
}

func (c *Server) routes() {
	services, err := c.BuildCatalog()
	if err != nil {
		c.Log.Info(err)
	}
	html, err := renderHTML(services)

	c.router.HandleFunc("/", c.ListHandlers(html, "html", "/"))
	for _, service := range services {
		c.router.PathPrefix(service.SwaggerUILink).Handler(service.handler)
	}

}

// Serve Runs the command and runs a webserver on catalogURL of a list of endpoints in the sysl file
func (c *Server) Serve() error {
	c.Log.Info("serving")
	log.Fatal(http.ListenAndServe(":8080", c))
	return nil
}

var catalogFields = `team,
team.slack,
owner.name,
owner.email,
file.version,
release.version,
description,
deploy.env1.url,
deploy.sit1.url,
deploy.sit2.url,
deploy.qa.url,
deploy.prod.url,
repo.url,
type,
confluence.url`

func serviceType(app *sysl.Application) string {
	if syslutil.HasPattern(app.GetAttrs(), "GRPC") {
		return "GRPC"
	} else if syslutil.HasPattern(app.GetAttrs(), "REST") {
		return "REST"
	}
	return ""
}

// BuildCatalog unpacks a sysl modules and registers a http handler for the grpcui or swagger ui
func (c *Server) BuildCatalog() ([]WebService, error) {
	var catalog []WebService
	var serviceCount int
	for _, m := range c.Modules {
		for serviceName, a := range m.GetApps() {
			if serviceMethod := serviceType(a); serviceMethod != "" {
				var err error
				serviceCount++
				var attr = make(map[string]string, 10)

				atts := a.GetAttrs()
				for key, value := range atts {
					attr[key] = value.GetS()
				}
				newService := WebService{
					App:           a,
					Fields:        c.Fields,
					Attrs:         attr,
					AppName:       serviceName,
					Log:           c.Log,
					SwaggerUILink: path.Join(c.BasePath, serviceName+strconv.Itoa(serviceCount)),
				}
				switch serviceMethod {
				case "GRPC":

					newService.handler, err = newService.GrpcUIHandler()

				case "REST":
					newService.handler, err = newService.SwaggerUIHandler()
				}
				if err != nil {
					c.Log.Infof("Added %s service: %s from %s failed: %s\n",
						newService.Attrs["type"],
						newService.AppName,
						newService.Attrs["deploy.prod.url"],
						err)
					continue
				}
				c.Log.Infof("Added %s service: %s from %s",
					newService.Attrs["type"],
					newService.AppName,
					newService.Attrs["deploy.prod.url"])
				catalog = append(catalog, newService)
			}
		}
	}
	if len(catalog) == 0 {
		return catalog, fmt.Errorf(`No services registered;
                        Make sure ~GRPC or ~REST are in the service definitions`)
	}
	return catalog, nil
}

// GrpcUIHandler creates and returns a http handler for a grpcui server
func (service *WebService) GrpcUIHandler() (http.Handler, error) {
	ctx := context.Background()

	var opts []grpc.DialOption
	network := "tcp"

	creds, err := grpcurl.ClientTransportCredentials(false, "", "", "")
	if err != nil {
		return nil, err
	}
	cc, err := grpcurl.BlockingDial(ctx, network, service.Attrs["deploy.prod.url"], creds, opts...)
	// If that failed, try an insecure dial
	if err != nil {
		cc, err = grpc.DialContext(ctx, service.Attrs["deploy.prod.url"], grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	}
	h, err := standalone.HandlerViaReflection(ctx, cc, service.SwaggerUILink)
	if err != nil {
		return nil, err
	}
	h = http.StripPrefix(service.SwaggerUILink, h)
	return h, nil

}

// SwaggerUIHandler creates and returns a http handler for a SwaggerUI server
// Todo: move this to its own file
func (service *WebService) SwaggerUIHandler() (http.Handler, error) {
	swaggerExporter := exporter.MakeSwaggerExporter(service.App, service.Log)
	err := swaggerExporter.GenerateSwagger()
	if err != nil {
		return nil, err
	}
	output, err := swaggerExporter.SerializeOutput("json")
	if err != nil {
		return nil, err
	}
	return service.SwaggerUI(output), nil
}

// ListHandlers registers handlers for both the homepage, if t is json the header will be set as json content type
func (c *Server) ListHandlers(contents []byte, t string, pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if t == "json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		}
		_, err := w.Write(contents)
		if err != nil {
			panic(err)
		}
	}
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
