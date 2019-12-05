textOnly: 'hello';

foo: 'foo';
simpleForward: textOnly;
simpleRequired: textOnly 'something' foo;

simpleChoice: textOnly | foo;

optional: foo?;
atLeastOne: foo+;
Astrix: foo*;

HardChoice: (atLeastOne | simpleForward)* | ForwardDeclTest;

ForwardDeclTest: UnnamedRule;


UnnamedRule: 'hello';


a: b+ | atLeastOne;
