// Generated from SyslParser.g4 by ANTLR 4.7.

package parser // SyslParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseSyslParserListener is a complete listener for a parse tree produced by SyslParser.
type BaseSyslParserListener struct{}

var _ SyslParserListener = &BaseSyslParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSyslParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSyslParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSyslParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSyslParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterModifier is called when production modifier is entered.
func (s *BaseSyslParserListener) EnterModifier(ctx *ModifierContext) {}

// ExitModifier is called when production modifier is exited.
func (s *BaseSyslParserListener) ExitModifier(ctx *ModifierContext) {}

// EnterSize_spec is called when production size_spec is entered.
func (s *BaseSyslParserListener) EnterSize_spec(ctx *Size_specContext) {}

// ExitSize_spec is called when production size_spec is exited.
func (s *BaseSyslParserListener) ExitSize_spec(ctx *Size_specContext) {}

// EnterModifier_list is called when production modifier_list is entered.
func (s *BaseSyslParserListener) EnterModifier_list(ctx *Modifier_listContext) {}

// ExitModifier_list is called when production modifier_list is exited.
func (s *BaseSyslParserListener) ExitModifier_list(ctx *Modifier_listContext) {}

// EnterSq_open is called when production sq_open is entered.
func (s *BaseSyslParserListener) EnterSq_open(ctx *Sq_openContext) {}

// ExitSq_open is called when production sq_open is exited.
func (s *BaseSyslParserListener) ExitSq_open(ctx *Sq_openContext) {}

// EnterModifiers is called when production modifiers is entered.
func (s *BaseSyslParserListener) EnterModifiers(ctx *ModifiersContext) {}

// ExitModifiers is called when production modifiers is exited.
func (s *BaseSyslParserListener) ExitModifiers(ctx *ModifiersContext) {}

// EnterName_str is called when production name_str is entered.
func (s *BaseSyslParserListener) EnterName_str(ctx *Name_strContext) {}

// ExitName_str is called when production name_str is exited.
func (s *BaseSyslParserListener) ExitName_str(ctx *Name_strContext) {}

// EnterReference is called when production reference is entered.
func (s *BaseSyslParserListener) EnterReference(ctx *ReferenceContext) {}

// ExitReference is called when production reference is exited.
func (s *BaseSyslParserListener) ExitReference(ctx *ReferenceContext) {}

// EnterDoc_string is called when production doc_string is entered.
func (s *BaseSyslParserListener) EnterDoc_string(ctx *Doc_stringContext) {}

// ExitDoc_string is called when production doc_string is exited.
func (s *BaseSyslParserListener) ExitDoc_string(ctx *Doc_stringContext) {}

// EnterQuoted_string is called when production quoted_string is entered.
func (s *BaseSyslParserListener) EnterQuoted_string(ctx *Quoted_stringContext) {}

// ExitQuoted_string is called when production quoted_string is exited.
func (s *BaseSyslParserListener) ExitQuoted_string(ctx *Quoted_stringContext) {}

// EnterArray_of_strings is called when production array_of_strings is entered.
func (s *BaseSyslParserListener) EnterArray_of_strings(ctx *Array_of_stringsContext) {}

// ExitArray_of_strings is called when production array_of_strings is exited.
func (s *BaseSyslParserListener) ExitArray_of_strings(ctx *Array_of_stringsContext) {}

// EnterArray_of_arrays is called when production array_of_arrays is entered.
func (s *BaseSyslParserListener) EnterArray_of_arrays(ctx *Array_of_arraysContext) {}

// ExitArray_of_arrays is called when production array_of_arrays is exited.
func (s *BaseSyslParserListener) ExitArray_of_arrays(ctx *Array_of_arraysContext) {}

// EnterNvp is called when production nvp is entered.
func (s *BaseSyslParserListener) EnterNvp(ctx *NvpContext) {}

// ExitNvp is called when production nvp is exited.
func (s *BaseSyslParserListener) ExitNvp(ctx *NvpContext) {}

