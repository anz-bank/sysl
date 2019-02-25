// Generated from sysl_js/SyslParser.g4 by ANTLR 4.7
// jshint ignore: start
var antlr4 = require('antlr4/index');

// This class defines a complete listener for a parse tree produced by SyslParser.
function SyslParserListener() {
	antlr4.tree.ParseTreeListener.call(this);
	return this;
}

SyslParserListener.prototype = Object.create(antlr4.tree.ParseTreeListener.prototype);
SyslParserListener.prototype.constructor = SyslParserListener;

// Enter a parse tree produced by SyslParser#modifier.
SyslParserListener.prototype.enterModifier = function(ctx) {
};

// Exit a parse tree produced by SyslParser#modifier.
SyslParserListener.prototype.exitModifier = function(ctx) {
};


// Enter a parse tree produced by SyslParser#size_spec.
SyslParserListener.prototype.enterSize_spec = function(ctx) {
};

// Exit a parse tree produced by SyslParser#size_spec.
SyslParserListener.prototype.exitSize_spec = function(ctx) {
};


// Enter a parse tree produced by SyslParser#modifier_list.
SyslParserListener.prototype.enterModifier_list = function(ctx) {
};

// Exit a parse tree produced by SyslParser#modifier_list.
SyslParserListener.prototype.exitModifier_list = function(ctx) {
};


// Enter a parse tree produced by SyslParser#modifiers.
SyslParserListener.prototype.enterModifiers = function(ctx) {
};

// Exit a parse tree produced by SyslParser#modifiers.
SyslParserListener.prototype.exitModifiers = function(ctx) {
};


// Enter a parse tree produced by SyslParser#name_str.
SyslParserListener.prototype.enterName_str = function(ctx) {
};

// Exit a parse tree produced by SyslParser#name_str.
SyslParserListener.prototype.exitName_str = function(ctx) {
};


// Enter a parse tree produced by SyslParser#reference.
SyslParserListener.prototype.enterReference = function(ctx) {
};

// Exit a parse tree produced by SyslParser#reference.
SyslParserListener.prototype.exitReference = function(ctx) {
};


// Enter a parse tree produced by SyslParser#doc_string.
SyslParserListener.prototype.enterDoc_string = function(ctx) {
};

// Exit a parse tree produced by SyslParser#doc_string.
SyslParserListener.prototype.exitDoc_string = function(ctx) {
};


// Enter a parse tree produced by SyslParser#quoted_string.
SyslParserListener.prototype.enterQuoted_string = function(ctx) {
};

// Exit a parse tree produced by SyslParser#quoted_string.
SyslParserListener.prototype.exitQuoted_string = function(ctx) {
};


// Enter a parse tree produced by SyslParser#array_of_strings.
SyslParserListener.prototype.enterArray_of_strings = function(ctx) {
};

// Exit a parse tree produced by SyslParser#array_of_strings.
SyslParserListener.prototype.exitArray_of_strings = function(ctx) {
};


// Enter a parse tree produced by SyslParser#array_of_arrays.
SyslParserListener.prototype.enterArray_of_arrays = function(ctx) {
};

// Exit a parse tree produced by SyslParser#array_of_arrays.
SyslParserListener.prototype.exitArray_of_arrays = function(ctx) {
};


// Enter a parse tree produced by SyslParser#nvp.
SyslParserListener.prototype.enterNvp = function(ctx) {
};

// Exit a parse tree produced by SyslParser#nvp.
SyslParserListener.prototype.exitNvp = function(ctx) {
};


// Enter a parse tree produced by SyslParser#attributes.
SyslParserListener.prototype.enterAttributes = function(ctx) {
};

// Exit a parse tree produced by SyslParser#attributes.
SyslParserListener.prototype.exitAttributes = function(ctx) {
};


// Enter a parse tree produced by SyslParser#entry.
SyslParserListener.prototype.enterEntry = function(ctx) {
};

// Exit a parse tree produced by SyslParser#entry.
SyslParserListener.prototype.exitEntry = function(ctx) {
};


// Enter a parse tree produced by SyslParser#attribs_or_modifiers.
SyslParserListener.prototype.enterAttribs_or_modifiers = function(ctx) {
};

// Exit a parse tree produced by SyslParser#attribs_or_modifiers.
SyslParserListener.prototype.exitAttribs_or_modifiers = function(ctx) {
};


// Enter a parse tree produced by SyslParser#user_defined_type.
SyslParserListener.prototype.enterUser_defined_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#user_defined_type.
SyslParserListener.prototype.exitUser_defined_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#types.
SyslParserListener.prototype.enterTypes = function(ctx) {
};

// Exit a parse tree produced by SyslParser#types.
SyslParserListener.prototype.exitTypes = function(ctx) {
};


// Enter a parse tree produced by SyslParser#set_of.
SyslParserListener.prototype.enterSet_of = function(ctx) {
};

// Exit a parse tree produced by SyslParser#set_of.
SyslParserListener.prototype.exitSet_of = function(ctx) {
};


// Enter a parse tree produced by SyslParser#set_type.
SyslParserListener.prototype.enterSet_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#set_type.
SyslParserListener.prototype.exitSet_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#sequence_of.
SyslParserListener.prototype.enterSequence_of = function(ctx) {
};

