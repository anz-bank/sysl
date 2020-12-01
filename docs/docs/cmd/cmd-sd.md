---
id: cmd-sd
title: Sequence Diagram
sidebar_label: sd
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

`sysl sd` lets you generate a sequence diagram originating from a single endpoint.

<img alt="Sequence diagram" src={useBaseUrl('img/diagramgen/seq-diagram-puml.png')}/>

## Usage

`usage: sysl sd [<flags>] <MODULE>...`

## Output Formats

The output file format can be specified via the extension passed into the -o flag.
Valid formats include .svg, .png, .uml, .puml, .plantuml, .html or .link

## Required Flags

- `-s, --endpoint=ENDPOINT ...`Include endpoint in sequence diagram
- `-a, --app=APP ...`Include all endpoints for app in sequence diagram (currently only works with
  templated --output). Use SYSL_SD_FILTERS env (a comma-list of shell globs) to limit
  the diagrams generated

## Optional Flags

Optional flags:

- `--endpoint_format="%(epname)"`
  Specify the format string for sequence diagram endpoints. May include %(epname),
  %(eplongname) and %(@foo) for attribute foo (default: %(epname))
- `--app_format="%(appname)"`Specify the format string for sequence diagram participants. May include %%(appname)
  and %%(@foo) for attribute foo (default: %(appname))
- `-t, --title=TITLE`diagram title
- `-p, --plantuml=PLANTUML`base url of PlantUML server (default: `SYSL_PLANTUML` or `http://localhost:8080/plantuml`
  see http://plantuml.com/server.html#install for more info)
- `-o, --output="%(epname).png"`output file (default: %(epname).png)
- `-b, --blackbox=BLACKBOX ...`Input blackboxes in the format App <- Endpoint=Some description, repeat '-b App <-
  Endpoint=Some description' to set multiple blackboxes
- `-g, --groupby=GROUPBY`Enter the groupby attribute (apps having the same attribute value are grouped
  together in one box

[More common optional flags](common-flags.md)
[Diagram format arguments](format-diagram.md)

## Arguments

Args:
`<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.

## Examples

### Simple Sequence Diagram

Command line

```bash
sysl sd -s "GroceryStore <- POST /checkout" GroceryStore.sysl -o checkout.png
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

```

<img alt="Sequence diagram" src={useBaseUrl('img/diagramgen/seq-diagram-puml.png')}/>