// EnterAttributes is called when production attributes is entered.
func (s *BaseSyslParserListener) EnterAttributes(ctx *AttributesContext) {}

// ExitAttributes is called when production attributes is exited.
func (s *BaseSyslParserListener) ExitAttributes(ctx *AttributesContext) {}

// EnterEntry is called when production entry is entered.
func (s *BaseSyslParserListener) EnterEntry(ctx *EntryContext) {}

// ExitEntry is called when production entry is exited.
func (s *BaseSyslParserListener) ExitEntry(ctx *EntryContext) {}

// EnterAttribs_or_modifiers is called when production attribs_or_modifiers is entered.
func (s *BaseSyslParserListener) EnterAttribs_or_modifiers(ctx *Attribs_or_modifiersContext) {}

// ExitAttribs_or_modifiers is called when production attribs_or_modifiers is exited.
func (s *BaseSyslParserListener) ExitAttribs_or_modifiers(ctx *Attribs_or_modifiersContext) {}

// EnterSet_type is called when production set_type is entered.
func (s *BaseSyslParserListener) EnterSet_type(ctx *Set_typeContext) {}

// ExitSet_type is called when production set_type is exited.
func (s *BaseSyslParserListener) ExitSet_type(ctx *Set_typeContext) {}

// EnterCollection_type is called when production collection_type is entered.
func (s *BaseSyslParserListener) EnterCollection_type(ctx *Collection_typeContext) {}

// ExitCollection_type is called when production collection_type is exited.
func (s *BaseSyslParserListener) ExitCollection_type(ctx *Collection_typeContext) {}

// EnterUser_defined_type is called when production user_defined_type is entered.
func (s *BaseSyslParserListener) EnterUser_defined_type(ctx *User_defined_typeContext) {}

// ExitUser_defined_type is called when production user_defined_type is exited.
func (s *BaseSyslParserListener) ExitUser_defined_type(ctx *User_defined_typeContext) {}

// EnterMulti_line_docstring is called when production multi_line_docstring is entered.
func (s *BaseSyslParserListener) EnterMulti_line_docstring(ctx *Multi_line_docstringContext) {}

// ExitMulti_line_docstring is called when production multi_line_docstring is exited.
func (s *BaseSyslParserListener) ExitMulti_line_docstring(ctx *Multi_line_docstringContext) {}

// EnterAnnotation_value is called when production annotation_value is entered.
func (s *BaseSyslParserListener) EnterAnnotation_value(ctx *Annotation_valueContext) {}

// ExitAnnotation_value is called when production annotation_value is exited.
func (s *BaseSyslParserListener) ExitAnnotation_value(ctx *Annotation_valueContext) {}

// EnterAnnotation is called when production annotation is entered.
func (s *BaseSyslParserListener) EnterAnnotation(ctx *AnnotationContext) {}

// ExitAnnotation is called when production annotation is exited.
func (s *BaseSyslParserListener) ExitAnnotation(ctx *AnnotationContext) {}

// EnterAnnotations is called when production annotations is entered.
func (s *BaseSyslParserListener) EnterAnnotations(ctx *AnnotationsContext) {}

// ExitAnnotations is called when production annotations is exited.
func (s *BaseSyslParserListener) ExitAnnotations(ctx *AnnotationsContext) {}

// EnterField_type is called when production field_type is entered.
func (s *BaseSyslParserListener) EnterField_type(ctx *Field_typeContext) {}

// ExitField_type is called when production field_type is exited.
func (s *BaseSyslParserListener) ExitField_type(ctx *Field_typeContext) {}

// EnterArray_size is called when production array_size is entered.
func (s *BaseSyslParserListener) EnterArray_size(ctx *Array_sizeContext) {}

// ExitArray_size is called when production array_size is exited.
func (s *BaseSyslParserListener) ExitArray_size(ctx *Array_sizeContext) {}

