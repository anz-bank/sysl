goFile: PackageClause Comment? '\n' ImportDecl? '\n' TopLevelDecl+ '\n';
PackageClause: 'package' PackageName '\n';

ImportDecl: 'import' '(\n' ImportSpec* '\n)\n';
ImportSpec: (Import | NamedImport) '\n';
NamedImport: Name Import;
TopLevelDecl: Comment '\n' (Declaration | FunctionDecl | MethodDecl);
Declaration: VarDecl | VarDeclWithVal | ConstDecl | StructType | InterfaceType | AliasDecl;
StructType : 'type' StructName 'struct' '{\n' FieldDecl* '}\n\n';
FieldDecl: '\t' identifier Type? Tag? '\n';
IdentifierList: identifier IdentifierListC*;
IdentifierListC: ',' identifier;

VarDeclWithVal: 'var' identifier '=' Expression '\n';
VarDecl: 'var' identifier TypeName '\n';
ConstDecl: 'const' '(\n'  ConstSpec '\n)\n';
ConstSpec: VarName TypeName '=' ConstValue '\n';

FunctionDecl   : 'func' FunctionName Signature? Block '\n\n';
Signature: Parameters Result?;
Parameters: '(' ParameterList? ')';
Result         : ReturnTypes | TypeName;
ReturnTypes: '(' TypeName ResultTypeList* ')';
ResultTypeList: ',' TypeName ;
TypeList:  TypeName;
ParameterList     : ParameterDecl ParameterDeclC*;
ParameterDecl  : Identifier TypeName;
ParameterDeclC: ',' ParameterDecl;

InterfaceType      : 'type' InterfaceName 'interface'  '{\n'  MethodSpec* '}\n\n' MethodDecl*;
MethodSpec         : '\t' MethodName Signature '\n' | InterfaceTypeName ;
MethodDecl: 'func' Receiver FunctionName Signature? Block? '\n\n';
Receiver: '(' ReceiverType ')';
AliasDecl: 'type' identifier Type? '\n\n';

Block: '{\n'  StatementList* '}\n';
StatementList: '\t' Statement '\n';
Statement: ReturnStmt |  DeclareAndAssignStmt | AssignStmt | IfElseStmt | IncrementVarByStmt | FunctionCall | VarDecl | ForLoop;

AssignStmt: Variables '=' Expression;
IfElseStmt: 'if' Expression Block;
IncrementVarByStmt: Variables '+=' Expression;
ReturnStmt: 'return' (PayLoad | Expression);
DeclareAndAssignStmt: Variables ':=' Expression;

Expression: FunctionCall | NewStruct | GetArg |  ValueExpr | NewSlice | Map;

GetArg: LHS '.' RHS;
NewSlice: '[]' TypeName '{' SliceValues? '}';
Map: 'map[' KeyType ']' ValType '{\n' KeyValuePairs? '\n}';
KeyValuePairs: KeyValuePair*;
KeyValuePair: Key ':' Value ',\n';
FunctionCall: FunctionName '(' FunctionArgs? ')';
FunctionArgs: Expression FuncArgsRest*;
FuncArgsRest: ',' Expression;
NewStruct: StructName '{}';
ForLoop: '\nfor' LoopVar ':=' 'range' Variable Block;
