---
id: cmd-codegen
title: Code Generation
sidebar_label: codegen
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

The `sysl codegen` command is used to generate the code specified by your sysl files. Currently Go is the only language fully supported by the codegen command, with the intent to expand this to other languages in the future (Swift, Kotlin, Java, JavaScript).

## Usage

```bash
usage: sysl codegen --transform=TRANSFORM --grammar=GRAMMAR [<flags>] <MODULE>
```

## Required Flags

- `--transform=TRANSFORM` path to transform file from the root transform directory
- `--grammar=GRAMMAR` path to grammar file

## Optional Flags

The remaining flags are all optional

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [off,debug,info,warn,trace]
- `-verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports must be relative).
- `--root-transform="."` sysl root directory for input transform file (default: .)
- `--app-name=""` name of the sysl app defined in sysl model. if there are multiple apps defined in sysl model, code will be generated only for the given app
- `--start="."` start rule for the grammar
- `--outdir="."` output directory
- `--dep-path=""` path passed to sysl transform
- `--basepath=""` base path for ReST output
- `--validate-only` Only Perform validation on the transform grammar
- `--disable-validator` Disable validation on the transform grammar
- `--debugger` Enable the evaluation debugger on error

## Arguments

Args:
`<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.
