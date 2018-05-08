parser grammar SyslParser;

options { tokenVocab=SyslLexer; }

modifier        : TILDE Name (PLUS Name)*;
size_spec       : OPEN_PAREN DIGITS ( DOT DIGITS)? CLOSE_PAREN;
modifier_list   : modifier (COMMA modifier)*;
sq_open: SQ_OPEN | SQ_OPEN_2;
modifiers       : sq_open modifier_list SQ_CLOSE;
name_str        : Name|TEXT_LINE;
reference       : app_name (DOT name_str)+;
doc_string      : PIPE TEXT;
quoted_string       : QSTRING;
array_of_strings    : sq_open quoted_string (COMMA quoted_string)* SQ_CLOSE;
array_of_arrays     : sq_open array_of_strings SQ_CLOSE;
nvp                 : Name EQ (quoted_string | array_of_strings| array_of_arrays);
attributes          : sq_open nvp (COMMA nvp)* SQ_CLOSE;
entry               : nvp | modifier ;
attribs_or_modifiers: sq_open entry (COMMA entry)* SQ_CLOSE;
set_type            : SET_OF (Name | reference | NativeDataTypes) size_spec?;
//TODO : allow for other collection types?
collection_type     : set_type;
user_defined_type       : name_str;
multi_line_docstring    :   COLON INDENT doc_string+ DEDENT;
annotation_value        :   QSTRING | array_of_strings | multi_line_docstring;
annotation      : AT VAR_NAME EQ annotation_value;
annotations     : INDENT annotation+ DEDENT;

field_type      : collection_type
                |  ((reference | NativeDataTypes | user_defined_type) (array_size | size_spec)? QN? attribs_or_modifiers? (COLON annotations)?) ;

array_size  :  OPEN_PAREN DIGITS DOTDOT DIGITS? CLOSE_PAREN;
inplace_tuple: INDENT field+ DEDENT;
field: name_str (array_size? LESS_COLON (field_type | inplace_tuple))?;

inplace_table : table;
table   :   SYSL_COMMENT*
            (TABLE | TYPE)
            name_str attribs_or_modifiers? COLON ( WHATEVER | INDENT (SYSL_COMMENT | field | annotation | inplace_table | WHATEVER )+ DEDENT)
        ;

package_name   : name_str;
sub_package    : NAME_SEP package_name;
app_name       : package_name sub_package*;

name_with_attribs       :   app_name  QSTRING? attribs_or_modifiers?;

model_name          :  Name COLON ;
inplace_table_def   :  COLON INDENT (Name attribs_or_modifiers?)+ DEDENT;
table_refs          :  (TABLE | TYPE) Name inplace_table_def?;
facade              :  SYSL_COMMENT* WRAP model_name INDENT table_refs+ DEDENT;

documentation_stmts     : AT Name EQ QSTRING NEWLINE;

var_in_curly    : CURLY_OPEN Name CURLY_CLOSE;
query_var       : Name EQ (NativeDataTypes | var_in_curly) QN?;
query_param     : QN query_var (AMP query_var)*;

http_path_part :name_str;
http_path_var_with_type : CURLY_OPEN http_path_part LESS_COLON (NativeDataTypes | name_str) CURLY_CLOSE;
http_path_static : http_path_part;
http_path_suffix : FORWARD_SLASH (http_path_static | http_path_var_with_type);
http_path       : FORWARD_SLASH | http_path_suffix+;

endpoint_name   : ( Name+ | TEXT_LINE);

ret_stmt        : RETURN TEXT;

target          : app_name;
target_endpoint : TEXT_VALUE;
call_arg : (Q_ARG | TEXT_VALUE)+ | (TEXT_VALUE LESS_COLON_2 TEXT_VALUE);
call_args: OPEN_PAREN_ARG call_arg (COMMA_ARG call_arg)* CLOSE_PAREN_ARG;
call_stmt       : (DOT_ARROW | target ARROW_LEFT) target_endpoint call_args?;

arg_value: TEXT_VALUE | Q_ARG;
if_stmt                 : IF arg_value COLON_ARG INDENT statements* DEDENT;
if_else                 : if_stmt (ELSE name_str? COLON INDENT statements* DEDENT)?;

for_cond                : arg_value;

for_stmt                : (UNTIL | FOR_EACH | FOR) for_cond COLON_ARG
                                INDENT statements* DEDENT;

http_method_comment     : SYSL_COMMENT;
//group_stmt              : GROUP TEXT_NAME COLON
//                                INDENT statements+ DEDENT;

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
                // | group_stmt
                | text_stmt
                | annotation
            )
            params?
            attribs_or_modifiers?
            ;

method_def: HTTP_VERBS params?
                  query_param? attribs_or_modifiers? COLON
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
collector_call_stmt:  app_name (ARROW_LEFT TEXT_VALUE)?;
collector_http_stmt_part: name_str | CURLY_OPEN name_str CURLY_CLOSE ;
collector_http_stmt: HTTP_VERBS (FORWARD_SLASH collector_http_stmt_part)+ collector_query_param?;
collector_stmts: (collector_call_stmt | collector_http_stmt) attribs_or_modifiers;

collector:  COLLECTOR COLON (WHATEVER | (INDENT collector_stmts+ DEDENT));

event: DISTANCE name_str params?
        attribs_or_modifiers? COLON (WHATEVER | INDENT statements+ DEDENT);

subscribe: app_name ARROW_RIGHT name_str attribs_or_modifiers? COLON (WHATEVER | INDENT statements+ DEDENT);

app_decl: INDENT  (table | facade | SYSL_COMMENT | rest_endpoint | simple_endpoint | collector | event | subscribe | annotation | mixin )+ DEDENT;

application:  SYSL_COMMENT*
                name_with_attribs
                COLON
                app_decl
                ;

path            : FORWARD_SLASH? Name (FORWARD_SLASH Name)* ;
import_stmt     : IMPORT SYSL_COMMENT*;
imports_decl    : import_stmt+;
sysl_file       : imports_decl? application+ EOF;
