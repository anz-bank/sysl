//nolint:revive,stylecheck,funlen
package parse // SyslParser

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	parser "github.com/anz-bank/sysl/pkg/grammar"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

const patternsKey = "patterns"

type importDef struct {
	filename string
	appname  string
	pkg      string
	mode     string
}

// TreeShapeListener ..
type TreeShapeListener struct {
	*parser.BaseSyslParserListener
	base         string
	sc           sourceCtxHelper
	imports      []importDef
	module       *sysl.Module
	endpointName string

	app             *sysl.Application
	currentTypePath PathStack
	viewName        string

	linter                *linterRecords
	fieldname             []string
	urlPrefixes           PathStack
	app_name              PathStack
	annotation_name       string
	typemap               map[string]*sysl.Type
	pendingDocString      bool
	lintMode              bool
	rest_attrs            []map[string]*sysl.Attribute
	rest_queryparams      []*sysl.Endpoint_RestParams_QueryParam
	rest_urlparams        []*sysl.Endpoint_RestParams_QueryParam
	method_urlparams      []*sysl.Endpoint_RestParams_QueryParam
	rest_queryparams_len  []int
	rest_urlparams_len    []int
	http_path_query_param string
	stmt_scope            []interface{} // Endpoint, if, if_else, loop
	expr_stack            []*sysl.Expr
	opmap                 map[string]sysl.Expr_BinExpr_Op
	currentMultiLineAnno  []string
	lastEnd               *sysl.SourceContext_Location
}

// NewTreeShapeListener ...
func NewTreeShapeListener() *TreeShapeListener {
	opmap := map[string]sysl.Expr_BinExpr_Op{
		"==":        sysl.Expr_BinExpr_EQ,
		"!=":        sysl.Expr_BinExpr_NE,
		"<":         sysl.Expr_BinExpr_LT,
		"<=":        sysl.Expr_BinExpr_LE,
		">":         sysl.Expr_BinExpr_GT,
		">=":        sysl.Expr_BinExpr_GE,
		"in":        sysl.Expr_BinExpr_IN,
		"contains":  sysl.Expr_BinExpr_CONTAINS,
		"!in":       sysl.Expr_BinExpr_NOT_IN,
		"!contains": sysl.Expr_BinExpr_NOT_CONTAINS,
	}

	return &TreeShapeListener{
		module: &sysl.Module{
			Apps: map[string]*sysl.Application{},
		},
		opmap:           opmap,
		currentTypePath: NewPathStack("."),
		urlPrefixes:     NewPathStack(""),
		app_name:        NewPathStack(" :: "),
	}
}

func (s *TreeShapeListener) currentApp() *sysl.Application {
	return s.app
}

func (s *TreeShapeListener) getSrcCtx(ctx antlr.ParserRuleContext) *sysl.SourceContext {
	return s.getSrcCtxFor(ctx.GetStart(), ctx.GetStop())
}

func (s *TreeShapeListener) getSrcCtxFor(start, end antlr.Token) *sysl.SourceContext {
	sc := s.sc.get(start, end)
	s.lastEnd = sc.End
	return sc
}

func (s *TreeShapeListener) getSrcCtxWithText(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	start := ctx.GetStart()
	end := ctx.GetStop()
	is := start.GetInputStream()
	endChar := end.GetStop()
	if endChar == 0 {
		endChar = is.Size() - 1
	}
	text := is.GetText(start.GetStart(), endChar)

	sc := s.getSrcCtx(ctx)
	sc.Text = text
	return sc
}

// EnterName_str is called when production name_str is entered.
func (s *TreeShapeListener) EnterName_str(ctx *parser.Name_strContext) {
	s.app_name.Push(ctx.GetText())
}

// EnterReference is called when production reference is entered.
func (s *TreeShapeListener) EnterReference(*parser.ReferenceContext) {
	context_app_part := s.currentApp().Name.Part
	context_path := s.currentTypePath.Parts()

	s.currentType().Type = &sysl.Type_TypeRef{
		TypeRef: &sysl.ScopedRef{
			Context: &sysl.Scope{
				Appname: &sysl.AppName{
					Part: context_app_part,
				},
				Path: context_path,
			},
		},
	}
	s.app_name.Reset()
}

// ExitReference is called when production reference is exited.
func (s *TreeShapeListener) ExitReference(ctx *parser.ReferenceContext) {
	s.currentType().GetTypeRef().Ref = makeScope(s.app_name.Parts(), ctx)
}

// EnterDoc_string is called when production doc_string is entered.
func (s *TreeShapeListener) EnterDoc_string(ctx *parser.Doc_stringContext) {
	if s.pendingDocString {
		s.pendingDocString = false

		space := ""
		text := ctx.TEXT().GetText()
		text = strings.ReplaceAll(text, `"`, `\"`)
		text = fromQString(`"` + text[1:] + `"`)

		if s.currentApp().Endpoints[s.endpointName].GetRestParams() != nil {
			if x := s.peekScope().(*sysl.Endpoint); x != nil && len(x.Stmt) == 0 {
				if len(x.Docstring) > 0 {
					space = " "
				}
				str := x.Docstring + space + text
				x.Docstring = str
				return
			}
		}

		stmt := s.lastStatement()
		if len(stmt.GetAction().Action) > 0 {
			space = " "
		}

		str := stmt.GetAction().Action + space + text
		stmt.GetAction().Action = str
		return
	}

	s.currentMultiLineAnno = append(s.currentMultiLineAnno, strings.TrimPrefix(ctx.TEXT().GetText(), " "))
}

// EnterTypes is called when production types is entered.
func (s *TreeShapeListener) EnterTypes(ctx *parser.TypesContext) {
	type1 := s.currentType()
	primType, constraints := primitiveFromNativeDataType(ctx.NativeDataTypes())
	if primType != nil {
		type1.Type = primType
		if constraints != nil {
			if type1.Constraint == nil {
				type1.Constraint = []*sysl.Type_Constraint{constraints}
			} else {
				type1.Constraint = append(type1.Constraint, constraints)
			}
		}
	}
}

// exitSetOrSequence_type is common between ExitSet_type and ExitSequence_type.
func (s *TreeShapeListener) exitSetOrSequence_type(
	sizeSpec parser.ISize_specContext,
	arraySpec parser.IArray_sizeContext,
) (type1, newType1 *sysl.Type) {
	type1 = s.currentType()
	newType1 = &sysl.Type{
		SourceContext:  type1.SourceContext, //nolint:staticcheck
		SourceContexts: type1.SourceContexts,
		Opt:            type1.Opt,
		Attrs:          type1.Attrs,
	}
	s.setCurrentType(newType1)

	type1.Opt = false

	if type1.Attrs != nil {
		type1.Attrs = nil
		s.popScope()
		s.pushScope(newType1)
	}

	if sizeSpec != nil {
		if type1.GetPrimitive() != sysl.Type_NO_Primitive {
			spec := sizeSpec.(*parser.Size_specContext)
			type1.Constraint = makeTypeConstraint(type1, spec)
		}
	} else if arraySpec != nil {
		if type1.GetPrimitive() != sysl.Type_NO_Primitive {
			spec := arraySpec.(*parser.Array_sizeContext)
			type1.Constraint = makeArrayConstraint(type1, spec)
		}
	}
	return type1, newType1
}

// ExitSet_type is called when production set_type is exited.
func (s *TreeShapeListener) ExitSet_type(ctx *parser.Set_typeContext) {
	type1, newType1 := s.exitSetOrSequence_type(ctx.Size_spec(), ctx.Array_size())
	newType1.Type = &sysl.Type_Set{Set: type1}
}

// EnterSequence_type is called when production set_type is entered.
func (s *TreeShapeListener) EnterSequence_type(*parser.Sequence_typeContext) {}

// ExitSequence_type is called when production set_type is exited.
func (s *TreeShapeListener) ExitSequence_type(ctx *parser.Sequence_typeContext) {
	type1, newType1 := s.exitSetOrSequence_type(ctx.Size_spec(), ctx.Array_size())
	newType1.Type = &sysl.Type_Sequence{Sequence: type1}
}

// EnterUser_defined_type is called when production user_defined_type is entered.
func (s *TreeShapeListener) EnterUser_defined_type(ctx *parser.User_defined_typeContext) {
	if len(s.fieldname) == 0 {
		return
	}
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	context_app_part := s.currentApp().Name.Part
	context_path := s.currentTypePath.Parts()
	ref_path := []string{MustUnescape(ctx.GetText())}

	type1.Type = &sysl.Type_TypeRef{
		TypeRef: &sysl.ScopedRef{
			Context: &sysl.Scope{
				Appname: &sysl.AppName{
					Part: context_app_part,
				},
				Path: context_path,
			},
			Ref: &sysl.Scope{
				Path: ref_path,
			},
		},
	}
}

func fromQString(str string) string {
	l := len(str)
	if l > 0 && strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
		return strings.Trim(str, "'")
	}
	if l > 0 && strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		var val string
		if json.Unmarshal([]byte(str), &val) == nil {
			return val
		}
	}
	return strings.Trim(str, `"`)
}

// EnterAnnotation_value is called when production annotation_value is entered.
func (s *TreeShapeListener) EnterAnnotation_value(ctx *parser.Annotation_valueContext) {
	sc := s.getSrcCtx(ctx.BaseParserRuleContext)
	attr := &sysl.Attribute{
		SourceContext:  sc,
		SourceContexts: []*sysl.SourceContext{sc},
	}
	switch {
	case ctx.Multi_line_docstring() != nil:
		attr.Attribute = &sysl.Attribute_S{}
	case ctx.QSTRING() != nil:
		attr.Attribute = &sysl.Attribute_S{
			S: fromQString(ctx.QSTRING().GetText()),
		}
	case ctx.Arrays() != nil:
		attr.Attribute = &sysl.Attribute_A{
			A: s.makeArraysAttribute(ctx.Arrays().(*parser.ArraysContext)).GetA(),
		}
	default:
		panic("unexpected annotation value")
	}
	addAttrWithPrecedence(s.peekAttrs(), s.annotation_name, attr)
}

// ExitAnnotation_value is called when production annotation_value is exited.
func (s *TreeShapeListener) ExitAnnotation_value(ctx *parser.Annotation_valueContext) {
	sc := s.getSrcCtx(ctx.BaseParserRuleContext)
	if ctx.Multi_line_docstring() != nil {
		addAttrWithPrecedence(s.peekAttrs(), s.annotation_name, &sysl.Attribute{
			Attribute: &sysl.Attribute_S{
				S: strings.TrimLeft(strings.Join(s.currentMultiLineAnno, "\n"), " ") + "\n",
			},
			SourceContext:  sc,
			SourceContexts: []*sysl.SourceContext{sc},
		})
		s.currentMultiLineAnno = []string{}
	}
}

// EnterAnnotation is called when production annotation is entered.
func (s *TreeShapeListener) EnterAnnotation(ctx *parser.AnnotationContext) {
	attr_name := ctx.VAR_NAME().GetText()
	attrs := s.peekAttrs()

	// if v.Attribute is not nil, there are already values defined for the
	// annotation which should not be replaced.
	if v, exists := attrs[attr_name]; !exists || v.Attribute == nil {
		attrs[attr_name] = &sysl.Attribute{}
	}
	s.annotation_name = attr_name
}

// ExitAnnotation is called when production annotation is exited.
func (s *TreeShapeListener) ExitAnnotation(*parser.AnnotationContext) {
	s.annotation_name = ""
}

// EnterField_type is called when production field_type is entered.
func (s *TreeShapeListener) EnterField_type(ctx *parser.Field_typeContext) {
	s.app_name.Reset()
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = mergeAttrsWithPrecendence(type1.Attrs, s.makeAttributeArray(attribs))
	}

	if ctx.QN() != nil {
		type1.Opt = true
	}

	if ctx.Annotations() != nil || ctx.Attribs_or_modifiers() != nil {
		if type1.Attrs == nil {
			type1.Attrs = map[string]*sysl.Attribute{}
		}
		s.pushScope(type1)
	}
}

func makeScope(app_name []string, ctx *parser.ReferenceContext) *sysl.Scope {
	scope := &sysl.Scope{}
	var dotCount int
	if ctx.DOT(0) != nil {
		dotCount = len(ctx.AllDOT())
	} else {
		dotCount = len(ctx.AllE_DOT())
	}
	appComponentCount := len(app_name) - dotCount

	if appComponentCount > 0 {
		scope.Appname = &sysl.AppName{
			Part: MustUnescapeStrings(app_name[:appComponentCount]),
		}
	}
	scope.Path = MustUnescapeStrings(app_name[appComponentCount:])
	return scope
}

// ExitField_type is called when production field_type is exited.
func (s *TreeShapeListener) ExitField_type(ctx *parser.Field_typeContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
	primitive_type := type1.GetPrimitive()
	if primitive_type != sysl.Type_NO_Primitive {
		size_spec, has_size_spec := ctx.Size_spec().(*parser.Size_specContext)
		array_spec, has_array_spec := ctx.Array_size().(*parser.Array_sizeContext)
		if has_size_spec {
			type1.Constraint = makeTypeConstraint(type1, size_spec)
		} else if has_array_spec {
			type1.Constraint = makeArrayConstraint(type1, array_spec)
		}
	}
	if ctx.Annotations() != nil || ctx.Attribs_or_modifiers() != nil {
		s.popScope()
	}
}

func makeTypeConstraint(t *sysl.Type, size_spec *parser.Size_specContext) []*sysl.Type_Constraint {
	c := []*sysl.Type_Constraint{}
	var err error
	var l int64

	bitWidth := getBitWidth(t)
	switch t.GetPrimitive() {
	case sysl.Type_DATE, sysl.Type_DATETIME, sysl.Type_INT, sysl.Type_STRING, sysl.Type_BYTES:
		val1 := size_spec.DIGITS(0).GetText()
		if l, err = strconv.ParseInt(val1, 10, 0); err == nil {
			if bitWidth > 0 {
				c = append(c, &sysl.Type_Constraint{
					BitWidth: bitWidth,
					Length:   &sysl.Type_Constraint_Length{Max: l},
				})
			} else {
				c = append(c, &sysl.Type_Constraint{
					Length: &sysl.Type_Constraint_Length{Max: l},
				})
			}
		}
	case sysl.Type_DECIMAL:
		val1 := size_spec.DIGITS(0).GetText()
		if l, err = strconv.ParseInt(val1, 10, 0); err == nil {
			c = append(c, &sysl.Type_Constraint{
				Length: &sysl.Type_Constraint_Length{Max: l},
			})
			if size_spec.DIGITS(1) != nil {
				c[0].Precision = int32(l)
				val1 = size_spec.DIGITS(1).GetText()
				if l, err = strconv.ParseInt(val1, 10, 0); err == nil {
					c[0].Scale = int32(l)
				}
			}
		}
	default:
		panic("should not be here")
	}

	syslutil.PanicOnErrorf(err, "should not happen: unable to parse size_spec")
	return c
}

