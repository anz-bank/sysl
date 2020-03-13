---
title: "template(beta)"
description: "Applies a model to a template"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Applies a model to a template"
toc: true
---

The `sysl template` command is used to apply a model to a template for custom text output
## Usage

`usage: sysl template --template=TEMPLATE --start=START [<flags>] <MODULE>...`

## Required Flags

Required flags:
*  `--template=TEMPLATE`  path to template file from the root transform directory
*  `--start=START      ` start rule for the template

## Optional Flags

Optional flags:

*  `    --help      `         Show context-sensitive help (also try --help-long and --help-man).
*  `    --version   `         Show application version.
*  `    --log="warn"`         log level: [info,warn,trace,off,debug]
*  `-v, --verbose   `         enable verbose logging
*  `    --root=ROOT `         sysl root directory for input model file. If root is not found, the module directory
                           becomes the root, but the module can not import with absolute paths (or imports must be
                           relative).
*  `    --root-template="."`  sysl root directory for input template file (default: .)
*  `-a, --app-name= ...    `  name of the sysl app defined in sysl model. if there are multiple apps defined in sysl
                           model, code will be generated only for the given app
*  `-o, --outdir="."       `  output directory

## Arguments

Args:
  <MODULE>  Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.