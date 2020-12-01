---
id: installation
title: Installation
---

`sysl` is a command-line tool written in [Go](https://golang.org). There are a variety of ways to install it, depending on your OS and use case.

## Summary

- Mac:
  ```
  brew install anz-bank/homebrew-sysl/sysl
  ```
- Docker:
  ```
  docker run --rm -it -v $HOME:$HOME -w $(pwd) anzbank/sysl:latest
  ```
- Go:
  ```
  GO111MODULE=on go get -u github.com/anz-bank/sysl/cmd/sysl
  ```
- Source:
  ```
  git clone https://github.com/anz-bank/sysl.git
  cd sysl
  make install
  ```
- Binary: download from the [GitHub releases page](https://github.com/anz-bank/sysl/releases) to your `PATH`

Check the installation with `sysl help`.

## Requirements

- [Golang](https://golang.org/doc/install) version >= 1.13 (check with `go version`).
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports).

## [(Mac) Homebrew](https://github.com/anz-bank/homebrew-sysl)

If you use a Mac and have [Homebrew](https://brew.sh/) installed, you can simply run the following commands in your terminal:

```
brew tap anz-bank/homebrew-sysl
brew install anz-bank/homebrew-sysl/sysl
```

## Pre-compiled binary

1. Download the pre-compiled binaries matching your OS from the [releases page](https://github.com/anz-bank/sysl/releases).

1. Uncompress the archive and move the `sysl` binary to your desired location:

   1. On your `PATH` to run it with `sysl`
   1. Elsewhere to run it with `./sysl`, or some other `path/to/sysl`

## Go get it

First make sure you've installed Go:

```bash
go version
```

Fetch the `sysl` command's Go module:

```bash
GO111MODULE=on go get -u github.com/anz-bank/sysl/cmd/sysl
```

:::caution
Do NOT run this from inside a Go source directory that is module enabled, otherwise it gets added to go.mod/go.sum.
:::

## Docker

You can use `sysl` within a [Docker container](https://hub.docker.com/r/anzbank/sysl) (created from [this Dockerfile](https://github.com/anz-bank/sysl/blob/master/Dockerfile)):

```bash
docker run --rm -it -v $HOME:$HOME -w $(pwd) anzbank/sysl:latest
```

For example:

```
docker run --rm \
  -v $PWD:/go/src/github.com/anz-bank/sysl \
  -w /go/src/github.com/anz-bank/sysl \
  anzbank/sysl:latest validate -v demo/examples/Modules/model_with_deps.sysl
```

Mac and Linux users can create an `alias` for the `sysl` command:

```
alias sysl="docker run --rm -it -v $HOME:$HOME -w $(pwd) anzbank/sysl:latest"
```

`sysl` can then be used from the same terminal window. Alternatively, add the `alias` to your `.bashrc` or `.zshrc` file to keep it permanently.

## Compile from source

Here you have two options:

1. If you want to contribute to the project, please follow the steps on our [contributing guide](https://github.com/anz-bank/sysl/blob/master/docs/CONTRIBUTING.md).
2. If just want to build from source, follow the steps below:

```
# clone it to create a local copy on your computer
git clone https://github.com/anz-bank/sysl.git
cd sysl
GOPATH=$(go env GOPATH) make install
```

## Try it out

If the installation worked, you should be able to run:

```bash
sysl
usage: sysl [<flags>] <command> [<args> ...]
...
```

You can always check yours setup of `sysl` with:

```bash
sysl --version
sysl info
sysl env
```

## VS Code Extension

Sysl has a VS Code extension which provides syntax highlighting for `.sysl` files. [Get it from here](https://marketplace.visualstudio.com/items?itemName=ANZ-BANK.vscode-sysl), or search Extensions for "sysl".
