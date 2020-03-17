# Grammar musings

Thinking about language embeddings: Sysl's own grammar could be defined as a Sysl grammar. This implies that the only grammar requiring low-level coding is the grammar syntax itself, which would be relatively simple:

- terminals (string literals `"Hi!"` and regexes `/REGEXP/`)¹
- non-terminal rules ²
- parenthetical groups ³
- union: `|` ⁴
- exclusion: `-` ⁵
- quantifiers: `?` `*` `+`
- lists, e.g.: `Arg:","` (shorthand for `Arg ("," Arg)*`) ⁶
- An escaping mechanism to jump in and out of embedded languages, e.g.: `!{java: ... :}`
- To declare a grammar `!grammar`
- Special tokens: `INDENT`, `OUTDENT` maybe `WHITESPACE` ⁷

There should be some kind of built-in support for implicitly skipping whitespace (an almost universal feature of programming languages (with the notable exceptions of XML and whitespace)) and indent/outdent to ease support for Python, YAML and Sysl itself.

The output would be a grammar protobuf, which is handed to a grammar interpreter to parse files written in the given language. Sysl could even parse multiple top-level formats like JSON and YAML by mapping from the file extension to a standard grammar expressed in Sysl, or with a command-line option: `--lang=json` or `--lang=lib/lang/java.sysl`.

```
¹ 	NUMBER:
		/[1-9][0-9]*/

² 	EXPRESSION:
		NUMBER OPERATOR NUMBER
		NUMBER OPERATOR EXPRESSION

³ 	(KEYWORD | BOOL_LITERAL | NULL_LITERAL)*

⁴ 	"[" (INT | "-")* "]"   # valid array type: [123 123 - 234]

⁵ 	IDENTIFIER:
		IDENTIFIER_CHARS - (KEYWORD | BOOL_LITERAL | NULL_LITERAL)

⁶ 	"[" (INT | "-"):","? "]" # valid array type: [123, 123, -, 234] []

⁷ 	PYTHON_FUNC: "def" ID ":" INDENT BODY DEDENT
```

## Example: Java Grammar

See [pkg/lib/lang/java.sysl](../../pkg/lib/lang/java.sysl).
