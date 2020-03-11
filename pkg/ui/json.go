package ui

import (
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

const GRPC = "grpc"
const REST = "rest"

// WebService is the struct which will be translated to JSON via reflection.
type WebService struct {
	Name       string
	Attributes map[string]string
	Endpoints  []Endpoint
	Type       string
	Path       string
}

// BuildWebService takes a sysl Application and returns a json-exportable representation of Sysl
func BuildWebService(a *sysl.Application) (*WebService, error) {
	// TODO: Ensure there are no duplicate service names.
	serviceMethod := serviceType(a)
	serviceName := a.GetName().GetPart()[0]
	attributes := mapAttributes(a.GetAttrs())
	endpoints := mapEndpoints(a.GetEndpoints())

	newService := &WebService{
		Name:       serviceName,
		Attributes: attributes,
		Endpoints:  endpoints,
		Type:       serviceMethod,
		Path:       makePath(serviceMethod, "/", serviceName),
	}
	return newService, nil
}

func serviceType(app *sysl.Application) string {
	if syslutil.HasPattern(app.GetAttrs(), GRPC) {
		return GRPC
	} else if syslutil.HasPattern(app.GetAttrs(), REST) {
		return REST
	}
	return ""
}

func mapAttributes(attributes map[string]*sysl.Attribute) map[string]string {
	var attr = make(map[string]string, 15)
	for key, value := range attributes {
		attr[key] = value.GetS()
	}
	return attr
}

type Endpoint struct {
	Path     string
	Request  string
	Response string
}

func mapEndpoints(ep map[string]*sysl.Endpoint) []Endpoint {
	var endpoints []Endpoint
	for key, value := range ep {
		endpoints = append(endpoints, Endpoint{
			Path:     key,
			Request:  value.GetName(),
			Response: value.GetName(),
		})
	}
	return endpoints
}