// Exit a parse tree produced by SyslParser#sequence_of.
SyslParserListener.prototype.exitSequence_of = function(ctx) {
};


// Enter a parse tree produced by SyslParser#sequence_type.
SyslParserListener.prototype.enterSequence_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#sequence_type.
SyslParserListener.prototype.exitSequence_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collection_type.
SyslParserListener.prototype.enterCollection_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collection_type.
SyslParserListener.prototype.exitCollection_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#multi_line_docstring.
SyslParserListener.prototype.enterMulti_line_docstring = function(ctx) {
};

// Exit a parse tree produced by SyslParser#multi_line_docstring.
SyslParserListener.prototype.exitMulti_line_docstring = function(ctx) {
};


// Enter a parse tree produced by SyslParser#annotation_value.
SyslParserListener.prototype.enterAnnotation_value = function(ctx) {
};

// Exit a parse tree produced by SyslParser#annotation_value.
SyslParserListener.prototype.exitAnnotation_value = function(ctx) {
};


// Enter a parse tree produced by SyslParser#annotation.
SyslParserListener.prototype.enterAnnotation = function(ctx) {
};

// Exit a parse tree produced by SyslParser#annotation.
SyslParserListener.prototype.exitAnnotation = function(ctx) {
};


// Enter a parse tree produced by SyslParser#annotations.
SyslParserListener.prototype.enterAnnotations = function(ctx) {
};

// Exit a parse tree produced by SyslParser#annotations.
SyslParserListener.prototype.exitAnnotations = function(ctx) {
};


// Enter a parse tree produced by SyslParser#field_type.
SyslParserListener.prototype.enterField_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#field_type.
SyslParserListener.prototype.exitField_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#array_size.
SyslParserListener.prototype.enterArray_size = function(ctx) {
};

// Exit a parse tree produced by SyslParser#array_size.
SyslParserListener.prototype.exitArray_size = function(ctx) {
};


// Enter a parse tree produced by SyslParser#inplace_tuple.
SyslParserListener.prototype.enterInplace_tuple = function(ctx) {
};

// Exit a parse tree produced by SyslParser#inplace_tuple.
SyslParserListener.prototype.exitInplace_tuple = function(ctx) {
};


// Enter a parse tree produced by SyslParser#field.
SyslParserListener.prototype.enterField = function(ctx) {
};

// Exit a parse tree produced by SyslParser#field.
SyslParserListener.prototype.exitField = function(ctx) {
};


// Enter a parse tree produced by SyslParser#inplace_table.
SyslParserListener.prototype.enterInplace_table = function(ctx) {
};

// Exit a parse tree produced by SyslParser#inplace_table.
SyslParserListener.prototype.exitInplace_table = function(ctx) {
};


// Enter a parse tree produced by SyslParser#table_stmts.
SyslParserListener.prototype.enterTable_stmts = function(ctx) {
};

// Exit a parse tree produced by SyslParser#table_stmts.
SyslParserListener.prototype.exitTable_stmts = function(ctx) {
};


// Enter a parse tree produced by SyslParser#table_def.
SyslParserListener.prototype.enterTable_def = function(ctx) {
};

// Exit a parse tree produced by SyslParser#table_def.
SyslParserListener.prototype.exitTable_def = function(ctx) {
};


// Enter a parse tree produced by SyslParser#table.
SyslParserListener.prototype.enterTable = function(ctx) {
};

// Exit a parse tree produced by SyslParser#table.
SyslParserListener.prototype.exitTable = function(ctx) {
};


// Enter a parse tree produced by SyslParser#union.
SyslParserListener.prototype.enterUnion = function(ctx) {
};

// Exit a parse tree produced by SyslParser#union.
SyslParserListener.prototype.exitUnion = function(ctx) {
};


// Enter a parse tree produced by SyslParser#package_name.
SyslParserListener.prototype.enterPackage_name = function(ctx) {
};

// Exit a parse tree produced by SyslParser#package_name.
SyslParserListener.prototype.exitPackage_name = function(ctx) {
};


// Enter a parse tree produced by SyslParser#sub_package.
SyslParserListener.prototype.enterSub_package = function(ctx) {
};

// Exit a parse tree produced by SyslParser#sub_package.
SyslParserListener.prototype.exitSub_package = function(ctx) {
};


// Enter a parse tree produced by SyslParser#app_name.
SyslParserListener.prototype.enterApp_name = function(ctx) {
};

// Exit a parse tree produced by SyslParser#app_name.
SyslParserListener.prototype.exitApp_name = function(ctx) {
};


// Enter a parse tree produced by SyslParser#name_with_attribs.
SyslParserListener.prototype.enterName_with_attribs = function(ctx) {
};

// Exit a parse tree produced by SyslParser#name_with_attribs.
SyslParserListener.prototype.exitName_with_attribs = function(ctx) {
};


// Enter a parse tree produced by SyslParser#model_name.
SyslParserListener.prototype.enterModel_name = function(ctx) {
};

// Exit a parse tree produced by SyslParser#model_name.
SyslParserListener.prototype.exitModel_name = function(ctx) {
};


