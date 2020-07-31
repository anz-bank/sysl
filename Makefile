.PHONY: all install grammar antlr build lint test coverage clean check-tidy golden

include VersionReport.mk

GOVERSION=$(shell go version | cut -d' ' -f3-4)
BIN_DIR := $(GOPATH)/bin

ifneq ("$(shell which gotestsum)", "")
	TESTEXE := gotestsum --
else
	TESTEXE := go test ./...
endif

all: lint build buildlsp install examples test coverage

TUTORIALS: $(wildcard ./demo/examples/*) $(wildcard ./demo/examples/*/*)

examples: TUTORIALS
	cd demo/examples/ && go run generate_website.go && cd ../../  && git --no-pager diff HEAD && test -z "$$(git status --porcelain)"

lint:
	golangci-lint run ./...

test:
	$(TESTEXE)

coverage:
	./scripts/test-with-coverage.sh 80

# Updates golden test files in pkg/parse.
# TODO: Extend to work for all golden files
golden:
	go test ./pkg/parse ./pkg/exporter ./pkg/importer -update

check-tidy: ## Check go.mod and go.sum is tidy
	go mod tidy && git --no-pager diff HEAD && test -z "$$(git status --porcelain)"

build:
	go build -o ./dist/sysl -ldflags=$(LDFLAGS) -v ./cmd/sysl

buildlsp:
	go build -o ./dist/sysllsp -ldflags=$(LDFLAGS) -v ./cmd/sysllsp

release:
	GOVERSION="$(GOVERSION)" goreleaser build --rm-dist --snapshot

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

proto: pkg/sysl/sysl.pb.go

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test-grammar:
	which wbnf || go install github.com/arr-ai/wbnf
	./scripts/test-grammar-wbnf.sh . | diff ./scripts/grammar-out.txt -

update-grammar-result:
	which wbnf || go install github.com/arr-ai/wbnf
	./scripts/test-grammar-wbnf.sh . > ./scripts/grammar-out.txt