// EnterInplace_tuple is called when production inplace_tuple is entered.
func (s *BaseSyslParserListener) EnterInplace_tuple(ctx *Inplace_tupleContext) {}

// ExitInplace_tuple is called when production inplace_tuple is exited.
func (s *BaseSyslParserListener) ExitInplace_tuple(ctx *Inplace_tupleContext) {}

// EnterField is called when production field is entered.
func (s *BaseSyslParserListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseSyslParserListener) ExitField(ctx *FieldContext) {}

// EnterInplace_table is called when production inplace_table is entered.
func (s *BaseSyslParserListener) EnterInplace_table(ctx *Inplace_tableContext) {}

// ExitInplace_table is called when production inplace_table is exited.
func (s *BaseSyslParserListener) ExitInplace_table(ctx *Inplace_tableContext) {}

// EnterTable is called when production table is entered.
func (s *BaseSyslParserListener) EnterTable(ctx *TableContext) {}

// ExitTable is called when production table is exited.
func (s *BaseSyslParserListener) ExitTable(ctx *TableContext) {}

// EnterPackage_name is called when production package_name is entered.
func (s *BaseSyslParserListener) EnterPackage_name(ctx *Package_nameContext) {}

// ExitPackage_name is called when production package_name is exited.
func (s *BaseSyslParserListener) ExitPackage_name(ctx *Package_nameContext) {}

// EnterSub_package is called when production sub_package is entered.
func (s *BaseSyslParserListener) EnterSub_package(ctx *Sub_packageContext) {}

// ExitSub_package is called when production sub_package is exited.
func (s *BaseSyslParserListener) ExitSub_package(ctx *Sub_packageContext) {}

// EnterApp_name is called when production app_name is entered.
func (s *BaseSyslParserListener) EnterApp_name(ctx *App_nameContext) {}

// ExitApp_name is called when production app_name is exited.
func (s *BaseSyslParserListener) ExitApp_name(ctx *App_nameContext) {}

// EnterName_with_attribs is called when production name_with_attribs is entered.
func (s *BaseSyslParserListener) EnterName_with_attribs(ctx *Name_with_attribsContext) {}

// ExitName_with_attribs is called when production name_with_attribs is exited.
func (s *BaseSyslParserListener) ExitName_with_attribs(ctx *Name_with_attribsContext) {}

// EnterModel_name is called when production model_name is entered.
func (s *BaseSyslParserListener) EnterModel_name(ctx *Model_nameContext) {}

// ExitModel_name is called when production model_name is exited.
func (s *BaseSyslParserListener) ExitModel_name(ctx *Model_nameContext) {}

// EnterInplace_table_def is called when production inplace_table_def is entered.
func (s *BaseSyslParserListener) EnterInplace_table_def(ctx *Inplace_table_defContext) {}

// ExitInplace_table_def is called when production inplace_table_def is exited.
func (s *BaseSyslParserListener) ExitInplace_table_def(ctx *Inplace_table_defContext) {}

// EnterTable_refs is called when production table_refs is entered.
func (s *BaseSyslParserListener) EnterTable_refs(ctx *Table_refsContext) {}

// ExitTable_refs is called when production table_refs is exited.
func (s *BaseSyslParserListener) ExitTable_refs(ctx *Table_refsContext) {}

// EnterFacade is called when production facade is entered.
func (s *BaseSyslParserListener) EnterFacade(ctx *FacadeContext) {}

// ExitFacade is called when production facade is exited.
func (s *BaseSyslParserListener) ExitFacade(ctx *FacadeContext) {}

// EnterDocumentation_stmts is called when production documentation_stmts is entered.
func (s *BaseSyslParserListener) EnterDocumentation_stmts(ctx *Documentation_stmtsContext) {}

// ExitDocumentation_stmts is called when production documentation_stmts is exited.
func (s *BaseSyslParserListener) ExitDocumentation_stmts(ctx *Documentation_stmtsContext) {}

