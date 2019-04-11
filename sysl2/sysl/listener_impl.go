package main // SyslParser

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/grammar"
	"github.com/sirupsen/logrus"
)

var _ = fmt.Println

// TreeShapeListener ..
type TreeShapeListener struct {
	*parser.BaseSyslParserListener
	base                  string
	fs                    http.FileSystem
	filename              string
	imports               []string
	module                *sysl.Module
	appname               string
	typename              string
	fieldname             []string
	url_prefix            []string
	app_name              []string
	annotation            string
	typemap               map[string]*sysl.Type
	prevLineEmpty         bool
	pendingDocString      bool
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
}

// NewTreeShapeListener ...
func NewTreeShapeListener(fs http.FileSystem) *TreeShapeListener {
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
		fs: fs,
		module: &sysl.Module{
			Apps: map[string]*sysl.Application{},
		},
		opmap: opmap,
	}
}

// VisitErrorNode is called when an error node is visited.
func (s *TreeShapeListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *TreeShapeListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterName_str is called when production name_str is entered.
func (s *TreeShapeListener) EnterName_str(ctx *parser.Name_strContext) {
	s.app_name = append(s.app_name, ctx.GetText())
}

// ExitName_str is called when production name_str is exited.
func (s *TreeShapeListener) ExitName_str(ctx *parser.Name_strContext) {}

// EnterModifier is called when production modifier is entered.
func (s *TreeShapeListener) EnterModifier(ctx *parser.ModifierContext) {}

// ExitModifier is called when production modifier is exited.
func (s *TreeShapeListener) ExitModifier(ctx *parser.ModifierContext) {}

// EnterSize_spec is called when production size_spec is entered.
func (s *TreeShapeListener) EnterSize_spec(ctx *parser.Size_specContext) {}

// ExitSize_spec is called when production size_spec is exited.
func (s *TreeShapeListener) ExitSize_spec(ctx *parser.Size_specContext) {}

// EnterModifier_list is called when production modifier_list is entered.
func (s *TreeShapeListener) EnterModifier_list(ctx *parser.Modifier_listContext) {}

// ExitModifier_list is called when production modifier_list is exited.
func (s *TreeShapeListener) ExitModifier_list(ctx *parser.Modifier_listContext) {}

// EnterModifiers is called when production modifiers is entered.
func (s *TreeShapeListener) EnterModifiers(ctx *parser.ModifiersContext) {}

// ExitModifiers is called when production modifiers is exited.
func (s *TreeShapeListener) ExitModifiers(ctx *parser.ModifiersContext) {}

// EnterReference is called when production reference is entered.
func (s *TreeShapeListener) EnterReference(ctx *parser.ReferenceContext) {
	context_app_part := s.module.Apps[s.appname].Name.Part
	context_path := strings.Split(s.typename, ".")
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	type1.Type = &sysl.Type_TypeRef{
		TypeRef: &sysl.ScopedRef{
			Context: &sysl.Scope{
				Appname: &sysl.AppName{
					Part: context_app_part,
				},
				Path: context_path,
			},
		},
	}
	s.app_name = []string{}
}

// ExitReference is called when production reference is exited.
func (s *TreeShapeListener) ExitReference(ctx *parser.ReferenceContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
	type1.GetTypeRef().Ref = makeScope(s.app_name, ctx)
}

func lastTwoChars(str string) string {
	return str[len(str)-2:]
}

// EnterDoc_string is called when production doc_string is entered.
func (s *TreeShapeListener) EnterDoc_string(ctx *parser.Doc_stringContext) {
	if s.pendingDocString {
		s.pendingDocString = false

		space := ""
		text := ctx.TEXT().GetText()
		text = strings.Replace(text, `"`, `\"`, -1)
		text = fromQString(`"` + text[1:] + `"`)

		if s.module.Apps[s.appname].Endpoints[s.typename].GetRestParams() != nil {
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
	attrs := s.peekAttrs()

	str := ctx.TEXT().GetText()
	if len(strings.TrimSpace(str)) == 0 {
		// hack to match legacy code, required to support excel sheets
		attrs[s.annotation].Attribute.(*sysl.Attribute_S).S += "\n\n"
		s.prevLineEmpty = true
		return
	}
	ss := attrs[s.annotation].Attribute.(*sysl.Attribute_S).S
	if s.prevLineEmpty && len(ss) > 2 && lastTwoChars(ss) == "\n\n" {
		attrs[s.annotation].Attribute.(*sysl.Attribute_S).S += strings.TrimSpace(str)
		s.prevLineEmpty = false
		return
	}
	s.prevLineEmpty = false
	attrs[s.annotation].Attribute.(*sysl.Attribute_S).S += str
}

// ExitDoc_string is called when production doc_string is exited.
func (s *TreeShapeListener) ExitDoc_string(ctx *parser.Doc_stringContext) {}

// EnterQuoted_string is called when production quoted_string is entered.
func (s *TreeShapeListener) EnterQuoted_string(ctx *parser.Quoted_stringContext) {}

// ExitQuoted_string is called when production quoted_string is exited.
func (s *TreeShapeListener) ExitQuoted_string(ctx *parser.Quoted_stringContext) {}

// EnterArray_of_strings is called when production array_of_strings is entered.
func (s *TreeShapeListener) EnterArray_of_strings(ctx *parser.Array_of_stringsContext) {}

// ExitArray_of_strings is called when production array_of_strings is exited.
func (s *TreeShapeListener) ExitArray_of_strings(ctx *parser.Array_of_stringsContext) {}

// EnterArray_of_arrays is called when production array_of_arrays is entered.
func (s *TreeShapeListener) EnterArray_of_arrays(ctx *parser.Array_of_arraysContext) {}

// ExitArray_of_arrays is called when production array_of_arrays is exited.
func (s *TreeShapeListener) ExitArray_of_arrays(ctx *parser.Array_of_arraysContext) {}

// EnterNvp is called when production nvp is entered.
func (s *TreeShapeListener) EnterNvp(ctx *parser.NvpContext) {}

// ExitNvp is called when production nvp is exited.
func (s *TreeShapeListener) ExitNvp(ctx *parser.NvpContext) {}

// EnterAttributes is called when production attributes is entered.
func (s *TreeShapeListener) EnterAttributes(ctx *parser.AttributesContext) {}

// ExitAttributes is called when production attributes is exited.
func (s *TreeShapeListener) ExitAttributes(ctx *parser.AttributesContext) {}

// EnterEntry is called when production entry is entered.
func (s *TreeShapeListener) EnterEntry(ctx *parser.EntryContext) {}

// ExitEntry is called when production entry is exited.
func (s *TreeShapeListener) ExitEntry(ctx *parser.EntryContext) {}

// EnterAttribs_or_modifiers is called when production attribs_or_modifiers is entered.
func (s *TreeShapeListener) EnterAttribs_or_modifiers(ctx *parser.Attribs_or_modifiersContext) {}

// ExitAttribs_or_modifiers is called when production attribs_or_modifiers is exited.
func (s *TreeShapeListener) ExitAttribs_or_modifiers(ctx *parser.Attribs_or_modifiersContext) {}

// EnterTypes is called when production types is entered.
func (s *TreeShapeListener) EnterTypes(ctx *parser.TypesContext) {
	native := ctx.NativeDataTypes()
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	if native != nil {
		primitive_type := sysl.Type_Primitive(sysl.Type_Primitive_value[strings.ToUpper(native.GetText())])
		var constraint *sysl.Type_Constraint
		if primitive_type == sysl.Type_NO_Primitive {
			if native.GetText() == "int32" {
				primitive_type = sysl.Type_Primitive(sysl.Type_Primitive_value["INT"])
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
			} else if native.GetText() == "int64" {
				primitive_type = sysl.Type_Primitive(sysl.Type_Primitive_value["INT"])
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
		}
		type1.Type = &sysl.Type_Primitive_{
			Primitive: primitive_type,
		}
		if type1.Constraint == nil {
			type1.Constraint = []*sysl.Type_Constraint{}
		}
		if constraint != nil {
			type1.Constraint = append(type1.Constraint, constraint)
		}
	}
}

// ExitTypes is called when production types is exited.
func (s *TreeShapeListener) ExitTypes(ctx *parser.TypesContext) {}

// EnterSet_type is called when production set_type is entered.
func (s *TreeShapeListener) EnterSet_type(ctx *parser.Set_typeContext) {}

// ExitSet_type is called when production set_type is exited.
func (s *TreeShapeListener) ExitSet_type(ctx *parser.Set_typeContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{
		Type: &sysl.Type_Set{
			Set: type1,
		},
		SourceContext: type1.SourceContext,
		Opt:           type1.Opt,
		Attrs:         type1.Attrs,
	}

	type1.Opt = false

	if type1.Attrs != nil {
		type1.Attrs = nil
		s.popScope()
		s.pushScope(s.typemap[s.fieldname[len(s.fieldname)-1]])
	}

	if ctx.Size_spec() != nil {
		if type1.GetPrimitive() != sysl.Type_NO_Primitive {
			spec := ctx.Size_spec().(*parser.Size_specContext)
			type1.Constraint = makeTypeConstraint(type1.GetPrimitive(), spec)
		}
	}
}

// EnterSequence_type is called when production set_type is entered.
func (s *TreeShapeListener) EnterSequence_type(ctx *parser.Sequence_typeContext) {}

// ExitSequence_type is called when production set_type is exited.
func (s *TreeShapeListener) ExitSequence_type(ctx *parser.Sequence_typeContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{
		Type: &sysl.Type_Sequence{
			Sequence: type1,
		},
		SourceContext: type1.SourceContext,
		Opt:           type1.Opt,
		Attrs:         type1.Attrs,
	}
	type1.Opt = false

	if type1.Attrs != nil {
		type1.Attrs = nil
		s.popScope()
		s.pushScope(s.typemap[s.fieldname[len(s.fieldname)-1]])
	}

	if ctx.Size_spec() != nil {
		if type1.GetPrimitive() != sysl.Type_NO_Primitive {
			spec := ctx.Size_spec().(*parser.Size_specContext)
			type1.Constraint = makeTypeConstraint(type1.GetPrimitive(), spec)
		}
	}
}

// EnterCollection_type is called when production collection_type is entered.
func (s *TreeShapeListener) EnterCollection_type(ctx *parser.Collection_typeContext) {}

// ExitCollection_type is called when production collection_type is exited.
func (s *TreeShapeListener) ExitCollection_type(ctx *parser.Collection_typeContext) {}

// EnterUser_defined_type is called when production user_defined_type is entered.
func (s *TreeShapeListener) EnterUser_defined_type(ctx *parser.User_defined_typeContext) {
	if len(s.fieldname) == 0 {
		return
	}
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	context_app_part := s.module.Apps[s.appname].Name.Part
	context_path := strings.Split(s.typename, ".")
	ref_path := []string{ctx.GetText()}

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

// ExitUser_defined_type is called when production user_defined_type is exited.
func (s *TreeShapeListener) ExitUser_defined_type(ctx *parser.User_defined_typeContext) {}

// EnterMulti_line_docstring is called when production multi_line_docstring is entered.
func (s *TreeShapeListener) EnterMulti_line_docstring(ctx *parser.Multi_line_docstringContext) {}

// ExitMulti_line_docstring is called when production multi_line_docstring is exited.
func (s *TreeShapeListener) ExitMulti_line_docstring(ctx *parser.Multi_line_docstringContext) {}

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
	attrs := s.peekAttrs()

	if ctx.QSTRING() != nil {
		attrs[s.annotation].Attribute = &sysl.Attribute_S{
			S: fromQString(ctx.QSTRING().GetText()),
		}
	} else if ctx.Multi_line_docstring() != nil {
		attrs[s.annotation].Attribute = &sysl.Attribute_S{}
	} else {
		arr := makeArrayOfStringsAttribute(ctx.Array_of_strings().(*parser.Array_of_stringsContext))

		attrs[s.annotation].Attribute = &sysl.Attribute_A{
			A: arr.GetA(),
		}
	}
}

// ExitAnnotation_value is called when production annotation_value is exited.
func (s *TreeShapeListener) ExitAnnotation_value(ctx *parser.Annotation_valueContext) {
	if ctx.Multi_line_docstring() != nil {
		attrs := s.peekAttrs()
		attrs[s.annotation].Attribute.(*sysl.Attribute_S).S =
			strings.TrimLeft(attrs[s.annotation].GetS(), " ")
	}
}

// EnterAnnotation is called when production annotation is entered.
func (s *TreeShapeListener) EnterAnnotation(ctx *parser.AnnotationContext) {
	attr_name := ctx.VAR_NAME().GetText()
	attrs := s.peekAttrs()
	attrs[attr_name] = &sysl.Attribute{}
	s.annotation = attr_name
	if t, ok := s.peekScope().(*sysl.Type); ok {
		if t.SourceContext != nil {
			t.SourceContext.Start.Line = int32(ctx.GetStart().GetLine())
		}
	}
}

// ExitAnnotation is called when production annotation is exited.
func (s *TreeShapeListener) ExitAnnotation(ctx *parser.AnnotationContext) {
	s.annotation = ""
}

// EnterAnnotations is called when production annotations is entered.
func (s *TreeShapeListener) EnterAnnotations(ctx *parser.AnnotationsContext) {}

// ExitAnnotations is called when production annotations is exited.
func (s *TreeShapeListener) ExitAnnotations(ctx *parser.AnnotationsContext) {}

// EnterField_type is called when production field_type is entered.
func (s *TreeShapeListener) EnterField_type(ctx *parser.Field_typeContext) {
	s.app_name = []string{}
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = makeAttributeArray(attribs)
	}

	if ctx.QN() != nil {
		type1.Opt = true
	}
	type1.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}

	if ctx.Annotations() != nil {
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
			Part: app_name[:appComponentCount],
		}
	}
	scope.Path = app_name[appComponentCount:]
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
			type1.Constraint = makeTypeConstraint(primitive_type, size_spec)
		} else if has_array_spec {
			type1.Constraint = makeArrayConstraint(primitive_type, array_spec)
		}
	}
	if ctx.Annotations() != nil {
		s.popScope()
	}
}

