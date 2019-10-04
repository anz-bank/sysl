package importer

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-openapi/spec"
)

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

func buildParam(param spec.Parameter, types TypeList, globals Parameters, logger *logrus.Logger) (Param, error) {
	fromType := func(t Type) Param {
		return Param{
			Field: Field{
				Name:     param.Name,
				Optional: !param.Required,
				Type:     t,
			},
			In: param.In,
		}
	}
	fromString := func(typeName string) Param {
		t, found := types.Find(typeName)
		if !found {
			logger.Errorf("referenced parameter type %s not found\n", typeName)
			t = NewStringAlias(typeName)
		}
		return fromType(t)
	}

	paramTypeName := param.Type
	if paramTypeName == "" {
		if param.Schema != nil {
			if ptype, ok := types.FindFromSchema(*param.Schema, &typeData{logger: logger}); ok {
				return fromType(ptype), nil
			}
			return Param{}, fmt.Errorf("referenced parameter type not found: %+v", param.Schema)
		} else if refURL := param.Ref.GetURL(); refURL != nil {
			refParamName := getReferenceFragment(refURL)
			if p, ok := globals.items[refParamName]; ok {
				return p, nil
			}
			return Param{}, fmt.Errorf("referenced parameter %s unknown", refParamName)
		}
	} else if paramTypeName == "string" {
		return fromString(paramTypeName), nil
	}
	return fromString(mapSwaggerTypeAndFormatToType(paramTypeName, param.Format, logger)), nil
}

func buildParameters(params []spec.Parameter, types TypeList,
	globals Parameters, baseParams Parameters, logger *logrus.Logger) Parameters {
	res := baseParams
	for _, param := range params {
		p, err := buildParam(param, types, globals, logger)
		if err != nil {
			logger.Errorf("%s", err.Error())
		} else {
			res.Add(p)
		}
	}

	return res
}

func buildGlobalParams(params map[string]spec.Parameter, types TypeList, logger *logrus.Logger) Parameters {
	res := Parameters{}
	for key, param := range params {
		p, err := buildParam(param, types, res, logger)
		if err != nil {
			logger.Errorf("%s\n", err.Error())
		} else {
			p.Name = key
			res.Add(p)
		}
	}
	return res
}
