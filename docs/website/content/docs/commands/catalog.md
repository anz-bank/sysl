---
title: "catalog"
description: "Generate documentation"
date: 2020-05-06T14:05:40+11:00
weight: 70
draft: false
bref: "Generate documentation"
toc: true
---

The `sysl catalog` command is used to generate documentation for systems specified in your sysl files.
Passing in the `-s` flag starts up a webserver to allow viewing of the documentation in the browser.
The generated output can be of two formats: markdown or html.

## Usage

```bash
usage: sysl catalog [<flags>] <MODULE>...
```

## Optional Flags

All flags are all optional

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [off,debug,info,warn,trace]
- `-verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports must be relative).
- `-p, --port=":6900"       `  host and port to serve on
- `-s, --server             `  start a server on port
- `-t, --outputType=markdown`  output type (markdown | html)
- `-o, --outputDir="/"      `  output directory to generate docs to

## Arguments

Args:
<MODULE> Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.

usage: sysl catalog [<flags>] <MODULE>...

Generate Documentation from your sysl definitions

Optional flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
      --version              Show application version.
      --log="warn"           log level: [off,debug,info,warn,trace]
  -v, --verbose              enable verbose logging
      --root=ROOT            sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports must be relative).

