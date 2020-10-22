include ./scripts/version-report.mk

.PHONY: all install grammar antlr build lint coverage clean check-tidy golden

GOPATH		= $(shell go env GOPATH)
GOVERSION	= $(shell go version | cut -d' ' -f3-4)
BIN_DIR		:= $(GOPATH)/bin
ARRAI		= docker run --rm -v $(CURDIR):/src -w /src anzbank/arrai

ifneq ("$(shell which gotestsum)", "")
	TESTEXE := gotestsum --
else
	TESTEXE := go test ./...
endif

all: lint build buildlsp install test coverage

lint: generate
	golangci-lint run ./...

lint-docker: generate
	docker run --rm \
		-v $(CURDIR):/app \
		-v $(CURDIR)/.lint-cache:/cache/go \
		-e GOCACHE=/cache/go \
		-e GOLANGCI_LINT_CACHE=/cache/go \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /app \
		golangci/golangci-lint:v1.30.0 \
			golangci-lint run -v

tidy:
	go mod tidy
	gofmt -s -w .
	goimports -w .

# Generates intermediate files for build.
generate: internal/arrai/bindata.go

test: generate
	$(TESTEXE)

coverage: generate
	./scripts/test-with-coverage.sh 80

# Updates golden test files in pkg/parse.
# TODO: Extend to work for all golden files
golden: generate
	go test ./pkg/parse ./pkg/exporter ./pkg/importer -update

check-tidy: generate
	git --no-pager diff HEAD && test -z "$$(git status --porcelain)"

build: generate
	go build -o ./dist/sysl -ldflags=$(LDFLAGS) -v ./cmd/sysl

build-windows: generate
	go build -o ./dist/sysl.exe -ldflags=$(LDFLAGS) -v ./cmd/sysl

buildlsp:
	go build -o ./dist/sysllsp -ldflags=$(LDFLAGS) -v ./cmd/sysllsp

build-docker: generate
	docker build -t sysl .

# Assumes that every arr.ai script depends on every other arr.ai script.
%.arraiz: %.arrai $(shell find . -name '*.arrai')
	$(ARRAI) bundle $< > $@

internal/arrai/bindata.go: \
		pkg/importer/avro/transformer_cli.arraiz \
		pkg/importer/spanner/import_spanner_sql.arraiz \
		pkg/importer/spanner/import_migrations.arraiz \
		pkg/exporter/spanner/spanner_cli.arraiz
	# Binary files in bindata.go have metadata like size, mode and modification time(modTime).
	# And modTime will be updated every time when arrai bundle file is regenerated, it will cause task check-tidy failed.
	# So add parameter `-modtime 1` to set a fixed modTime.
	# Add `-mode 0644` for similar reason as files' mode are possible different in CI and local development environments.
	go-bindata -mode 0644 -modtime 1 -pkg arrai -o $@ $^

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
