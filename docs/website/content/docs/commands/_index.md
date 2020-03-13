---
title: "Sysl Commands"
description: "Sysl can be used via a simple command-line interface"
date: 2018-02-27T15:51:27+11:00
weight: 20
draft: false
bref: "Sysl can be used via a simple command-line interface"
layout: "commandslist"
toc: true
---

Sysl has a number of subcommands which can perform various tasks with **.sysl** files.

To view a list of available commands, run sysl with no arguments:

This documentation refers to [v0.8.0](https://github.com/anz-bank/sysl/releases/tag/v0.8.0) of the sysl CLI tool.

```
$ sysl
usage: sysl [<flags>] <command> [<args> ...]

System Modelling Language Toolkit

Optional flags:
      --help        Show context-sensitive help (also try --help-long and --help-man).
      --version     Show application version.
      --log="warn"  log level: [warn,trace,off,debug,info]
  -v, --verbose     enable verbose logging
      --root=ROOT   sysl root directory for input model file. If root is not found, the module directory
                    becomes the root, but the module can not import with absolute paths (or imports must be
                    relative).

Commands:
  help [<command>...]
    Show help.

  codegen --transform=TRANSFORM --grammar=GRAMMAR [<flags>] <MODULE>...
    Generate code

  datamodel [<flags>] <MODULE>...
    Generate data models

  env
    Print sysl environment information.

  export [<flags>] <MODULE>...
    Export sysl to external types. Supported types: Swagger

  generate-db-scripts [<flags>] <MODULE>...
    Generate db script

  generate-db-scripts-delta [<flags>] <MODULE>...
    Generate delta db scripts

  import --input=INPUT --app-name=APP-NAME [<flags>]
    Import foreign type to sysl. Supported types: [grammar, openapi, swagger, xsd]

  info
    Show binary information

  integrations [<flags>] <MODULE>...
    Generate integrations

  mod init [<name>]
    initializes and writes a new go.mod to the current directory

  protobuf [<flags>] <MODULE>...
    Generate textpb/json

  repl
    Enter a sysl REPL

  sd [<flags>] <MODULE>...
    Generate Sequence Diagram

  template --template=TEMPLATE --start=START [<flags>] <MODULE>...
    Apply a model to a template for custom text output

  test-rig --template=TEMPLATE [<flags>] <MODULE>...
    Generate test rig

  ui [<flags>] <MODULE>...
    Starts the Sysl UI which displays a visual view of apps defined in Sysl.

  validate <MODULE>...
    Validate the sysl file	
```

To get help for a specific command, pass `--help` into the subcommand to find out more about its usage.

e.g:
```
$ sysl import --help
usage: sysl import --input=INPUT --app-name=APP-NAME [<flags>]

Import foreign type to sysl. Supported types: [grammar, openapi, swagger, xsd]

Required flags:
  -i, --input=INPUT        input filename
  -a, --app-name=APP-NAME  name of the sysl app to define in sysl model.

Optional flags:
      --help                  Show context-sensitive help (also try --help-long and --help-man).
      --version               Show application version.
      --log="warn"            log level: [off,debug,info,warn,trace]
  -v, --verbose               enable verbose logging
      --root=ROOT             sysl root directory for input model file. If root is not found, the module
                              directory becomes the root, but the module can not import with absolute paths (or
                              imports must be relative).
  -p, --package=PACKAGE       name of the sysl package to define in sysl model.
  -o, --output="output.sysl"  output filename
  -f, --format=auto           format of the input filename, options: [auto, grammar, openapi, swagger, xsd]
```
