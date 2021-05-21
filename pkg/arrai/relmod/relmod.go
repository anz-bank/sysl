package relmod

import (
	"context"
	"fmt"

	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/wbnf/parser"
	"github.com/pkg/errors"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/arrai/rel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const tagAttr = "patterns"
const placeholder = "..."

type ctxKey int

const (
	payloadParseKey ctxKey = iota
)

// NormalizeSpec returns a relational form of a Sysl module loaded from a Sysl spec.
func NormalizeSpec(ctx context.Context, root, path string) (*Schema, error) {
	m, _, err := loader.LoadSyslModule(root, path, afero.NewOsFs(), logrus.StandardLogger())
	if err != nil {
		return nil, err
	}
	return Normalize(ctx, m)
}

// withPayloadParser precompiles the arr.ai function for parsing return statement payloads.
func withPayloadParser(ctx context.Context) (context.Context, error) {
	expr, err := buildPayloadParser()
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, payloadParseKey, expr), nil
}

// getPayloadParser returns the payload parser function stored in ctx.
func getPayloadParser(ctx context.Context) rel.Expr {
	return ctx.Value(payloadParseKey).(*rel.BinExpr)
}

// attrToValue returns the attribute value as a rel.Value.
//
// If the attribute value contains any arrays, they will be serialized as arrays, rather than sets.
func attrToValue(a *sysl.Attribute) rel.Value {
	switch a.Attribute.(type) {
	case *sysl.Attribute_S:
		return rel.NewString([]rune(a.GetS()))
	case *sysl.Attribute_I:
		return rel.NewNumber(float64(a.GetI()))
	case *sysl.Attribute_N:
		return rel.NewNumber(a.GetN())
	case *sysl.Attribute_A:
		as := a.GetA().Elt
		vs := make([]rel.Value, 0, len(as))
		for _, elt := range as {
			vs = append(vs, attrToValue(elt))
		}
		return rel.NewArray(vs...)
	default:
		panic(fmt.Errorf(fmt.Sprintf("unknown attr type: %x", a)))
	}
}

// tags extracts the tags from a set of attributes.
func tags(attrs map[string]*sysl.Attribute) []string {
	var tags []string
	for attrName, attr := range attrs {
		if attrName == tagAttr {
			if _, ok := attr.GetAttribute().(*sysl.Attribute_A); !ok {
				panic(fmt.Errorf(fmt.Sprintf("patterns attr not an array: %x", attr)))
			}
			for _, elt := range attr.GetA().Elt {
				if _, ok := elt.GetAttribute().(*sysl.Attribute_S); !ok {
					panic(fmt.Errorf(fmt.Sprintf("pattern value not a string: %x", elt)))
				}
				tags = append(tags, elt.GetS())
			}
		}
	}
	return tags
}

// annos extracts the annotations from a set of attributes.
func annos(attrs map[string]*sysl.Attribute) map[string]interface{} {
	annos := map[string]interface{}{}
	for name, attr := range attrs {
		if name == tagAttr {
			continue
		}
		annos[name] = attrToValue(attr)
	}
	return annos
}

// buildPayloadParser constructs an expression to parse return statement payloads.
func buildPayloadParser() (rel.Expr, error) {
	ctx := arraictx.InitRunCtx(context.Background())
	parse := `
		\payload //grammar.parse({://grammar.lang.wbnf:
			payload -> (status ("<:" type)? | (status "<:")? type) attr?;
			type		-> sequence | set | PRIMITIVE | ref $ | raw $;
			sequence	-> "sequence of " type;
			set			-> "set of " type;
			ref			-> (app=([^\s.:]+):"::" ".")? type=[^\s.]+;
			raw			-> [^\[\n]+\b;
			PRIMITIVE	-> "int" | "int32" | "int64" | "float" | "float32" | "float64" | "decimal"
						 | "bool" | "bytes" | "string" | "date" | "datetime" | "any";
			status -> ("ok"|"error"|[1-5][0-9][0-9]);
			attr -> %!Array(nvp|modifier);
			nvp_item -> str | array=%!Array(nvp_item) | dict=%!Dict(nvp_item);
			nvp ->  name=\w+ "=" nvp_item;
			modifier -> "~" name=[\w\+]+;
			str -> ('"' ([^"\\] | [\\][\\brntu'"])* '"' | "'" ([^''])* "'") {
				.wrapRE -> /{()};
			};
			.wrapRE -> /{\s*()\s*};
			.macro Array(child) {
				"[" (child):"," "]"
			}
			.macro Dict(child) {
				"{" entry=(key=child ":" value=child):"," "}"
			}
		:}, "payload", payload)
	`
	tx := `
		\ast
		let rec buildNvp = \nvp cond nvp {
			(array: (nvp_item: i, ...), ...): (a: i => (:.@, @item: buildNvp(.@item))),
			(dict: (entry: i, ...), ...): (d: i => (
				@     : buildNvp(.@item.key.nvp_item),
				@value: buildNvp(.@item.value.nvp_item)
			)),
			(str: ('': s, ...), ...): //eval.value(//seq.join('', s)),
			_: //eval.value(//seq.join('', nvp.''))
		};
		let rec type = \t cond t {
			(:set, ...): (set: type(set.type)),
			(:sequence, ...): (sequence: type(sequence.type)),
			(:PRIMITIVE, ...): (primitive: PRIMITIVE.'' rank (:.@)),
			(:ref, ...): (
				appName: ref.app?:[] >> (.'' rank (:.@)),
				typePath: [ref.type.'' rank (:.@)],
			),
			(:raw, ...): (primitive: 'any'), # TODO: encode raw.'' rank (:.@)
			_: t,
		};
		(
			status: ast.status?.'':'' rank (:.@),
			type: type(ast.type?:()),
			nvp: ast.attr?.nvp?:{} => (@: (.@item.name.'' rank (:.@)), @value: buildNvp(.@item.nvp_item)),
			modifier: ast.attr?.modifier?:{} => (.@item.name.'' rank (:.@))
		)
	`

	txFn, err := syntax.EvaluateExpr(ctx, ".", tx)
	if err != nil {
		return nil, err
	}
	if _, is := txFn.(rel.Closure); !is {
		return nil, errors.New("tx not a function")
	}

	parseFn, err := syntax.EvaluateExpr(ctx, ".", parse)
	if err != nil {
		return nil, err
	}
	if _, is := parseFn.(rel.Closure); !is {
		return nil, errors.New("parse not a function")
	}

	wrapFn, err := syntax.EvaluateExpr(ctx, ".", `\parse \tx \payload tx(parse(payload))`)
	if err != nil {
		return nil, err
	}
	if _, is := parseFn.(rel.Closure); !is {
		return nil, errors.New("wrap not a function")
	}

	scanner := *parser.NewScanner("")
	return rel.NewCallExprCurry(scanner, wrapFn, parseFn, txFn), nil
}

