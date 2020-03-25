package parser

import (
	"sort"

	"github.com/arr-ai/wbnf/ast"
)

func Leafs(node ast.Node) []ast.Leaf {
	vals := walkNodes(node)

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Scanner().Offset() < vals[j].Scanner().Offset()
	})

	return vals
}

func walkNodes(node ast.Node) []ast.Leaf {
	var vals []ast.Leaf
	switch n := node.(type) {
	case ast.Branch:
		for _, val := range n {
			switch c := val.(type) {
			case ast.One:
				vals = append(vals, walkNodes(c.Node)...)
			case ast.Many:
				for _, c := range c {
					vals = append(vals, walkNodes(c)...)
				}
			}
		}
	case ast.Leaf:
		return []ast.Leaf{n}
	}
	return vals
}