// Enter a parse tree produced by SyslParser#inplace_table_def.
SyslParserListener.prototype.enterInplace_table_def = function(ctx) {
};

// Exit a parse tree produced by SyslParser#inplace_table_def.
SyslParserListener.prototype.exitInplace_table_def = function(ctx) {
};


// Enter a parse tree produced by SyslParser#table_refs.
SyslParserListener.prototype.enterTable_refs = function(ctx) {
};

// Exit a parse tree produced by SyslParser#table_refs.
SyslParserListener.prototype.exitTable_refs = function(ctx) {
};


// Enter a parse tree produced by SyslParser#facade.
SyslParserListener.prototype.enterFacade = function(ctx) {
};

// Exit a parse tree produced by SyslParser#facade.
SyslParserListener.prototype.exitFacade = function(ctx) {
};


// Enter a parse tree produced by SyslParser#documentation_stmts.
SyslParserListener.prototype.enterDocumentation_stmts = function(ctx) {
};

// Exit a parse tree produced by SyslParser#documentation_stmts.
SyslParserListener.prototype.exitDocumentation_stmts = function(ctx) {
};


// Enter a parse tree produced by SyslParser#var_in_curly.
SyslParserListener.prototype.enterVar_in_curly = function(ctx) {
};

// Exit a parse tree produced by SyslParser#var_in_curly.
SyslParserListener.prototype.exitVar_in_curly = function(ctx) {
};


// Enter a parse tree produced by SyslParser#query_var.
SyslParserListener.prototype.enterQuery_var = function(ctx) {
};

// Exit a parse tree produced by SyslParser#query_var.
SyslParserListener.prototype.exitQuery_var = function(ctx) {
};


// Enter a parse tree produced by SyslParser#query_param.
SyslParserListener.prototype.enterQuery_param = function(ctx) {
};

// Exit a parse tree produced by SyslParser#query_param.
SyslParserListener.prototype.exitQuery_param = function(ctx) {
};


// Enter a parse tree produced by SyslParser#http_path_part.
SyslParserListener.prototype.enterHttp_path_part = function(ctx) {
};

// Exit a parse tree produced by SyslParser#http_path_part.
SyslParserListener.prototype.exitHttp_path_part = function(ctx) {
};


// Enter a parse tree produced by SyslParser#http_path_var_with_type.
SyslParserListener.prototype.enterHttp_path_var_with_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#http_path_var_with_type.
SyslParserListener.prototype.exitHttp_path_var_with_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#http_path_static.
SyslParserListener.prototype.enterHttp_path_static = function(ctx) {
};

// Exit a parse tree produced by SyslParser#http_path_static.
SyslParserListener.prototype.exitHttp_path_static = function(ctx) {
};


// Enter a parse tree produced by SyslParser#http_path_suffix.
SyslParserListener.prototype.enterHttp_path_suffix = function(ctx) {
};

// Exit a parse tree produced by SyslParser#http_path_suffix.
SyslParserListener.prototype.exitHttp_path_suffix = function(ctx) {
};


// Enter a parse tree produced by SyslParser#http_path.
SyslParserListener.prototype.enterHttp_path = function(ctx) {
};

// Exit a parse tree produced by SyslParser#http_path.
SyslParserListener.prototype.exitHttp_path = function(ctx) {
};


// Enter a parse tree produced by SyslParser#endpoint_name.
SyslParserListener.prototype.enterEndpoint_name = function(ctx) {
};

// Exit a parse tree produced by SyslParser#endpoint_name.
SyslParserListener.prototype.exitEndpoint_name = function(ctx) {
};


// Enter a parse tree produced by SyslParser#ret_stmt.
SyslParserListener.prototype.enterRet_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#ret_stmt.
SyslParserListener.prototype.exitRet_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#target.
SyslParserListener.prototype.enterTarget = function(ctx) {
};

// Exit a parse tree produced by SyslParser#target.
SyslParserListener.prototype.exitTarget = function(ctx) {
};


// Enter a parse tree produced by SyslParser#target_endpoint.
SyslParserListener.prototype.enterTarget_endpoint = function(ctx) {
};

// Exit a parse tree produced by SyslParser#target_endpoint.
SyslParserListener.prototype.exitTarget_endpoint = function(ctx) {
};


// Enter a parse tree produced by SyslParser#call_arg.
SyslParserListener.prototype.enterCall_arg = function(ctx) {
};

// Exit a parse tree produced by SyslParser#call_arg.
SyslParserListener.prototype.exitCall_arg = function(ctx) {
};


// Enter a parse tree produced by SyslParser#call_args.
SyslParserListener.prototype.enterCall_args = function(ctx) {
};

// Exit a parse tree produced by SyslParser#call_args.
SyslParserListener.prototype.exitCall_args = function(ctx) {
};


// Enter a parse tree produced by SyslParser#call_stmt.
SyslParserListener.prototype.enterCall_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#call_stmt.
SyslParserListener.prototype.exitCall_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#if_stmt.
SyslParserListener.prototype.enterIf_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#if_stmt.
SyslParserListener.prototype.exitIf_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#else_stmt.
SyslParserListener.prototype.enterElse_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#else_stmt.
SyslParserListener.prototype.exitElse_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#if_else.
SyslParserListener.prototype.enterIf_else = function(ctx) {
};

