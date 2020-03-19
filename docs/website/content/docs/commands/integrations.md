---
title: "integrations"
description: "Generate integration diagrams"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generate integration diagrams"
toc: true
---

The `sysl integrations` command generates integration diagrams. For an example, refer to <https://sysl.io/docs/byexample/integration-diagrams/>

## Usage

```bash
usage: sysl integrations [<flags>] <MODULE>
```

Aliases

- ints `usage: sysl ints [<flags>] <MODULE>`

## Optional Flags

All flags are all optional

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [warn,trace,off,debug,info]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory
- `` becomes the root, but the module can not import with absolute paths (or imports must
- `` be relative).
- `-t, --title=TITLE` diagram title
- `-p, --plantuml=PLANTUML` base url of plantuml server (default: SYSL_PLANTUML or http://localhost:8080/plantuml
- `` see http://plantuml.com/server.html#install for more info)
- `-o, --output="%(epname).png"` output file(default: %(epname).png)
- `-j, --project=PROJECT` project pseudo-app to render
- `--filter=FILTER` Only generate diagrams whose output paths match a pattern
- `-e, --exclude=EXCLUDE ...` apps to exclude
- `-c, --clustered` group integration components into clusters
- `--epa` produce and EPA integration view

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.
