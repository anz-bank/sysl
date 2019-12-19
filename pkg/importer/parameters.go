package importer

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
