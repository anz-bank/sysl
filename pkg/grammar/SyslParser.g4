parser grammar SyslParser;

options { tokenVocab=SyslLexer; }

modifier        : TILDE Name (PLUS Name)*;
size_spec       : OPEN_PAREN DIGITS ( DOT DIGITS)? CLOSE_PAREN;
name_str        : Name|TEXT_LINE|E_Name;

reference       : app_name ((E_DOT | DOT) name_str)+;
doc_string      : PIPE TEXT;
quoted_string       : QSTRING;
array_of_strings    : SQ_OPEN quoted_string (COMMA quoted_string)* SQ_CLOSE;
array_of_arrays     : SQ_OPEN array_of_strings (COMMA array_of_strings)* SQ_CLOSE;
nvp                 : Name EQ (quoted_string | array_of_strings| array_of_arrays);
entry               : nvp | modifier ;
attribs_or_modifiers: SQ_OPEN entry (COMMA entry)* SQ_CLOSE;

user_defined_type       : name_str;
types: user_defined_type | reference | NativeDataTypes;
set_of: SET_OF | (E_SET_OF);
set_type            : set_of types size_spec?;
sequence_of: SEQUENCE_OF | (E_SEQUENCE_OF);
sequence_type            : sequence_of types size_spec?;
//TODO : allow for other collection types?
collection_type     : set_type | sequence_type;
multi_line_docstring    :   COLON INDENT doc_string+ DEDENT;
annotation_value        :   QSTRING | array_of_strings | multi_line_docstring;
annotation      : AT VAR_NAME EQ annotation_value;
annotations     : INDENT annotation+ DEDENT;

field_type      : (collection_type |  (types (array_size | size_spec)?))
                    QN? attribs_or_modifiers? (COLON annotations)?;

array_size  :  OPEN_PAREN DIGITS DOTDOT DIGITS? CLOSE_PAREN;
inplace_tuple: INDENT field+ DEDENT;
field: name_str (array_size? LESS_COLON (field_type | inplace_tuple) QSTRING? )?;

inplace_table : table;

table_stmts:  INDENT
                (SYSL_COMMENT | field | annotation | inplace_table | WHATEVER)+
              DEDENT;

table_def: attribs_or_modifiers? COLON ( WHATEVER | table_stmts);

table   :  SYSL_COMMENT* (TYPE | TABLE) name_str table_def;

union   :   SYSL_COMMENT*
            UNION
            name_str attribs_or_modifiers? COLON ( WHATEVER | INDENT (SYSL_COMMENT | user_defined_type | annotation | WHATEVER )+ DEDENT)
        ;

package_name   : name_str;
sub_package    : NAME_SEP package_name;
app_name       : package_name sub_package*;

name_with_attribs       :   app_name  QSTRING? attribs_or_modifiers?;

model_name          :  Name COLON ;
inplace_table_def   :  COLON INDENT (Name attribs_or_modifiers?)+ DEDENT;
table_refs          :  (TABLE | TYPE | UNION) Name inplace_table_def?;
facade              :  SYSL_COMMENT* WRAP model_name INDENT table_refs+ DEDENT;

var_in_curly    : CURLY_OPEN Name CURLY_CLOSE;
query_var       : Name EQ (NativeDataTypes | name_str | var_in_curly) QN?;
query_param     : QN query_var (AMP query_var)*;

http_path_part :name_str | DIGITS;
http_path_var_with_type : CURLY_OPEN http_path_part LESS_COLON (NativeDataTypes | name_str | reference) CURLY_CLOSE;
http_path_static : http_path_part;
http_path_suffix : FORWARD_SLASH (http_path_static | http_path_var_with_type);
http_path       : (FORWARD_SLASH | http_path_suffix+);

endpoint_name   : name_str (FORWARD_SLASH name_str)*;

ret_stmt        : RETURN TEXT;

target          : app_name;
target_endpoint : name_str;
call_arg : (QSTRING | name_str)+ | (name_str LESS_COLON (name_str|NativeDataTypes));
call_args: OPEN_PAREN call_arg (COMMA call_arg)* CLOSE_PAREN;
call_stmt       : (DOT_ARROW | target ARROW_LEFT) target_endpoint call_args?;

