---
id: common-flags
title: Common Flags
sidebar_label: Common Flags
keywords:
  - command
---

## Summary

This page lists all common optional flags which can be used in Sysl commands.

## Common Optional Flags

- `--help` Show context-sensitive help (also try --help-long).
- `--version` Show application version.
- `--log="warn"` Log level: [off,debug,info,warn,trace].
- `-v, --verbose` Enable verbose logging.
- `--root=ROOT` Sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports must be relative).

## Examples

```
sysl --help
```

```
sysl --version
```

```
sysl import -v --input=simple-api.yaml --app-name=Simple --output=simple-api.sysl
```

```
sysl import --log="info" --input=simple-api.yaml --app-name=Simple --output=simple-api.sysl
```

```
sysl export --format="openapi3" --output="%(appname).yaml" --root=./demos simple-openapi3.sysl
```
