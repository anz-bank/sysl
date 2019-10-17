package importer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

// Response is either going to be freetext or a type
type Response struct {
	Text string
	Type Type
}

type Endpoint struct {
	Path        string
	Description string

	Params Parameters

	Responses []Response
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
		Responses:   buildResponses(path, op.Responses, types, logger),
		Params:      buildParameters(op.Parameters, types, globals, apiParams, logger),
	}
	return res
}

func buildResponses(path string, responses *spec.Responses, types TypeList, logger *logrus.Logger) []Response {
	var outs []Response

	for statusCode, response := range responses.StatusCodeResponses {
		if schema := response.Schema; schema != nil {
			t, found := types.FindFromSchema(*schema, &typeData{logger: logger})
			if !found {
				logger.Errorf("Responses type for code %d not found, endpoint: %s, skipping", statusCode, path)
				continue
			}
			outs = append(outs, Response{Type: t})
		} else {
			outs = append(outs, Response{Text: fmt.Sprintf("%d", statusCode)})
		}
	}
	if responses.Extensions != nil {
		logger.Warnf("x-* responses not implemented, endpoint: %s", path)
		for key := range responses.Extensions {
			outs = append(outs, Response{Text: key})
		}
	}
	if responses.Default != nil {
		logger.Warnf("default responses not implemented, endpoint: %s", path)
		outs = append(outs, Response{Text: "default"})
	}

	return outs
}

func InitEndpoints(doc *spec.Swagger, types TypeList, globals Parameters, logger *logrus.Logger) []MethodEndpoints {
	epMap := map[string][]Endpoint{}

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
				epMap[method] = append(epMap[method], initEndpoint(path, op, item.Parameters, types, globals, logger))
			}
		}
	}

	for key := range epMap {
		key := key
		sort.SliceStable(epMap[key], func(i, j int) bool {
			return strings.Compare(epMap[key][i].Path, epMap[key][j].Path) < 0
		})
	}

	var result []MethodEndpoints
	for _, method := range methodDisplayOrder {
		if eps, ok := epMap[method]; ok {
			result = append(result, MethodEndpoints{
				Method:    method,
				Endpoints: eps,
			})
		}
	}

	return result
}
