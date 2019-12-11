// Generated from pkg/grammar/SyslParser.g4 by ANTLR 4.7.

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

// EnterEntry is called when production entry is entered.
func (s *BaseSyslParserListener) EnterEntry(ctx *EntryContext) {}

// ExitEntry is called when production entry is exited.
func (s *BaseSyslParserListener) ExitEntry(ctx *EntryContext) {}

// EnterAttribs_or_modifiers is called when production attribs_or_modifiers is entered.
func (s *BaseSyslParserListener) EnterAttribs_or_modifiers(ctx *Attribs_or_modifiersContext) {}

// ExitAttribs_or_modifiers is called when production attribs_or_modifiers is exited.
func (s *BaseSyslParserListener) ExitAttribs_or_modifiers(ctx *Attribs_or_modifiersContext) {}

// EnterUser_defined_type is called when production user_defined_type is entered.
func (s *BaseSyslParserListener) EnterUser_defined_type(ctx *User_defined_typeContext) {}

// ExitUser_defined_type is called when production user_defined_type is exited.
func (s *BaseSyslParserListener) ExitUser_defined_type(ctx *User_defined_typeContext) {}

// EnterTypes is called when production types is entered.
func (s *BaseSyslParserListener) EnterTypes(ctx *TypesContext) {}

// ExitTypes is called when production types is exited.
func (s *BaseSyslParserListener) ExitTypes(ctx *TypesContext) {}

// EnterSet_of is called when production set_of is entered.
func (s *BaseSyslParserListener) EnterSet_of(ctx *Set_ofContext) {}

// ExitSet_of is called when production set_of is exited.
func (s *BaseSyslParserListener) ExitSet_of(ctx *Set_ofContext) {}

// EnterSet_type is called when production set_type is entered.
func (s *BaseSyslParserListener) EnterSet_type(ctx *Set_typeContext) {}

// ExitSet_type is called when production set_type is exited.
func (s *BaseSyslParserListener) ExitSet_type(ctx *Set_typeContext) {}

// EnterSequence_of is called when production sequence_of is entered.
func (s *BaseSyslParserListener) EnterSequence_of(ctx *Sequence_ofContext) {}

// ExitSequence_of is called when production sequence_of is exited.
func (s *BaseSyslParserListener) ExitSequence_of(ctx *Sequence_ofContext) {}

// EnterSequence_type is called when production sequence_type is entered.
func (s *BaseSyslParserListener) EnterSequence_type(ctx *Sequence_typeContext) {}

// ExitSequence_type is called when production sequence_type is exited.
func (s *BaseSyslParserListener) ExitSequence_type(ctx *Sequence_typeContext) {}

// EnterCollection_type is called when production collection_type is entered.
func (s *BaseSyslParserListener) EnterCollection_type(ctx *Collection_typeContext) {}

// ExitCollection_type is called when production collection_type is exited.
func (s *BaseSyslParserListener) ExitCollection_type(ctx *Collection_typeContext) {}

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

// EnterTable_stmts is called when production table_stmts is entered.
func (s *BaseSyslParserListener) EnterTable_stmts(ctx *Table_stmtsContext) {}

// ExitTable_stmts is called when production table_stmts is exited.
func (s *BaseSyslParserListener) ExitTable_stmts(ctx *Table_stmtsContext) {}

// EnterTable_def is called when production table_def is entered.
func (s *BaseSyslParserListener) EnterTable_def(ctx *Table_defContext) {}

// ExitTable_def is called when production table_def is exited.
func (s *BaseSyslParserListener) ExitTable_def(ctx *Table_defContext) {}

// EnterTable is called when production table is entered.
func (s *BaseSyslParserListener) EnterTable(ctx *TableContext) {}

// ExitTable is called when production table is exited.
func (s *BaseSyslParserListener) ExitTable(ctx *TableContext) {}

// EnterUnion is called when production union is entered.
func (s *BaseSyslParserListener) EnterUnion(ctx *UnionContext) {}

// ExitUnion is called when production union is exited.
func (s *BaseSyslParserListener) ExitUnion(ctx *UnionContext) {}

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

