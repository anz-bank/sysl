---
title: "sd"
description: "Generates a sequence diagram"
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Generates a sequence diagram"
toc: true
---

The `sysl sd` command is used to generate a sequence diagram. For an example, refer to <https://sysl.io/docs/byexample/sequence-diagrams/>

## Usage

`usage: sysl sd [<flags>] <MODULE>...`

## Required Flags


## Optional Flags

Optional flags:

*  `    --help                     `Show context-sensitive help (also try --help-long and --help-man).
*  `    --version                  `Show application version.
*  `    --log="warn"               `log level: [debug,info,warn,trace,off]
*  `-v, --verbose                  `enable verbose logging
*  `    --root=ROOT                `sysl root directory for input model file. If root is not found, the module directory
                                    becomes the root, but the module can not import with absolute paths (or imports must
                                    be relative).
*  `    --endpoint_format="%(epname)"`
                                    Specify the format string for sequence diagram endpoints. May include %(epname),
                                    %(eplongname) and %(@foo) for attribute foo (default: %(epname))
*  `    --app_format="%(appname)"  `Specify the format string for sequence diagram participants. May include %%(appname)
                                    and %%(@foo) for attribute foo (default: %(appname))
*  `-t, --title=TITLE              `diagram title
*  `-p, --plantuml=PLANTUML        `base url of plantuml server (default: SYSL_PLANTUML or http://localhost:8080/plantuml
                                    see http://plantuml.com/server.html#install for more info)
*  `-o, --output="%(epname).png"   `output file (default: %(epname).png)
*  `-s, --endpoint=ENDPOINT ...    `Include endpoint in sequence diagram
*  `-a, --app=APP ...              `Include all endpoints for app in sequence diagram (currently only works with
                                    templated --output). Use SYSL_SD_FILTERS env (a comma-list of shell globs) to limit
                                    the diagrams generated
*  `-b, --blackbox=BLACKBOX ...    `Input blackboxes in the format App <- Endpoint=Some description, repeat '-b App <-
                                    Endpoint=Some description' to set multiple blackboxes
*  `-g, --groupby=GROUPBY          `Enter the groupby attribute (apps having the same attribute value are grouped
                                 together in one box

## Arguments

Args:
  <MODULE>  Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` filetype is optional.