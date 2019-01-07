#!/bin/bash

set -e -x

flake8
pytest
pytest test/e2e --syslexe=sysl --reljamexe=reljam
gradle test -b test/java/build.gradle
go test -coverprofile=coverage.txt -covermode=atomic github.com/anz-bank/sysl/sysl2/sysl
npm test --prefix sysl2/sysl/sysl_js
GOOS=darwin GOARCH=amd64 go build -o gosysl/gosysl-darwin github.com/anz-bank/sysl/sysl2/sysl
GOOS=linux GOARCH=amd64 go build -o gosysl/gosysl-linux github.com/anz-bank/sysl/sysl2/sysl
