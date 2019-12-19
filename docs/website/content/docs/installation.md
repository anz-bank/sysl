---
title: "Installation"
description: "Sysl can be installed on Windows, MacOS and Linux - follow this guide."
date: 2018-02-27T15:51:27+11:00
weight: 10
draft: false
bref: "Sysl can be installed on Windows, MacOS and Linux - follow this guide"
toc: true
---

Sysl is a CLI (Command Line Interface) that executes with the `sysl` command.

Install it with

    GO111MODULE=on go get github.com/anz-bank/sysl/cmd/sysl@v0.4.0

Note: Do NOT run it from inside a Go source directory that is module enabled,
otherwise it gets added to go.mod/go.sum.

## Prerequisites

- [Go 1.13](https://golang.org/doc/install)
- PlantUML environment variable for diagram generation:
  `export SYSL_PLANTUML=http://www.plantuml.com/plantuml`

If the external PlantUML service is not suitable for your use case, you can run
a local instance of the [PlantUML
server](https://hub.docker.com/r/plantuml/plantuml-server/) using the docker
image and point the `SYSL_PLANTUML` environment variable to that instance.
