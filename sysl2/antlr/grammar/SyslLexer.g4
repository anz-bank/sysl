lexer grammar SyslLexer;

tokens { INDENT, DEDENT}

// @header {
// // import (
// //     "anz.com/sysl/antlr/impl"
// // )
// // }

@lexer::members {

var spaces int
var linenum int
var in_sq_brackets int

var gotNewLine bool
var prevTokenIndex = -1

func (l *SyslLexer) NextToken() antlr.Token {
    return GetNextToken(l)
}

}
//     java.util.LinkedList<Token> prevToken = new java.util.LinkedList<Token>();
//     static java.util.ArrayList<String> keywords = new java.util.ArrayList<String>();

//     static {
//         // level 0
//         keywords.add("import");
//         // level 1
//         keywords.add("set of");
//         // level 2
//         keywords.add("return");
//         keywords.add("for");
//         keywords.add("one of");
//         keywords.add("else");
//         keywords.add("if");
//         keywords.add("loop");
//         keywords.add("group");
//     }
    



//     @Override
// 	public Token nextToken() {
//         if(prevToken.size() > 0) {
//             Token next = prevToken.poll();
//             String tokenName = tokenName(next);
//             System.out.println(linenum + "-" + level.size() + " :Q " + tokenName);
//             return next;
//         }

//         Token next = super.nextToken();
//         String tokenName = tokenName(next);
        
//         // return NEWLINE
//         if(gotNewLine && next.getType() == NEWLINE)  {
//             System.out.println(linenum + "-" + level.size() + " :NL " + tokenName);
//             return next;
//         }
//         // regular whitespace, return as is.
//         // return from here only when we encounter HIDDEN after INDENT has been generated
//         // after processing NL.
//         if(!gotNewLine && next.getChannel() == HIDDEN) {
//             // System.out.println("returning hidden token " + spaces);
//             System.out.println(linenum + "-" + level.size() + " :HIDDEN " + tokenName);
//             spaces = 0;
//             return next;
//         }
//         else if (next.getType() == SYSL_COMMENT) {
//             System.out.println(linenum + "-" + level.size() + " :SYSL_COMMENT " + tokenName);
//             spaces = 0;
//             return next;            
//         }

//         if(next.getType() == EOF ) {
//             spaces = 0; // done with the file
//         }
//         else if(!gotNewLine) {
//             System.out.println(linenum + "-" + level.size() + " :NOT_NL " + tokenName + "(" +next.getText()+ ")" );
//             return next;
//         }

//         if(spaces != getPreviousIndent()) {
//             System.out.println(linenum + " :will_que " + tokenName);
//         }
//         while(spaces != getPreviousIndent()) {
            
//             if(spaces > getPreviousIndent()) {
//                 System.out.println(linenum + " :next level of indent : " + spaces);
//                 level.push(spaces);
//                 prevToken.offer(createIndentToken());
//             }
//             else {
//                 int popped = level.pop();
//                 System.out.println(linenum + " :decreasing indent" + popped);
//                 prevToken.offer(createDedentToken());
//             }
//         }

//         gotNewLine = false;
//         prevToken.offer(next);
//         Token temp = prevToken.poll();
//         tokenName = tokenName(temp);
//         System.out.println(linenum + "-" + level.size() + " :LAST " + tokenName + "(" +next.getText()+ ")" );
//         return temp;
//     }
// }


NativeDataTypes     : 'int' | 'string' | 'date' | 'bool' | 'decimal' ;
HTTP_VERBS          : 'GET' | 'POST' | 'DELETE' | 'PUT' | 'PATCH';

WRAP                : '!wrap';
TABLE               : '!table';
TYPE                : '!type';
IMPORT              : 'import'              -> pushMode(NOT_NEWLINE);
RETURN              : ('return' | 'RETURN') -> pushMode(NOT_NEWLINE); //revisit this?
IF                  : ([iI][fF])            -> pushMode(FREE_TEXT_NAME);
ELSE                : ('else' | 'Else')     ;//-> pushMode(FREE_TEXT_NAME);//  do we support else if()??

FOR                 : ('For' | 'for')       -> pushMode(FREE_TEXT_NAME);

LOOP                : 'Loop' | 'loop';
GROUP               : ('Group' | 'group') -> pushMode(FREE_TEXT_NAME);
WHATEVER            : '...';
SET_OF              : 'set of';
ONE_OF              : [oO]'ne of'      ;//-> pushMode(FREE_TEXT_NAME);
DISTANCE            : '<->'         -> pushMode(EVENT_NAME_MODE);
NAME_SEP            : '::'          -> pushMode(FREE_TEXT_NAME);
LESS_COLON          : '<:';
MEMBER              : '<-'          -> pushMode(FREE_TEXT_NAME);
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
OPEN_ANGLE          : '<';
CLOSE_ANGLE         : '>';
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
WITHIN_DBL_QTS        : (~[\r\t"])+;

fragment 
WITHIN_SNGL_QTS        : (~[\r\t'])+;

QSTRING: (DBL_QT WITHIN_DBL_QTS DBL_QT) | SINGLE_QT WITHIN_SNGL_QTS SINGLE_QT;

NEWLINE             : [\r\n]
                    {gotNewLine = true; spaces=0; linenum++;}
                     -> channel(HIDDEN)
                    ;

SYSL_COMMENT    : HASH TEXT -> channel(HIDDEN);

// '<' -> will be required for transformation syntax
// so add '>'
fragment
PRINTABLE       :   ~[ \r\n@!|'"#[:\\/<>\]]+;

// defined before Name
TEXT_LINE       :  PRINTABLE (' ' PRINTABLE)+
                { in_sq_brackets == 0 }?
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
VAR_NAME        : [a-zA-Z][a-zA-Z0-9_-]*;

mode EVENT_NAME_MODE;
SKIP_WS_2         : [ ]   -> skip;
EVENT_NAME       : ~[ ](~[\r\n:[\]])+  -> popMode;
