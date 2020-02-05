// Package catalog takes a sysl module with attributes defined (catalogFields) and serves a webserver listing the applications and endpoints
// It also uses GRPCUI and Redoc in order to generate an interactive page to interact with all the endpoints
// GRPC currently uses server reflection TODO: Support gpcui directly from swagger files
package catalog

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/fullstorydev/grpcui"
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
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
type Server struct {
	Fs       afero.Fs
	Log      *logrus.Logger
	Modules  []*sysl.Module
	Fields   []string
	BasePath string `long:"base-path" description:"the base path to serve the spec and UI at"`
	Path     string
	Resource string
	Flavor   string `short:"F" long:"flavor" description:"the flavor of docs, can be swagger or redoc" default:"redoc" choice:"redoc,swagger"` //nolint: lll
	DocURL   string `long:"doc-url" description:"override the url which takes a url query param to render the doc ui"`
	NoOpen   bool   `long:"no-open" description:"when present won't open the the browser to show the url"`
	NoUI     bool   `long:"no-ui" description:"when present, only the swagger spec will be served"`
	Flatten  bool   `long:"flatten" description:"when present, flatten the swagger spec before serving it"`
	Port     string `long:"port" short:"p" description:"the port to serve this site" env:"PORT"`
	Host     string `long:"host" description:"the interface to serve this site, defaults to 0.0.0.0" env:"HOST"`
}

// WebService is the type which will be rendered on the home page of the html/json as a row
type WebService struct {
	App           *sysl.Application `json:"-"`
	Fields        []string          `json:"-"`
	Attrs         map[string]string
	AppName       string
	SwaggerUILink string
}

// ListHandlers registers handlers for both the homepage, if t is json the header will be set as json content type
func (c *Server) ListHandlers(contents []byte, t string, pattern string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if t == "json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		}
		_, err := w.Write(contents)
		if err != nil {
			panic(err)
		}
	})
}

// Serve Runs the command and runs a webserver on catalogURL of a list of endpoints in the sysl file
func (c *Server) Serve() error {
	services, err := c.BuildCatalog()
	if err != nil {
		return err
	}
	json, err := json.Marshal(services)
	if err != nil {
		return err
	}
	html, err := renderHTML(services)
	if err != nil {
		return err
	}
	c.ListHandlers(json, "json", "/json")
	c.ListHandlers(html, "html", "/")
	err = http.ListenAndServe(c.Host, nil)
	return err
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
func (c *Server) BuildCatalog() ([]WebService, error) {
	var catalog []WebService
	var h http.Handler
	var serviceCount int
	var err error
	for _, m := range c.Modules {
		for serviceName, a := range m.GetApps() {

			if serviceMethod := serviceType(a); serviceMethod != "" {
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
					SwaggerUILink: "/" + serviceName + strconv.Itoa(serviceCount),
				}

				switch serviceMethod {
				case "GRPC":
					h, err = c.GrpcUIHandler(newService)
				case "REST":
					h, err = c.SwaggerUIHandler(newService)
				}
				if err != nil {
					return nil, err
				}
				h = http.StripPrefix(newService.SwaggerUILink, h)
				http.Handle(newService.SwaggerUILink+"/", h)
				catalog = append(catalog, newService)

				c.Log.Infof("Added %s service: %s from %s",
					newService.Attrs["type"],
					newService.AppName,
					newService.Attrs["deploy.prod.url"])
			}
		}
	}
	return catalog, nil
}

// GrpcUIHandler creates and returns a http handler for a grpcui server
func (c *Server) GrpcUIHandler(service WebService) (http.Handler, error) {
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
	methods, err := grpcui.AllMethodsViaReflection(ctx, cc)
	if err != nil {
		return nil, err
	}
	err = c.GrpcUIHTML(methods)
	if err != nil {
		return nil, err
	}
	return standalone.HandlerViaReflection(ctx, cc, service.SwaggerUILink)

}

// GrpcUIHTML Writes all the static files from grpcui to serve
func (c *Server) GrpcUIHTML(methods []*desc.MethodDescriptor) error {
	file, err := c.Fs.Create("index.html")
	if err != nil {
		return err
	}
	_, err = file.Write(grpcui.WebFormContents("invoke", "metadata", methods))
	if err != nil {
		return err
	}
	file, err = c.Fs.Create("grpc-web-form.js")
	if err != nil {
		return err
	}
	_, err = file.Write(grpcui.WebFormScript())
	if err != nil {
		return err
	}
	file, err = c.Fs.Create("grpc-web-form.css")
	if err != nil {
		return err
	}
	_, err = file.Write(grpcui.WebFormSampleCSS())
	return err

}

// SwaggerUIHandler creates and returns a http handler for a SwaggerUI server
func (c *Server) SwaggerUIHandler(service WebService) (http.Handler, error) {
	c.Resource = service.SwaggerUILink
	swaggerExporter := exporter.MakeSwaggerExporter(service.App, c.Log)
	err := swaggerExporter.GenerateSwagger()
	if err != nil {
		return nil, err
	}
	output, err := swaggerExporter.SerializeOutput("json")
	if err != nil {
		return nil, err
	}
	return c.SwaggerUI(output), nil
}
