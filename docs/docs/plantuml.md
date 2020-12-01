---
id: plantuml
title: PlantUML Setup
sidebar_label: PlantUML Setup
---

The `SYSL_PLANTUML` environment variable must be configured or passed in when using the following sysl commands:

- `sysl sd`
- `sysl integrations`
- `sysl datamodel`

To configure your `SYSL_PLANTUML` environment variable, you can run the following

```bash
export SYSL_PLANTUML=<PLANTUML_SERVER_ADDRESS>
```

Alternatively, it can be passed into the command like so:

```bash
sysl sd -p PLANTUML_SERVER_ADDRESS
```

or set it before you run the command with

```bash
SYSL_PLANTUML=PLANTUML_SERVER_ADDRESS sysl sd ...
```

You can choose to use the public PlantUML service or run your own locally.

## Local PlantUML Service

```bash
docker run -d -p 8080:8080 --name plantuml plantuml/plantuml-server:jetty-v1.2020.14
export SYSL_PLANTUML=http://localhost:8080
```

## Public PlantUML Service

```bash
export SYSL_PLANTUML=http://www.plantuml.com/plantuml
```

:::warning
We strongly recommend you run a local version of PlantUML when working with sensitive system specifications.
:::
