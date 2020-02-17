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
	- Sysl depends upon [PlantUML](http://plantuml.com/) for diagram generation. Some of the automated tests require a PlantUML dependency. Provide PlantUML access either via local installation or URL to remote service. Warning, for sensitive data the public service at www.plantuml.com is not suitable. You can use one of the following options to set up your environment:
		- execute `SYSL_PLANTUML=http://www.plantuml.com/plantuml`
		- add `export SYSL_PLANTUML=http://www.plantuml.com/plantuml` to your `.bashrc`
		  or similar
		- [install PlantUML](http://plantuml.com/starting) locally or run a local instance of the [PlantUML server](https://hub.docker.com/r/plantuml/plantuml-server/) using the docker image and run on port 8080. Otherwise you can refer to the [plantuml server guide](docs/plantUML_server.md)

---

Here are several approaches to get start using Sysl:

## Install the pre-compiled binary

Download the pre-compiled binaries from the [releases page](https://github.com/anz-bank/sysl/releases) and copy to the desired location.

## Go get it

```
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

```
$ docker run --rm anzbank/sysl:latest help
```

```
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

```
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