func makeTypeConstraint(t sysl.Type_Primitive, size_spec *parser.Size_specContext) []*sysl.Type_Constraint {
	c := []*sysl.Type_Constraint{}
	var err error
	var l int

	switch t {
	case sysl.Type_DATE, sysl.Type_DATETIME, sysl.Type_INT, sysl.Type_STRING:
		val1 := size_spec.DIGITS(0).GetText()
		if l, err = strconv.Atoi(val1); err == nil {
			c = append(c, &sysl.Type_Constraint{
				Length: &sysl.Type_Constraint_Length{
					Max: int64(l),
				},
			})
		}
	case sysl.Type_DECIMAL:
		val1 := size_spec.DIGITS(0).GetText()
		if l, err = strconv.Atoi(val1); err == nil {
			c = append(c, &sysl.Type_Constraint{
				Length: &sysl.Type_Constraint_Length{
					Max: int64(l),
				},
			})
			if size_spec.DIGITS(1) != nil {
				c[0].Precision = int32(l)
				val1 = size_spec.DIGITS(1).GetText()
				if l, err = strconv.Atoi(val1); err == nil {
					c[0].Scale = int32(l)
				}
			}
		}
	default:
		panic("should not be here")
	}

	if err != nil {
		panic("should not happen: unable to parse size_spec")
	}
	return c
}

func makeArrayConstraint(t sysl.Type_Primitive, array_size *parser.Array_sizeContext) []*sysl.Type_Constraint {
	c := []*sysl.Type_Constraint{}
	var err error
	var l int

	switch t {
	case sysl.Type_DATE, sysl.Type_DATETIME, sysl.Type_INT, sysl.Type_DECIMAL, sysl.Type_STRING:
		ct := &sysl.Type_Constraint{
			Length: &sysl.Type_Constraint_Length{},
		}
		if t != sysl.Type_STRING && array_size.DIGITS(0) != nil {
			val := array_size.DIGITS(0).GetText()
			if l, err = strconv.Atoi(val); err == nil && l != 0 {
				ct.Length.Min = int64(l)
			}
		}

		if array_size.DIGITS(1) != nil {
			val1 := array_size.DIGITS(1).GetText()
			if l, err = strconv.Atoi(val1); err == nil && l != 0 {
				ct.Length.Max = int64(l)
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

func makeArrayOfStringsAttribute(array_strings *parser.Array_of_stringsContext) *sysl.Attribute {
	arr := []*sysl.Attribute{}
	for _, ars := range array_strings.AllQuoted_string() {
		str := ars.(*parser.Quoted_stringContext)

		arr = append(arr, &sysl.Attribute{
			Attribute: &sysl.Attribute_S{
				S: fromQString(str.QSTRING().GetText()),
			},
		})
	}
	return &sysl.Attribute{
		Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{
				Elt: arr,
			},
		},
	}
}

func makeAttributeArray(attribs *parser.Attribs_or_modifiersContext) map[string]*sysl.Attribute {
	patterns := []*sysl.Attribute{}
	attributes := map[string]*sysl.Attribute{}

	for _, e := range attribs.AllEntry() {
		entry := e.(*parser.EntryContext)
		if nvp, ok := entry.Nvp().(*parser.NvpContext); ok {
			if nvp.Quoted_string() != nil {
				qs := nvp.Quoted_string().(*parser.Quoted_stringContext)
				attributes[nvp.Name().GetText()] = &sysl.Attribute{
					Attribute: &sysl.Attribute_S{
						S: fromQString(qs.QSTRING().GetText()),
					},
				}
			} else if nvp.Array_of_strings() != nil {
				array_strings := nvp.Array_of_strings().(*parser.Array_of_stringsContext)
				attributes[nvp.Name().GetText()] = makeArrayOfStringsAttribute(array_strings)
			} else if nvp.Array_of_arrays() != nil {
				arr := nvp.Array_of_arrays().(*parser.Array_of_arraysContext)
				attrArray := sysl.Attribute_Array{
					Elt: []*sysl.Attribute{},
				}
				for _, astrings := range arr.AllArray_of_strings() {
					array_strings := astrings.(*parser.Array_of_stringsContext)
					elt := makeArrayOfStringsAttribute(array_strings)
					attrArray.Elt = append(attrArray.Elt, elt)
				}

				attributes[nvp.Name().GetText()] = &sysl.Attribute{
					Attribute: &sysl.Attribute_A{
						A: &attrArray,
					},
				}
			}
		} else if mod, ok := entry.Modifier().(*parser.ModifierContext); ok {
			patterns = append(patterns, &sysl.Attribute{
				Attribute: &sysl.Attribute_S{
					S: mod.GetText()[1:],
				},
			})
		}
	}
	if len(patterns) > 0 {
		attributes["patterns"] = &sysl.Attribute{
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
func (s *TreeShapeListener) EnterInplace_tuple(ctx *parser.Inplace_tupleContext) {
	s.typename = s.typename + "." + s.fieldname[len(s.fieldname)-1]
	s.typemap = map[string]*sysl.Type{}
	s.module.Apps[s.appname].Types[s.typename] = &sysl.Type{
		Type: &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{
				AttrDefs: s.typemap,
			},
		},
	}
}

// ExitInplace_tuple is called when production inplace_tuple is exited.
func (s *TreeShapeListener) ExitInplace_tuple(ctx *parser.Inplace_tupleContext) {
	fixFieldDefinitions(s.module.Apps[s.appname].Types[s.typename])
	l := strings.LastIndex(s.typename, ".")
	s.typename = s.typename[:l]
	s.typemap = s.module.Apps[s.appname].Types[s.typename].GetTuple().GetAttrDefs()
}

// EnterField is called when production field is entered.
func (s *TreeShapeListener) EnterField(ctx *parser.FieldContext) {
	fieldName := ctx.Name_str().GetText()
	s.fieldname = append(s.fieldname, fieldName)
	type1, has := s.typemap[fieldName]
	if has {
		logrus.Warnf("%s) %s.%s defined multiple times",
			s.filename, s.typename, fieldName)
	} else {
		type1 = &sysl.Type{}
	}

	type1.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}
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
	s.app_name = []string{}
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

// EnterInplace_table is called when production inplace_table is entered.
func (s *TreeShapeListener) EnterInplace_table(ctx *parser.Inplace_tableContext) {}

// ExitInplace_table is called when production inplace_table is exited.
func (s *TreeShapeListener) ExitInplace_table(ctx *parser.Inplace_tableContext) {}

// EnterTable_stmts is called when production table_stmts is entered.
func (s *TreeShapeListener) EnterTable_stmts(ctx *parser.Table_stmtsContext) {
	type1 := s.module.Apps[s.appname].Types[s.typename]
	if ctx.Annotation(0) != nil {
		if type1.Attrs == nil {
			type1.Attrs = map[string]*sysl.Attribute{}
		}
	}
	if ctx.Field(0) == nil {
		type1.Type = nil
	}
}

// ExitTable_stmts is called when production table_stmts is exited.
func (s *TreeShapeListener) ExitTable_stmts(ctx *parser.Table_stmtsContext) {}

// EnterTable_def is called when production table_def is entered.
func (s *TreeShapeListener) EnterTable_def(ctx *parser.Table_defContext) {
	type1 := s.module.Apps[s.appname].Types[s.typename]
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = makeAttributeArray(attribs)
	}
	if ctx.WHATEVER() != nil {
		type1.Type = nil
	}
}

// ExitTable_def is called when production table_def is exited.
func (s *TreeShapeListener) ExitTable_def(ctx *parser.Table_defContext) {}

// EnterTable is called when production table is entered.
func (s *TreeShapeListener) EnterTable(ctx *parser.TableContext) {
	if s.typename == "" {
		s.typename = ctx.Name_str().GetText()
	} else {
		s.typename = s.typename + "." + ctx.Name_str().GetText()
	}
	s.typemap = map[string]*sysl.Type{}

	types := s.module.Apps[s.appname].Types
	if ctx.TABLE() != nil {
		if types[s.typename].GetRelation().GetAttrDefs() != nil {
			panic("not implemented yet")
		}

		types[s.typename] = &sysl.Type{
			Type: &sysl.Type_Relation_{
				Relation: &sysl.Type_Relation{
					AttrDefs: s.typemap,
				},
			},
		}
	}
	if ctx.TYPE() != nil {
		types[s.typename] = &sysl.Type{
			Type: &sysl.Type_Tuple_{
				Tuple: &sysl.Type_Tuple{
					AttrDefs: s.typemap,
				},
			},
		}
	}
	type1 := types[s.typename]
	s.pushScope(type1)
	type1.SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
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

func fixFieldDefinitions(collection *sysl.Type) {
	for name, f := range attributesForType(collection) {
		if f.Type == nil {
			continue
		}
		if f.GetPrimitive() == sysl.Type_NO_Primitive {
			var type1 *sysl.ScopedRef
			switch t := f.GetType().(type) {
			case *sysl.Type_TypeRef:
				type1 = t.TypeRef
			case *sysl.Type_Sequence:
				type1 = t.Sequence.GetTypeRef()
			case *sysl.Type_Set:
				type1 = t.Set.GetTypeRef()
			case *sysl.Type_List_:
				type1 = t.List.GetType().GetTypeRef()
			default:
				panic("unhandled type:" + name)
			}

			if type1 != nil && type1.Ref != nil && type1.Ref.Appname != nil {
				l := len(type1.Ref.Appname.Part)
				str := []string{strings.TrimSpace(type1.Ref.Appname.Part[l-1])}
				type1.Ref.Path = append(str, type1.Ref.Path...)
				type1.Ref.Appname.Part = type1.Ref.Appname.Part[:l-1]
				if len(type1.Ref.Appname.Part) == 0 {
					type1.Ref.Appname = nil
				}
			}
		}
	}
}

// ExitTable is called when production table is exited.
func (s *TreeShapeListener) ExitTable(ctx *parser.TableContext) {
	// wire up primary key
	if rel := s.module.Apps[s.appname].Types[s.typename].GetRelation(); rel != nil {
		pks := []string{}
		for _, name := range s.fieldname {
			f := rel.GetAttrDefs()[name]
			if patterns, has := f.GetAttrs()["patterns"]; has {
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
	tableDef := ctx.Table_def().(*parser.Table_defContext)
	if tableStmts := tableDef.Table_stmts(); tableStmts != nil {
		s.applyAnnotations(tableStmts.(*parser.Table_stmtsContext).AllAnnotation())
	}
	s.popScope()

	// Match legacy behavior
	fixFieldDefinitions(s.module.Apps[s.appname].Types[s.typename])
	// End

	l := strings.LastIndex(s.typename, ".")
	if l > 0 {
		s.typename = s.typename[:l]
	} else {
		s.typename = ""
	}

	s.fieldname = []string{}
	s.typemap = nil
}

func (s *TreeShapeListener) applyAnnotations(
	annotations []parser.IAnnotationContext,
) {
	// Match legacy behavior
	// Copy the annotations from the parent (tuple or relation) to each child
	collection := s.module.Apps[s.appname].Types[s.typename]

	for _, annotation := range annotations {
		varname := annotation.(*parser.AnnotationContext).VAR_NAME().GetText()
		attr := collection.Attrs[varname]
		for _, name := range s.fieldname {
			f := attributesForType(collection)[name]

			if f.Attrs == nil {
				f.Attrs = map[string]*sysl.Attribute{}
			}
			f.Attrs[varname] = attr
		}
	}
}

func (s *TreeShapeListener) pushTypename(typename string) string {
	if s.typename != "" {
		s.typename += "."
	}
	s.typename += typename
	return s.typename
}

func (s *TreeShapeListener) popTypename() string {
	if lastDot := strings.LastIndex(s.typename, "."); lastDot != -1 {
		s.typename = s.typename[:lastDot]
	} else {
		s.typename = ""
	}
	return s.typename
}

// EnterUnion is called when production union is entered.
func (s *TreeShapeListener) EnterUnion(ctx *parser.UnionContext) {
	s.pushTypename(ctx.Name_str().GetText())
	s.typemap = map[string]*sysl.Type{}

	types := s.module.Apps[s.appname].Types

	if types[s.typename].GetOneOf().GetType() != nil {
		panic("not implemented yet")
	}

	types[s.typename] = &sysl.Type{
		Type: &sysl.Type_OneOf_{
			OneOf: &sysl.Type_OneOf{},
		},
	}

	type1 := types[s.typename]
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = makeAttributeArray(attribs)
	}
	if ctx.Annotation(0) != nil {
		if type1.Attrs == nil {
			type1.Attrs = map[string]*sysl.Attribute{}
		}
	}
	s.pushScope(type1)
	type1.SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
}

// ExitUnion is called when production union is exited.
func (s *TreeShapeListener) ExitUnion(ctx *parser.UnionContext) {
	s.applyAnnotations(ctx.AllAnnotation())
	s.popScope()

	context_app_part := s.module.Apps[s.appname].Name.Part
	context_path := strings.Split(s.typename, ".")

	oneof := s.module.Apps[s.appname].Types[s.typename].GetOneOf()
	for _, ref := range ctx.AllUser_defined_type() {
		oneof.Type = append(oneof.Type, &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: context_app_part,
						},
						Path: context_path,
					},
					Ref: &sysl.Scope{
						Path: []string{ref.(*parser.User_defined_typeContext).Name_str().GetText()},
					},
				},
			},
		})
	}

	s.popTypename()

	s.fieldname = []string{}
	s.typemap = nil
}

// EnterPackage_name is called when production package_name is entered.
func (s *TreeShapeListener) EnterPackage_name(ctx *parser.Package_nameContext) {}

// ExitPackage_name is called when production package_name is exited.
func (s *TreeShapeListener) ExitPackage_name(ctx *parser.Package_nameContext) {}

// EnterSub_package is called when production sub_package is entered.
func (s *TreeShapeListener) EnterSub_package(ctx *parser.Sub_packageContext) {
	top := len(s.app_name) - 1
	str := ctx.NAME_SEP().GetText()
	sp := strings.Split(str, "::")
	s.app_name[top] = s.app_name[top] + sp[0]
}

// ExitSub_package is called when production sub_package is exited.
func (s *TreeShapeListener) ExitSub_package(ctx *parser.Sub_packageContext) {
	top := len(s.app_name) - 1
	str := ctx.NAME_SEP().GetText()
	sp := strings.Split(str, "::")
	s.app_name[top] = sp[1] + s.app_name[top]
}

// EnterApp_name is called when production app_name is entered.
func (s *TreeShapeListener) EnterApp_name(ctx *parser.App_nameContext) {
	s.app_name = []string{}
}

// ExitApp_name is called when production app_name is exited.
func (s *TreeShapeListener) ExitApp_name(ctx *parser.App_nameContext) {}

// EnterName_with_attribs is called when production name_with_attribs is entered.
func (s *TreeShapeListener) EnterName_with_attribs(ctx *parser.Name_with_attribsContext) {
	s.appname = ctx.App_name().GetText()
	if _, has := s.module.Apps[s.appname]; !has {
		s.module.Apps[s.appname] = &sysl.Application{
			Name: &sysl.AppName{},
		}
	}

	if ctx.QSTRING() != nil {
		s.module.Apps[s.appname].LongName = fromQString(ctx.QSTRING().GetText())
	}

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		attrs := makeAttributeArray(attribs)
		if s.module.Apps[s.appname].Attrs == nil {
			s.module.Apps[s.appname].Attrs = attrs
		} else {
			mergeAttrs(attrs, s.module.Apps[s.appname].Attrs)
		}
	}
	s.module.Apps[s.appname].SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()+1)
}