// Exit a parse tree produced by SyslParser#if_else.
SyslParserListener.prototype.exitIf_else = function(ctx) {
};


// Enter a parse tree produced by SyslParser#for_stmt.
SyslParserListener.prototype.enterFor_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#for_stmt.
SyslParserListener.prototype.exitFor_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#http_method_comment.
SyslParserListener.prototype.enterHttp_method_comment = function(ctx) {
};

// Exit a parse tree produced by SyslParser#http_method_comment.
SyslParserListener.prototype.exitHttp_method_comment = function(ctx) {
};


// Enter a parse tree produced by SyslParser#group_stmt.
SyslParserListener.prototype.enterGroup_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#group_stmt.
SyslParserListener.prototype.exitGroup_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#one_of_case_label.
SyslParserListener.prototype.enterOne_of_case_label = function(ctx) {
};

// Exit a parse tree produced by SyslParser#one_of_case_label.
SyslParserListener.prototype.exitOne_of_case_label = function(ctx) {
};


// Enter a parse tree produced by SyslParser#one_of_cases.
SyslParserListener.prototype.enterOne_of_cases = function(ctx) {
};

// Exit a parse tree produced by SyslParser#one_of_cases.
SyslParserListener.prototype.exitOne_of_cases = function(ctx) {
};


// Enter a parse tree produced by SyslParser#one_of_stmt.
SyslParserListener.prototype.enterOne_of_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#one_of_stmt.
SyslParserListener.prototype.exitOne_of_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#text_stmt.
SyslParserListener.prototype.enterText_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#text_stmt.
SyslParserListener.prototype.exitText_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#mixin.
SyslParserListener.prototype.enterMixin = function(ctx) {
};

// Exit a parse tree produced by SyslParser#mixin.
SyslParserListener.prototype.exitMixin = function(ctx) {
};


// Enter a parse tree produced by SyslParser#param.
SyslParserListener.prototype.enterParam = function(ctx) {
};

// Exit a parse tree produced by SyslParser#param.
SyslParserListener.prototype.exitParam = function(ctx) {
};


// Enter a parse tree produced by SyslParser#param_list.
SyslParserListener.prototype.enterParam_list = function(ctx) {
};

// Exit a parse tree produced by SyslParser#param_list.
SyslParserListener.prototype.exitParam_list = function(ctx) {
};


// Enter a parse tree produced by SyslParser#params.
SyslParserListener.prototype.enterParams = function(ctx) {
};

// Exit a parse tree produced by SyslParser#params.
SyslParserListener.prototype.exitParams = function(ctx) {
};


// Enter a parse tree produced by SyslParser#statements.
SyslParserListener.prototype.enterStatements = function(ctx) {
};

// Exit a parse tree produced by SyslParser#statements.
SyslParserListener.prototype.exitStatements = function(ctx) {
};


// Enter a parse tree produced by SyslParser#method_def.
SyslParserListener.prototype.enterMethod_def = function(ctx) {
};

// Exit a parse tree produced by SyslParser#method_def.
SyslParserListener.prototype.exitMethod_def = function(ctx) {
};


// Enter a parse tree produced by SyslParser#shortcut.
SyslParserListener.prototype.enterShortcut = function(ctx) {
};

// Exit a parse tree produced by SyslParser#shortcut.
SyslParserListener.prototype.exitShortcut = function(ctx) {
};


// Enter a parse tree produced by SyslParser#simple_endpoint.
SyslParserListener.prototype.enterSimple_endpoint = function(ctx) {
};

// Exit a parse tree produced by SyslParser#simple_endpoint.
SyslParserListener.prototype.exitSimple_endpoint = function(ctx) {
};


// Enter a parse tree produced by SyslParser#rest_endpoint.
SyslParserListener.prototype.enterRest_endpoint = function(ctx) {
};

// Exit a parse tree produced by SyslParser#rest_endpoint.
SyslParserListener.prototype.exitRest_endpoint = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_query_var.
SyslParserListener.prototype.enterCollector_query_var = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_query_var.
SyslParserListener.prototype.exitCollector_query_var = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_query_param.
SyslParserListener.prototype.enterCollector_query_param = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_query_param.
SyslParserListener.prototype.exitCollector_query_param = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_call_stmt.
SyslParserListener.prototype.enterCollector_call_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_call_stmt.
SyslParserListener.prototype.exitCollector_call_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_http_stmt_part.
SyslParserListener.prototype.enterCollector_http_stmt_part = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_http_stmt_part.
SyslParserListener.prototype.exitCollector_http_stmt_part = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_http_stmt_suffix.
SyslParserListener.prototype.enterCollector_http_stmt_suffix = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_http_stmt_suffix.
SyslParserListener.prototype.exitCollector_http_stmt_suffix = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_http_stmt.
SyslParserListener.prototype.enterCollector_http_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_http_stmt.
SyslParserListener.prototype.exitCollector_http_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#publisher.
SyslParserListener.prototype.enterPublisher = function(ctx) {
};