func unpackType(typ rel.Tuple, appName []string) interface{} {
	tuple := typ.(rel.Tuple)
	//nolint:gocritic
	if name, ok := tuple.Get("primitive"); ok {
		return TypePrimitive{Primitive: name.String()}
	} else if set, ok := tuple.Get("set"); ok {
		return TypeSet{unpackType(set.(rel.Tuple), appName)}
	} else if seq, ok := tuple.Get("sequence"); ok {
		return TypeSequence{unpackType(seq.(rel.Tuple), appName)}
	} else if tuple.HasName("typePath") {
		ctx := context.Background()
		name := arrai.ToStrings(tuple.MustGet("appName").Export(ctx))
		if len(name) == 0 {
			name = appName
		}
		return TypeRef{
			AppName:  name,
			TypePath: arrai.ToStrings(tuple.MustGet("typePath").Export(ctx)),
		}
	} else {
		panic(fmt.Errorf("unknown type: %T %s", typ, typ))
	}
}

func parseReturnPayload(ctx context.Context, payload string, appName []string) (StatementReturn, error) {
	parse := getPayloadParser(ctx)

	scanner := *parser.NewScanner("")
	out, err := rel.NewCallExpr(scanner, parse, rel.NewString([]rune(payload))).Eval(ctx, rel.EmptyScope)
	if err != nil {
		return StatementReturn{}, err
	}
	t := out.(rel.Tuple)

	r := StatementReturn{
		Status: t.MustGet("status").String(),
		Attr: StatementReturnAttrs{
			Modifier: arrai.ToStrings(t.MustGet("modifier").Export(ctx)),
			Nvp:      arrai.ToInterfaceMap(t.MustGet("nvp").Export(ctx)),
		},
	}

	if typ, ok := t.Get("type"); ok && typ.IsTrue() {
		r.Type = unpackType(typ.(rel.Tuple), appName)
	}
	return r, nil
}

func parseFieldType(appName []string, t *sysl.Type) interface{} {
	switch t.Type.(type) {
	case *sysl.Type_Primitive_:
		return TypePrimitive{Primitive: t.GetPrimitive().String()}
	case *sysl.Type_Tuple_:
		return TypeTuple{Tuple: t.GetTuple()}
	case *sysl.Type_TypeRef:
		ref := t.GetTypeRef()
		if ref.Ref.Appname != nil {
			return TypeRef{AppName: ref.Ref.Appname.Part, TypePath: ref.Ref.Path}
		} else if ref.Context != nil {
			return TypeRef{AppName: ref.Context.Appname.Part, TypePath: ref.Ref.Path}
		}
		return TypeRef{AppName: appName, TypePath: ref.Ref.Path}
	case *sysl.Type_Set:
		ft := parseFieldType(appName, t.GetSet())
		return TypeSet{Set: ft}
	case *sysl.Type_Sequence:
		ft := parseFieldType(appName, t.GetSequence())
		return TypeSequence{Sequence: ft}
	case *sysl.Type_NoType_:
		return nil
	case *sysl.Type_List_:
		return parseFieldType(appName, t.GetList().Type)
	default:
		return nil
	}
}

func parseRestPath(p *sysl.Endpoint_RestParams) (Rest, error) {
	if p == nil {
		return Rest{}, nil
	}
	return Rest{
		Method: p.Method.String(),
		Path:   p.Path,
	}, nil
}

func relmodSourceContexts(contexts []*sysl.SourceContext) []SourceContext {
	var srcs []SourceContext
	for _, context := range contexts {
		srcs = append(srcs, relmodSourceContext(context))
	}
	return srcs
}

func relmodSourceContext(context *sysl.SourceContext) SourceContext {
	return SourceContext{
		File: context.File,
		Start: Location{
			Line: int(context.Start.Line),
			Col:  int(context.Start.Col),
		},
		End: Location{
			Line: int(context.End.Line),
			Col:  int(context.End.Col),
		},
	}
}
