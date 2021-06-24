package importer

import "github.com/anz-bank/sysl/pkg/syslutil"

type Endpoint struct {
	Path        string
	Description string

	Params Parameters

	Responses []Response
}

// nolint:gochecknoglobals
var methodDisplayOrder = []string{
	syslutil.Method_GET,
	syslutil.Method_PUT,
	syslutil.Method_POST,
	syslutil.Method_DELETE,
	syslutil.Method_PATCH,
}

// Response is going to be either freetext or a type, or both
type Response struct {
	Text string
	Type Type
}

type Param struct {
	Field
	In string
}

type Parameters struct {
	items       map[string]Param
	insertOrder []string
}

func (p *Parameters) Add(param Param) {
	if p.items == nil {
		p.items = map[string]Param{}
	}
	if _, found := p.items[param.Name]; !found {
		p.insertOrder = append(p.insertOrder, param.Name)
	}
	p.items[param.Name] = param
}

func (p *Parameters) Extend(others Parameters) Parameters {
	res := Parameters{}
	for _, name := range p.insertOrder {
		res.Add(p.items[name])
	}

	for _, name := range others.insertOrder {
		res.Add(others.items[name])
	}
	return res
}

func (p Parameters) findParams(where string) []Param {
	var res []Param
	for _, name := range p.insertOrder {
		item := p.items[name]
		if item.In == where {
			res = append(res, item)
		}
	}
	return res
}
func (p Parameters) QueryParams() []Param {
	return p.findParams("query")
}
func (p Parameters) HeaderParams() []Param {
	return p.findParams("header")
}
func (p Parameters) BodyParams() []Param {
	return p.findParams("body")
}
func (p Parameters) PathParams() []Param {
	return p.findParams("path")
}
func (p Parameters) CookieParams() []Param {
	return p.findParams("cookie")
}