// EnterVar_in_curly is called when production var_in_curly is entered.
func (s *BaseSyslParserListener) EnterVar_in_curly(ctx *Var_in_curlyContext) {}

// ExitVar_in_curly is called when production var_in_curly is exited.
func (s *BaseSyslParserListener) ExitVar_in_curly(ctx *Var_in_curlyContext) {}

// EnterQuery_var is called when production query_var is entered.
func (s *BaseSyslParserListener) EnterQuery_var(ctx *Query_varContext) {}

// ExitQuery_var is called when production query_var is exited.
func (s *BaseSyslParserListener) ExitQuery_var(ctx *Query_varContext) {}

// EnterQuery_param is called when production query_param is entered.
func (s *BaseSyslParserListener) EnterQuery_param(ctx *Query_paramContext) {}

// ExitQuery_param is called when production query_param is exited.
func (s *BaseSyslParserListener) ExitQuery_param(ctx *Query_paramContext) {}

// EnterHttp_path_part is called when production http_path_part is entered.
func (s *BaseSyslParserListener) EnterHttp_path_part(ctx *Http_path_partContext) {}

// ExitHttp_path_part is called when production http_path_part is exited.
func (s *BaseSyslParserListener) ExitHttp_path_part(ctx *Http_path_partContext) {}

// EnterHttp_path_var_with_type is called when production http_path_var_with_type is entered.
func (s *BaseSyslParserListener) EnterHttp_path_var_with_type(ctx *Http_path_var_with_typeContext) {}

// ExitHttp_path_var_with_type is called when production http_path_var_with_type is exited.
func (s *BaseSyslParserListener) ExitHttp_path_var_with_type(ctx *Http_path_var_with_typeContext) {}

// EnterHttp_path_static is called when production http_path_static is entered.
func (s *BaseSyslParserListener) EnterHttp_path_static(ctx *Http_path_staticContext) {}

// ExitHttp_path_static is called when production http_path_static is exited.
func (s *BaseSyslParserListener) ExitHttp_path_static(ctx *Http_path_staticContext) {}

// EnterHttp_path_suffix is called when production http_path_suffix is entered.
func (s *BaseSyslParserListener) EnterHttp_path_suffix(ctx *Http_path_suffixContext) {}

// ExitHttp_path_suffix is called when production http_path_suffix is exited.
func (s *BaseSyslParserListener) ExitHttp_path_suffix(ctx *Http_path_suffixContext) {}

// EnterHttp_path is called when production http_path is entered.
func (s *BaseSyslParserListener) EnterHttp_path(ctx *Http_pathContext) {}

// ExitHttp_path is called when production http_path is exited.
func (s *BaseSyslParserListener) ExitHttp_path(ctx *Http_pathContext) {}

// EnterEndpoint_name is called when production endpoint_name is entered.
func (s *BaseSyslParserListener) EnterEndpoint_name(ctx *Endpoint_nameContext) {}

// ExitEndpoint_name is called when production endpoint_name is exited.
func (s *BaseSyslParserListener) ExitEndpoint_name(ctx *Endpoint_nameContext) {}

// EnterRet_stmt is called when production ret_stmt is entered.
func (s *BaseSyslParserListener) EnterRet_stmt(ctx *Ret_stmtContext) {}

// ExitRet_stmt is called when production ret_stmt is exited.
func (s *BaseSyslParserListener) ExitRet_stmt(ctx *Ret_stmtContext) {}

// EnterTarget is called when production target is entered.
func (s *BaseSyslParserListener) EnterTarget(ctx *TargetContext) {}

// ExitTarget is called when production target is exited.
func (s *BaseSyslParserListener) ExitTarget(ctx *TargetContext) {}

// EnterTarget_endpoint is called when production target_endpoint is entered.
func (s *BaseSyslParserListener) EnterTarget_endpoint(ctx *Target_endpointContext) {}

// ExitTarget_endpoint is called when production target_endpoint is exited.
func (s *BaseSyslParserListener) ExitTarget_endpoint(ctx *Target_endpointContext) {}

