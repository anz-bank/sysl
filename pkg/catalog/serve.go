package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fullstorydev/grpcui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"net/http"

	"github.com/anz-bank/sysl/pkg/exporter"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/fullstorydev/grpcui/standalone"
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

// Server to set context of catalog
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

// WebService is the type which will be rendered on the home page of the html/json as a row
type WebService struct {
	App           *sysl.Application `json:"-"`
	Fields        []string          `json:"-"`
	Attrs         map[string]string
	ServiceName   string
	SwaggerUILink string
}

func (c Server) String() string {
	return "Server:" + c.BasePath + c.Path + string(c.Port)
}
func (c WebService) String() string {
	return "WebService:" + c.ServiceName + c.SwaggerUILink
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
	var serviceCount int
	for _, m := range c.Modules {
		for i, a := range m.GetApps() {
			serviceCount++
			atts := a.GetAttrs()
			serviceName := strings.Join(a.Name.GetPart(), "")
			var re = regexp.MustCompile(`(?m)\w*`)
			serviceName = strings.Join(re.FindAllString(serviceName, -1), "")
			serviceName = strings.ToLower(serviceName) + strconv.Itoa(serviceCount)
			// serviceName = strings.ReplaceAll(serviceName, " ", "")
			fmt.Println(serviceName)

			var attr = make(map[string]string, 10)

			for key, value := range atts {
				attr[key] = value.GetS()
			}
			c.Log.Infof("eofn: %s", i)
			newService := WebService{
				App:           a,
				Fields:        catalogFields,
				Attrs:         attr,
				ServiceName:   serviceName,
				SwaggerUILink: "/" + serviceName,
			}
			c.Log.Infof("Adding %s service: %s from %s", newService.Attrs["type"], newService.ServiceName, newService.Attrs["deploy.prod.url"])
			switch strings.ToUpper(newService.Attrs["type"]) {
			case "GRPC":
				h, err = c.GrpcUIHandler(newService)
			case "REST":
				c.Log.Info("Hello")
				c.Log.Info(c)
				c.Log.Info(newService)
				h, err = c.SwaggerUIHandler(newService)
				c.Log.Info("elvns")

			}
			if err != nil {
				c.Log.Errorf(err.Error())
			}
			h = http.StripPrefix(newService.SwaggerUILink, h)
			http.Handle(newService.SwaggerUILink+"/", h)
			c.Log.Errorf("newService.SwaggerUILink" + newService.SwaggerUILink + "/")
			ser = append(ser, newService)
			c.Log.Infof("Added %s service: %s from %s", newService.Attrs["type"], newService.ServiceName, newService.Attrs["deploy.prod.url"])
		}
	}
	return ser, nil
}

// GrpcUIHandler creates and returns a http handler for a grpcui server
func (c *Server) grpcGrpcUIHandler(service WebService) (http.Handler, error) {
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
	swag := Server{BasePath: basePath, Port: c.Port, Path: "/", Resource: service.SwaggerUILink}
	swaggerExporter := exporter.MakeSwaggerExporter(service.App, c.Log)
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
