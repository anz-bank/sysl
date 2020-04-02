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

all: test lint build buildlsp coverage examples ## test, lint, build, coverage test run

TUTORIALS: $(wildcard ./demo/examples/*) $(wildcard ./demo/examples/*/*)

examples: TUTORIALS
	cd demo/examples/ && go run generate_website.go && cd ../../

.PHONY: all deps install grammar antlr build lint test coverage clean check-tidy
lint: ## Run golangci-lint
	golangci-lint run

coverage: ## Run tests and verify the test coverage remains high
	./scripts/test-with-coverage.sh 85

check-tidy: ## Check go.mod and go.sum is tidy
	go mod tidy && test -z "$$(git status --porcelain)"

test: ## Run tests without coverage
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

buildlsp: ## Build sysllsp into the ./dist folder
	go build -o ./dist/sysllsp -ldflags="$(LDFLAGS)" -v ./cmd/sysllsp

## This option is used for Sysl UI to bundle static files into the go binary.
## For more details on pkger, refer to https://github.com/markbates/pkger
.PHONY: resources
resources:
	cd ui && npm run build
	pkger
	mv pkged.go pkg/ui/pkged.go
	## Replaces the package declaration 'main' with'ui' due to bug in pkger
	## Remove once https://github.com/markbates/pkger/pull/67 has been merged in
	sed -i '' 's/main/ui/' pkg/ui/pkged.go

deps: ## Download the project dependencies with `go get`
	go mod tidy
ifneq ("$(shell git status --porcelain)", "")
	## GoReleaser has to make sure go.mod is up to date before release sysl binary.
	## Keep everyone remembering to update go.mod to avoid release failure.
	echo "git is currently in a dirty state, please check in your pipeline what can be changing the following files:$(shell git diff)"
	exit 1
endif
	go get -v -t -d ./...

.PHONY: release
release: $(PLATFORMS) ## Build release binaries for all supported platforms into ./release

install: build ## Install the sysl binary into $(GOPATH)/bin
	cp ./dist/sysl $(GOPATH)/bin

clean: ## Clean temp and build files
	rm -rf release dist

# Autogen rules
ANTLR = java -jar pkg/antlr-4.7-complete.jar
ANTLR_GO = $(ANTLR) -Dlanguage=Go -lib $(@D) $<

grammar: pkg/grammar/sysl_lexer.go pkg/grammar/sysl_parser.go ## Regenerate the grammars

pkg/grammar/sysl_parser.go: pkg/grammar/SyslParser.g4 pkg/grammar/SyslLexer.tokens
	$(ANTLR_GO)

pkg/grammar/sysl_lexer.go pkg/grammar/SyslLexer.tokens &: pkg/grammar/SyslLexer.g4
	$(ANTLR_GO)

pkg/proto_old/sysl.pb.go: pkg/proto_old/sysl.proto
	protoc -I pkg/proto_old -I $(GOPATH)/src --go_out=pkg/proto_old/ sysl.proto

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
