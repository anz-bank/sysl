EXE = $(GOPATH)/bin/sysl
JS = sysl_js/SyslParser.js
ANTLR = java -jar lib/antlr-4.7-complete.jar

GRAMMARS = SyslParser.g4 SyslLexer.g4

.PHONY : all
all: $(EXE) $(JS)

$(EXE): grammar/sysl_parser.go grammar/sysl_lexer.go sysl.go sysl_listener_impl.go
	go install -v -buildmode=exe .

grammar/sysl_parser.go grammar/sysl_lexer.go: grammar/SyslParser.g4 grammar/SyslLexer.g4
	$(ANTLR) -Dlanguage=Go -lib grammar $^

$(JS): sysl_js/SyslParser.g4 sysl_js/SyslLexer.g4
	$(ANTLR) -Dlanguage=JavaScript -lib sysl_js $^

# make run file=path/to/file.sysl args=-gui
.PHONY : run
run: antlr
	java -cp "./classes:./lib/antlr-4.7-complete.jar" org.antlr.v4.gui.TestRig Sysl sysl_file $(file) $(args)
