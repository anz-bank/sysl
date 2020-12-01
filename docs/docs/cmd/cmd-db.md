---
id: cmd-db
title: Database Script (beta)
sidebar_label: generate-db-scripts
keywords:
  - command
---

:::caution
WIP

**TODO:**

- Update and polish content.
- Move referenced assets to a permanent directory on GitHub and update links.

:::

---

The `sysl generate-db-scripts` command is used to generate database scripts.

## Usage

```bash
usage: sysl protobuf [<flags>] <MODULE>
```

## Required Flags

- `-a, --app-names=APP-NAMES` application names to parse

## Optional Flags

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [off,debug,info,warn,trace]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory
  becomes the root, but the module can not import with absolute paths (or imports must be
  relative).
- `-t, --title=TITLE` file title
- `-o, --output-dir=OUTPUT-DIR` output directory for generated file
- `-d, --db-type=DB-TYPE` database type e.g postgres

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.
