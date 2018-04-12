package main // SyslParser

import (
	"strconv"
	"strings"

	"anz-bank/sysl/src/proto"
	"anz-bank/sysl/sysl2/antlr/grammar"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// TreeShapeListener ..
type TreeShapeListener struct {
	*parser.BaseSyslParserListener
	root      string
	module    *sysl.Module
	appname   string
	typename  string
	fieldname string
	typemap   map[string]*sysl.Type
}

// NewTreeShapeListener ...
func NewTreeShapeListener(root string) *TreeShapeListener {
	return &TreeShapeListener{root: root}
}

// VisitErrorNode is called when an error node is visited.
func (s *TreeShapeListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *TreeShapeListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

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

// EnterDoc_string is called when production doc_string is entered.
func (s *TreeShapeListener) EnterDoc_string(ctx *parser.Doc_stringContext) {}

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
func (s *TreeShapeListener) EnterAnnotation_value(ctx *parser.Annotation_valueContext) {}

// ExitAnnotation_value is called when production annotation_value is exited.
func (s *TreeShapeListener) ExitAnnotation_value(ctx *parser.Annotation_valueContext) {}

// EnterAnnotation is called when production annotation is entered.
func (s *TreeShapeListener) EnterAnnotation(ctx *parser.AnnotationContext) {}

// ExitAnnotation is called when production annotation is exited.
func (s *TreeShapeListener) ExitAnnotation(ctx *parser.AnnotationContext) {}

// EnterAnnotations is called when production annotations is entered.
func (s *TreeShapeListener) EnterAnnotations(ctx *parser.AnnotationsContext) {}

// ExitAnnotations is called when production annotations is exited.
func (s *TreeShapeListener) ExitAnnotations(ctx *parser.AnnotationsContext) {}

// EnterField_type is called when production field_type is entered.
func (s *TreeShapeListener) EnterField_type(ctx *parser.Field_typeContext) {}

// ExitField_type is called when production field_type is exited.
func (s *TreeShapeListener) ExitField_type(ctx *parser.Field_typeContext) {}

func mapNativeDataTypeToType_Primitive(val string) sysl.Type_Primitive {

	switch sysl.Type_Primitive_value[strings.ToUpper(val)] {
	case 3:
		return sysl.Type_BOOL
	case 4:
		return sysl.Type_INT
	case 5:
		return sysl.Type_FLOAT
	case 6:
		return sysl.Type_STRING
	case 12:
		return sysl.Type_DECIMAL
	case 9:
		return sysl.Type_DATE
	case 10:
		return sysl.Type_DATETIME
	default:
		// panic("handle other cases too!")
	}
	return sysl.Type_NO_Primitive
}

func makeTypeConstraint(t sysl.Type_Primitive, size_spec *parser.Size_specContext) []*sysl.Type_Constraint {
	c := []*sysl.Type_Constraint{}
	var err error
	var l int

	switch t {
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

func makeAttributeArray(attribs *parser.Attribs_or_modifiersContext) []*sysl.Attribute {
	elt := make([]*sysl.Attribute, 0)

	for _, e := range attribs.AllEntry() {
		entry := e.(*parser.EntryContext)
		if entry.Nvp() != nil {
			panic("got nvp")
		} else if entry.Modifier() != nil {
			mod := entry.Modifier().(*parser.ModifierContext)
			elt = append(elt, &sysl.Attribute{
				Attribute: &sysl.Attribute_S{
					S: mod.GetText()[1:],
				},
			})
		}
	}
	return elt
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

// EnterField is called when production field is entered.
func (s *TreeShapeListener) EnterField(ctx *parser.FieldContext) {
	s.fieldname = ctx.Name().GetText()

	field_type := ctx.Field_type().(*parser.Field_typeContext)
	size_spec, has_size_spec := field_type.Size_spec().(*parser.Size_specContext)
	native := field_type.NativeDataTypes()

	var type1 *sysl.Type

	if native != nil {
		primitive_type := mapNativeDataTypeToType_Primitive(native.GetText())
		type1 = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: primitive_type,
			},
		}
		if has_size_spec {
			type1.Constraint = makeTypeConstraint(primitive_type, size_spec)
		}
	}
	if field_type.Reference() != nil {
		refContext := field_type.Reference().(*parser.ReferenceContext)
		context_app_part := []string{s.appname}
		context_path := []string{s.typename}
		ref_path := []string{
			refContext.GetParent_ref().GetText(),
			refContext.GetMember().GetText(),
		}

		type1 = &sysl.Type{
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
	if attribs, ok := field_type.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		type1.Attrs = make(map[string]*sysl.Attribute)
		attrs := makeAttributeArray(attribs)
		type1.Attrs["patterns"] = &sysl.Attribute{
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: attrs,
				},
			},
		}
	}
	type1.SourceContext = &sysl.SourceContext{
		Start: &sysl.SourceContext_Location{
			Line: int32(ctx.GetStart().GetLine()),
		},
	}
	s.typemap[s.fieldname] = type1
}

// ExitField is called when production field is exited.
func (s *TreeShapeListener) ExitField(ctx *parser.FieldContext) {}

// EnterTable is called when production table is entered.
func (s *TreeShapeListener) EnterTable(ctx *parser.TableContext) {
	s.typemap = make(map[string]*sysl.Type)
	s.typename = ctx.Name().GetText()

	// fmt.Printf("Enter Table: %s %s\n", s.appname, s.typename)

	if ctx.TABLE() != nil {
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
}

// ExitTable is called when production table is exited.
func (s *TreeShapeListener) ExitTable(ctx *parser.TableContext) {

	rel := s.module.Apps[s.appname].Types[s.typename].GetRelation()
	if rel == nil {
		return
	}
	pks := make([]string, 0)
	for name, f := range rel.AttrDefs {
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
	s.typename = ""
}

// EnterPackage_name is called when production package_name is entered.
func (s *TreeShapeListener) EnterPackage_name(ctx *parser.Package_nameContext) {
	s.module.Apps[s.appname].Name.Part = append(s.module.Apps[s.appname].Name.Part, ctx.GetText())
}

// ExitPackage_name is called when production package_name is exited.
func (s *TreeShapeListener) ExitPackage_name(ctx *parser.Package_nameContext) {}

// EnterSub_package is called when production sub_package is entered.
func (s *TreeShapeListener) EnterSub_package(ctx *parser.Sub_packageContext) {}

// ExitSub_package is called when production sub_package is exited.
func (s *TreeShapeListener) ExitSub_package(ctx *parser.Sub_packageContext) {}

// EnterApp_name is called when production app_name is entered.
func (s *TreeShapeListener) EnterApp_name(ctx *parser.App_nameContext) {}

// ExitApp_name is called when production app_name is exited.
func (s *TreeShapeListener) ExitApp_name(ctx *parser.App_nameContext) {}

// EnterName_with_attribs is called when production name_with_attribs is entered.
func (s *TreeShapeListener) EnterName_with_attribs(ctx *parser.Name_with_attribsContext) {
	s.appname = ctx.App_name().GetText()
	s.module.Apps[s.appname] = &sysl.Application{
		Name: &sysl.AppName{
			Part: make([]string, 0),
		},
	}

	if attribs, ok := ctx.Attribs_or_modifiers().(*parser.Attribs_or_modifiersContext); ok {
		attribMap := make(map[string]*sysl.Attribute)
		s.module.Apps[s.appname].Attrs = attribMap
		for _, e := range attribs.AllEntry() {
			entry := e.(*parser.EntryContext)
			if nvp, ok := entry.Nvp().(*parser.NvpContext); ok {
				if qs, ok := nvp.Quoted_string().(*parser.Quoted_stringContext); ok {
					attribMap[nvp.Name().GetText()] = &sysl.Attribute{
						Attribute: &sysl.Attribute_S{
							S: strings.Trim(qs.GetText(), `"`),
						},
					}
				}
			}
		}
	}
}

// ExitName_with_attribs is called when production name_with_attribs is exited.
func (s *TreeShapeListener) ExitName_with_attribs(ctx *parser.Name_with_attribsContext) {}

// EnterModel_name is called when production model_name is entered.
func (s *TreeShapeListener) EnterModel_name(ctx *parser.Model_nameContext) {
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

// EnterVariable_substitution is called when production variable_substitution is entered.
func (s *TreeShapeListener) EnterVariable_substitution(ctx *parser.Variable_substitutionContext) {}

// ExitVariable_substitution is called when production variable_substitution is exited.
func (s *TreeShapeListener) ExitVariable_substitution(ctx *parser.Variable_substitutionContext) {}

// EnterStatic_path is called when production static_path is entered.
func (s *TreeShapeListener) EnterStatic_path(ctx *parser.Static_pathContext) {}

// ExitStatic_path is called when production static_path is exited.
func (s *TreeShapeListener) ExitStatic_path(ctx *parser.Static_pathContext) {}

// EnterQuery_var is called when production query_var is entered.
func (s *TreeShapeListener) EnterQuery_var(ctx *parser.Query_varContext) {}

// ExitQuery_var is called when production query_var is exited.
func (s *TreeShapeListener) ExitQuery_var(ctx *parser.Query_varContext) {}

// EnterQuery_param is called when production query_param is entered.
func (s *TreeShapeListener) EnterQuery_param(ctx *parser.Query_paramContext) {}

// ExitQuery_param is called when production query_param is exited.
func (s *TreeShapeListener) ExitQuery_param(ctx *parser.Query_paramContext) {}

// EnterHttp_path is called when production http_path is entered.
func (s *TreeShapeListener) EnterHttp_path(ctx *parser.Http_pathContext) {}

// ExitHttp_path is called when production http_path is exited.
func (s *TreeShapeListener) ExitHttp_path(ctx *parser.Http_pathContext) {}

// EnterEndpoint_name is called when production endpoint_name is entered.
func (s *TreeShapeListener) EnterEndpoint_name(ctx *parser.Endpoint_nameContext) {}

// ExitEndpoint_name is called when production endpoint_name is exited.
func (s *TreeShapeListener) ExitEndpoint_name(ctx *parser.Endpoint_nameContext) {}

// EnterRet_stmt is called when production ret_stmt is entered.
func (s *TreeShapeListener) EnterRet_stmt(ctx *parser.Ret_stmtContext) {}

// ExitRet_stmt is called when production ret_stmt is exited.
func (s *TreeShapeListener) ExitRet_stmt(ctx *parser.Ret_stmtContext) {}

// EnterTarget is called when production target is entered.
func (s *TreeShapeListener) EnterTarget(ctx *parser.TargetContext) {}

// ExitTarget is called when production target is exited.
func (s *TreeShapeListener) ExitTarget(ctx *parser.TargetContext) {}

// EnterTarget_endpoint is called when production target_endpoint is entered.
func (s *TreeShapeListener) EnterTarget_endpoint(ctx *parser.Target_endpointContext) {}

// ExitTarget_endpoint is called when production target_endpoint is exited.
func (s *TreeShapeListener) ExitTarget_endpoint(ctx *parser.Target_endpointContext) {}

// EnterCall_stmt is called when production call_stmt is entered.
func (s *TreeShapeListener) EnterCall_stmt(ctx *parser.Call_stmtContext) {}

// ExitCall_stmt is called when production call_stmt is exited.
func (s *TreeShapeListener) ExitCall_stmt(ctx *parser.Call_stmtContext) {}

// EnterIf_stmt is called when production if_stmt is entered.
func (s *TreeShapeListener) EnterIf_stmt(ctx *parser.If_stmtContext) {}

// ExitIf_stmt is called when production if_stmt is exited.
func (s *TreeShapeListener) ExitIf_stmt(ctx *parser.If_stmtContext) {}

// EnterIf_else is called when production if_else is entered.
func (s *TreeShapeListener) EnterIf_else(ctx *parser.If_elseContext) {}

// ExitIf_else is called when production if_else is exited.
func (s *TreeShapeListener) ExitIf_else(ctx *parser.If_elseContext) {}

// EnterFor_cond is called when production for_cond is entered.
func (s *TreeShapeListener) EnterFor_cond(ctx *parser.For_condContext) {}

// ExitFor_cond is called when production for_cond is exited.
func (s *TreeShapeListener) ExitFor_cond(ctx *parser.For_condContext) {}

// EnterFor_stmt is called when production for_stmt is entered.
func (s *TreeShapeListener) EnterFor_stmt(ctx *parser.For_stmtContext) {}

// ExitFor_stmt is called when production for_stmt is exited.
func (s *TreeShapeListener) ExitFor_stmt(ctx *parser.For_stmtContext) {}

// EnterHttp_method_comment is called when production http_method_comment is entered.
func (s *TreeShapeListener) EnterHttp_method_comment(ctx *parser.Http_method_commentContext) {}

// ExitHttp_method_comment is called when production http_method_comment is exited.
func (s *TreeShapeListener) ExitHttp_method_comment(ctx *parser.Http_method_commentContext) {}

// EnterGroup_stmt is called when production group_stmt is entered.
func (s *TreeShapeListener) EnterGroup_stmt(ctx *parser.Group_stmtContext) {}

// ExitGroup_stmt is called when production group_stmt is exited.
func (s *TreeShapeListener) ExitGroup_stmt(ctx *parser.Group_stmtContext) {}

// EnterOne_of_case_label is called when production one_of_case_label is entered.
func (s *TreeShapeListener) EnterOne_of_case_label(ctx *parser.One_of_case_labelContext) {}

// ExitOne_of_case_label is called when production one_of_case_label is exited.
func (s *TreeShapeListener) ExitOne_of_case_label(ctx *parser.One_of_case_labelContext) {}

// EnterOne_of_cases is called when production one_of_cases is entered.
func (s *TreeShapeListener) EnterOne_of_cases(ctx *parser.One_of_casesContext) {}

// ExitOne_of_cases is called when production one_of_cases is exited.
func (s *TreeShapeListener) ExitOne_of_cases(ctx *parser.One_of_casesContext) {}

// EnterOne_of_stmt is called when production one_of_stmt is entered.
func (s *TreeShapeListener) EnterOne_of_stmt(ctx *parser.One_of_stmtContext) {}

// ExitOne_of_stmt is called when production one_of_stmt is exited.
func (s *TreeShapeListener) ExitOne_of_stmt(ctx *parser.One_of_stmtContext) {}

// EnterText_stmt is called when production text_stmt is entered.
func (s *TreeShapeListener) EnterText_stmt(ctx *parser.Text_stmtContext) {}

// ExitText_stmt is called when production text_stmt is exited.
func (s *TreeShapeListener) ExitText_stmt(ctx *parser.Text_stmtContext) {}

// EnterHttp_statements is called when production http_statements is entered.
func (s *TreeShapeListener) EnterHttp_statements(ctx *parser.Http_statementsContext) {}

// ExitHttp_statements is called when production http_statements is exited.
func (s *TreeShapeListener) ExitHttp_statements(ctx *parser.Http_statementsContext) {}

// EnterMethod_def is called when production method_def is entered.
func (s *TreeShapeListener) EnterMethod_def(ctx *parser.Method_defContext) {}

// ExitMethod_def is called when production method_def is exited.
func (s *TreeShapeListener) ExitMethod_def(ctx *parser.Method_defContext) {}

// EnterEndpoint_decl is called when production endpoint_decl is entered.
func (s *TreeShapeListener) EnterEndpoint_decl(ctx *parser.Endpoint_declContext) {}

// ExitEndpoint_decl is called when production endpoint_decl is exited.
func (s *TreeShapeListener) ExitEndpoint_decl(ctx *parser.Endpoint_declContext) {}

// EnterShortcut is called when production shortcut is entered.
func (s *TreeShapeListener) EnterShortcut(ctx *parser.ShortcutContext) {}

// ExitShortcut is called when production shortcut is exited.
func (s *TreeShapeListener) ExitShortcut(ctx *parser.ShortcutContext) {}

// EnterApi_endpoint is called when production api_endpoint is entered.
func (s *TreeShapeListener) EnterApi_endpoint(ctx *parser.Api_endpointContext) {
	if ctx.WHATEVER() != nil {
		s.module.Apps[s.appname].Endpoints[ctx.WHATEVER().GetText()] = &sysl.Endpoint{
			Name: ctx.WHATEVER().GetText(),
		}
		return
	}
}

// ExitApi_endpoint is called when production api_endpoint is exited.
func (s *TreeShapeListener) ExitApi_endpoint(ctx *parser.Api_endpointContext) {}

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
func (s *TreeShapeListener) EnterEvent(ctx *parser.EventContext) {}

// ExitEvent is called when production event is exited.
func (s *TreeShapeListener) ExitEvent(ctx *parser.EventContext) {}

// EnterApp_decl is called when production app_decl is entered.
func (s *TreeShapeListener) EnterApp_decl(ctx *parser.App_declContext) {
	if len(ctx.AllTable()) > 0 {
		s.module.Apps[s.appname].Types = make(map[string]*sysl.Type)
	}
	if len(ctx.AllApi_endpoint()) > 0 {
		s.module.Apps[s.appname].Endpoints = make(map[string]*sysl.Endpoint)
	}
	if len(ctx.AllFacade()) > 0 {
		s.module.Apps[s.appname].Wrapped = &sysl.Application{
			Types: make(map[string]*sysl.Type),
		}
	}
	// fmt.Println(len(ctx.AllEvent()))
}

// ExitApp_decl is called when production app_decl is exited.
func (s *TreeShapeListener) ExitApp_decl(ctx *parser.App_declContext) {}

// EnterApplication is called when production application is entered.
func (s *TreeShapeListener) EnterApplication(ctx *parser.ApplicationContext) {}

// ExitApplication is called when production application is exited.
func (s *TreeShapeListener) ExitApplication(ctx *parser.ApplicationContext) {}

// EnterPath is called when production path is entered.
func (s *TreeShapeListener) EnterPath(ctx *parser.PathContext) {}

// ExitPath is called when production path is exited.
func (s *TreeShapeListener) ExitPath(ctx *parser.PathContext) {}

// EnterImport_stmt is called when production import_stmt is entered.
func (s *TreeShapeListener) EnterImport_stmt(ctx *parser.Import_stmtContext) {}

// ExitImport_stmt is called when production import_stmt is exited.
func (s *TreeShapeListener) ExitImport_stmt(ctx *parser.Import_stmtContext) {}

// EnterImports_decl is called when production imports_decl is entered.
func (s *TreeShapeListener) EnterImports_decl(ctx *parser.Imports_declContext) {}

// ExitImports_decl is called when production imports_decl is exited.
func (s *TreeShapeListener) ExitImports_decl(ctx *parser.Imports_declContext) {}

// EnterSysl_file is called when production sysl_file is entered.
func (s *TreeShapeListener) EnterSysl_file(ctx *parser.Sysl_fileContext) {
	s.module = &sysl.Module{
		Apps: make(map[string]*sysl.Application),
	}
}

// ExitSysl_file is called when production sysl_file is exited.
func (s *TreeShapeListener) ExitSysl_file(ctx *parser.Sysl_fileContext) {}
