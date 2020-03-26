package parser

import (
	"strconv"

	"github.com/anz-bank/sysl/pkg/sysl"
)

type AttributeMap map[string]*sysl.Attribute

const patternsKey = "patterns"

func (p *pscope) buildAttributes(node AttribsNode) AttributeMap {
	attrs := AttributeMap{}
	var patterns []*sysl.Attribute
	WalkerOps{
		EnterAttribsAttrNode: func(node AttribsAttrNode) Stopper {
			s, _ := strconv.Unquote(node.OneQstring().String())
			attrs[node.OneName().String()] = &sysl.Attribute{
				Attribute: &sysl.Attribute_S{
					S: s,
				},
			}
			return nil
		},
		EnterAttribsPatternNode: func(node AttribsPatternNode) Stopper {
			for _, n := range node.AllName() {
				patterns = append(patterns, &sysl.Attribute{
					Attribute: &sysl.Attribute_S{
						S: n.String(),
					},
				})
			}
			return nil
		}}.Walk(node)

	if len(patterns) > 0 {
		attrs[patternsKey] = &sysl.Attribute{
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: patterns,
				},
			},
		}
	}
	return attrs
}

func hasPattern(attrs AttributeMap, which string) bool {
	if pats, ok := attrs[patternsKey]; ok {
		for _, elt := range pats.Attribute.(*sysl.Attribute_A).A.Elt {
			if s, ok := elt.Attribute.(*sysl.Attribute_S); ok && s.S == which {
				return true
			}
		}
	}
	return false
}
