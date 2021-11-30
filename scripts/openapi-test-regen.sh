#!/bin/bash
# set -e

# This file is used to regenerate all test data related to OpenAPI importer changes.
# Run this in the root of the repository.

make -B internal/bundles/assets/import_openapi_cli.arraiz

FILES=($(find ./pkg/importer/tests/openapi3 -depth 1 -type f -name "*.yaml"))
for i in "${FILES[@]}"
do
    ext=".yaml"
    OUT=${i%"$ext"}
    OUT="$OUT.sysl"
    echo go run ./cmd/sysl/... import --input "$i" --format openapi3 --app-name TestApp -o "$OUT" --package=com.example.package
    go run ./cmd/sysl/... import --input "$i" --format openapi3 --app-name TestApp -o "$OUT" --package=com.example.package
done

go run ./cmd/sysl/... pb --mode textpb --root pkg/parse -o pkg/parse/tests/openapi3.sysl.golden.textpb tests/openapi3.sysl
go run ./cmd/sysl/... pb --mode textpb --root pkg/parse -o pkg/parse/tests/inplace_tuple.sysl.golden.textpb tests/inplace_tuple.sysl
