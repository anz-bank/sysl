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
	object   = "object"
	format   = "format"
	dataType = "type"
)

func (v *EndpointExporter) exportChildStmts(returnStatusMap map[int]spec.Response, endpoint *proto.Endpoint) {
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
					v.log.Warnf("Type matching failed %s", err)
					continue
				}
				res.ResponseProps.Description = http.StatusText(status)
				if v.typeEx.isCompositeString(retValues[1]) {
					str := strings.Split(retValues[1], " ")
					res.ResponseProps.Schema.ExtraProps["$ref"] = "#/definitions/" + str[len(str)-1]
				} else {
					res.ResponseProps.Schema.ExtraProps["$ref"] = "#/definitions/" + retValues[1]
				}
				returnStatusMap[status] = *res
			case hasStatusCode:
				status, err := strconv.Atoi(retValues[0])
				if err != nil {
					v.log.Warnf("Error converting return code %s", err)
					continue
				}
				res.ResponseProps.Description = http.StatusText(status)
				returnStatusMap[status] = *res
			default:
				res.ResponseProps.Description = http.StatusText(200)
				if v.typeEx.isCompositeString(retValues[0]) {
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

func (v *EndpointExporter) exportEndpoint(path string, endpoint *proto.Endpoint, paths map[string]spec.PathItem) error {
	// extract the endpoint info and populate spec.PathItem
	var pathItem spec.PathItem
	var pathExists bool
	if pathItem, pathExists = paths[strings.Split(path, " ")[1]]; !pathExists {
		pathItem = spec.PathItem{}
	}
	returnStatusMap := map[int]spec.Response{}
	v.exportChildStmts(returnStatusMap, endpoint)
	op := v.parseHTTPMethod(path, endpoint, &pathItem)
	op.Description = endpoint.GetDocstring()
	op.Summary = endpoint.GetDocstring()
	op.Responses.StatusCodeResponses = returnStatusMap
	op.Produces = []string{"application/json"}
	op.Consumes = []string{"application/json"}

	pathParamError := v.parsePathParams(endpoint, op)
	if pathParamError != nil {
		return pathParamError
	}
	queryParamError := v.parseQueryParams(endpoint, op)
	if queryParamError != nil {
		return queryParamError
	}
	paramError := v.parseOtherParams(endpoint, op)
	if paramError != nil {
		return paramError
	}
	paths[strings.Split(path, " ")[1]] = pathItem
	return nil
}

func (v *EndpointExporter) parseHTTPMethod(path string, endpoint *proto.Endpoint, pathItem *spec.PathItem,
) *spec.Operation {
	endpointTokens := strings.Split(path, " ")
	switch endpointTokens[0] {
	case `GET`:
		pathItem.PathItemProps.Get = &spec.Operation{}
		pathItem.PathItemProps.Get.Description = endpoint.GetLongName()
		pathItem.PathItemProps.Get.Responses = &spec.Responses{}
		return pathItem.PathItemProps.Get
	case `POST`:
		pathItem.PathItemProps.Post = &spec.Operation{}
		pathItem.PathItemProps.Post.Description = endpoint.GetLongName()
		pathItem.PathItemProps.Post.Responses = &spec.Responses{}
		return pathItem.PathItemProps.Post
	case `PUT`:
		pathItem.PathItemProps.Put = &spec.Operation{}
		pathItem.PathItemProps.Put.Description = endpoint.GetLongName()
		pathItem.PathItemProps.Put.Responses = &spec.Responses{}
		return pathItem.PathItemProps.Put
	case `DELETE`:
		pathItem.PathItemProps.Delete = &spec.Operation{}
		pathItem.PathItemProps.Delete.Description = endpoint.GetLongName()
		pathItem.PathItemProps.Delete.Responses = &spec.Responses{}
		return pathItem.PathItemProps.Delete
	case `PATCH`:
		pathItem.PathItemProps.Patch = &spec.Operation{}
		pathItem.PathItemProps.Patch.Description = endpoint.GetLongName()
		pathItem.PathItemProps.Patch.Responses = &spec.Responses{}
		return pathItem.PathItemProps.Patch
	}
	return nil
}

func (v *EndpointExporter) parseQueryParams(endpoint *proto.Endpoint, op *spec.Operation) error {
	for _, inParam := range endpoint.GetRestParams().GetUrlParam() {
		param := spec.PathParam(inParam.GetName())
		valueMap, err := v.typeEx.findSwaggerType(inParam.GetType())
		if err != nil {
			return err
		}
		param.Format = valueMap[format]
		param.Type = valueMap[dataType]
		param.Name = inParam.GetName()
		param.Required = true
		if param.Type == object {
			param.Schema.ExtraProps["$ref"] = "#/definitions/" + param.Format
		}
		op.Parameters = append(op.Parameters, *param)
	}
	return nil
}

func (v *EndpointExporter) parsePathParams(endpoint *proto.Endpoint, op *spec.Operation) error {
	for _, inParam := range endpoint.GetRestParams().GetQueryParam() {
		param := spec.QueryParam(inParam.GetName())
		valueMap, err := v.typeEx.findSwaggerType(inParam.GetType())
		if err != nil {
			return err
		}
		param.Format = valueMap[format]
		param.Type = valueMap[dataType]
		param.Name = inParam.GetName()
		param.Required = true
		if param.Type == object {
			param.Schema.ExtraProps["$ref"] = "#/definitions/" + param.Format
		}
		op.Parameters = append(op.Parameters, *param)
	}
	return nil
}

func (v *EndpointExporter) parseOtherParams(endpoint *proto.Endpoint, op *spec.Operation) error {
	var attrMap map[string]*proto.Attribute
	for _, inParam := range endpoint.GetParam() {
		attrMap = inParam.GetType().GetAttrs()
		param := &spec.Parameter{}
		valueMap, err := v.typeEx.findSwaggerType(inParam.GetType())
		if err != nil {
			return err
		}
		if _, ok := attrMap["header"]; ok {
			param = spec.HeaderParam(attrMap["name"].GetS())
		} else if _, ok := attrMap["body"]; ok {
			param = spec.BodyParam(attrMap["name"].GetS(), param.ParamProps.Schema)
		} else if _, ok := attrMap["array"]; ok {
			param = spec.SimpleArrayParam(attrMap["name"].GetS(), valueMap["type"], valueMap["format"])
		} else {
			param = spec.HeaderParam(attrMap["name"].GetS())
		}
		if _, ok := attrMap["required"]; ok {
			param.Required = true
		} else {
			param = param.AsOptional()
		}
		param.Format = valueMap[format]
		param.Type = valueMap[dataType]
		if param.Type == object {
			param.ParamProps.Schema = &spec.Schema{}
			param.ParamProps.Schema.ExtraProps = map[string]interface{}{
				"$ref": "#/definitions/" + param.Format}
		}
		op.Parameters = append(op.Parameters, *param)
	}
	return nil
}
