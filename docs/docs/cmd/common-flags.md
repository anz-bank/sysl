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

## Standard Input

In some contexts, Sysl source is not available to be read from disk and must be passed to Sysl via stdin. This can be done like so:

```sh
echo '[{"path": "path/to/foo.sysl", "content": "$(cat path/to/foo.sysl)"}]' | sysl cmd`
```

The format of the stdin data is JSON-encoded array of files, where each file has `path` and `content` properties.

- `content` is the source to parse
- `path` is where the source should be assumed to live (whether or not is actually on disk). The parser will use this when resolving relative imports in the source.

Note that modules cannot be provided via both stdin and the `MODULES` arg. If both are present, `MODULES` will take precedence and stdin will be ignored.

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
