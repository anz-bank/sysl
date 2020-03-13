---
title: "protobuf"
description: "Generates a sequence diagram"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generates a sequence diagram"
toc: true
---

The `sysl protobuf` command is used to generate a proto representation of the Sysl file.

## Usage

```
usage: sysl protobuf [<flags>] <MODULE>
```

## Optional Flags

Optional flags:

*  `    --help       `  Show context-sensitive help (also try --help-long and --help-man).
*  `    --version    `  Show application version.
*  `    --log="warn" `  log level: [info,warn,trace,off,debug]
*  `-v, --verbose    `  enable verbose logging
*  `    --root=ROOT  `  sysl root directory for input model file. If root is not found, the module directory becomes the
*  `                 `  root, but the module can not import with absolute paths (or imports must be relative).
*  `-o, --output="-" `  output file name
*  `    --mode=textpb`  output mode: [textpb,json]

## Arguments

Args:
*  `<MODULE>`  Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.