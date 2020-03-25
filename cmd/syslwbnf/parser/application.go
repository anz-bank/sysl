package parser

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func buildApplication(node ApplicationNode) (string, *sysl.Application) {
	app := &sysl.Application{}
	app.Name = appName(*node.OneAppname())

	if qs := node.OneQstring(); qs != nil {
		app.LongName = qs.String()
	}
	if attrs := node.OneAttribs(); attrs != nil {
		app.Attrs = buildAttributes(*attrs)
	}

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
