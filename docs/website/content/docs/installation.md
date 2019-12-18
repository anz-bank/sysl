---
title: "Installation"
description: "Sysl can be installed on Windows, MacOS and Linux - follow this guide."
date: 2018-02-27T15:51:27+11:00
weight: 10
draft: false
bref: "Sysl can be installed on Windows, MacOS and Linux - follow this guide"
toc: true
---

Sysl is a CLI (Command Line Interface) that excecutes with the `sysl` command.

    GO111MODULE=on go get github.com/anz-bank/sysl/cmd/sysl@v0.4.0

Note: Do NOT run it from inside a Go source directory that is module enabled,
otherwise it gets added to go.mod/go.sum.

## Prerequisites

- [Go](https://golang.org)
- [PlantUML](https://hub.docker.com/r/plantuml/plantuml-server/) server for diagram generation for use if using the [external service](http://www.plantuml.com/plantuml/) is not appropriate

## Setting PlantUML Environment variable

In order to be able to generate diagrams the `SYSL_PLANTUML` Environment variable needs to be set

    export SYSL_PLANTUML=http://www.plantuml.com/plantuml
