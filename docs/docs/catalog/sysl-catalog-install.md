---
id: sysl-catalog-install
title: Sysl Catalog Installation
sidebar_label: Installation
---

## Setup

### Docker

If you have Docker, you can run sysl-catalog with

```bash
docker run --rm -p 6900:6900 -v $(pwd):/usr/:ro anzbank/sysl-catalog:latest input.sysl --serve
```

:::warning
This is still work in progress, please use the method below for now.
:::

### Go Get

Install `sysl-catalog` if you have `go` installed using:

```bash
go get -u -v github.com/anz-bank/sysl-catalog
```
