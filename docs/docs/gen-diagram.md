---
id: gen-diagram
title: Diagram Generation
keywords:
  - diagram
  - mermaid
  - plantuml
  - generation
---

import useBaseUrl from '@docusaurus/useBaseUrl';

:::info
We are currently in the process of migrating from PlantUML to Mermaid for our diagram generation. This will remove the external dependency on PlantUML and offer a better user experience. Diagram generation with Mermaid is currently supported for integration diagrams and sequence diagrams only. For more details, check out [sysl diagram](cmd/cmd-diagram.md).
:::

---

Sysl lets you generate various diagrams from your specifications so that you can visualise your design as it evolves. These capabilities become more and more valuable as your project grows to include multiple services and complex dependencies.

## Integration Diagrams

Integration diagrams shows you which applications which make up your architecture and how they interact with each other.

<img alt="Integration diagram" src={useBaseUrl('img/diagramgen/int-diagram-mermaid.svg')}/>

For more details, refer to [Integration Diagram](cmd/cmd-integrations.md)

## Sequence Diagrams

Sequence diagrams show how a call to an endpoint propagates through your system.

<img alt="Sequence Diagram" src={useBaseUrl('img/diagramgen/attribs-Seq.png')} max-width="80%"/>

For more details, refer to [Sequence Diagram](cmd/cmd-sd.md)

## Data Model Diagrams

Data Model Diagrams show the relationship between your data types.

<img alt="Data Model diagram" src={useBaseUrl('img/diagramgen/data-diagram-mermaid.svg')}/>

For details on the command, refer to [Datamodel Diagram](cmd/cmd-datamodel.md)
