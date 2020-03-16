---
title: "import"
description: "Imports sysl from other formats"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Imports sysl from other formats"
toc: true
---

The `sysl import` command imports foreign types to sysl.
Supported types: [openapi, swagger, xsd]

Currently, the supported formats include:

- OpenAPI 3.0 `openapi`
- OpenAPI 2.0 `swagger`
- XSD `xsd`

Note: The grammar importer type is deprecated.

## Usage

```bash
usage: sysl import --input=INPUT --app-name=APP-NAME [<flags>]
```

## Required Flags

Required flags:
-i, --input=INPUT input filename
-a, --app-name=APP-NAME name of the sysl app to define in sysl model.

## Optional Flags

The remaining flags are all optional

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [info,warn,trace,off,debug]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory
  becomes the root, but the module can not import with absolute paths (or imports must be
  relative).
- `-p, --package=PACKAGE` name of the sysl package to define in sysl model.
- `-o, --output="output.sysl"` output filename
- `-f, --format=auto` format of the input filename, options: [auto, grammar, openapi, swagger, xsd]
