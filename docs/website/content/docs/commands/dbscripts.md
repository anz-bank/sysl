---
title: "generate-db-scripts(beta)"
description: "Generates database scripts"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generates database scripts"
toc: true
---

The `sysl generate-db-scripts` command is used to generate database scripts.

## Usage

```
usage: sysl protobuf [<flags>] <MODULE>
```

## Required Flags
*  `-a, --app-names=APP-NAMES`    application names to parse

## Optional Flags
Optional flags:

*  `    --help      `             Show context-sensitive help (also try --help-long and --help-man).
*  `    --version   `             Show application version.
*  `    --log="warn"`             log level: [off,debug,info,warn,trace]
*  `-v, --verbose   `             enable verbose logging
*  `    --root=ROOT `             sysl root directory for input model file. If root is not found, the module directory
                               becomes the root, but the module can not import with absolute paths (or imports must be
                               relative).
*  `-t, --title=TITLE          `  file title
*  `-o, --output-dir=OUTPUT-DIR`  output directory for generated file
*  `-d, --db-type=DB-TYPE      `  database type e.g postgres

## Arguments

Args:
*  `<MODULE>`  Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.
