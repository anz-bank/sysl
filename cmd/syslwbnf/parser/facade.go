package parser

import "github.com/anz-bank/sysl/pkg/sysl"

func (p *pscope) buildFacade(node FacadeNode) *sysl.Application {
	app := &sysl.Application{}
	app.Name = &sysl.AppName{
		Part: []string{node.OneName().String()},
	}
	app.SourceContext = buildSourceContext(node.Node)
	app.Types = map[string]*sysl.Type{}

	for _, child := range node.AllCovering() {
		t := &sysl.Type{}
		// todo: inplace_table_def
		app.Types[child.AllName()[0].String()] = t
	}

	return app
}