// EnterGroup_stmt is called when production group_stmt is entered.
func (s *BaseSyslParserListener) EnterGroup_stmt(ctx *Group_stmtContext) {}

// ExitGroup_stmt is called when production group_stmt is exited.
func (s *BaseSyslParserListener) ExitGroup_stmt(ctx *Group_stmtContext) {}

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

// EnterPublisher is called when production publisher is entered.
func (s *BaseSyslParserListener) EnterPublisher(ctx *PublisherContext) {}

// ExitPublisher is called when production publisher is exited.
func (s *BaseSyslParserListener) ExitPublisher(ctx *PublisherContext) {}

// EnterSubscriber is called when production subscriber is entered.
func (s *BaseSyslParserListener) EnterSubscriber(ctx *SubscriberContext) {}

// ExitSubscriber is called when production subscriber is exited.
func (s *BaseSyslParserListener) ExitSubscriber(ctx *SubscriberContext) {}

// EnterCollector_pubsub_call is called when production collector_pubsub_call is entered.
func (s *BaseSyslParserListener) EnterCollector_pubsub_call(ctx *Collector_pubsub_callContext) {}

// ExitCollector_pubsub_call is called when production collector_pubsub_call is exited.
func (s *BaseSyslParserListener) ExitCollector_pubsub_call(ctx *Collector_pubsub_callContext) {}

// EnterCollector_action_stmt is called when production collector_action_stmt is entered.
func (s *BaseSyslParserListener) EnterCollector_action_stmt(ctx *Collector_action_stmtContext) {}

// ExitCollector_action_stmt is called when production collector_action_stmt is exited.
func (s *BaseSyslParserListener) ExitCollector_action_stmt(ctx *Collector_action_stmtContext) {}

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

// EnterView_type_spec is called when production view_type_spec is entered.
func (s *BaseSyslParserListener) EnterView_type_spec(ctx *View_type_specContext) {}

// ExitView_type_spec is called when production view_type_spec is exited.
func (s *BaseSyslParserListener) ExitView_type_spec(ctx *View_type_specContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BaseSyslParserListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BaseSyslParserListener) ExitLiteral(ctx *LiteralContext) {}

// EnterExpr_table_of_op is called when production expr_table_of_op is entered.
func (s *BaseSyslParserListener) EnterExpr_table_of_op(ctx *Expr_table_of_opContext) {}

// ExitExpr_table_of_op is called when production expr_table_of_op is exited.
func (s *BaseSyslParserListener) ExitExpr_table_of_op(ctx *Expr_table_of_opContext) {}

// EnterFunc_arg is called when production func_arg is entered.
func (s *BaseSyslParserListener) EnterFunc_arg(ctx *Func_argContext) {}

// ExitFunc_arg is called when production func_arg is exited.
func (s *BaseSyslParserListener) ExitFunc_arg(ctx *Func_argContext) {}

// EnterFunc_args is called when production func_args is entered.
func (s *BaseSyslParserListener) EnterFunc_args(ctx *Func_argsContext) {}

// ExitFunc_args is called when production func_args is exited.
func (s *BaseSyslParserListener) ExitFunc_args(ctx *Func_argsContext) {}

// EnterExpr_func is called when production expr_func is entered.
func (s *BaseSyslParserListener) EnterExpr_func(ctx *Expr_funcContext) {}

// ExitExpr_func is called when production expr_func is exited.
func (s *BaseSyslParserListener) ExitExpr_func(ctx *Expr_funcContext) {}

// EnterRank_expr is called when production rank_expr is entered.
func (s *BaseSyslParserListener) EnterRank_expr(ctx *Rank_exprContext) {}

// ExitRank_expr is called when production rank_expr is exited.
func (s *BaseSyslParserListener) ExitRank_expr(ctx *Rank_exprContext) {}

// EnterRank_expr_list is called when production rank_expr_list is entered.
func (s *BaseSyslParserListener) EnterRank_expr_list(ctx *Rank_expr_listContext) {}

// ExitRank_expr_list is called when production rank_expr_list is exited.
func (s *BaseSyslParserListener) ExitRank_expr_list(ctx *Rank_expr_listContext) {}

// EnterExpr_rank_func is called when production expr_rank_func is entered.
func (s *BaseSyslParserListener) EnterExpr_rank_func(ctx *Expr_rank_funcContext) {}

