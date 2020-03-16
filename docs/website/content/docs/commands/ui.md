---
title: "ui(beta)"
description: "Generate code"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generate code"
toc: true
---

The `sysl ui` command is used to start a webserver which displays a visual view of apps defined in Sysl.

## Usage

```bash
usage: sysl ui [<flags>] <MODULE>
```

## Optional Flags

All flags are all optional

Optional flags:

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--version` Show application version.
- `--log="warn"` log level: [off,debug,info,warn,trace]
- `-v, --verbose` enable verbose logging
- `--root=ROOT` sysl root directory for input model file. If root is not found, the module directory becomes the
-                       root, but the module can not import with absolute paths (or imports must be relative).
- `-h, --host=":8080"` host and port to serve on
- `-f, --fields="\nteam,\nteam.slack,\nowner.name,\nowner.email,\nfile.version,\nrelease.version,\ndescription,\ndeploy.env1.url,\ndeploy.sit1.url,\ndeploy.sit2.url,\ndeploy.* qa.url,\ndeploy.prod.url,\nrepo.url,\ndocs.url,\ntype"`
-                       fields to display on the UI, separated by comma
- `--grpcui` enables the grpcUI handlers

## Arguments

Args:
<MODULE> Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.