func makeArrayConstraint(t *sysl.Type, array_size *parser.Array_sizeContext) []*sysl.Type_Constraint {
	c := []*sysl.Type_Constraint{}
	var err error
	var l int64

	bitWidth := getBitWidth(t)
	switch t.GetPrimitive() {
	case sysl.Type_DATE, sysl.Type_DATETIME, sysl.Type_INT, sysl.Type_DECIMAL, sysl.Type_STRING, sysl.Type_BYTES:
		var ct *sysl.Type_Constraint
		if bitWidth > 0 {
			ct = &sysl.Type_Constraint{
				BitWidth: bitWidth,
				Length:   &sysl.Type_Constraint_Length{},
			}
		} else {
			ct = &sysl.Type_Constraint{
				Length: &sysl.Type_Constraint_Length{},
			}
		}
		if array_size.DIGITS(0) != nil {
			val := array_size.DIGITS(0).GetText()
			if l, err = strconv.ParseInt(val, 10, 0); err == nil && l != 0 {
				ct.Length.Min = l
			}
		}

		if array_size.DIGITS(1) != nil {
			val1 := array_size.DIGITS(1).GetText()
			if l, err = strconv.ParseInt(val1, 10, 0); err == nil && l != 0 {
				ct.Length.Max = l
			}
		}
		c = append(c, ct)
	default:
		panic("should not be here")
	}

	if err != nil {
		panic("should not happen: unable to parse array_size")
	}
	return c
}

func (s *TreeShapeListener) makeArraysAttribute(arrays *parser.ArraysContext) *sysl.Attribute {
	sc := s.getSrcCtx(arrays.BaseParserRuleContext)

	// return empty array
	if arrays.AllArrays_item() == nil || len(arrays.AllArrays_item()) == 0 {
		return &sysl.Attribute{
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{Elt: []*sysl.Attribute{}},
			},
			SourceContext:  sc,
			SourceContexts: []*sysl.SourceContext{sc},
		}
	}

	elements := make([]*sysl.Attribute, 0, len(arrays.AllArrays_item()))
	for _, item := range arrays.AllArrays_item() {
		itemCtx := item.GetRuleContext().(*parser.Arrays_itemContext)
		var itemSrc *sysl.SourceContext
		var element *sysl.Attribute
		switch {
		case itemCtx.Quoted_string() != nil:
			i := itemCtx.Quoted_string().(*parser.Quoted_stringContext)
			itemSrc = s.getSrcCtx(i.BaseParserRuleContext)
			element = &sysl.Attribute{
				Attribute: &sysl.Attribute_S{
					S: fromQString(i.QSTRING().GetText()),
				},
			}
		case itemCtx.Arrays() != nil:
			i := itemCtx.Arrays().(*parser.ArraysContext)
			itemSrc = s.getSrcCtx(i.BaseParserRuleContext)
			element = s.makeArraysAttribute(i)
		default:
			panic(fmt.Errorf("unhandled annotation arrays element"))
		}

		element.SourceContext = itemSrc //nolint:staticcheck
		element.SourceContexts = []*sysl.SourceContext{itemSrc}

		elements = append(elements, element)
	}
	return &sysl.Attribute{
		Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{
				Elt: elements,
			},
		},
		SourceContext:  sc,
		SourceContexts: []*sysl.SourceContext{sc},
	}
}

func addAttrWithPrecedence(
	attrs map[string]*sysl.Attribute,
	key string, attr *sysl.Attribute,
) map[string]*sysl.Attribute {
	if attrs == nil {
		attrs = make(map[string]*sysl.Attribute)
	}
	if patterns, hasPatterns := attrs[patternsKey]; hasPatterns && key == patternsKey {
		currPatterns := patterns.Attribute.(*sysl.Attribute_A)
		newPatterns := attr.Attribute.(*sysl.Attribute_A)
		currPatterns.A.Elt = append(currPatterns.A.GetElt(), newPatterns.A.GetElt()...)
		return attrs
	}
	// This ensures that the first non-empty annotation has the highest precedence and does not get replaced
	if v, exists := attrs[key]; exists && v.Attribute != nil {
		switch x := v.Attribute.(type) {
		case *sysl.Attribute_S:
			if x.S != "" {
				v.SourceContexts = append(v.SourceContexts, attr.SourceContexts...)
				return attrs
			}
		case *sysl.Attribute_A:
			if len(x.A.GetElt()) > 0 {
				v.SourceContexts = append(v.SourceContexts, attr.SourceContexts...)
				return attrs
			}
		}
	}
	attrs[key] = attr
	return attrs
}

func mergeAttrsWithPrecendence(currentAttrs, newAttrs map[string]*sysl.Attribute) map[string]*sysl.Attribute {
	if currentAttrs == nil {
		currentAttrs = make(map[string]*sysl.Attribute)
	}
	for k, v := range newAttrs {
		currentAttrs = addAttrWithPrecedence(currentAttrs, k, v)
	}
	return currentAttrs
}

func (s *TreeShapeListener) makeAttributeArray(attribs *parser.Attribs_or_modifiersContext) map[string]*sysl.Attribute {
	patterns := []*sysl.Attribute{}
	attributes := map[string]*sysl.Attribute{}

	for _, e := range attribs.AllEntry() {
		entry := e.(*parser.EntryContext)
		if nvp, ok := entry.Nvp().(*parser.NvpContext); ok {
			switch {
			case nvp.Quoted_string() != nil:
				qs := nvp.Quoted_string().(*parser.Quoted_stringContext)
				sc := s.getSrcCtx(qs.BaseParserRuleContext)
				attributes[nvp.Name().GetText()] = &sysl.Attribute{
					Attribute: &sysl.Attribute_S{
						S: fromQString(qs.QSTRING().GetText()),
					},
					SourceContext:  sc,
					SourceContexts: []*sysl.SourceContext{sc},
				}
			case nvp.Arrays() != nil:
				attributes[nvp.Name().GetText()] = s.makeArraysAttribute(nvp.Arrays().(*parser.ArraysContext))
			}
		} else if mod, ok := entry.Modifier().(*parser.ModifierContext); ok {
			sc := s.getSrcCtx(mod.BaseParserRuleContext)
			patterns = append(patterns, &sysl.Attribute{
				Attribute: &sysl.Attribute_S{
					S: mod.GetText()[1:],
				},
				SourceContext:  sc,
				SourceContexts: []*sysl.SourceContext{sc},
			})
		}
	}
	if len(patterns) > 0 {
		attributes[patternsKey] = &sysl.Attribute{
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: patterns,
				},
			},
		}
	}
	return attributes
}

// EnterInplace_tuple is called when production inplace_tuple is entered.
func (s *TreeShapeListener) EnterInplace_tuple(*parser.Inplace_tupleContext) {
	s.currentTypePath.Push(MustUnescape(s.fieldname[len(s.fieldname)-1]))
	s.typemap = map[string]*sysl.Type{}
	s.currentApp().Types[s.currentTypePath.Get()] = &sysl.Type{
		Type: &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{
				AttrDefs: s.typemap,
			},
		},
	}
}

// ExitInplace_tuple is called when production inplace_tuple is exited.
func (s *TreeShapeListener) ExitInplace_tuple(*parser.Inplace_tupleContext) {
	s.currentTypePath.Pop()
	s.typemap = s.currentApp().Types[s.currentTypePath.Get()].GetTuple().GetAttrDefs()
}

// EnterField is called when production field is entered.
func (s *TreeShapeListener) EnterField(ctx *parser.FieldContext) {
	fieldName := MustUnescape(ctx.Name_str().GetText())
	s.fieldname = append(s.fieldname, fieldName)
	type1, has := s.typemap[fieldName]
	if has {
		logrus.Debugf("%s) %s.%s defined multiple times, merging field definitions",
			s.sc.filename, s.currentTypePath.Get(), fieldName)
	} else {
		type1 = &sysl.Type{}
		type1.Type = &sysl.Type_NoType_{
			NoType: &sysl.Type_NoType{},
		}
	}

	type1.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)             //nolint:staticcheck
	type1.SourceContexts = append(type1.SourceContexts, type1.SourceContext) //nolint:staticcheck
	if ctx.QSTRING() != nil {
		type1.Docstring = fromQString(ctx.QSTRING().GetText())
	}

	if ctx.Inplace_tuple() != nil {
		type1.Type = &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Ref: &sysl.Scope{
					Path: []string{fieldName},
				},
			},
		}
	}
	s.typemap[fieldName] = type1
	s.app_name.Reset()
}

// ExitField is called when production field is exited.
func (s *TreeShapeListener) ExitField(ctx *parser.FieldContext) {
	if ctx.Array_size() != nil {
		name := s.fieldname[len(s.fieldname)-1]
		itemType := s.typemap[name]
		if ctx.Inplace_tuple() != nil {
			name = ctx.Name_str().GetText()
			itemType = s.typemap[name]
		}
		if itemType == nil {
			itemType = &sysl.Type{}
		}
		type1 := &sysl.Type{
			Type: &sysl.Type_List_{
				List: &sysl.Type_List{
					Type: itemType,
				},
			},
		}
		s.typemap[name] = type1
	}
}

// EnterTable_stmts is called when production table_stmts is entered.
func (s *TreeShapeListener) EnterTable_stmts(ctx *parser.Table_stmtsContext) {
	type1 := s.currentApp().Types[s.currentTypePath.Get()]
	if ctx.Annotation(0) != nil {
		if type1.Attrs == nil {
			type1.Attrs = map[string]*sysl.Attribute{}
		}
	}
	if ctx.Field(0) == nil {
		if type1.Attrs == nil {
			switch t := type1.Type.(type) {
			case *sysl.Type_Relation_:
				if t.Relation.AttrDefs == nil {
					type1.Type = nil
				}
			case *sysl.Type_Tuple_:
				if t.Tuple.AttrDefs == nil {
					type1.Type = nil
				}
			}
		}
	}
}

// EnterTable_def is called when production table_def is entered.
func (s *TreeShapeListener) EnterTable_def(ctx *parser.Table_defContext) {
	type1 := s.currentApp().Types[s.currentTypePath.Get()]
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		if attrs := s.makeAttributeArray(attribs); type1.Attrs == nil {
			type1.Attrs = attrs
		} else {
			for k, v := range attrs {
				if k == patternsKey {
					if patterns, hasPatterns := type1.GetAttrs()[patternsKey]; hasPatterns {
						patterns.GetA().Elt = append(patterns.GetA().Elt, v.GetA().Elt...)
						continue
					}
				}
				type1.Attrs[k] = v
			}
		}
	}
	if ctx.WHATEVER() != nil {
		type1.Type = nil
	}
}

// EnterTable is called when production table is entered.
func (s *TreeShapeListener) EnterTable(ctx *parser.TableContext) {
	s.currentTypePath.Push(MustUnescape(ctx.Name_str().GetText()))
	s.typemap = map[string]*sysl.Type{}

	types := s.currentApp().Types
	if existing, ok := types[s.currentTypePath.Get()]; ok {
		switch e := existing.Type.(type) {
		case *sysl.Type_Relation_:
			s.typemap = e.Relation.AttrDefs
		case *sysl.Type_Tuple_:
			s.typemap = e.Tuple.AttrDefs
		default:
		}
	}
	if ctx.TABLE() != nil {
		if _, exists := types[s.currentTypePath.Get()]; !exists {
			types[s.currentTypePath.Get()] = &sysl.Type{
				Type: &sysl.Type_Relation_{
					Relation: &sysl.Type_Relation{
						AttrDefs: s.typemap,
					},
				},
			}
		}
	}
	if ctx.TYPE() != nil {
		if _, exists := types[s.currentTypePath.Get()]; !exists {
			types[s.currentTypePath.Get()] = &sysl.Type{
				Type: &sysl.Type_Tuple_{
					Tuple: &sysl.Type_Tuple{
						AttrDefs: s.typemap,
					},
				},
			}
		}
	}
	type1 := types[s.currentTypePath.Get()]
	s.pushScope(type1)
	type1.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)             //nolint:staticcheck
	type1.SourceContexts = append(type1.SourceContexts, type1.SourceContext) //nolint:staticcheck
}

func attributesForType(collection *sysl.Type) map[string]*sysl.Type {
	var attrs map[string]*sysl.Type

	switch x := collection.Type.(type) {
	case *sysl.Type_Relation_:
		attrs = x.Relation.AttrDefs
	case *sysl.Type_Tuple_:
		attrs = x.Tuple.AttrDefs
	}
	return attrs
}

// ExitTable is called when production table is exited.
func (s *TreeShapeListener) ExitTable(ctx *parser.TableContext) {
	t := s.currentApp().Types[s.currentTypePath.Get()]
	// wire up primary key
	if rel := t.GetRelation(); rel != nil {
		pks := []string{}
		for _, name := range s.fieldname {
			f := rel.GetAttrDefs()[name]
			if patterns, has := f.GetAttrs()[patternsKey]; has {
				for _, a := range patterns.GetA().Elt {
					if a.GetS() == "pk" {
						pks = append(pks, name)
					}
				}
			}
		}
		if len(pks) > 0 {
			rel.PrimaryKey = &sysl.Type_Relation_Key{
				AttrName: pks,
			}
		}
	}
	s.popScope()

	t.SourceContext.End = s.lastEnd //nolint:staticcheck
	t.SourceContexts[len(t.SourceContexts)-1].End = s.lastEnd

	s.currentTypePath.Pop()
	s.fieldname = []string{}
	s.typemap = nil
}

func (s *TreeShapeListener) applyAnnotations(
	annotations []parser.IAnnotationContext,
) {
	// Match legacy behavior
	// Copy the annotations from the parent (tuple or relation) to each child
	collection := s.currentApp().Types[s.currentTypePath.Get()]

	if collection == nil {
		return
	}
	switch collection.Type.(type) {
	case *sysl.Type_Tuple_, *sysl.Type_Relation_:
	default:
		// skip if type is not tuple or relation
		return
	}

	for _, annotation := range annotations {
		varname := annotation.(*parser.AnnotationContext).VAR_NAME().GetText()
		attr := collection.Attrs[varname]
		for _, name := range s.fieldname {
			f := attributesForType(collection)[name]

			if f.Attrs == nil {
				f.Attrs = map[string]*sysl.Attribute{}
			}
			// if varname is already in lower scope, it shouldn't replace the value
			if _, exists := f.Attrs[varname]; !exists {
				f.Attrs[varname] = attr
			}
		}
	}
}

