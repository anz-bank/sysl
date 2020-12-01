---
id: cmd-datamodel
title: Data Model Diagram
sidebar_label: datamodel
keywords:
  - command
---

import useBaseUrl from '@docusaurus/useBaseUrl';

:::info
We are currently in the process of migrating from PlantUML to Mermaid for our diagram generation. This will remove the external dependency on PlantUML and offer a better user experience. Diagram generation with Mermaid is currently supported for integration diagrams and sequence diagrams only. For more details, check out [sysl diagram](cmd-diagram.md).
:::

:::info
This command requires the `SYSL_PLANTUML` environment variable to be set or passed in as a flag. Follow the instructions [here](../plantuml.md) for more details.
:::

---

`sysl datamodel` generates data model diagrams for types defined in Sysl.

## Usage

```
usage: sysl datamodel [<flags>] <MODULE>
```

## Output Formats

The output file format can be specified via the extension of the filename passed into the `-o` flag.

Valid extensions include `.svg`, `.png`, `.uml`, `.puml`, `.plantuml`, `.html` or `.link`.

## Required Flags

Either the `project` flag or the `direct` flag must be passed in. The `direct` flag generates data model diagrams for ALL types defined in the input Sysl file. The `project` flag specifies a subset of applications to generate diagrams for. Refer to the [Project Datamodel Diagram](#project-datamodel-diagram) example for more info.

- `-j, --project=PROJECT` Generate diagrams only for applications specified in the specified project
- `-d, --direct` Generate diagrams for all applications and types

## Optional Flags

All flags are all optional.

Optional flags:

- `--class_format="%(classname)"`
- ` Specify the format string for data diagram participants. May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))
- `-t, --title=TITLE` Diagram title
- `-p, --plantuml=PLANTUML` base url of PlantUML server (default: `SYSL_PLANTUML` or `http://localhost:8080/plantuml` see http://plantuml.com/server.html#install for more info)
- `-o, --output="%(epname).png"` Output file (default: %(epname).png)
- `-f, --filter=FILTER` Only generate diagrams whose names match a pattern

[More common optional flags](common-flags.md)

[Diagram format arguments](format-diagram.md)

## Arguments

Args:

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.

## Examples

### Simple Datamodel Diagram

```bash
sysl datamodel -d CardInfo.sysl -o CardInfo.svg
```

```sysl title="Input Sysl file: CardInfo.sysl"
Payment:
    !type CardInfo:
        credit_card_number <: string:
            @sensitive="true"
        type <: string
```

<img alt="Payment Datamodel diagram" src={useBaseUrl('img/diagramgen/data-diagram-puml.svg')} />

### Compound Datamodel Diagram

```bash
sysl datamodel -d Payment.sysl -o Payment.svg
```

```sysl title="Input Sysl file: Payment.sysl"
Payment:
    !type CardInfo:
        credit_card_number <: string:
            @sensitive="true"
        type <: string
    !type Payment:
        CardInfo <: CardInfo
        Amount <: int
```

<img alt="Payment Datamodel diagram" src={useBaseUrl('img/diagramgen/data-diagram-compound-puml.svg')} />

### Project Datamodel Diagram

```bash
sysl datamodel -j Project Payment.sysl -o "%(epname).svg"
```

In this example we generate two diagrams using the -j flag to specify which applications to display

```sysl title="Input Sysl file: Payment.sysl"
Payment:
    !type CardInfo:
        credit_card_number <: string:
            @sensitive="true"
        type <: string
    !type Payment:
        CardInfo <: CardInfo
        Amount <: int

PaymentService:
    !type PaymentProvider:
        Provider <: string

Project:
    PaymentService:
        PaymentService
    Payment:
        Payment

```

This diagram only shows the types within the **Payment** application

<img alt="Payment Datamodel diagram" src={useBaseUrl('img/diagramgen/data-diagram-payment.svg')} />

This diagram only shows the types within the **PaymentService** application

<img alt="Payment Info Datamodel diagram" src={useBaseUrl('img/diagramgen/data-diagram-paymentinfo.svg')} />
