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

## Required inputs

- `<MODULE>` input file without `.sysl` extension and with leading `/`. For example `/project_dir/my_models/model`. Combine with `--root` if needed.

## Optional flags

- `-a, --app-name=APP-NAME` Name of the sysl App defined in the sysl model. If there are multiple Apps defined in the sysl model, exported file will be generated only for the given app.
- `-f, --format="swagger"` Format of export, supported options: (swagger | openapi3).
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
