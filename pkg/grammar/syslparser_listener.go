// Generated from pkg/grammar/SyslParser.g4 by ANTLR 4.7.

package parser // SyslParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// SyslParserListener is a complete listener for a parse tree produced by SyslParser.
type SyslParserListener interface {
	antlr.ParseTreeListener

	// EnterModifier is called when entering the modifier production.
	EnterModifier(c *ModifierContext)

	// EnterSize_spec is called when entering the size_spec production.
	EnterSize_spec(c *Size_specContext)

	// EnterName_str is called when entering the name_str production.
	EnterName_str(c *Name_strContext)

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

	// EnterEntry is called when entering the entry production.
	EnterEntry(c *EntryContext)

	// EnterAttribs_or_modifiers is called when entering the attribs_or_modifiers production.
	EnterAttribs_or_modifiers(c *Attribs_or_modifiersContext)

	// EnterUser_defined_type is called when entering the user_defined_type production.
	EnterUser_defined_type(c *User_defined_typeContext)

	// EnterTypes is called when entering the types production.
	EnterTypes(c *TypesContext)

	// EnterSet_of is called when entering the set_of production.
	EnterSet_of(c *Set_ofContext)

	// EnterSet_type is called when entering the set_type production.
	EnterSet_type(c *Set_typeContext)

	// EnterSequence_of is called when entering the sequence_of production.
	EnterSequence_of(c *Sequence_ofContext)

	// EnterSequence_type is called when entering the sequence_type production.
	EnterSequence_type(c *Sequence_typeContext)

	// EnterCollection_type is called when entering the collection_type production.
	EnterCollection_type(c *Collection_typeContext)

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

	// EnterInplace_tuple is called when entering the inplace_tuple production.
	EnterInplace_tuple(c *Inplace_tupleContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterInplace_table is called when entering the inplace_table production.
	EnterInplace_table(c *Inplace_tableContext)

	// EnterTable_stmts is called when entering the table_stmts production.
	EnterTable_stmts(c *Table_stmtsContext)

	// EnterTable_def is called when entering the table_def production.
	EnterTable_def(c *Table_defContext)

	// EnterTable is called when entering the table production.
	EnterTable(c *TableContext)

	// EnterUnion is called when entering the union production.
	EnterUnion(c *UnionContext)

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

	// EnterCall_arg is called when entering the call_arg production.
	EnterCall_arg(c *Call_argContext)

	// EnterCall_args is called when entering the call_args production.
	EnterCall_args(c *Call_argsContext)

	// EnterCall_stmt is called when entering the call_stmt production.
	EnterCall_stmt(c *Call_stmtContext)

	// EnterIf_stmt is called when entering the if_stmt production.
	EnterIf_stmt(c *If_stmtContext)

	// EnterElse_stmt is called when entering the else_stmt production.
	EnterElse_stmt(c *Else_stmtContext)

	// EnterIf_else is called when entering the if_else production.
	EnterIf_else(c *If_elseContext)

	// EnterFor_stmt is called when entering the for_stmt production.
	EnterFor_stmt(c *For_stmtContext)

	// EnterHttp_method_comment is called when entering the http_method_comment production.
	EnterHttp_method_comment(c *Http_method_commentContext)

	// EnterGroup_stmt is called when entering the group_stmt production.
	EnterGroup_stmt(c *Group_stmtContext)

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

	// EnterParam is called when entering the param production.
	EnterParam(c *ParamContext)

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

	// EnterCollector_query_var is called when entering the collector_query_var production.
	EnterCollector_query_var(c *Collector_query_varContext)

	// EnterCollector_query_param is called when entering the collector_query_param production.
	EnterCollector_query_param(c *Collector_query_paramContext)

	// EnterCollector_call_stmt is called when entering the collector_call_stmt production.
	EnterCollector_call_stmt(c *Collector_call_stmtContext)

	// EnterCollector_http_stmt_part is called when entering the collector_http_stmt_part production.
	EnterCollector_http_stmt_part(c *Collector_http_stmt_partContext)

	// EnterCollector_http_stmt_suffix is called when entering the collector_http_stmt_suffix production.
	EnterCollector_http_stmt_suffix(c *Collector_http_stmt_suffixContext)

	// EnterCollector_http_stmt is called when entering the collector_http_stmt production.
	EnterCollector_http_stmt(c *Collector_http_stmtContext)

	// EnterPublisher is called when entering the publisher production.
	EnterPublisher(c *PublisherContext)

	// EnterSubscriber is called when entering the subscriber production.
	EnterSubscriber(c *SubscriberContext)

	// EnterCollector_pubsub_call is called when entering the collector_pubsub_call production.
	EnterCollector_pubsub_call(c *Collector_pubsub_callContext)

	// EnterCollector_action_stmt is called when entering the collector_action_stmt production.
	EnterCollector_action_stmt(c *Collector_action_stmtContext)

	// EnterCollector_stmts is called when entering the collector_stmts production.
	EnterCollector_stmts(c *Collector_stmtsContext)

	// EnterCollector is called when entering the collector production.
	EnterCollector(c *CollectorContext)

	// EnterEvent is called when entering the event production.
	EnterEvent(c *EventContext)

	// EnterSubscribe is called when entering the subscribe production.
	EnterSubscribe(c *SubscribeContext)

	// EnterView_type_spec is called when entering the view_type_spec production.
	EnterView_type_spec(c *View_type_specContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterExpr_table_of_op is called when entering the expr_table_of_op production.
	EnterExpr_table_of_op(c *Expr_table_of_opContext)

	// EnterFunc_arg is called when entering the func_arg production.
	EnterFunc_arg(c *Func_argContext)

	// EnterFunc_args is called when entering the func_args production.
	EnterFunc_args(c *Func_argsContext)

	// EnterExpr_func is called when entering the expr_func production.
	EnterExpr_func(c *Expr_funcContext)

	// EnterRank_expr is called when entering the rank_expr production.
	EnterRank_expr(c *Rank_exprContext)

	// EnterRank_expr_list is called when entering the rank_expr_list production.
	EnterRank_expr_list(c *Rank_expr_listContext)

	// EnterExpr_rank_func is called when entering the expr_rank_func production.
	EnterExpr_rank_func(c *Expr_rank_funcContext)

	// EnterExpr_agg_func is called when entering the expr_agg_func production.
	EnterExpr_agg_func(c *Expr_agg_funcContext)

	// EnterE_scope_var is called when entering the e_scope_var production.
	EnterE_scope_var(c *E_scope_varContext)

	// EnterFirst_func_target is called when entering the first_func_target production.
	EnterFirst_func_target(c *First_func_targetContext)

	// EnterExpr_first_func is called when entering the expr_first_func production.
	EnterExpr_first_func(c *Expr_first_funcContext)

	// EnterE_single_arg_func is called when entering the e_single_arg_func production.
	EnterE_single_arg_func(c *E_single_arg_funcContext)

	// EnterExpr_single_arg_func is called when entering the expr_single_arg_func production.
	EnterExpr_single_arg_func(c *Expr_single_arg_funcContext)

	// EnterExpr_any_func is called when entering the expr_any_func production.
	EnterExpr_any_func(c *Expr_any_funcContext)

	// EnterExpr_single_or_null is called when entering the expr_single_or_null production.
	EnterExpr_single_or_null(c *Expr_single_or_nullContext)

	// EnterExpr_snapshot is called when entering the expr_snapshot production.
	EnterExpr_snapshot(c *Expr_snapshotContext)

	// EnterExpr_count is called when entering the expr_count production.
	EnterExpr_count(c *Expr_countContext)

	// EnterExpr_navigate_attr is called when entering the expr_navigate_attr production.
	EnterExpr_navigate_attr(c *Expr_navigate_attrContext)

	// EnterExpr_navigate is called when entering the expr_navigate production.
	EnterExpr_navigate(c *Expr_navigateContext)

	// EnterMatching_rhs is called when entering the matching_rhs production.
	EnterMatching_rhs(c *Matching_rhsContext)

	// EnterSquiggly_args is called when entering the squiggly_args production.
	EnterSquiggly_args(c *Squiggly_argsContext)

	// EnterExpr_matching is called when entering the expr_matching production.
	EnterExpr_matching(c *Expr_matchingContext)

	// EnterRelop is called when entering the relop production.
	EnterRelop(c *RelopContext)

	// EnterList_item is called when entering the list_item production.
	EnterList_item(c *List_itemContext)

	// EnterExpr_list is called when entering the expr_list production.
	EnterExpr_list(c *Expr_listContext)

	// EnterExpr_set is called when entering the expr_set production.
	EnterExpr_set(c *Expr_setContext)

	// EnterEmpty_tuple is called when entering the empty_tuple production.
	EnterEmpty_tuple(c *Empty_tupleContext)

	// EnterAtom_dot_relop is called when entering the atom_dot_relop production.
	EnterAtom_dot_relop(c *Atom_dot_relopContext)

	// EnterAtomT_implied_dot is called when entering the atomT_implied_dot production.
	EnterAtomT_implied_dot(c *AtomT_implied_dotContext)

	// EnterAtomT_name is called when entering the atomT_name production.
	EnterAtomT_name(c *AtomT_nameContext)

	// EnterAtomT_paren is called when entering the atomT_paren production.
	EnterAtomT_paren(c *AtomT_parenContext)

	// EnterExpr_atom_list is called when entering the expr_atom_list production.
	EnterExpr_atom_list(c *Expr_atom_listContext)

	// EnterAtomT is called when entering the atomT production.
	EnterAtomT(c *AtomTContext)

	// EnterAtom is called when entering the atom production.
	EnterAtom(c *AtomContext)

	// EnterPowerT is called when entering the powerT production.
	EnterPowerT(c *PowerTContext)

	// EnterPower is called when entering the power production.
	EnterPower(c *PowerContext)

	// EnterUnaryTerm is called when entering the unaryTerm production.
	EnterUnaryTerm(c *UnaryTermContext)

	// EnterTermT is called when entering the termT production.
	EnterTermT(c *TermTContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterBinexprT is called when entering the binexprT production.
	EnterBinexprT(c *BinexprTContext)

	// EnterBinexpr is called when entering the binexpr production.
	EnterBinexpr(c *BinexprContext)

	// EnterE_compare_ops is called when entering the e_compare_ops production.
	EnterE_compare_ops(c *E_compare_opsContext)

	// EnterExpr_rel is called when entering the expr_rel production.
	EnterExpr_rel(c *Expr_relContext)

	// EnterExpr_bitand is called when entering the expr_bitand production.
	EnterExpr_bitand(c *Expr_bitandContext)

	// EnterExpr_bitxor is called when entering the expr_bitxor production.
	EnterExpr_bitxor(c *Expr_bitxorContext)

	// EnterExpr_bitor is called when entering the expr_bitor production.
	EnterExpr_bitor(c *Expr_bitorContext)

	// EnterExpr_and is called when entering the expr_and production.
	EnterExpr_and(c *Expr_andContext)

	// EnterExpr_or is called when entering the expr_or production.
	EnterExpr_or(c *Expr_orContext)

	// EnterExpr_but_not is called when entering the expr_but_not production.
	EnterExpr_but_not(c *Expr_but_notContext)

	// EnterExpr_coalesce is called when entering the expr_coalesce production.
	EnterExpr_coalesce(c *Expr_coalesceContext)

	// EnterIf_one_liner is called when entering the if_one_liner production.
	EnterIf_one_liner(c *If_one_linerContext)

	// EnterElse_block_stmt is called when entering the else_block_stmt production.
	EnterElse_block_stmt(c *Else_block_stmtContext)

	// EnterControl_item is called when entering the control_item production.
	EnterControl_item(c *Control_itemContext)

	// EnterIf_controls is called when entering the if_controls production.
	EnterIf_controls(c *If_controlsContext)

	// EnterCond_block is called when entering the cond_block production.
	EnterCond_block(c *Cond_blockContext)

	// EnterFinal_else is called when entering the final_else production.
	EnterFinal_else(c *Final_elseContext)

	// EnterIfvar is called when entering the ifvar production.
	EnterIfvar(c *IfvarContext)

	// EnterIf_multiple_lines is called when entering the if_multiple_lines production.
	EnterIf_multiple_lines(c *If_multiple_linesContext)

	// EnterExpr_if_else is called when entering the expr_if_else production.
	EnterExpr_if_else(c *Expr_if_elseContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterExpr_assign is called when entering the expr_assign production.
	EnterExpr_assign(c *Expr_assignContext)

	// EnterExpr_simple_assign is called when entering the expr_simple_assign production.
	EnterExpr_simple_assign(c *Expr_simple_assignContext)

	// EnterExpr_let_statement is called when entering the expr_let_statement production.
	EnterExpr_let_statement(c *Expr_let_statementContext)

	// EnterExpr_table_of_statement is called when entering the expr_table_of_statement production.
	EnterExpr_table_of_statement(c *Expr_table_of_statementContext)

	// EnterExpr_dot_assign is called when entering the expr_dot_assign production.
	EnterExpr_dot_assign(c *Expr_dot_assignContext)

	// EnterExpr_statement_no_nl is called when entering the expr_statement_no_nl production.
	EnterExpr_statement_no_nl(c *Expr_statement_no_nlContext)

	// EnterExpr_statement is called when entering the expr_statement production.
	EnterExpr_statement(c *Expr_statementContext)

	// EnterExpr_inject_stmt is called when entering the expr_inject_stmt production.
	EnterExpr_inject_stmt(c *Expr_inject_stmtContext)

	// EnterExpr_stmt is called when entering the expr_stmt production.
	EnterExpr_stmt(c *Expr_stmtContext)

	// EnterTransform_return_type is called when entering the transform_return_type production.
	EnterTransform_return_type(c *Transform_return_typeContext)

	// EnterView_return_type is called when entering the view_return_type production.
	EnterView_return_type(c *View_return_typeContext)

	// EnterTransform_scope_var is called when entering the transform_scope_var production.
	EnterTransform_scope_var(c *Transform_scope_varContext)

	// EnterTransform_arg is called when entering the transform_arg production.
	EnterTransform_arg(c *Transform_argContext)

	// EnterTransform is called when entering the transform production.
	EnterTransform(c *TransformContext)

	// EnterExpr_block is called when entering the expr_block production.
	EnterExpr_block(c *Expr_blockContext)

	// EnterView_param is called when entering the view_param production.
	EnterView_param(c *View_paramContext)

	// EnterView_params is called when entering the view_params production.
	EnterView_params(c *View_paramsContext)

	// EnterAbstract_view is called when entering the abstract_view production.
	EnterAbstract_view(c *Abstract_viewContext)

	// EnterView is called when entering the view production.
	EnterView(c *ViewContext)

	// EnterAlias is called when entering the alias production.
	EnterAlias(c *AliasContext)

	// EnterApp_decl is called when entering the app_decl production.
	EnterApp_decl(c *App_declContext)

	// EnterApplication is called when entering the application production.
	EnterApplication(c *ApplicationContext)

	// EnterImport_mode is called when entering the import_mode production.
	EnterImport_mode(c *Import_modeContext)

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

	// ExitName_str is called when exiting the name_str production.
	ExitName_str(c *Name_strContext)

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

	// ExitEntry is called when exiting the entry production.
	ExitEntry(c *EntryContext)

	// ExitAttribs_or_modifiers is called when exiting the attribs_or_modifiers production.
	ExitAttribs_or_modifiers(c *Attribs_or_modifiersContext)

	// ExitUser_defined_type is called when exiting the user_defined_type production.
	ExitUser_defined_type(c *User_defined_typeContext)

	// ExitTypes is called when exiting the types production.
	ExitTypes(c *TypesContext)

	// ExitSet_of is called when exiting the set_of production.
	ExitSet_of(c *Set_ofContext)

	// ExitSet_type is called when exiting the set_type production.
	ExitSet_type(c *Set_typeContext)

	// ExitSequence_of is called when exiting the sequence_of production.
	ExitSequence_of(c *Sequence_ofContext)

	// ExitSequence_type is called when exiting the sequence_type production.
	ExitSequence_type(c *Sequence_typeContext)

	// ExitCollection_type is called when exiting the collection_type production.
	ExitCollection_type(c *Collection_typeContext)

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

	// ExitInplace_tuple is called when exiting the inplace_tuple production.
	ExitInplace_tuple(c *Inplace_tupleContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitInplace_table is called when exiting the inplace_table production.
	ExitInplace_table(c *Inplace_tableContext)

	// ExitTable_stmts is called when exiting the table_stmts production.
	ExitTable_stmts(c *Table_stmtsContext)

	// ExitTable_def is called when exiting the table_def production.
	ExitTable_def(c *Table_defContext)

	// ExitTable is called when exiting the table production.
	ExitTable(c *TableContext)

	// ExitUnion is called when exiting the union production.
	ExitUnion(c *UnionContext)

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

	// ExitCall_arg is called when exiting the call_arg production.
	ExitCall_arg(c *Call_argContext)

	// ExitCall_args is called when exiting the call_args production.
	ExitCall_args(c *Call_argsContext)

	// ExitCall_stmt is called when exiting the call_stmt production.
	ExitCall_stmt(c *Call_stmtContext)

	// ExitIf_stmt is called when exiting the if_stmt production.
	ExitIf_stmt(c *If_stmtContext)

	// ExitElse_stmt is called when exiting the else_stmt production.
	ExitElse_stmt(c *Else_stmtContext)

	// ExitIf_else is called when exiting the if_else production.
	ExitIf_else(c *If_elseContext)

	// ExitFor_stmt is called when exiting the for_stmt production.
	ExitFor_stmt(c *For_stmtContext)

	// ExitHttp_method_comment is called when exiting the http_method_comment production.
	ExitHttp_method_comment(c *Http_method_commentContext)

	// ExitGroup_stmt is called when exiting the group_stmt production.
	ExitGroup_stmt(c *Group_stmtContext)

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

	// ExitParam is called when exiting the param production.
	ExitParam(c *ParamContext)

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

	// ExitCollector_query_var is called when exiting the collector_query_var production.
	ExitCollector_query_var(c *Collector_query_varContext)

	// ExitCollector_query_param is called when exiting the collector_query_param production.
	ExitCollector_query_param(c *Collector_query_paramContext)

	// ExitCollector_call_stmt is called when exiting the collector_call_stmt production.
	ExitCollector_call_stmt(c *Collector_call_stmtContext)

	// ExitCollector_http_stmt_part is called when exiting the collector_http_stmt_part production.
	ExitCollector_http_stmt_part(c *Collector_http_stmt_partContext)

	// ExitCollector_http_stmt_suffix is called when exiting the collector_http_stmt_suffix production.
	ExitCollector_http_stmt_suffix(c *Collector_http_stmt_suffixContext)

	// ExitCollector_http_stmt is called when exiting the collector_http_stmt production.
	ExitCollector_http_stmt(c *Collector_http_stmtContext)

	// ExitPublisher is called when exiting the publisher production.
	ExitPublisher(c *PublisherContext)

	// ExitSubscriber is called when exiting the subscriber production.
	ExitSubscriber(c *SubscriberContext)

	// ExitCollector_pubsub_call is called when exiting the collector_pubsub_call production.
	ExitCollector_pubsub_call(c *Collector_pubsub_callContext)

	// ExitCollector_action_stmt is called when exiting the collector_action_stmt production.
	ExitCollector_action_stmt(c *Collector_action_stmtContext)

	// ExitCollector_stmts is called when exiting the collector_stmts production.
	ExitCollector_stmts(c *Collector_stmtsContext)

	// ExitCollector is called when exiting the collector production.
	ExitCollector(c *CollectorContext)

	// ExitEvent is called when exiting the event production.
	ExitEvent(c *EventContext)

	// ExitSubscribe is called when exiting the subscribe production.
	ExitSubscribe(c *SubscribeContext)

	// ExitView_type_spec is called when exiting the view_type_spec production.
	ExitView_type_spec(c *View_type_specContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitExpr_table_of_op is called when exiting the expr_table_of_op production.
	ExitExpr_table_of_op(c *Expr_table_of_opContext)

	// ExitFunc_arg is called when exiting the func_arg production.
	ExitFunc_arg(c *Func_argContext)

	// ExitFunc_args is called when exiting the func_args production.
	ExitFunc_args(c *Func_argsContext)

	// ExitExpr_func is called when exiting the expr_func production.
	ExitExpr_func(c *Expr_funcContext)

	// ExitRank_expr is called when exiting the rank_expr production.
	ExitRank_expr(c *Rank_exprContext)

	// ExitRank_expr_list is called when exiting the rank_expr_list production.
	ExitRank_expr_list(c *Rank_expr_listContext)

	// ExitExpr_rank_func is called when exiting the expr_rank_func production.
	ExitExpr_rank_func(c *Expr_rank_funcContext)

	// ExitExpr_agg_func is called when exiting the expr_agg_func production.
	ExitExpr_agg_func(c *Expr_agg_funcContext)

	// ExitE_scope_var is called when exiting the e_scope_var production.
	ExitE_scope_var(c *E_scope_varContext)

	// ExitFirst_func_target is called when exiting the first_func_target production.
	ExitFirst_func_target(c *First_func_targetContext)

	// ExitExpr_first_func is called when exiting the expr_first_func production.
	ExitExpr_first_func(c *Expr_first_funcContext)

	// ExitE_single_arg_func is called when exiting the e_single_arg_func production.
	ExitE_single_arg_func(c *E_single_arg_funcContext)

	// ExitExpr_single_arg_func is called when exiting the expr_single_arg_func production.
	ExitExpr_single_arg_func(c *Expr_single_arg_funcContext)

	// ExitExpr_any_func is called when exiting the expr_any_func production.
	ExitExpr_any_func(c *Expr_any_funcContext)

	// ExitExpr_single_or_null is called when exiting the expr_single_or_null production.
	ExitExpr_single_or_null(c *Expr_single_or_nullContext)

	// ExitExpr_snapshot is called when exiting the expr_snapshot production.
	ExitExpr_snapshot(c *Expr_snapshotContext)

	// ExitExpr_count is called when exiting the expr_count production.
	ExitExpr_count(c *Expr_countContext)

	// ExitExpr_navigate_attr is called when exiting the expr_navigate_attr production.
	ExitExpr_navigate_attr(c *Expr_navigate_attrContext)

	// ExitExpr_navigate is called when exiting the expr_navigate production.
	ExitExpr_navigate(c *Expr_navigateContext)

	// ExitMatching_rhs is called when exiting the matching_rhs production.
	ExitMatching_rhs(c *Matching_rhsContext)

	// ExitSquiggly_args is called when exiting the squiggly_args production.
	ExitSquiggly_args(c *Squiggly_argsContext)

	// ExitExpr_matching is called when exiting the expr_matching production.
	ExitExpr_matching(c *Expr_matchingContext)

	// ExitRelop is called when exiting the relop production.
	ExitRelop(c *RelopContext)

	// ExitList_item is called when exiting the list_item production.
	ExitList_item(c *List_itemContext)

	// ExitExpr_list is called when exiting the expr_list production.
	ExitExpr_list(c *Expr_listContext)

	// ExitExpr_set is called when exiting the expr_set production.
	ExitExpr_set(c *Expr_setContext)

	// ExitEmpty_tuple is called when exiting the empty_tuple production.
	ExitEmpty_tuple(c *Empty_tupleContext)

	// ExitAtom_dot_relop is called when exiting the atom_dot_relop production.
	ExitAtom_dot_relop(c *Atom_dot_relopContext)

	// ExitAtomT_implied_dot is called when exiting the atomT_implied_dot production.
	ExitAtomT_implied_dot(c *AtomT_implied_dotContext)

	// ExitAtomT_name is called when exiting the atomT_name production.
	ExitAtomT_name(c *AtomT_nameContext)

	// ExitAtomT_paren is called when exiting the atomT_paren production.
	ExitAtomT_paren(c *AtomT_parenContext)

	// ExitExpr_atom_list is called when exiting the expr_atom_list production.
	ExitExpr_atom_list(c *Expr_atom_listContext)

	// ExitAtomT is called when exiting the atomT production.
	ExitAtomT(c *AtomTContext)

	// ExitAtom is called when exiting the atom production.
	ExitAtom(c *AtomContext)

	// ExitPowerT is called when exiting the powerT production.
	ExitPowerT(c *PowerTContext)

	// ExitPower is called when exiting the power production.
	ExitPower(c *PowerContext)

	// ExitUnaryTerm is called when exiting the unaryTerm production.
	ExitUnaryTerm(c *UnaryTermContext)

	// ExitTermT is called when exiting the termT production.
	ExitTermT(c *TermTContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitBinexprT is called when exiting the binexprT production.
	ExitBinexprT(c *BinexprTContext)

	// ExitBinexpr is called when exiting the binexpr production.
	ExitBinexpr(c *BinexprContext)

	// ExitE_compare_ops is called when exiting the e_compare_ops production.
	ExitE_compare_ops(c *E_compare_opsContext)

	// ExitExpr_rel is called when exiting the expr_rel production.
	ExitExpr_rel(c *Expr_relContext)

	// ExitExpr_bitand is called when exiting the expr_bitand production.
	ExitExpr_bitand(c *Expr_bitandContext)

	// ExitExpr_bitxor is called when exiting the expr_bitxor production.
	ExitExpr_bitxor(c *Expr_bitxorContext)

	// ExitExpr_bitor is called when exiting the expr_bitor production.
	ExitExpr_bitor(c *Expr_bitorContext)

	// ExitExpr_and is called when exiting the expr_and production.
	ExitExpr_and(c *Expr_andContext)

	// ExitExpr_or is called when exiting the expr_or production.
	ExitExpr_or(c *Expr_orContext)

	// ExitExpr_but_not is called when exiting the expr_but_not production.
	ExitExpr_but_not(c *Expr_but_notContext)

	// ExitExpr_coalesce is called when exiting the expr_coalesce production.
	ExitExpr_coalesce(c *Expr_coalesceContext)

	// ExitIf_one_liner is called when exiting the if_one_liner production.
	ExitIf_one_liner(c *If_one_linerContext)

	// ExitElse_block_stmt is called when exiting the else_block_stmt production.
	ExitElse_block_stmt(c *Else_block_stmtContext)

	// ExitControl_item is called when exiting the control_item production.
	ExitControl_item(c *Control_itemContext)

	// ExitIf_controls is called when exiting the if_controls production.
	ExitIf_controls(c *If_controlsContext)

	// ExitCond_block is called when exiting the cond_block production.
	ExitCond_block(c *Cond_blockContext)

	// ExitFinal_else is called when exiting the final_else production.
	ExitFinal_else(c *Final_elseContext)

	// ExitIfvar is called when exiting the ifvar production.
	ExitIfvar(c *IfvarContext)

	// ExitIf_multiple_lines is called when exiting the if_multiple_lines production.
	ExitIf_multiple_lines(c *If_multiple_linesContext)

	// ExitExpr_if_else is called when exiting the expr_if_else production.
	ExitExpr_if_else(c *Expr_if_elseContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitExpr_assign is called when exiting the expr_assign production.
	ExitExpr_assign(c *Expr_assignContext)

	// ExitExpr_simple_assign is called when exiting the expr_simple_assign production.
	ExitExpr_simple_assign(c *Expr_simple_assignContext)

	// ExitExpr_let_statement is called when exiting the expr_let_statement production.
	ExitExpr_let_statement(c *Expr_let_statementContext)

	// ExitExpr_table_of_statement is called when exiting the expr_table_of_statement production.
	ExitExpr_table_of_statement(c *Expr_table_of_statementContext)

	// ExitExpr_dot_assign is called when exiting the expr_dot_assign production.
	ExitExpr_dot_assign(c *Expr_dot_assignContext)

	// ExitExpr_statement_no_nl is called when exiting the expr_statement_no_nl production.
	ExitExpr_statement_no_nl(c *Expr_statement_no_nlContext)

	// ExitExpr_statement is called when exiting the expr_statement production.
	ExitExpr_statement(c *Expr_statementContext)

	// ExitExpr_inject_stmt is called when exiting the expr_inject_stmt production.
	ExitExpr_inject_stmt(c *Expr_inject_stmtContext)

	// ExitExpr_stmt is called when exiting the expr_stmt production.
	ExitExpr_stmt(c *Expr_stmtContext)

	// ExitTransform_return_type is called when exiting the transform_return_type production.
	ExitTransform_return_type(c *Transform_return_typeContext)

	// ExitView_return_type is called when exiting the view_return_type production.
	ExitView_return_type(c *View_return_typeContext)

	// ExitTransform_scope_var is called when exiting the transform_scope_var production.
	ExitTransform_scope_var(c *Transform_scope_varContext)

	// ExitTransform_arg is called when exiting the transform_arg production.
	ExitTransform_arg(c *Transform_argContext)

	// ExitTransform is called when exiting the transform production.
	ExitTransform(c *TransformContext)

	// ExitExpr_block is called when exiting the expr_block production.
	ExitExpr_block(c *Expr_blockContext)

	// ExitView_param is called when exiting the view_param production.
	ExitView_param(c *View_paramContext)

	// ExitView_params is called when exiting the view_params production.
	ExitView_params(c *View_paramsContext)

	// ExitAbstract_view is called when exiting the abstract_view production.
	ExitAbstract_view(c *Abstract_viewContext)

	// ExitView is called when exiting the view production.
	ExitView(c *ViewContext)

	// ExitAlias is called when exiting the alias production.
	ExitAlias(c *AliasContext)

	// ExitApp_decl is called when exiting the app_decl production.
	ExitApp_decl(c *App_declContext)

	// ExitApplication is called when exiting the application production.
	ExitApplication(c *ApplicationContext)

	// ExitImport_mode is called when exiting the import_mode production.
	ExitImport_mode(c *Import_modeContext)

	// ExitImport_stmt is called when exiting the import_stmt production.
	ExitImport_stmt(c *Import_stmtContext)

	// ExitImports_decl is called when exiting the imports_decl production.
	ExitImports_decl(c *Imports_declContext)

	// ExitSysl_file is called when exiting the sysl_file production.
	ExitSysl_file(c *Sysl_fileContext)
}
