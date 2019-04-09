lexer grammar SyslLexer;

tokens { INDENT, DEDENT}

@lexer::members {

var keywords = ["sequence of", "set of", "return", "for", "one of", "else", "if", "loop", "until", "alt", "while"];

SyslLexer.prototype.linenum = 0;
SyslLexer.prototype.in_sq_brackets = false;
SyslLexer.prototype.gotHttpVerb = false;
SyslLexer.prototype.gotNewLine = false;
SyslLexer.prototype.gotView = false;
SyslLexer.prototype.prevToken = [];
SyslLexer.prototype.spaces = 0;
SyslLexer.prototype.parens = 0;
SyslLexer.prototype.level = []; // stack

function calcSpaces(text) {
    if (text == undefined) return 0;

	let s = 0;

	for (i = 0; i < text.length; i++) {
		if (text[i] == ' ') {
			s++
		}
		if (text[i] == '\t') {
			s += 4
		}
	}
	return s
}

function startsWithKeyword(s) {
	var lower = s.toLowerCase(s)

	for ( k in keywords) {
		if (lower.indexOf(keywords[k]) == 0) {
			return true
		}
	}
	return false
}

function createDedentToken(source) {
	let t = new antlr4.CommonToken(source, SyslLexer.DEDENT, 0, 0, 0)
	return t
}

function createIndentToken(source) {
	let t = new antlr4.CommonToken(source, SyslLexer.INDENT, 0, 0, 0)
	return t
}

function getPreviousIndent(level) {
	if (level.length == 0) {
		return 0
	}
	// peek, read but not remove HEAD
	return level[level.length-1]
}

SyslLexer.prototype.nextToken = function () {
	if (this.prevToken.length > 0) {
		// poll, retrieve head
		let nextTok = this.prevToken.shift();
    let tokenName = this.symbolicNames[nextTok.type]
    //console.log("Q: " + tokenName + " " + nextTok.text);
		return nextTok
	}

    let next = antlr4.Lexer.prototype.nextToken.call(this);
    let tokenName = this.symbolicNames[next.type]

	// return NEWLINE
    if (this.gotNewLine &&
        (next.type == SyslLexer.NEWLINE
            || next.type == SyslLexer.NEWLINE_2
            || next.type == SyslLexer.E_NL
            || next.type == SyslLexer.EMPTY_LINE
            || next.type == SyslLexer.INDENTED_COMMENT
            || next.type == SyslLexer.E_INDENTED_COMMENT
            || next.type == SyslLexer.E_EMPTY_LINE
            || next.type == SyslLexer.E_DOT_NAME_NL
            || next.type == SyslLexer.EMPTY_COMMENT)) {
		return next
    }

	// regular whitespace, return as is.
	// return from here only when we encounter HIDDEN after INDENT has been generated
	// after processing NL.
	if (!this.gotNewLine && next.channel == antlr4.Token.HIDDEN_CHANNEL) {
		this.spaces = 0
		return next
	} else if (next.type == SyslLexer.SYSL_COMMENT) {
		this.spaces = 0
		return next
	}

	if (next.type == antlr4.Token.EOF) {
		this.spaces = 0 // done with the file
	} else if (!this.gotNewLine) {
    //console.log("NOT NL: " + tokenName + " " + next.text);
		return next
	}

	while(this.spaces != getPreviousIndent(this.level)) {
		if (this.spaces > getPreviousIndent(this.level)) {
			this.level.push(this.spaces)
			this.prevToken.push(createIndentToken(next.source))
		} else {
			this.level.pop()
			this.prevToken.push(createDedentToken(next.source))
		}
	}

	this.gotNewLine = false
	this.prevToken.push(next);
	var temp = this.prevToken.shift();
	return temp
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
                    { this.in_sq_brackets == 0 }?
                    ;

HTTP_VERBS          : ('GET' | 'POST' | 'DELETE' | 'PUT' | 'PATCH' ) [ \t]*
                    { this.gotHttpVerb = true; }
                    ;

WRAP                : '!wrap';
TABLE               : '!table';
TYPE                : '!type';
ALIAS               : '!alias';
UNION               : '!union';
VIEW                : '!view' { this.gotView = true;};

fragment
IMPORT_KEY: 'import';

fragment
SUB_PATH_NAME: ~[ \r\n\t\\/:]+ ;

IMPORT              : IMPORT_KEY ' '+ (SUB_PATH_NAME |   ('/' SUB_PATH_NAME)+) [ \t]* NEWLINE
                     { this.gotNewLine = true; this.spaces=0; this.gotHttpVerb=false; this.linenum++;}
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
                    { this.gotView}?
                    '[~abstract]'
                    { this.gotView = false;}
                    ;
TILDE               : '~';
COMMA               : ',';
EQ                  : '=';
FORWARD_SLASH       : '/'
                    { this.gotHttpVerb = true; }
                    ;
COLON               : ':';
DOT                 : '.';
QN                  : '?';
AT                  : '@'       -> pushMode(AT_VAR_DECL);
AMP                 : '&' { this.gotHttpVerb }? ;
SQ_OPEN             : '['   { this.in_sq_brackets++;};
SQ_CLOSE            : ']'   { this.in_sq_brackets--;};
CURLY_OPEN          : '{';
CURLY_CLOSE         : '}';
OPEN_PAREN          : '(';
CLOSE_PAREN         : ')';
// OPEN_ANGLE          : '<';
// CLOSE_ANGLE         : '>';
EMPTY_COMMENT       : ('#' '\r'? '\n')
                    { this.gotNewLine = true; this.spaces=0; this.gotHttpVerb=false; this.linenum++;}
                    -> channel(HIDDEN)
                    ;
HASH                : '#'       -> pushMode(NOT_NEWLINE);
PIPE                : '|'       -> pushMode(NOT_NEWLINE);

EMPTY_LINE          : ([ \t]+ ( [\r\n] | EOF ))
                    { this.gotNewLine = true; this.spaces=0; this.gotHttpVerb=false; this.linenum++;}
                    -> channel(HIDDEN)
                    ;
// added '#' to skip comments that start with an indent
INDENTED_COMMENT    : ([ \t]+ '#' ~[\n]* ('\n' | EOF))
                    { this.gotNewLine = true; this.spaces=0; this.gotHttpVerb=false; this.linenum++; }
                    -> channel(HIDDEN)
                    ;
DIGITS              : [0-9][0-9]*;

fragment
DOUBLE_QUOTE_STRING: ["] (~["\\] | [\\][\\brntu'"])* ["];
fragment
SINGLE_QUOTE_STRING: ['] (~['])* ['];

QSTRING     : DOUBLE_QUOTE_STRING | SINGLE_QUOTE_STRING;

NEWLINE     : '\r'? '\n'
            {this.gotNewLine = true; this.gotHttpVerb=false; this.spaces=0; this.linenum++;}
            { if (this.gotView) { this.pushMode(SyslLexer.VIEW_TRANSFORM);}}
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
                { this.in_sq_brackets == 0 }?
                { startsWithKeyword(this.text) == false}?
                ;

Name            : [a-zA-Z_][-a-zA-Z0-9_]*;
/// end--textline & name

// cim.sysl has spaces and tab in the same line.
// TAB             : [\t]+
//                 { setChannel(HIDDEN);   spaces += (getText().length() * 4);}
//                 ;

WS              : [ \t]+
                { this.spaces =  calcSpaces(this.text);}
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
                { this.type = SyslLexer.Name}
                { this.text = this.text.trim()}
                -> popMode
                ;

NEWLINE_2           : '\r'? '\n'
                    {this.gotNewLine = true; this.gotHttpVerb=false; this.spaces=0; this.linenum++;}
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
                    { this.type = SyslLexer.NativeDataTypes;}
                    ;

E_INDENTED_COMMENT    : ([ \t]+ '#' ~[\n]* ('\n' | EOF))
                    { this.gotNewLine = true; this.spaces=0; this.gotHttpVerb=false; this.linenum++; }
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
E_OPEN_PAREN    : '(' { this.parens++;};
E_CLOSE_PAREN   : ')' { this.parens--;};
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
                { this.spaces > 1 }?
                NAME? '.' NAME [ \t]* '\r'? '\n'
                { this.gotNewLine = true; this.gotHttpVerb=false; this.linenum++;}
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
                { this.spaces =  calcSpaces(this.text);}
                -> channel(HIDDEN)
                ;

E_EMPTY_LINE    :
                {this.gotNewLine}?
                [ \t]* (( '\r'? '\n') | EOF )
                { this.gotNewLine = true; this.spaces=0; this.gotHttpVerb=false; this.linenum++;}
                -> channel(HIDDEN)
                ;

E_NL  : '\r'? '\n'
      {this.gotNewLine = true; this.gotHttpVerb=false; this.spaces=0; this.linenum++; if (this.parens==0) {this.gotView=false; this.popMode();}}
      ;