// ExitExpr_rank_func is called when production expr_rank_func is exited.
func (s *BaseSyslParserListener) ExitExpr_rank_func(ctx *Expr_rank_funcContext) {}

// EnterExpr_agg_func is called when production expr_agg_func is entered.
func (s *BaseSyslParserListener) EnterExpr_agg_func(ctx *Expr_agg_funcContext) {}

// ExitExpr_agg_func is called when production expr_agg_func is exited.
func (s *BaseSyslParserListener) ExitExpr_agg_func(ctx *Expr_agg_funcContext) {}

// EnterE_scope_var is called when production e_scope_var is entered.
func (s *BaseSyslParserListener) EnterE_scope_var(ctx *E_scope_varContext) {}

// ExitE_scope_var is called when production e_scope_var is exited.
func (s *BaseSyslParserListener) ExitE_scope_var(ctx *E_scope_varContext) {}

// EnterFirst_func_target is called when production first_func_target is entered.
func (s *BaseSyslParserListener) EnterFirst_func_target(ctx *First_func_targetContext) {}

// ExitFirst_func_target is called when production first_func_target is exited.
func (s *BaseSyslParserListener) ExitFirst_func_target(ctx *First_func_targetContext) {}

// EnterExpr_first_func is called when production expr_first_func is entered.
func (s *BaseSyslParserListener) EnterExpr_first_func(ctx *Expr_first_funcContext) {}

// ExitExpr_first_func is called when production expr_first_func is exited.
func (s *BaseSyslParserListener) ExitExpr_first_func(ctx *Expr_first_funcContext) {}

// EnterE_single_arg_func is called when production e_single_arg_func is entered.
func (s *BaseSyslParserListener) EnterE_single_arg_func(ctx *E_single_arg_funcContext) {}

// ExitE_single_arg_func is called when production e_single_arg_func is exited.
func (s *BaseSyslParserListener) ExitE_single_arg_func(ctx *E_single_arg_funcContext) {}

// EnterExpr_single_arg_func is called when production expr_single_arg_func is entered.
func (s *BaseSyslParserListener) EnterExpr_single_arg_func(ctx *Expr_single_arg_funcContext) {}

// ExitExpr_single_arg_func is called when production expr_single_arg_func is exited.
func (s *BaseSyslParserListener) ExitExpr_single_arg_func(ctx *Expr_single_arg_funcContext) {}

// EnterExpr_any_func is called when production expr_any_func is entered.
func (s *BaseSyslParserListener) EnterExpr_any_func(ctx *Expr_any_funcContext) {}

// ExitExpr_any_func is called when production expr_any_func is exited.
func (s *BaseSyslParserListener) ExitExpr_any_func(ctx *Expr_any_funcContext) {}

// EnterExpr_single_or_null is called when production expr_single_or_null is entered.
func (s *BaseSyslParserListener) EnterExpr_single_or_null(ctx *Expr_single_or_nullContext) {}

// ExitExpr_single_or_null is called when production expr_single_or_null is exited.
func (s *BaseSyslParserListener) ExitExpr_single_or_null(ctx *Expr_single_or_nullContext) {}

// EnterExpr_snapshot is called when production expr_snapshot is entered.
func (s *BaseSyslParserListener) EnterExpr_snapshot(ctx *Expr_snapshotContext) {}

// ExitExpr_snapshot is called when production expr_snapshot is exited.
func (s *BaseSyslParserListener) ExitExpr_snapshot(ctx *Expr_snapshotContext) {}

// EnterExpr_count is called when production expr_count is entered.
func (s *BaseSyslParserListener) EnterExpr_count(ctx *Expr_countContext) {}

// ExitExpr_count is called when production expr_count is exited.
func (s *BaseSyslParserListener) ExitExpr_count(ctx *Expr_countContext) {}

// EnterExpr_navigate_attr is called when production expr_navigate_attr is entered.
func (s *BaseSyslParserListener) EnterExpr_navigate_attr(ctx *Expr_navigate_attrContext) {}

