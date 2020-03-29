package parser

import (
	"github.com/anz-bank/sysl/pkg/sysl"
)

type EndpointMap map[string]*sysl.Endpoint

type epScope struct {
	EndpointMap
	*pscope
}

func (p *epScope) add(ep *sysl.Endpoint) {
	p.EndpointMap[ep.Name] = ep
}

func (p *pscope) buildEndpointMap(nodes []EndpointNode) EndpointMap {
	scope := &epScope{
		EndpointMap: EndpointMap{},
		pscope:      p,
	}
	for _, node := range nodes {
		scope.buildEndpoint(node)
	}
	return scope.EndpointMap
}

func (p *epScope) buildEndpoint(node EndpointNode) {
	if x := node.OneRestEndpoint(); x != nil {
		p.buildRestEndpoint(x, nil)
	} else if x := node.OneSimpleEndpoint(); x != nil {
		p.buildSimpleEndpoint(x)
	} else if x := node.OneCollector(); x != nil {
	} else if x := node.OneEvent(); x != nil {
	}
}

func (p *epScope) buildSimpleEndpoint(node *SimpleEndpointNode) {
	ep := &sysl.Endpoint{}

	if attrs := node.OneAttribs(); attrs != nil {
		ep.Attrs = p.buildAttributes(*attrs)
	}

	if qs := node.OneQstring(); qs != nil {
		ep.LongName = qs.String()
	}
	name := appName(*node.OneEndpointName())
	ep.Name = name.Part[0] //fixme
	ep.SourceContext = buildSourceContext(node.Node)
	for _, stmt := range node.AllStmt() {
		ep.Stmt = append(ep.Stmt, p.buildStatement(stmt))
	}

	p.add(ep)
}
