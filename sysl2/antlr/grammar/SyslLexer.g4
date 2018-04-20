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

NativeDataTypes     : 'int' | 'string' | 'date' | 'bool' | 'decimal' | 'datetime' ;
HTTP_VERBS          : ('GET' | 'POST' | 'DELETE' | 'PUT' | 'PATCH')
                    { gotHttpVerb = true}
                    ;

WRAP                : '!wrap';
TABLE               : '!table';
TYPE                : '!type';

fragment
IMPORT_KEY: 'import';

fragment
SUB_PATH_NAME: ~[ \r\n\t\\/:]+ ;

IMPORT              : IMPORT_KEY ' '+ (SUB_PATH_NAME |   ('/' SUB_PATH_NAME)+) [ \t]* NEWLINE ;

RETURN              : ( [rR][eE][tT][uU][rR][nN] )  -> pushMode(NOT_NEWLINE); //revisit this?
IF                  : ([iI][fF])                    -> pushMode(FREE_TEXT_NAME);
ELSE                : ([eE][lL][sS][eE]);
FOR                 : ( [fF][oO][rR])               -> pushMode(FREE_TEXT_NAME);
LOOP                : [lL][oO][oO][pP];
//GROUP               : ('Group' | 'group') -> pushMode(FREE_TEXT_NAME);
WHATEVER            : '...';
SET_OF              : 'set of';
ONE_OF              : [oO]'ne of'      ;//-> pushMode(FREE_TEXT_NAME);
DISTANCE            : '<->'         -> pushMode(EVENT_NAME_MODE);
NAME_SEP            : '::'          -> pushMode(FREE_TEXT_NAME);
LESS_COLON          : '<:';
ARROW_LEFT          : '<-'          -> pushMode(FREE_TEXT_NAME);
ARROW_RIGHT         : '->';
COLLECTOR           : '.. * <- *';
PLUS                : '+';
TILDE               : '~';
COMMA               : ',';
EQ                  : '=';
DOLLAR              : '$';
FORWARD_SLASH       : '/';
MINUS               : '-';
STAR                : '*';
COLON               : ':';
PERCENT             : '%';
DOT                 : '.';
EXCLAIM             : '!';
QN                  : '?';
AT                  : '@'       -> pushMode(AT_VAR_DECL);
AMP                 : '&';
SQ_OPEN             : '['   { in_sq_brackets++;};
SQ_CLOSE            : ']'   { in_sq_brackets--;};
CURLY_OPEN          : '{';
CURLY_CLOSE         : '}';
OPEN_PAREN          : '(';
CLOSE_PAREN         : ')';
// OPEN_ANGLE          : '<';
// CLOSE_ANGLE         : '>';
EMPTY_COMMENT       : ('#' '\r'? '\n')
                     {gotNewLine = true; spaces=0; gotHttpVerb=false;linenum++;}
                    -> channel(HIDDEN);

HASH                : '#'       -> pushMode(NOT_NEWLINE);
PIPE                : '|'       -> pushMode(NOT_NEWLINE);
DBL_QT              : ["];
SINGLE_QT           : ['];


EMPTY_LINE          : [ \t]+ [\r\n]
                    {  l.Skip(); spaces=0; linenum++;}
                    ;

// added '#' to skip comments that start with an indent
INDENTED_COMMENT    : [ \t]+ '#' ~[\n]+? '\n'
                    {  l.Skip(); spaces=0; linenum++;}
                    ;

DIGITS              : [1-9][0-9]*;

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
// '=' '{' '}' '&' '?'
fragment
PRINTABLE       :   ~[ \n\r!"#'&\-/:=<?@[\]{}|]+;

// defined before Name
TEXT_LINE       :
                { !gotHttpVerb}?
                PRINTABLE ([ &\-]+ PRINTABLE)+
                { in_sq_brackets == 0 }?
                { !gotHttpVerb}?
                { startsWithKeyword(p.GetText()) == false}?
                //  { (_input.LA(1) == '\n') | ((_input.LA(1) == '\r') && (_input.LA(2) == '\n')) }?
                ;

Name            : [a-zA-Z][a-zA-Z0-9_-]*;
/// end--textline & name

// cim.sysl has spaces and tab in the same line.
// TAB             : [\t]+
//                 { setChannel(HIDDEN);   spaces += (getText().length() * 4);}
//                 ;

WS              : [ \t]+
                { spaces =  calcSpaces(l.GetText());}
                -> channel(HIDDEN)
                ;


mode NOT_NEWLINE;
TEXT            : (~[\r\n])+        -> popMode ;

mode FREE_TEXT_NAME;
SKIP_WS         : [ ]   -> skip;
TEXT_NAME       : ~[ ](~[\r\n:[\]<])+  -> popMode;

// either add '=' into TEXT_LINE
// or have this special mode
mode AT_VAR_DECL;
POP_WS          : [ ]   -> skip, popMode;
VAR_NAME        : [a-zA-Z][a-zA-Z0-9._-]*;

mode EVENT_NAME_MODE;
SKIP_WS_2         : [ ]   -> skip;
EVENT_NAME       : ~[ ](~[\r\n:[\]])+  -> popMode;
