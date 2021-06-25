package parse

import (
	"math"
	"net/url"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

func makeExpr(ctx *sysl.SourceContext) *sysl.Expr {
	return &sysl.Expr{
		SourceContext:  ctx,
		SourceContexts: []*sysl.SourceContext{ctx},
	}
}

type sourceCtxHelper struct {
	filename string
	version  string
}

func (s *sourceCtxHelper) get(start, end antlr.Token) *sysl.SourceContext {
	// In Antlr, lines are 1-based, columns are 0-based.
	// In Sysl SourceContext, both lines and columns are 0-based.
	ctx := &sysl.SourceContext{
		File: s.filename,
		Start: &sysl.SourceContext_Location{
			Line: int32(start.GetLine() - 1),
			Col:  int32(start.GetColumn()),
		},
		End: &sysl.SourceContext_Location{
			Line: int32(end.GetLine() - 1),
			Col:  int32(end.GetColumn()),
		},
		Version: s.version,
	}
	if text := end.GetText(); text != "" {
		ctx.End.Col += int32(len(text))
	}
	return ctx
}

// lastToken returns the last token in the AST.
func lastToken(tree antlr.Tree) antlr.Token {
	switch t := tree.(type) {
	case antlr.TerminalNode:
		return t.GetSymbol()
	default:
		return lastToken(tree.GetChild(tree.GetChildCount() - 1))
	}
}

// getBitWidth returns the bit width of the type, or 0 if no such constraint exists.
func getBitWidth(t *sysl.Type) int32 {
	if t.Constraint != nil {
		for _, item := range t.Constraint {
			if item != nil && item.BitWidth > 0 {
				return item.BitWidth
			}
		}
	}

	return 0
}

func primitiveFromNativeDataType(native antlr.TerminalNode) (*sysl.Type_Primitive_, *sysl.Type_Constraint) {
	if native == nil {
		return nil, nil
	}
	text := strings.ToUpper(native.GetText())
	primitiveType := sysl.Type_Primitive(sysl.Type_Primitive_value[text])
	if primitiveType != sysl.Type_NO_Primitive {
		return &sysl.Type_Primitive_{Primitive: primitiveType}, nil
	}

	var constraint *sysl.Type_Constraint
	switch text {
	case "INT32":
		primitiveType = sysl.Type_Primitive(sysl.Type_Primitive_value["INT"])
		constraint = &sysl.Type_Constraint{
			BitWidth: 32,
			Range: &sysl.Type_Constraint_Range{
				Min: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MinInt32,
					},
				},
				Max: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MaxInt32,
					},
				},
			},
		}

	case "INT64":
		primitiveType = sysl.Type_Primitive(sysl.Type_Primitive_value["INT"])
		constraint = &sysl.Type_Constraint{
			BitWidth: 64,
			Range: &sysl.Type_Constraint_Range{
				Min: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MinInt64,
					},
				},
				Max: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MaxInt64,
					},
				},
			},
		}

	case "FLOAT32":
		primitiveType = sysl.Type_Primitive(sysl.Type_Primitive_value["FLOAT"])
		constraint = &sysl.Type_Constraint{BitWidth: 32}
	case "FLOAT64":
		primitiveType = sysl.Type_Primitive(sysl.Type_Primitive_value["FLOAT"])
		constraint = &sysl.Type_Constraint{BitWidth: 64}
	}
	return &sysl.Type_Primitive_{Primitive: primitiveType}, constraint
}

type PathStack struct {
	sep   string
	parts []string

	path   string
	prefix string
}

func NewPathStack(sep string) PathStack {
	return PathStack{
		sep:   sep,
		parts: []string{},
	}
}

func (p PathStack) Get() string {
	return p.path
}

func (p *PathStack) Push(items ...string) string {
	for _, i := range items {
		p.parts = append(p.parts, strings.TrimSpace(i))
	}
	return p.update()
}

func (p *PathStack) Pop() string {
	p.parts = p.parts[:len(p.parts)-1]
	return p.update()
}

func (p *PathStack) Reset() string {
	p.parts = []string{}
	return p.update()
}

// Parts() clones the path items into a new slice so that it is safe to store
func (p PathStack) Parts() []string {
	return append([]string{}, p.parts...)
}

func (p *PathStack) update() string {
	p.path = p.prefix + strings.Join(p.parts, p.sep)
	return p.path
}

// MustUnescape will escape URL escaped characters.
func MustUnescape(endpoint string) string {
	s, err := url.PathUnescape(endpoint)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(s)
}
