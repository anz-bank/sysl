package catalog

import (
	"bytes"
	"context"
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

// CatalogServer is a server for the catalog command which hosts an interactive web ui
type CatalogServer struct {
	Port   int
	Host   string
	Type   string
	Fs     afero.Fs
	Log    *logrus.Logger
	Module *sysl.Module
}

// service is the type which will be rendered on the home page of the html/json as a row
type webService struct {
	App         *sysl.Application `json:"-"`
	Filename    string            `josn:"filename"`
	Team        string            `josn:"team"`
	Owner       string            `josn:"owner"`
	Email       string            `josn:"email"`
	URL         string            `josn:"url"`
	ServiceName string            `josn:"servicename"`
	Type        string            `josn:"type"`
	Link        string            `josn:"link"`
}

// Serve Runs the command and runs a webserver on catalogURL of a list of endpoints in the sysl file
func (c *CatalogServer) Serve() error {

	var contents bytes.Buffer
	var err error

	services, err := c.buildCatalog(c.Module)

	switch c.Type {
	case "json":
		err = renderJSON(services, &contents)
	case "html":
		err = renderHTML(services, &contents)
	}

	if err != nil {
		c.Log.Errorf(err.Error())
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(contents.Bytes())
	})

	addr := c.Host + ":" + strconv.Itoa(c.Port)
	fmt.Println(addr)
	err = http.ListenAndServe(addr, nil)

	c.Log.Errorf(err.Error())
	return err

}

// buildCatalog unpacks a sysl modules and registers a http handler for the grpcui or swagger ui
func (c *CatalogServer) buildCatalog(m *sysl.Module) ([]webService, error) {
	var ser []webService
	var h http.Handler
	var err error
	for _, a := range m.GetApps() {
		atts := a.GetAttrs()
		serviceName := strings.Join(a.Name.GetPart(), ",")
		serviceName = strings.ReplaceAll(serviceName, "/", "")
		serviceName = strings.ReplaceAll(serviceName, " ", "")
		serviceName = strings.ToLower(serviceName)

		newService := webService{
			App:         a,
			Filename:    atts["proto"].GetS(),
			URL:         atts["url"].GetS(),
			ServiceName: serviceName,
			Type:        atts["type"].GetS(),
			Owner:       atts["owner"].GetS(),
			Email:       atts["email"].GetS(),
			Link:        "/" + serviceName,
		}
		switch newService.Type {
		case "GRPC":
			h, err = c.registerGrpcUI(newService)
		case "REST":
			h, err = c.registerSwaggerUI(newService)
		}
		if err != nil {
			c.Log.Errorf(err.Error())
			return nil, err
		}
		h = http.StripPrefix(newService.Link, h)
		http.Handle(newService.Link+"/", h)
		ser = append(ser, newService)

	}
	return ser, nil
}

// registerGrpcUI creates and returns a http handler for a grpcui server
func (c *CatalogServer) registerGrpcUI(service webService) (http.Handler, error) {
	ctx := context.Background()
	cc, err := grpc.DialContext(ctx, service.URL, grpc.WithInsecure())
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	methods, err := grpcui.AllMethodsViaReflection(ctx, cc)
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err
	}
	file, _ := c.Fs.Create("index.html")
	file.Write(grpcui.WebFormContents("invoke", "metadata", methods))
	file, _ = c.Fs.Create("grpc-web-form.js")
	file.Write(grpcui.WebFormScript())
	file, _ = c.Fs.Create("grpc-web-form.css")
	file.Write(grpcui.WebFormSampleCSS())
	h, err := standalone.HandlerViaReflection(ctx, cc, service.URL)
	if err != nil {
		c.Log.Errorf(err.Error())
		return nil, err

	}

	return h, nil
}

// registerSwaggerUI creates and returns a http handler for a SwaggerUI server
func (c *CatalogServer) registerSwaggerUI(service webService) (http.Handler, error) {
	basePath := "/"
	swag := ServeCmd{BasePath: basePath, Port: c.Port, Path: "/", Resource: service.Link}
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
	h := swag.SwaggerUI(output)
	return h, nil
}