// ExitName_with_attribs is called when production name_with_attribs is exited.
func (s *TreeShapeListener) ExitName_with_attribs(ctx *parser.Name_with_attribsContext) {
	for i := range s.app_name {
		s.app_name[i] = strings.TrimSpace(s.app_name[i])
	}
	s.module.Apps[s.appname].GetName().Part = s.app_name
}

// EnterModel_name is called when production model_name is entered.
func (s *TreeShapeListener) EnterModel_name(ctx *parser.Model_nameContext) {
	if s.module.Apps[s.appname].Wrapped.Name != nil {
		panic("not implemented yet?")
	}

	s.module.Apps[s.appname].Wrapped.Name = &sysl.AppName{
		Part: []string{ctx.Name().GetText()},
	}
}

// ExitModel_name is called when production model_name is exited.
func (s *TreeShapeListener) ExitModel_name(ctx *parser.Model_nameContext) {}

// EnterInplace_table_def is called when production inplace_table_def is entered.
func (s *TreeShapeListener) EnterInplace_table_def(ctx *parser.Inplace_table_defContext) {}

// ExitInplace_table_def is called when production inplace_table_def is exited.
func (s *TreeShapeListener) ExitInplace_table_def(ctx *parser.Inplace_table_defContext) {}

// EnterTable_refs is called when production table_refs is entered.
func (s *TreeShapeListener) EnterTable_refs(ctx *parser.Table_refsContext) {
	s.module.Apps[s.appname].Wrapped.Types[ctx.Name().GetText()] = &sysl.Type{}
}

// ExitTable_refs is called when production table_refs is exited.
func (s *TreeShapeListener) ExitTable_refs(ctx *parser.Table_refsContext) {}

// EnterFacade is called when production facade is entered.
func (s *TreeShapeListener) EnterFacade(ctx *parser.FacadeContext) {}

// ExitFacade is called when production facade is exited.
func (s *TreeShapeListener) ExitFacade(ctx *parser.FacadeContext) {}

// EnterDocumentation_stmts is called when production documentation_stmts is entered.
func (s *TreeShapeListener) EnterDocumentation_stmts(ctx *parser.Documentation_stmtsContext) {}

// ExitDocumentation_stmts is called when production documentation_stmts is exited.
func (s *TreeShapeListener) ExitDocumentation_stmts(ctx *parser.Documentation_stmtsContext) {}

// EnterQuery_var is called when production query_var is entered.
func (s *TreeShapeListener) EnterQuery_var(ctx *parser.Query_varContext) {
	var_name := ctx.Name().GetText()
	var type1 *sysl.Type
	var ref_path []string

	if ctx.Var_in_curly() != nil {
		ref_path = append(ref_path, ctx.Var_in_curly().GetText())
		type1 = &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: s.module.Apps[s.appname].Name.Part,
						},
					},
					Ref: &sysl.Scope{
						Path: ref_path,
					},
				},
			},
		}
	} else if ctx.NativeDataTypes() != nil {
		type_str := strings.ToUpper(ctx.NativeDataTypes().GetText())
		primitive_type := sysl.Type_Primitive(sysl.Type_Primitive_value[type_str])
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: primitive_type,
			},
		}
	} else if ctx.Name_str() != nil {
		type1 = &sysl.Type{}
	}

	rest_param := &sysl.Endpoint_RestParams_QueryParam{
		Name: var_name,
		Type: type1,
	}

	if ctx.QN() != nil {
		rest_param.Type.Opt = true
	}

	rest_param.Type.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}
	s.method_urlparams = append(s.method_urlparams, rest_param)
}

// ExitQuery_var is called when production query_var is exited.
func (s *TreeShapeListener) ExitQuery_var(ctx *parser.Query_varContext) {}

// EnterQuery_param is called when production query_param is entered.
func (s *TreeShapeListener) EnterQuery_param(ctx *parser.Query_paramContext) {}

// ExitQuery_param is called when production query_param is exited.
func (s *TreeShapeListener) ExitQuery_param(ctx *parser.Query_paramContext) {}

// EnterHttp_path_var_with_type is called when production http_path_var_with_type is entered.
func (s *TreeShapeListener) EnterHttp_path_var_with_type(ctx *parser.Http_path_var_with_typeContext) {
	var_name := ctx.Http_path_part().GetText()
	var type1 *sysl.Type
	if ctx.NativeDataTypes() != nil {
		type_str := strings.ToUpper(ctx.NativeDataTypes().GetText())
		primitive_type := sysl.Type_Primitive(sysl.Type_Primitive_value[type_str])
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: primitive_type,
			},
		}
	} else if ctx.Reference() != nil {
		s.fieldname = append(s.fieldname, var_name)
		type1 = &sysl.Type{}
		s.typemap[s.fieldname[len(s.fieldname)-1]] = type1
	} else {
		ref_path := []string{ctx.Name_str().GetText()}

		type1 = &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: s.module.Apps[s.appname].Name.Part,
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
		Name: var_name,
		Type: type1,
	}

	rest_param.Type.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}

	s.rest_urlparams = append(s.rest_urlparams, rest_param)
	s.typename += ctx.CURLY_OPEN().GetText() + var_name + ctx.CURLY_CLOSE().GetText()
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
	s.typename += ctx.GetText()
}

// ExitHttp_path_static is called when production http_path_static is exited.
func (s *TreeShapeListener) ExitHttp_path_static(ctx *parser.Http_path_staticContext) {}

// EnterHttp_path_suffix is called when production http_path_suffix is entered.
func (s *TreeShapeListener) EnterHttp_path_suffix(ctx *parser.Http_path_suffixContext) {
	s.typename += ctx.FORWARD_SLASH().GetText()
}

// ExitHttp_path_suffix is called when production http_path_suffix is exited.
func (s *TreeShapeListener) ExitHttp_path_suffix(ctx *parser.Http_path_suffixContext) {}

// EnterHttp_path is called when production http_path is entered.
func (s *TreeShapeListener) EnterHttp_path(ctx *parser.Http_pathContext) {
	s.typename = ""
	if ctx.FORWARD_SLASH() != nil {
		s.typename = ctx.GetText()
	}
}

// ExitHttp_path is called when production http_path is exited.
func (s *TreeShapeListener) ExitHttp_path(ctx *parser.Http_pathContext) {
	// s.typename is built along as we enter http_path/http_path_suffix/http_path_var_with_type
	// commit this value to url_prefix
	s.url_prefix = append(s.url_prefix, s.typename)
}

// EnterEndpoint_name is called when production endpoint_name is entered.
func (s *TreeShapeListener) EnterEndpoint_name(ctx *parser.Endpoint_nameContext) {}

// ExitEndpoint_name is called when production endpoint_name is exited.
func (s *TreeShapeListener) ExitEndpoint_name(ctx *parser.Endpoint_nameContext) {}

// EnterRet_stmt is called when production ret_stmt is entered.
func (s *TreeShapeListener) EnterRet_stmt(ctx *parser.Ret_stmtContext) {
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Ret{
			Ret: &sysl.Return{
				Payload: strings.Trim(ctx.TEXT().GetText(), " "),
			},
		},
	})
}

// ExitRet_stmt is called when production ret_stmt is exited.
func (s *TreeShapeListener) ExitRet_stmt(ctx *parser.Ret_stmtContext) {}

// EnterTarget is called when production target is entered.
func (s *TreeShapeListener) EnterTarget(ctx *parser.TargetContext) {
	s.app_name = []string{s.appname}
}

// ExitTarget is called when production target is exited.
func (s *TreeShapeListener) ExitTarget(ctx *parser.TargetContext) {
	for i := range s.app_name {
		s.app_name[i] = strings.TrimSpace(s.app_name[i])
	}
	s.lastStatement().GetCall().Target.Part = s.app_name
	s.app_name = []string{}
}

// EnterTarget_endpoint is called when production target_endpoint is entered.
func (s *TreeShapeListener) EnterTarget_endpoint(ctx *parser.Target_endpointContext) {}

// ExitTarget_endpoint is called when production target_endpoint is exited.
func (s *TreeShapeListener) ExitTarget_endpoint(ctx *parser.Target_endpointContext) {}

// EnterCall_arg is called when production call_arg is entered.
func (s *TreeShapeListener) EnterCall_arg(ctx *parser.Call_argContext) {
	arg := &sysl.Call_Arg{
		Name: ctx.GetText(),
	}
	s.lastStatement().GetCall().Arg = append(s.lastStatement().GetCall().Arg, arg)
}

// ExitCall_arg is called when production call_arg is exited.
func (s *TreeShapeListener) ExitCall_arg(ctx *parser.Call_argContext) {}

// EnterCall_args is called when production call_args is entered.
func (s *TreeShapeListener) EnterCall_args(ctx *parser.Call_argsContext) {}

// ExitCall_args is called when production call_args is exited.
func (s *TreeShapeListener) ExitCall_args(ctx *parser.Call_argsContext) {}

// EnterCall_stmt is called when production call_stmt is entered.
func (s *TreeShapeListener) EnterCall_stmt(ctx *parser.Call_stmtContext) {
	appName := &sysl.AppName{}
	if ctx.DOT_ARROW() != nil {
		appName.Part = s.module.Apps[s.appname].Name.Part
	}
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: ctx.Target_endpoint().GetText(),
			},
		},
	})
	if ctx.Call_args() != nil {
		s.lastStatement().GetCall().Arg = []*sysl.Call_Arg{}
	}
}

// ExitCall_stmt is called when production call_stmt is exited.
func (s *TreeShapeListener) ExitCall_stmt(ctx *parser.Call_stmtContext) {}

// EnterIf_stmt is called when production if_stmt is entered.
func (s *TreeShapeListener) EnterIf_stmt(ctx *parser.If_stmtContext) {
	if_stmt := &sysl.Statement{
		Stmt: &sysl.Statement_Cond{
			Cond: &sysl.Cond{
				Test: ctx.PREDICATE_VALUE().GetText(),
				Stmt: []*sysl.Statement{},
			},
		},
	}
	s.addToCurrentScope(if_stmt)
	s.pushScope(if_stmt.GetCond())
}

// ExitIf_stmt is called when production if_stmt is exited.
func (s *TreeShapeListener) ExitIf_stmt(ctx *parser.If_stmtContext) {
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
	}
	s.addToCurrentScope(else_stmt)
	s.pushScope(else_stmt.GetGroup())
}

// ExitElse_stmt is called when production else_stmt is exited.
func (s *TreeShapeListener) ExitElse_stmt(ctx *parser.Else_stmtContext) {
	s.popScope()
}

// EnterIf_else is called when production if_else is entered.
func (s *TreeShapeListener) EnterIf_else(ctx *parser.If_elseContext) {
}

// ExitIf_else is called when production if_else is exited.
func (s *TreeShapeListener) ExitIf_else(ctx *parser.If_elseContext) {
}