// EnterCall_arg is called when production call_arg is entered.
func (s *BaseSyslParserListener) EnterCall_arg(ctx *Call_argContext) {}

// ExitCall_arg is called when production call_arg is exited.
func (s *BaseSyslParserListener) ExitCall_arg(ctx *Call_argContext) {}

// EnterCall_args is called when production call_args is entered.
func (s *BaseSyslParserListener) EnterCall_args(ctx *Call_argsContext) {}

// ExitCall_args is called when production call_args is exited.
func (s *BaseSyslParserListener) ExitCall_args(ctx *Call_argsContext) {}

// EnterCall_stmt is called when production call_stmt is entered.
func (s *BaseSyslParserListener) EnterCall_stmt(ctx *Call_stmtContext) {}

// ExitCall_stmt is called when production call_stmt is exited.
func (s *BaseSyslParserListener) ExitCall_stmt(ctx *Call_stmtContext) {}

// EnterIf_stmt is called when production if_stmt is entered.
func (s *BaseSyslParserListener) EnterIf_stmt(ctx *If_stmtContext) {}

// ExitIf_stmt is called when production if_stmt is exited.
func (s *BaseSyslParserListener) ExitIf_stmt(ctx *If_stmtContext) {}

// EnterElse_stmt is called when production else_stmt is entered.
func (s *BaseSyslParserListener) EnterElse_stmt(ctx *Else_stmtContext) {}

// ExitElse_stmt is called when production else_stmt is exited.
func (s *BaseSyslParserListener) ExitElse_stmt(ctx *Else_stmtContext) {}

// EnterIf_else is called when production if_else is entered.
func (s *BaseSyslParserListener) EnterIf_else(ctx *If_elseContext) {}

// ExitIf_else is called when production if_else is exited.
func (s *BaseSyslParserListener) ExitIf_else(ctx *If_elseContext) {}

// EnterFor_stmt is called when production for_stmt is entered.
func (s *BaseSyslParserListener) EnterFor_stmt(ctx *For_stmtContext) {}

// ExitFor_stmt is called when production for_stmt is exited.
func (s *BaseSyslParserListener) ExitFor_stmt(ctx *For_stmtContext) {}

// EnterHttp_method_comment is called when production http_method_comment is entered.
func (s *BaseSyslParserListener) EnterHttp_method_comment(ctx *Http_method_commentContext) {}

// ExitHttp_method_comment is called when production http_method_comment is exited.
func (s *BaseSyslParserListener) ExitHttp_method_comment(ctx *Http_method_commentContext) {}

// EnterOne_of_case_label is called when production one_of_case_label is entered.
func (s *BaseSyslParserListener) EnterOne_of_case_label(ctx *One_of_case_labelContext) {}

// ExitOne_of_case_label is called when production one_of_case_label is exited.
func (s *BaseSyslParserListener) ExitOne_of_case_label(ctx *One_of_case_labelContext) {}

// EnterOne_of_cases is called when production one_of_cases is entered.
func (s *BaseSyslParserListener) EnterOne_of_cases(ctx *One_of_casesContext) {}

// ExitOne_of_cases is called when production one_of_cases is exited.
func (s *BaseSyslParserListener) ExitOne_of_cases(ctx *One_of_casesContext) {}

// EnterOne_of_stmt is called when production one_of_stmt is entered.
func (s *BaseSyslParserListener) EnterOne_of_stmt(ctx *One_of_stmtContext) {}

// ExitOne_of_stmt is called when production one_of_stmt is exited.
func (s *BaseSyslParserListener) ExitOne_of_stmt(ctx *One_of_stmtContext) {}

// EnterText_stmt is called when production text_stmt is entered.
func (s *BaseSyslParserListener) EnterText_stmt(ctx *Text_stmtContext) {}

// ExitText_stmt is called when production text_stmt is exited.
func (s *BaseSyslParserListener) ExitText_stmt(ctx *Text_stmtContext) {}

