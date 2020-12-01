---
id: cmd-integrations
title: Integration Diagram
sidebar_label: integrations
keywords:
  - command
---

import useBaseUrl from '@docusaurus/useBaseUrl';

:::info
We are currently in the process of migrating from PlantUML to Mermaid for our diagram generation. This will remove the external dependency on PlantUML and offer a better user experience. Diagram generation with mermaid is currently supported for integration diagrams and sequence diagrams only. For more details, check out [sysl diagram](cmd-diagram.md)
:::

:::info
This command requires the `SYSL_PLANTUML` environment variable to be set or passed in as a flag. Follow the instructions [here](../plantuml.md) for more details
:::

---

`sysl integrations` lets you generate integration diagrams. The command requires a project to be specified to produce an integration diagram. Refer to the examples for more details.

## Usage

```bash
usage: sysl integrations [<flags>] <MODULE>
```

Aliases

```bash
sysl ints [<flags>] <MODULE>
```

## Output Formats

The output file format can be specified via the extension passed into the -o flag.
Valid formats include .svg, .png, .uml, .puml, .plantuml, .html or .link

## Required Flags

- `-j, --project=PROJECT` project pseudo-app to render

## Optional Flags

- `-t, --title=TITLE` diagram title
- `-p, --plantuml=PLANTUML` base url of PlantUML server (default: `SYSL_PLANTUML` or `http://localhost:8080/plantuml`
- `` see http://plantuml.com/server.html#install for more info)
- `-o, --output="%(epname).png"` output file(default: %(epname).png)

- `--filter=FILTER` Only generate diagrams whose output paths match a pattern
- `-e, --exclude=EXCLUDE ...` apps to exclude
- `-c, --clustered` group integration components into clusters
- `--epa` produce and EPA integration view

[More common optional flags](common-flags.md)

[Diagram format syntax](format-diagram.md)

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.

## Examples

### Simple Integration Diagram

Command line

```bash
sysl integrations -o epa.png --project Project GroceryStore.sysl
```

```sysl title="Input Sysl file: GroceryStore.sysl"
GroceryStore:
    /checkout:
        POST?payment_info=string:
            Payment <- POST /validate
            Payment <- POST /pay
            | Checks out the specified cart
            return ok <: string

Payment:
    /validate:
        POST?payment_info=string:
            | Validates payment information
            return 200 <: string

    /pay:
        POST:
            | Processes a payment
            return ok <: string

Project [appfmt="%(appname)"]:
    _:
        GroceryStore
        Payment
```

<img alt="Integration diagram" src={useBaseUrl('img/diagramgen/integration-diagram-puml.png')}/>

### Endpoint Analysis Diagram

Command line

```bash
sysl integrations -o epa.png --project Project --epa GroceryStore.sysl
```

```sysl title="Input Sysl file: GroceryStore.sysl"
GroceryStore:
    /checkout:
        POST?payment_info=string:
            Payment <- POST /validate
            Payment <- POST /pay
            | Checks out the specified cart
            return ok <: string

Payment:
    /validate:
        POST?payment_info=string:
            | Validates payment information
            return 200 <: string

    /pay:
        POST:
            | Processes a payment
            return ok <: string

Project [appfmt="%(appname)"]:
    _:
        GroceryStore
        Payment
```

<img alt="EPA diagram" src={useBaseUrl('img/diagramgen/epa-diagram-puml.png')}/>
