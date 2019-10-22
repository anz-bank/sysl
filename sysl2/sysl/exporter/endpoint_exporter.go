package exporter

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	proto "github.com/anz-bank/sysl/src/proto"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

type EndpointExporter struct {
	typeEx *TypeExporter
	log    *logrus.Logger
}

func makeEndpointExporter(typeEx *TypeExporter, logger *logrus.Logger) *EndpointExporter {
	return &EndpointExporter{
		typeEx: typeEx,
		log:    logger,
	}
}

const (
	object = "object"
)

func (e *EndpointExporter) exportChildStmts(returnStatusMap map[int]spec.Response, endpoint *proto.Endpoint) {
	regex := regexp.MustCompile(`^\d{3}$`)
	var retValues []string
	for _, stmt := range endpoint.GetStmt() {
		if ret, ok := stmt.Stmt.(*proto.Statement_Ret); ok {
			retValues = strings.Split(ret.Ret.GetPayload(), " <: ")
			res := &spec.Response{}
			res.Schema = &spec.Schema{}
			res.Schema.SchemaProps = spec.SchemaProps{}
			res.Schema.ExtraProps = map[string]interface{}{}
			hasStatusCode := regex.MatchString(retValues[0])
			switch {
			case len(retValues) > 1:
				status, err := strconv.Atoi(retValues[0])
				if err != nil {
					// log and ignore type
					e.log.Warnf("Type matching failed %s", err)
					continue
				}
				res.ResponseProps.Description = http.StatusText(status)
				if e.typeEx.isCompositeString(retValues[1]) {
					str := strings.Split(retValues[1], " ")
					res.ResponseProps.Schema.ExtraProps["$ref"] = "#/definitions/" + str[len(str)-1]
				} else {
					res.ResponseProps.Schema.ExtraProps["$ref"] = "#/definitions/" + retValues[1]
				}
				returnStatusMap[status] = *res
			case hasStatusCode:
				status, err := strconv.Atoi(retValues[0])
				if err != nil {
					e.log.Warnf("Error converting return code %s", err)
					continue
				}
				res.ResponseProps.Description = http.StatusText(status)
				returnStatusMap[status] = *res
			default:
				res.ResponseProps.Description = http.StatusText(200)
				if e.typeEx.isCompositeString(retValues[0]) {
					str := strings.Split(retValues[0], " ")
					res.ResponseProps.Schema.ExtraProps["$ref"] = "#/definitions/" + str[len(str)-1]
				} else {
					res.ResponseProps.Schema.ExtraProps["$ref"] = "#/definitions/" + retValues[0]
				}
				returnStatusMap[200] = *res
			}
		}
	}
}

func (e *EndpointExporter) populateEndpoint(path string, endpoint *proto.Endpoint, paths map[string]spec.PathItem,
) error {
	// extract the endpoint info and populate spec.PathItem
	var pathItem spec.PathItem
	var pathExists bool
	if pathItem, pathExists = paths[strings.Split(path, " ")[1]]; !pathExists {
		pathItem = spec.PathItem{}
	}
	returnStatusMap := map[int]spec.Response{}
	e.exportChildStmts(returnStatusMap, endpoint)
	op := e.setHTTPMethod(path, endpoint, &pathItem)
	op.Description = endpoint.GetDocstring()
	op.Summary = endpoint.GetDocstring()
	op.Responses.StatusCodeResponses = returnStatusMap
	op.Produces = []string{"application/json"}
	op.Consumes = []string{"application/json"}

	pathParamError := e.setPathParams(endpoint, op)
	if pathParamError != nil {
		return pathParamError
	}
	queryParamError := e.setQueryParams(endpoint, op)
	if queryParamError != nil {
		return queryParamError
	}
	paramError := e.setOtherParams(endpoint, op)
	if paramError != nil {
		return paramError
	}
	paths[strings.Split(path, " ")[1]] = pathItem
	return nil
}

