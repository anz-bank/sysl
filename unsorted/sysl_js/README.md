# SyslJS

Javascript Library to parse Sysl language. Sysl is hosted at [https://github.com/anz-bank/sysl].

```
npm install sysljs
```

## Usage

```javascript
const listener = new SyslParserErrorListener();
const text = "App:\n\t...\n";
sysl.SyslParse(text, listener);
expect(listener.hasErrors == false);
```

## Antlr generation commands

```
java -cp ../../pkg/antlr-4.7-complete.jar org.antlr.v4.Tool -Dlanguage=JavaScript  SyslLexer.g4
```

```
java -cp ../../pkg/antlr-4.7-complete.jar org.antlr.v4.Tool -Dlanguage=JavaScript  SyslParser.g4
```