// EnterUnion is called when production union is entered.
func (s *TreeShapeListener) EnterUnion(ctx *parser.UnionContext) {
	s.currentTypePath.Push(MustUnescape(ctx.Name_str().GetText()))
	s.typemap = map[string]*sysl.Type{}

	types := s.currentApp().Types
	if types == nil {
		types = map[string]*sysl.Type{}
		s.currentApp().Types = types
	}

	if types[s.currentTypePath.Get()].GetOneOf().GetType() != nil {
		logrus.Warnf("oneof is not implemented yet")
	}

	types[s.currentTypePath.Get()] = &sysl.Type{
		Type: &sysl.Type_OneOf_{
			OneOf: &sysl.Type_OneOf{},
		},
	}

	type1 := types[s.currentTypePath.Get()]
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = s.makeAttributeArray(attribs)
	}
	if ctx.Annotation(0) != nil {
		if type1.Attrs == nil {
			type1.Attrs = map[string]*sysl.Attribute{}
		}
	}
	s.pushScope(type1)
	type1.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)             //nolint:staticcheck
	type1.SourceContexts = append(type1.SourceContexts, type1.SourceContext) //nolint:staticcheck
}

// EnterUnion_type is called when production union_type is entered.
func (s *TreeShapeListener) EnterUnion_type(ctx *parser.Union_typeContext) {
	fieldName := ctx.GetText()
	s.fieldname = append(s.fieldname, fieldName)
	dataType, has := s.typemap[fieldName]
	if has {
		logrus.Debugf("%s) %s.(%s) defined multiple times, redefinition ignored",
			s.sc.filename, s.currentTypePath.Get(), fieldName)
	} else {
		dataType = &sysl.Type{}
		dataType.Type = &sysl.Type_NoType_{
			NoType: &sysl.Type_NoType{},
		}
	}

	dataType.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)                   //nolint:staticcheck
	dataType.SourceContexts = append(dataType.SourceContexts, dataType.SourceContext) //nolint:staticcheck
	s.typemap[fieldName] = dataType
}

// ExitUnion is called when production union is exited.
func (s *TreeShapeListener) ExitUnion(ctx *parser.UnionContext) {
	s.applyAnnotations(ctx.AllAnnotation())
	s.popScope()

	oneof := s.currentApp().Types[s.currentTypePath.Get()].GetOneOf()
	for _, i := range ctx.AllUnion_type() {
		oneof.Type = append(oneof.Type, s.typemap[i.GetText()])
	}

	s.currentTypePath.Pop()
	s.fieldname = []string{}
	s.typemap = nil
}

// EnterApp_name is called when production app_name is entered.
func (s *TreeShapeListener) EnterApp_name(*parser.App_nameContext) {
	s.app_name.Reset()
}

// EnterName_with_attribs is called when production name_with_attribs is entered.
func (s *TreeShapeListener) EnterName_with_attribs(ctx *parser.Name_with_attribsContext) {
	appName := MustUnescape(syslutil.CleanAppName(ctx.App_name().GetText()))
	if _, has := s.module.Apps[appName]; !has {
		s.module.Apps[appName] = &sysl.Application{
			Name: &sysl.AppName{},
		}
	}
	s.app = s.module.Apps[appName]

	if ctx.QSTRING() != nil {
		s.currentApp().LongName = fromQString(ctx.QSTRING().GetText())
	}

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		attrs := s.makeAttributeArray(attribs)
		if s.currentApp().Attrs == nil {
			s.currentApp().Attrs = attrs
		} else {
			mergeAttrs(attrs, s.currentApp().Attrs)
		}
	}

	sc := s.getSrcCtx(ctx.BaseParserRuleContext)
	s.currentApp().SourceContext = sc //nolint:staticcheck
	s.currentApp().SourceContexts = append(s.currentApp().SourceContexts, sc)
}

// ExitName_with_attribs is called when production name_with_attribs is exited.
func (s *TreeShapeListener) ExitName_with_attribs(*parser.Name_with_attribsContext) {
	s.currentApp().Name.Part = MustUnescapeStrings(s.app_name.Parts())
}

// EnterModel_name is called when production model_name is entered.
func (s *TreeShapeListener) EnterModel_name(ctx *parser.Model_nameContext) {
	if s.currentApp().Wrapped.Name != nil {
		panic("not implemented yet?")
	}

	s.currentApp().Wrapped.Name = &sysl.AppName{
		Part: []string{ctx.Name().GetText()},
	}
}

// EnterTable_refs is called when production table_refs is entered.
func (s *TreeShapeListener) EnterTable_refs(ctx *parser.Table_refsContext) {
	s.currentApp().Wrapped.Types[ctx.Name().GetText()] = &sysl.Type{}
}

// EnterQuery_var is called when production query_var is entered.
func (s *TreeShapeListener) EnterQuery_var(ctx *parser.Query_varContext) {
	var_name := ctx.Name().GetText()
	var type1 *sysl.Type
	var ref_path []string

	switch {
	case ctx.Var_in_curly() != nil:
		cleanString := strings.Replace(ctx.Var_in_curly().GetText(), "{", "", 1)
		cleanString = strings.Replace(cleanString, "}", "", 1)
		ref_path = append(ref_path, cleanString)
		type1 = &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: s.currentApp().Name.Part,
						},
					},
					Ref: &sysl.Scope{
						Path: ref_path,
					},
				},
			},
		}
	case ctx.NativeDataTypes() != nil:
		primType, constraints := primitiveFromNativeDataType(ctx.NativeDataTypes())
		type1 = &sysl.Type{
			Type: primType,
		}
		if constraints != nil {
			type1.Constraint = []*sysl.Type_Constraint{constraints}
		}
	case ctx.Name_str() != nil:
		type1 = &sysl.Type{}
	}

	rest_param := &sysl.Endpoint_RestParams_QueryParam{
		Name: MustUnescape(var_name),
		Type: type1,
	}

	if ctx.QN() != nil {
		rest_param.Type.Opt = true
	}

	//nolint:staticcheck
	rest_param.Type.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)
	//nolint:staticcheck
	rest_param.Type.SourceContexts = append(rest_param.Type.SourceContexts, rest_param.Type.SourceContext)
	s.method_urlparams = append(s.method_urlparams, rest_param)
}

// EnterHttp_path_var_with_type is called when production http_path_var_with_type is entered.
func (s *TreeShapeListener) EnterHttp_path_var_with_type(ctx *parser.Http_path_var_with_typeContext) {
	var_name := ctx.Http_path_part().GetText()
	var type1 *sysl.Type
	switch {
	case ctx.NativeDataTypes() != nil:
		primType, constraints := primitiveFromNativeDataType(ctx.NativeDataTypes())
		type1 = &sysl.Type{
			Type: primType,
		}
		if constraints != nil {
			type1.Constraint = []*sysl.Type_Constraint{constraints}
		}
	case ctx.Reference() != nil:
		s.fieldname = append(s.fieldname, var_name)
		type1 = &sysl.Type{}
		s.typemap[s.fieldname[len(s.fieldname)-1]] = type1
	default:
		ref_path := []string{ctx.Name_str().GetText()}

		type1 = &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: s.currentApp().Name.Part,
						},
					},
					Ref: &sysl.Scope{
						Path: ref_path,
					},
				},
			},
		}
	}
	rest_param := &sysl.Endpoint_RestParams_QueryParam{
		Name: MustUnescape(var_name),
		Type: type1,
	}
	//nolint:staticcheck
	rest_param.Type.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)
	//nolint:staticcheck
	rest_param.Type.SourceContexts = append(rest_param.Type.SourceContexts, rest_param.Type.SourceContext)

	s.rest_urlparams = append(s.rest_urlparams, rest_param)
	s.endpointName += ctx.CURLY_OPEN().GetText() + var_name + ctx.CURLY_CLOSE().GetText()
}

// ExitHttp_path_var_with_type is called when production http_path_var_with_type is exited.
func (s *TreeShapeListener) ExitHttp_path_var_with_type(ctx *parser.Http_path_var_with_typeContext) {
	if ctx.Reference() != nil {
		type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
		type1.GetTypeRef().Context.Path = nil
		type1.GetTypeRef().Ref.Path = append(type1.GetTypeRef().Ref.Appname.Part, type1.GetTypeRef().Ref.Path...)
		type1.GetTypeRef().Ref.Appname = nil
	}
}

// EnterHttp_path_static is called when production http_path_static is entered.
func (s *TreeShapeListener) EnterHttp_path_static(ctx *parser.Http_path_staticContext) {
	s.endpointName += ctx.GetText()
}

// EnterHttp_path_suffix is called when production http_path_suffix is entered.
func (s *TreeShapeListener) EnterHttp_path_suffix(ctx *parser.Http_path_suffixContext) {
	s.endpointName += ctx.FORWARD_SLASH().GetText()
}

// EnterHttp_path is called when production http_path is entered.
func (s *TreeShapeListener) EnterHttp_path(ctx *parser.Http_pathContext) {
	s.endpointName = ""
	if ctx.FORWARD_SLASH() != nil {
		s.endpointName = ctx.GetText()
	}
}

// ExitHttp_path is called when production http_path is exited.
func (s *TreeShapeListener) ExitHttp_path(*parser.Http_pathContext) {
	// s.endpointName is built along as we enter http_path/http_path_suffix/http_path_var_with_type
	// commit this value to urlPrefixes
	s.urlPrefixes.Push(MustUnescape(s.endpointName))
}

// EnterRet_stmt is called when production ret_stmt is entered.
func (s *TreeShapeListener) EnterRet_stmt(ctx *parser.Ret_stmtContext) {
	payload := MustUnescape(strings.Trim(ctx.TEXT().GetText(), " "))
	if s.lintMode {
		s.lintRetStmt(payload, ctx)
	}
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Ret{
			Ret: &sysl.Return{
				Payload: payload,
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	})
}

// EnterTarget is called when production target is entered.
func (s *TreeShapeListener) EnterTarget(*parser.TargetContext) {
	s.app_name.Reset()
}

// ExitTarget is called when production target is exited.
func (s *TreeShapeListener) ExitTarget(*parser.TargetContext) {
	s.lastStatement().GetCall().Target.Part = MustUnescapeStrings(s.app_name.Parts())
	s.app_name.Reset()
}

// EnterCall_arg is called when production call_arg is entered.
func (s *TreeShapeListener) EnterCall_arg(ctx *parser.Call_argContext) {
	arg := &sysl.Call_Arg{
		Name: ctx.GetText(),
	}
	s.lastStatement().GetCall().Arg = append(s.lastStatement().GetCall().Arg, arg)
}

// EnterCall_stmt is called when production call_stmt is entered.
func (s *TreeShapeListener) EnterCall_stmt(ctx *parser.Call_stmtContext) {
	var name string
	//FIXME: app target isn't recorded at all
	appName := &sysl.AppName{}
	if ctx.DOT_ARROW() != nil {
		appName.Part = s.currentApp().Name.Part
		name = s.getFullAppName()
	} else {
		//TODO: handle target with sub packages
		name = MustUnescape(ctx.Target().GetText())
	}
	endpoint := MustUnescape(ctx.Target_endpoint().GetText())
	methodRe := regexp.MustCompile(`^[ \t]*(GET|POST|DELETE|PUT|PATCH)[ \t]*`)
	location := s.createLocation(ctx.GetStart())
	if methodRe.MatchString(endpoint) {
		method := methodRe.FindString(endpoint)
		urlPath := strings.TrimPrefix(endpoint, method)
		s.recordCall(name, strings.TrimSpace(urlPath), strings.TrimSpace(method), location)
	} else {
		s.recordCall(name, strings.TrimSpace(endpoint), "", location)
	}
	sc := s.getSrcCtxFor(ctx.BaseParserRuleContext.GetStart(), lastToken(ctx))
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: endpoint,
			},
		},
		SourceContext:  sc,
		SourceContexts: []*sysl.SourceContext{sc},
	})

	if ctx.Call_args() != nil {
		s.lastStatement().GetCall().Arg = []*sysl.Call_Arg{}
	}
}

