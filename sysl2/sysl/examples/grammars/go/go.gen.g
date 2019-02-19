goFile: PackageClause  ImportDecl?  TopLevelDecl+;
PackageClause: 'package' PackageName ';\n';

ImportDecl: 'import' '(\n' ImportSpec* '\n)\n';
ImportSpec: Import '\n';
TopLevelDecl: Declaration | FunctionDecl | MethodDecl;
Declaration: VarDecl | ConstDecl | StructType | InterfaceType ;
StructType : Comment? '\n' 'type' StructName 'struct' '{\n' FieldDecl* '}\n';
FieldDecl: '\t' identifier Type? Tag? '\n';
IdentifierList: identifier IdentifierListC*;
IdentifierListC: ',' identifier;

VarDecl: 'var' IdentifierList TypeName;
ConstDecl: 'const' '(\n'  ConstSpec '\n)\n';
ConstSpec: VarName TypeName '=' ConstValue '\n';

FunctionDecl   : 'func' FunctionName Signature? Block? ;
Signature: Parameters Result?;
Parameters: '(' ParameterList* ')';
Result         : Parameters | Type;
ParameterList     : ParameterDecl ParameterDeclC*;
ParameterDecl  : Identifier TypeName;
ParameterDeclC: ',' ParameterDecl;

InterfaceType      : 'type' InterfaceName 'interface'  '{\n'  MethodSpec* '}\n' MethodDecl*;
MethodSpec         : '\t' MethodName Signature? '\n' | InterfaceTypeName ;
MethodDecl: 'func' Receiver FunctionName Signature? Block? ;
Receiver: '(' ReceiverType ')';

Block: '{' StatementList* '}\n';
StatementList: Statement ';';
