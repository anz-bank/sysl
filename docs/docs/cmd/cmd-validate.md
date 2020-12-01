---
id: cmd-validate
title: Validate
sidebar_label: validate
keywords:
  - command
  - validate
---

## Summary

The `sysl validate` command is used to verify that sysl files are valid.

## Usage

```bash
usage: sysl validate [<flags>] <MODULE>
```

## Required Flags

## Optional Flags

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [info,warn,trace,off,debug]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory becomes the
  root, but the module can not import with absolute paths (or imports must be relative).

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.

## Linter

The validator also has some linting features that helps with writing
a good sysl file.

The format of the warnings are as follows:

1. When there is only one location:

```
lint path/to/file:lineNumber:colNumber: linter message
```

2. When there are multiple locations:

```
lint: linter message:
additional message:path/to/file1:lineNumber:colNumber
additional message:path/to/file2:lineNumber:colNumber
additional message:path/to/file3:lineNumber:colNumber
...
```

Currently, the linter checks for the following things:

### Return Statements

This linter checks for correct use of return statements in sysl.
Return statements in sysl have to use either the keyword `ok` or
`error` or HTTP status code followed by an optional type.

```sysl
...
return ok
return error
return 200 <: some_type
...
```

This linter checks for statements such as

```sysl
return some_type
```

An example of a warning is the following:

```log
lint path/to/file.sysl:3:4: 'return some_type' not supported, use
'return ok <: some_type' instead
```

### Case-Sensitive Application Redefinition

This linter checks for case-sensitive application redefinition.
For example:

```sysl
App:
    ...

aPP:
    ...

ApP:
    ...
```

An example warning of this linter is the following:

```log
lint: case-sensitive redefinitions detected:
ApP:path/to/file.sysl:7:1
App:path/to/file.sysl:1:1
aPP:path/to/file.sysl:4:1
```

### Call Statements

This linter checks for the validity of call statements. This linter specifically
checks for whether a call statements calls to a defined endpoint or not. For
example, given the following sysl:

```sysl
App:
    Endpoint:
        . <- Endpoint2
```

As `App Endpoint2` is not defined yet, the linter will show warnings about this.

For example:

```log
lint path/to/file.sysl:3:8: Endpoint 'Endpoint2' does not exist for call 'App <- Endpoint2'
```

and for REST endpoints:

```sysl
App:
    /Endpoint:
        GET:
            . <- GET Endpoint2
```

The warning is as follows:

```log
lint path/to/file.sysl:4:12: Endpoint 'Endpoint2' does not exist for call 'App <- GET Endpoint2'
```