// EnterIf_stmt is called when production if_stmt is entered.
func (s *TreeShapeListener) EnterIf_stmt(ctx *parser.If_stmtContext) {
	if_stmt := &sysl.Statement{
		Stmt: &sysl.Statement_Cond{
			Cond: &sysl.Cond{
				Test: ctx.PREDICATE_VALUE().GetText(),
				Stmt: []*sysl.Statement{},
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	}
	s.addToCurrentScope(if_stmt)
	s.pushScope(if_stmt.GetCond())
}

// ExitIf_stmt is called when production if_stmt is exited.
func (s *TreeShapeListener) ExitIf_stmt(*parser.If_stmtContext) {
	s.popScope()
}

// EnterElse_stmt is called when production else_stmt is entered.
func (s *TreeShapeListener) EnterElse_stmt(ctx *parser.Else_stmtContext) {
	else_cond := ctx.ELSE().GetText()
	if ctx.PREDICATE_VALUE() != nil {
		else_cond += ctx.PREDICATE_VALUE().GetText()
	}
	else_cond = strings.TrimSpace(else_cond)
	else_stmt := &sysl.Statement{
		Stmt: &sysl.Statement_Group{
			Group: &sysl.Group{
				Title: else_cond,
				Stmt:  []*sysl.Statement{},
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	}
	s.addToCurrentScope(else_stmt)
	s.pushScope(else_stmt.GetGroup())
}

// ExitElse_stmt is called when production else_stmt is exited.
func (s *TreeShapeListener) ExitElse_stmt(*parser.Else_stmtContext) {
	s.popScope()
}

// EnterFor_stmt is called when production for_stmt is entered.
func (s *TreeShapeListener) EnterFor_stmt(ctx *parser.For_stmtContext) {
	stmt := &sysl.Statement{
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	}
	s.addToCurrentScope(stmt)

	switch {
	case ctx.FOR() != nil || ctx.LOOP() != nil:
		var text string
		if ctx.FOR() != nil {
			text = ctx.FOR().GetText()
		} else {
			text = ctx.LOOP().GetText()
		}
		text = strings.TrimSpace(text + ctx.PREDICATE_VALUE().GetText())
		stmt.Stmt = &sysl.Statement_Group{
			Group: &sysl.Group{
				Title: text,
				Stmt:  []*sysl.Statement{},
			},
		}
		s.pushScope(stmt.GetGroup())
	case ctx.UNTIL() != nil || ctx.WHILE() != nil:
		mode := sysl.Loop_UNTIL
		if ctx.WHILE() != nil {
			mode = sysl.Loop_WHILE
		}
		stmt.Stmt = &sysl.Statement_Loop{
			Loop: &sysl.Loop{
				Mode:      mode,
				Criterion: ctx.PREDICATE_VALUE().GetText(),
				Stmt:      []*sysl.Statement{},
			},
		}
		s.pushScope(stmt.GetLoop())
	case ctx.FOR_EACH() != nil:
		stmt.Stmt = &sysl.Statement_Foreach{
			Foreach: &sysl.Foreach{
				Collection: ctx.PREDICATE_VALUE().GetText(),
				Stmt:       []*sysl.Statement{},
			},
		}
		s.pushScope(stmt.GetForeach())
	case ctx.ALT() != nil:
		text := ctx.ALT().GetText() + ctx.PREDICATE_VALUE().GetText()
		text = strings.TrimSpace(text)
		stmt.Stmt = &sysl.Statement_Group{
			Group: &sysl.Group{
				Title: text,
				Stmt:  []*sysl.Statement{},
			},
		}
		s.pushScope(stmt.GetGroup())
	}
}

// ExitFor_stmt is called when production for_stmt is exited.
func (s *TreeShapeListener) ExitFor_stmt(*parser.For_stmtContext) {
	s.popScope()
}

// EnterGroup_stmt is called when production group_stmt is entered.
func (s *TreeShapeListener) EnterGroup_stmt(ctx *parser.Group_stmtContext) {
	stmt := &sysl.Statement{
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	}

	text := ctx.Name_str().GetText()
	text = strings.TrimSpace(text)
	stmt.Stmt = &sysl.Statement_Group{
		Group: &sysl.Group{
			Title: text,
			Stmt:  []*sysl.Statement{},
		},
	}
	s.addToCurrentScope(stmt)
	s.pushScope(stmt.GetGroup())
}

// ExitGroup_stmt is called when production group_stmt is exited.
func (s *TreeShapeListener) ExitGroup_stmt(*parser.Group_stmtContext) {
	s.popScope()
}

// EnterOne_of_cases is called when production one_of_cases is entered.
func (s *TreeShapeListener) EnterOne_of_cases(ctx *parser.One_of_casesContext) {
	alt := s.peekScope().(*sysl.Alt)
	choice := &sysl.Alt_Choice{
		Stmt: []*sysl.Statement{},
	}
	if ctx.One_of_case_label() != nil {
		choice.Cond = ctx.One_of_case_label().GetText()
	}
	alt.Choice = append(alt.Choice, choice)
	s.pushScope(choice)
}

// ExitOne_of_cases is called when production one_of_cases is exited.
func (s *TreeShapeListener) ExitOne_of_cases(*parser.One_of_casesContext) {
	s.popScope()
}

// EnterOne_of_stmt is called when production one_of_stmt is entered.
func (s *TreeShapeListener) EnterOne_of_stmt(ctx *parser.One_of_stmtContext) {
	alt := &sysl.Statement_Alt{
		Alt: &sysl.Alt{
			Choice: []*sysl.Alt_Choice{},
		},
	}
	s.addToCurrentScope(&sysl.Statement{
		Stmt:           alt,
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	})
	s.pushScope(alt.Alt)
}

// ExitOne_of_stmt is called when production one_of_stmt is exited.
func (s *TreeShapeListener) ExitOne_of_stmt(*parser.One_of_stmtContext) {
	s.popScope()
}

func withQuotesQString(str string) string {
	s := str[0]
	e := str[len(str)-1]
	syslutil.Assert(s == '\'' && e == '\'' || s == '"' && e == '"', "%q", str)
	return str
}

// EnterText_stmt is called when production text_stmt is entered.
func (s *TreeShapeListener) EnterText_stmt(ctx *parser.Text_stmtContext) {
	// Need to Coalesce multiple doc_string's into one
	// See enterdoc_string.
	if ctx.Doc_string() == nil {
		s.pendingDocString = false
		str := ctx.GetText()
		if ctx.QSTRING() != nil {
			str = withQuotesQString(str)
		}

		// Adjust the SourceContext end to account for the stmt text
		sc := s.getSrcCtx(ctx.BaseParserRuleContext)
		sc.End.Col = sc.Start.Col + int32(len(str))
		s.lastEnd = sc.End

		s.addToCurrentScope(&sysl.Statement{
			Stmt: &sysl.Statement_Action{
				Action: &sysl.Action{
					Action: str,
				},
			},
			SourceContext:  sc,
			SourceContexts: []*sysl.SourceContext{sc},
		})
	} else {
		s.pendingDocString = true

		if s.currentApp().Endpoints[s.endpointName].GetRestParams() != nil {
			if x := s.peekScope().(*sysl.Endpoint); x != nil && len(x.Stmt) == 0 {
				return
			}
		}
		// if laststatement is nil, add
		// if laststatement is not text statement add
		// if last statement is text statement does not start with  '|',  add
		x := s.lastStatement()
		add_stmt := x == nil || x.GetAction() == nil || !strings.HasPrefix(x.GetAction().Action, "|")

		if add_stmt {
			sc := s.getSrcCtx(ctx.BaseParserRuleContext)
			s.addToCurrentScope(&sysl.Statement{
				Stmt: &sysl.Statement_Action{
					Action: &sysl.Action{
						Action: "|",
					},
				},
				SourceContext:  sc,
				SourceContexts: []*sysl.SourceContext{sc},
			})
		}
	}
}

// EnterMixin is called when production mixin is entered.
func (s *TreeShapeListener) EnterMixin(*parser.MixinContext) {
	if s.currentApp().Mixin2 == nil {
		s.currentApp().Mixin2 = []*sysl.Application{}
	}
}

// ExitMixin is called when production mixin is exited.
func (s *TreeShapeListener) ExitMixin(*parser.MixinContext) {
	s.currentApp().Mixin2 = append(s.currentApp().Mixin2, &sysl.Application{
		Name: &sysl.AppName{
			Part: MustUnescapeStrings(s.app_name.Parts()),
		},
	})
}

// EnterParam is called when production param is entered.
func (s *TreeShapeListener) EnterParam(ctx *parser.ParamContext) {
	if ctx.Reference() != nil {
		s.fieldname = append(s.fieldname, ctx.Reference().GetText())
		type1 := &sysl.Type{}
		s.typemap[s.fieldname[len(s.fieldname)-1]] = type1
	}
	s.app_name.Reset()
}

// ExitParam is called when production param is exited.
func (s *TreeShapeListener) ExitParam(ctx *parser.ParamContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	if ctx.Reference() != nil {
		type1.Type = &sysl.Type_NoType_{
			NoType: &sysl.Type_NoType{},
		}
	}
}

// EnterParams is called when production params is entered.
func (s *TreeShapeListener) EnterParams(*parser.ParamsContext) {
	s.typemap = map[string]*sysl.Type{}
	s.fieldname = []string{}
}

// ExitParams is called when production params is exited.
func (s *TreeShapeListener) ExitParams(*parser.ParamsContext) {
	params := []*sysl.Param{}

	for _, fieldname := range s.fieldname {
		type1 := s.typemap[fieldname]
		switch t := type1.Type.(type) {
		case *sysl.Type_Set:
			t.Set.GetTypeRef().Context = nil
			t.Set.SourceContext = nil //nolint:staticcheck
			t.Set.SourceContexts = nil
		case *sysl.Type_TypeRef:
			t.TypeRef.Context = nil
			ref := t.TypeRef.GetRef()
			if ref.Appname != nil && ref.Appname.Part != nil {
				for i := range ref.Appname.Part {
					ref.Appname.Part[i] = strings.TrimSpace(ref.Appname.Part[i])
				}
			}
		case nil:
			type1.Type = &sysl.Type_NoType_{
				NoType: &sysl.Type_NoType{},
			}
		}
		type1.SourceContext = nil //nolint:staticcheck
		type1.SourceContexts = nil

		p := sysl.Param{
			Name: fieldname,
			Type: type1,
		}
		params = append(params, &p)
	}
	ep := s.currentApp().Endpoints[s.endpointName]
	if len(params) > 0 {
		if ep.Param == nil {
			ep.Param = params
		} else {
			ep.Param = append(ep.Param, params...)
		}
	}
	s.typemap = nil
	s.fieldname = []string{}
}

// ExitStatements is called when production statements is exited.
func (s *TreeShapeListener) ExitStatements(ctx *parser.StatementsContext) {
	if ctx.Attribs_or_modifiers() != nil {
		stmt := s.lastStatement()
		stmt.Attrs = s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}
}

func (s *TreeShapeListener) PushExpr(expr *sysl.Expr) {
	s.expr_stack = append(s.expr_stack, expr)
}

func (s *TreeShapeListener) popExpr() *sysl.Expr {
	l := len(s.expr_stack) - 1
	ret := s.expr_stack[l]
	s.expr_stack = s.expr_stack[:l]
	return ret
}

func (s *TreeShapeListener) TopExpr() *sysl.Expr {
	l := len(s.expr_stack) - 1
	ret := s.expr_stack[l]
	return ret
}

func (s *TreeShapeListener) pushScope(scope interface{}) {
	s.stmt_scope = append(s.stmt_scope, scope)
}

func (s *TreeShapeListener) popScope() {
	top := s.lastStatement()
	if top != nil {
		top.SourceContext.End = s.lastEnd //nolint:staticcheck
		top.SourceContexts[len(top.SourceContexts)-1].End = s.lastEnd
	}

	l := len(s.stmt_scope) - 1
	s.stmt_scope = s.stmt_scope[:l]
}

func (s *TreeShapeListener) peekScope() interface{} {
	l := len(s.stmt_scope) - 1
	return s.stmt_scope[l]
}

func (s *TreeShapeListener) lastStatement() *sysl.Statement {
	switch scope := s.peekScope().(type) {
	case *sysl.Application, *sysl.Type, *sysl.Alt, *sysl.View:
		return nil
	case *sysl.Endpoint:
		l := len(scope.Stmt) - 1
		if l < 0 {
			return nil
		}
		return scope.Stmt[l]
	case *sysl.Cond:
		l := len(scope.Stmt) - 1
		if l < 0 {
			return nil
		}
		return scope.Stmt[l]
	case *sysl.Alt_Choice:
		l := len(scope.Stmt) - 1
		if l < 0 {
			return nil
		}
		return scope.Stmt[l]
	case *sysl.Group:
		l := len(scope.Stmt) - 1
		if l < 0 {
			return nil
		}
		return scope.Stmt[l]
	case *sysl.Loop:
		l := len(scope.Stmt) - 1
		if l < 0 {
			return nil
		}
		return scope.Stmt[l]
	case *sysl.Foreach:
		l := len(scope.Stmt) - 1
		if l < 0 {
			return nil
		}
		return scope.Stmt[l]
	default:
		fmt.Printf("got unexpected %T\n", scope)
		panic("not implemented")
	}
}

func (s *TreeShapeListener) peekAttrs() map[string]*sysl.Attribute {
	switch t := s.peekScope().(type) {
	case *sysl.Application:
		return t.Attrs
	case *sysl.Type:
		return t.Attrs
	case *sysl.Endpoint:
		return t.Attrs
	case *sysl.View:
		return t.Attrs
	default:
		fmt.Printf("got unexpected %T\n", t)
		panic("not implemented")
	}
}

func (s *TreeShapeListener) addToCurrentScope(stmt *sysl.Statement) {
	top := len(s.stmt_scope) - 1
	switch scope := s.stmt_scope[top].(type) {
	case *sysl.Endpoint:
		scope.Stmt = append(scope.Stmt, stmt)
	case *sysl.Cond:
		scope.Stmt = append(scope.Stmt, stmt)
	case *sysl.Group:
		scope.Stmt = append(scope.Stmt, stmt)
	case *sysl.Alt_Choice:
		scope.Stmt = append(scope.Stmt, stmt)
	case *sysl.Loop:
		scope.Stmt = append(scope.Stmt, stmt)
	case *sysl.Foreach:
		scope.Stmt = append(scope.Stmt, stmt)
	default:
		fmt.Printf("got unexpected %T\n", scope)
		panic("not implemented")
	}
}

func mergeAttrs(src map[string]*sysl.Attribute, dst map[string]*sysl.Attribute) {
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

// EnterMethod_def is called when production method_def is entered.
func (s *TreeShapeListener) EnterMethod_def(ctx *parser.Method_defContext) {
	url := s.urlPrefixes.Get()
	method := strings.TrimSpace(ctx.HTTP_VERBS().GetText())
	s.endpointName = method + " " + url
	s.recordMethod(url, method, s.createLocation(ctx.GetStart()))
	s.method_urlparams = []*sysl.Endpoint_RestParams_QueryParam{}
	if _, has := s.currentApp().Endpoints[s.endpointName]; !has {
		s.currentApp().Endpoints[s.endpointName] = &sysl.Endpoint{
			Name: s.endpointName,
			RestParams: &sysl.Endpoint_RestParams{
				Method: sysl.Endpoint_RestParams_Method(sysl.Endpoint_RestParams_Method_value[method]),
				Path:   url,
			},
			Stmt: []*sysl.Statement{},
		}
	}
	restEndpoint := s.currentApp().Endpoints[s.endpointName]
	s.pushScope(restEndpoint)

	attrs := map[string]*sysl.Attribute{
		patternsKey: {
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: []*sysl.Attribute{
						{
							Attribute: &sysl.Attribute_S{
								S: "rest",
							},
						},
					},
				},
			},
		},
	}
	for _, parentAttrs := range s.rest_attrs {
		mergeAttrs(parentAttrs, attrs)
	}

	if ctx.Attribs_or_modifiers() != nil {
		mergeAttrs(s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext)), attrs)
	}

	if restEndpoint.Attrs == nil {
		restEndpoint.Attrs = attrs
	} else {
		mergeAttrs(attrs, restEndpoint.Attrs)
	}

	if len(s.rest_queryparams) > 0 {
		qparams := []*sysl.Endpoint_RestParams_QueryParam{}
		for i := range s.rest_queryparams {
			q := s.rest_queryparams[i]
			qcopy := &sysl.Endpoint_RestParams_QueryParam{
				Name: MustUnescape(q.Name),
				Type: &sysl.Type{
					Type:           q.Type.Type,
					SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
					SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
				},
			}
			qparams = append(qparams, qcopy)
		}
		restEndpoint.RestParams.QueryParam = qparams
	}

	if len(s.rest_urlparams) > 0 {
		qparams := []*sysl.Endpoint_RestParams_QueryParam{}
		for i := range s.rest_urlparams {
			qparams = append(qparams, s.rest_urlparams[i])
		}
		restEndpoint.RestParams.UrlParam = qparams
	}
	//nolint:staticcheck
	restEndpoint.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)
	//nolint:staticcheck
	restEndpoint.SourceContexts = append(restEndpoint.SourceContexts, restEndpoint.SourceContext)
}

// ExitMethod_def is called when production method_def is exited.
func (s *TreeShapeListener) ExitMethod_def(*parser.Method_defContext) {
	if len(s.method_urlparams) > 0 {
		qparams := s.currentApp().Endpoints[s.endpointName].RestParams.QueryParam
		if qparams == nil {
			qparams = []*sysl.Endpoint_RestParams_QueryParam{}
		}
		qparams = append(qparams, s.method_urlparams...)
		s.currentApp().Endpoints[s.endpointName].RestParams.QueryParam = qparams
	}
	if s.currentApp().Endpoints[s.endpointName].Param != nil {
		for i, p := range s.currentApp().Endpoints[s.endpointName].Param {
			if p.GetType().GetNoType() != nil {
				s.currentApp().Endpoints[s.endpointName].Param[i] = &sysl.Param{}
			}
		}
	}

	if len(s.currentApp().Endpoints[s.endpointName].Stmt) == 0 {
		s.currentApp().Endpoints[s.endpointName].Stmt = nil
	}

	s.popScope()
	s.endpointName = ""
	s.http_path_query_param = ""
}

