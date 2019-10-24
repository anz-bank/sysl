#!/usr/bin/env bash


export SYSL_PLANTUML=http://www.plantuml.com/plantuml


run_tests() {
if [[ -n `which gotestsum` ]]; then
    go test -json ./... | gotestsum
else
    go test ./...
fi
}

golangci-lint run &&
go generate ./sysl2/... && run_tests &&
./scripts/check-coverage.sh 80 #&&
#go build -o tmp/sysl ./sysl2/sysl  &&
#find ./sysl2/ -iname "*.sysl" -exec tmp/sysl validate '{}' \;
