---
title: "generate-db-scripts-delta(beta)"
description: "Generates database change scripts"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generates database change scripts"
toc: true
---

The `sysl generate-db-scripts-delta` command is used to generate database change scripts.

## Usage

Two sysl modules are required

```bash
usage: sysl generate-db-scripts-delta [<flags>] <MODULE> <MODULE>
```

## Required Flags

## Optional Flags

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [debug,info,warn,trace,off]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory
  becomes the root, but the module can not import with absolute paths (or imports must be
  relative).
- `-t, --title=TITLE` file title
- `-o, --output-dir=OUTPUT-DIR` output directory
- `-a, --app-names=APP-NAMES` application names to read
- `-d, --db-type=DB-TYPE` database type e.g postgres

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.
