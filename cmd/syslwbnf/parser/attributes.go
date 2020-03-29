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

func addPattern(attrs AttributeMap, which string) {
	pats, ok := attrs[patternsKey]
	if !ok {
		attrs[patternsKey] = &sysl.Attribute{
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: []*sysl.Attribute{{Attribute: &sysl.Attribute_S{S: which}}},
				},
			},
		}
	} else {
		attrs := pats.Attribute.(*sysl.Attribute_A).A.Elt
		for _, elt := range attrs {
			if s, ok := elt.Attribute.(*sysl.Attribute_S); ok && s.S == which {
				return
			}
		}
		pats.Attribute.(*sysl.Attribute_A).A.Elt = append(attrs, &sysl.Attribute{Attribute: &sysl.Attribute_S{S: which}})
	}
}

func mergeAttrs(src AttributeMap, dst AttributeMap) {
	for k, v := range src {
		if _, has := dst[k]; !has {
			dst[k] = v
		} else {
			dstAttr, dstOK := dst[k].Attribute.(*sysl.Attribute_A)
			vAttr, vOK := v.Attribute.(*sysl.Attribute_A)
			if dstOK && vOK {
				dstAttr.A.Elt = append(dstAttr.A.Elt, vAttr.A.Elt...)
			} else {
				dst[k] = v
			}
		}
	}
}
