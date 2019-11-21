textOnly: 'hello';

foo: 'foo';
simpleForward: textOnly;
simpleChoice: textOnly 'something' foo;

optional: foo?;
atLeastOne: foo+;
Asterix: foo*;

HardChoice: (atLeastOne | simpleForward)* | ForwardDeclTest;

ForwardDeclTest: UnnamedRule;


UnnamedRule: 'hello';
