include ./scripts/version-report.mk

.PHONY: all install grammar antlr build lint test test-arrai coverage clean check-clean gen generate golden

GOPATH		= $(shell go env GOPATH)
GOVERSION	= $(shell go version | cut -d' ' -f3-4)
BIN_DIR		:= $(GOPATH)/bin
ARRAI		= go run github.com/arr-ai/arrai/cmd/arrai

ifneq ("$(shell which gotestsum)", "")
	TESTEXE := gotestsum --
else
	TESTEXE := go test ./...
endif

all: lint build buildlsp install test coverage

lint: generate
	golangci-lint run ./...
	make -C docs lint

lint-docker: generate
	docker run --rm \
		-v $(CURDIR):/app \
		-e GOCACHE=/cache/go \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /app \
		golangci/golangci-lint:v1.41.1 \
			golangci-lint run -v

tidy:
	go mod tidy
	gofmt -s -w .
	goimports -w .
	make -C docs tidy

# Generates intermediate files for build.
generate: internal/bundles/bundles.go plugins bundled-proto
	go generate ./pkg/lsp/...

gen: generate

test: test-arrai coverage

test-arrai:
	$(ARRAI) test

coverage: generate
	./scripts/test-with-coverage.sh 75

# Updates golden test files in pkg/parse.
golden: generate
	go test ./pkg/parse ./pkg/exporter ./pkg/importer -update
	# TODO: Extend the -update flag to work for all golden files
	sysl pb --mode=textpb tests/args.sysl > tests/args.sysl.golden.textpb
	sysl pb --mode=json tests/args.sysl > tests/args.sysl.golden.json
	sysl pb --mode=textpb tests/type_merge1.sysl > tests/type_merge1.sysl.golden.textpb
	sysl pb --mode=textpb tests/type_merge2.sysl > tests/type_merge2.sysl.golden.textpb
	sysl pb --mode=textpb tests/file_merge.sysl > tests/file_merge.sysl.golden.textpb
	sysl pb --mode=json tests/file_merge.sysl > tests/file_merge.sysl.golden.json
	sysl pb --mode=textpb tests/namespace_merge.sysl > tests/namespace_merge.sysl.golden.textpb

check-clean: generate
	git --no-pager diff HEAD && test -z "$$(git status --porcelain)"

build: generate
	go build -o ./dist/sysl -ldflags=$(LDFLAGS) -v ./cmd/sysl

build-windows: generate
	go build -o ./dist/sysl.exe -ldflags=$(LDFLAGS) -v ./cmd/sysl

buildlsp: generate
	go build -o ./dist/sysllsp -ldflags=$(LDFLAGS) -v ./cmd/sysllsp

build-docker: generate
	docker build -t sysl .

build-sysl-version-diff-docker: generate
	docker build -t sysl-version-diff -f sysl-version-diff/Dockerfile .

# Assumes that every arr.ai script depends on every other arr.ai script.
%.arraiz: %.arrai $(shell find . -name '*.arrai')
	$(ARRAI) bundle $< > $@

.PHONY: plugins
plugins: \
		pkg/plugins/integration_model_plugin.arraiz

internal/bundles/assets/transformer_cli.arraiz: pkg/importer/avro/transformer_cli.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/assets/import_sql_cli.arraiz: pkg/importer/sql/import_sql_cli.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/assets/import_openapi_cli.arraiz: pkg/importer/openapi/import_openapi_cli.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/assets/spanner_cli.arraiz: pkg/exporter/spanner/spanner_cli.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/assets/import_proto_cli.arraiz: pkg/importer/proto/import_proto_cli.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/exporters/%/transform.arraiz: transforms/exporters/%/transform.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/importers/%/transform.arraiz: transforms/importers/%/transform.arrai
	$(ARRAI) bundle $< > $@

internal/bundles/bundles.go: \
		internal/bundles/assets/transformer_cli.arraiz \
		internal/bundles/assets/import_sql_cli.arraiz \
		internal/bundles/assets/import_openapi_cli.arraiz \
		internal/bundles/assets/spanner_cli.arraiz \
		internal/bundles/assets/import_proto_cli.arraiz \
		internal/bundles/exporters/proto/transform.arraiz \
		internal/bundles/importers/jsonschema/transform.arraiz

pkg/importer/proto/bundled_files/local_imports.arrai: pkg/importer/proto/bundled_files/bundler.arrai
	$(ARRAI) run $< > pkg/importer/proto/bundled_files/tmp.arrai && \
	mv -f pkg/importer/proto/bundled_files/tmp.arrai $@

.PHONY: bundled-proto
bundled-proto: pkg/importer/proto/bundled_files/local_imports.arrai

release:
	GOVERSION="$(GOVERSION)" goreleaser build --rm-dist --snapshot -f .goreleaser.yml

install: build ## Install the sysl binary into $(GOPATH)/bin. We don't use go install because we need to pass in LDFLAGS.
	test -n "$(GOPATH)"  # $$GOPATH
	cp ./dist/sysl $(GOPATH)/bin

clean:
	rm -rf dist

# Autogen rules
ANTLR = java -jar pkg/antlr-4.7-complete.jar
ANTLR_GO = $(ANTLR) -Dlanguage=Go -lib $(@D) $<

grammar: pkg/grammar/sysl_lexer.go pkg/grammar/sysl_parser.go ## Regenerate the grammars

pkg/grammar/sysl_parser.go: pkg/grammar/SyslParser.g4 pkg/grammar/SyslLexer.tokens
	$(ANTLR_GO)

pkg/grammar/sysl_lexer.go pkg/grammar/SyslLexer.tokens &: pkg/grammar/SyslLexer.g4
	$(ANTLR_GO)

pkg/sysl/sysl.pb.go: pkg/sysl/sysl.proto
	protoc -I pkg/sysl -I $(GOPATH)/src --go_out=pkg/sysl/ sysl.proto

pkg/sysl/sysl.pb: pkg/sysl/sysl.proto
	protoc -o $@ $<

proto: pkg/sysl/sysl.pb.go pkg/sysl/sysl.pb

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test-grammar:
	which wbnf || go install github.com/arr-ai/wbnf
	./scripts/test-grammar-wbnf.sh . | diff ./scripts/grammar-out.txt -

update-grammar-result:
	which wbnf || go install github.com/arr-ai/wbnf
	./scripts/test-grammar-wbnf.sh . > ./scripts/grammar-out.txt
