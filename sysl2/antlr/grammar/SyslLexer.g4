lexer grammar SyslLexer;

tokens { INDENT, DEDENT}

@lexer::header {
    import "encoding/json"
}

@lexer::members {

var spaces int
var linenum int
var in_sq_brackets int

var gotNewLine bool
var gotHttpVerb bool
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
                    { in_sq_brackets == 0 }?
                    ((I N T) |( S T R I N G) | (D A T E) | (B O O L) | (D E C I M A L) | (D A T E T I M E) | (X M L));

HTTP_VERBS          : ('GET' | 'POST' | 'DELETE' | 'PUT' | 'PATCH' )
                    { gotHttpVerb = true; }
                    ;

WRAP                : '!wrap';
TABLE               : '!table';
TYPE                : '!type';

fragment
IMPORT_KEY: 'import';

fragment
SUB_PATH_NAME: ~[ \r\n\t\\/:]+ ;

IMPORT              : IMPORT_KEY ' '+ (SUB_PATH_NAME |   ('/' SUB_PATH_NAME)+) [ \t]* NEWLINE
                     { gotNewLine = true; spaces=0; gotHttpVerb=false;linenum++;}
                    ;

RETURN              : ( R E T U R N )           -> pushMode(NOT_NEWLINE); //revisit this?
IF                  : (I F)  [ \t]*             -> pushMode(ARGS);
FOR_EACH            : (F O R) [ \t]* (E A C H) [ \t]* -> pushMode(ARGS);
FOR                 : (F O R) [ \t]*            -> pushMode(ARGS);
UNTIL               : (U N T I L) [ \t]*        -> pushMode(ARGS);
ELSE                : (E L S E);
LOOP                : (L O O P);
//GROUP               : ('Group' | 'group') -> pushMode(FREE_TEXT_NAME);
WHATEVER            : '...';
DOTDOT              : '..';
SET_OF              : S E T [ \t]* O F;
ONE_OF              : O N E [ \t]* O F      ;//-> pushMode(FREE_TEXT_NAME);
MIXIN               : '-' '|' '>';
DISTANCE            : '<->';
DOT_ARROW           : '.' [ \t]+ '<-' -> pushMode(ARGS); // for " '. <-' name" syntax, change mode to all  ". <- GET /rest/api/calls"
NAME_SEP            : [ \t]* '::' [ \t]*;
LESS_COLON          : '<:';
ARROW_LEFT          : '<-'  -> pushMode(ARGS); // Added for: 'server <- GET /http/path' calls
ARROW_RIGHT         : [ \t]* '->' [ \t]* ;
COLLECTOR           : '.. * <- *';
PLUS                : '+';
TILDE               : '~';
COMMA               : ',';
EQ                  : '=';
DOLLAR              : '$';
FORWARD_SLASH       : '/';
STAR                : '*';
COLON               : ':';
PERCENT             : '%';
DOT                 : '.';
EXCLAIM             : '!';
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
DBL_QT              : ["];
SINGLE_QT           : ['];


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
WITHIN_DBL_QTS        : (~[\r\n"])*;

fragment
WITHIN_SNGL_QTS        : (~[\r\n'])*;

QSTRING: (
            (DBL_QT WITHIN_DBL_QTS DBL_QT)
            |
            (SINGLE_QT WITHIN_SNGL_QTS SINGLE_QT)
        )
        {
            var val string
            if json.Unmarshal([]byte(l.GetText()), &val) == nil {
                l.SetText(val)
            } else {
                val =l.GetText()[1:]
                val =val[:len(val)-1]
                l.SetText(val)
            }
        }
    ;

NEWLINE             : '\r'? '\n'
                    {gotNewLine = true; gotHttpVerb=false; spaces=0; linenum++;}
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

Name            : [a-zA-Z][-a-zA-Z0-9_]*;
/// end--textline & name

// cim.sysl has spaces and tab in the same line.
// TAB             : [\t]+
//                 { setChannel(HIDDEN);   spaces += (getText().length() * 4);}
//                 ;

WS              : [ \t]+
                { spaces =  calcSpaces(l.GetText());}
                -> channel(HIDDEN)
                ;


mode ARGS;
SKIP_WS_ARG         : [ ]   -> skip;
LESS_COLON_2          : '<:';
SQ_OPEN_2             : '['   { in_sq_brackets++;} -> popMode;

Q_ARG: (
            (DBL_QT WITHIN_DBL_QTS DBL_QT)
            |
            (SINGLE_QT WITHIN_SNGL_QTS SINGLE_QT)
        );

TEXT_VALUE      : (~[,'"()\r\n:[\]<])+;
OPEN_PAREN_ARG  : '(';
CLOSE_PAREN_ARG : ')'   -> popMode;
COLON_ARG       : ':'   -> popMode;
COMMA_ARG       : ',' [ \t]*;

NEWLINE_2           : '\r'? '\n'
                    {gotNewLine = true; gotHttpVerb=false; spaces=0; linenum++;}
                    -> channel(HIDDEN), popMode
                    ;

mode NOT_NEWLINE;
TEXT            : (~[\r\n])+        -> popMode ;

// mode FREE_TEXT_NAME;
// SKIP_WS_1         : [ ]   -> skip;
// TEXT_NAME       : ~['"()\r\n:[\]<]+  -> popMode;

mode AT_VAR_DECL;
POP_WS          : [ ]   -> skip, popMode;
VAR_NAME        : [a-zA-Z][a-zA-Z0-9._-]* -> popMode;
