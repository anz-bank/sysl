---
id: cmd-export
title: Export
sidebar_label: export
keywords:
  - command
---

## Summary

`sysl export` command converts `.sysl` files to other formats.

## Usage

```bash
sysl export [<flags>] <MODULE>
```

Currently, the supported formats include:

- OpenAPI 2.0 (fka Swagger): `swagger`
- OpenAPI 3.0 `openapi`
- Spanner `spanner`
- Proto `proto`

## Required inputs

- `<MODULE>` input file without `.sysl` extension and with leading `/`. For example `/project_dir/my_models/model`. Combine with `--root` if needed.

## Optional flags

- `-a, --app-name=APP-NAME` Name of the sysl App defined in the sysl model. If there are multiple Apps defined in the sysl model, exported file will be generated only for the given app.
- `-f, --format="swagger"` Format of export, supported options: (swagger | openapi3 | spanner | proto).
- `-o, --output="%(appname).yaml"` Output file path, supported file extensions: (yaml | json) and default value is `%(appname).yaml`. Note, appname is the one specified by flag `--app-name`.

[More common optional flags](common-flags.md)

## Arguments

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.

## Examples

### OpenAPI

Command line

```bash
sysl export --format="openapi3" --app-name="SimpleOpenAPI3" --output="simple-openapi3.yaml" simple-openapi3.sysl
```

```sysl title="Input Sysl file: simple-openapi3.sysl"
SimpleOpenAPI3 "SimpleOpenAPI3":
    @description =:
        | Simple demo for openapi file export

    /test:
        GET:
            | Endpoint for testing GET
            return error
            return ok <: SimpleObj

    #---------------------------------------------------------------------------
    # definitions

    !type SimpleObj:
        name <: string?:
            @json_tag = "name"

    !type SimpleObj2:
        name <: SimpleObj?:
            @json_tag = "name"
```

```yaml title="Output OpenAPI3 file: simple-openapi3.yaml"
components:
  schemas:
    SimpleObj:
      properties:
        name:
          type: string
      type: object
    SimpleObj2:
      properties:
        name:
          $ref: "#/components/schemas/SimpleObj"
      type: object
info:
  contact: {}
  description: |
    Simple demo for openapi file export
  title: SimpleOpenAPI3
  version: ""
openapi: 3.0.0
paths:
  /test:
    get:
      description: Endpoint for testing GET
      responses:
        "200":
          content:
            application/json: {}
          description: "200"
servers:
  - url: ""
```

### Swagger

Command line

```bash
sysl export --format="swagger" --output="simple-swagger.yaml" --app-name="SimpleSwagger" simple-swagger.sysl
```

```sysl title="Input Sysl file: simple-swagger.sysl"
SimpleSwagger "SimpleSwagger":
    @description =:
        | Simple demo for swagger file export

    /test:
        GET:
            | Endpoint for testing GET
            return 200 <: SimpleObj

        DELETE:
            | Endpoint for testing DELETE
            return 203

    /tests:
        GET:
            | Endpoint for testing GET
            return set of SimpleObj

    #---------------------------------------------------------------------------
    # definitions

    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
```

```yaml title="Output Swagger file: simple-swagger.yaml"
definitions:
  SimpleObj:
    format: tuple
    properties:
      name:
        format: string
        type: string
    type: object
info:
  description: |
    Simple demo for swagger file export
  title: SimpleSwagger
  version: 0.0.0
paths:
  /test:
    delete:
      consumes:
        - application/json
      description: Endpoint for testing DELETE
      produces:
        - application/json
      responses:
        "203":
          description: Non-Authoritative Information
      summary: Endpoint for testing DELETE
    get:
      consumes:
        - application/json
      description: Endpoint for testing GET
      produces:
        - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: "#/definitions/SimpleObj"
      summary: Endpoint for testing GET
  /tests:
    get:
      consumes:
        - application/json
      description: Endpoint for testing GET
      produces:
        - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: "#/definitions/SimpleObj"
      summary: Endpoint for testing GET
swagger: "2.0"
```

### Proto

Exporting to proto from sysl requires 2 additional environmental variables to be provided with the command.

| Variable | Description |
| -- | -- |
| SYSL_PROTO_NAMESPACE | The Sysl namespace to apply to the contents of the protobuf file. |
| SYSL_PROTO_PACKAGE | The name to use for the protobuf package |

Command line

```bash
SYSL_PROTO_NAMESPACE="Name :: Space" SYSL_PROTO_PACKAGE="Foo" sysl export --format="proto" --output="example.proto" example.sysl
```

```proto title="Output Proto file: example.proto
// {"sysl": {"namespace": "Name :: Space"}}

// File generated by Sysl. DO NOT EDIT.

syntax = "proto3";

package Foo;

option go_package = "Foo";

message Request {
  float64 decimal_with_precision = 1;
  string nativeTypeField = 2;
  int32 numbered = 42;
  string optional = 4;
  string primaryKey = 5;
  Type reference = 6;
  string sequence = 7;
  string set = 8;
  string string_max_constraint = 9;
  string string_range_constraint = 10;
  string with_anno = 11;
}

message Type {
  int32 foo = 1;
}
```