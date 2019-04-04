lexer grammar SyslLexer;

tokens { INDENT, DEDENT}

@lexer::members {

var spaces int
var linenum int
var in_sq_brackets int
var parens int

var gotNewLine bool
var gotHttpVerb bool
var gotView bool
var prevTokenIndex = -1

func (l *SyslLexer) NextToken() antlr.Token {
    return GetNextToken(l)
}

}

fragment A : [aA]; // match either an 'a' or 'A'
fragment B : [bB];
fragment C : [cC];
fragment D : [dD];
fragment E : [eE];
fragment F : [fF];
fragment G : [gG];
fragment H : [hH];
fragment I : [iI];
fragment J : [jJ];
fragment K : [kK];
fragment L : [lL];
fragment M : [mM];
fragment N : [nN];
fragment O : [oO];
fragment P : [pP];
fragment Q : [qQ];
fragment R : [rR];
fragment S : [sS];
fragment T : [tT];
fragment U : [uU];
fragment V : [vV];
fragment W : [wW];
fragment X : [xX];
fragment Y : [yY];
fragment Z : [zZ];

NativeDataTypes     :
                    ( (I N T '3' '2') | (I N T '6' '4') | (I N T) | (F L O A T) | ( S T R I N G) | (D A T E) | (B O O L) | (D E C I M A L) | (D A T E T I M E) | (X M L) | (A N Y))
                    { in_sq_brackets == 0 }?
                    ;

HTTP_VERBS          : ('GET' | 'POST' | 'DELETE' | 'PUT' | 'PATCH' ) [ \t]*
                    { gotHttpVerb = true; }
                    ;

WRAP                : '!wrap';
TABLE               : '!table';
TYPE                : '!type';
ALIAS               : '!alias';
UNION               : '!union';
VIEW                : '!view' { gotView = true;};

fragment
IMPORT_KEY: 'import';

fragment
SUB_PATH_NAME: ~[ \r\n\t\\/:]+ ;

IMPORT              : IMPORT_KEY ' '+ (SUB_PATH_NAME |   ('/' SUB_PATH_NAME)+) [ \t]* NEWLINE
                     { gotNewLine = true; spaces=0; gotHttpVerb=false;linenum++;}
                    ;

RETURN              : ( R E T U R N )           -> pushMode(NOT_NEWLINE); //revisit this?
IF                  : (I F)  [ \t]*             -> pushMode(PREDICATE);
FOR_EACH            : (F O R) [ \t]* (E A C H) [ \t]* -> pushMode(PREDICATE);
FOR                 : (F O R) [ \t]*            -> pushMode(PREDICATE);
UNTIL               : (U N T I L) [ \t]*        -> pushMode(PREDICATE);
ELSE                : (E L S E) [ \t]*          -> pushMode(PREDICATE);
LOOP                : (L O O P) [ \t]*          -> pushMode(PREDICATE);
ALT                 : (A L T) [ \t]*            -> pushMode(PREDICATE);
WHILE                : (W H I L E) [ \t]*       -> pushMode(PREDICATE);
WHATEVER            : '...';
DOTDOT              : '..';
SEQUENCE_OF         : S E Q U E N C E [ \t]* O F;
SET_OF              : S E T [ \t]* O F;
ONE_OF              : O N E [ \t]* O F      ;//-> pushMode(FREE_TEXT_NAME);
MIXIN               : '-' '|' '>';
DISTANCE            : '<->';
DOT_ARROW           : '.' [ \t]+ '<-' -> pushMode(ARGS); // for " '. <-' name" syntax, change mode to all  ". <- GET /rest/api/calls"
NAME_SEP            : [ \t]* '::' [ \t]*;
LESS_COLON          : [ \t]* '<:' [ \t]*;
ARROW_LEFT          : '<-' [ \t]* -> pushMode(ARGS); // Added for: 'server <- GET /http/path' calls
ARROW_RIGHT         : [ \t]* '->' [ \t]* ;
COLLECTOR           : '.. * <- *';
PLUS                : '+';
ABSTRACT            :
                    { gotView}?
                    '[~abstract]'
                    { gotView = false;}
                    ;
TILDE               : '~';
COMMA               : ',';
EQ                  : '=';
FORWARD_SLASH       : '/'
                    { gotHttpVerb = true; }
                    ;