if_stmt                 : IF PREDICATE_VALUE COLON INDENT statements* DEDENT;
else_stmt               : ELSE PREDICATE_VALUE? COLON INDENT statements* DEDENT;
if_else                 : if_stmt else_stmt*;

for_stmt                : (ALT | UNTIL | FOR_EACH | FOR | LOOP | WHILE ) PREDICATE_VALUE COLON
                                INDENT statements* DEDENT;

http_method_comment     : SYSL_COMMENT;
group_stmt              : name_str COLON
                               INDENT statements+ DEDENT;

one_of_case_label: (Name | TEXT_LINE | QSTRING)+;

one_of_cases: one_of_case_label? COLON
                    INDENT statements+ DEDENT;

one_of_stmt             : ONE_OF COLON
                           INDENT one_of_cases+ DEDENT;

text_stmt               : doc_string | QSTRING | app_name (ARROW_RIGHT name_str)?  | WHATEVER ;

mixin:  MIXIN app_name;

param: reference | field;

param_list: param (COMMA param)*;

params : OPEN_PAREN param_list CLOSE_PAREN;

statements: ( if_else
                | for_stmt
                | ret_stmt
                | call_stmt
                | one_of_stmt
                | http_method_comment
                | group_stmt
                | text_stmt
                | annotation
            )
            attribs_or_modifiers?
            ;

method_def: HTTP_VERBS params? query_param? attribs_or_modifiers? COLON
                        INDENT statements+ DEDENT
                ;

shortcut        : WHATEVER;

simple_endpoint :
                WHATEVER
                | (
                    endpoint_name QSTRING? params? attribs_or_modifiers? COLON
                    (  shortcut
                        | (INDENT statements+ DEDENT)
                    )
                )
                ;


rest_endpoint: http_path attribs_or_modifiers? COLON
                                    (INDENT (method_def | rest_endpoint)+ DEDENT)
                ;

collector_query_var: name_str EQ (NativeDataTypes | name_str);
collector_query_param: QN collector_query_var (AMP collector_query_var)*;
collector_call_stmt:  target ARROW_LEFT target_endpoint;

collector_http_stmt_part: name_str | CURLY_OPEN name_str CURLY_CLOSE ;
collector_http_stmt_suffix: (FORWARD_SLASH collector_http_stmt_part)+ collector_query_param?;
collector_http_stmt: HTTP_VERBS collector_http_stmt_suffix;

publisher: app_name;
subscriber: app_name;
collector_pubsub_call: subscriber ARROW_LEFT publisher ARROW_RIGHT name_str;

collector_action_stmt: name_str;
collector_stmts: (collector_action_stmt | collector_call_stmt | collector_http_stmt | collector_pubsub_call) attribs_or_modifiers;

collector:  COLLECTOR COLON (WHATEVER | (INDENT collector_stmts+ DEDENT));

event: DISTANCE name_str params?
        attribs_or_modifiers? COLON (WHATEVER | INDENT statements+ DEDENT);

subscribe: app_name ARROW_RIGHT name_str attribs_or_modifiers? COLON (WHATEVER | INDENT statements+ DEDENT);

view_type_spec: collection_type | types;

literal: E_DIGITS | E_DECIMAL | E_STRING_DBL | E_STRING_SINGLE | E_NULL | E_TRUE | E_FALSE;

// expr_qualified_name: E_QN? E_RefName;
expr_table_of_op: (E_DOT | E_NULLSAFE_DOT) E_TABLE_OF? E_Name;

func_arg: expr;
func_args: func_arg (E_COMMA func_arg)*;
expr_func: ( E_FUNC | E_Name | NativeDataTypes) E_OPEN_PAREN func_args? E_CLOSE_PAREN;

rank_expr: expr (E_ASC | E_DESC)?;
rank_expr_list: rank_expr (E_COMMA rank_expr)*;

expr_rank_func: E_RELOPS_RANK
                (E_ANGLE_OPEN view_type_spec E_ANGLE_CLOSE)?
                E_OPEN_PAREN rank_expr_list E_AS E_Name E_CLOSE_PAREN;

