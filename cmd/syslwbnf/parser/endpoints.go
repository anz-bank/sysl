package parser

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/wbnf/ast"
)

func (p *pscope) buildEndpoint(node EndpointNode) *sysl.Endpoint {
	if x := node.OneRestEndpoint(); x != nil {
		return p.buildRestEndpoint(x, "")
	} else if x := node.OneSimpleEndpoint(); x != nil {
		return p.buildSimpleEndpoint(x, "")
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

func (p *pscope) buildRestEndpoint(node *RestEndpointNode, pathPrefix string) *sysl.Endpoint {
	ep := &sysl.Endpoint{}

	if attrs := node.OneAttribs(); attrs != nil {
		ep.Attrs = p.buildAttributes(*attrs)
	}

	ep.RestParams = &sysl.Endpoint_RestParams{
		Path:   buildHttpPath(pathPrefix, *node.OneHttpPath()),
		Method: sysl.Endpoint_RestParams_Method(sysl.Endpoint_RestParams_Method_value["GET"]),
	}

	ep.Name = fmt.Sprintf("%s %s", ep.RestParams.Method, ep.RestParams.Path)

	return ep
}

func (p *pscope) buildSimpleEndpoint(node *SimpleEndpointNode, pathPrefix string) *sysl.Endpoint {
	ep := &sysl.Endpoint{}

	if attrs := node.OneAttribs(); attrs != nil {
		ep.Attrs = p.buildAttributes(*attrs)
	}

	if qs := node.OneQstring(); qs != nil {
		ep.LongName = qs.String()
	}
	name := appName(*node.OneEndpointName())
	ep.Name = name.Part[0] //fixme
	ep.Stmt = []*sysl.Statement{}
	ep.SourceContext = buildSourceContext(node.Node)
	for _, stmt := range node.AllStmt() {
		ep.Stmt = append(ep.Stmt, p.buildStatement(stmt))
	}

	return ep
}