// Exit a parse tree produced by SyslParser#publisher.
SyslParserListener.prototype.exitPublisher = function(ctx) {
};


// Enter a parse tree produced by SyslParser#subscriber.
SyslParserListener.prototype.enterSubscriber = function(ctx) {
};

// Exit a parse tree produced by SyslParser#subscriber.
SyslParserListener.prototype.exitSubscriber = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_pubsub_call.
SyslParserListener.prototype.enterCollector_pubsub_call = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_pubsub_call.
SyslParserListener.prototype.exitCollector_pubsub_call = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_action_stmt.
SyslParserListener.prototype.enterCollector_action_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_action_stmt.
SyslParserListener.prototype.exitCollector_action_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector_stmts.
SyslParserListener.prototype.enterCollector_stmts = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector_stmts.
SyslParserListener.prototype.exitCollector_stmts = function(ctx) {
};


// Enter a parse tree produced by SyslParser#collector.
SyslParserListener.prototype.enterCollector = function(ctx) {
};

// Exit a parse tree produced by SyslParser#collector.
SyslParserListener.prototype.exitCollector = function(ctx) {
};


// Enter a parse tree produced by SyslParser#event.
SyslParserListener.prototype.enterEvent = function(ctx) {
};

// Exit a parse tree produced by SyslParser#event.
SyslParserListener.prototype.exitEvent = function(ctx) {
};


// Enter a parse tree produced by SyslParser#subscribe.
SyslParserListener.prototype.enterSubscribe = function(ctx) {
};

// Exit a parse tree produced by SyslParser#subscribe.
SyslParserListener.prototype.exitSubscribe = function(ctx) {
};


// Enter a parse tree produced by SyslParser#view_type_spec.
SyslParserListener.prototype.enterView_type_spec = function(ctx) {
};

// Exit a parse tree produced by SyslParser#view_type_spec.
SyslParserListener.prototype.exitView_type_spec = function(ctx) {
};


// Enter a parse tree produced by SyslParser#literal.
SyslParserListener.prototype.enterLiteral = function(ctx) {
};

// Exit a parse tree produced by SyslParser#literal.
SyslParserListener.prototype.exitLiteral = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_table_of_op.
SyslParserListener.prototype.enterExpr_table_of_op = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_table_of_op.
SyslParserListener.prototype.exitExpr_table_of_op = function(ctx) {
};


// Enter a parse tree produced by SyslParser#func_arg.
SyslParserListener.prototype.enterFunc_arg = function(ctx) {
};

// Exit a parse tree produced by SyslParser#func_arg.
SyslParserListener.prototype.exitFunc_arg = function(ctx) {
};


// Enter a parse tree produced by SyslParser#func_args.
SyslParserListener.prototype.enterFunc_args = function(ctx) {
};

// Exit a parse tree produced by SyslParser#func_args.
SyslParserListener.prototype.exitFunc_args = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_func.
SyslParserListener.prototype.enterExpr_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_func.
SyslParserListener.prototype.exitExpr_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#rank_expr.
SyslParserListener.prototype.enterRank_expr = function(ctx) {
};

// Exit a parse tree produced by SyslParser#rank_expr.
SyslParserListener.prototype.exitRank_expr = function(ctx) {
};


// Enter a parse tree produced by SyslParser#rank_expr_list.
SyslParserListener.prototype.enterRank_expr_list = function(ctx) {
};

// Exit a parse tree produced by SyslParser#rank_expr_list.
SyslParserListener.prototype.exitRank_expr_list = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_rank_func.
SyslParserListener.prototype.enterExpr_rank_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_rank_func.
SyslParserListener.prototype.exitExpr_rank_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_agg_func.
SyslParserListener.prototype.enterExpr_agg_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_agg_func.
SyslParserListener.prototype.exitExpr_agg_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#e_scope_var.
SyslParserListener.prototype.enterE_scope_var = function(ctx) {
};

// Exit a parse tree produced by SyslParser#e_scope_var.
SyslParserListener.prototype.exitE_scope_var = function(ctx) {
};


// Enter a parse tree produced by SyslParser#first_func_target.
SyslParserListener.prototype.enterFirst_func_target = function(ctx) {
};

// Exit a parse tree produced by SyslParser#first_func_target.
SyslParserListener.prototype.exitFirst_func_target = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_first_func.
SyslParserListener.prototype.enterExpr_first_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_first_func.
SyslParserListener.prototype.exitExpr_first_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#e_single_arg_func.
SyslParserListener.prototype.enterE_single_arg_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#e_single_arg_func.
SyslParserListener.prototype.exitE_single_arg_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_single_arg_func.
SyslParserListener.prototype.enterExpr_single_arg_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_single_arg_func.
SyslParserListener.prototype.exitExpr_single_arg_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_any_func.
SyslParserListener.prototype.enterExpr_any_func = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_any_func.
SyslParserListener.prototype.exitExpr_any_func = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_single_or_null.
SyslParserListener.prototype.enterExpr_single_or_null = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_single_or_null.
SyslParserListener.prototype.exitExpr_single_or_null = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_snapshot.
SyslParserListener.prototype.enterExpr_snapshot = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_snapshot.
SyslParserListener.prototype.exitExpr_snapshot = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_count.
SyslParserListener.prototype.enterExpr_count = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_count.
SyslParserListener.prototype.exitExpr_count = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_navigate_attr.
SyslParserListener.prototype.enterExpr_navigate_attr = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_navigate_attr.
SyslParserListener.prototype.exitExpr_navigate_attr = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_navigate.
SyslParserListener.prototype.enterExpr_navigate = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_navigate.
SyslParserListener.prototype.exitExpr_navigate = function(ctx) {
};


