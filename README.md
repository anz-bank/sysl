# Sysl
Match your system implementation and design as consistent as possible

[![Latest Release](https://img.shields.io/github/v/release/anz-bank/sysl?color=%2300ADD8)](https://github.com/anz-bank/sysl/releases)
[![Codecov](https://img.shields.io/codecov/c/github/anz-bank/sysl/master.svg)](https://codecov.io/gh/anz-bank/sysl/branch/master)

[![GitHub Actions Release status](https://github.com/anz-bank/sysl/workflows/Release/badge.svg)](https://github.com/anz-bank/sysl/actions?query=workflow%3ARelease)
[![GitHub Actions Go-Darwin status](https://github.com/anz-bank/sysl/workflows/Go-Darwin/badge.svg)](https://github.com/anz-bank/sysl/actions?query=workflow%3AGo-Darwin)
[![GitHub Actions Go-Linux status](https://github.com/anz-bank/sysl/workflows/Go-Linux/badge.svg)](https://github.com/anz-bank/sysl/actions?query=workflow%3AGo-Linux)
[![GitHub Actions Go-Windows status](https://github.com/anz-bank/sysl/workflows/Go-Windows/badge.svg)](https://github.com/anz-bank/sysl/actions?query=workflow%3AGo-Windows)

Sysl (pronounced "sizzle") is a open source system specification language. Using Sysl, you
can specify systems, endpoints, endpoint behaviour, data models and data
transformations. The Sysl compiler automatically generates sequence diagrams,
integrations, and other views. It also offers a range of code generation
options, all from one common Sysl specification.

The set of outputs is open-ended and allows for your own extensions. Sysl has
been created with extensibility in mind and it will grow to support other
representations over time.

## Usage Examples

[Sysl by Example](https://github.service.anz/pages/sysl/syslbyexample/docs/byexample/) is a hands-on introduction to Sysl using annotated examples.


## Installation

Here are several approach to get start using Sysl:

### Install the pre-compiled binary

Download the pre-compiled binaries from the [releases page](https://github.com/anz-bank/sysl/releases) and copy to the desired location.

### Go get it

```bash
# make sure you've installed go in your computer at first
$ go version

# go get it
$ go get -u github.com/anz-bank/sysl/cmd/sysl

# check it works
$ sysl help
```

### Running with Docker

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


### Compiling from source

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

## Documentation

Documentation is hosted live at [https://sysl.io](https://sysl.io).

## Contributing

We encourage contributions to this project! Please have a look at the
[contributing guide](docs/CONTRIBUTING.md) for more information.

## Contributors

This project exists thanks to [all the people who contribute](https://github.com/anz-bank/sysl/graphs/contributors).

## Versioning

We use [Semver](https://semver.org/) for versioning. For the versions available, see the [releases](https://github.com/anz-bank/sysl/releases) on this repository.

## License

[![License](https://img.shields.io/github/license/anz-bank/sysl)](https://github.com/anz-bank/sysl/blob/master/LICENSE)

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details



