parser grammar SyslParser;

options { tokenVocab=SyslLexer; }

modifier        : TILDE Name (PLUS Name)*;
size_spec       : OPEN_PAREN DIGITS ( DOT DIGITS)? CLOSE_PAREN;
modifier_list   : modifier (COMMA modifier)*;
modifiers       : SQ_OPEN modifier_list SQ_CLOSE;
reference       : parent_ref=Name '.' member=Name;
doc_string      : PIPE TEXT;
quoted_string       : QSTRING;
array_of_strings    : SQ_OPEN quoted_string (COMMA quoted_string)* SQ_CLOSE;
array_of_arrays     : SQ_OPEN array_of_strings SQ_CLOSE;
nvp                 : Name EQ (quoted_string | array_of_strings| array_of_arrays);
attributes          : SQ_OPEN nvp (COMMA nvp)* SQ_CLOSE;
entry               : nvp | modifier ;
attribs_or_modifiers: SQ_OPEN entry (COMMA entry)* SQ_CLOSE;
set_type            : SET_OF (Name | NativeDataTypes) size_spec?;
//TODO : allow for other collection types?
collection_type     : set_type;
user_defined_type       : Name;
multi_line_docstring    :   COLON INDENT doc_string+ DEDENT;
annotation_value        :   QSTRING | array_of_strings | multi_line_docstring;
annotation      : AT VAR_NAME EQ annotation_value;
annotations     : INDENT annotation+ DEDENT;

field_type      : collection_type
                |  ((reference | NativeDataTypes | user_defined_type) size_spec? QN? attribs_or_modifiers? (COLON annotations)?) ;

field: Name LESS_COLON field_type;

table   :   SYSL_COMMENT*
            (TABLE | TYPE)
            Name attribs_or_modifiers? COLON
            INDENT (SYSL_COMMENT | field )+ DEDENT
        ;

package_name   : Name (DOT Name)? | TEXT_LINE | TEXT_NAME;
sub_package    : NAME_SEP package_name;
app_name       : package_name sub_package*;

name_with_attribs       :   app_name  QSTRING? attribs_or_modifiers?;

model_name          :  Name COLON ;
inplace_table_def   :  COLON INDENT (Name attribs_or_modifiers?)+ DEDENT;
table_refs          :  (TABLE | TYPE) Name inplace_table_def?;
facade              :  SYSL_COMMENT* WRAP model_name INDENT table_refs+ DEDENT;

documentation_stmts     : AT Name EQ QSTRING NEWLINE;

var_in_curly    : CURLY_OPEN Name CURLY_CLOSE;
query_var       : Name EQ (Name | NativeDataTypes | var_in_curly);
query_param     : QN query_var (AMP query_var)*;

http_path_var_with_type : CURLY_OPEN Name LESS_COLON (NativeDataTypes | Name) CURLY_CLOSE;
http_path_static : Name;
http_path_suffix : FORWARD_SLASH (http_path_static | http_path_var_with_type);
http_path       : FORWARD_SLASH | http_path_suffix+;

endpoint_name   : ( DOT? Name+ | TEXT_LINE);

ret_stmt        : RETURN TEXT;

target          : (DOT | app_name);
target_endpoint : (TEXT_NAME | Name);
call_stmt       : target (MEMBER target_endpoint);

if_stmt                 : IF TEXT_NAME COLON INDENT statements* DEDENT;
if_else                 : if_stmt (ELSE COLON INDENT statements* DEDENT)?;

for_cond                : TEXT_NAME COLON;

for_stmt                : FOR for_cond
                                INDENT statements* DEDENT;

http_method_comment     : SYSL_COMMENT;
//group_stmt              : GROUP TEXT_NAME COLON
//                                INDENT statements+ DEDENT;

one_of_case_label: Name*;

one_of_cases: one_of_case_label COLON
                    INDENT statements+ DEDENT;

one_of_stmt             : ONE_OF COLON
                           INDENT one_of_cases+ DEDENT;

text_stmt               : doc_string | TEXT_LINE | WHATEVER;

statements: ( if_else
                | for_stmt
                | ret_stmt
                | call_stmt
                | one_of_stmt
                | http_method_comment
                // | group_stmt
                | QSTRING
                | text_stmt
                | annotation
            )
            attribs_or_modifiers?
            ;

method_def: HTTP_VERBS
                  query_param? attributes? COLON
                        INDENT statements+ DEDENT
                ;

shortcut        : WHATEVER;

simple_endpoint :
                WHATEVER
                | (
                    endpoint_name QSTRING? attribs_or_modifiers? COLON
                    (  shortcut
                        | (INDENT statements+ DEDENT)
                    )
                )
                ;


rest_endpoint: http_path attribs_or_modifiers? COLON
                                    (INDENT (method_def | rest_endpoint)+ DEDENT)
                ;

collector_stmt: call_stmt
                | (HTTP_VERBS http_path);

collector_stmts: collector_stmt attribs_or_modifiers;

collector:  COLLECTOR COLON (WHATEVER | (INDENT collector_stmts+ DEDENT));

event: DISTANCE EVENT_NAME
        attribs_or_modifiers? COLON (WHATEVER | INDENT statements+ DEDENT);

subscribe: Name SUBSCRIBE attribs_or_modifiers? COLON (WHATEVER | INDENT statements+ DEDENT);

app_decl: INDENT  (table | facade | SYSL_COMMENT | rest_endpoint | simple_endpoint | collector | event | subscribe | annotation )+ DEDENT;

application:  SYSL_COMMENT*
                name_with_attribs
                COLON
                app_decl
                ;

path            : FORWARD_SLASH? Name (FORWARD_SLASH Name)* ;
import_stmt     : IMPORT SYSL_COMMENT*;
imports_decl    : import_stmt+;
sysl_file       : imports_decl? application+ EOF;