// Enter a parse tree produced by SyslParser#matching_rhs.
SyslParserListener.prototype.enterMatching_rhs = function(ctx) {
};

// Exit a parse tree produced by SyslParser#matching_rhs.
SyslParserListener.prototype.exitMatching_rhs = function(ctx) {
};


// Enter a parse tree produced by SyslParser#squiggly_args.
SyslParserListener.prototype.enterSquiggly_args = function(ctx) {
};

// Exit a parse tree produced by SyslParser#squiggly_args.
SyslParserListener.prototype.exitSquiggly_args = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_matching.
SyslParserListener.prototype.enterExpr_matching = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_matching.
SyslParserListener.prototype.exitExpr_matching = function(ctx) {
};


// Enter a parse tree produced by SyslParser#relop.
SyslParserListener.prototype.enterRelop = function(ctx) {
};

// Exit a parse tree produced by SyslParser#relop.
SyslParserListener.prototype.exitRelop = function(ctx) {
};


// Enter a parse tree produced by SyslParser#list_item.
SyslParserListener.prototype.enterList_item = function(ctx) {
};

// Exit a parse tree produced by SyslParser#list_item.
SyslParserListener.prototype.exitList_item = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_list.
SyslParserListener.prototype.enterExpr_list = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_list.
SyslParserListener.prototype.exitExpr_list = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_set.
SyslParserListener.prototype.enterExpr_set = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_set.
SyslParserListener.prototype.exitExpr_set = function(ctx) {
};


// Enter a parse tree produced by SyslParser#empty_tuple.
SyslParserListener.prototype.enterEmpty_tuple = function(ctx) {
};

// Exit a parse tree produced by SyslParser#empty_tuple.
SyslParserListener.prototype.exitEmpty_tuple = function(ctx) {
};


// Enter a parse tree produced by SyslParser#atom_dot_relop.
SyslParserListener.prototype.enterAtom_dot_relop = function(ctx) {
};

// Exit a parse tree produced by SyslParser#atom_dot_relop.
SyslParserListener.prototype.exitAtom_dot_relop = function(ctx) {
};


// Enter a parse tree produced by SyslParser#atomT_implied_dot.
SyslParserListener.prototype.enterAtomT_implied_dot = function(ctx) {
};

// Exit a parse tree produced by SyslParser#atomT_implied_dot.
SyslParserListener.prototype.exitAtomT_implied_dot = function(ctx) {
};


// Enter a parse tree produced by SyslParser#atomT_name.
SyslParserListener.prototype.enterAtomT_name = function(ctx) {
};

// Exit a parse tree produced by SyslParser#atomT_name.
SyslParserListener.prototype.exitAtomT_name = function(ctx) {
};


// Enter a parse tree produced by SyslParser#atomT_paren.
SyslParserListener.prototype.enterAtomT_paren = function(ctx) {
};

// Exit a parse tree produced by SyslParser#atomT_paren.
SyslParserListener.prototype.exitAtomT_paren = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_atom_list.
SyslParserListener.prototype.enterExpr_atom_list = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_atom_list.
SyslParserListener.prototype.exitExpr_atom_list = function(ctx) {
};


// Enter a parse tree produced by SyslParser#atomT.
SyslParserListener.prototype.enterAtomT = function(ctx) {
};

// Exit a parse tree produced by SyslParser#atomT.
SyslParserListener.prototype.exitAtomT = function(ctx) {
};


// Enter a parse tree produced by SyslParser#atom.
SyslParserListener.prototype.enterAtom = function(ctx) {
};

// Exit a parse tree produced by SyslParser#atom.
SyslParserListener.prototype.exitAtom = function(ctx) {
};


// Enter a parse tree produced by SyslParser#powerT.
SyslParserListener.prototype.enterPowerT = function(ctx) {
};

// Exit a parse tree produced by SyslParser#powerT.
SyslParserListener.prototype.exitPowerT = function(ctx) {
};


// Enter a parse tree produced by SyslParser#power.
SyslParserListener.prototype.enterPower = function(ctx) {
};

// Exit a parse tree produced by SyslParser#power.
SyslParserListener.prototype.exitPower = function(ctx) {
};


// Enter a parse tree produced by SyslParser#unaryTerm.
SyslParserListener.prototype.enterUnaryTerm = function(ctx) {
};

// Exit a parse tree produced by SyslParser#unaryTerm.
SyslParserListener.prototype.exitUnaryTerm = function(ctx) {
};


// Enter a parse tree produced by SyslParser#termT.
SyslParserListener.prototype.enterTermT = function(ctx) {
};

// Exit a parse tree produced by SyslParser#termT.
SyslParserListener.prototype.exitTermT = function(ctx) {
};


