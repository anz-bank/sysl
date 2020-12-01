---
id: cmd-diagram
title: Mermaid Diagram
sidebar_label: diagram
keywords:
  - command
  - mermaid
---

import useBaseUrl from '@docusaurus/useBaseUrl';

:::info
We are currently in the process of migrating from PlantUML to Mermaid for our diagram generation. This will remove the external dependency on PlantUML and offer a better user experience. Diagram generation with mermaid is currently supported for integration diagrams and sequence diagrams only.
:::

---

## Summary

`sysl diagram` lets you generate sequence or integration diagrams from your specification files using [MermaidJS](https://mermaidjs.github.io/#/).

## Usage

```bash
sysl diagram [<flags>] <MODULE>
```

## Output Formats

Currently SVG is the only supported output format.

## Optional Flags

- `-i, --integrationdiagram=INTEGRATIONDIAGRAM` Generate an integration diagram (Specify the application name)
- `-s, --sequencediagram=SEQUENCEDIAGRAM` Generate a sequence diagram (Specify 'appname->endpoint')
- `-e, --endpointanalysis` Generate an integration diagram with its endpoints (Specify 'true')
- `-d, --datadiagram` Generate a Data model diagram (Specify 'true')
- `-a, --app=APP` Optional flag to specify specific application
- `-e, --endpoint=ENDPOINT` Optional flag to specify endpoint
- `-o, --output="diagram.svg"` Output file (Default: diagram.svg)

[More common optional flags](common-flags.md)

## Sequence Diagram

```bash
sysl diagram -s grocerystore.sysl --app GroceryStore --endpoint "POST /checkout"
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

<img alt="Sequence Diagram" src={useBaseUrl('img/diagramgen/seq-diagram-mermaid.svg')} />

## Integration Diagram

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

```bash
sysl diagram -i grocerystore.sysl --app GroceryStore
```

<img alt="Integration Diagram" src={useBaseUrl('img/diagramgen/int-diagram-mermaid.svg')} />

## Data Model Diagram

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

```

```bash
sysl diagram -d Payment.sysl
```

<img alt="Data Model Diagram" src={useBaseUrl('img/diagramgen/data-diagram-mermaid.svg')} />
