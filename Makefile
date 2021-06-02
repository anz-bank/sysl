include ./scripts/version-report.mk

.PHONY: all install grammar antlr build lint test test-arrai coverage clean check-clean golden

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
		golangci/golangci-lint:v1.30.0 \
			golangci-lint run -v

tidy:
	go mod tidy
	gofmt -s -w .
	goimports -w .
	make -C docs tidy

# Generates intermediate files for build.
generate: internal/arrai/bindata.go
gen: generate

test: test-arrai coverage

test-arrai:
	$(ARRAI) test

coverage: generate
	./scripts/test-with-coverage.sh 80

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

buildlsp:
	go build -o ./dist/sysllsp -ldflags=$(LDFLAGS) -v ./cmd/sysllsp

build-docker: generate
	docker build -t sysl .

build-sysl-version-diff-docker: generate
	docker build -t sysl-version-diff -f sysl-version-diff/Dockerfile .

# Assumes that every arr.ai script depends on every other arr.ai script.
%.arraiz: %.arrai $(shell find . -name '*.arrai')
	$(ARRAI) bundle $< > $@

internal/arrai/bindata.go: \
		pkg/importer/avro/transformer_cli.arraiz \
		pkg/importer/sql/import_cli.arraiz
	# Binary files in bindata.go have metadata like size, mode and modification time(modTime).
	# And modTime will be updated every time when arrai bundle file is regenerated, it will cause task check-clean failed.
	# So add parameter `-modtime 1` to set a fixed modTime.
	# Add `-mode 0644` for similar reason as files' mode are possible different in CI and local development environments.
	go-bindata -mode 0644 -modtime 1 -pkg arrai -o $@ $^
	gofmt -s -w $@

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
