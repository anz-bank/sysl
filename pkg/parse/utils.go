package parse

import (
	"math"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

func makeExpr(ctx *sysl.SourceContext) *sysl.Expr {
	return &sysl.Expr{
		SourceContext: ctx,
	}
}

type sourceCtxHelper struct {
	filename string
}

func (s sourceCtxHelper) Get(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	return s.get(ctx.GetStart(), ctx.GetStop(), false)
}

func (s sourceCtxHelper) GetWithText(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	return s.get(ctx.GetStart(), ctx.GetStop(), true)
}

func (s sourceCtxHelper) get(start, end antlr.Token, withText bool) *sysl.SourceContext {
	var text string
	if withText {
		is := start.GetInputStream()
		endChar := end.GetStop()
		if endChar == 0 {
			endChar = is.Size() - 1
		}
		text = is.GetText(start.GetStart(), endChar)
	}

	return &sysl.SourceContext{
		File: s.filename,
		Start: &sysl.SourceContext_Location{
			Line: int32(start.GetLine()),
			Col:  int32(start.GetColumn()),
		},
		End: &sysl.SourceContext_Location{
			Line: int32(end.GetLine()),
			Col:  int32(end.GetColumn()),
		},
		Text: text,
	}
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

// Attempt to find the referenced referenced sysl.Type by walking the current scope
// Simplified Rules
// 1) if the Reference has an AppName ignore the scope context path, only attempt to resolve with the scopes AppName
// 2) If the Reference has no AppName then the full Scope's Appname must be walked up with the refed appname appended
//
func resolveTypeRef(module *sysl.Module, ref *sysl.ScopedRef) *sysl.Type {
	testScope := func(appname []string, path []string) *sysl.Type {
		appName := &sysl.AppName{Part: appname}
		if app, ok := module.Apps[syslutil.GetAppName(appName)]; ok {
			if t, ok := app.Types[strings.Join(path, ".")]; ok {
				return t
			}
		}
		return nil
	}
	var refAppName []string
	if ref.Ref.Appname != nil {
		refAppName = ref.Ref.Appname.Part
	}

	// Only check the scope path if the ref doesnt have an app name
	if len(refAppName) == 0 {
		for i := len(ref.Context.Path); i >= 0; i-- {
			ctx := make([]string, i)
			copy(ctx, ref.Context.Path[:i])
			path := append(ctx, ref.Ref.Path...)
			if t := testScope(ref.Context.Appname.Part, path); t != nil {
				return t
			}
		}
	}

	// now find the Apps which could possibly contain the correct type:
	// Scope: A::B.C Ref: X::Y.Z should check for:
	//   A::B::X.Y.Z
	//   A::X:Y.Z
	//   X::Y.Z
	for i := len(ref.Context.Appname.Part); i >= 0; i-- {
		ctx := make([]string, i)
		copy(ctx, ref.Context.Appname.Part[:i])
		appName := append(ctx, refAppName...)
		if t := testScope(appName, ref.Ref.Path); t != nil {
			return t
		}
	}
	return nil
}
