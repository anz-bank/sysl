---
id: gen-docs
title: Documentation Generation
keywords:
  - documentation
  - docs
  - catalog
  - generation
---

import useBaseUrl from '@docusaurus/useBaseUrl';

[Reference](https://github.com/anz-bank/sysl-catalog)

## Summary

Sysl lets you generate documentation from your specifications so that they are always up-to-date. This is done using the `sysl-catalog` tool.

More Resources:

- [Check out a live demo](https://demo.sysl.io/)
- To install it, [follow the instructions here](catalog/sysl-catalog-install.md)
- [Click here for more details on the sysl-catalog CLI](catalog/sysl-catalog-cmd.md)

## Features

`sysl-catalog` generates high level integration diagrams to give a high level overview of your applications and how they integrate.

<img alt="sysl-catalog-integration" src={useBaseUrl('img/sysl/sysl-catalog-integration.png')} />

For each service, you can view:

- Every endpoint and their request and response types
- A sequence diagram for each endpoint
- An integration diagram with adjacent services
  <img alt="sysl-catalog-endpoint" src={useBaseUrl('img/sysl/sysl-catalog-endpoint.png')} />

For REST API specs, a [ReDoc](https://github.com/Redocly/redoc) view is also generated

<img alt="sysl-catalog-redoc" src={useBaseUrl('img/sysl/sysl-catalog-redoc.png')} />