// EnterFor_stmt is called when production for_stmt is entered.
func (s *TreeShapeListener) EnterFor_stmt(ctx *parser.For_stmtContext) {
	stmt := &sysl.Statement{}
	s.addToCurrentScope(stmt)

	if ctx.FOR() != nil || ctx.LOOP() != nil {
		var text string
		if ctx.FOR() != nil {
			text = ctx.FOR().GetText()
		} else {
			text = ctx.LOOP().GetText()
		}
		text = text + ctx.PREDICATE_VALUE().GetText()
		text = strings.TrimSpace(text)
		stmt.Stmt = &sysl.Statement_Group{
			Group: &sysl.Group{
				Title: text,
				Stmt:  []*sysl.Statement{},
			},
		}
		s.pushScope(stmt.GetGroup())
	} else if ctx.UNTIL() != nil || ctx.WHILE() != nil {
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
	} else if ctx.FOR_EACH() != nil {
		stmt.Stmt = &sysl.Statement_Foreach{
			Foreach: &sysl.Foreach{
				Collection: ctx.PREDICATE_VALUE().GetText(),
				Stmt:       []*sysl.Statement{},
			},
		}
		s.pushScope(stmt.GetForeach())
	} else if ctx.ALT() != nil {
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
func (s *TreeShapeListener) ExitFor_stmt(ctx *parser.For_stmtContext) {
	s.popScope()
}

// EnterHttp_method_comment is called when production http_method_comment is entered.
func (s *TreeShapeListener) EnterHttp_method_comment(ctx *parser.Http_method_commentContext) {}

// ExitHttp_method_comment is called when production http_method_comment is exited.
func (s *TreeShapeListener) ExitHttp_method_comment(ctx *parser.Http_method_commentContext) {}

// EnterGroup_stmt is called when production group_stmt is entered.
func (s *TreeShapeListener) EnterGroup_stmt(ctx *parser.Group_stmtContext) {
	stmt := &sysl.Statement{}

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
func (s *TreeShapeListener) ExitGroup_stmt(ctx *parser.Group_stmtContext) {
	s.popScope()
}

// EnterOne_of_case_label is called when production one_of_case_label is entered.
func (s *TreeShapeListener) EnterOne_of_case_label(ctx *parser.One_of_case_labelContext) {}

// ExitOne_of_case_label is called when production one_of_case_label is exited.
func (s *TreeShapeListener) ExitOne_of_case_label(ctx *parser.One_of_case_labelContext) {}

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
func (s *TreeShapeListener) ExitOne_of_cases(ctx *parser.One_of_casesContext) {
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
		Stmt: alt,
	})
	s.pushScope(alt.Alt)
}

// ExitOne_of_stmt is called when production one_of_stmt is exited.
func (s *TreeShapeListener) ExitOne_of_stmt(ctx *parser.One_of_stmtContext) {
	s.popScope()
}

func withQuotesQString(str string) string {
	l := len(str)
	if l > 0 && str[:1] == "'" && str[l-1:] == "'" {
		return str
	}
	l = len(str)
	if l > 0 && str[:1] == `"` && str[l-1:] == `"` {
		return str
	}
	return `"` + str + `"`
}

// EnterText_stmt is called when production text_stmt is entered.
func (s *TreeShapeListener) EnterText_stmt(ctx *parser.Text_stmtContext) {
	// Need to Coalesce multiple doc_string's into one
	// See enterdoc_string.
	if ctx.Doc_string() == nil {
		str := ctx.GetText()
		if ctx.QSTRING() != nil {
			str = withQuotesQString(str)
		}
		s.addToCurrentScope(&sysl.Statement{
			Stmt: &sysl.Statement_Action{
				Action: &sysl.Action{
					Action: str,
				},
			},
		})
		s.pendingDocString = false
	} else {
		s.pendingDocString = true

		if s.module.Apps[s.appname].Endpoints[s.typename].GetRestParams() != nil {
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
			s.addToCurrentScope(&sysl.Statement{
				Stmt: &sysl.Statement_Action{
					Action: &sysl.Action{
						Action: "|",
					},
				},
			})
		}

	}
}

// ExitText_stmt is called when production text_stmt is exited.
func (s *TreeShapeListener) ExitText_stmt(ctx *parser.Text_stmtContext) {}

// EnterMixin is called when production mixin is entered.
func (s *TreeShapeListener) EnterMixin(ctx *parser.MixinContext) {
	if s.module.Apps[s.appname].Mixin2 == nil {
		s.module.Apps[s.appname].Mixin2 = []*sysl.Application{}
	}
}

// ExitMixin is called when production mixin is exited.
func (s *TreeShapeListener) ExitMixin(ctx *parser.MixinContext) {
	s.module.Apps[s.appname].Mixin2 = append(s.module.Apps[s.appname].Mixin2, &sysl.Application{
		Name: &sysl.AppName{
			Part: s.app_name,
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
	s.app_name = []string{}
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

// EnterParam_list is called when production param_list is entered.
func (s *TreeShapeListener) EnterParam_list(ctx *parser.Param_listContext) {}

// ExitParam_list is called when production param_list is exited.
func (s *TreeShapeListener) ExitParam_list(ctx *parser.Param_listContext) {}

// EnterParams is called when production params is entered.
func (s *TreeShapeListener) EnterParams(ctx *parser.ParamsContext) {
	s.typemap = map[string]*sysl.Type{}
	s.fieldname = []string{}
}

// ExitParams is called when production params is exited.
func (s *TreeShapeListener) ExitParams(ctx *parser.ParamsContext) {
	params := []*sysl.Param{}

	for _, fieldname := range s.fieldname {
		type1 := s.typemap[fieldname]
		if type1.GetSet() != nil {
			type1.GetSet().GetTypeRef().Context = nil
			type1.GetSet().SourceContext = nil
			ref := type1.GetSet().GetTypeRef().GetRef()
			if ref.Appname == nil {
				ref.Appname = &sysl.AppName{
					Part: ref.Path,
				}
				ref.Path = nil
			}
		} else if type1.GetTypeRef() != nil {
			type1.GetTypeRef().Context = nil
			ref := type1.GetTypeRef().GetRef()
			if ref.Appname == nil {
				ref.Appname = &sysl.AppName{
					Part: ref.Path,
				}
				ref.Path = nil
			}
			for i := range ref.Appname.Part {
				ref.Appname.Part[i] = strings.TrimSpace(ref.Appname.Part[i])
			}
		} else if type1.Type == nil {
			type1.Type = &sysl.Type_NoType_{
				NoType: &sysl.Type_NoType{},
			}
		}
		type1.SourceContext = nil

		p := sysl.Param{
			Name: fieldname,
			Type: type1,
		}
		params = append(params, &p)
	}
	if s.module.Apps[s.appname].Endpoints[s.typename].Param == nil {
		s.module.Apps[s.appname].Endpoints[s.typename].Param = params
	} else {
		s.module.Apps[s.appname].Endpoints[s.typename].Param = append(s.module.Apps[s.appname].Endpoints[s.typename].Param, params...)
	}
	s.typemap = nil
	s.fieldname = []string{}
}

// EnterStatements is called when production statements is entered.
func (s *TreeShapeListener) EnterStatements(ctx *parser.StatementsContext) {}

// ExitStatements is called when production statements is exited.
func (s *TreeShapeListener) ExitStatements(ctx *parser.StatementsContext) {
	if ctx.Attribs_or_modifiers() != nil {
		stmt := s.lastStatement()
		stmt.Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}
}

func (s *TreeShapeListener) urlString() string {
	var url string
	for _, str := range s.url_prefix {
		url += str
	}
	return url
}

func (s *TreeShapeListener) pushExpr(expr *sysl.Expr) {
	s.expr_stack = append(s.expr_stack, expr)
}

func (s *TreeShapeListener) popExpr() *sysl.Expr {
	l := len(s.expr_stack) - 1
	ret := s.expr_stack[l]
	s.expr_stack = s.expr_stack[:l]
	return ret
}

func (s *TreeShapeListener) topExpr() *sysl.Expr {
	l := len(s.expr_stack) - 1
	ret := s.expr_stack[l]
	return ret
}

func (s *TreeShapeListener) pushScope(scope interface{}) {
	s.stmt_scope = append(s.stmt_scope, scope)
}

func (s *TreeShapeListener) popScope() {
	l := len(s.stmt_scope) - 1
	s.stmt_scope = s.stmt_scope[:l]
}

func (s *TreeShapeListener) peekScope() interface{} {
	l := len(s.stmt_scope) - 1
	return s.stmt_scope[l]
}

func (s *TreeShapeListener) lastStatement() *sysl.Statement {
	top := len(s.stmt_scope) - 1
	switch scope := s.stmt_scope[top].(type) {
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
			if dst[k].GetA() != nil && v.GetA() != nil {
				dst[k].GetA().Elt = append(dst[k].GetA().Elt, v.GetA().Elt...)
			} else {
				dst[k] = v
			}
		}
	}
}

// EnterMethod_def is called when production method_def is entered.
func (s *TreeShapeListener) EnterMethod_def(ctx *parser.Method_defContext) {
	url := s.urlString()
	method := strings.TrimSpace(ctx.HTTP_VERBS().GetText())
	s.typename = method + " " + url
	s.method_urlparams = []*sysl.Endpoint_RestParams_QueryParam{}
	if _, has := s.module.Apps[s.appname].Endpoints[s.typename]; !has {
		s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
			Name: s.typename,
			RestParams: &sysl.Endpoint_RestParams{
				Method: sysl.Endpoint_RestParams_Method(sysl.Endpoint_RestParams_Method_value[method]),
				Path:   url,
			},
			Stmt: []*sysl.Statement{},
		}
	}
	restEndpoint := s.module.Apps[s.appname].Endpoints[s.typename]
	s.pushScope(restEndpoint)

	attrs := map[string]*sysl.Attribute{}
	elt := []*sysl.Attribute{&sysl.Attribute{
		Attribute: &sysl.Attribute_S{
			S: "rest",
		},
	}}

	attrs["patterns"] = &sysl.Attribute{
		Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{
				Elt: elt,
			},
		},
	}
	for _, parentAttrs := range s.rest_attrs {
		mergeAttrs(parentAttrs, attrs)
	}

	if ctx.Attribs_or_modifiers() != nil {
		mergeAttrs(makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext)), attrs)
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
				Name: q.Name,
				Type: &sysl.Type{
					Type: q.Type.Type,
					SourceContext: &sysl.SourceContext{
						Start: &sysl.SourceContext_Location{
							Line: int32(ctx.GetStart().GetLine()),
						},
					},
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
	restEndpoint.SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
}

// ExitMethod_def is called when production method_def is exited.
func (s *TreeShapeListener) ExitMethod_def(ctx *parser.Method_defContext) {
	if len(s.method_urlparams) > 0 {
		qparams := s.module.Apps[s.appname].Endpoints[s.typename].RestParams.QueryParam
		if qparams == nil {
			qparams = []*sysl.Endpoint_RestParams_QueryParam{}
		}
		qparams = append(qparams, s.method_urlparams...)
		s.module.Apps[s.appname].Endpoints[s.typename].RestParams.QueryParam = qparams
	}
	if s.module.Apps[s.appname].Endpoints[s.typename].Param != nil {
		for i, p := range s.module.Apps[s.appname].Endpoints[s.typename].Param {
			if p.GetType().GetNoType() != nil {
				s.module.Apps[s.appname].Endpoints[s.typename].Param[i] = &sysl.Param{}
			}
		}
	}

	if len(s.module.Apps[s.appname].Endpoints[s.typename].Stmt) == 0 {
		s.module.Apps[s.appname].Endpoints[s.typename].Stmt = nil
	}

	s.popScope()
	s.typename = ""
	s.http_path_query_param = ""
}

// EnterShortcut is called when production shortcut is entered.
func (s *TreeShapeListener) EnterShortcut(ctx *parser.ShortcutContext) {}

// ExitShortcut is called when production shortcut is exited.
func (s *TreeShapeListener) ExitShortcut(ctx *parser.ShortcutContext) {}

// EnterSimple_endpoint is called when production api_endpoint is entered.
func (s *TreeShapeListener) EnterSimple_endpoint(ctx *parser.Simple_endpointContext) {
	if ctx.WHATEVER() != nil {
		s.module.Apps[s.appname].Endpoints[ctx.WHATEVER().GetText()] = &sysl.Endpoint{
			Name: ctx.WHATEVER().GetText(),
		}
		return
	}
	s.typename = ctx.Endpoint_name().GetText()
	ep := s.module.Apps[s.appname].Endpoints[s.typename]

	if ep == nil {
		ep = &sysl.Endpoint{
			Name: s.typename,
		}
		s.module.Apps[s.appname].Endpoints[s.typename] = ep
	}

	if ctx.QSTRING() != nil {
		ep.LongName = fromQString(ctx.QSTRING().GetText())
	}

	if ctx.Attribs_or_modifiers() != nil {
		attrs := makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
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

		s.pushScope(s.module.Apps[s.appname].Endpoints[s.typename])
	}
	ep.SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
}

// ExitSimple_endpoint is called when production simple_endpoint is exited.
func (s *TreeShapeListener) ExitSimple_endpoint(ctx *parser.Simple_endpointContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
	ep := s.module.Apps[s.appname].Endpoints[s.typename]
	if ep != nil && len(ep.Attrs) == 0 {
		ep.Attrs = nil
	}
	s.typename = ""
}

// EnterRest_endpoint is called when production rest_endpoint is entered.
func (s *TreeShapeListener) EnterRest_endpoint(ctx *parser.Rest_endpointContext) {
	s.rest_queryparams_len = append(s.rest_queryparams_len, len(s.rest_queryparams))
	s.rest_urlparams_len = append(s.rest_urlparams_len, len(s.rest_urlparams))

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		s.rest_attrs = append(s.rest_attrs, makeAttributeArray(attribs))
	} else {
		s.rest_attrs = append(s.rest_attrs, nil)
	}
}

// ExitRest_endpoint is called when production rest_endpoint is exited.
func (s *TreeShapeListener) ExitRest_endpoint(ctx *parser.Rest_endpointContext) {
	s.url_prefix = s.url_prefix[:len(s.url_prefix)-1]
	ltop := len(s.rest_urlparams_len) - 1
	if ltop >= 0 {
		l := s.rest_urlparams_len[ltop]
		s.rest_urlparams = s.rest_urlparams[:l]
		s.rest_urlparams_len = s.rest_urlparams_len[:ltop]
	}
	s.rest_attrs = s.rest_attrs[:len(s.rest_attrs)-1]

	if len(s.url_prefix) != len(s.rest_attrs) {
		panic("something is wrong")
	}
}

// EnterCollector_query_var is called when production collector_query_var is entered.
func (s *TreeShapeListener) EnterCollector_query_var(ctx *parser.Collector_query_varContext) {}

// ExitCollector_query_var is called when production collector_query_var is exited.
func (s *TreeShapeListener) ExitCollector_query_var(ctx *parser.Collector_query_varContext) {}

// EnterCollector_query_param is called when production collector_query_param is entered.
func (s *TreeShapeListener) EnterCollector_query_param(ctx *parser.Collector_query_paramContext) {}

// ExitCollector_query_param is called when production collector_query_param is exited.
func (s *TreeShapeListener) ExitCollector_query_param(ctx *parser.Collector_query_paramContext) {}

// EnterCollector_call_stmt is called when production collector_call_stmt is entered.
func (s *TreeShapeListener) EnterCollector_call_stmt(ctx *parser.Collector_call_stmtContext) {
	appName := &sysl.AppName{}
	s.app_name = []string{}
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: strings.TrimSpace(ctx.Target_endpoint().GetText()),
			},
		},
	})
}

// ExitCollector_call_stmt is called when production collector_call_stmt is exited.
func (s *TreeShapeListener) ExitCollector_call_stmt(ctx *parser.Collector_call_stmtContext) {}

// EnterCollector_http_stmt_part is called when production collector_http_stmt_part is entered.
func (s *TreeShapeListener) EnterCollector_http_stmt_part(ctx *parser.Collector_http_stmt_partContext) {
}

// ExitCollector_http_stmt_part is called when production collector_http_stmt_part is exited.
func (s *TreeShapeListener) ExitCollector_http_stmt_part(ctx *parser.Collector_http_stmt_partContext) {
}

// EnterCollector_http_stmt is called when production collector_http_stmt is entered.
func (s *TreeShapeListener) EnterCollector_http_stmt(ctx *parser.Collector_http_stmtContext) {
	text := strings.TrimSpace(ctx.HTTP_VERBS().GetText()) + " " + ctx.Collector_http_stmt_suffix().GetText()

	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Action{
			Action: &sysl.Action{
				Action: text,
			},
		},
	})
}

// ExitCollector_http_stmt is called when production collector_http_stmt is exited.
func (s *TreeShapeListener) ExitCollector_http_stmt(ctx *parser.Collector_http_stmtContext) {}

// EnterCollector_stmts is called when production collector_stmts is entered.
func (s *TreeShapeListener) EnterCollector_stmts(ctx *parser.Collector_stmtsContext) {}