// ExitExpr_navigate_attr is called when production expr_navigate_attr is exited.
func (s *BaseSyslParserListener) ExitExpr_navigate_attr(ctx *Expr_navigate_attrContext) {}

// EnterExpr_navigate is called when production expr_navigate is entered.
func (s *BaseSyslParserListener) EnterExpr_navigate(ctx *Expr_navigateContext) {}

// ExitExpr_navigate is called when production expr_navigate is exited.
func (s *BaseSyslParserListener) ExitExpr_navigate(ctx *Expr_navigateContext) {}

// EnterMatching_rhs is called when production matching_rhs is entered.
func (s *BaseSyslParserListener) EnterMatching_rhs(ctx *Matching_rhsContext) {}

// ExitMatching_rhs is called when production matching_rhs is exited.
func (s *BaseSyslParserListener) ExitMatching_rhs(ctx *Matching_rhsContext) {}

// EnterSquiggly_args is called when production squiggly_args is entered.
func (s *BaseSyslParserListener) EnterSquiggly_args(ctx *Squiggly_argsContext) {}

// ExitSquiggly_args is called when production squiggly_args is exited.
func (s *BaseSyslParserListener) ExitSquiggly_args(ctx *Squiggly_argsContext) {}

// EnterExpr_matching is called when production expr_matching is entered.
func (s *BaseSyslParserListener) EnterExpr_matching(ctx *Expr_matchingContext) {}

// ExitExpr_matching is called when production expr_matching is exited.
func (s *BaseSyslParserListener) ExitExpr_matching(ctx *Expr_matchingContext) {}

// EnterRelop is called when production relop is entered.
func (s *BaseSyslParserListener) EnterRelop(ctx *RelopContext) {}

// ExitRelop is called when production relop is exited.
func (s *BaseSyslParserListener) ExitRelop(ctx *RelopContext) {}

// EnterList_item is called when production list_item is entered.
func (s *BaseSyslParserListener) EnterList_item(ctx *List_itemContext) {}

// ExitList_item is called when production list_item is exited.
func (s *BaseSyslParserListener) ExitList_item(ctx *List_itemContext) {}

// EnterExpr_list is called when production expr_list is entered.
func (s *BaseSyslParserListener) EnterExpr_list(ctx *Expr_listContext) {}

// ExitExpr_list is called when production expr_list is exited.
func (s *BaseSyslParserListener) ExitExpr_list(ctx *Expr_listContext) {}

// EnterExpr_set is called when production expr_set is entered.
func (s *BaseSyslParserListener) EnterExpr_set(ctx *Expr_setContext) {}

// ExitExpr_set is called when production expr_set is exited.
func (s *BaseSyslParserListener) ExitExpr_set(ctx *Expr_setContext) {}

// EnterEmpty_tuple is called when production empty_tuple is entered.
func (s *BaseSyslParserListener) EnterEmpty_tuple(ctx *Empty_tupleContext) {}

// ExitEmpty_tuple is called when production empty_tuple is exited.
func (s *BaseSyslParserListener) ExitEmpty_tuple(ctx *Empty_tupleContext) {}

// EnterAtom_dot_relop is called when production atom_dot_relop is entered.
func (s *BaseSyslParserListener) EnterAtom_dot_relop(ctx *Atom_dot_relopContext) {}

// ExitAtom_dot_relop is called when production atom_dot_relop is exited.
func (s *BaseSyslParserListener) ExitAtom_dot_relop(ctx *Atom_dot_relopContext) {}

// EnterAtomT_implied_dot is called when production atomT_implied_dot is entered.
func (s *BaseSyslParserListener) EnterAtomT_implied_dot(ctx *AtomT_implied_dotContext) {}

// ExitAtomT_implied_dot is called when production atomT_implied_dot is exited.
func (s *BaseSyslParserListener) ExitAtomT_implied_dot(ctx *AtomT_implied_dotContext) {}

// EnterAtomT_name is called when production atomT_name is entered.
func (s *BaseSyslParserListener) EnterAtomT_name(ctx *AtomT_nameContext) {}

// ExitAtomT_name is called when production atomT_name is exited.
func (s *BaseSyslParserListener) ExitAtomT_name(ctx *AtomT_nameContext) {}