// EnterSimple_endpoint is called when production api_endpoint is entered.
func (s *TreeShapeListener) EnterSimple_endpoint(ctx *parser.Simple_endpointContext) {
	if ctx.WHATEVER() != nil {
		s.currentApp().Endpoints[ctx.WHATEVER().GetText()] = &sysl.Endpoint{
			Name: ctx.WHATEVER().GetText(),
		}
		return
	}
	s.endpointName = ctx.Endpoint_name().GetText()
	s.recordEndpoint(s.endpointName, s.createLocation(ctx.GetStart()))
	ep := s.currentApp().Endpoints[s.endpointName]

	if ep == nil {
		ep = &sysl.Endpoint{
			Name: s.endpointName,
		}
		s.currentApp().Endpoints[s.endpointName] = ep
	}

	if ctx.QSTRING() != nil {
		ep.LongName = fromQString(ctx.QSTRING().GetText())
	}

	if ctx.Attribs_or_modifiers() != nil {
		attrs := s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		if ep.Attrs == nil {
			ep.Attrs = attrs
		} else {
			mergeAttrs(attrs, ep.Attrs)
		}
	}
	if ep.Attrs == nil {
		ep.Attrs = map[string]*sysl.Attribute{}
	}

	if ctx.Statements(0) != nil {
		if ep.Stmt == nil {
			ep.Stmt = []*sysl.Statement{}
		}

		s.pushScope(s.currentApp().Endpoints[s.endpointName])
	}
	ep.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)       //nolint:staticcheck
	ep.SourceContexts = append(ep.SourceContexts, ep.SourceContext) //nolint:staticcheck
}

// ExitSimple_endpoint is called when production simple_endpoint is exited.
func (s *TreeShapeListener) ExitSimple_endpoint(ctx *parser.Simple_endpointContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
	ep := s.currentApp().Endpoints[s.endpointName]
	if ep != nil {
		if len(ep.Attrs) == 0 {
			ep.Attrs = nil
		}

		ep.SourceContext.End = s.lastEnd //nolint:staticcheck
		ep.SourceContexts[len(ep.SourceContexts)-1].End = s.lastEnd
	}
	s.endpointName = ""
}

// EnterRest_endpoint is called when production rest_endpoint is entered.
func (s *TreeShapeListener) EnterRest_endpoint(ctx *parser.Rest_endpointContext) {
	s.rest_queryparams_len = append(s.rest_queryparams_len, len(s.rest_queryparams))
	s.rest_urlparams_len = append(s.rest_urlparams_len, len(s.rest_urlparams))

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		s.rest_attrs = append(s.rest_attrs, s.makeAttributeArray(attribs))
	} else {
		s.rest_attrs = append(s.rest_attrs, nil)
	}
}

// ExitRest_endpoint is called when production rest_endpoint is exited.
func (s *TreeShapeListener) ExitRest_endpoint(*parser.Rest_endpointContext) {
	s.urlPrefixes.Pop()
	ltop := len(s.rest_urlparams_len) - 1
	if ltop >= 0 {
		l := s.rest_urlparams_len[ltop]
		s.rest_urlparams = s.rest_urlparams[:l]
		s.rest_urlparams_len = s.rest_urlparams_len[:ltop]
	}
	s.rest_attrs = s.rest_attrs[:len(s.rest_attrs)-1]

	if len(s.urlPrefixes.parts) != len(s.rest_attrs) {
		panic("something is wrong")
	}
}

// EnterCollector_call_stmt is called when production collector_call_stmt is entered.
func (s *TreeShapeListener) EnterCollector_call_stmt(ctx *parser.Collector_call_stmtContext) {
	appName := &sysl.AppName{}
	s.app_name.Reset()
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: MustUnescape(ctx.Target_endpoint().GetText()),
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	})
}

// EnterCollector_http_stmt is called when production collector_http_stmt is entered.
func (s *TreeShapeListener) EnterCollector_http_stmt(ctx *parser.Collector_http_stmtContext) {
	text := strings.TrimSpace(ctx.HTTP_VERBS().GetText()) + " " +
		MustUnescape(ctx.Collector_http_stmt_suffix().GetText())

	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Action{
			Action: &sysl.Action{
				Action: text,
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	})
}

// ExitSubscriber is called when production subscriber is exited.
func (s *TreeShapeListener) ExitSubscriber(*parser.SubscriberContext) {
	s.lastStatement().GetCall().Target.Part = s.app_name.Parts()
	s.app_name.Reset()
}

// EnterCollector_pubsub_call is called when production collector_pubsub_call is entered.
func (s *TreeShapeListener) EnterCollector_pubsub_call(ctx *parser.Collector_pubsub_callContext) {
	appName := &sysl.AppName{}
	s.app_name.Reset()
	publisher := ctx.Publisher().GetText() + ctx.ARROW_RIGHT().GetText()
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: strings.TrimSpace(publisher + ctx.Name_str().GetText()),
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	})
}

// EnterCollector_action_stmt is called when production collector_action_stmt is entered.
func (s *TreeShapeListener) EnterCollector_action_stmt(ctx *parser.Collector_action_stmtContext) {
	text := ctx.Name_str().GetText()
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Action{
			Action: &sysl.Action{
				Action: text,
			},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	})
}

// ExitCollector_stmts is called when production collector_stmts is exited.
func (s *TreeShapeListener) ExitCollector_stmts(ctx *parser.Collector_stmtsContext) {
	if ctx.Attribs_or_modifiers() != nil {
		stmt := s.lastStatement()
		stmt.Attrs = s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}
}

// EnterCollector is called when production collector is entered.
func (s *TreeShapeListener) EnterCollector(ctx *parser.CollectorContext) {
	collector := ctx.COLLECTOR().GetText()
	ep := s.currentApp().Endpoints[collector]

	if ep == nil {
		ep = &sysl.Endpoint{
			Name: collector,
		}
		s.currentApp().Endpoints[collector] = ep
	}

	if ctx.Collector_stmts(0) != nil {
		ep.Stmt = []*sysl.Statement{}
		if ep.Attrs == nil {
			ep.Attrs = map[string]*sysl.Attribute{}
		}
		s.pushScope(ep)
	}
	ep.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)       //nolint:staticcheck
	ep.SourceContexts = append(ep.SourceContexts, ep.SourceContext) //nolint:staticcheck
}

// ExitCollector is called when production collector is exited.
func (s *TreeShapeListener) ExitCollector(ctx *parser.CollectorContext) {
	collector := ctx.COLLECTOR().GetText()
	if ctx.Collector_stmts(0) != nil {
		s.popScope()
	}
	ep := s.currentApp().Endpoints[collector]
	if len(ep.Attrs) == 0 {
		ep.Attrs = nil
	}
}

// EnterEvent is called when production event is entered.
func (s *TreeShapeListener) EnterEvent(ctx *parser.EventContext) {
	if ctx.Name_str() != nil {
		s.endpointName = ctx.Name_str().GetText()
		// fmt.Printf("Event: %s\n", s.endpointName)
		ep := s.currentApp().Endpoints[s.endpointName]
		if ep == nil {
			ep = &sysl.Endpoint{
				Name:           s.endpointName,
				IsPubsub:       true,
				SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
				SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
			}
			s.currentApp().Endpoints[s.endpointName] = ep
		}
		if ctx.Attribs_or_modifiers() != nil {
			ep.Attrs = s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
		if ctx.Statements(0) != nil && ep.Stmt == nil {
			ep.Stmt = []*sysl.Statement{}
		}
		if ctx.Statements(0) != nil {
			s.pushScope(ep)
		}
	}
}

// ExitEvent is called when production event is exited.
func (s *TreeShapeListener) ExitEvent(ctx *parser.EventContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
	s.endpointName = ""
}

// EnterSubscribe is called when production subscribe is entered.
func (s *TreeShapeListener) EnterSubscribe(ctx *parser.SubscribeContext) {
	if ctx.App_name() != nil {
		eventName := ctx.Name_str().GetText()
		s.endpointName = syslutil.CleanAppName(MustUnescape(ctx.App_name().GetText())) +
			ctx.ARROW_RIGHT().GetText() + eventName
		// fmt.Printf("\t%s Subscriber: %s\n", s.appName, s.endpointName)
		str := strings.Split(ctx.App_name().GetText(), "::")
		for i := range str {
			str[i] = MustUnescape(strings.TrimSpace(str[i]))
		}
		app_src := &sysl.AppName{
			Part: str,
		}
		typeEndpoint := &sysl.Endpoint{
			Name:           s.endpointName,
			Source:         app_src,
			Attrs:          map[string]*sysl.Attribute{},
			SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
			SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
		}
		if ctx.Attribs_or_modifiers() != nil {
			typeEndpoint.Attrs = s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
		if ctx.Statements(0) != nil {
			typeEndpoint.Stmt = []*sysl.Statement{}
			s.pushScope(typeEndpoint)
		}
		s.currentApp().Endpoints[s.endpointName] = typeEndpoint
		srcAppName := syslutil.GetAppName(app_src)
		srcApp := syslutil.GetApp(app_src, s.module)
		if srcApp == nil {
			s.module.Apps[srcAppName] = &sysl.Application{
				Name:      app_src,
				Endpoints: map[string]*sysl.Endpoint{},
			}
			srcApp = s.module.Apps[srcAppName]
		}
		ep := srcApp.Endpoints[eventName]
		if ep == nil {
			srcApp.Endpoints[eventName] = &sysl.Endpoint{
				Name:     eventName,
				Stmt:     []*sysl.Statement{},
				IsPubsub: true,
			}
			ep = srcApp.Endpoints[eventName]
		}
		if ep.Stmt == nil {
			ep.Stmt = []*sysl.Statement{}
		}
		stmt := &sysl.Statement{
			Stmt: &sysl.Statement_Call{
				Call: &sysl.Call{
					Target:   s.currentApp().Name,
					Endpoint: s.endpointName,
				},
			},
			SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
			SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
		}
		ep.Stmt = append(ep.Stmt, stmt)
	}
}

// ExitSubscribe is called when production subscribe is exited.
func (s *TreeShapeListener) ExitSubscribe(ctx *parser.SubscribeContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
	s.endpointName = ""
}

// ExitLiteral is called when production literal is exited.
func (s *TreeShapeListener) ExitLiteral(ctx *parser.LiteralContext) {
	val := &sysl.Value{}
	txt := ctx.GetText()
	var type1 *sysl.Type
	switch {
	case ctx.E_DECIMAL() != nil:
		val.Value = &sysl.Value_Decimal{
			Decimal: txt,
		}
	case ctx.E_DIGITS() != nil:
		iVal, err := strconv.ParseInt(txt, 10, 0)
		syslutil.PanicOnError(err)
		val.Value = &sysl.Value_I{
			I: iVal,
		}
	case ctx.E_TRUE() != nil || ctx.E_FALSE() != nil:
		val.Value = &sysl.Value_B{
			B: txt == "true",
		}
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_BOOL,
			},
		}

	case ctx.E_NULL() != nil:
		val.Value = &sysl.Value_Null_{
			Null: &sysl.Value_Null{},
		}
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_EMPTY,
			},
		}
	default:
		val.Value = &sysl.Value_S{
			S: fromQString(txt),
		}
	}

	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Literal{
		Literal: val,
	}
	if type1 != nil {
		expr.Type = type1
	}
	s.PushExpr(expr)
}

func makeGetAttr(arg *sysl.Expr, name string, nullsafe, setof bool, ctx *sysl.SourceContext) *sysl.Expr {
	expr := makeExpr(ctx)
	expr.Expr = &sysl.Expr_GetAttr_{
		GetAttr: &sysl.Expr_GetAttr{
			Arg:      arg,
			Attr:     name,
			Setof:    setof,
			Nullsafe: nullsafe,
		},
	}
	return expr
}

func makeLiteralString(val string, ctx *sysl.SourceContext) *sysl.Expr {
	expr := makeExpr(ctx)
	expr.Expr = &sysl.Expr_Literal{Literal: &sysl.Value{Value: &sysl.Value_S{S: val}}}
	return expr
}

func makeRelExpr(op sysl.Expr_RelExpr_Op, ctx *sysl.SourceContext) *sysl.Expr {
	expr := makeExpr(ctx)
	expr.Expr = &sysl.Expr_Relexpr{
		Relexpr: &sysl.Expr_RelExpr{
			Op: op,
		},
	}
	return expr
}

func makeExprName(name string, ctx *sysl.SourceContext) *sysl.Expr {
	expr := makeExpr(ctx)
	expr.Expr = &sysl.Expr_Name{
		Name: name,
	}
	return expr
}

func makeUnaryExpr(op sysl.Expr_UnExpr_Op, arg *sysl.Expr, ctx *sysl.SourceContext) *sysl.Expr {
	expr := makeExpr(ctx)
	expr.Expr = &sysl.Expr_Unexpr{
		Unexpr: &sysl.Expr_UnExpr{
			Op:  op,
			Arg: arg,
		},
	}
	return expr
}

func makeBinaryExpr(op sysl.Expr_BinExpr_Op, lhs, rhs *sysl.Expr, ctx *sysl.SourceContext) *sysl.Expr {
	expr := makeExpr(ctx)
	expr.Expr = &sysl.Expr_Binexpr{
		Binexpr: &sysl.Expr_BinExpr{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		},
	}
	return expr
}

func addStmt(ifelse *sysl.Expr, stmt *sysl.Expr) {
	if ifelse == nil {
		return
	}
	if ifelse.GetIfelse().IfTrue == nil {
		ifelse.GetIfelse().IfTrue = stmt
	}
	addStmt(ifelse.GetIfelse().IfFalse, stmt)
}

func makeIfElseExpr() *sysl.Expr {
	return &sysl.Expr{
		Expr: &sysl.Expr_Ifelse{
			Ifelse: &sysl.Expr_IfElse{},
		},
	}
}

func addIfElseControl(ifelse *sysl.Expr, cond *sysl.Expr) {
	root := ifelse.GetIfelse()

	if root.IfFalse == nil {
		root.IfFalse = makeIfElseExpr()
		root = root.IfFalse.GetIfelse()
		root.Cond = cond
		return
	}
	addIfElseControl(root.IfFalse, cond)
}

// ExitExpr_table_of_op is called when production expr_table_of_op is exited.
func (s *TreeShapeListener) ExitExpr_table_of_op(ctx *parser.Expr_table_of_opContext) {
	nullsafe := ctx.E_NULLSAFE_DOT() != nil
	table_of := ctx.E_TABLE_OF() != nil
	attrName := ctx.E_Name().GetText()
	arg := s.popExpr()
	s.PushExpr(makeGetAttr(arg, attrName, nullsafe, table_of, s.getSrcCtx(ctx.BaseParserRuleContext)))
}

