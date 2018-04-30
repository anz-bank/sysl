package main // SyslParser

import (
	"fmt"
	"strconv"
	"strings"

	"anz-bank/sysl/src/proto"
	"anz-bank/sysl/sysl2/antlr/grammar"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var _ = fmt.Println

// TreeShapeListener ..
type TreeShapeListener struct {
	*parser.BaseSyslParserListener
	base                 string
	root                 string
	imports              []string
	module               *sysl.Module
	appname              string
	typename             string
	fieldname            []string
	url_prefix           []string
	app_name             []string
	annotation           string
	typemap              map[string]*sysl.Type
	prevLineEmpty        bool
	rest_queryparams     []*sysl.Endpoint_RestParams_QueryParam
	method_queryparams   []*sysl.Endpoint_RestParams_QueryParam
	rest_queryparams_len []int
	stmt_scope           []interface{} // Endpoint, if, if_else, loop
}

// NewTreeShapeListener ...
func NewTreeShapeListener(base, root string) *TreeShapeListener {
	return &TreeShapeListener{
		base: base,
		root: root,
		module: &sysl.Module{
			Apps: make(map[string]*sysl.Application),
		},
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
func (s *TreeShapeListener) EnterReference(ctx *parser.ReferenceContext) {}

// ExitReference is called when production reference is exited.
func (s *TreeShapeListener) ExitReference(ctx *parser.ReferenceContext) {}

func lastChar(str string) string {
	return str[len(str)-2:]
}

// EnterDoc_string is called when production doc_string is entered.
func (s *TreeShapeListener) EnterDoc_string(ctx *parser.Doc_stringContext) {
	if s.typemap == nil {
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
	if s.prevLineEmpty && len(ss) > 2 && lastChar(ss) == "\n\n" {
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

// EnterSet_type is called when production set_type is entered.
func (s *TreeShapeListener) EnterSet_type(ctx *parser.Set_typeContext) {}

// ExitSet_type is called when production set_type is exited.
func (s *TreeShapeListener) ExitSet_type(ctx *parser.Set_typeContext) {}

// EnterCollection_type is called when production collection_type is entered.
func (s *TreeShapeListener) EnterCollection_type(ctx *parser.Collection_typeContext) {}

// ExitCollection_type is called when production collection_type is exited.
func (s *TreeShapeListener) ExitCollection_type(ctx *parser.Collection_typeContext) {}

// EnterUser_defined_type is called when production user_defined_type is entered.
func (s *TreeShapeListener) EnterUser_defined_type(ctx *parser.User_defined_typeContext) {}

// ExitUser_defined_type is called when production user_defined_type is exited.
func (s *TreeShapeListener) ExitUser_defined_type(ctx *parser.User_defined_typeContext) {}

// EnterMulti_line_docstring is called when production multi_line_docstring is entered.
func (s *TreeShapeListener) EnterMulti_line_docstring(ctx *parser.Multi_line_docstringContext) {}

// ExitMulti_line_docstring is called when production multi_line_docstring is exited.
func (s *TreeShapeListener) ExitMulti_line_docstring(ctx *parser.Multi_line_docstringContext) {}

// EnterAnnotation_value is called when production annotation_value is entered.
func (s *TreeShapeListener) EnterAnnotation_value(ctx *parser.Annotation_valueContext) {
	attrs := s.peekAttrs()

	if ctx.QSTRING() != nil {
		attrs[s.annotation].Attribute = &sysl.Attribute_S{
			S: ctx.QSTRING().GetText(),
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
	size_spec, has_size_spec := ctx.Size_spec().(*parser.Size_specContext)
	array_spec, has_array_spec := ctx.Array_size().(*parser.Array_sizeContext)

	native := ctx.NativeDataTypes()
	type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
	// fmt.Printf("enter field type:")
	// fmt.Println(s.fieldname[len(s.fieldname)-1])

	if native != nil {
		primitive_type := sysl.Type_Primitive(sysl.Type_Primitive_value[strings.ToUpper(native.GetText())])
		type1.Type = &sysl.Type_Primitive_{
			Primitive: primitive_type,
		}
		if has_size_spec {
			type1.Constraint = makeTypeConstraint(primitive_type, size_spec)
		} else if has_array_spec {
			type1.Constraint = makeArrayConstraint(primitive_type, array_spec)
		}
	} else if ctx.Reference() != nil {
		// refContext := ctx.Reference().(*parser.ReferenceContext)
		context_app_part := []string{s.appname}
		context_path := []string{s.typename}

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
	} else if ctx.User_defined_type() != nil {
		ctxt := ctx.User_defined_type().(*parser.User_defined_typeContext)
		context_app_part := []string{s.appname}
		context_path := []string{s.typename}
		ref_path := []string{
			ctxt.GetText(),
		}

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

	} else if ctx.Collection_type() != nil {
		ctxt := ctx.Collection_type().(*parser.Collection_typeContext)
		setCtxt := ctxt.Set_type().(*parser.Set_typeContext)
		native := setCtxt.NativeDataTypes()

		var contained_type *sysl.Type

		if native != nil {
			primitive_type := sysl.Type_Primitive(sysl.Type_Primitive_value[strings.ToUpper(native.GetText())])
			contained_type = &sysl.Type{
				Type: &sysl.Type_Primitive_{
					Primitive: primitive_type,
				},
			}
		} else {
			context_app_part := []string{s.appname}
			context_path := []string{s.typename}
			ref_path := []string{
				setCtxt.Name().GetText(),
			}

			contained_type = &sysl.Type{
				Type: &sysl.Type_TypeRef{
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
				},
			}
		}

		contained_type.SourceContext = &sysl.SourceContext{
			Start: &sysl.SourceContext_Location{
				Line: int32(setCtxt.GetStart().GetLine()),
			},
		}

		type1.Type = &sysl.Type_Set{
			Set: contained_type,
		}
	}
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = makeAttributeArray(attribs)
	}

	if ctx.QN() != nil {
		type1.Opt = true
	}

	if ctx.Annotations() != nil {
		if type1.Attrs == nil {
			type1.Attrs = make(map[string]*sysl.Attribute)
		}
		s.pushScope(type1)
	}
}

// ExitField_type is called when production field_type is exited.
func (s *TreeShapeListener) ExitField_type(ctx *parser.Field_typeContext) {
	// fmt.Printf("exit field type")
	// fmt.Println(s.fieldname[len(s.fieldname)-1])
	if ctx.Reference() != nil {
		type1 := s.typemap[s.fieldname[len(s.fieldname)-1]]
		l := len(s.app_name)
		ref_path := s.app_name[l-2:]
		appname := s.app_name[:l-2]
		type1.GetTypeRef().Ref = &sysl.Scope{
			Path: ref_path,
		}
		if len(appname) > 0 {
			// added to match legacy
			appname[0] = appname[0] + " "
			for i := 1; i < len(appname); i++ {
				appname[1] = " " + appname[1] + " "
			}
			// end hack
			type1.GetTypeRef().Ref.Appname = &sysl.AppName{
				Part: appname,
			}
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
	case sysl.Type_DATE:
		fallthrough
	case sysl.Type_DATETIME:
		fallthrough
	case sysl.Type_INT:
		fallthrough
	case sysl.Type_STRING:
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
				Precision: int32(l),
			})
			if size_spec.DIGITS(1) != nil {
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
	case sysl.Type_DATE:
		fallthrough
	case sysl.Type_DATETIME:
		fallthrough
	case sysl.Type_INT:
		fallthrough
	case sysl.Type_STRING:
		val1 := array_size.DIGITS(1).GetText()
		if l, err = strconv.Atoi(val1); err == nil {
			c = append(c, &sysl.Type_Constraint{
				Length: &sysl.Type_Constraint_Length{
					Max: int64(l),
				},
			})
		}
	case sysl.Type_DECIMAL:
		val0 := array_size.DIGITS(0).GetText()
		val1 := array_size.DIGITS(1).GetText()
		l0, _ := strconv.Atoi(val0)
		l1, _ := strconv.Atoi(val1)
		ct := &sysl.Type_Constraint{
			Length: &sysl.Type_Constraint_Length{
				Max: int64(l1),
			},
		}
		if l0 != 0 {
			ct.GetLength().Min = int64(l0)
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
	arr := make([]*sysl.Attribute, 0)
	for _, ars := range array_strings.AllQuoted_string() {
		str := ars.(*parser.Quoted_stringContext)
		arr = append(arr, &sysl.Attribute{
			Attribute: &sysl.Attribute_S{
				S: str.QSTRING().GetText(),
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
	patterns := make([]*sysl.Attribute, 0)
	attributes := make(map[string]*sysl.Attribute)

	for _, e := range attribs.AllEntry() {
		entry := e.(*parser.EntryContext)
		if entry.Nvp() != nil {
			nvp := entry.Nvp().(*parser.NvpContext)

			if nvp.Quoted_string() != nil {
				qs := nvp.Quoted_string().(*parser.Quoted_stringContext)
				// fmt.Printf("attrib: %s %s\n", nvp.Name().GetText(), qs.QSTRING().GetText())
				attributes[nvp.Name().GetText()] = &sysl.Attribute{
					Attribute: &sysl.Attribute_S{
						S: qs.QSTRING().GetText(),
					},
				}
			} else if nvp.Array_of_strings() != nil {
				array_strings := nvp.Array_of_strings().(*parser.Array_of_stringsContext)
				attributes[nvp.Name().GetText()] = makeArrayOfStringsAttribute(array_strings)
			} else {
				panic("array of arrays: not handled yet")
			}
		} else if entry.Modifier() != nil {
			mod := entry.Modifier().(*parser.ModifierContext)
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

func search(attr string, attrs []*sysl.Attribute) bool {
	for _, a := range attrs {
		if s, ok := a.Attribute.(*sysl.Attribute_S); ok {
			if s.S == attr {
				return true
			}
		}
	}
	return false
}

// EnterInplace_tuple is called when production inplace_tuple is entered.
func (s *TreeShapeListener) EnterInplace_tuple(ctx *parser.Inplace_tupleContext) {
	s.typename = s.typename + "." + s.fieldname[len(s.fieldname)-1]
	s.typemap = make(map[string]*sysl.Type)
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
	l := strings.LastIndex(s.typename, ".")
	s.typename = s.typename[:l]
	s.typemap = s.module.Apps[s.appname].Types[s.typename].GetTuple().GetAttrDefs()
}

// EnterField is called when production field is entered.
func (s *TreeShapeListener) EnterField(ctx *parser.FieldContext) {
	s.fieldname = append(s.fieldname, ctx.Name_str().GetText())
	// fmt.Printf("%d) %s %s\n", len(s.fieldname), s.typename, ctx.Name_str().GetText())
	type1 := &sysl.Type{}

	type1.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}
	if ctx.Inplace_tuple() != nil {
		type1.Type = &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Ref: &sysl.Scope{
					Path: []string{s.fieldname[len(s.fieldname)-1]},
				},
			},
		}
		if ctx.Array_size() != nil {
			type1.Type = &sysl.Type_List_{
				List: &sysl.Type_List{
					Type: &sysl.Type{
						Type:          type1.Type,
						SourceContext: type1.SourceContext,
					},
				},
			}
			type1.SourceContext = nil
		}
	}
	s.typemap[s.fieldname[len(s.fieldname)-1]] = type1
	s.app_name = make([]string, 0)
}

// ExitField is called when production field is exited.
func (s *TreeShapeListener) ExitField(ctx *parser.FieldContext) {}

// EnterTable is called when production table is entered.
func (s *TreeShapeListener) EnterTable(ctx *parser.TableContext) {
	s.typename = ctx.Name_str().GetText()
	s.typemap = make(map[string]*sysl.Type)

	if ctx.TABLE() != nil {
		if s.module.Apps[s.appname].Types[s.typename].GetRelation().GetAttrDefs() != nil {
			panic("not implemented yet")
		}

		s.module.Apps[s.appname].Types[s.typename] = &sysl.Type{
			Type: &sysl.Type_Relation_{
				Relation: &sysl.Type_Relation{
					AttrDefs: s.typemap,
				},
			},
		}
	}
	if ctx.TYPE() != nil {
		s.module.Apps[s.appname].Types[s.typename] = &sysl.Type{
			Type: &sysl.Type_Tuple_{
				Tuple: &sysl.Type_Tuple{
					AttrDefs: s.typemap,
				},
			},
		}
	}
	type1 := s.module.Apps[s.appname].Types[s.typename]
	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = makeAttributeArray(attribs)
	}
	if ctx.Annotation(0) != nil {
		if type1.Attrs == nil {
			type1.Attrs = make(map[string]*sysl.Attribute)
		}
		s.pushScope(type1)
	}
}

// ExitTable is called when production table is exited.
func (s *TreeShapeListener) ExitTable(ctx *parser.TableContext) {
	// wire up primary key
	if rel := s.module.Apps[s.appname].Types[s.typename].GetRelation(); rel != nil {
		pks := make([]string, 0)
		for _, name := range s.fieldname {
			f := rel.GetAttrDefs()[name]
			patterns, has := f.GetAttrs()["patterns"]
			if has {
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

	if ctx.Annotation(0) != nil {
		s.popScope()
	}

	s.typename = ""
	s.fieldname = []string{}
	s.typemap = nil
}

// EnterPackage_name is called when production package_name is entered.
func (s *TreeShapeListener) EnterPackage_name(ctx *parser.Package_nameContext) {

}

// ExitPackage_name is called when production package_name is exited.
func (s *TreeShapeListener) ExitPackage_name(ctx *parser.Package_nameContext) {
	// fmt.Printf("package_name: ")
	// fmt.Println(s.app_name)
}

// EnterSub_package is called when production sub_package is entered.
func (s *TreeShapeListener) EnterSub_package(ctx *parser.Sub_packageContext) {}

// ExitSub_package is called when production sub_package is exited.
func (s *TreeShapeListener) ExitSub_package(ctx *parser.Sub_packageContext) {}

// EnterApp_name is called when production app_name is entered.
func (s *TreeShapeListener) EnterApp_name(ctx *parser.App_nameContext) {
	s.app_name = make([]string, 0)
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
		s.module.Apps[s.appname].LongName = ctx.QSTRING().GetText()
	}

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		if s.module.Apps[s.appname].Attrs == nil {
			s.module.Apps[s.appname].Attrs = makeAttributeArray(attribs)
		}
	}
}

// ExitName_with_attribs is called when production name_with_attribs is exited.
func (s *TreeShapeListener) ExitName_with_attribs(ctx *parser.Name_with_attribsContext) {
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
	context_app_part := []string{s.appname}
	var ref_path []string

	if ctx.Var_in_curly() != nil {
		ref_path = append(ref_path, ctx.Var_in_curly().GetText())
		type1 = &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: context_app_part,
						},
					},
					Ref: &sysl.Scope{
						Path: ref_path,
					},
				},
			},
		}
	} else {
		type_str := strings.ToUpper(ctx.NativeDataTypes().GetText())
		primitive_type := sysl.Type_Primitive(sysl.Type_Primitive_value[type_str])
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: primitive_type,
			},
		}
	}

	rest_param := &sysl.Endpoint_RestParams_QueryParam{
		Name: var_name,
		Type: type1,
		Loc:  true,
	}

	if ctx.QN() != nil {
		rest_param.Type.Opt = true
	}

	rest_param.Type.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}
	s.method_queryparams = append(s.method_queryparams, rest_param)

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
	} else {
		context_app_part := []string{s.appname}
		ref_path := []string{
			ctx.Name_str().GetText(),
		}

		type1 = &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Context: &sysl.Scope{
						Appname: &sysl.AppName{
							Part: context_app_part,
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
		Loc:  true,
	}

	rest_param.Type.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}

	s.rest_queryparams = append(s.rest_queryparams, rest_param)
	s.typename += ctx.CURLY_OPEN().GetText() + var_name + ctx.CURLY_CLOSE().GetText()
}

// ExitHttp_path_var_with_type is called when production http_path_var_with_type is exited.
func (s *TreeShapeListener) ExitHttp_path_var_with_type(ctx *parser.Http_path_var_with_typeContext) {}

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
	s.lastStatement().GetCall().Target.Part = s.app_name
}

// EnterTarget_endpoint is called when production target_endpoint is entered.
func (s *TreeShapeListener) EnterTarget_endpoint(ctx *parser.Target_endpointContext) {}

// ExitTarget_endpoint is called when production target_endpoint is exited.
func (s *TreeShapeListener) ExitTarget_endpoint(ctx *parser.Target_endpointContext) {}

// EnterCall_arg is called when production call_arg is entered.
func (s *TreeShapeListener) EnterCall_arg(ctx *parser.Call_argContext) {
	fmt.Println(ctx.GetText())
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
	fmt.Println(ctx.Target_endpoint().GetText())
	appName := &sysl.AppName{}
	if ctx.DOT_ARROW() != nil {
		appName.Part = []string{s.appname}
	}
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Call{
			Call: &sysl.Call{
				Target:   appName,
				Endpoint: strings.TrimSpace(ctx.Target_endpoint().GetText()),
			},
		},
	})
	if ctx.Call_args() != nil {
		s.lastStatement().GetCall().Arg = make([]*sysl.Call_Arg, 0)
	}
}

// ExitCall_stmt is called when production call_stmt is exited.
func (s *TreeShapeListener) ExitCall_stmt(ctx *parser.Call_stmtContext) {}

// EnterIf_stmt is called when production if_stmt is entered.
func (s *TreeShapeListener) EnterIf_stmt(ctx *parser.If_stmtContext) {
}

// ExitIf_stmt is called when production if_stmt is exited.
func (s *TreeShapeListener) ExitIf_stmt(ctx *parser.If_stmtContext) {
	s.popScope()
}

// EnterIf_else is called when production if_else is entered.
func (s *TreeShapeListener) EnterIf_else(ctx *parser.If_elseContext) {
	ifstmt := ctx.If_stmt().(*parser.If_stmtContext)

	if_stmt := &sysl.Statement{
		Stmt: &sysl.Statement_Cond{
			Cond: &sysl.Cond{
				Test: ifstmt.Arg_value().GetText(),
				Stmt: make([]*sysl.Statement, 0),
			},
		},
	}
	s.addToCurrentScope(if_stmt)

	// else statements
	if ctx.Statements(0) != nil {
		else_stmt := &sysl.Statement{
			Stmt: &sysl.Statement_Group{
				Group: &sysl.Group{
					Title: ctx.ELSE().GetText(),
					Stmt:  make([]*sysl.Statement, 0),
				},
			},
		}
		s.addToCurrentScope(else_stmt)
		s.pushScope(else_stmt.GetGroup())
	}
	// if stmt is on top
	s.pushScope(if_stmt.GetCond())
}

// ExitIf_else is called when production if_else is exited.
func (s *TreeShapeListener) ExitIf_else(ctx *parser.If_elseContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
}

// EnterFor_cond is called when production for_cond is entered.
func (s *TreeShapeListener) EnterFor_cond(ctx *parser.For_condContext) {}

// ExitFor_cond is called when production for_cond is exited.
func (s *TreeShapeListener) ExitFor_cond(ctx *parser.For_condContext) {}

// EnterFor_stmt is called when production for_stmt is entered.
func (s *TreeShapeListener) EnterFor_stmt(ctx *parser.For_stmtContext) {
	for_stmt := &sysl.Statement{
		Stmt: &sysl.Statement_Group{
			Group: &sysl.Group{
				Title: ctx.FOR().GetText() + ctx.For_cond().GetText(),
				Stmt:  make([]*sysl.Statement, 0),
			},
		},
	}
	s.addToCurrentScope(for_stmt)
	s.pushScope(for_stmt.GetGroup())

}

// ExitFor_stmt is called when production for_stmt is exited.
func (s *TreeShapeListener) ExitFor_stmt(ctx *parser.For_stmtContext) {
	s.popScope()
}

// EnterHttp_method_comment is called when production http_method_comment is entered.
func (s *TreeShapeListener) EnterHttp_method_comment(ctx *parser.Http_method_commentContext) {}

// ExitHttp_method_comment is called when production http_method_comment is exited.
func (s *TreeShapeListener) ExitHttp_method_comment(ctx *parser.Http_method_commentContext) {}

// EnterOne_of_case_label is called when production one_of_case_label is entered.
func (s *TreeShapeListener) EnterOne_of_case_label(ctx *parser.One_of_case_labelContext) {}

// ExitOne_of_case_label is called when production one_of_case_label is exited.
func (s *TreeShapeListener) ExitOne_of_case_label(ctx *parser.One_of_case_labelContext) {}

// EnterOne_of_cases is called when production one_of_cases is entered.
func (s *TreeShapeListener) EnterOne_of_cases(ctx *parser.One_of_casesContext) {
	alt := s.peekScope().(*sysl.Alt)
	choice := &sysl.Alt_Choice{
		Stmt: make([]*sysl.Statement, 0),
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
			Choice: make([]*sysl.Alt_Choice, 0),
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

// EnterText_stmt is called when production text_stmt is entered.
func (s *TreeShapeListener) EnterText_stmt(ctx *parser.Text_stmtContext) {
	s.addToCurrentScope(&sysl.Statement{
		Stmt: &sysl.Statement_Action{
			Action: &sysl.Action{
				Action: ctx.GetText(),
			},
		},
	})
}

// ExitText_stmt is called when production text_stmt is exited.
func (s *TreeShapeListener) ExitText_stmt(ctx *parser.Text_stmtContext) {}

// EnterMixin is called when production mixin is entered.
func (s *TreeShapeListener) EnterMixin(ctx *parser.MixinContext) {
	if s.module.Apps[s.appname].Mixin2 == nil {
		s.module.Apps[s.appname].Mixin2 = make([]*sysl.Application, 0)
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

// EnterParam_list is called when production param_list is entered.
func (s *TreeShapeListener) EnterParam_list(ctx *parser.Param_listContext) {}

// ExitParam_list is called when production param_list is exited.
func (s *TreeShapeListener) ExitParam_list(ctx *parser.Param_listContext) {}

// EnterParams is called when production params is entered.
func (s *TreeShapeListener) EnterParams(ctx *parser.ParamsContext) {
	s.typemap = make(map[string]*sysl.Type)
	// fmt.Println("enter param")
}

// ExitParams is called when production params is exited.
func (s *TreeShapeListener) ExitParams(ctx *parser.ParamsContext) {
	params := make([]*sysl.Param, 0)

	for _, fieldname := range s.fieldname {
		// fieldname := s.fieldname[i]
		fmt.Printf("%s.%s: enter param \n", s.typename, fieldname)

		type1 := s.typemap[fieldname]
		if type1.GetTypeRef() != nil {
			type1.GetTypeRef().Context = nil
			ref := type1.GetTypeRef().GetRef()
			if ref.Appname == nil {
				ref.Appname = &sysl.AppName{
					Part: []string{},
				}
			}
			if len(ref.Path) == 1 {
				ref.Appname.Part = ref.Path
				ref.Path = nil
			} else {
				ref.Appname.Part = append(ref.Appname.Part, ref.Path[:len(ref.Path)-1]...)
				ref.Path = ref.Path[len(ref.Path)-1:]
				// Hack, to match with legacy
				// remove the extra space added earlier.
				for i := range ref.Path {
					ref.Path[i] = strings.TrimSpace(ref.Path[i])
				}
				for i := range ref.Appname.Part {
					ref.Appname.Part[i] = strings.TrimSpace(ref.Appname.Part[i])
				}
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

	s.module.Apps[s.appname].Endpoints[s.typename].Param = params
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
		return scope.Stmt[l]
	case *sysl.Cond:
		l := len(scope.Stmt) - 1
		return scope.Stmt[l]
	case *sysl.Alt_Choice:
		l := len(scope.Stmt) - 1
		return scope.Stmt[l]
	case *sysl.Group:
		l := len(scope.Stmt) - 1
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
	default:
		fmt.Printf("got unexpected %T\n", scope)
		panic("not implemented")
	}
}

// EnterMethod_def is called when production method_def is entered.
func (s *TreeShapeListener) EnterMethod_def(ctx *parser.Method_defContext) {
	url := s.urlString()
	s.typename = ctx.HTTP_VERBS().GetText() + " " + url
	s.method_queryparams = make([]*sysl.Endpoint_RestParams_QueryParam, 0)

	s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
		Name: s.typename,
		RestParams: &sysl.Endpoint_RestParams{
			Method: sysl.Endpoint_RestParams_Method(sysl.Endpoint_RestParams_Method_value[ctx.HTTP_VERBS().GetText()]),
			Path:   url,
		},
		Stmt: make([]*sysl.Statement, 0),
	}
	s.pushScope(s.module.Apps[s.appname].Endpoints[s.typename])

	if ctx.Attribs_or_modifiers() != nil {
		s.module.Apps[s.appname].Endpoints[s.typename].Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
	} else {
		s.module.Apps[s.appname].Endpoints[s.typename].Attrs = make(map[string]*sysl.Attribute)
	}

	attrs := s.module.Apps[s.appname].Endpoints[s.typename].Attrs
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

	if len(s.rest_queryparams) > 0 {
		qparams := make([]*sysl.Endpoint_RestParams_QueryParam, 0)
		for i := len(s.rest_queryparams) - 1; i >= 0; i-- {
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
				Loc: true,
			}
			qparams = append(qparams, qcopy)
		}
		s.module.Apps[s.appname].Endpoints[s.typename].RestParams.QueryParam = qparams
	}
}

// ExitMethod_def is called when production method_def is exited.
func (s *TreeShapeListener) ExitMethod_def(ctx *parser.Method_defContext) {
	if len(s.method_queryparams) > 0 {
		qparams := s.module.Apps[s.appname].Endpoints[s.typename].RestParams.QueryParam
		if qparams == nil {
			qparams = make([]*sysl.Endpoint_RestParams_QueryParam, 0)
		}
		qparams = append(qparams, s.method_queryparams...)
		s.module.Apps[s.appname].Endpoints[s.typename].RestParams.QueryParam = qparams
	}
	// HACK to match legacy
	if len(s.module.Apps[s.appname].Endpoints[s.typename].Param) > 0 {
		s.module.Apps[s.appname].Endpoints[s.typename].Param = []*sysl.Param{&sysl.Param{}}
	}
	s.popScope()
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
	if ctx.Endpoint_name() != nil {
		s.typename = ctx.Endpoint_name().GetText()
		s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
			Name: s.typename,
			Stmt: make([]*sysl.Statement, 0),
		}
		if ctx.Attribs_or_modifiers() != nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
	}
	if ctx.Statements(0) != nil {
		if s.module.Apps[s.appname].Endpoints[s.typename].Attrs == nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Attrs = make(map[string]*sysl.Attribute)
		}
		s.pushScope(s.module.Apps[s.appname].Endpoints[s.typename])
	}
}

// ExitSimple_endpoint is called when production simple_endpoint is exited.
func (s *TreeShapeListener) ExitSimple_endpoint(ctx *parser.Simple_endpointContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
}

// EnterRest_endpoint is called when production rest_endpoint is entered.
func (s *TreeShapeListener) EnterRest_endpoint(ctx *parser.Rest_endpointContext) {
	s.rest_queryparams_len = append(s.rest_queryparams_len, len(s.rest_queryparams))
}

// ExitRest_endpoint is called when production rest_endpoint is exited.
func (s *TreeShapeListener) ExitRest_endpoint(ctx *parser.Rest_endpointContext) {
	s.url_prefix = s.url_prefix[:len(s.url_prefix)-1]
	ltop := len(s.rest_queryparams_len) - 1
	if ltop >= 0 {
		l := s.rest_queryparams_len[ltop]
		s.rest_queryparams = s.rest_queryparams[:l]
		s.rest_queryparams_len = s.rest_queryparams_len[:ltop]
	}
}

// EnterCollector_stmt is called when production collector_stmt is entered.
func (s *TreeShapeListener) EnterCollector_stmt(ctx *parser.Collector_stmtContext) {}

// ExitCollector_stmt is called when production collector_stmt is exited.
func (s *TreeShapeListener) ExitCollector_stmt(ctx *parser.Collector_stmtContext) {}

// EnterCollector_stmts is called when production collector_stmts is entered.
func (s *TreeShapeListener) EnterCollector_stmts(ctx *parser.Collector_stmtsContext) {}

// ExitCollector_stmts is called when production collector_stmts is exited.
func (s *TreeShapeListener) ExitCollector_stmts(ctx *parser.Collector_stmtsContext) {}

// EnterCollector is called when production collector is entered.
func (s *TreeShapeListener) EnterCollector(ctx *parser.CollectorContext) {}

// ExitCollector is called when production collector is exited.
func (s *TreeShapeListener) ExitCollector(ctx *parser.CollectorContext) {}

// EnterEvent is called when production event is entered.
func (s *TreeShapeListener) EnterEvent(ctx *parser.EventContext) {
	if ctx.EVENT_NAME() != nil {
		s.typename = ctx.EVENT_NAME().GetText()
		s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
			Name:     s.typename,
			Stmt:     make([]*sysl.Statement, 0),
			IsPubsub: true,
		}
		if ctx.Attribs_or_modifiers() != nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
	}
}

// ExitEvent is called when production event is exited.
func (s *TreeShapeListener) ExitEvent(ctx *parser.EventContext) {}

// EnterSubscribe is called when production subscribe is entered.
func (s *TreeShapeListener) EnterSubscribe(ctx *parser.SubscribeContext) {
	if ctx.App_name() != nil {
		s.typename = ctx.App_name().GetText() + ctx.ARROW_RIGHT().GetText() + ctx.Name_str().GetText()
		app_src := &sysl.AppName{
			Part: []string{ctx.App_name().GetText()},
		}
		s.module.Apps[s.appname].Endpoints[s.typename] = &sysl.Endpoint{
			Name:   s.typename,
			Stmt:   make([]*sysl.Statement, 0),
			Source: app_src,
		}
		if ctx.Attribs_or_modifiers() != nil {
			s.module.Apps[s.appname].Endpoints[s.typename].Attrs = makeAttributeArray(ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext))
		}
		if ctx.Statements(0) != nil {
			if s.module.Apps[s.appname].Endpoints[s.typename].Attrs == nil {
				s.module.Apps[s.appname].Endpoints[s.typename].Attrs = make(map[string]*sysl.Attribute)
			}
			s.pushScope(s.module.Apps[s.appname].Endpoints[s.typename])
		}
	}
}

// ExitSubscribe is called when production subscribe is exited.
func (s *TreeShapeListener) ExitSubscribe(ctx *parser.SubscribeContext) {
	if ctx.Statements(0) != nil {
		s.popScope()
	}
}

// EnterApp_decl is called when production app_decl is entered.
func (s *TreeShapeListener) EnterApp_decl(ctx *parser.App_declContext) {
	if s.module.Apps[s.appname].Types == nil && len(ctx.AllTable()) > 0 {
		s.module.Apps[s.appname].Types = make(map[string]*sysl.Type)
	}
	if s.module.Apps[s.appname].Endpoints == nil && (ctx.Simple_endpoint(0) != nil || ctx.Rest_endpoint(0) != nil || ctx.Event(0) != nil || ctx.Subscribe(0) != nil) {
		s.module.Apps[s.appname].Endpoints = make(map[string]*sysl.Endpoint)
		s.url_prefix = []string{""}
		s.rest_queryparams = make([]*sysl.Endpoint_RestParams_QueryParam, 0)
		s.rest_queryparams_len = []int{0}
	}
	if s.module.Apps[s.appname].Wrapped == nil && len(ctx.AllFacade()) > 0 {
		s.module.Apps[s.appname].Wrapped = &sysl.Application{
			Types: make(map[string]*sysl.Type),
		}
	}
	if ctx.Annotation(0) != nil {
		if s.module.Apps[s.appname].Attrs == nil {
			s.module.Apps[s.appname].Attrs = make(map[string]*sysl.Attribute, 0)
		}
		s.pushScope(s.module.Apps[s.appname])
	}
}

// ExitApp_decl is called when production app_decl is exited.
func (s *TreeShapeListener) ExitApp_decl(ctx *parser.App_declContext) {
	if ctx.Annotation(0) != nil {
		s.popScope()
	}
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
	if path[0] == '/' {
		path = s.root + path
	} else {
		path = s.base + "/" + path
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