// EnterAtomT_paren is called when production atomT_paren is entered.
func (s *BaseSyslParserListener) EnterAtomT_paren(ctx *AtomT_parenContext) {}

// ExitAtomT_paren is called when production atomT_paren is exited.
func (s *BaseSyslParserListener) ExitAtomT_paren(ctx *AtomT_parenContext) {}

// EnterExpr_atom_list is called when production expr_atom_list is entered.
func (s *BaseSyslParserListener) EnterExpr_atom_list(ctx *Expr_atom_listContext) {}

// ExitExpr_atom_list is called when production expr_atom_list is exited.
func (s *BaseSyslParserListener) ExitExpr_atom_list(ctx *Expr_atom_listContext) {}

// EnterAtomT is called when production atomT is entered.
func (s *BaseSyslParserListener) EnterAtomT(ctx *AtomTContext) {}

// ExitAtomT is called when production atomT is exited.
func (s *BaseSyslParserListener) ExitAtomT(ctx *AtomTContext) {}

// EnterAtom is called when production atom is entered.
func (s *BaseSyslParserListener) EnterAtom(ctx *AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *BaseSyslParserListener) ExitAtom(ctx *AtomContext) {}

// EnterPowerT is called when production powerT is entered.
func (s *BaseSyslParserListener) EnterPowerT(ctx *PowerTContext) {}

// ExitPowerT is called when production powerT is exited.
func (s *BaseSyslParserListener) ExitPowerT(ctx *PowerTContext) {}

// EnterPower is called when production power is entered.
func (s *BaseSyslParserListener) EnterPower(ctx *PowerContext) {}

// ExitPower is called when production power is exited.
func (s *BaseSyslParserListener) ExitPower(ctx *PowerContext) {}

// EnterUnaryTerm is called when production unaryTerm is entered.
func (s *BaseSyslParserListener) EnterUnaryTerm(ctx *UnaryTermContext) {}

// ExitUnaryTerm is called when production unaryTerm is exited.
func (s *BaseSyslParserListener) ExitUnaryTerm(ctx *UnaryTermContext) {}

// EnterTermT is called when production termT is entered.
func (s *BaseSyslParserListener) EnterTermT(ctx *TermTContext) {}

// ExitTermT is called when production termT is exited.
func (s *BaseSyslParserListener) ExitTermT(ctx *TermTContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseSyslParserListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseSyslParserListener) ExitTerm(ctx *TermContext) {}

// EnterBinexprT is called when production binexprT is entered.
func (s *BaseSyslParserListener) EnterBinexprT(ctx *BinexprTContext) {}

// ExitBinexprT is called when production binexprT is exited.
func (s *BaseSyslParserListener) ExitBinexprT(ctx *BinexprTContext) {}

// EnterBinexpr is called when production binexpr is entered.
func (s *BaseSyslParserListener) EnterBinexpr(ctx *BinexprContext) {}

// ExitBinexpr is called when production binexpr is exited.
func (s *BaseSyslParserListener) ExitBinexpr(ctx *BinexprContext) {}

// EnterE_compare_ops is called when production e_compare_ops is entered.
func (s *BaseSyslParserListener) EnterE_compare_ops(ctx *E_compare_opsContext) {}

// ExitE_compare_ops is called when production e_compare_ops is exited.
func (s *BaseSyslParserListener) ExitE_compare_ops(ctx *E_compare_opsContext) {}

// EnterExpr_rel is called when production expr_rel is entered.
func (s *BaseSyslParserListener) EnterExpr_rel(ctx *Expr_relContext) {}

// ExitExpr_rel is called when production expr_rel is exited.
func (s *BaseSyslParserListener) ExitExpr_rel(ctx *Expr_relContext) {}

// EnterExpr_bitand is called when production expr_bitand is entered.
func (s *BaseSyslParserListener) EnterExpr_bitand(ctx *Expr_bitandContext) {}

// ExitExpr_bitand is called when production expr_bitand is exited.
func (s *BaseSyslParserListener) ExitExpr_bitand(ctx *Expr_bitandContext) {}

// EnterExpr_bitxor is called when production expr_bitxor is entered.
func (s *BaseSyslParserListener) EnterExpr_bitxor(ctx *Expr_bitxorContext) {}

