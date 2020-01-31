package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fullstorydev/grpcui"
	"github.com/sirupsen/logrus"

	"net/http"

	"github.com/anz-bank/sysl/pkg/exporter"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

var catalogFields = []string{
	"team",
	"team.slack",
	"owner.name",
	"owner.email",
	"file.version",
	"release.version",
	"description",
	"deploy.env1.url",
	"deploy.sit1.url",
	"deploy.sit2.url",
	"deploy.qa.url",
	"deploy.prod.url",
	"repo.url",
	"type",
	"confluence.url",
}

// Server is a server for the catalog command which hosts an interactive web ui
type Server struct {
	Port    int
	Host    string
	Fs      afero.Fs
	Log     *logrus.Logger
	Modules []*sysl.Module
}

// WebService is the type which will be rendered on the home page of the html/json as a row
type WebService struct {
	App           *sysl.Application `json:"-"`
	Fields        []string          `json:"-"`
	Attrs         map[string]string
	ServiceName   string
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
			c.Log.Errorf(err.Error())
			panic(err)
		}
	})
}

// Serve Runs the command and runs a webserver on catalogURL of a list of endpoints in the sysl file
func (c *Server) Serve() error {
	var err error
	services, err := c.BuildCatalog()
	if err != nil {
		c.Log.Errorf(err.Error())
		return err
	}
	json, err := json.Marshal(services)
	if err != nil {
		c.Log.Errorf(err.Error())
		return err
	}
	html, err := renderHTML(services)
	if err != nil {
		c.Log.Errorf(err.Error())
		return err
	}
	c.ListHandlers(json, "json", "/json")
	c.ListHandlers(html, "html", "/")
	addr := c.Host + ":" + strconv.Itoa(c.Port)
	fmt.Println(addr)
	err = http.ListenAndServe(addr, nil)
	c.Log.Errorf(err.Error())
	return err
}

// BuildCatalog unpacks a sysl modules and registers a http handler for the grpcui or swagger ui
func (c *Server) BuildCatalog() ([]WebService, error) {
	var ser []WebService
	var h http.Handler
	var err error
	for _, m := range c.Modules {
		for _, a := range m.GetApps() {
			atts := a.GetAttrs()
			serviceName := strings.Join(a.Name.GetPart(), ",")
			serviceName = strings.ReplaceAll(serviceName, "/", "")
			serviceName = strings.ReplaceAll(serviceName, " ", "")
			serviceName = strings.ToLower(serviceName)
			var attr = make(map[string]string, 10)

			for key, value := range atts {
				attr[key] = value.GetS()
			}

			newService := WebService{
				App:           a,
				Fields:        catalogFields,
				Attrs:         attr,
				ServiceName:   serviceName,
				SwaggerUILink: "/" + serviceName,
			}
			switch newService.Attrs["type"] {
			case "GRPC":
				h, err = c.GrpcUIHandler(newService)
			case "REST":
				h, err = c.SwaggerUIHandler(newService)
			}
			if err != nil {
				c.Log.Errorf(err.Error())
				return nil, err
			}
			h = http.StripPrefix(newService.SwaggerUILink, h)
			http.Handle(newService.SwaggerUILink+"/", h)
			ser = append(ser, newService)
		}
	}
	return ser, nil
}

// GrpcUIHandler creates and returns a http handler for a grpcui server
func (c *Server) GrpcUIHandler(service WebService) (http.Handler, error) {
	ctx := context.Background()
	cc, err := grpc.DialContext(ctx, service.Attrs["deploy.prod.url"], grpc.WithInsecure())
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	methods, err := grpcui.AllMethodsViaReflection(ctx, cc)
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	file, err := c.Fs.Create("index.html")
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	_, err = file.Write(grpcui.WebFormContents("invoke", "metadata", methods))
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	file, err = c.Fs.Create("grpc-web-form.js")
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	_, err = file.Write(grpcui.WebFormScript())
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	file, err = c.Fs.Create("grpc-web-form.css")
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	_, err = file.Write(grpcui.WebFormSampleCSS())
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	h, err := standalone.HandlerViaReflection(ctx, cc, service.SwaggerUILink)
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	return h, nil
}

// SwaggerUIHandler creates and returns a http handler for a SwaggerUI server
func (c *Server) SwaggerUIHandler(service WebService) (http.Handler, error) {
	basePath := "/"
	swag := ServeCmd{BasePath: basePath, Port: c.Port, Path: "/", Resource: service.SwaggerUILink}
	swaggerExporter := exporter.MakeSwaggerExporter(service.App, nil)
	err := swaggerExporter.GenerateSwagger()
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	output, err := swaggerExporter.SerializeOutput("json")
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	return swag.SwaggerUI(output), nil
}