// Enter a parse tree produced by SyslParser#term.
SyslParserListener.prototype.enterTerm = function(ctx) {
};

// Exit a parse tree produced by SyslParser#term.
SyslParserListener.prototype.exitTerm = function(ctx) {
};


// Enter a parse tree produced by SyslParser#binexprT.
SyslParserListener.prototype.enterBinexprT = function(ctx) {
};

// Exit a parse tree produced by SyslParser#binexprT.
SyslParserListener.prototype.exitBinexprT = function(ctx) {
};


// Enter a parse tree produced by SyslParser#binexpr.
SyslParserListener.prototype.enterBinexpr = function(ctx) {
};

// Exit a parse tree produced by SyslParser#binexpr.
SyslParserListener.prototype.exitBinexpr = function(ctx) {
};


// Enter a parse tree produced by SyslParser#e_compare_ops.
SyslParserListener.prototype.enterE_compare_ops = function(ctx) {
};

// Exit a parse tree produced by SyslParser#e_compare_ops.
SyslParserListener.prototype.exitE_compare_ops = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_rel.
SyslParserListener.prototype.enterExpr_rel = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_rel.
SyslParserListener.prototype.exitExpr_rel = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_bitand.
SyslParserListener.prototype.enterExpr_bitand = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_bitand.
SyslParserListener.prototype.exitExpr_bitand = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_bitxor.
SyslParserListener.prototype.enterExpr_bitxor = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_bitxor.
SyslParserListener.prototype.exitExpr_bitxor = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_bitor.
SyslParserListener.prototype.enterExpr_bitor = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_bitor.
SyslParserListener.prototype.exitExpr_bitor = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_and.
SyslParserListener.prototype.enterExpr_and = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_and.
SyslParserListener.prototype.exitExpr_and = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_or.
SyslParserListener.prototype.enterExpr_or = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_or.
SyslParserListener.prototype.exitExpr_or = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_but_not.
SyslParserListener.prototype.enterExpr_but_not = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_but_not.
SyslParserListener.prototype.exitExpr_but_not = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_coalesce.
SyslParserListener.prototype.enterExpr_coalesce = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_coalesce.
SyslParserListener.prototype.exitExpr_coalesce = function(ctx) {
};


// Enter a parse tree produced by SyslParser#if_one_liner.
SyslParserListener.prototype.enterIf_one_liner = function(ctx) {
};

// Exit a parse tree produced by SyslParser#if_one_liner.
SyslParserListener.prototype.exitIf_one_liner = function(ctx) {
};


// Enter a parse tree produced by SyslParser#else_block_stmt.
SyslParserListener.prototype.enterElse_block_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#else_block_stmt.
SyslParserListener.prototype.exitElse_block_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#control_item.
SyslParserListener.prototype.enterControl_item = function(ctx) {
};

// Exit a parse tree produced by SyslParser#control_item.
SyslParserListener.prototype.exitControl_item = function(ctx) {
};


// Enter a parse tree produced by SyslParser#if_controls.
SyslParserListener.prototype.enterIf_controls = function(ctx) {
};

// Exit a parse tree produced by SyslParser#if_controls.
SyslParserListener.prototype.exitIf_controls = function(ctx) {
};


// Enter a parse tree produced by SyslParser#cond_block.
SyslParserListener.prototype.enterCond_block = function(ctx) {
};

// Exit a parse tree produced by SyslParser#cond_block.
SyslParserListener.prototype.exitCond_block = function(ctx) {
};


// Enter a parse tree produced by SyslParser#final_else.
SyslParserListener.prototype.enterFinal_else = function(ctx) {
};

// Exit a parse tree produced by SyslParser#final_else.
SyslParserListener.prototype.exitFinal_else = function(ctx) {
};


// Enter a parse tree produced by SyslParser#ifvar.
SyslParserListener.prototype.enterIfvar = function(ctx) {
};

// Exit a parse tree produced by SyslParser#ifvar.
SyslParserListener.prototype.exitIfvar = function(ctx) {
};


// Enter a parse tree produced by SyslParser#if_multiple_lines.
SyslParserListener.prototype.enterIf_multiple_lines = function(ctx) {
};

// Exit a parse tree produced by SyslParser#if_multiple_lines.
SyslParserListener.prototype.exitIf_multiple_lines = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_if_else.
SyslParserListener.prototype.enterExpr_if_else = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_if_else.
SyslParserListener.prototype.exitExpr_if_else = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr.
SyslParserListener.prototype.enterExpr = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr.
SyslParserListener.prototype.exitExpr = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_assign.
SyslParserListener.prototype.enterExpr_assign = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_assign.
SyslParserListener.prototype.exitExpr_assign = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_simple_assign.
SyslParserListener.prototype.enterExpr_simple_assign = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_simple_assign.
SyslParserListener.prototype.exitExpr_simple_assign = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_let_statement.
SyslParserListener.prototype.enterExpr_let_statement = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_let_statement.
SyslParserListener.prototype.exitExpr_let_statement = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_table_of_statement.
SyslParserListener.prototype.enterExpr_table_of_statement = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_table_of_statement.
SyslParserListener.prototype.exitExpr_table_of_statement = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_dot_assign.
SyslParserListener.prototype.enterExpr_dot_assign = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_dot_assign.
SyslParserListener.prototype.exitExpr_dot_assign = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_statement_no_nl.
SyslParserListener.prototype.enterExpr_statement_no_nl = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_statement_no_nl.
SyslParserListener.prototype.exitExpr_statement_no_nl = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_statement.
SyslParserListener.prototype.enterExpr_statement = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_statement.
SyslParserListener.prototype.exitExpr_statement = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_inject_stmt.
SyslParserListener.prototype.enterExpr_inject_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_inject_stmt.
SyslParserListener.prototype.exitExpr_inject_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_stmt.
SyslParserListener.prototype.enterExpr_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_stmt.
SyslParserListener.prototype.exitExpr_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#transform_return_type.
SyslParserListener.prototype.enterTransform_return_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#transform_return_type.
SyslParserListener.prototype.exitTransform_return_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#view_return_type.
SyslParserListener.prototype.enterView_return_type = function(ctx) {
};