// ExitExpr_bitxor is called when production expr_bitxor is exited.
func (s *BaseSyslParserListener) ExitExpr_bitxor(ctx *Expr_bitxorContext) {}

// EnterExpr_bitor is called when production expr_bitor is entered.
func (s *BaseSyslParserListener) EnterExpr_bitor(ctx *Expr_bitorContext) {}

// ExitExpr_bitor is called when production expr_bitor is exited.
func (s *BaseSyslParserListener) ExitExpr_bitor(ctx *Expr_bitorContext) {}

// EnterExpr_and is called when production expr_and is entered.
func (s *BaseSyslParserListener) EnterExpr_and(ctx *Expr_andContext) {}

// ExitExpr_and is called when production expr_and is exited.
func (s *BaseSyslParserListener) ExitExpr_and(ctx *Expr_andContext) {}

// EnterExpr_or is called when production expr_or is entered.
func (s *BaseSyslParserListener) EnterExpr_or(ctx *Expr_orContext) {}

// ExitExpr_or is called when production expr_or is exited.
func (s *BaseSyslParserListener) ExitExpr_or(ctx *Expr_orContext) {}

// EnterExpr_but_not is called when production expr_but_not is entered.
func (s *BaseSyslParserListener) EnterExpr_but_not(ctx *Expr_but_notContext) {}

// ExitExpr_but_not is called when production expr_but_not is exited.
func (s *BaseSyslParserListener) ExitExpr_but_not(ctx *Expr_but_notContext) {}

// EnterExpr_coalesce is called when production expr_coalesce is entered.
func (s *BaseSyslParserListener) EnterExpr_coalesce(ctx *Expr_coalesceContext) {}

// ExitExpr_coalesce is called when production expr_coalesce is exited.
func (s *BaseSyslParserListener) ExitExpr_coalesce(ctx *Expr_coalesceContext) {}

// EnterIf_one_liner is called when production if_one_liner is entered.
func (s *BaseSyslParserListener) EnterIf_one_liner(ctx *If_one_linerContext) {}

// ExitIf_one_liner is called when production if_one_liner is exited.
func (s *BaseSyslParserListener) ExitIf_one_liner(ctx *If_one_linerContext) {}

// EnterElse_block_stmt is called when production else_block_stmt is entered.
func (s *BaseSyslParserListener) EnterElse_block_stmt(ctx *Else_block_stmtContext) {}

// ExitElse_block_stmt is called when production else_block_stmt is exited.
func (s *BaseSyslParserListener) ExitElse_block_stmt(ctx *Else_block_stmtContext) {}

// EnterControl_item is called when production control_item is entered.
func (s *BaseSyslParserListener) EnterControl_item(ctx *Control_itemContext) {}

// ExitControl_item is called when production control_item is exited.
func (s *BaseSyslParserListener) ExitControl_item(ctx *Control_itemContext) {}

// EnterIf_controls is called when production if_controls is entered.
func (s *BaseSyslParserListener) EnterIf_controls(ctx *If_controlsContext) {}

// ExitIf_controls is called when production if_controls is exited.
func (s *BaseSyslParserListener) ExitIf_controls(ctx *If_controlsContext) {}

// EnterCond_block is called when production cond_block is entered.
func (s *BaseSyslParserListener) EnterCond_block(ctx *Cond_blockContext) {}

// ExitCond_block is called when production cond_block is exited.
func (s *BaseSyslParserListener) ExitCond_block(ctx *Cond_blockContext) {}

// EnterFinal_else is called when production final_else is entered.
func (s *BaseSyslParserListener) EnterFinal_else(ctx *Final_elseContext) {}

// ExitFinal_else is called when production final_else is exited.
func (s *BaseSyslParserListener) ExitFinal_else(ctx *Final_elseContext) {}

// EnterIfvar is called when production ifvar is entered.
func (s *BaseSyslParserListener) EnterIfvar(ctx *IfvarContext) {}

// ExitIfvar is called when production ifvar is exited.
func (s *BaseSyslParserListener) ExitIfvar(ctx *IfvarContext) {}