// EnterMixin is called when production mixin is entered.
func (s *BaseSyslParserListener) EnterMixin(ctx *MixinContext) {}

// ExitMixin is called when production mixin is exited.
func (s *BaseSyslParserListener) ExitMixin(ctx *MixinContext) {}

// EnterParam is called when production param is entered.
func (s *BaseSyslParserListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BaseSyslParserListener) ExitParam(ctx *ParamContext) {}

// EnterParam_list is called when production param_list is entered.
func (s *BaseSyslParserListener) EnterParam_list(ctx *Param_listContext) {}

// ExitParam_list is called when production param_list is exited.
func (s *BaseSyslParserListener) ExitParam_list(ctx *Param_listContext) {}

// EnterParams is called when production params is entered.
func (s *BaseSyslParserListener) EnterParams(ctx *ParamsContext) {}

// ExitParams is called when production params is exited.
func (s *BaseSyslParserListener) ExitParams(ctx *ParamsContext) {}

// EnterStatements is called when production statements is entered.
func (s *BaseSyslParserListener) EnterStatements(ctx *StatementsContext) {}

// ExitStatements is called when production statements is exited.
func (s *BaseSyslParserListener) ExitStatements(ctx *StatementsContext) {}

// EnterMethod_def is called when production method_def is entered.
func (s *BaseSyslParserListener) EnterMethod_def(ctx *Method_defContext) {}

// ExitMethod_def is called when production method_def is exited.
func (s *BaseSyslParserListener) ExitMethod_def(ctx *Method_defContext) {}

// EnterShortcut is called when production shortcut is entered.
func (s *BaseSyslParserListener) EnterShortcut(ctx *ShortcutContext) {}

// ExitShortcut is called when production shortcut is exited.
func (s *BaseSyslParserListener) ExitShortcut(ctx *ShortcutContext) {}

// EnterSimple_endpoint is called when production simple_endpoint is entered.
func (s *BaseSyslParserListener) EnterSimple_endpoint(ctx *Simple_endpointContext) {}

// ExitSimple_endpoint is called when production simple_endpoint is exited.
func (s *BaseSyslParserListener) ExitSimple_endpoint(ctx *Simple_endpointContext) {}

// EnterRest_endpoint is called when production rest_endpoint is entered.
func (s *BaseSyslParserListener) EnterRest_endpoint(ctx *Rest_endpointContext) {}

// ExitRest_endpoint is called when production rest_endpoint is exited.
func (s *BaseSyslParserListener) ExitRest_endpoint(ctx *Rest_endpointContext) {}

// EnterCollector_query_var is called when production collector_query_var is entered.
func (s *BaseSyslParserListener) EnterCollector_query_var(ctx *Collector_query_varContext) {}

// ExitCollector_query_var is called when production collector_query_var is exited.
func (s *BaseSyslParserListener) ExitCollector_query_var(ctx *Collector_query_varContext) {}

// EnterCollector_query_param is called when production collector_query_param is entered.
func (s *BaseSyslParserListener) EnterCollector_query_param(ctx *Collector_query_paramContext) {}

// ExitCollector_query_param is called when production collector_query_param is exited.
func (s *BaseSyslParserListener) ExitCollector_query_param(ctx *Collector_query_paramContext) {}

// EnterCollector_call_stmt is called when production collector_call_stmt is entered.
func (s *BaseSyslParserListener) EnterCollector_call_stmt(ctx *Collector_call_stmtContext) {}

// ExitCollector_call_stmt is called when production collector_call_stmt is exited.
func (s *BaseSyslParserListener) ExitCollector_call_stmt(ctx *Collector_call_stmtContext) {}

// EnterCollector_http_stmt_part is called when production collector_http_stmt_part is entered.
func (s *BaseSyslParserListener) EnterCollector_http_stmt_part(ctx *Collector_http_stmt_partContext) {}

