EXE = $(GOPATH)/bin/sysl
SD_EXE = $(GOPATH)/bin/sd
JS = sysl_js/SyslParser.js
ANTLR = java -jar lib/antlr-4.7-complete.jar
SYSL_PYTHON_BIN = /usr/local/bin/sysl

GRAMMARS = SyslParser.g4 SyslLexer.g4

.PHONY : all
all: tests $(EXE) $(JS)

CODE_FILES = \
	ValueUtil.go \
	codegen.go \
	constants.go \
	error_listener.go \
	exit.go \
	exprEval.go \
	exprOp.go \
	fs.go \
	helpers.go \
	listener_impl.go \
	output.go \
	parse.go \
	sysl.go \
	testtemp.go \
	binexprEval.go \
	valueType.go \
	unaryEval.go \
	goFuncs.go \
	grammar/sysl_parser.go \
	grammar/sysl_lexer.go

$(EXE): $(CODE_FILES)
	go install -v

grammar/sysl_parser.go grammar/sysl_lexer.go: grammar/SyslParser.g4 grammar/SyslLexer.g4
	$(ANTLR) -Dlanguage=Go -lib grammar $^

$(JS): sysl_js/SyslParser.g4 sysl_js/SyslLexer.g4
	$(ANTLR) -Dlanguage=JavaScript -lib sysl_js $^

# make run file=path/to/file.sysl args=-gui
.PHONY : run
run: antlr
	java -cp "./classes:./lib/antlr-4.7-complete.jar" org.antlr.v4.gui.TestRig Sysl sysl_file $(file) $(args)

SYSL_TEST_FILES = \
	attribs \
	bad_order \
	collector \
	crash \
	docstrings \
	duplicate \
	ep_params \
	eventing \
	for_loop \
	funcs \
	group_stmt \
	if_else \
	implied \
	inplace_tuple \
	library \
	matching \
	math \
	mixin \
	navigate \
	oneof \
	project \
	pubsub_collector \
	rank \
	rest_api_query_params \
	school \
	stmts \
	tableof \
	test1 \
	test2 \
	test4 \
	test_rest_api \
	transform \
	until_loop \
	with_spaces

tests: $(patsubst %,tests/%.sysl.golden.textpb,$(SYSL_TEST_FILES))

tests/school.sysl.golden.pb: root = tests

tests/%.sysl.golden.textpb: tests/%.sysl $(SYSL_PYTHON_BIN)
	$(SYSL_PYTHON_BIN) $(if $(root),--root $(root),) textpb -o $@ $<

seqs:
	go build -o $(SD_EXE)