// EnterIf_multiple_lines is called when production if_multiple_lines is entered.
func (s *BaseSyslParserListener) EnterIf_multiple_lines(ctx *If_multiple_linesContext) {}

// ExitIf_multiple_lines is called when production if_multiple_lines is exited.
func (s *BaseSyslParserListener) ExitIf_multiple_lines(ctx *If_multiple_linesContext) {}

// EnterExpr_if_else is called when production expr_if_else is entered.
func (s *BaseSyslParserListener) EnterExpr_if_else(ctx *Expr_if_elseContext) {}

// ExitExpr_if_else is called when production expr_if_else is exited.
func (s *BaseSyslParserListener) ExitExpr_if_else(ctx *Expr_if_elseContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseSyslParserListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseSyslParserListener) ExitExpr(ctx *ExprContext) {}

// EnterExpr_assign is called when production expr_assign is entered.
func (s *BaseSyslParserListener) EnterExpr_assign(ctx *Expr_assignContext) {}

// ExitExpr_assign is called when production expr_assign is exited.
func (s *BaseSyslParserListener) ExitExpr_assign(ctx *Expr_assignContext) {}

// EnterExpr_simple_assign is called when production expr_simple_assign is entered.
func (s *BaseSyslParserListener) EnterExpr_simple_assign(ctx *Expr_simple_assignContext) {}

// ExitExpr_simple_assign is called when production expr_simple_assign is exited.
func (s *BaseSyslParserListener) ExitExpr_simple_assign(ctx *Expr_simple_assignContext) {}

// EnterExpr_let_statement is called when production expr_let_statement is entered.
func (s *BaseSyslParserListener) EnterExpr_let_statement(ctx *Expr_let_statementContext) {}

// ExitExpr_let_statement is called when production expr_let_statement is exited.
func (s *BaseSyslParserListener) ExitExpr_let_statement(ctx *Expr_let_statementContext) {}

// EnterExpr_table_of_statement is called when production expr_table_of_statement is entered.
func (s *BaseSyslParserListener) EnterExpr_table_of_statement(ctx *Expr_table_of_statementContext) {}

// ExitExpr_table_of_statement is called when production expr_table_of_statement is exited.
func (s *BaseSyslParserListener) ExitExpr_table_of_statement(ctx *Expr_table_of_statementContext) {}

// EnterExpr_dot_assign is called when production expr_dot_assign is entered.
func (s *BaseSyslParserListener) EnterExpr_dot_assign(ctx *Expr_dot_assignContext) {}

// ExitExpr_dot_assign is called when production expr_dot_assign is exited.
func (s *BaseSyslParserListener) ExitExpr_dot_assign(ctx *Expr_dot_assignContext) {}

// EnterExpr_statement_no_nl is called when production expr_statement_no_nl is entered.
func (s *BaseSyslParserListener) EnterExpr_statement_no_nl(ctx *Expr_statement_no_nlContext) {}

// ExitExpr_statement_no_nl is called when production expr_statement_no_nl is exited.
func (s *BaseSyslParserListener) ExitExpr_statement_no_nl(ctx *Expr_statement_no_nlContext) {}

// EnterExpr_statement is called when production expr_statement is entered.
func (s *BaseSyslParserListener) EnterExpr_statement(ctx *Expr_statementContext) {}

// ExitExpr_statement is called when production expr_statement is exited.
func (s *BaseSyslParserListener) ExitExpr_statement(ctx *Expr_statementContext) {}

// EnterExpr_inject_stmt is called when production expr_inject_stmt is entered.
func (s *BaseSyslParserListener) EnterExpr_inject_stmt(ctx *Expr_inject_stmtContext) {}

// ExitExpr_inject_stmt is called when production expr_inject_stmt is exited.
func (s *BaseSyslParserListener) ExitExpr_inject_stmt(ctx *Expr_inject_stmtContext) {}

// EnterExpr_stmt is called when production expr_stmt is entered.
func (s *BaseSyslParserListener) EnterExpr_stmt(ctx *Expr_stmtContext) {}

// ExitExpr_stmt is called when production expr_stmt is exited.
func (s *BaseSyslParserListener) ExitExpr_stmt(ctx *Expr_stmtContext) {}

