parser grammar SyslParser;

options { tokenVocab=SyslLexer; }

modifier        : TILDE Name (PLUS Name)*;
size_spec       : OPEN_PAREN DIGITS ( DOT DIGITS)? CLOSE_PAREN;
modifier_list   : modifier (COMMA modifier)*;
modifiers       : SQ_OPEN modifier_list SQ_CLOSE;
reference       : parent_ref=Name '.' member=Name
                // {if(tables.contains($parent.text)) throw new RuntimeException();}
                // {System.out.println("ref :"+ $parent.text + " . " + $member.text + " " + tables.contains($parent.text));}
                ;

doc_string      : PIPE TEXT;

// TODO: refactor value
// value: '"' Name ('.' Name)* '"';
// nvp: Name EQ value;
quoted_string       : QSTRING //{System.out.println("quoted_string :"+ $QSTRING.text);}
                        ;
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

field: Name LESS_COLON field_type  //{System.out.println("field :" +  $Name.text + " " + $field_type.text);}
                                ;

table   :       SYSL_COMMENT*
                (TABLE | TYPE)
                Name attribs_or_modifiers? COLON
                INDENT (SYSL_COMMENT | field )+ DEDENT
               // { tables.add( $Name.text);}
                // {System.out.println("table :"+ $Name.text);}
        ;


package_name   : Name | TEXT_LINE | TEXT_NAME;
sub_package    : NAME_SEP package_name;
app_name       : package_name sub_package*;

name_with_attribs       :   app_name  QSTRING? attribs_or_modifiers?
                        //     {System.out.println("Name (app_name):"+ $app_name.text);}
                        ;

model_name          :  Name COLON ;
inplace_table_def   :  COLON INDENT (Name attribs_or_modifiers?)+ DEDENT;
table_refs          :  (TABLE | TYPE) Name inplace_table_def?;
facade              :  SYSL_COMMENT* WRAP model_name INDENT table_refs+ DEDENT;


documentation_stmts     : AT Name EQ QSTRING NEWLINE;
// added Name for next1.sysl, GET http path
variable_substitution   : FORWARD_SLASH CURLY_OPEN (field | Name) CURLY_CLOSE;

static_path     : FORWARD_SLASH (MINUS | Name)+;
query_var       : Name EQ (Name | NativeDataTypes);
query_param     : QN query_var (AMP query_var)*;
http_path       : (static_path | variable_substitution)+ query_param?;

endpoint_name   : (http_path | DOT? Name+ | TEXT_LINE) QSTRING? attribs_or_modifiers? COLON;

ret_stmt        : RETURN TEXT;

target          : (DOT | app_name);
target_endpoint : (TEXT_NAME | Name);
call_stmt       : target (MEMBER target_endpoint)?;

if_stmt                 : IF TEXT_NAME COLON INDENT http_statements* DEDENT;
if_else                 : if_stmt (ELSE COLON INDENT http_statements* DEDENT)?;

for_cond                : TEXT_NAME COLON;

for_stmt                : FOR for_cond
                                INDENT http_statements* DEDENT;

http_method_comment     : SYSL_COMMENT;
group_stmt              : GROUP TEXT_NAME COLON
                                INDENT http_statements+ DEDENT;

one_of_case_label: Name*;

one_of_cases: one_of_case_label COLON
                    INDENT http_statements+ DEDENT;

one_of_stmt             : ONE_OF COLON
                           INDENT one_of_cases+ DEDENT;

text_stmt               : TEXT_LINE;

http_statements: doc_string
                | if_else
                | for_stmt
                | ret_stmt
                | call_stmt
                | one_of_stmt
                | http_method_comment
                | group_stmt
                | QSTRING
                | WHATEVER
                | text_stmt
                | annotation
                ;

method_def [string url_path]     : HTTP_VERBS //{ System.out.println($HTTP_VERBS.text + " " + $url_path);}
                  query_param? attributes? COLON
                        INDENT http_statements* DEDENT
                ;

endpoint_decl [string prefix] : api_endpoint[$prefix]
                 | method_def[$prefix]
                 | http_statements;

shortcut        : WHATEVER  //    {System.out.println("shortcut");}
                        ;

api_endpoint [string prefix] :
                WHATEVER
                | (
                    endpoint_name
                    // { System.out.println("endpoint: " + $endpoint_name.text);}
                    (  shortcut
                    | (INDENT endpoint_decl[$prefix+$endpoint_name.text]+ DEDENT)
                    )
                )
                ;

collector_stmt: call_stmt
                | (HTTP_VERBS http_path);

collector_stmts: collector_stmt attribs_or_modifiers;

collector:  COLLECTOR COLON (WHATEVER | (INDENT collector_stmts+ DEDENT));
// event_name: Name | TEXT_LINE;

event: DISTANCE EVENT_NAME //{System.out.println("event_name :"+ $EVENT_NAME.text);}
        attribs_or_modifiers? COLON (WHATEVER | INDENT http_statements+ DEDENT);

app_decl: INDENT  (table | facade | SYSL_COMMENT | api_endpoint[""] | collector | event | annotation )+ DEDENT;

application:  SYSL_COMMENT*
                name_with_attribs      // {System.out.println("application :"+ $name_with_attribs.text);}
                COLON
                app_decl               // {System.out.println("app_decl");}
                ;

path            : FORWARD_SLASH? Name (FORWARD_SLASH Name)* ;
import_stmt     : IMPORT TEXT SYSL_COMMENT*;
imports_decl    : import_stmt+;
sysl_file       : imports_decl? application+ EOF;
