ImportDecl: 'import' '(\n' ImportSpec* '\n)\n';
ImportSpec: (Import | NamedImport) '\n';
ImportSpec2: Import | NamedImport;
ImportSpec3: (Import | NamedImport);

ImportDecl2: 'import' '(\n' ImportSpec* '\n)\n';
