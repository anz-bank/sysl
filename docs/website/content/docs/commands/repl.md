---
title: "repl"
description: "Enters a sysl REPL"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Enters a sysl REPL"
toc: true
---

The `sysl repl` command creates an interactive session where Sysl syntax can be evaluated.

## Usage

```
usage: sysl repl
```

## Optional Flags

All flags are all optional.

Optional flags:

*  `    --help      `  Show context-sensitive help (also try --help-long and --help-man).
*  `    --version   `  Show application version.
*  `    --log="warn"`  log level: [debug,info,warn,trace,off]
*  `-v, --verbose   `  enable verbose logging
*  `    --root=ROOT `  sysl root directory for input model file. If root is not found, the module directory becomes the
                    root, but the module can not import with absolute paths (or imports must be relative).