// EnterPublisher is called when production publisher is entered.
func (s *TreeShapeListener) EnterPublisher(ctx *parser.PublisherContext) {}

// ExitPublisher is called when production publisher is exited.
func (s *TreeShapeListener) ExitPublisher(ctx *parser.PublisherContext) {}

// EnterSubscriber is called when production subscriber is entered.
func (s *TreeShapeListener) EnterSubscriber(ctx *parser.SubscriberContext) {}

// ExitSubscriber is called when production subscriber is exited.
func (s *TreeShapeListener) ExitSubscriber(ctx *parser.SubscriberContext) {
	for i := range s.app_name {
		s.app_name[i] = strings.TrimSpace(s.app_name[i])
	}

	s.lastStatement().GetCall().Target.Part = s.app_name
	s.app_name = []string{}
}

// EnterCollector_pubsub_call is called when production collector_pubsub_call is entered.
func (s *TreeShapeListener) EnterCollector_pubsub_call(ctx *parser.Collector_pubsub_callContext) {
	appName := &sysl.AppName{}
	s.app_name = []string{}
	publisher := ctx.Publisher().GetText() + ctx.ARROW_RIGHT().GetText()
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: strings.TrimSpace(publisher + ctx.Name_str().GetText()),
			},
		},
	})
}

// ExitCollector_pubsub_call is called when production collector_pubsub_call is exited.
func (s *TreeShapeListener) ExitCollector_pubsub_call(ctx *parser.Collector_pubsub_callContext) {}

// EnterCollector_action_stmt is called when production collector_action_stmt is entered.
func (s *TreeShapeListener) EnterCollector_action_stmt(ctx *parser.Collector_action_stmtContext) {
	text := ctx.Name_str().GetText()
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Action{
			Action: &sysl.Action{
				Action: text,
			},
		},
	})
}

// ExitCollector_action_stmt is called when production collector_action_stmt is exited.
func (s *TreeShapeListener) ExitCollector_action_stmt(ctx *parser.Collector_action_stmtContext) {}

// ExitCollector_stmts is called when production collector_stmts is exited.
func (s *TreeShapeListener) ExitCollector_stmts(ctx *parser.Collector_stmtsContext) {
	if ctx.Attribs_or_modifiers() != nil {
		stmt := s.lastStatement()
		stmt.Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}
}

// EnterCollector is called when production collector is entered.
func (s *TreeShapeListener) EnterCollector(ctx *parser.CollectorContext) {
	s.typename = ctx.COLLECTOR().GetText()
	ep := s.module.Apps[s.appname].Endpoints[s.typename]

	if ep == nil {
		ep = &sysl.Endpoint{
			Name: s.typename,
		}
		s.module.Apps[s.appname].Endpoints[s.typename] = ep
	}

	if ctx.Collector_stmts(0) != nil {
		ep.Stmt = []*sysl.Statement{}
		if ep.Attrs == nil {
			ep.Attrs = map[string]*sysl.Attribute{}
		}
		s.pushScope(ep)
	}
	ep.SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
}

// ExitCollector is called when production collector is exited.
func (s *TreeShapeListener) ExitCollector(ctx *parser.CollectorContext) {
	if ctx.Collector_stmts(0) != nil {
		s.popScope()
	}
	ep := s.module.Apps[s.appname].Endpoints[s.typename]
	if len(ep.Attrs) == 0 {
		ep.Attrs = nil
	}
	s.typename = ""
}

// EnterEvent is called when production event is entered.
func (s *TreeShapeListener) EnterEvent(ctx *parser.EventContext) {
	if ctx.Name_str() != nil {
		s.typename = ctx.Name_str().GetText()
		// fmt.Printf("Event: %s\n", s.typename)
		if s.module.Apps[s.appname].Endpoints[s.typename] == nil {
			s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
				Name:          s.typename,
				IsPubsub:      true,
				SourceContext: buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()),
			}
		}
		if ctx.Attribs_or_modifiers() != nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
		if ctx.Statements(0) != nil && s.module.Apps[s.appname].Endpoints[s.typename].Stmt == nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Stmt = []*sysl.Statement{}
		}
		if ctx.Statements(0) != nil {
			s.pushScope(s.module.Apps[s.appname].Endpoints[s.typename])
		}
	}
}

// ExitEvent is called when production event is exited.
func (s *TreeShapeListener) ExitEvent(ctx *parser.EventContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
	s.typename = ""
}

// EnterSubscribe is called when production subscribe is entered.
func (s *TreeShapeListener) EnterSubscribe(ctx *parser.SubscribeContext) {
	if ctx.App_name() != nil {
		eventName := ctx.Name_str().GetText()
		s.typename = ctx.App_name().GetText() + ctx.ARROW_RIGHT().GetText() + eventName
		// fmt.Printf("\t%s Subscriber: %s\n", s.appname, s.typename)
		str := strings.Split(ctx.App_name().GetText(), "::")
		for i := range str {
			str[i] = strings.TrimSpace(str[i])
		}
		app_src := &sysl.AppName{
			Part: str,
		}
		s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
			Name:          s.typename,
			Source:        app_src,
			SourceContext: buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()),
		}
		if ctx.Attribs_or_modifiers() != nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
		if ctx.Statements(0) != nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Stmt = []*sysl.Statement{}
			s.pushScope(s.module.Apps[s.appname].Endpoints[s.typename])
		}
		srcAppName := getAppName(app_src)
		srcApp := getApp(app_src, s.module)
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
					Target:   s.module.Apps[s.appname].Name,
					Endpoint: s.typename,
				},
			},
		}
		ep.Stmt = append(ep.Stmt, stmt)
	}
}

// ExitSubscribe is called when production subscribe is exited.
func (s *TreeShapeListener) ExitSubscribe(ctx *parser.SubscribeContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
	s.typename = ""
}

// EnterView_type_spec is called when production view_type_spec is entered.
func (s *TreeShapeListener) EnterView_type_spec(ctx *parser.View_type_specContext) {}

// ExitView_type_spec is called when production view_type_spec is exited.
func (s *TreeShapeListener) ExitView_type_spec(ctx *parser.View_type_specContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
	if type1.GetSet() != nil {
		type1 = type1.GetSet()
	}
	if type1.GetTypeRef() != nil {
		tr := type1.GetTypeRef()
		if tr.Ref.Appname == nil && len(tr.Ref.Path) == 1 {
			tr.Ref.Appname = &sysl.AppName{
				Part: tr.Ref.Path,
			}
			tr.Ref.Path = nil
		}
	}
}

// EnterLiteral is called when production literal is entered.
func (s *TreeShapeListener) EnterLiteral(ctx *parser.LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *TreeShapeListener) ExitLiteral(ctx *parser.LiteralContext) {
	val := &sysl.Value{}
	txt := ctx.GetText()
	var type1 *sysl.Type
	if ctx.E_DECIMAL() != nil {
		val.Value = &sysl.Value_Decimal{
			Decimal: txt,
		}
	} else if ctx.E_DIGITS() != nil {
		iVal, _ := strconv.Atoi(txt)
		val.Value = &sysl.Value_I{
			I: int64(iVal),
		}
	} else if ctx.E_TRUE() != nil || ctx.E_FALSE() != nil {
		val.Value = &sysl.Value_B{
			B: txt == "true",
		}
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_BOOL,
			},
		}

	} else if ctx.E_NULL() != nil {
		val.Value = &sysl.Value_Null_{
			Null: &sysl.Value_Null{},
		}
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_EMPTY,
			},
		}
	} else {
		val.Value = &sysl.Value_S{
			S: fromQString(txt),
		}
	}

	expr := &sysl.Expr{
		Expr: &sysl.Expr_Literal{
			Literal: val,
		},
	}
	if type1 != nil {
		expr.Type = type1
	}
	s.pushExpr(expr)
}

func makeGetAttr(arg *sysl.Expr, name string, nullsafe, setof bool) *sysl.Expr {
	return &sysl.Expr{
		Expr: &sysl.Expr_GetAttr_{
			GetAttr: &sysl.Expr_GetAttr{
				Arg:      arg,
				Attr:     name,
				Setof:    setof,
				Nullsafe: nullsafe,
			},
		},
	}
}

func makeRelExpr(op sysl.Expr_RelExpr_Op) *sysl.Expr {
	return &sysl.Expr{
		Expr: &sysl.Expr_Relexpr{
			Relexpr: &sysl.Expr_RelExpr{
				Op: op,
			},
		},
	}
}

func makeExprName(name string) *sysl.Expr {
	return &sysl.Expr{
		Expr: &sysl.Expr_Name{
			Name: name,
		},
	}
}

func makeUnaryExpr(op sysl.Expr_UnExpr_Op, expr *sysl.Expr) *sysl.Expr {
	return &sysl.Expr{
		Expr: &sysl.Expr_Unexpr{
			Unexpr: &sysl.Expr_UnExpr{
				Op:  op,
				Arg: expr,
			},
		},
	}
}

func makeBinaryExpr(op sysl.Expr_BinExpr_Op, lhs, rhs *sysl.Expr) *sysl.Expr {
	return &sysl.Expr{
		Expr: &sysl.Expr_Binexpr{
			Binexpr: &sysl.Expr_BinExpr{
				Op:  op,
				Lhs: lhs,
				Rhs: rhs,
			},
		},
	}
}

