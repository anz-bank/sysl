package parser

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func (p *pscope) buildApplication(node ApplicationNode) (string, *sysl.Application) {
	app := &sysl.Application{}
	app.Name = appName(*node.OneAppname())
	app.Endpoints = map[string]*sysl.Endpoint{}
	app.SourceContext = buildSourceContext(node.Node)
	app.Types = map[string]*sysl.Type{}

	if qs := node.OneQstring(); qs != nil {
		app.LongName = qs.String()
	}
	if attrs := node.OneAttribs(); attrs != nil {
		app.Attrs = p.buildAttributes(*attrs)
	}

	WalkerOps{
		EnterAnnotationNode: func(node AnnotationNode) Stopper {
			return nil
		},
		EnterEndpointNode: func(node EndpointNode) Stopper {
			ep := p.buildEndpoint(node)
			app.Endpoints[ep.Name] = ep
			return nil
		},
		EnterTypeDeclNode: func(node TypeDeclNode) Stopper {
			name, t := p.buildType(node)
			app.Types[name] = t
			return NodeExiter
		},
		EnterFacadeNode: func(node FacadeNode) Stopper {
			if app.Wrapped == nil {
				app.Wrapped = p.buildFacade(node)
			} else {
				panic("only one facade supported per application")
			}
			return NodeExiter
		},
	}.Walk(node)

	return strings.Join(app.Name.Part, "::"), app
}

func appName(node AppnameNode) *sysl.AppName {
	res := &sysl.AppName{}
	WalkerOps{
		EnterAppnamePartNode: func(node AppnamePartNode) Stopper {
			switch node.Choice() {
			case 0:
				res.Part = append(res.Part, node.AllToken()[0])
			case 1:
				res.Part = append(res.Part, node.AllQstring()[0].String())
			default:
				panic("oops")
			}
			return nil
		},
	}.Walk(node)
	return res
}