expr_agg_func: E_RELOPS_AGG
                E_OPEN_PAREN e_scope_var? func_args E_CLOSE_PAREN;

e_scope_var: E_Name E_COLON;
first_func_target: E_NULL | expr;
expr_first_func: E_RELOPS_FIRST first_func_target E_BY E_OPEN_PAREN e_scope_var? rank_expr_list E_CLOSE_PAREN;

e_single_arg_func: E_RELOPS_WHERE | E_RELOPS_FLATTEN;
expr_single_arg_func: e_single_arg_func
                E_OPEN_PAREN e_scope_var? expr E_CLOSE_PAREN;

expr_any_func: E_RELOPS_ANY
                E_OPEN_PAREN expr E_CLOSE_PAREN;

expr_single_or_null: E_RELOPS_SINGLE_NULL;
expr_snapshot: E_RELOPS_SNAPSHOT;
expr_count: E_RELOPS_COUNT;
expr_navigate_attr: (E_DOT? E_Name);
expr_navigate: E_QN? E_ARROW_RIGHT (E_SET_OF | E_SEQUENCE_OF)? expr_navigate_attr (E_VIA E_Name)?;

matching_rhs: expr_table_of_op
    | expr_navigate
    | atomT_paren
    | E_Name;

squiggly_args: E_SQ_OPEN E_Name (E_COMMA E_Name)* E_SQ_CLOSE;
expr_matching: E_NOT? E_TILDE squiggly_args? E_ANGLE_CLOSE matching_rhs;

relop:  expr_table_of_op
        | expr_rank_func
        | expr_agg_func
        | expr_single_arg_func
        | expr_any_func
        | expr_count
        | expr_single_or_null
        | expr_snapshot
        | expr_first_func
        | expr_navigate
        | expr_matching
        ;

list_item: expr;
expr_list: list_item (E_COMMA list_item)*;

expr_set: E_CURLY_OPEN expr_list E_CURLY_CLOSE;

empty_tuple: E_EMPTY_TUPLE;

atom_dot_relop: E_DOT
                (expr_rank_func // 1
                  | expr_agg_func  // 4
                  | expr_single_arg_func // 2
                  | expr_any_func // 1
                  | expr_count // 1
                  | expr_single_or_null // 2
                  | expr_snapshot // 1
                  | expr_first_func // 1
                  | expr_navigate
                  | expr_matching
                )
              ;

atomT_implied_dot: atom_dot_relop
    | expr_navigate
    | expr_matching
    | expr_table_of_op
    ;

atomT_name : E_Name | E_WHATEVER | E_DOT;

atomT_paren: E_OPEN_PAREN expr E_CLOSE_PAREN;
expr_atom_list: E_SQ_OPEN expr_list E_SQ_CLOSE;
atomT: expr_func
    | expr_set
    | expr_atom_list
    | empty_tuple
    | atomT_paren
    | literal
    | atomT_implied_dot
    | atomT_name
    ;

atom: atomT relop*;

powerT: E_POW unaryTerm;
power: atom powerT?;

unaryTerm: (E_PLUS | E_MINUS | E_NOT | E_TILDE )? power;

termT: (E_STAR | E_DIVIDE | E_MOD ) unaryTerm;
term: unaryTerm termT*;

binexprT: (E_PLUS | E_MINUS) term;
binexpr: term binexprT*;
e_compare_ops: E_REL | E_ANGLE_OPEN | E_ANGLE_CLOSE | E_DOUBLE_EQ;
expr_rel: binexpr ( e_compare_ops binexpr)*;
expr_bitand: expr_rel ((E_AMP | E_AND) expr_rel)*;
expr_bitxor: expr_bitand (E_XOR expr_bitand)*;
expr_bitor: expr_bitxor (E_BITOR expr_bitxor)*;
expr_and: expr_bitor (E_DOUBLE_AMP expr_bitor)*;
expr_or: expr_and (E_LOGIC_OR expr_and)*;
expr_but_not: expr_or (E_BUTNOT expr_or)*;
expr_coalesce: expr_but_not (E_COALESCE expr_but_not)*;