func makeGetAttrFromRefName(arg *sysl.Expr, names []string, nullsafe, setof bool) *sysl.Expr {
	for _, name := range names {
		arg = makeGetAttr(arg, name, nullsafe, setof)
		nullsafe = false
		setof = false
	}
	return arg
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

// EnterExpr_table_of_op is called when production expr_table_of_op is entered.
func (s *TreeShapeListener) EnterExpr_table_of_op(ctx *parser.Expr_table_of_opContext) {}

// ExitExpr_table_of_op is called when production expr_table_of_op is exited.
func (s *TreeShapeListener) ExitExpr_table_of_op(ctx *parser.Expr_table_of_opContext) {
	nullsafe := ctx.E_NULLSAFE_DOT() != nil
	table_of := ctx.E_TABLE_OF() != nil
	attrName := ctx.E_Name().GetText()
	arg := s.popExpr()
	s.pushExpr(makeGetAttr(arg, attrName, nullsafe, table_of))
}

// EnterFunc_arg is called when production func_arg is entered.
func (s *TreeShapeListener) EnterFunc_arg(ctx *parser.Func_argContext) {}

// ExitFunc_arg is called when production func_arg is exited.
func (s *TreeShapeListener) ExitFunc_arg(ctx *parser.Func_argContext) {
	arg := s.popExpr()
	if s.topExpr().GetRelexpr() != nil {
		expr := s.topExpr().GetRelexpr()
		expr.Arg = append(expr.Arg, arg)
	} else if s.topExpr().GetCall() != nil {
		expr := s.topExpr().GetCall()
		expr.Arg = append(expr.Arg, arg)
	} else {
		fmt.Printf("%v\n", arg)
		fmt.Printf("%v\n", s.topExpr())
		panic("ExitFunc_arg: should not be here")
	}
}

// EnterFunc_args is called when production func_args is entered.
func (s *TreeShapeListener) EnterFunc_args(ctx *parser.Func_argsContext) {}

// ExitFunc_args is called when production func_args is exited.
func (s *TreeShapeListener) ExitFunc_args(ctx *parser.Func_argsContext) {}

// EnterExpr_func is called when production expr_func is entered.
func (s *TreeShapeListener) EnterExpr_func(ctx *parser.Expr_funcContext) {
	var funcName string
	if ctx.E_FUNC() != nil {
		funcName = ctx.E_FUNC().GetText()
	} else if ctx.E_Name() != nil {
		funcName = ctx.E_Name().GetText()
	} else {
		funcName = ctx.NativeDataTypes().GetText()
	}

	expr := &sysl.Expr{
		Expr: &sysl.Expr_Call_{
			Call: &sysl.Expr_Call{
				Func: funcName,
				Arg:  []*sysl.Expr{},
			},
		},
	}
	s.pushExpr(expr)
}

// ExitExpr_func is called when production expr_func is exited.
func (s *TreeShapeListener) ExitExpr_func(ctx *parser.Expr_funcContext) {}

// EnterRank_expr is called when production rank_expr is entered.
func (s *TreeShapeListener) EnterRank_expr(ctx *parser.Rank_exprContext) {}

// ExitRank_expr is called when production rank_expr is exited.
func (s *TreeShapeListener) ExitRank_expr(ctx *parser.Rank_exprContext) {
	expr := s.popExpr()
	relexpr := s.topExpr().GetRelexpr()
	relexpr.Arg = append(relexpr.Arg, expr)
	desc := false
	if ctx.E_DESC() != nil {
		desc = true
	}
	relexpr.Descending = append(relexpr.Descending, desc)
}

// EnterRank_expr_list is called when production rank_expr_list is entered.
func (s *TreeShapeListener) EnterRank_expr_list(ctx *parser.Rank_expr_listContext) {}

// ExitRank_expr_list is called when production rank_expr_list is exited.
func (s *TreeShapeListener) ExitRank_expr_list(ctx *parser.Rank_expr_listContext) {}

// EnterExpr_rank_func is called when production expr_rank_func is entered.
func (s *TreeShapeListener) EnterExpr_rank_func(ctx *parser.Expr_rank_funcContext) {
	target := s.popExpr()
	s.pushExpr(makeRelExpr(sysl.Expr_RelExpr_RANK))
	relexpr := s.topExpr().GetRelexpr()
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
func (s *TreeShapeListener) ExitExpr_rank_func(ctx *parser.Expr_rank_funcContext) {
	expr := s.topExpr()
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
	s.pushExpr(makeRelExpr(op))
	relexpr := s.topExpr().GetRelexpr()
	relexpr.Target = expr
}

// ExitExpr_agg_func is called when production expr_agg_func is exited.
func (s *TreeShapeListener) ExitExpr_agg_func(ctx *parser.Expr_agg_funcContext) {
	relexpr := s.topExpr().GetRelexpr()
	if relexpr.Scopevar == "" {
		relexpr.Scopevar = "."
	}
}

// EnterE_scope_var is called when production e_scope_var is entered.
func (s *TreeShapeListener) EnterE_scope_var(ctx *parser.E_scope_varContext) {}

// ExitE_scope_var is called when production e_scope_var is exited.
func (s *TreeShapeListener) ExitE_scope_var(ctx *parser.E_scope_varContext) {
	expr := s.topExpr()
	if expr.GetBinexpr() != nil {
		binexpr := expr.GetBinexpr()
		binexpr.Scopevar = ctx.E_Name().GetText()
	} else if expr.GetRelexpr() != nil {
		relexpr := expr.GetRelexpr()
		relexpr.Scopevar = ctx.E_Name().GetText()
	}
}

// EnterFirst_func_target is called when production first_func_target is entered.
func (s *TreeShapeListener) EnterFirst_func_target(ctx *parser.First_func_targetContext) {}

// ExitFirst_func_target is called when production first_func_target is exited.
func (s *TreeShapeListener) ExitFirst_func_target(ctx *parser.First_func_targetContext) {
	var expr *sysl.Expr
	if ctx.E_NULL() != nil {
		val := &sysl.Value{}
		val.Value = &sysl.Value_Null_{
			Null: &sysl.Value_Null{},
		}

		expr = &sysl.Expr{
			Expr: &sysl.Expr_Literal{
				Literal: val,
			},
			Type: &sysl.Type{
				Type: &sysl.Type_Primitive_{
					Primitive: sysl.Type_EMPTY,
				},
			},
		}
	} else {
		expr = s.popExpr()
	}
	relexpr := s.topExpr().GetRelexpr()
	relexpr.Arg = append(relexpr.Arg, expr)
}

// EnterExpr_first_func is called when production expr_first_func is entered.
func (s *TreeShapeListener) EnterExpr_first_func(ctx *parser.Expr_first_funcContext) {
	expr := makeRelExpr(sysl.Expr_RelExpr_FIRST_BY)
	relexpr := expr.GetRelexpr()
	relexpr.Target = s.popExpr()
	s.pushExpr(expr)
}

// ExitExpr_first_func is called when production expr_first_func is exited.
func (s *TreeShapeListener) ExitExpr_first_func(ctx *parser.Expr_first_funcContext) {
	relexpr := s.topExpr().GetRelexpr()
	if relexpr.Scopevar == "" {
		relexpr.Scopevar = "."
	}
	// Legacy has double the length of the arg array.
	relexpr.Arg = append(relexpr.Arg, relexpr.Arg...)
}

// EnterE_single_arg_func is called when production e_single_arg_func is entered.
func (s *TreeShapeListener) EnterE_single_arg_func(ctx *parser.E_single_arg_funcContext) {}

// ExitE_single_arg_func is called when production e_single_arg_func is exited.
func (s *TreeShapeListener) ExitE_single_arg_func(ctx *parser.E_single_arg_funcContext) {
	op := sysl.Expr_BinExpr_NO_Op
	if ctx.E_RELOPS_WHERE() != nil {
		op = sysl.Expr_BinExpr_WHERE
	} else if ctx.E_RELOPS_FLATTEN() != nil {
		op = sysl.Expr_BinExpr_FLATTEN
	}
	s.pushExpr(makeBinaryExpr(op, nil, nil))
}

// EnterExpr_single_arg_func is called when production expr_single_arg_func is entered.
func (s *TreeShapeListener) EnterExpr_single_arg_func(ctx *parser.Expr_single_arg_funcContext) {}

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
	s.pushExpr(expr)
}

// EnterExpr_any_func is called when production expr_any_func is entered.
func (s *TreeShapeListener) EnterExpr_any_func(ctx *parser.Expr_any_funcContext) {}

// ExitExpr_any_func is called when production expr_any_func is exited.
func (s *TreeShapeListener) ExitExpr_any_func(ctx *parser.Expr_any_funcContext) {
	limit := s.popExpr()
	arg := s.popExpr()
	expr := &sysl.Expr{
		Expr: &sysl.Expr_Call_{
			Call: &sysl.Expr_Call{
				Func: ".any",
				Arg:  []*sysl.Expr{arg, limit},
			},
		},
	}
	s.pushExpr(expr)
}

// EnterExpr_single_or_null is called when production expr_single_or_null is entered.
func (s *TreeShapeListener) EnterExpr_single_or_null(ctx *parser.Expr_single_or_nullContext) {}

// ExitExpr_single_or_null is called when production expr_single_or_null is exited.
func (s *TreeShapeListener) ExitExpr_single_or_null(ctx *parser.Expr_single_or_nullContext) {
	op := sysl.Expr_UnExpr_SINGLE_OR_NULL
	if ctx.GetText() == "single" {
		op = sysl.Expr_UnExpr_SINGLE
	}
	expr := makeUnaryExpr(op, s.popExpr())
	s.pushExpr(expr)
}

// EnterExpr_snapshot is called when production expr_snapshot is entered.
func (s *TreeShapeListener) EnterExpr_snapshot(ctx *parser.Expr_snapshotContext) {}

// ExitExpr_snapshot is called when production expr_snapshot is exited.
func (s *TreeShapeListener) ExitExpr_snapshot(ctx *parser.Expr_snapshotContext) {
	expr := s.popExpr()
	s.pushExpr(makeRelExpr(sysl.Expr_RelExpr_SNAPSHOT))
	relexpr := s.topExpr().GetRelexpr()
	relexpr.Target = expr
}

// EnterExpr_count is called when production expr_count is entered.
func (s *TreeShapeListener) EnterExpr_count(ctx *parser.Expr_countContext) {}

// ExitExpr_count is called when production expr_count is exited.
func (s *TreeShapeListener) ExitExpr_count(ctx *parser.Expr_countContext) {
	expr := s.popExpr()
	callExpr := &sysl.Expr{
		Expr: &sysl.Expr_Call_{
			Call: &sysl.Expr_Call{
				Func: ".count",
				Arg:  []*sysl.Expr{expr},
			},
		},
	}
	s.pushExpr(callExpr)
}

// EnterExpr_navigate_attr is called when production expr_navigate_attr is entered.
func (s *TreeShapeListener) EnterExpr_navigate_attr(ctx *parser.Expr_navigate_attrContext) {}

// ExitExpr_navigate_attr is called when production expr_navigate_attr is exited.
func (s *TreeShapeListener) ExitExpr_navigate_attr(ctx *parser.Expr_navigate_attrContext) {
	nav := s.topExpr().GetNavigate()
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

	expr := &sysl.Expr{
		Expr: &sysl.Expr_Navigate_{
			Navigate: nav,
		},
	}
	s.pushExpr(expr)
}

// ExitExpr_navigate is called when production expr_navigate is exited.
func (s *TreeShapeListener) ExitExpr_navigate(ctx *parser.Expr_navigateContext) {}

// EnterMatching_rhs is called when production matching_rhs is entered.
func (s *TreeShapeListener) EnterMatching_rhs(ctx *parser.Matching_rhsContext) {
	if ctx.E_Name() == nil && ctx.AtomT_paren() == nil {
		s.pushExpr(makeExprName("."))
	}
}

// ExitMatching_rhs is called when production matching_rhs is exited.
func (s *TreeShapeListener) ExitMatching_rhs(ctx *parser.Matching_rhsContext) {
	if ctx.E_Name() != nil {
		s.pushExpr(makeExprName(ctx.E_Name().GetText()))
	}
}

// EnterSquiggly_args is called when production squiggly_args is entered.
func (s *TreeShapeListener) EnterSquiggly_args(ctx *parser.Squiggly_argsContext) {
	binexpr := s.topExpr().GetBinexpr()
	names := ctx.AllE_Name()
	for _, name := range names {
		binexpr.AttrName = append(binexpr.AttrName, name.GetText())
	}
}

// ExitSquiggly_args is called when production squiggly_args is exited.
func (s *TreeShapeListener) ExitSquiggly_args(ctx *parser.Squiggly_argsContext) {}

// EnterExpr_matching is called when production expr_matching is entered.
func (s *TreeShapeListener) EnterExpr_matching(ctx *parser.Expr_matchingContext) {
	op := sysl.Expr_BinExpr_TO_MATCHING
	if ctx.E_NOT() != nil {
		op = sysl.Expr_BinExpr_TO_NOT_MATCHING
	}
	lhs := s.popExpr()
	s.pushExpr(makeBinaryExpr(op, lhs, nil))
}

// ExitExpr_matching is called when production expr_matching is exited.
func (s *TreeShapeListener) ExitExpr_matching(ctx *parser.Expr_matchingContext) {
	var binexpr *sysl.Expr_BinExpr
	rhs := s.popExpr()
	binexpr = s.topExpr().GetBinexpr()
	binexpr.Rhs = rhs
	if len(binexpr.AttrName) == 0 {
		binexpr.AttrName = []string{"*"}
	}
}

// EnterRelop is called when production relop is entered.
func (s *TreeShapeListener) EnterRelop(ctx *parser.RelopContext) {}

// ExitRelop is called when production relop is exited.
func (s *TreeShapeListener) ExitRelop(ctx *parser.RelopContext) {}

// EnterList_item is called when production list_item is entered.
func (s *TreeShapeListener) EnterList_item(ctx *parser.List_itemContext) {}

// ExitList_item is called when production list_item is exited.
func (s *TreeShapeListener) ExitList_item(ctx *parser.List_itemContext) {
	item := s.popExpr()

	if s.topExpr().GetSet() != nil {
		list := s.topExpr().GetSet()
		list.Expr = append(list.Expr, item)
	} else if s.topExpr().GetList() != nil {
		list := s.topExpr().GetList()
		list.Expr = append(list.Expr, item)
	}
}

// EnterExpr_list is called when production expr_list is entered.
func (s *TreeShapeListener) EnterExpr_list(ctx *parser.Expr_listContext) {}

// ExitExpr_list is called when production expr_list is exited.
func (s *TreeShapeListener) ExitExpr_list(ctx *parser.Expr_listContext) {
	expr := s.topExpr()
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
	expr := &sysl.Expr{
		Expr: &sysl.Expr_Set{
			Set: &sysl.Expr_List{
				Expr: []*sysl.Expr{},
			},
		},
		Type: &sysl.Type{
			Type: &sysl.Type_Set{
				Set: &sysl.Type{},
			},
		},
	}
	s.pushExpr(expr)
}

// ExitExpr_set is called when production expr_set is exited.
func (s *TreeShapeListener) ExitExpr_set(ctx *parser.Expr_setContext) {}

// EnterEmpty_tuple is called when production empty_tuple is entered.
func (s *TreeShapeListener) EnterEmpty_tuple(ctx *parser.Empty_tupleContext) {}

// ExitEmpty_tuple is called when production empty_tuple is exited.
func (s *TreeShapeListener) ExitEmpty_tuple(ctx *parser.Empty_tupleContext) {
	s.pushExpr(&sysl.Expr{
		Expr: &sysl.Expr_Tuple_{
			Tuple: &sysl.Expr_Tuple{},
		},
		Type: &sysl.Type{
			Type: &sysl.Type_Tuple_{
				Tuple: &sysl.Type_Tuple{},
			},
		},
	})
}

// EnterAtom_dot_relop is called when production atom_dot_relop is entered.
func (s *TreeShapeListener) EnterAtom_dot_relop(ctx *parser.Atom_dot_relopContext) {}

// ExitAtom_dot_relop is called when production atom_dot_relop is exited.
func (s *TreeShapeListener) ExitAtom_dot_relop(ctx *parser.Atom_dot_relopContext) {}

// EnterAtomT_implied_dot is called when production atomT_implied_dot is entered.
func (s *TreeShapeListener) EnterAtomT_implied_dot(ctx *parser.AtomT_implied_dotContext) {
	s.pushExpr(makeExprName("."))
}

// ExitAtomT_implied_dot is called when production atomT_implied_dot is exited.
func (s *TreeShapeListener) ExitAtomT_implied_dot(ctx *parser.AtomT_implied_dotContext) {}

// EnterAtomT_name is called when production atomT_name is entered.
func (s *TreeShapeListener) EnterAtomT_name(ctx *parser.AtomT_nameContext) {}

// ExitAtomT_name is called when production atomT_name is exited.
func (s *TreeShapeListener) ExitAtomT_name(ctx *parser.AtomT_nameContext) {
	if ctx.E_Name() != nil {
		txt := ctx.E_Name().GetText()
		s.pushExpr(makeExprName(txt))
	} else if ctx.E_WHATEVER() != nil {
		s.pushExpr(&sysl.Expr{})
	} else if ctx.E_DOT() != nil {
		s.pushExpr(makeExprName("."))
	} else {
		panic("ExitAtomT_name: should not be here")
	}
}

// EnterAtomT_paren is called when production atomT_paren is entered.
func (s *TreeShapeListener) EnterAtomT_paren(ctx *parser.AtomT_parenContext) {}

// ExitAtomT_paren is called when production atomT_paren is exited.
func (s *TreeShapeListener) ExitAtomT_paren(ctx *parser.AtomT_parenContext) {}

// EnterExpr_atom_list is called when production expr_atom_list is entered.
func (s *TreeShapeListener) EnterExpr_atom_list(ctx *parser.Expr_atom_listContext) {
	expr := &sysl.Expr{
		Expr: &sysl.Expr_List_{
			List: &sysl.Expr_List{
				Expr: []*sysl.Expr{},
			},
		},
	}
	s.pushExpr(expr)
}

// ExitExpr_atom_list is called when production expr_atom_list is exited.
func (s *TreeShapeListener) ExitExpr_atom_list(ctx *parser.Expr_atom_listContext) {}

// EnterAtomT is called when production atomT is entered.
func (s *TreeShapeListener) EnterAtomT(ctx *parser.AtomTContext) {}

// ExitAtomT is called when production atomT is exited.
func (s *TreeShapeListener) ExitAtomT(ctx *parser.AtomTContext) {}

// EnterAtom is called when production atom is entered.
func (s *TreeShapeListener) EnterAtom(ctx *parser.AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *TreeShapeListener) ExitAtom(ctx *parser.AtomContext) {}

// EnterPowerT is called when production powerT is entered.
func (s *TreeShapeListener) EnterPowerT(ctx *parser.PowerTContext) {}

// ExitPowerT is called when production powerT is exited.
func (s *TreeShapeListener) ExitPowerT(ctx *parser.PowerTContext) {}

// EnterPower is called when production power is entered.
func (s *TreeShapeListener) EnterPower(ctx *parser.PowerContext) {}

// ExitPower is called when production power is exited.
func (s *TreeShapeListener) ExitPower(ctx *parser.PowerContext) {
	if ctx.PowerT() != nil {
		rhs := s.popExpr()
		lhs := s.popExpr()
		s.pushExpr(makeBinaryExpr(sysl.Expr_BinExpr_POW, lhs, rhs))
	}
}

// EnterUnaryTerm is called when production unaryTerm is entered.
func (s *TreeShapeListener) EnterUnaryTerm(ctx *parser.UnaryTermContext) {}

// ExitUnaryTerm is called when production unaryTerm is exited.
func (s *TreeShapeListener) ExitUnaryTerm(ctx *parser.UnaryTermContext) {
	op := sysl.Expr_UnExpr_NO_Op

	if ctx.E_TILDE() != nil {
		op = sysl.Expr_UnExpr_INV
	} else if ctx.E_NOT() != nil {
		op = sysl.Expr_UnExpr_NOT
	} else if ctx.E_MINUS() != nil {
		op = sysl.Expr_UnExpr_NEG
	} else if ctx.E_PLUS() != nil {
		op = sysl.Expr_UnExpr_POS
	}
	if op != sysl.Expr_UnExpr_NO_Op {
		expr := makeUnaryExpr(op, s.popExpr())
		s.pushExpr(expr)
	}
}

// EnterTermT is called when production termT is entered.
func (s *TreeShapeListener) EnterTermT(ctx *parser.TermTContext) {}

// ExitTermT is called when production termT is exited.
func (s *TreeShapeListener) ExitTermT(ctx *parser.TermTContext) {
	op := sysl.Expr_BinExpr_MOD
	if ctx.E_STAR() != nil {
		op = sysl.Expr_BinExpr_MUL
	} else if ctx.E_DIVIDE() != nil {
		op = sysl.Expr_BinExpr_DIV
	}
	rhs := s.popExpr()
	lhs := s.popExpr()

	s.pushExpr(makeBinaryExpr(op, lhs, rhs))
}

// EnterTerm is called when production term is entered.
func (s *TreeShapeListener) EnterTerm(ctx *parser.TermContext) {}

// ExitTerm is called when production term is exited.
func (s *TreeShapeListener) ExitTerm(ctx *parser.TermContext) {}

// EnterBinexprT is called when production binexprT is entered.
func (s *TreeShapeListener) EnterBinexprT(ctx *parser.BinexprTContext) {}

// ExitBinexprT is called when production binexprT is exited.
func (s *TreeShapeListener) ExitBinexprT(ctx *parser.BinexprTContext) {
	op := sysl.Expr_BinExpr_NO_Op

	if ctx.E_PLUS() != nil {
		op = sysl.Expr_BinExpr_ADD
	} else {
		op = sysl.Expr_BinExpr_SUB
	}
	rhs := s.popExpr()
	lhs := s.popExpr()
	s.pushExpr(makeBinaryExpr(op, lhs, rhs))
}

// EnterBinexpr is called when production binexpr is entered.
func (s *TreeShapeListener) EnterBinexpr(ctx *parser.BinexprContext) {}

// ExitBinexpr is called when production binexpr is exited.
func (s *TreeShapeListener) ExitBinexpr(ctx *parser.BinexprContext) {}

// EnterE_compare_ops is called when production e_compare_ops is entered.
func (s *TreeShapeListener) EnterE_compare_ops(ctx *parser.E_compare_opsContext) {}

// ExitE_compare_ops is called when production e_compare_ops is exited.
func (s *TreeShapeListener) ExitE_compare_ops(ctx *parser.E_compare_opsContext) {}

// EnterExpr_rel is called when production expr_rel is entered.
func (s *TreeShapeListener) EnterExpr_rel(ctx *parser.Expr_relContext) {}

// ExitExpr_rel is called when production expr_rel is exited.
func (s *TreeShapeListener) ExitExpr_rel(ctx *parser.Expr_relContext) {
	if ctx.E_compare_ops(0) != nil {
		for i := len(ctx.AllE_compare_ops()) - 1; i >= 0; i-- {
			op := s.opmap[ctx.E_compare_ops(i).GetText()]
			rhs := s.popExpr()
			lhs := s.popExpr()
			s.pushExpr(makeBinaryExpr(op, lhs, rhs))
		}
	}
}

// EnterExpr_bitand is called when production expr_bitand is entered.
func (s *TreeShapeListener) EnterExpr_bitand(ctx *parser.Expr_bitandContext) {}

// ExitExpr_bitand is called when production expr_bitand is exited.
func (s *TreeShapeListener) ExitExpr_bitand(ctx *parser.Expr_bitandContext) {
	if ctx.E_AMP(0) != nil || ctx.E_AND(0) != nil {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_BITAND, len(ctx.AllExpr_rel())-1))
	}
}