func (e *EndpointExporter) setHTTPMethod(path string, endpoint *proto.Endpoint, pathItem *spec.PathItem,
) *spec.Operation {
	endpointTokens := strings.Split(path, " ")
	op := &spec.Operation{}
	op.Description = endpoint.GetLongName()
	op.Responses = &spec.Responses{}
	switch endpointTokens[0] {
	case `GET`:
		pathItem.PathItemProps.Get = op
	case `POST`:
		pathItem.PathItemProps.Post = op
	case `PUT`:
		pathItem.PathItemProps.Put = op
	case `DELETE`:
		pathItem.PathItemProps.Delete = op
	case `PATCH`:
		pathItem.PathItemProps.Patch = op
	}
	return op
}

func (e *EndpointExporter) setCommonAttributes(
	name string,
	param *spec.Parameter,
	attrMap map[string]*proto.Attribute,
	valueMap protoType,
) {
	param.Format = valueMap.Format
	param.Type = valueMap.Type
	param.Name = name
	if _, ok := attrMap["required"]; ok {
		param.Required = true
	} else {
		param = param.AsOptional()
	}
	if param.Type == object {
		param.Schema.ExtraProps["$ref"] = "#/definitions/" + param.Format
	}
}

func (e *EndpointExporter) setPathParams(endpoint *proto.Endpoint, op *spec.Operation) error {
	var attrMap map[string]*proto.Attribute
	for _, inParam := range endpoint.GetRestParams().GetUrlParam() {
		attrMap = inParam.GetType().GetAttrs()
		param := spec.PathParam(inParam.GetName())
		valueMap, err := e.typeEx.findSwaggerType(inParam.GetType())
		if err != nil {
			e.log.Warnf("Setting path params failed %s", err)
			return err
		}
		e.setCommonAttributes(inParam.GetName(), param, attrMap, valueMap)
		op.Parameters = append(op.Parameters, *param)
	}
	return nil
}

func (e *EndpointExporter) setQueryParams(endpoint *proto.Endpoint, op *spec.Operation) error {
	var attrMap map[string]*proto.Attribute
	for _, inParam := range endpoint.GetRestParams().GetQueryParam() {
		attrMap = inParam.GetType().GetAttrs()
		param := spec.QueryParam(inParam.GetName())
		valueMap, err := e.typeEx.findSwaggerType(inParam.GetType())
		if err != nil {
			e.log.Warnf("Setting query params failed %s", err)
			return err
		}
		e.setCommonAttributes(inParam.GetName(), param, attrMap, valueMap)
		op.Parameters = append(op.Parameters, *param)
	}
	return nil
}

func (e *EndpointExporter) setOtherParams(endpoint *proto.Endpoint, op *spec.Operation) error {
	var attrMap map[string]*proto.Attribute
	for _, inParam := range endpoint.GetParam() {
		attrMap = inParam.GetType().GetAttrs()
		param := &spec.Parameter{}
		valueMap, err := e.typeEx.findSwaggerType(inParam.GetType())
		if err != nil {
			e.log.Warnf("Setting header params failed %s", err)
			return err
		}
		if _, ok := attrMap["header"]; ok {
			param = spec.HeaderParam(attrMap["name"].GetS())
		} else if _, ok := attrMap["body"]; ok {
			param = spec.BodyParam(attrMap["name"].GetS(), param.ParamProps.Schema)
		} else if _, ok := attrMap["array"]; ok {
			param = spec.SimpleArrayParam(attrMap["name"].GetS(), valueMap.Type, valueMap.Format)
		} else {
			param = spec.HeaderParam(attrMap["name"].GetS())
		}
		if _, ok := attrMap["required"]; ok {
			param.Required = true
		} else {
			param = param.AsOptional()
		}
		param.Format = valueMap.Format
		param.Type = valueMap.Type
		if param.Type == object {
			param.ParamProps.Schema = &spec.Schema{}
			param.ParamProps.Schema.ExtraProps = map[string]interface{}{
				"$ref": "#/definitions/" + param.Format}
		}
		op.Parameters = append(op.Parameters, *param)
	}
	return nil
}