// ExitFunc_arg is called when production func_arg is exited.
func (s *TreeShapeListener) ExitFunc_arg(*parser.Func_argContext) {
	arg := s.popExpr()
	top := s.TopExpr()
	var args *[]*sysl.Expr
	switch {
	case top.GetRelexpr() != nil:
		args = &top.GetRelexpr().Arg
	case top.GetCall() != nil:
		args = &top.GetCall().Arg
	}
	syslutil.Assert(args != nil, "ExitFunc_arg: should not be here") // nolint:staticcheck
	*args = append(*args, arg)                                       // nolint:staticcheck
}

// EnterExpr_func is called when production expr_func is entered.
func (s *TreeShapeListener) EnterExpr_func(ctx *parser.Expr_funcContext) {
	var funcName string
	switch {
	case ctx.E_FUNC() != nil:
		funcName = ctx.E_FUNC().GetText()
	case ctx.E_Name() != nil:
		funcName = ctx.E_Name().GetText()
	default:
		funcName = ctx.NativeDataTypes().GetText()
	}
	expr := makeExpr(s.getSrcCtxWithText(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Call_{
		Call: &sysl.Expr_Call{
			Func: funcName,
			Arg:  []*sysl.Expr{},
		},
	}
	s.PushExpr(expr)
}

// ExitRank_expr is called when production rank_expr is exited.
func (s *TreeShapeListener) ExitRank_expr(ctx *parser.Rank_exprContext) {
	expr := s.popExpr()
	relexpr := s.TopExpr().GetRelexpr()
	relexpr.Arg = append(relexpr.Arg, expr)
	relexpr.Descending = append(relexpr.Descending, ctx.E_DESC() != nil)
}

// EnterExpr_rank_func is called when production expr_rank_func is entered.
func (s *TreeShapeListener) EnterExpr_rank_func(ctx *parser.Expr_rank_funcContext) {
	target := s.popExpr()
	s.PushExpr(makeRelExpr(sysl.Expr_RelExpr_RANK, s.getSrcCtx(ctx.BaseParserRuleContext)))
	relexpr := s.TopExpr().GetRelexpr()
	relexpr.AttrName = append(relexpr.AttrName, ctx.E_Name().GetText())
	relexpr.Target = target
	s.fieldname = append(s.fieldname, "rank-type-spec")
	type1 := &sysl.Type{}
	if ctx.View_type_spec() == nil {
		type1.Type = &sysl.Type_Set{
			Set: &sysl.Type{},
		}
	}
	s.typemap[s.fieldname[len(s.fieldname)-1]] = type1
}

// ExitExpr_rank_func is called when production expr_rank_func is exited.
func (s *TreeShapeListener) ExitExpr_rank_func(*parser.Expr_rank_funcContext) {
	expr := s.TopExpr()
	relexpr := expr.GetRelexpr()
	if relexpr.Scopevar == "" {
		relexpr.Scopevar = "."
	}

	// Legacy has double the length of the arg array.
	relexpr.Arg = append(relexpr.Arg, relexpr.Arg...)

	expr.Type = s.typemap[s.fieldname[len(s.fieldname)-1]]
	if expr.Type.GetTypeRef() != nil {
		expr.Type.GetTypeRef().Context = nil
		setType := &sysl.Type{
			Type: &sysl.Type_Set{
				Set: expr.Type,
			},
		}
		expr.Type = setType
	}
}

// EnterExpr_agg_func is called when production expr_agg_func is entered.
func (s *TreeShapeListener) EnterExpr_agg_func(ctx *parser.Expr_agg_funcContext) {
	var op sysl.Expr_RelExpr_Op
	switch ctx.E_RELOPS_AGG().GetText() {
	case "min":
		op = sysl.Expr_RelExpr_MIN
	case "max":
		op = sysl.Expr_RelExpr_MAX
	case "sum":
		op = sysl.Expr_RelExpr_SUM
	case "average":
		op = sysl.Expr_RelExpr_AVERAGE
	}
	expr := s.popExpr()
	s.PushExpr(makeRelExpr(op, s.getSrcCtx(ctx.BaseParserRuleContext)))
	relexpr := s.TopExpr().GetRelexpr()
	relexpr.Target = expr
}

// ExitExpr_agg_func is called when production expr_agg_func is exited.
func (s *TreeShapeListener) ExitExpr_agg_func(*parser.Expr_agg_funcContext) {
	relexpr := s.TopExpr().GetRelexpr()
	if relexpr.Scopevar == "" {
		relexpr.Scopevar = "."
	}
}

// ExitE_scope_var is called when production e_scope_var is exited.
func (s *TreeShapeListener) ExitE_scope_var(ctx *parser.E_scope_varContext) {
	expr := s.TopExpr()
	if expr.GetBinexpr() != nil {
		binexpr := expr.GetBinexpr()
		binexpr.Scopevar = ctx.E_Name().GetText()
	} else if expr.GetRelexpr() != nil {
		relexpr := expr.GetRelexpr()
		relexpr.Scopevar = ctx.E_Name().GetText()
	}
}

// ExitFirst_func_target is called when production first_func_target is exited.
func (s *TreeShapeListener) ExitFirst_func_target(ctx *parser.First_func_targetContext) {
	var expr *sysl.Expr
	if ctx.E_NULL() != nil {
		val := &sysl.Value{}
		val.Value = &sysl.Value_Null_{
			Null: &sysl.Value_Null{},
		}

		expr = makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
		expr.Expr = &sysl.Expr_Literal{
			Literal: val,
		}
		expr.Type = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_EMPTY,
			},
		}
	} else {
		expr = s.popExpr()
	}
	relexpr := s.TopExpr().GetRelexpr()
	relexpr.Arg = append(relexpr.Arg, expr)
}

// EnterExpr_first_func is called when production expr_first_func is entered.
func (s *TreeShapeListener) EnterExpr_first_func(ctx *parser.Expr_first_funcContext) {
	expr := makeRelExpr(sysl.Expr_RelExpr_FIRST_BY, s.getSrcCtx(ctx.BaseParserRuleContext))
	relexpr := expr.GetRelexpr()
	relexpr.Target = s.popExpr()
	s.PushExpr(expr)
}

// ExitExpr_first_func is called when production expr_first_func is exited.
func (s *TreeShapeListener) ExitExpr_first_func(*parser.Expr_first_funcContext) {
	relexpr := s.TopExpr().GetRelexpr()
	if relexpr.Scopevar == "" {
		relexpr.Scopevar = "."
	}
	// Legacy has double the length of the arg array.
	relexpr.Arg = append(relexpr.Arg, relexpr.Arg...)
}

// ExitE_single_arg_func is called when production e_single_arg_func is exited.
func (s *TreeShapeListener) ExitE_single_arg_func(ctx *parser.E_single_arg_funcContext) {
	op := sysl.Expr_BinExpr_NO_Op
	if ctx.E_RELOPS_WHERE() != nil {
		op = sysl.Expr_BinExpr_WHERE
	} else if ctx.E_RELOPS_FLATTEN() != nil {
		op = sysl.Expr_BinExpr_FLATTEN
	}
	s.PushExpr(makeBinaryExpr(op, nil, nil, s.getSrcCtx(ctx.BaseParserRuleContext)))
}

// ExitExpr_single_arg_func is called when production expr_single_arg_func is exited.
func (s *TreeShapeListener) ExitExpr_single_arg_func(ctx *parser.Expr_single_arg_funcContext) {
	rhs := s.popExpr()
	expr := s.popExpr()
	lhs := s.popExpr()
	binexpr := expr.GetBinexpr()
	if ctx.E_scope_var() == nil {
		binexpr.Scopevar = "."
	}
	binexpr.Lhs = lhs
	binexpr.Rhs = rhs
	s.PushExpr(expr)
}

// ExitExpr_any_func is called when production expr_any_func is exited.
func (s *TreeShapeListener) ExitExpr_any_func(ctx *parser.Expr_any_funcContext) {
	limit := s.popExpr()
	arg := s.popExpr()
	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Call_{
		Call: &sysl.Expr_Call{
			Func: ".any",
			Arg:  []*sysl.Expr{arg, limit},
		},
	}
	s.PushExpr(expr)
}

// ExitExpr_single_or_null is called when production expr_single_or_null is exited.
func (s *TreeShapeListener) ExitExpr_single_or_null(ctx *parser.Expr_single_or_nullContext) {
	op := sysl.Expr_UnExpr_SINGLE_OR_NULL
	if ctx.GetText() == "single" {
		op = sysl.Expr_UnExpr_SINGLE
	}
	expr := makeUnaryExpr(op, s.popExpr(), s.getSrcCtx(ctx.BaseParserRuleContext))
	s.PushExpr(expr)
}

// ExitExpr_snapshot is called when production expr_snapshot is exited.
func (s *TreeShapeListener) ExitExpr_snapshot(ctx *parser.Expr_snapshotContext) {
	expr := s.popExpr()
	s.PushExpr(makeRelExpr(sysl.Expr_RelExpr_SNAPSHOT, s.getSrcCtx(ctx.BaseParserRuleContext)))
	relexpr := s.TopExpr().GetRelexpr()
	relexpr.Target = expr
}

// ExitExpr_count is called when production expr_count is exited.
func (s *TreeShapeListener) ExitExpr_count(ctx *parser.Expr_countContext) {
	expr := s.popExpr()
	callExpr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	callExpr.Expr = &sysl.Expr_Call_{
		Call: &sysl.Expr_Call{
			Func: ".count",
			Arg:  []*sysl.Expr{expr},
		},
	}
	s.PushExpr(callExpr)
}

// ExitExpr_navigate_attr is called when production expr_navigate_attr is exited.
func (s *TreeShapeListener) ExitExpr_navigate_attr(ctx *parser.Expr_navigate_attrContext) {
	nav := s.TopExpr().GetNavigate()
	if ctx.E_DOT() != nil {
		nav.Attr = "."
	}
	nav.Attr += ctx.E_Name().GetText()
}

// EnterExpr_navigate is called when production expr_navigate is entered.
func (s *TreeShapeListener) EnterExpr_navigate(ctx *parser.Expr_navigateContext) {
	arg := s.popExpr()

	nav := &sysl.Expr_Navigate{
		Arg: arg,
	}

	if ctx.E_QN() != nil {
		nav.Nullsafe = true
	}

	if ctx.E_SET_OF() != nil {
		nav.Setof = true
	}

	if ctx.E_Name() != nil {
		nav.Via = ctx.E_Name().GetText()
	}

	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Navigate_{
		Navigate: nav,
	}
	s.PushExpr(expr)
}