// EnterExpr_bitxor is called when production expr_bitxor is entered.
func (s *TreeShapeListener) EnterExpr_bitxor(ctx *parser.Expr_bitxorContext) {}

// ExitExpr_bitxor is called when production expr_bitxor is exited.
func (s *TreeShapeListener) ExitExpr_bitxor(ctx *parser.Expr_bitxorContext) {
	if ctx.E_XOR(0) != nil {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_BITXOR, len(ctx.AllE_XOR())))
	}
}

// EnterExpr_bitor is called when production expr_bitor is entered.
func (s *TreeShapeListener) EnterExpr_bitor(ctx *parser.Expr_bitorContext) {}

// ExitExpr_bitor is called when production expr_bitor is exited.
func (s *TreeShapeListener) ExitExpr_bitor(ctx *parser.Expr_bitorContext) {
	if ctx.E_BITOR(0) != nil {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_BITOR, len(ctx.AllE_BITOR())))
	}
}

// EnterExpr_and is called when production expr_and is entered.
func (s *TreeShapeListener) EnterExpr_and(ctx *parser.Expr_andContext) {}

// ExitExpr_and is called when production expr_and is exited.
func (s *TreeShapeListener) ExitExpr_and(ctx *parser.Expr_andContext) {
	if ctx.E_DOUBLE_AMP(0) != nil {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_AND, len(ctx.AllE_DOUBLE_AMP())))
	}
}

// EnterExpr_or is called when production expr_or is entered.
func (s *TreeShapeListener) EnterExpr_or(ctx *parser.Expr_orContext) {}

// ExitExpr_or is called when production expr_or is exited.
func (s *TreeShapeListener) ExitExpr_or(ctx *parser.Expr_orContext) {
	if ctx.E_LOGIC_OR(0) != nil {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_OR, len(ctx.AllE_LOGIC_OR())))
	}
}

// EnterExpr_but_not is called when production expr_but_not is entered.
func (s *TreeShapeListener) EnterExpr_but_not(ctx *parser.Expr_but_notContext) {}

// ExitExpr_but_not is called when production expr_but_not is exited.
func (s *TreeShapeListener) ExitExpr_but_not(ctx *parser.Expr_but_notContext) {
	if ctx.E_BUTNOT(0) != nil {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_BUTNOT, len(ctx.AllE_BUTNOT())))
	}
}

// EnterExpr_coalesce is called when production expr_coalesce is entered.
func (s *TreeShapeListener) EnterExpr_coalesce(ctx *parser.Expr_coalesceContext) {}

func (s *TreeShapeListener) reverseOp(op sysl.Expr_BinExpr_Op, count int) *sysl.Expr {
	if count == 0 {
		return s.popExpr()
	}
	rhs := s.popExpr()
	lhs := s.reverseOp(op, count-1)
	return makeBinaryExpr(op, lhs, rhs)
}

// ExitExpr_coalesce is called when production expr_coalesce is exited.
func (s *TreeShapeListener) ExitExpr_coalesce(ctx *parser.Expr_coalesceContext) {
	exprs := ctx.AllExpr_but_not()
	if len(exprs) > 1 {
		s.pushExpr(s.reverseOp(sysl.Expr_BinExpr_COALESCE, len(exprs)-1))
	}
}

// EnterIf_one_liner is called when production if_one_liner is entered.
func (s *TreeShapeListener) EnterIf_one_liner(ctx *parser.If_one_linerContext) {}

// ExitIf_one_liner is called when production if_one_liner is exited.
func (s *TreeShapeListener) ExitIf_one_liner(ctx *parser.If_one_linerContext) {
	if_false := s.popExpr()
	if_true := s.popExpr()
	cond := s.popExpr()
	expr := &sysl.Expr{
		Expr: &sysl.Expr_Ifelse{
			Ifelse: &sysl.Expr_IfElse{
				Cond:     cond,
				IfTrue:   if_true,
				IfFalse:  if_false,
				Nullsafe: ctx.E_QN() != nil,
			},
		},
	}
	s.pushExpr(expr)
}

// EnterElse_block_stmt is called when production else_block_stmt is entered.
func (s *TreeShapeListener) EnterElse_block_stmt(ctx *parser.Else_block_stmtContext) {}

// ExitElse_block_stmt is called when production else_block_stmt is exited.
func (s *TreeShapeListener) ExitElse_block_stmt(ctx *parser.Else_block_stmtContext) {
	stmt := s.popExpr()
	ifelse := s.topExpr()
	addStmt(ifelse, stmt)
}

// EnterControl_item is called when production control_item is entered.
func (s *TreeShapeListener) EnterControl_item(ctx *parser.Control_itemContext) {}

// ExitControl_item is called when production control_item is exited.
func (s *TreeShapeListener) ExitControl_item(ctx *parser.Control_itemContext) {
	control := s.popExpr()
	expr := s.topExpr()
	ifelse := expr.GetIfelse()

	if ifelse.Cond.GetBinexpr() != nil && ifelse.Cond.GetBinexpr().Op == sysl.Expr_BinExpr_EQ {
		if ifelse.Cond.GetBinexpr().Rhs == nil {
			ifelse.Cond.GetBinexpr().Rhs = control
		} else {
			newCond := makeBinaryExpr(sysl.Expr_BinExpr_EQ, ifelse.Cond.GetBinexpr().Lhs, control)
			addIfElseControl(expr, newCond)
		}
	} else {
		if len(ifelse.Cond.GetCall().Arg) == 0 {
			ifelse.Cond.GetCall().Arg = append(ifelse.Cond.GetCall().Arg, control)
		} else {
			newCond := &sysl.Expr{
				Expr: &sysl.Expr_Call_{
					Call: &sysl.Expr_Call{
						Func: "bool",
						Arg:  []*sysl.Expr{control},
					},
				},
			}
			addIfElseControl(expr, newCond)
		}
	}
}

// EnterIf_controls is called when production if_controls is entered.
func (s *TreeShapeListener) EnterIf_controls(ctx *parser.If_controlsContext) {}

// ExitIf_controls is called when production if_controls is exited.
func (s *TreeShapeListener) ExitIf_controls(ctx *parser.If_controlsContext) {}

// EnterCond_block is called when production cond_block is entered.
func (s *TreeShapeListener) EnterCond_block(ctx *parser.Cond_blockContext) {}

// ExitCond_block is called when production cond_block is exited.
func (s *TreeShapeListener) ExitCond_block(ctx *parser.Cond_blockContext) {}

// EnterFinal_else is called when production final_else is entered.
func (s *TreeShapeListener) EnterFinal_else(ctx *parser.Final_elseContext) {}

// ExitFinal_else is called when production final_else is exited.
func (s *TreeShapeListener) ExitFinal_else(ctx *parser.Final_elseContext) {
	elseStmt := s.popExpr()
	ifelse := s.topExpr()
	for ifelse.GetIfelse().IfFalse != nil {
		ifelse = ifelse.GetIfelse().IfFalse
	}
	ifelse.GetIfelse().IfFalse = elseStmt
}

// EnterIfvar is called when production ifvar is entered.
func (s *TreeShapeListener) EnterIfvar(ctx *parser.IfvarContext) {}

// ExitIfvar is called when production ifvar is exited.
func (s *TreeShapeListener) ExitIfvar(ctx *parser.IfvarContext) {
	lhs := s.popExpr()
	ifelse := s.topExpr()
	ifelse.GetIfelse().Cond = makeBinaryExpr(sysl.Expr_BinExpr_EQ, lhs, nil)
}

// EnterIf_multiple_lines is called when production if_multiple_lines is entered.
func (s *TreeShapeListener) EnterIf_multiple_lines(ctx *parser.If_multiple_linesContext) {
	expr := makeIfElseExpr()
	if ctx.Ifvar() == nil {
		expr.GetIfelse().Cond = &sysl.Expr{
			Expr: &sysl.Expr_Call_{
				Call: &sysl.Expr_Call{
					Func: "bool",
					Arg:  []*sysl.Expr{},
				},
			},
		}
	}
	s.pushExpr(expr)
}

// ExitIf_multiple_lines is called when production if_multiple_lines is exited.
func (s *TreeShapeListener) ExitIf_multiple_lines(ctx *parser.If_multiple_linesContext) {}

// EnterExpr_if_else is called when production expr_if_else is entered.
func (s *TreeShapeListener) EnterExpr_if_else(ctx *parser.Expr_if_elseContext) {}

// ExitExpr_if_else is called when production expr_if_else is exited.
func (s *TreeShapeListener) ExitExpr_if_else(ctx *parser.Expr_if_elseContext) {}

// EnterExpr is called when production expr is entered.
func (s *TreeShapeListener) EnterExpr(ctx *parser.ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *TreeShapeListener) ExitExpr(ctx *parser.ExprContext) {}

// EnterExpr_assign is called when production expr_assign is entered.
func (s *TreeShapeListener) EnterExpr_assign(ctx *parser.Expr_assignContext) {}

// ExitExpr_assign is called when production expr_assign is exited.
func (s *TreeShapeListener) ExitExpr_assign(ctx *parser.Expr_assignContext) {}

// EnterExpr_simple_assign is called when production expr_simple_assign is entered.
func (s *TreeShapeListener) EnterExpr_simple_assign(ctx *parser.Expr_simple_assignContext) {}