// ExitCollector_http_stmt_part is called when production collector_http_stmt_part is exited.
func (s *BaseSyslParserListener) ExitCollector_http_stmt_part(ctx *Collector_http_stmt_partContext) {}

// EnterCollector_http_stmt_suffix is called when production collector_http_stmt_suffix is entered.
func (s *BaseSyslParserListener) EnterCollector_http_stmt_suffix(ctx *Collector_http_stmt_suffixContext) {
}

// ExitCollector_http_stmt_suffix is called when production collector_http_stmt_suffix is exited.
func (s *BaseSyslParserListener) ExitCollector_http_stmt_suffix(ctx *Collector_http_stmt_suffixContext) {
}

// EnterCollector_http_stmt is called when production collector_http_stmt is entered.
func (s *BaseSyslParserListener) EnterCollector_http_stmt(ctx *Collector_http_stmtContext) {}

// ExitCollector_http_stmt is called when production collector_http_stmt is exited.
func (s *BaseSyslParserListener) ExitCollector_http_stmt(ctx *Collector_http_stmtContext) {}

// EnterCollector_stmts is called when production collector_stmts is entered.
func (s *BaseSyslParserListener) EnterCollector_stmts(ctx *Collector_stmtsContext) {}

// ExitCollector_stmts is called when production collector_stmts is exited.
func (s *BaseSyslParserListener) ExitCollector_stmts(ctx *Collector_stmtsContext) {}

// EnterCollector is called when production collector is entered.
func (s *BaseSyslParserListener) EnterCollector(ctx *CollectorContext) {}

// ExitCollector is called when production collector is exited.
func (s *BaseSyslParserListener) ExitCollector(ctx *CollectorContext) {}

// EnterEvent is called when production event is entered.
func (s *BaseSyslParserListener) EnterEvent(ctx *EventContext) {}

// ExitEvent is called when production event is exited.
func (s *BaseSyslParserListener) ExitEvent(ctx *EventContext) {}

// EnterSubscribe is called when production subscribe is entered.
func (s *BaseSyslParserListener) EnterSubscribe(ctx *SubscribeContext) {}

// ExitSubscribe is called when production subscribe is exited.
func (s *BaseSyslParserListener) ExitSubscribe(ctx *SubscribeContext) {}

// EnterApp_decl is called when production app_decl is entered.
func (s *BaseSyslParserListener) EnterApp_decl(ctx *App_declContext) {}

// ExitApp_decl is called when production app_decl is exited.
func (s *BaseSyslParserListener) ExitApp_decl(ctx *App_declContext) {}

// EnterApplication is called when production application is entered.
func (s *BaseSyslParserListener) EnterApplication(ctx *ApplicationContext) {}

// ExitApplication is called when production application is exited.
func (s *BaseSyslParserListener) ExitApplication(ctx *ApplicationContext) {}

// EnterPath is called when production path is entered.
func (s *BaseSyslParserListener) EnterPath(ctx *PathContext) {}

// ExitPath is called when production path is exited.
func (s *BaseSyslParserListener) ExitPath(ctx *PathContext) {}

// EnterImport_stmt is called when production import_stmt is entered.
func (s *BaseSyslParserListener) EnterImport_stmt(ctx *Import_stmtContext) {}

// ExitImport_stmt is called when production import_stmt is exited.
func (s *BaseSyslParserListener) ExitImport_stmt(ctx *Import_stmtContext) {}

// EnterImports_decl is called when production imports_decl is entered.
func (s *BaseSyslParserListener) EnterImports_decl(ctx *Imports_declContext) {}

// ExitImports_decl is called when production imports_decl is exited.
func (s *BaseSyslParserListener) ExitImports_decl(ctx *Imports_declContext) {}

// EnterSysl_file is called when production sysl_file is entered.
func (s *BaseSyslParserListener) EnterSysl_file(ctx *Sysl_fileContext) {}

// ExitSysl_file is called when production sysl_file is exited.
func (s *BaseSyslParserListener) ExitSysl_file(ctx *Sysl_fileContext) {}
