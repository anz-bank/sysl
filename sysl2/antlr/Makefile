EXE=$(GOPATH)/bin/antlr
JS=js/SyslParser.js

.PHONY : all
all: $(EXE) $(JS)

$(EXE): grammar/sysl_parser.go grammar/sysl_lexer.go sysl.go sysl_listener_impl.go
	go install -v -buildmode=exe .

grammar/sysl_parser.go: grammar/SyslParser.g4 grammar/SyslLexer.g4
	@echo $(EXE)
	java -jar lib/antlr-4.7-complete.jar -Dlanguage=Go -lib grammar grammar/SyslParser.g4 grammar/SyslLexer.g4

js/SyslParser.js: js/SyslParser.g4 js/SyslLexer.g4
	@echo $(JS)
	java -jar lib/antlr-4.7-complete.jar -Dlanguage=JavaScript -lib js js/SyslLexer.g4 js/SyslParser.g4

# make run file=path/to/file.sysl args=-gui
.PHONY : run
run: antlr
	java -cp "./classes:./lib/antlr-4.7-complete.jar" org.antlr.v4.gui.TestRig Sysl sysl_file $(file) $(args)
