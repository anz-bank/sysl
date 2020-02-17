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

## Prerequisites

- [Go 1.13](https://golang.org/doc/install)
- There're extra prerequisites for several subcommands like `sysl sd`(sequence diagram generation):
	- `export SYSL_PLANTUML=http://www.plantuml.com/plantuml`

		Set the PlantUML environment variable for diagram generation.	
		If the external PlantUML service is not suitable for your use case, you can run a local instance of the [PlantUML server](https://hub.docker.com/r/plantuml/plantuml-server/) using the docker image and point the `SYSL_PLANTUML` environment variable to that instance.

---

Here are several approaches to get start using Sysl:

## Install the pre-compiled binary

Download the pre-compiled binaries from the [releases page](https://github.com/anz-bank/sysl/releases) and copy to the desired location.

## Go get it

```bash
# make sure you've installed go in your computer at first
$ go version

# go get it
$ GO111MODULE=on go get -u github.com/anz-bank/sysl/cmd/sysl

# check it works
$ sysl help
```

> Note: Do NOT run it from inside a Go source directory that is module enabled, otherwise it gets added to go.mod/go.sum.

## Running with Docker

You can also use it within a [Docker container](https://hub.docker.com/r/anzbank/sysl). To do that, you’ll need to execute something more-or-less like the following:

```bash
$ docker run --rm anzbank/sysl:latest help
```

```bash
$ docker run --rm \
  -v $PWD:/go/src/github.com/anz-bank/sysl \
  -w /go/src/github.com/anz-bank/sysl \
  anzbank/sysl:latest validate -v ./demo/examples/Modules/model_with_deps.sysl
```
We have used this [Dockerfile](Dockerfile) to create this image.


## Compiling from source

Here you have two options:

1. If you want to contribute to the project, please follow the steps on our [contributing guide](docs/CONTRIBUTING.md).
2. If just want to build from source for whatever reason, follow the steps bellow:

```bash
# clone it to create a local copy on your computer
$ git clone https://github.com/anz-bank/sysl.git
$ cd sysl

# get dependencies using go modules (needs go 1.11+)
$ go get ./...

# build
$ go build -o sysl ./cmd/sysl

# check it works
$ ./sysl help
```
