---
title: "datamodel"
description: "Generate datamodel diagrams"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generate datamodel diagrams"
toc: true
---

The `sysl datamodel` command is generate data model diagrams defined in Sysl. For an example, refer to <https://sysl.io/docs/byexample/data-model-diagrams/>

## Usage

```
usage: sysl datamodel [<flags>] <MODULE>
```

## Optional Flags

All flags are all optional.

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [trace,off,debug,info,warn]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports must be relative).
- `--class_format="%(classname)"`
- ` Specify the format string for data diagram participants. May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))
- `-t, --title=TITLE` Diagram title
- `-p, --plantuml=PLANTUML` base url of plantuml server (default: SYSL_PLANTUML or http://localhost:8080/plantuml see http://plantuml.com/server.html#install for more info)
- `-o, --output="%(epname).png"` Output file (default: %(epname).png)
- `-j, --project=PROJECT` Project pseudo-app to render
- `-d, --direct` Process data model directly without project manner
- `-f, --filter=FILTER` Only generate diagrams whose names match a pattern

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.
