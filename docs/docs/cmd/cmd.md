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

URL of PlantUML server. Sysl depends upon [PlantUML](http://plantuml.com/) for diagram generation.

- `SYSL_MODULES`

```
export SYSL_MODULES=on
```

Whether the sysl modules is enabled. Enable by default, set to "off" to disable sysl modules.

- `SYSL_CACHE`

Cache location in current directory, defaults to "sysl-modules" if SYSL_MODULES is enabled

- `SYSL_PROXY`

Proxy service to use, won't use SYSL_PROXY if not set

- `SYSL_TOKENS`

```
export SYSL_TOKENS=github.com:<GITHUB-PAT>
```

Setting `SYSL_TOKENS` with tokens (e.g. [GitHub personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)) for sysl to import specifications from private source via token.

- `SYSL_SSH_PRIVATE_KEY` and `SYSL_SSH_PASSPHRASE`

```
export SYSL_SSH_PRIVATE_KEY="/ssh/private/key/filepath"
export SYSL_SSH_PASSPHRASE="abcdef"
```

Setting `SYSL_SSH_PRIVATE_KEY` with filepath to your [SSH private key](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for sysl to import specifications from private source via SSH.

For more details, refer to [Installation doc](../installation.md)
