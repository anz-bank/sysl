PKGS := $(shell go list ./... | grep -v /vendor)

BIN_DIR := $(GOPATH)/bin

# Try to detect current branch if not provided from environment
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

# Commit hash from git
COMMIT=$(shell git rev-parse --short HEAD)

# Tag on this commit
TAG = $(shell git tag --points-at HEAD)


ifneq ("$(shell which gotestsum)", "")
	TESTEXE := gotestsum --
else
	TESTEXE := go test ./...
endif

BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
VERSION := $(or $(TAG),$(COMMIT)-$(BRANCH)-$(BUILD_DATE))

LDFLAGS = -X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)

all: test lint build coverage ## test, lint, build, coverage test run


.PHONY: all deps install grammar antlr build lint test coverage clean
lint: ## Run golangci-lint
	golangci-lint run

coverage: ## Verify the test coverage remains high
	./scripts/check-coverage.sh 80

test: ## Run tests without coverage
	rm -rf pkg/grammar/temp
	$(TESTEXE)

BINARY := sysl
PLATFORMS := windows linux darwin
.PHONY: $(PLATFORMS)
$(PLATFORMS): build
	mkdir -p release
	GOOS=$@ GOARCH=amd64 \
		go build -o release/$(BINARY)-$(VERSION)-$@$(shell test $@ = windows && echo .exe) \
		-ldflags="$(LDFLAGS)" \
		-v \
		./cmd/sysl

build: ## Build sysl into the ./dist folder
	go build -o ./dist/sysl -ldflags="$(LDFLAGS)" -v ./cmd/sysl

deps: ## Download the project dependencies with `go get`
	go get -v -t -d ./...

.PHONY: release
release: $(PLATFORMS) ## Build release binaries for all supported platforms into ./release

install: build ## Install the sysl binary into $(GOPATH)/bin
	cp ./dist/sysl $(GOPATH)/bin

clean: ## Clean temp and build files
	rm -rf release dist pkg/grammar/temp

# Autogen rules
ANTLR = java -jar pkg/antlr-4.7-complete.jar
GRAMMARS = pkg/grammar/SyslParser.g4 \
		   pkg/grammar/SyslLexer.g4

antlr: $(GRAMMARS)
	$(ANTLR) -Dlanguage=Go -lib pkg/grammar -o pkg/grammar/temp $^

pkg/grammar/sysl_parser.go: antlr
	cp pkg/grammar/temp/pkg/grammar/sysl*parser*.go pkg/grammar
	git apply pkg/grammar/antlr4-datarace-fix-parser.go.diff

pkg/grammar/sysl_lexer.go: antlr
	cp pkg/grammar/temp/pkg/grammar/sysl*lexer*.go pkg/grammar
	git apply pkg/grammar/antlr4-datarace-fix-lexer.go.diff


grammar: pkg/grammar/sysl_lexer.go pkg/grammar/sysl_parser.go pkg/parser/grammar.pb.go ## Regenerate the grammars

pkg/parser/grammar.pb.go: pkg/parser/grammar.proto
	protoc -I pkg/parser -I $(GOPATH)/src --go_out=pkg/parser grammar.proto

pkg/proto_old/sysl.pb.go: pkg/proto_old/sysl.proto
	protoc -I pkg/proto_old -I $(GOPATH)/src --go_out=pkg/proto_old/ sysl.proto

proto: pkg/sysl/sysl.pb.go

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
