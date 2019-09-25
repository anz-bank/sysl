package swagger

import (
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
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
	types TypeList, globals Parameters, logger *logrus.Logger) Endpoint {

	apiParams := buildParameters(params, types, globals, Parameters{}, logger)

	res := Endpoint{
		Path:        path,
		Description: op.Description,
		Responses:   op.Responses,
		Params:      buildParameters(op.Parameters, types, globals, apiParams, logger),
	}
	return res
}

func InitEndpoints(doc *spec.Swagger, types TypeList, globals Parameters, logger *logrus.Logger) map[string][]Endpoint {
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
				res[method] = append(res[method], initEndpoint(path, op, item.Parameters, types, globals, logger))
			}
		}

	}

	for key := range res {
		key := key
		sort.SliceStable(res[key], func(i, j int) bool {
			return strings.Compare(res[key][i].Path, res[key][j].Path) < 0
		})
	}

	return res
}