// EnterTransform_return_type is called when production transform_return_type is entered.
func (s *BaseSyslParserListener) EnterTransform_return_type(ctx *Transform_return_typeContext) {}

// ExitTransform_return_type is called when production transform_return_type is exited.
func (s *BaseSyslParserListener) ExitTransform_return_type(ctx *Transform_return_typeContext) {}

// EnterView_return_type is called when production view_return_type is entered.
func (s *BaseSyslParserListener) EnterView_return_type(ctx *View_return_typeContext) {}

// ExitView_return_type is called when production view_return_type is exited.
func (s *BaseSyslParserListener) ExitView_return_type(ctx *View_return_typeContext) {}

// EnterTransform_scope_var is called when production transform_scope_var is entered.
func (s *BaseSyslParserListener) EnterTransform_scope_var(ctx *Transform_scope_varContext) {}

// ExitTransform_scope_var is called when production transform_scope_var is exited.
func (s *BaseSyslParserListener) ExitTransform_scope_var(ctx *Transform_scope_varContext) {}

// EnterTransform_arg is called when production transform_arg is entered.
func (s *BaseSyslParserListener) EnterTransform_arg(ctx *Transform_argContext) {}

// ExitTransform_arg is called when production transform_arg is exited.
func (s *BaseSyslParserListener) ExitTransform_arg(ctx *Transform_argContext) {}

// EnterTransform is called when production transform is entered.
func (s *BaseSyslParserListener) EnterTransform(ctx *TransformContext) {}

// ExitTransform is called when production transform is exited.
func (s *BaseSyslParserListener) ExitTransform(ctx *TransformContext) {}

// EnterExpr_block is called when production expr_block is entered.
func (s *BaseSyslParserListener) EnterExpr_block(ctx *Expr_blockContext) {}

// ExitExpr_block is called when production expr_block is exited.
func (s *BaseSyslParserListener) ExitExpr_block(ctx *Expr_blockContext) {}

// EnterView_param is called when production view_param is entered.
func (s *BaseSyslParserListener) EnterView_param(ctx *View_paramContext) {}

// ExitView_param is called when production view_param is exited.
func (s *BaseSyslParserListener) ExitView_param(ctx *View_paramContext) {}

// EnterView_params is called when production view_params is entered.
func (s *BaseSyslParserListener) EnterView_params(ctx *View_paramsContext) {}

// ExitView_params is called when production view_params is exited.
func (s *BaseSyslParserListener) ExitView_params(ctx *View_paramsContext) {}

// EnterAbstract_view is called when production abstract_view is entered.
func (s *BaseSyslParserListener) EnterAbstract_view(ctx *Abstract_viewContext) {}

// ExitAbstract_view is called when production abstract_view is exited.
func (s *BaseSyslParserListener) ExitAbstract_view(ctx *Abstract_viewContext) {}

// EnterView is called when production view is entered.
func (s *BaseSyslParserListener) EnterView(ctx *ViewContext) {}

// ExitView is called when production view is exited.
func (s *BaseSyslParserListener) ExitView(ctx *ViewContext) {}

// EnterAlias is called when production alias is entered.
func (s *BaseSyslParserListener) EnterAlias(ctx *AliasContext) {}

// ExitAlias is called when production alias is exited.
func (s *BaseSyslParserListener) ExitAlias(ctx *AliasContext) {}

// EnterApp_decl is called when production app_decl is entered.
func (s *BaseSyslParserListener) EnterApp_decl(ctx *App_declContext) {}

// ExitApp_decl is called when production app_decl is exited.
func (s *BaseSyslParserListener) ExitApp_decl(ctx *App_declContext) {}

// EnterApplication is called when production application is entered.
func (s *BaseSyslParserListener) EnterApplication(ctx *ApplicationContext) {}

// ExitApplication is called when production application is exited.
func (s *BaseSyslParserListener) ExitApplication(ctx *ApplicationContext) {}

// EnterImport_mode is called when production import_mode is entered.
func (s *BaseSyslParserListener) EnterImport_mode(ctx *Import_modeContext) {}

// ExitImport_mode is called when production import_mode is exited.
func (s *BaseSyslParserListener) ExitImport_mode(ctx *Import_modeContext) {}

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
