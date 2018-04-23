// Generated from SyslParser.g4 by ANTLR 4.7.

package parser // SyslParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// SyslParserListener is a complete listener for a parse tree produced by SyslParser.
type SyslParserListener interface {
	antlr.ParseTreeListener

	// EnterModifier is called when entering the modifier production.
	EnterModifier(c *ModifierContext)

	// EnterSize_spec is called when entering the size_spec production.
	EnterSize_spec(c *Size_specContext)

	// EnterModifier_list is called when entering the modifier_list production.
	EnterModifier_list(c *Modifier_listContext)

	// EnterModifiers is called when entering the modifiers production.
	EnterModifiers(c *ModifiersContext)

	// EnterReference is called when entering the reference production.
	EnterReference(c *ReferenceContext)

	// EnterDoc_string is called when entering the doc_string production.
	EnterDoc_string(c *Doc_stringContext)

	// EnterQuoted_string is called when entering the quoted_string production.
	EnterQuoted_string(c *Quoted_stringContext)

	// EnterArray_of_strings is called when entering the array_of_strings production.
	EnterArray_of_strings(c *Array_of_stringsContext)

	// EnterArray_of_arrays is called when entering the array_of_arrays production.
	EnterArray_of_arrays(c *Array_of_arraysContext)

	// EnterNvp is called when entering the nvp production.
	EnterNvp(c *NvpContext)

	// EnterAttributes is called when entering the attributes production.
	EnterAttributes(c *AttributesContext)

	// EnterEntry is called when entering the entry production.
	EnterEntry(c *EntryContext)

	// EnterAttribs_or_modifiers is called when entering the attribs_or_modifiers production.
	EnterAttribs_or_modifiers(c *Attribs_or_modifiersContext)

	// EnterSet_type is called when entering the set_type production.
	EnterSet_type(c *Set_typeContext)

	// EnterCollection_type is called when entering the collection_type production.
	EnterCollection_type(c *Collection_typeContext)

	// EnterUser_defined_type is called when entering the user_defined_type production.
	EnterUser_defined_type(c *User_defined_typeContext)

	// EnterMulti_line_docstring is called when entering the multi_line_docstring production.
	EnterMulti_line_docstring(c *Multi_line_docstringContext)

	// EnterAnnotation_value is called when entering the annotation_value production.
	EnterAnnotation_value(c *Annotation_valueContext)

	// EnterAnnotation is called when entering the annotation production.
	EnterAnnotation(c *AnnotationContext)

	// EnterAnnotations is called when entering the annotations production.
	EnterAnnotations(c *AnnotationsContext)

	// EnterField_type is called when entering the field_type production.
	EnterField_type(c *Field_typeContext)

	// EnterArray_size is called when entering the array_size production.
	EnterArray_size(c *Array_sizeContext)

	// EnterInplace_field is called when entering the inplace_field production.
	EnterInplace_field(c *Inplace_fieldContext)

	// EnterInplace_tuple is called when entering the inplace_tuple production.
	EnterInplace_tuple(c *Inplace_tupleContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterTable is called when entering the table production.
	EnterTable(c *TableContext)

	// EnterPackage_name is called when entering the package_name production.
	EnterPackage_name(c *Package_nameContext)

	// EnterSub_package is called when entering the sub_package production.
	EnterSub_package(c *Sub_packageContext)

	// EnterApp_name is called when entering the app_name production.
	EnterApp_name(c *App_nameContext)

	// EnterName_with_attribs is called when entering the name_with_attribs production.
	EnterName_with_attribs(c *Name_with_attribsContext)

	// EnterModel_name is called when entering the model_name production.
	EnterModel_name(c *Model_nameContext)

	// EnterInplace_table_def is called when entering the inplace_table_def production.
	EnterInplace_table_def(c *Inplace_table_defContext)

	// EnterTable_refs is called when entering the table_refs production.
	EnterTable_refs(c *Table_refsContext)

	// EnterFacade is called when entering the facade production.
	EnterFacade(c *FacadeContext)

	// EnterDocumentation_stmts is called when entering the documentation_stmts production.
	EnterDocumentation_stmts(c *Documentation_stmtsContext)

	// EnterVar_in_curly is called when entering the var_in_curly production.
	EnterVar_in_curly(c *Var_in_curlyContext)

	// EnterQuery_var is called when entering the query_var production.
	EnterQuery_var(c *Query_varContext)

	// EnterQuery_param is called when entering the query_param production.
	EnterQuery_param(c *Query_paramContext)

	// EnterHttp_path_part is called when entering the http_path_part production.
	EnterHttp_path_part(c *Http_path_partContext)

	// EnterHttp_path_var_with_type is called when entering the http_path_var_with_type production.
	EnterHttp_path_var_with_type(c *Http_path_var_with_typeContext)

	// EnterHttp_path_static is called when entering the http_path_static production.
	EnterHttp_path_static(c *Http_path_staticContext)

	// EnterHttp_path_suffix is called when entering the http_path_suffix production.
	EnterHttp_path_suffix(c *Http_path_suffixContext)

	// EnterHttp_path is called when entering the http_path production.
	EnterHttp_path(c *Http_pathContext)

	// EnterEndpoint_name is called when entering the endpoint_name production.
	EnterEndpoint_name(c *Endpoint_nameContext)

	// EnterRet_stmt is called when entering the ret_stmt production.
	EnterRet_stmt(c *Ret_stmtContext)

	// EnterTarget is called when entering the target production.
	EnterTarget(c *TargetContext)

	// EnterTarget_endpoint is called when entering the target_endpoint production.
	EnterTarget_endpoint(c *Target_endpointContext)

	// EnterCall_stmt is called when entering the call_stmt production.
	EnterCall_stmt(c *Call_stmtContext)

	// EnterIf_stmt is called when entering the if_stmt production.
	EnterIf_stmt(c *If_stmtContext)

	// EnterIf_else is called when entering the if_else production.
	EnterIf_else(c *If_elseContext)

	// EnterFor_cond is called when entering the for_cond production.
	EnterFor_cond(c *For_condContext)

	// EnterFor_stmt is called when entering the for_stmt production.
	EnterFor_stmt(c *For_stmtContext)

	// EnterHttp_method_comment is called when entering the http_method_comment production.
	EnterHttp_method_comment(c *Http_method_commentContext)

	// EnterOne_of_case_label is called when entering the one_of_case_label production.
	EnterOne_of_case_label(c *One_of_case_labelContext)

	// EnterOne_of_cases is called when entering the one_of_cases production.
	EnterOne_of_cases(c *One_of_casesContext)

	// EnterOne_of_stmt is called when entering the one_of_stmt production.
	EnterOne_of_stmt(c *One_of_stmtContext)

	// EnterText_stmt is called when entering the text_stmt production.
	EnterText_stmt(c *Text_stmtContext)

	// EnterMixin is called when entering the mixin production.
	EnterMixin(c *MixinContext)

	// EnterParam_list is called when entering the param_list production.
	EnterParam_list(c *Param_listContext)

	// EnterParams is called when entering the params production.
	EnterParams(c *ParamsContext)

	// EnterStatements is called when entering the statements production.
	EnterStatements(c *StatementsContext)

	// EnterMethod_def is called when entering the method_def production.
	EnterMethod_def(c *Method_defContext)

	// EnterShortcut is called when entering the shortcut production.
	EnterShortcut(c *ShortcutContext)

	// EnterSimple_endpoint is called when entering the simple_endpoint production.
	EnterSimple_endpoint(c *Simple_endpointContext)

	// EnterRest_endpoint is called when entering the rest_endpoint production.
	EnterRest_endpoint(c *Rest_endpointContext)

	// EnterCollector_stmt is called when entering the collector_stmt production.
	EnterCollector_stmt(c *Collector_stmtContext)

	// EnterCollector_stmts is called when entering the collector_stmts production.
	EnterCollector_stmts(c *Collector_stmtsContext)

	// EnterCollector is called when entering the collector production.
	EnterCollector(c *CollectorContext)

	// EnterEvent is called when entering the event production.
	EnterEvent(c *EventContext)

	// EnterSubscribe is called when entering the subscribe production.
	EnterSubscribe(c *SubscribeContext)

	// EnterApp_decl is called when entering the app_decl production.
	EnterApp_decl(c *App_declContext)

	// EnterApplication is called when entering the application production.
	EnterApplication(c *ApplicationContext)

	// EnterPath is called when entering the path production.
	EnterPath(c *PathContext)

	// EnterImport_stmt is called when entering the import_stmt production.
	EnterImport_stmt(c *Import_stmtContext)

	// EnterImports_decl is called when entering the imports_decl production.
	EnterImports_decl(c *Imports_declContext)

	// EnterSysl_file is called when entering the sysl_file production.
	EnterSysl_file(c *Sysl_fileContext)

	// ExitModifier is called when exiting the modifier production.
	ExitModifier(c *ModifierContext)

	// ExitSize_spec is called when exiting the size_spec production.
	ExitSize_spec(c *Size_specContext)

	// ExitModifier_list is called when exiting the modifier_list production.
	ExitModifier_list(c *Modifier_listContext)

	// ExitModifiers is called when exiting the modifiers production.
	ExitModifiers(c *ModifiersContext)

	// ExitReference is called when exiting the reference production.
	ExitReference(c *ReferenceContext)

	// ExitDoc_string is called when exiting the doc_string production.
	ExitDoc_string(c *Doc_stringContext)

	// ExitQuoted_string is called when exiting the quoted_string production.
	ExitQuoted_string(c *Quoted_stringContext)

	// ExitArray_of_strings is called when exiting the array_of_strings production.
	ExitArray_of_strings(c *Array_of_stringsContext)

	// ExitArray_of_arrays is called when exiting the array_of_arrays production.
	ExitArray_of_arrays(c *Array_of_arraysContext)

	// ExitNvp is called when exiting the nvp production.
	ExitNvp(c *NvpContext)

	// ExitAttributes is called when exiting the attributes production.
	ExitAttributes(c *AttributesContext)

	// ExitEntry is called when exiting the entry production.
	ExitEntry(c *EntryContext)

	// ExitAttribs_or_modifiers is called when exiting the attribs_or_modifiers production.
	ExitAttribs_or_modifiers(c *Attribs_or_modifiersContext)

	// ExitSet_type is called when exiting the set_type production.
	ExitSet_type(c *Set_typeContext)

	// ExitCollection_type is called when exiting the collection_type production.
	ExitCollection_type(c *Collection_typeContext)

	// ExitUser_defined_type is called when exiting the user_defined_type production.
	ExitUser_defined_type(c *User_defined_typeContext)

	// ExitMulti_line_docstring is called when exiting the multi_line_docstring production.
	ExitMulti_line_docstring(c *Multi_line_docstringContext)

	// ExitAnnotation_value is called when exiting the annotation_value production.
	ExitAnnotation_value(c *Annotation_valueContext)

	// ExitAnnotation is called when exiting the annotation production.
	ExitAnnotation(c *AnnotationContext)

	// ExitAnnotations is called when exiting the annotations production.
	ExitAnnotations(c *AnnotationsContext)

	// ExitField_type is called when exiting the field_type production.
	ExitField_type(c *Field_typeContext)

	// ExitArray_size is called when exiting the array_size production.
	ExitArray_size(c *Array_sizeContext)

	// ExitInplace_field is called when exiting the inplace_field production.
	ExitInplace_field(c *Inplace_fieldContext)

	// ExitInplace_tuple is called when exiting the inplace_tuple production.
	ExitInplace_tuple(c *Inplace_tupleContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitTable is called when exiting the table production.
	ExitTable(c *TableContext)

	// ExitPackage_name is called when exiting the package_name production.
	ExitPackage_name(c *Package_nameContext)

	// ExitSub_package is called when exiting the sub_package production.
	ExitSub_package(c *Sub_packageContext)

	// ExitApp_name is called when exiting the app_name production.
	ExitApp_name(c *App_nameContext)

	// ExitName_with_attribs is called when exiting the name_with_attribs production.
	ExitName_with_attribs(c *Name_with_attribsContext)

	// ExitModel_name is called when exiting the model_name production.
	ExitModel_name(c *Model_nameContext)

	// ExitInplace_table_def is called when exiting the inplace_table_def production.
	ExitInplace_table_def(c *Inplace_table_defContext)

	// ExitTable_refs is called when exiting the table_refs production.
	ExitTable_refs(c *Table_refsContext)

	// ExitFacade is called when exiting the facade production.
	ExitFacade(c *FacadeContext)

	// ExitDocumentation_stmts is called when exiting the documentation_stmts production.
	ExitDocumentation_stmts(c *Documentation_stmtsContext)

	// ExitVar_in_curly is called when exiting the var_in_curly production.
	ExitVar_in_curly(c *Var_in_curlyContext)

	// ExitQuery_var is called when exiting the query_var production.
	ExitQuery_var(c *Query_varContext)

	// ExitQuery_param is called when exiting the query_param production.
	ExitQuery_param(c *Query_paramContext)

	// ExitHttp_path_part is called when exiting the http_path_part production.
	ExitHttp_path_part(c *Http_path_partContext)

	// ExitHttp_path_var_with_type is called when exiting the http_path_var_with_type production.
	ExitHttp_path_var_with_type(c *Http_path_var_with_typeContext)

	// ExitHttp_path_static is called when exiting the http_path_static production.
	ExitHttp_path_static(c *Http_path_staticContext)

	// ExitHttp_path_suffix is called when exiting the http_path_suffix production.
	ExitHttp_path_suffix(c *Http_path_suffixContext)

	// ExitHttp_path is called when exiting the http_path production.
	ExitHttp_path(c *Http_pathContext)

	// ExitEndpoint_name is called when exiting the endpoint_name production.
	ExitEndpoint_name(c *Endpoint_nameContext)

	// ExitRet_stmt is called when exiting the ret_stmt production.
	ExitRet_stmt(c *Ret_stmtContext)

	// ExitTarget is called when exiting the target production.
	ExitTarget(c *TargetContext)

	// ExitTarget_endpoint is called when exiting the target_endpoint production.
	ExitTarget_endpoint(c *Target_endpointContext)

	// ExitCall_stmt is called when exiting the call_stmt production.
	ExitCall_stmt(c *Call_stmtContext)

	// ExitIf_stmt is called when exiting the if_stmt production.
	ExitIf_stmt(c *If_stmtContext)

	// ExitIf_else is called when exiting the if_else production.
	ExitIf_else(c *If_elseContext)

	// ExitFor_cond is called when exiting the for_cond production.
	ExitFor_cond(c *For_condContext)

	// ExitFor_stmt is called when exiting the for_stmt production.
	ExitFor_stmt(c *For_stmtContext)

	// ExitHttp_method_comment is called when exiting the http_method_comment production.
	ExitHttp_method_comment(c *Http_method_commentContext)

	// ExitOne_of_case_label is called when exiting the one_of_case_label production.
	ExitOne_of_case_label(c *One_of_case_labelContext)

	// ExitOne_of_cases is called when exiting the one_of_cases production.
	ExitOne_of_cases(c *One_of_casesContext)

	// ExitOne_of_stmt is called when exiting the one_of_stmt production.
	ExitOne_of_stmt(c *One_of_stmtContext)

	// ExitText_stmt is called when exiting the text_stmt production.
	ExitText_stmt(c *Text_stmtContext)

	// ExitMixin is called when exiting the mixin production.
	ExitMixin(c *MixinContext)

	// ExitParam_list is called when exiting the param_list production.
	ExitParam_list(c *Param_listContext)

	// ExitParams is called when exiting the params production.
	ExitParams(c *ParamsContext)

	// ExitStatements is called when exiting the statements production.
	ExitStatements(c *StatementsContext)

	// ExitMethod_def is called when exiting the method_def production.
	ExitMethod_def(c *Method_defContext)

	// ExitShortcut is called when exiting the shortcut production.
	ExitShortcut(c *ShortcutContext)

	// ExitSimple_endpoint is called when exiting the simple_endpoint production.
	ExitSimple_endpoint(c *Simple_endpointContext)

	// ExitRest_endpoint is called when exiting the rest_endpoint production.
	ExitRest_endpoint(c *Rest_endpointContext)

	// ExitCollector_stmt is called when exiting the collector_stmt production.
	ExitCollector_stmt(c *Collector_stmtContext)

	// ExitCollector_stmts is called when exiting the collector_stmts production.
	ExitCollector_stmts(c *Collector_stmtsContext)

	// ExitCollector is called when exiting the collector production.
	ExitCollector(c *CollectorContext)

	// ExitEvent is called when exiting the event production.
	ExitEvent(c *EventContext)

	// ExitSubscribe is called when exiting the subscribe production.
	ExitSubscribe(c *SubscribeContext)

	// ExitApp_decl is called when exiting the app_decl production.
	ExitApp_decl(c *App_declContext)

	// ExitApplication is called when exiting the application production.
	ExitApplication(c *ApplicationContext)

	// ExitPath is called when exiting the path production.
	ExitPath(c *PathContext)

	// ExitImport_stmt is called when exiting the import_stmt production.
	ExitImport_stmt(c *Import_stmtContext)

	// ExitImports_decl is called when exiting the imports_decl production.
	ExitImports_decl(c *Imports_declContext)

	// ExitSysl_file is called when exiting the sysl_file production.
	ExitSysl_file(c *Sysl_fileContext)
}