// EnterMatching_rhs is called when production matching_rhs is entered.
func (s *TreeShapeListener) EnterMatching_rhs(ctx *parser.Matching_rhsContext) {
	if ctx.E_Name() == nil && ctx.AtomT_paren() == nil {
		s.PushExpr(makeExprName(".", s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitMatching_rhs is called when production matching_rhs is exited.
func (s *TreeShapeListener) ExitMatching_rhs(ctx *parser.Matching_rhsContext) {
	if ctx.E_Name() != nil {
		s.PushExpr(makeExprName(ctx.E_Name().GetText(), s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// EnterSquiggly_args is called when production squiggly_args is entered.
func (s *TreeShapeListener) EnterSquiggly_args(ctx *parser.Squiggly_argsContext) {
	binexpr := s.TopExpr().GetBinexpr()
	names := ctx.AllE_Name()
	for _, name := range names {
		binexpr.AttrName = append(binexpr.AttrName, name.GetText())
	}
}

// EnterExpr_matching is called when production expr_matching is entered.
func (s *TreeShapeListener) EnterExpr_matching(ctx *parser.Expr_matchingContext) {
	op := sysl.Expr_BinExpr_TO_MATCHING
	if ctx.E_NOT() != nil {
		op = sysl.Expr_BinExpr_TO_NOT_MATCHING
	}
	lhs := s.popExpr()
	s.PushExpr(makeBinaryExpr(op, lhs, nil, s.getSrcCtx(ctx.BaseParserRuleContext)))
}

// ExitExpr_matching is called when production expr_matching is exited.
func (s *TreeShapeListener) ExitExpr_matching(*parser.Expr_matchingContext) {
	var binexpr *sysl.Expr_BinExpr
	rhs := s.popExpr()
	binexpr = s.TopExpr().GetBinexpr()
	binexpr.Rhs = rhs
	if len(binexpr.AttrName) == 0 {
		binexpr.AttrName = []string{"*"}
	}
}

// ExitList_item is called when production list_item is exited.
func (s *TreeShapeListener) ExitList_item(*parser.List_itemContext) {
	item := s.popExpr()

	if s.TopExpr().GetSet() != nil {
		list := s.TopExpr().GetSet()
		list.Expr = append(list.Expr, item)
	} else if s.TopExpr().GetList() != nil {
		list := s.TopExpr().GetList()
		list.Expr = append(list.Expr, item)
	}
}

// ExitExpr_list is called when production expr_list is exited.
func (s *TreeShapeListener) ExitExpr_list(*parser.Expr_listContext) {
	expr := s.TopExpr()
	list := expr.GetSet()
	if list != nil && len(list.Expr) == 1 && list.Expr[0].Type.GetTuple() != nil {
		type1 := expr.Type.GetSet()
		type1.Type = &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{},
		}
	}
}

// EnterExpr_set is called when production expr_set is entered.
func (s *TreeShapeListener) EnterExpr_set(ctx *parser.Expr_setContext) {
	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Set{
		Set: &sysl.Expr_List{
			Expr: []*sysl.Expr{},
		},
	}
	expr.Type = &sysl.Type{
		Type: &sysl.Type_Set{
			Set: &sysl.Type{},
		},
	}
	s.PushExpr(expr)
}

// ExitEmpty_tuple is called when production empty_tuple is exited.
func (s *TreeShapeListener) ExitEmpty_tuple(ctx *parser.Empty_tupleContext) {
	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Tuple_{
		Tuple: &sysl.Expr_Tuple{},
	}
	expr.Type = &sysl.Type{
		Type: &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{},
		},
	}
	s.PushExpr(expr)
}

// EnterAtomT_implied_dot is called when production atomT_implied_dot is entered.
func (s *TreeShapeListener) EnterAtomT_implied_dot(ctx *parser.AtomT_implied_dotContext) {
	s.PushExpr(makeExprName(".", s.getSrcCtx(ctx.BaseParserRuleContext)))
}

// ExitAtomT_name is called when production atomT_name is exited.
func (s *TreeShapeListener) ExitAtomT_name(ctx *parser.AtomT_nameContext) {
	switch {
	case ctx.E_Name() != nil:
		txt := ctx.E_Name().GetText()
		s.PushExpr(makeExprName(txt, s.getSrcCtx(ctx.BaseParserRuleContext)))
	case ctx.E_WHATEVER() != nil:
		s.PushExpr(makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext)))
	case ctx.E_DOT() != nil:
		s.PushExpr(makeExprName(".", s.getSrcCtx(ctx.BaseParserRuleContext)))
	default:
		panic("ExitAtomT_name: should not be here")
	}
}

// EnterExpr_atom_list is called when production expr_atom_list is entered.
func (s *TreeShapeListener) EnterExpr_atom_list(ctx *parser.Expr_atom_listContext) {
	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_List_{
		List: &sysl.Expr_List{
			Expr: []*sysl.Expr{},
		},
	}
	s.PushExpr(expr)
}

// ExitPower is called when production power is exited.
func (s *TreeShapeListener) ExitPower(ctx *parser.PowerContext) {
	if ctx.PowerT() != nil {
		rhs := s.popExpr()
		lhs := s.popExpr()
		s.PushExpr(makeBinaryExpr(sysl.Expr_BinExpr_POW, lhs, rhs, s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitUnaryTerm is called when production unaryTerm is exited.
func (s *TreeShapeListener) ExitUnaryTerm(ctx *parser.UnaryTermContext) {
	op := sysl.Expr_UnExpr_NO_Op

	switch {
	case ctx.E_TILDE() != nil:
		op = sysl.Expr_UnExpr_INV
	case ctx.E_NOT() != nil:
		op = sysl.Expr_UnExpr_NOT
	case ctx.E_MINUS() != nil:
		op = sysl.Expr_UnExpr_NEG
	case ctx.E_PLUS() != nil:
		op = sysl.Expr_UnExpr_POS
	}
	if op != sysl.Expr_UnExpr_NO_Op {
		expr := makeUnaryExpr(op, s.popExpr(), s.getSrcCtx(ctx.BaseParserRuleContext))
		s.PushExpr(expr)
	}
}

// ExitTermT is called when production termT is exited.
func (s *TreeShapeListener) ExitTermT(ctx *parser.TermTContext) {
	op := sysl.Expr_BinExpr_MOD
	switch {
	case ctx.E_STAR() != nil:
		op = sysl.Expr_BinExpr_MUL
	case ctx.E_DIVIDE() != nil:
		op = sysl.Expr_BinExpr_DIV
	}
	rhs := s.popExpr()
	lhs := s.popExpr()

	s.PushExpr(makeBinaryExpr(op, lhs, rhs, s.getSrcCtx(ctx.BaseParserRuleContext)))
}

// ExitBinexprT is called when production binexprT is exited.
func (s *TreeShapeListener) ExitBinexprT(ctx *parser.BinexprTContext) {
	var op sysl.Expr_BinExpr_Op

	if ctx.E_PLUS() != nil {
		op = sysl.Expr_BinExpr_ADD
	} else {
		op = sysl.Expr_BinExpr_SUB
	}
	rhs := s.popExpr()
	lhs := s.popExpr()
	s.PushExpr(makeBinaryExpr(op, lhs, rhs, s.getSrcCtx(ctx.BaseParserRuleContext)))
}

// ExitExpr_rel is called when production expr_rel is exited.
func (s *TreeShapeListener) ExitExpr_rel(ctx *parser.Expr_relContext) {
	if ctx.E_compare_ops(0) != nil {
		for i := len(ctx.AllE_compare_ops()) - 1; i >= 0; i-- {
			op := s.opmap[ctx.E_compare_ops(i).GetText()]
			rhs := s.popExpr()
			lhs := s.popExpr()
			s.PushExpr(makeBinaryExpr(op, lhs, rhs, s.getSrcCtx(ctx.BaseParserRuleContext)))
		}
	}
}

// ExitExpr_bitand is called when production expr_bitand is exited.
func (s *TreeShapeListener) ExitExpr_bitand(ctx *parser.Expr_bitandContext) {
	if ctx.E_AMP(0) != nil || ctx.E_AND(0) != nil {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_BITAND, len(ctx.AllExpr_rel())-1, s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitExpr_bitxor is called when production expr_bitxor is exited.
func (s *TreeShapeListener) ExitExpr_bitxor(ctx *parser.Expr_bitxorContext) {
	if ctx.E_XOR(0) != nil {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_BITXOR, len(ctx.AllE_XOR()), s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitExpr_bitor is called when production expr_bitor is exited.
func (s *TreeShapeListener) ExitExpr_bitor(ctx *parser.Expr_bitorContext) {
	if ctx.E_BITOR(0) != nil {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_BITOR, len(ctx.AllE_BITOR()), s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitExpr_and is called when production expr_and is exited.
func (s *TreeShapeListener) ExitExpr_and(ctx *parser.Expr_andContext) {
	if ctx.E_DOUBLE_AMP(0) != nil {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_AND, len(ctx.AllE_DOUBLE_AMP()), s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitExpr_or is called when production expr_or is exited.
func (s *TreeShapeListener) ExitExpr_or(ctx *parser.Expr_orContext) {
	if ctx.E_LOGIC_OR(0) != nil {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_OR, len(ctx.AllE_LOGIC_OR()), s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitExpr_but_not is called when production expr_but_not is exited.
func (s *TreeShapeListener) ExitExpr_but_not(ctx *parser.Expr_but_notContext) {
	if ctx.E_BUTNOT(0) != nil {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_BUTNOT, len(ctx.AllE_BUTNOT()), s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

func (s *TreeShapeListener) reverseOp(op sysl.Expr_BinExpr_Op, count int, ctx *sysl.SourceContext) *sysl.Expr {
	if count == 0 {
		return s.popExpr()
	}
	rhs := s.popExpr()
	lhs := s.reverseOp(op, count-1, ctx)
	return makeBinaryExpr(op, lhs, rhs, ctx)
}

// ExitExpr_coalesce is called when production expr_coalesce is exited.
func (s *TreeShapeListener) ExitExpr_coalesce(ctx *parser.Expr_coalesceContext) {
	exprs := ctx.AllExpr_but_not()
	if len(exprs) > 1 {
		s.PushExpr(s.reverseOp(sysl.Expr_BinExpr_COALESCE, len(exprs)-1, s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

// ExitIf_one_liner is called when production if_one_liner is exited.
func (s *TreeShapeListener) ExitIf_one_liner(ctx *parser.If_one_linerContext) {
	if_false := s.popExpr()
	if_true := s.popExpr()
	cond := s.popExpr()
	expr := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	expr.Expr = &sysl.Expr_Ifelse{
		Ifelse: &sysl.Expr_IfElse{
			Cond:     cond,
			IfTrue:   if_true,
			IfFalse:  if_false,
			Nullsafe: ctx.E_QN() != nil,
		},
	}
	s.PushExpr(expr)
}

// ExitElse_block_stmt is called when production else_block_stmt is exited.
func (s *TreeShapeListener) ExitElse_block_stmt(*parser.Else_block_stmtContext) {
	stmt := s.popExpr()
	ifelse := s.TopExpr()
	addStmt(ifelse, stmt)
}

// ExitControl_item is called when production control_item is exited.
func (s *TreeShapeListener) ExitControl_item(ctx *parser.Control_itemContext) {
	control := s.popExpr()
	expr := s.TopExpr()
	ifelse := expr.GetIfelse()

	if ifelse.Cond.GetBinexpr() != nil && ifelse.Cond.GetBinexpr().Op == sysl.Expr_BinExpr_EQ {
		if ifelse.Cond.GetBinexpr().Rhs == nil {
			ifelse.Cond.GetBinexpr().Rhs = control
		} else {
			lhs := ifelse.Cond.GetBinexpr().Lhs
			newCond := makeBinaryExpr(sysl.Expr_BinExpr_EQ, lhs, control, s.getSrcCtx(ctx.BaseParserRuleContext))
			addIfElseControl(expr, newCond)
		}
	} else {
		if len(ifelse.Cond.GetCall().Arg) == 0 {
			ifelse.Cond.GetCall().Arg = append(ifelse.Cond.GetCall().Arg, control)
		} else {
			newCond := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
			newCond.Expr = &sysl.Expr_Call_{
				Call: &sysl.Expr_Call{
					Func: "bool",
					Arg:  []*sysl.Expr{control},
				},
			}
			addIfElseControl(expr, newCond)
		}
	}
}

// ExitFinal_else is called when production final_else is exited.
func (s *TreeShapeListener) ExitFinal_else(*parser.Final_elseContext) {
	elseStmt := s.popExpr()
	ifelse := s.TopExpr()
	for ifelse.GetIfelse().IfFalse != nil {
		ifelse = ifelse.GetIfelse().IfFalse
	}
	ifelse.GetIfelse().IfFalse = elseStmt
}

// ExitIfvar is called when production ifvar is exited.
func (s *TreeShapeListener) ExitIfvar(ctx *parser.IfvarContext) {
	lhs := s.popExpr()
	ifelse := s.TopExpr()
	ifelse.GetIfelse().Cond = makeBinaryExpr(sysl.Expr_BinExpr_EQ, lhs, nil, s.getSrcCtx(ctx.BaseParserRuleContext))
}

// EnterIf_multiple_lines is called when production if_multiple_lines is entered.
func (s *TreeShapeListener) EnterIf_multiple_lines(ctx *parser.If_multiple_linesContext) {
	expr := makeIfElseExpr()
	if ctx.Ifvar() == nil {
		newCond := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
		newCond.Expr = &sysl.Expr_Call_{
			Call: &sysl.Expr_Call{
				Func: "bool",
				Arg:  []*sysl.Expr{},
			},
		}
		expr.GetIfelse().Cond = newCond
	}
	s.PushExpr(expr)
}

// ExitExpr_simple_assign is called when production expr_simple_assign is exited.
func (s *TreeShapeListener) ExitExpr_simple_assign(ctx *parser.Expr_simple_assignContext) {
	expr := s.popExpr()
	sc := s.getSrcCtxWithText(ctx.BaseParserRuleContext)
	expr.SourceContext = sc                         //nolint:staticcheck
	expr.SourceContexts = []*sysl.SourceContext{sc} //nolint:staticcheck
	tx := s.TopExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v\n", s.TopExpr())
		panic("ExitExpr_simple_assign: Unexpected expression!")
	}
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Assign_{
			Assign: &sysl.Expr_Transform_Stmt_Assign{
				Name: ctx.E_Name().GetText(),
				Expr: expr,
			},
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// ExitExpr_let_statement is called when production expr_let_statement is exited.
func (s *TreeShapeListener) ExitExpr_let_statement(ctx *parser.Expr_let_statementContext) {
	expr := s.popExpr()
	sc := s.getSrcCtxWithText(ctx.BaseParserRuleContext)
	expr.SourceContext = sc                         //nolint:staticcheck
	expr.SourceContexts = []*sysl.SourceContext{sc} //nolint:staticcheck
	tx := s.TopExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.TopExpr())
		panic("ExitExpr_let_statement: Unexpected expression!")
	}
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Let{
			Let: &sysl.Expr_Transform_Stmt_Assign{
				Name: ctx.E_Name().GetText(),
				Expr: expr,
			},
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// ExitExpr_table_of_statement is called when production expr_table_of_statement is exited.
func (s *TreeShapeListener) ExitExpr_table_of_statement(ctx *parser.Expr_table_of_statementContext) {
	expr := s.popExpr()
	tx := s.TopExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.TopExpr())
		panic("ExitExpr_table_of_statement: Unexpected expression!")
	}
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Assign_{
			Assign: &sysl.Expr_Transform_Stmt_Assign{
				Name:  ctx.E_Name().GetText(),
				Expr:  expr,
				Table: true,
			},
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// ExitExpr_dot_assign is called when production expr_dot_assign is exited.
func (s *TreeShapeListener) ExitExpr_dot_assign(ctx *parser.Expr_dot_assignContext) {
	tx := s.TopExpr().GetTransform()
	names := strings.Split(strings.TrimRight(ctx.GetText(), " \r\n"), ".")
	var arg *sysl.Expr

	if names[0] == "" {
		names[0] = "."
	}
	arg = makeExprName(names[0], s.getSrcCtx(ctx.BaseParserRuleContext))
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Assign_{
			Assign: &sysl.Expr_Transform_Stmt_Assign{
				Name: names[1],
				Expr: makeGetAttr(arg, names[1], false, false, s.getSrcCtx(ctx.BaseParserRuleContext)),
			},
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// ExitExpr_inject_stmt is called when production expr_inject_stmt is exited.
func (s *TreeShapeListener) ExitExpr_inject_stmt(ctx *parser.Expr_inject_stmtContext) {
	expr := s.popExpr()
	tx := s.TopExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.TopExpr())
		panic("ExitExpr_inject_stmt: Unexpected expression!")
	}
	expr.GetCall().Arg = append(expr.GetCall().Arg, makeExprName("out", s.getSrcCtx(ctx.BaseParserRuleContext)))
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Inject{
			Inject: expr,
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// EnterTemplate_expression is called when production template_expression is entered.
func (s *TreeShapeListener) EnterTemplate_expression(ctx *parser.Template_expressionContext) {
}

// ExitTemplate_expression is called when production template_expression is exited.
func (s *TreeShapeListener) ExitTemplate_expression(ctx *parser.Template_expressionContext) {
	var lhs, rhs *sysl.Expr

	sourceCtx := s.getSrcCtx(ctx.BaseParserRuleContext)

	if txt := ctx.TMPL_TEXT(); txt != nil {
		text := strings.NewReplacer("{{", "{", "}}", "}").Replace(txt.GetText())
		rhs = makeLiteralString(text, sourceCtx)
	} else {
		rhs = makeUnaryExpr(sysl.Expr_UnExpr_STRING, s.popExpr(), sourceCtx)
	}

	if top := s.TopExpr(); top.GetTransform() == nil {
		lhs = s.popExpr()
	} else {
		lhs = makeLiteralString("", nil)
	}
	s.PushExpr(makeBinaryExpr(sysl.Expr_BinExpr_ADD, lhs, rhs, sourceCtx))
}

// EnterTemplate_statement is called when production template_statement is entered.
func (s *TreeShapeListener) EnterTemplate_statement(ctx *parser.Template_statementContext) {
	if len(ctx.AllTemplate_expression()) == 0 {
		s.PushExpr(makeLiteralString("", s.getSrcCtx(ctx.BaseParserRuleContext)))
	}
}

const TemplateImpliedResult = "__$"
const TemplateLogString = "__$Log"

// ExitTemplate_statement is called when production template_statement is exited.
func (s *TreeShapeListener) ExitTemplate_statement(ctx *parser.Template_statementContext) {
	statements := s.popExpr()
	tx := s.TopExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.TopExpr())
		panic("ExitTemplate_statement: Unexpected expression!")
	}

	getRHS := func(expr *sysl.Expr) *sysl.Expr {
		if out, ok := expr.Expr.(*sysl.Expr_Binexpr); ok {
			return out.Binexpr.Rhs
		}
		return nil
	}
	rhs := statements
	for {
		if x := getRHS(rhs); x != nil {
			rhs = x
		} else {
			break
		}
	}
	wantNewLine := true
	if e, ok := rhs.Expr.(*sysl.Expr_Literal); ok {
		if strings.HasSuffix(e.Literal.GetS(), "\\") {
			wantNewLine = false
			e.Literal = &sysl.Value{
				Value: &sysl.Value_S{
					S: strings.TrimSuffix(e.Literal.GetS(), "\\"),
				},
			}
		}
	}

	sc := s.getSrcCtx(ctx.BaseParserRuleContext)
	if wantNewLine {
		expr := makeLiteralString("\n", sc)
		statements = makeBinaryExpr(sysl.Expr_BinExpr_ADD, statements, expr, sc)
	}

	resultName := TemplateImpliedResult
	if ctx.TMPL_DEBUG() != nil {
		resultName = TemplateLogString
	}
	output := makeExprName(resultName, sc)
	output.Type = &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING}}
	statements = makeBinaryExpr(sysl.Expr_BinExpr_ADD, output, statements, sc)
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Let{
			Let: &sysl.Expr_Transform_Stmt_Assign{
				Name: resultName,
				Expr: statements,
			},
		},
	}

	tx.Stmt = append(tx.Stmt, stmt)
}

// EnterTransform_return_type is called when production transform_return_type is entered.
func (s *TreeShapeListener) EnterTransform_return_type(ctx *parser.Transform_return_typeContext) {
	if ctx.Set_of() == nil {
		s.fieldname = append(s.fieldname, "transform-return")
		s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{}
	}
}

// ExitTransform_return_type is called when production transform_return_type is exited.
func (s *TreeShapeListener) ExitTransform_return_type(ctx *parser.Transform_return_typeContext) {
	if ctx.Set_of() == nil {
		type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
		expr := s.TopExpr()
		expr.Type = type1
		if type1.GetSet() != nil {
			type1.SourceContext = nil //nolint:staticcheck
			type1.SourceContexts = nil
			type1 = type1.GetSet()
		}
		if type1.GetTypeRef() != nil {
			type1.GetTypeRef().Context.Path = nil
		}
	}
}

// EnterView_return_type is called when production view_return_type is entered.
func (s *TreeShapeListener) EnterView_return_type(ctx *parser.View_return_typeContext) {
	s.fieldname = append(s.fieldname, "view-return"+s.viewName)
	s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{}
	if ctx.QN() != nil {
		s.typemap[s.fieldname[len(s.fieldname)-1]].Opt = true
	}
}

// ExitView_return_type is called when production view_return_type is exited.
func (s *TreeShapeListener) ExitView_return_type(*parser.View_return_typeContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
	if type1.GetSet() != nil {
		type1.SourceContext = nil //nolint:staticcheck
		type1.SourceContexts = nil
		type1 = type1.GetSet()
	}
	if type1.GetTypeRef() != nil {
		type1.GetTypeRef().Context = nil
	}
}

// EnterTransform_scope_var is called when production transform_scope_var is entered.
func (s *TreeShapeListener) EnterTransform_scope_var(ctx *parser.Transform_scope_varContext) {
	tx := s.TopExpr().GetTransform()
	tx.Scopevar = ctx.GetText()
}

// ExitTransform_arg is called when production transform_arg is exited.
func (s *TreeShapeListener) ExitTransform_arg(*parser.Transform_argContext) {
	expr := s.popExpr()
	tx := s.TopExpr().GetTransform()
	tx.Arg = expr
}

// EnterTransform is called when production transform is entered.
func (s *TreeShapeListener) EnterTransform(ctx *parser.TransformContext) {
	tx := makeExpr(s.getSrcCtx(ctx.BaseParserRuleContext))
	tx.Expr = &sysl.Expr_Transform_{
		Transform: &sysl.Expr_Transform{
			Stmt: []*sysl.Expr_Transform_Stmt{},
		},
	}
	s.PushExpr(tx)
}

// ExitTransform is called when production transform is exited.
func (s *TreeShapeListener) ExitTransform(ctx *parser.TransformContext) {
	tx := s.TopExpr().GetTransform()
	if tx.Arg == nil {
		tx.Arg = makeExprName(".", s.getSrcCtx(ctx.BaseParserRuleContext))
	}
	if tx.Scopevar == "" {
		tx.Scopevar = "."
	}
}

// EnterView_param is called when production view_param is entered.
func (s *TreeShapeListener) EnterView_param(ctx *parser.View_paramContext) {
	s.fieldname = append(s.fieldname, ctx.Name_str().GetText())
	s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{}
	if ctx.QN() != nil {
		s.typemap[s.fieldname[len(s.fieldname)-1]].Opt = true
	}
}

// ExitView_param is called when production view_param is exited.
func (s *TreeShapeListener) ExitView_param(*parser.View_paramContext) {
	view := s.currentApp().Views[s.viewName]
	paramName := s.fieldname[len(s.fieldname)-1]
	type1 := s.typemap[paramName]

	if type1.GetSet() != nil && type1.GetSet().GetTypeRef() != nil {
		type1.GetSet().GetTypeRef().Context = nil
		type1.SourceContext = nil //nolint:staticcheck
		type1.SourceContexts = nil
	}
	if type1.GetTypeRef() != nil {
		type1.GetTypeRef().Context = nil
	}

	p := &sysl.Param{
		Name: paramName,
		Type: type1,
	}
	view.Param = append(view.Param, p)
}

// EnterView is called when production view is entered.
func (s *TreeShapeListener) EnterView(ctx *parser.ViewContext) {
	s.viewName = ctx.Name_str().GetText()

	// There is a bug in the antlr4 go runtime that means we cant directly get the text between 2 tokens!
	// To work around this, take the full text and the just cut it off at the expected text index
	sc := s.getSrcCtxWithText(ctx.BaseParserRuleContext)

	for index := 0; index < len(sc.Text); {
		next := strings.Index(sc.Text[index:], ":")
		if next < 1 {
			break
		}
		index += next
		if sc.Text[index-1] != '<' {
			sc.Text = sc.Text[:index+1]
			break
		}
		index++
	}

	s.currentApp().Views[s.viewName] = &sysl.View{
		Param:          []*sysl.Param{},
		RetType:        &sysl.Type{},
		SourceContext:  sc,
		SourceContexts: []*sysl.SourceContext{sc},
	}
	if ctx.Attribs_or_modifiers() != nil {
		v := s.currentApp().Views[s.viewName]
		v.Attrs = s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}

	s.fieldname = []string{}
	s.typemap = map[string]*sysl.Type{}

	v := s.currentApp().Views[s.viewName]
	if v.Attrs == nil {
		v.Attrs = map[string]*sysl.Attribute{}
	}
	s.pushScope(v)
}

// ExitView is called when production view is exited.
func (s *TreeShapeListener) ExitView(ctx *parser.ViewContext) {
	view := s.currentApp().Views[s.viewName]
	if ctx.Abstract_view() == nil {
		view.Expr = s.popExpr()
	} else {
		attributes := map[string]*sysl.Attribute{
			patternsKey: {
				Attribute: &sysl.Attribute_A{
					A: &sysl.Attribute_Array{
						Elt: []*sysl.Attribute{
							{Attribute: &sysl.Attribute_S{S: "abstract"}},
						},
					},
				},
			},
		}
		mergeAttrs(attributes, view.Attrs)
	}
	t1, has := s.typemap["view-return"+s.viewName]
	if has {
		view.RetType = t1
	} else {
		view.RetType = view.Expr.Type
	}

	s.popScope()
	s.fieldname = []string{}
	s.typemap = map[string]*sysl.Type{}
}

// EnterEnum is called when production enum is entered.
func (s *TreeShapeListener) EnterEnum(ctx *parser.EnumContext) {
	s.currentTypePath.Push(MustUnescape(ctx.Name().GetText()))
	defer s.currentTypePath.Pop()

	// Build a map of enumerations. If there are none (due to WHATEVER
	// as the only content), don't add the enum type to the sysl AST.
	items := map[string]int64{}
	for _, enumeration := range ctx.AllEnumeration() {
		e := enumeration.(*parser.EnumerationContext)
		var name string
		if e.Name() != nil {
			name = e.Name().GetText()
		} else {
			name = e.NativeDataTypes().GetText()
		}
		valstr := e.DIGITS().GetText()
		if val, err := strconv.ParseInt(valstr, 10, 64); err == nil {
			items[name] = val
		}
	}
	if len(items) == 0 {
		return
	}

	enumType := &sysl.Type{
		Type: &sysl.Type_Enum_{
			Enum: &sysl.Type_Enum{Items: items},
		},
		SourceContext:  s.getSrcCtx(ctx.BaseParserRuleContext),
		SourceContexts: []*sysl.SourceContext{s.getSrcCtx(ctx.BaseParserRuleContext)},
	}
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		enumType.Attrs = s.makeAttributeArray(attribs)
	}
	if ctx.Annotation(0) != nil {
		if enumType.Attrs == nil {
			enumType.Attrs = map[string]*sysl.Attribute{}
		}
	}
	s.pushScope(enumType)

	s.currentApp().Types[s.currentTypePath.Get()] = enumType
}

// ExitEnum is called when production enum is exited.
func (s *TreeShapeListener) ExitEnum(ctx *parser.EnumContext) {
	if s.currentApp().Types[s.currentTypePath.Get()] == nil {
		return
	}
	s.applyAnnotations(ctx.AllAnnotation())
	s.popScope()
}

// EnterAlias is called when production alias is entered.
func (s *TreeShapeListener) EnterAlias(ctx *parser.AliasContext) {
	s.currentTypePath.Push(MustUnescape(ctx.Name_str().GetText()))
	type1 := &sysl.Type{}

	s.typemap = map[string]*sysl.Type{
		s.currentTypePath.Get(): type1,
	}
	s.fieldname = []string{s.currentTypePath.Get()}
	s.currentApp().Types[s.currentTypePath.Get()] = type1

	if ctx.Attribs_or_modifiers() != nil {
		type1.Attrs = s.makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}
	if ctx.Annotation(0) != nil {
		if type1.Attrs == nil {
			type1.Attrs = map[string]*sysl.Attribute{}
		}
	}
	s.pushScope(type1)
	type1.SourceContext = s.getSrcCtx(ctx.BaseParserRuleContext)             //nolint:staticcheck
	type1.SourceContexts = append(type1.SourceContexts, type1.SourceContext) //nolint:staticcheck
}

// ExitAlias is called when production alias is exited.
func (s *TreeShapeListener) ExitAlias(ctx *parser.AliasContext) {
	s.currentApp().Types[s.currentTypePath.Get()] = s.typemap[s.fieldname[len(s.fieldname)-1]]

	s.currentTypePath.Pop()
	s.fieldname = []string{}
	s.typemap = map[string]*sysl.Type{}
	if ctx.Annotation(0) != nil {
		s.popScope()
	}
}

// EnterApp_decl is called when production app_decl is entered.
func (s *TreeShapeListener) EnterApp_decl(ctx *parser.App_declContext) {
	if s.currentApp().Types == nil && (ctx.Table(0) != nil || ctx.Alias(0) != nil || ctx.Enum(0) != nil) {
		s.currentApp().Types = map[string]*sysl.Type{}
	}
	if s.currentApp().Endpoints == nil {
		s.currentApp().Endpoints = map[string]*sysl.Endpoint{}
	}
	if s.currentApp().Wrapped == nil && len(ctx.AllFacade()) > 0 {
		s.currentApp().Wrapped = &sysl.Application{
			Types: map[string]*sysl.Type{},
		}
	}
	if s.currentApp().Views == nil && ctx.View(0) != nil {
		s.currentApp().Views = map[string]*sysl.View{}
	}
	if ctx.Annotation(0) != nil {
		if s.currentApp().Attrs == nil {
			s.currentApp().Attrs = map[string]*sysl.Attribute{}
		}
		s.pushScope(s.currentApp())
	}
	s.recordApp(s.createLocation(ctx.GetStart()))
	s.rest_queryparams = []*sysl.Endpoint_RestParams_QueryParam{}
	s.rest_queryparams_len = []int{0}
	s.rest_attrs = []map[string]*sysl.Attribute{}
	s.typemap = map[string]*sysl.Type{}
}

// ExitApp_decl is called when production app_decl is exited.
func (s *TreeShapeListener) ExitApp_decl(ctx *parser.App_declContext) {
	if ctx.Annotation(0) != nil {
		s.popScope()
	}

	s.currentApp().SourceContext.End = s.lastEnd //nolint:staticcheck
	s.currentApp().SourceContexts[len(s.currentApp().SourceContexts)-1].End = s.lastEnd
}

// EnterImport_stmt is called when production import_stmt is entered.
func (s *TreeShapeListener) EnterImport_stmt(ctx *parser.Import_stmtContext) {
	raw := strings.TrimSpace(ctx.IMPORT_PATH().GetText())
	filename := raw

	// append .sysl to filename ignoring the version
	parts := strings.Split(filename, "@")
	namePos := len(parts) - 2
	if namePos < 0 {
		namePos = 0
	}
	if filepath.Ext(parts[namePos]) == "" {
		parts[namePos] += syslExt
	}
	filename = strings.Join(parts, "@")

	if !syslutil.IsRemoteImport(filename) {
		base := s.base
		if strings.HasPrefix(filename, "/") {
			if syslutil.IsRemoteImport(s.base) {
				base = syslutil.GetRemoteRepoRoot(s.base)
			} else {
				base = "."
			}
		}

		if syslutil.IsRemoteImport(s.base) {
			filename = "/" + path.Join(base, filename)
		} else {
			filename = filepath.Join(base, filename)
		}
		if !strings.Contains(filename, "@") && s.sc.version != "" {
			filename += "@" + s.sc.version
		}
	}

	id := importDef{
		filename: filename,
	}

	if ctx.AS() != nil {
		if ctx.App_name() != nil {
			id.appname = ctx.App_name().GetText()
		} else {
			parts := make([]string, len(ctx.AllName()))
			for i, name := range ctx.AllName() {
				parts[i] = name.GetText()
			}

			id.appname = parts[len(parts)-1]
			id.pkg = strings.Join(parts[:len(parts)-1], ".")
		}
	}

	if ctx.Import_mode() != nil {
		mode := ctx.Import_mode()
		id.mode = mode.GetText()
	}

	s.imports = append(s.imports, id)

	im := &sysl.Import{
		Target:        raw,
		SourceContext: s.getSrcCtx(ctx.BaseParserRuleContext),
	}

	if len(id.appname) > 0 {
		parts := strings.Split(id.appname, "::")
		for i := range parts {
			parts[i] = MustUnescape(strings.TrimSpace(parts[i]))
		}
		im.Name = &sysl.AppName{Part: parts}
	}

	s.module.Imports = append(s.module.Imports, im)
}

func (s *TreeShapeListener) currentType() *sysl.Type {
	return s.typemap[s.fieldname[len(s.fieldname)-1]]
}

func (s *TreeShapeListener) setCurrentType(type1 *sysl.Type) {
	s.typemap[s.fieldname[len(s.fieldname)-1]] = type1
}