// Exit a parse tree produced by SyslParser#view_return_type.
SyslParserListener.prototype.exitView_return_type = function(ctx) {
};


// Enter a parse tree produced by SyslParser#transform_scope_var.
SyslParserListener.prototype.enterTransform_scope_var = function(ctx) {
};

// Exit a parse tree produced by SyslParser#transform_scope_var.
SyslParserListener.prototype.exitTransform_scope_var = function(ctx) {
};


// Enter a parse tree produced by SyslParser#transform_arg.
SyslParserListener.prototype.enterTransform_arg = function(ctx) {
};

// Exit a parse tree produced by SyslParser#transform_arg.
SyslParserListener.prototype.exitTransform_arg = function(ctx) {
};


// Enter a parse tree produced by SyslParser#transform.
SyslParserListener.prototype.enterTransform = function(ctx) {
};

// Exit a parse tree produced by SyslParser#transform.
SyslParserListener.prototype.exitTransform = function(ctx) {
};


// Enter a parse tree produced by SyslParser#expr_block.
SyslParserListener.prototype.enterExpr_block = function(ctx) {
};

// Exit a parse tree produced by SyslParser#expr_block.
SyslParserListener.prototype.exitExpr_block = function(ctx) {
};


// Enter a parse tree produced by SyslParser#view_param.
SyslParserListener.prototype.enterView_param = function(ctx) {
};

// Exit a parse tree produced by SyslParser#view_param.
SyslParserListener.prototype.exitView_param = function(ctx) {
};


// Enter a parse tree produced by SyslParser#view_params.
SyslParserListener.prototype.enterView_params = function(ctx) {
};

// Exit a parse tree produced by SyslParser#view_params.
SyslParserListener.prototype.exitView_params = function(ctx) {
};


// Enter a parse tree produced by SyslParser#abstract_view.
SyslParserListener.prototype.enterAbstract_view = function(ctx) {
};

// Exit a parse tree produced by SyslParser#abstract_view.
SyslParserListener.prototype.exitAbstract_view = function(ctx) {
};


// Enter a parse tree produced by SyslParser#view.
SyslParserListener.prototype.enterView = function(ctx) {
};

// Exit a parse tree produced by SyslParser#view.
SyslParserListener.prototype.exitView = function(ctx) {
};


// Enter a parse tree produced by SyslParser#alias.
SyslParserListener.prototype.enterAlias = function(ctx) {
};

// Exit a parse tree produced by SyslParser#alias.
SyslParserListener.prototype.exitAlias = function(ctx) {
};


// Enter a parse tree produced by SyslParser#app_decl.
SyslParserListener.prototype.enterApp_decl = function(ctx) {
};

// Exit a parse tree produced by SyslParser#app_decl.
SyslParserListener.prototype.exitApp_decl = function(ctx) {
};


// Enter a parse tree produced by SyslParser#application.
SyslParserListener.prototype.enterApplication = function(ctx) {
};

// Exit a parse tree produced by SyslParser#application.
SyslParserListener.prototype.exitApplication = function(ctx) {
};


// Enter a parse tree produced by SyslParser#path.
SyslParserListener.prototype.enterPath = function(ctx) {
};

// Exit a parse tree produced by SyslParser#path.
SyslParserListener.prototype.exitPath = function(ctx) {
};


// Enter a parse tree produced by SyslParser#import_stmt.
SyslParserListener.prototype.enterImport_stmt = function(ctx) {
};

// Exit a parse tree produced by SyslParser#import_stmt.
SyslParserListener.prototype.exitImport_stmt = function(ctx) {
};


// Enter a parse tree produced by SyslParser#imports_decl.
SyslParserListener.prototype.enterImports_decl = function(ctx) {
};

// Exit a parse tree produced by SyslParser#imports_decl.
SyslParserListener.prototype.exitImports_decl = function(ctx) {
};


// Enter a parse tree produced by SyslParser#sysl_file.
SyslParserListener.prototype.enterSysl_file = function(ctx) {
};

// Exit a parse tree produced by SyslParser#sysl_file.
SyslParserListener.prototype.exitSysl_file = function(ctx) {
};



exports.SyslParserListener = SyslParserListener;