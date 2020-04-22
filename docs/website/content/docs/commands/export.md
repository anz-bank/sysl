---
title: "export"
description: "Exports sysl to other formats"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Exports sysl to other formats"
toc: true
---

The `sysl export` command is convert sysl files to other formats. All flags are optional if only a single App is defined in the input sysl file.
Currently, the only supported format is OpenAPI 2.0 (formerly Swagger 2.0).

Note: Types named `Empty` in Sysl are treated specially and are not exported.

## Usage

```bash
usage: sysl export [<flags>] <MODULE>
```

## Optional Flags

All other flags are all optional

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [debug,info,warn,trace,off]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory
  becomes the root, but the module can not import with absolute paths (or imports must be
  relative).
- `-f, --format="swagger"` format of export, supported options; swagger
- `-o, --output="%(appname).yaml"`
  output filepath.format(yaml | json) (default: %(appname).yaml)
- `-a, --app-name=APP-NAME` name of the sysl app defined in sysl model. if there are multiple apps defined in sysl
  model, swagger will be generated only for the given app

## Arguments

Args:
<MODULE> Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.
