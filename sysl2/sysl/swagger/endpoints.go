package swagger

import (
	"github.com/go-openapi/spec"
)

type Endpoint struct {
	Path        string
	Description string

	Params Parameters

	Responses *spec.Responses
}

// nolint:gochecknoglobals
var methodDisplayOrder = []string{"GET", "PUT", "POST", "DELETE", "PATCH"}

func initEndpoint(path string,
	op *spec.Operation, params []spec.Parameter,
	types TypeList, globals Parameters) Endpoint {

	apiParams := buildParameters(params, types, globals, Parameters{})

	res := Endpoint{
		Path:        path,
		Description: op.Description,
		Responses:   op.Responses,
		Params:      buildParameters(op.Parameters, types, globals, apiParams),
	}
	return res
}

func InitEndpoints(doc *spec.Swagger, types TypeList, globals Parameters) map[string][]Endpoint {
	res := map[string][]Endpoint{}

	for path, item := range doc.Paths.Paths {

		ops := map[string]*spec.Operation{
			"GET":    item.Get,
			"PUT":    item.Put,
			"POST":   item.Post,
			"DELETE": item.Delete,
			"PATCH":  item.Patch,
		}

		for method, op := range ops {
			if op != nil {
				res[method] = append(res[method], initEndpoint(path, op, item.Parameters, types, globals))
			}
		}

	}

	return res
}
