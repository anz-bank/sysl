package parser

import (
	"strconv"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func buildAttributes(node AttribsNode) map[string]*sysl.Attribute {
	attrs := map[string]*sysl.Attribute{}
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
			var patterns []*sysl.Attribute
			for _, n := range node.AllName() {
				patterns = append(patterns, &sysl.Attribute{
					Attribute: &sysl.Attribute_S{
						S: n.String(),
					},
				})
			}
			attrs["patterns"] = &sysl.Attribute{
				Attribute: &sysl.Attribute_A{
					A: &sysl.Attribute_Array{
						Elt: patterns,
					},
				},
			}
			return nil
		}}.Walk(node)
	return attrs
}