COLON               : ':';
DOT                 : '.';
QN                  : '?';
AT                  : '@'       -> pushMode(AT_VAR_DECL);
AMP                 : '&' { gotHttpVerb }? ;
SQ_OPEN             : '['   { in_sq_brackets++;};
SQ_CLOSE            : ']'   { in_sq_brackets--;};
CURLY_OPEN          : '{';
CURLY_CLOSE         : '}';
OPEN_PAREN          : '(';
CLOSE_PAREN         : ')';
// OPEN_ANGLE          : '<';
// CLOSE_ANGLE         : '>';
EMPTY_COMMENT       : ('#' '\r'? '\n')
                    { gotNewLine = true; spaces=0; gotHttpVerb=false;linenum++;}
                    -> channel(HIDDEN)
                    ;
HASH                : '#'       -> pushMode(NOT_NEWLINE);
PIPE                : '|'       -> pushMode(NOT_NEWLINE);

EMPTY_LINE          : ([ \t]+ ( [\r\n] | EOF ))
                    { gotNewLine = true; spaces=0; gotHttpVerb=false; linenum++;}
                    -> channel(HIDDEN)
                    ;
// added '#' to skip comments that start with an indent
INDENTED_COMMENT    : ([ \t]+ '#' ~[\n]* ('\n' | EOF))
                    { gotNewLine = true; spaces=0; gotHttpVerb=false; linenum++; }
                    -> channel(HIDDEN)
                    ;
DIGITS              : [0-9][0-9]*;

