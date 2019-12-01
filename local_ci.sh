#!/usr/bin/env bash
set -x

export SYSL_PLANTUML=http://www.plantuml.com/plantuml

GOTEST_FLAGS='-race'

run_tests() {
if [[ -n `which gotestsum` ]]; then
    go test -json ./... ${GOTEST_FLAGS} | gotestsum
else
    go test ./... ${GOTEST_FLAGS}
fi
}

golangci-lint run &&
go generate ./pkg/... && run_tests &&
./scripts/check-coverage.sh 80 #&&
#go build -o tmp/sysl ./pkg/sysl  &&
#find ./pkg/ -iname "*.sysl" -exec tmp/sysl validate '{}' \;