if_one_liner: expr E_QN? E_THEN expr E_ELSE expr;

else_block_stmt returns [bool nested]:
      expr
      { $nested=$expr.nested;}
      ;

control_item: expr;
if_controls: control_item (E_COMMA control_item)*;
cond_block: if_controls E_EQ_RIGHT
              else_block_stmt
              ( {$else_block_stmt.nested == true}?| E_NL)
              ;

final_else: E_ELSE
        expr
        ( {$expr.nested == true}?| E_NL)
        ;
ifvar: expr E_DOUBLE_EQ;
if_multiple_lines: ifvar? E_COLON E_NL
                    INDENT cond_block+ final_else? DEDENT;

expr_if_else returns [bool nested] :
        E_IF (
            if_one_liner
            | if_multiple_lines { $nested=true;}
            )
            ;

//
//  EXPR
//
expr returns [bool nested]:
      expr_if_else { $nested=$expr_if_else.nested;}
      | expr_coalesce
      ;

expr_assign
    returns [bool nested]
            : E_EQ (expr { $nested=$expr.nested;} | transform { $nested=true;});
expr_simple_assign returns [bool nested] : E_Name expr_assign {$nested = $expr_assign.nested;};
expr_let_statement returns [bool nested] : E_LET E_Name expr_assign {$nested = $expr_assign.nested;};
expr_table_of_statement returns [bool nested]: E_TABLE_OF E_Name expr_assign {$nested = $expr_assign.nested;};
expr_dot_assign: E_DOT_NAME_NL;

expr_statement_no_nl: expr_dot_assign;

expr_statement locals [bool nested]:
                  (expr_let_statement {$nested = $expr_let_statement.nested;}
                | expr_table_of_statement {$nested = $expr_table_of_statement.nested;}
                // NL is not required when we got nested expression (like if_multiple_lines or transform)
                // followed by new statement
                | expr_simple_assign {$nested = $expr_simple_assign.nested;}
                )
                ( {$nested == true}? | E_NL)
                ;

expr_inject_stmt: expr_func E_DOT E_STAR E_NL;
expr_stmt: expr_statement_no_nl | expr_statement | expr_inject_stmt;

transform_return_type:  set_of | sequence_of | view_type_spec;
view_return_type: view_type_spec;

transform_scope_var: E_Name;
transform_arg: expr;

transform: transform_arg?
    E_ARROW_RIGHT
    (E_ANGLE_OPEN transform_return_type E_ANGLE_CLOSE)?
    E_OPEN_PAREN transform_scope_var? E_COLON E_NL
      INDENT
        (expr_stmt)+
      DEDENT E_CLOSE_PAREN
      E_NL;

expr_block: INDENT transform DEDENT;

view_param: name_str LESS_COLON view_type_spec;

view_params: view_param (COMMA view_param)*;

abstract_view: ABSTRACT;

view
  returns [bool abstractView]:
  VIEW name_str
    OPEN_PAREN view_params CLOSE_PAREN
    (ARROW_RIGHT view_return_type)?
    ( attribs_or_modifiers? COLON expr_block | abstract_view {$abstractView=true;} );

alias: ALIAS name_str attribs_or_modifiers? COLON (
        (annotation* (types | collection_type))
        | (INDENT annotation* (types | collection_type) DEDENT)
        );

app_decl
  locals [bool check]:
    INDENT  (
        alias
        | annotation
        | collector
        | event
        | facade
        | mixin
        | rest_endpoint
        | simple_endpoint
        | subscribe
        | SYSL_COMMENT
        | union
        | view { $check = $view.abstractView;}
        | table
      )+
      ( {$check}? | DEDENT );

application:  SYSL_COMMENT*
                name_with_attribs
                COLON
                app_decl
                ;

import_mode     : TILDE Name;
import_stmt     : IMPORT IMPORT_PATH (AS Name (DOT Name)*)? WS* import_mode? (SYSL_COMMENT*|NEWLINE);
imports_decl    : import_stmt+;

sysl_file       : imports_decl? application+ EOF;