// ExitExpr_simple_assign is called when production expr_simple_assign is exited.
func (s *TreeShapeListener) ExitExpr_simple_assign(ctx *parser.Expr_simple_assignContext) {
	expr := s.popExpr()
	tx := s.topExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v\n", s.topExpr())
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

// EnterExpr_let_statement is called when production expr_let_statement is entered.
func (s *TreeShapeListener) EnterExpr_let_statement(ctx *parser.Expr_let_statementContext) {}

// ExitExpr_let_statement is called when production expr_let_statement is exited.
func (s *TreeShapeListener) ExitExpr_let_statement(ctx *parser.Expr_let_statementContext) {
	expr := s.popExpr()
	tx := s.topExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.topExpr())
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

// EnterExpr_table_of_statement is called when production expr_table_of_statement is entered.
func (s *TreeShapeListener) EnterExpr_table_of_statement(ctx *parser.Expr_table_of_statementContext) {}

// ExitExpr_table_of_statement is called when production expr_table_of_statement is exited.
func (s *TreeShapeListener) ExitExpr_table_of_statement(ctx *parser.Expr_table_of_statementContext) {
	expr := s.popExpr()
	tx := s.topExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.topExpr())
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

// EnterExpr_dot_assign is called when production expr_dot_assign is entered.
func (s *TreeShapeListener) EnterExpr_dot_assign(ctx *parser.Expr_dot_assignContext) {}

// ExitExpr_dot_assign is called when production expr_dot_assign is exited.
func (s *TreeShapeListener) ExitExpr_dot_assign(ctx *parser.Expr_dot_assignContext) {
	tx := s.topExpr().GetTransform()
	names := strings.Split(strings.TrimRight(ctx.GetText(), " \r\n"), ".")
	var arg *sysl.Expr

	if names[0] == "" {
		names[0] = "."
	}
	arg = makeExprName(names[0])
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Assign_{
			Assign: &sysl.Expr_Transform_Stmt_Assign{
				Name: names[1],
				Expr: makeGetAttr(arg, names[1], false, false),
			},
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// EnterExpr_statement_no_nl is called when production expr_statement_no_nl is entered.
func (s *TreeShapeListener) EnterExpr_statement_no_nl(ctx *parser.Expr_statement_no_nlContext) {}

// ExitExpr_statement_no_nl is called when production expr_statement_no_nl is exited.
func (s *TreeShapeListener) ExitExpr_statement_no_nl(ctx *parser.Expr_statement_no_nlContext) {}

// EnterExpr_statement is called when production expr_statement is entered.
func (s *TreeShapeListener) EnterExpr_statement(ctx *parser.Expr_statementContext) {}

// ExitExpr_statement is called when production expr_statement is exited.
func (s *TreeShapeListener) ExitExpr_statement(ctx *parser.Expr_statementContext) {}

// EnterExpr_inject_stmt is called when production expr_inject_stmt is entered.
func (s *TreeShapeListener) EnterExpr_inject_stmt(ctx *parser.Expr_inject_stmtContext) {}

// ExitExpr_inject_stmt is called when production expr_inject_stmt is exited.
func (s *TreeShapeListener) ExitExpr_inject_stmt(ctx *parser.Expr_inject_stmtContext) {
	expr := s.popExpr()
	tx := s.topExpr().GetTransform()
	if tx == nil {
		fmt.Printf("%v", s.topExpr())
		panic("ExitExpr_inject_stmt: Unexpected expression!")
	}
	expr.GetCall().Arg = append(expr.GetCall().Arg, makeExprName("out"))
	stmt := &sysl.Expr_Transform_Stmt{
		Stmt: &sysl.Expr_Transform_Stmt_Inject{
			Inject: expr,
		},
	}
	tx.Stmt = append(tx.Stmt, stmt)
}

// EnterExpr_stmt is called when production expr_stmt is entered.
func (s *TreeShapeListener) EnterExpr_stmt(ctx *parser.Expr_stmtContext) {}

// ExitExpr_stmt is called when production expr_stmt is exited.
func (s *TreeShapeListener) ExitExpr_stmt(ctx *parser.Expr_stmtContext) {}

// EnterSet_of is called when production set_of is entered.
func (s *TreeShapeListener) EnterSet_of(ctx *parser.Set_ofContext) {}

// ExitSet_of is called when production set_of is exited.
func (s *TreeShapeListener) ExitSet_of(ctx *parser.Set_ofContext) {}

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
		expr := s.topExpr()
		expr.Type = type1
		if type1.GetSet() != nil {
			type1.SourceContext = nil
			type1 = type1.GetSet()
		}
		if type1.GetTypeRef() != nil {
			type1.GetTypeRef().Context.Path = nil
		}
	}
}

// EnterView_return_type is called when production view_return_type is entered.
func (s *TreeShapeListener) EnterView_return_type(ctx *parser.View_return_typeContext) {
	s.fieldname = append(s.fieldname, "view-return"+s.typename)
	s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{}
}

// ExitView_return_type is called when production view_return_type is exited.
func (s *TreeShapeListener) ExitView_return_type(ctx *parser.View_return_typeContext) {
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
	if type1.GetSet() != nil {
		type1.SourceContext = nil
		type1 = type1.GetSet()
	}
	if type1.GetTypeRef() != nil {
		type1.GetTypeRef().Context = nil
	}
}

// EnterTransform_scope_var is called when production transform_scope_var is entered.
func (s *TreeShapeListener) EnterTransform_scope_var(ctx *parser.Transform_scope_varContext) {
	tx := s.topExpr().GetTransform()
	tx.Scopevar = ctx.GetText()
}

// ExitTransform_scope_var is called when production transform_scope_var is exited.
func (s *TreeShapeListener) ExitTransform_scope_var(ctx *parser.Transform_scope_varContext) {}

// EnterTransform_arg is called when production transform_arg is entered.
func (s *TreeShapeListener) EnterTransform_arg(ctx *parser.Transform_argContext) {}

// ExitTransform_arg is called when production transform_arg is exited.
func (s *TreeShapeListener) ExitTransform_arg(ctx *parser.Transform_argContext) {
	expr := s.popExpr()
	tx := s.topExpr().GetTransform()
	tx.Arg = expr
}

// EnterTransform is called when production transform is entered.
func (s *TreeShapeListener) EnterTransform(ctx *parser.TransformContext) {
	tx := &sysl.Expr{
		Expr: &sysl.Expr_Transform_{
			Transform: &sysl.Expr_Transform{
				Stmt: []*sysl.Expr_Transform_Stmt{},
			},
		},
	}
	s.pushExpr(tx)
}

// ExitTransform is called when production transform is exited.
func (s *TreeShapeListener) ExitTransform(ctx *parser.TransformContext) {
	tx := s.topExpr().GetTransform()
	if tx.Arg == nil {
		tx.Arg = makeExprName(".")
	}
	if tx.Scopevar == "" {
		tx.Scopevar = "."
	}
}

// EnterExpr_block is called when production expr_block is entered.
func (s *TreeShapeListener) EnterExpr_block(ctx *parser.Expr_blockContext) {}

// ExitExpr_block is called when production expr_block is exited.
func (s *TreeShapeListener) ExitExpr_block(ctx *parser.Expr_blockContext) {}

// EnterView_param is called when production view_param is entered.
func (s *TreeShapeListener) EnterView_param(ctx *parser.View_paramContext) {
	s.fieldname = append(s.fieldname, ctx.Name_str().GetText())
	s.typemap[s.fieldname[len(s.fieldname)-1]] = &sysl.Type{}
}

// ExitView_param is called when production view_param is exited.
func (s *TreeShapeListener) ExitView_param(ctx *parser.View_paramContext) {
	view := s.module.Apps[s.appname].Views[s.typename]
	paramName := s.fieldname[len(s.fieldname)-1]
	type1 := s.typemap[paramName]

	if type1.GetSet() != nil && type1.GetSet().GetTypeRef() != nil {
		type1.GetSet().GetTypeRef().Context = nil
		type1.SourceContext = nil
	}
	if type1.GetTypeRef() != nil {
		type1.GetTypeRef().Context = nil
	}

	p := &sysl.Param{
		Name: paramName,
		Type: type1,
	}
	s.module.Apps[s.appname].Views[s.typename].Param = append(view.Param, p)
}

// EnterView_params is called when production view_params is entered.
func (s *TreeShapeListener) EnterView_params(ctx *parser.View_paramsContext) {}

// ExitView_params is called when production view_params is exited.
func (s *TreeShapeListener) ExitView_params(ctx *parser.View_paramsContext) {}

// EnterView is called when production view is entered.
func (s *TreeShapeListener) EnterView(ctx *parser.ViewContext) {
	viewName := ctx.Name_str().GetText()
	s.module.Apps[s.appname].Views[viewName] = &sysl.View{
		Param:         []*sysl.Param{},
		RetType:       &sysl.Type{},
		SourceContext: buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()),
	}
	if ctx.Attribs_or_modifiers() != nil {
		v := s.module.Apps[s.appname].Views[viewName]
		v.Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}

	s.typename = viewName
	s.fieldname = []string{}
	s.typemap = map[string]*sysl.Type{}
}

// ExitView is called when production view is exited.
func (s *TreeShapeListener) ExitView(ctx *parser.ViewContext) {
	view := s.module.Apps[s.appname].Views[s.typename]
	if ctx.Abstract_view() == nil {
		view.Expr = s.popExpr()
	} else {
		attributes := map[string]*sysl.Attribute{}
		patterns := []*sysl.Attribute{}
		patterns = append(patterns, &sysl.Attribute{
			Attribute: &sysl.Attribute_S{
				S: "abstract",
			},
		})

		attributes["patterns"] = &sysl.Attribute{
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: patterns,
				},
			},
		}
		view.Attrs = attributes
	}
	t1, has := s.typemap["view-return"+s.typename]
	if has {
		view.RetType = t1
	} else {
		view.RetType = view.Expr.Type
	}
	s.fieldname = []string{}
	s.typemap = map[string]*sysl.Type{}
	s.typename = ""
}

func buildSourceContext(filename string, line int, col int) *sysl.SourceContext {
	return &sysl.SourceContext{
		File: filename,
		Start: &sysl.SourceContext_Location{
			Line: int32(line),
			Col:  int32(col),
		},
	}
}

// EnterAlias is called when production alias is entered.
func (s *TreeShapeListener) EnterAlias(ctx *parser.AliasContext) {
	if s.typename == "" {
		s.typename = ctx.Name_str().GetText()
	} else {
		s.typename = s.typename + "." + ctx.Name_str().GetText()
	}
	type1 := &sysl.Type{}

	s.typemap = map[string]*sysl.Type{
		s.typename: type1,
	}
	s.fieldname = []string{s.typename}
	s.module.Apps[s.appname].Types[s.typename] = type1

	if ctx.Attribs_or_modifiers() != nil {
		type1.Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	}
	if ctx.Annotation(0) != nil {
		s.pushScope(type1)
	}
	type1.SourceContext = buildSourceContext(s.filename, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
}

// ExitAlias is called when production alias is exited.
func (s *TreeShapeListener) ExitAlias(ctx *parser.AliasContext) {
	s.module.Apps[s.appname].Types[s.typename] = s.typemap[s.fieldname[len(s.fieldname)-1]]

	s.typename = ""
	s.fieldname = []string{}
	s.typemap = map[string]*sysl.Type{}
	if ctx.Annotation(0) != nil {
		s.popScope()
	}
}

// EnterApp_decl is called when production app_decl is entered.
func (s *TreeShapeListener) EnterApp_decl(ctx *parser.App_declContext) {
	if s.module.Apps[s.appname].Types == nil && (ctx.Table(0) != nil || ctx.Alias(0) != nil) {
		s.module.Apps[s.appname].Types = map[string]*sysl.Type{}
	}
	has_stmts := (ctx.Simple_endpoint(0) != nil || ctx.Rest_endpoint(0) != nil || ctx.Event(0) != nil || ctx.Subscribe(0) != nil || ctx.Collector(0) != nil)
	if s.module.Apps[s.appname].Endpoints == nil && has_stmts {
		s.module.Apps[s.appname].Endpoints = map[string]*sysl.Endpoint{}
	}
	if s.module.Apps[s.appname].Wrapped == nil && len(ctx.AllFacade()) > 0 {
		s.module.Apps[s.appname].Wrapped = &sysl.Application{
			Types: map[string]*sysl.Type{},
		}
	}
	if s.module.Apps[s.appname].Views == nil && ctx.View(0) != nil {
		s.module.Apps[s.appname].Views = map[string]*sysl.View{}
	}
	if ctx.Annotation(0) != nil {
		if s.module.Apps[s.appname].Attrs == nil {
			s.module.Apps[s.appname].Attrs = map[string]*sysl.Attribute{}
		}
		s.pushScope(s.module.Apps[s.appname])
	}

	s.url_prefix = []string{""}
	s.rest_queryparams = []*sysl.Endpoint_RestParams_QueryParam{}
	s.rest_queryparams_len = []int{0}
	s.rest_attrs = []map[string]*sysl.Attribute{nil}
	s.typemap = map[string]*sysl.Type{}
}

// ExitApp_decl is called when production app_decl is exited.
func (s *TreeShapeListener) ExitApp_decl(ctx *parser.App_declContext) {
	if ctx.Annotation(0) != nil {
		s.popScope()
	}
	s.typename = ""
}

// EnterApplication is called when production application is entered.
func (s *TreeShapeListener) EnterApplication(ctx *parser.ApplicationContext) {}

// ExitApplication is called when production application is exited.
func (s *TreeShapeListener) ExitApplication(ctx *parser.ApplicationContext) {}

// EnterPath is called when production path is entered.
func (s *TreeShapeListener) EnterPath(ctx *parser.PathContext) {}

// ExitPath is called when production path is exited.
func (s *TreeShapeListener) ExitPath(ctx *parser.PathContext) {}

// EnterImport_stmt is called when production import_stmt is entered.
func (s *TreeShapeListener) EnterImport_stmt(ctx *parser.Import_stmtContext) {
	p := strings.Split(strings.TrimSpace(ctx.IMPORT().GetText()), " ")
	path := p[len(p)-1]
	if path[0] != '/' {
		path = filepath.ToSlash(s.base) + "/" + path
	}
	path += ".sysl"
	s.imports = append(s.imports, path)
}

// ExitImport_stmt is called when production import_stmt is exited.
func (s *TreeShapeListener) ExitImport_stmt(ctx *parser.Import_stmtContext) {}

// EnterImports_decl is called when production imports_decl is entered.
func (s *TreeShapeListener) EnterImports_decl(ctx *parser.Imports_declContext) {}

// ExitImports_decl is called when production imports_decl is exited.
func (s *TreeShapeListener) ExitImports_decl(ctx *parser.Imports_declContext) {}

// EnterSysl_file is called when production sysl_file is entered.
func (s *TreeShapeListener) EnterSysl_file(ctx *parser.Sysl_fileContext) {
}

// ExitSysl_file is called when production sysl_file is exited.
func (s *TreeShapeListener) ExitSysl_file(ctx *parser.Sysl_fileContext) {
	s.appname = ""
}
