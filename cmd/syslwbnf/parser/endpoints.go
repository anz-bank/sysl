package parser

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/wbnf/ast"
)

func buildEndpoint(node EndpointNode) *sysl.Endpoint {
	if x := node.OneRestEndpoint(); x != nil {
		return buildRestEndpoint(x, "")
	} else if x := node.OneSimpleEndpoint(); x != nil {
		return buildSimpleEndpoint(x, "")
	} else if x := node.OneCollector(); x != nil {
	} else if x := node.OneEvent(); x != nil {
	}

	return nil
}

func buildHttpPath(prefix string, node HttpPathNode) string {
	var parts []string
	if prefix != "" {
		parts = append(parts, prefix)
	}
	for _, p := range ast.Leafs(node.Node) {
		parts = append(parts, p.Scanner().String())
	}
	return strings.Join(parts, "/")
}

func buildRestEndpoint(node *RestEndpointNode, pathPrefix string) *sysl.Endpoint {
	ep := &sysl.Endpoint{}

	if attrs := node.OneAttribs(); attrs != nil {
		ep.Attrs = buildAttributes(*attrs)
	}

	ep.RestParams = &sysl.Endpoint_RestParams{
		Path:   buildHttpPath(pathPrefix, *node.OneHttpPath()),
		Method: sysl.Endpoint_RestParams_Method(sysl.Endpoint_RestParams_Method_value["GET"]),
	}

	ep.Name = fmt.Sprintf("%s %s", ep.RestParams.Method, ep.RestParams.Path)

	return ep
}

func buildSimpleEndpoint(node *SimpleEndpointNode, pathPrefix string) *sysl.Endpoint {
	ep := &sysl.Endpoint{}

	if attrs := node.OneAttribs(); attrs != nil {
		ep.Attrs = buildAttributes(*attrs)
	}

	if qs := node.OneQstring(); qs != nil {
		ep.LongName = qs.String()
	}
	ep.Name = node.OneEndpointName().String()

	return ep
}
