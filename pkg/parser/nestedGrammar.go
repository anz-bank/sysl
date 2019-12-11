package parser

type parserStack []*parser

func (s *parserStack) push(p *parser) {
	*s = append(*s, p)
}

func (s *parserStack) pop() {
	l := len(*s)
	*s = (*s)[:l-1]
}

func (s *parserStack) top() (p *parser) {
	l := len(*s)
	return (*s)[l-1]
}

type nestedParser struct {
	stack    parserStack
	grammars map[string]*parser
	nested   *parser
}

func makeNestedGrammarParser(text string, grammars ...*Grammar) *nestedParser {
	nested := nestedParser{
		stack:    parserStack{},
		grammars: map[string]*parser{},
		nested:   makeParser(makeNestedGrammar(), text),
	}

	for _, g := range grammars {
		p := makeParser(g, text)
		nested.grammars[p.g.Name] = p
	}
	return &nested
}

func (n *nestedParser) getLexer() *lexer {
	return n.nested.l
}

func (n *nestedParser) pushGrammar(name string) {
	name = name[1 : len(name)-1]
	n.stack.push(n.grammars[name])
}

func (n *nestedParser) popGrammar() {
	n.stack.pop()
}
