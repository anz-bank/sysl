// Package catalog takes a sysl module with attributes defined (catalogFields)
// and serves a webserver listing the applications and endpoints
// It also uses GRPCUI and Redoc in order to generate an interactive page to interact with all the endpoints
// GRPC currently uses server reflection TODO: Support gpcui directly from swagger files
package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strconv"

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

func (s *Server) Setup() error {
	if s.BasePath == "" {
		s.BasePath = "/"
	}
	s.router = mux.NewRouter()
	return s.routes()
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() error {
	var err error
	s.services, err = s.BuildCatalog()
	if err != nil {
		return err
	}
	html, err := renderHTML(s.services)
	if err != nil {
		return err
	}
	s.router.HandleFunc("/", s.ListHandlers(html, "html", "/"))

	jsonBytes, err := json.Marshal(s.services)
	if err != nil {
		return err
	}
	jsonHandler, err := s.ListHandlers()
	s.router.HandleFunc("/json", s.ListHandlers(html, "json", "/json"))
	for _, service := range s.services {
		s.router.PathPrefix(service.SwaggerUILink).Handler(service.handler)
	}
	return nil
}

// Serve Runs the command and runs a webserver on catalogURL of a list of endpoints in the sysl file
func (s *Server) Serve() error {
	s.Log.Info("serving")
	log.Fatal(http.ListenAndServe(":8080", s))
	return nil
}

func serviceType(app *sysl.Application) string {
	if syslutil.HasPattern(app.GetAttrs(), "GRPC") {
		return "GRPC"
	} else if syslutil.HasPattern(app.GetAttrs(), "REST") {
		return "REST"
	}
	return ""
}

// BuildCatalog unpacks a sysl modules and registers a http handler for the grpcui or swagger ui
func (s *Server) BuildCatalog() ([]*WebService, error) {
	var catalog []*WebService
	var serviceCount int
	for _, m := range s.Modules {
		for serviceName, a := range m.GetApps() {
			if serviceMethod := serviceType(a); serviceMethod != "" {
				var err error
				serviceCount++
				var attr = make(map[string]string, 10)

				atts := a.GetAttrs()
				for key, value := range atts {
					attr[key] = value.GetS()
				}
				newService := &WebService{
					App:           a,
					Fields:        s.Fields,
					Attrs:         attr,
					AppName:       serviceName,
					Log:           s.Log,
					SwaggerUILink: path.Join(s.BasePath, serviceName+strconv.Itoa(serviceCount)),
				}
				switch serviceMethod {
				case "GRPC":

					newService.handler, err = newService.GrpcUIHandler()

				case "REST":
					newService.handler, err = newService.SwaggerUIHandler()
				}
				if err != nil {
					s.Log.Infof("Added %s service: %s from %s failed: %s\n",
						newService.Attrs["type"],
						newService.AppName,
						newService.Attrs["deploy.prod.url"],
						err)
					continue
				}
				s.Log.Infof("Added %s service: %s from %s",
					newService.Attrs["type"],
					newService.AppName,
					newService.Attrs["deploy.prod.url"])
				catalog = append(catalog, newService)
			}
		}
	}
	if len(catalog) == 0 {
		return catalog, fmt.Errorf(`no services registered; Make sure ~GRPC or ~REST are in the service definitions`) //nolint:lll
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
	return service.SwaggerUI(output)
}

// ListHandlers registers handlers for both the homepage, if t is json the header will be set as json content type
func (s *Server) ListHandlers(contents []byte, t string, pattern string) http.HandlerFunc {
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
