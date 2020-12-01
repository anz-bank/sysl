---
id: cmd
title: Command Documentation
sidebar_label: Command Documentation
keywords:
  - command
---

```
$ sysl help
usage: sysl [<flags>] <command> [<args> ...]

System Modelling Language Toolkit

Optional flags:
      --help        Show context-sensitive help (also try --help-long and --help-man).
      --version     Show application version.
      --log="warn"  log level: [off,debug,info,warn,trace]
  -v, --verbose     enable verbose logging
      --root=ROOT   sysl root directory for input model file. If root is not found, the module directory becomes the root, but the module can not import with absolute paths (or imports
                    must be relative).

Commands:
  help [<command>...]
    Show help.

  codegen --transform=TRANSFORM --grammar=GRAMMAR [<flags>] <MODULE>...
    Generate code

  datamodel [<flags>] <MODULE>...
    Generate data models

  diagram [<flags>] <MODULE>...
    Generate mermaid diagrams

  env
    Print sysl environment information.

  export [<flags>] <MODULE>...
    Export sysl to external types. Supported types: Swagger,openapi2,openapi3

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
    Starts the Sysl UI which displays a visual view of Apps defined in the sysl model.

  validate <MODULE>...
    Validate the sysl file
```

## Environment Variables

Several commands require environment variables to be set before they are able to correctly work.

- `SYSL_PLANTUML`

```
export SYSL_PLANTUML=http://www.plantuml.com/plantuml
```

- `SYSL_MODULES`

```
export SYSL_MODULES=on
```

Setting `SYSL_MODULES` to `on` means Sysl modules are enabled, `off` means disabled. By default, if this is not declared, Sysl modules are enabled.

For more details, refer to [Installation doc](../installation.md)