fragment
DOUBLE_QUOTE_STRING: ["] (~["\\] | [\\][\\brntu'"])* ["];
fragment
SINGLE_QUOTE_STRING: ['] (~['])* ['];

QSTRING     : DOUBLE_QUOTE_STRING | SINGLE_QUOTE_STRING;

NEWLINE     : '\r'? '\n'
            {gotNewLine = true; gotHttpVerb=false; spaces=0; linenum++;}
            { if (gotView) { l.PushMode(SyslLexerVIEW_TRANSFORM);}}
            -> channel(HIDDEN)
            ;

SYSL_COMMENT    : HASH TEXT -> channel(HIDDEN);

// '<-', required for transformation syntax
// '->', required for events
// '/', for rest api
// ':', for everything
// '=' '{' '}'  '?'
// '(' ',' ')' for passing params
// removed '&' as its only required when we get HTTP verb
// removed '=' as its required only inside sq brackets
// add '.', required for decimal, reference syntax 'app.epname'
//      DOT can be in the app or epname!!
fragment
PRINTABLE       :   ~[ \t.\-<>,()\n\r!"#'/:?@[\]{}|]+;

fragment
IN_ANGLE        : '<' PRINTABLE '>';

TEXT_LINE       :
                PRINTABLE ([ \-]+ (PRINTABLE | IN_ANGLE))+
                { in_sq_brackets == 0 }?
                { startsWithKeyword(p.GetText()) == false}?
                ;

Name            : [a-zA-Z_][-a-zA-Z0-9_]*;
/// end--textline & name

// cim.sysl has spaces and tab in the same line.
// TAB             : [\t]+
//                 { setChannel(HIDDEN);   spaces += (getText().length() * 4);}
//                 ;

WS              : [ \t]+
                { spaces =  calcSpaces(l.GetText());}
                -> channel(HIDDEN)
                ;

ErrorChar   :
            .
            ;

mode PREDICATE;
PREDICATE_VALUE      : (~[\r\n:])* -> popMode;

mode ARGS;
SKIP_WS_ARG         : [ ]   -> skip;

fragment
TVALUE          : (~[()\][\r\n'"])+
                ;

TEXT_VALUE      : TVALUE
                { l.SetType(SyslLexerName)}
                { l.SetText(TrimText(l))}
                -> popMode
                ;

NEWLINE_2           : '\r'? '\n'
                    {gotNewLine = true; gotHttpVerb=false; spaces=0; linenum++;}
                    -> channel(HIDDEN), popMode
                    ;

mode NOT_NEWLINE;
TEXT            : (~[\r\n])*        -> popMode ;

mode AT_VAR_DECL;
POP_WS          : [ ]   -> skip, popMode;
VAR_NAME        : [a-zA-Z][a-zA-Z0-9._-]* -> popMode;

mode VIEW_TRANSFORM;
// Add 'any' if required.
E_NativeDataTypes     :
                    ( (I N T '3' '2') | (I N T '6' '4') | (I N T) | (F L O A T) | ( S T R I N G) | (D A T E) | (B O O L) | (D E C I M A L) | (D A T E T I M E) | (X M L))
                    { l.SetType(SyslLexerNativeDataTypes)}
                    ;

E_INDENTED_COMMENT    : ([ \t]+ '#' ~[\n]* ('\n' | EOF))
                    { gotNewLine = true; spaces=0; gotHttpVerb=false; linenum++; }
                    -> channel(HIDDEN)
                    ;
E_WHATEVER      : '...';
E_EQ_RIGHT      : '=>';
E_ARROW_RIGHT   : '->';
E_DOUBLE_EQ     : '==';
E_REL           : '!'?'in'|'!'?'contains'|'!='|'<='|'>=';
E_SQ_OPEN       : '[';
E_SQ_CLOSE      : ']';
E_ANGLE_OPEN    : '<';
E_ANGLE_CLOSE   : '>';
E_EMPTY_TUPLE   : '{:}';
E_CURLY_OPEN    : '{';
E_CURLY_CLOSE   : '}';
E_NULLSAFE_DOT  : '?.';
E_TABLE_OF      : 'table' [ \t]+ 'of';
E_POW           : '**';
E_COALESCE      : [ \t]* '??' [ \t]*;
E_COLON         : ':';
E_OPEN_PAREN    : '(' { parens++;};
E_CLOSE_PAREN   : ')' { parens--;};
E_COMMA         : ',';
E_EQ            : '=';
E_PLUS          : '+';
E_DIVIDE        : '/';
E_MOD           : '%';
E_MINUS         : '-';
E_QN            : '?';
E_TILDE         : '~';
E_NOT           : '!';
E_XOR           : '^' | 'xor';
E_LOGIC_OR      : '||';
E_DOUBLE_AMP    : '&&';
E_AMP           : '&';
E_BY            : 'by';
E_AND           : 'and';
E_BITOR         : 'or'  | '|';
E_STAR          : '*';
E_AS            : 'as';
E_VIA           : 'via';
E_IF            : 'if';
E_THEN          : 'then';
E_ELSE          : 'else';
E_LET           : 'let';
E_TRUE          : 'true';
E_FALSE         : 'false';
E_NULL          : 'null';
E_BUTNOT        : 'but' [ \t]+ 'not';
E_SEQUENCE_OF   : 'sequence' [ \t]+ 'of' ;
E_SET_OF        : 'set' [ \t]+ 'of' ;
E_ASC           : 'asc';
E_DESC          : 'desc';
E_RELOPS_RANK          : 'rank';
E_RELOPS_AGG           : 'sum'|'min'|'max'|'average';
E_RELOPS_ANY           : 'any';
E_RELOPS_SINGLE_NULL   : 'singleOrNull'|'single';
E_RELOPS_SNAPSHOT      : 'snapshot';
E_RELOPS_WHERE         : 'where';
E_RELOPS_COUNT         : 'count';
E_RELOPS_FLATTEN       : 'flatten';
E_RELOPS_FIRST         : 'first';
E_FUNC          : 'autoinc' | 'str' | 'substr';

E_STRING_DBL           : ["] (~["\\] | [\\][\\brntu'"])* ["];
E_STRING_SINGLE        : ['] ~[']* ['];

fragment
F_DIGITS   : [0-9];

E_DECIMAL       : F_DIGITS F_DIGITS* '.' F_DIGITS F_DIGITS*;
E_DIGITS        : F_DIGITS F_DIGITS* ;

fragment
NAME            : [a-zA-Z_][a-zA-Z0-9_-]*;

E_DOT_NAME_NL    :
                { spaces > 1 }?
                NAME? '.' NAME [ \t]* '\r'? '\n'
                { gotNewLine = true; gotHttpVerb=false; linenum++;}
                ;
// does not work for
//  'input.table of....'
// E_RefName       : NAME ('.' NAME)+
//                 ;
E_Name          : NAME;
// expr = '.'
// want a greedy
E_DOT           : '.';
E_WS            : [ \t]+
                { spaces =  calcSpaces(l.GetText());}
                -> channel(HIDDEN)
                ;

E_EMPTY_LINE    :
                {gotNewLine}?
                [ \t]* (( '\r'? '\n') | EOF )
                { gotNewLine = true; spaces=0; gotHttpVerb=false; linenum++;}
                -> channel(HIDDEN)
                ;

E_NL  : '\r'? '\n'
      {gotNewLine = true; gotHttpVerb=false; spaces=0; linenum++;}
      {if (parens==0) { gotView=false; l.PopMode();} }
      ;
