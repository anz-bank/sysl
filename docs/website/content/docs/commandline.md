---
title: "Command line"
description: "Learn to use sysl and reljam command line and its arguments"
date: 2018-02-27T15:55:46+11:00
weight: 20
draft: false
bref: "Sysl and reljam command line arguments"
toc: true
---

Sysl consists of two executables: `sysl` and `reljam`. `sysl` is mainly concerned with diagram creation whereas `reljam` generates different types of source code output.

[//]: # "TODO(juliaogris):"
[//]: # "* Explain `--root`"
[//]: # "* Positional `/module` argument e.g. `sysl pb --help`"
[//]: # "* `sysl|reljam <subcommand> --help`"
[//]: # "* `--help`"
[//]: # "* `--version`"
[//]: # "* `--trace` (missing for `reljam`!)"

sysl
----
```
> sysl --help
usage: sysl [-h] [--no-validations] [--root ROOT] [--version] [--trace]
            {pb,textpb,data,ints,sd} ...

System Modelling Language Toolkit

positional arguments:
  {pb,textpb,data,ints,sd}
                        sub-commands
                        more help with: sysl <sub-command> --help
                        eg: sysl pb --help

optional arguments:
  -h, --help            show this help message and exit
  --no-validations, --nv
                        suppress validations
  --root ROOT, -r ROOT  sysl root directory for input files (default: .)
  --version, -v         show version number (semver.org standard)
  --trace, -t
```

reljam
------
```
> reljam --help
usage: reljam [-h] [--root ROOT] [--out OUT] [--entities ENTITIES]
              [--package PACKAGE] [--serializers SERIALIZERS] [--version]
              {model,facade,view,xsd,swagger,spring-rest-service} module app

sysl relational Java Model exporter

positional arguments:
  {model,facade,view,xsd,swagger,spring-rest-service}
                        Code generation mode
  module                Module to load
  app                   Application to export

optional arguments:
  -h, --help            show this help message and exit
  --root ROOT, -r ROOT  sysl system root directory
  --out OUT, -o OUT     Output root directory
  --entities ENTITIES   Commalist of entities that are expected to have
                        corresponding output files generated. This is for
                        verification only. It doesn’t determine which files
                        are output.
  --package PACKAGE     Package expected to be used for generated classes.
                        This is for verification only. It doesn’t determine
                        the package used.
  --serializers SERIALIZERS
                        Control output of XML and JSON serialization code.
  --version, -v         show version number (semver.org standard)
```
