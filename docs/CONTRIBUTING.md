# Contributing

We want to make contributing to this project as easy and transparent as possible, whether it's:

- Reporting a bug
- Suggesting a feature or enhancement
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## We Develop with Github

We use github to host code, to track issues and feature requests, as well as accept pull requests.

## How to propose changes to the codebase

### Setup your machine

`Sysl` is written in [Go](https://golang.org/).

Prerequisites:

- `make`
- [Go 1.13+](https://golang.org/doc/install)
- `GOPATH` env var set (this is used by the make install command)
- [golangci-lint 1.23.8+](https://github.com/golangci/golangci-lint)

Clone `sysl` to create a local copy on your computer:

```sh
$ git clone https://github.com/anz-bank/sysl.git
```

Get dependencies using go modules (needs go 1.11+)

```sh
$ go get ./...
```

A good way of making sure everything is all right is running the test suite:

```sh
$ make test
```

Run all the linters(we use [golangci-lint](https://github.com/golangci/golangci-lint)):

```sh
$ make lint
```

### Commit your changes

We use [Github Flow](https://guides.github.com/introduction/flow/index.html), so all code changes happen through Pull Requests.

1. Fork the repo and create your branch from `master`.
2. Git commit your changes
   - If you've added code that should be tested, add tests.
   - If you've added code that makes the documentation out-of-date, update the documentation.
   - Ensure the test suite passes.
   - Make sure your code lints.
3. Git push and open a pull request against the master branch(attach **WIP** tag when the PR is still work in progress).
4. Merge it after it's reviewed and approved!

> The codebase structure refers to [this standard](https://github.com/golang-standards/project-layout)

> Commit messages should be well formatted, and to make that "standardized", we are using Conventional Commits.
> You can follow the documentation on [their website](https://www.conventionalcommits.org).

## How to report a bug or suggest a feature

We use GitHub [issues](https://github.com/anz-bank/sysl/issues) to track public bugs and collect enhancement suggestions. Report a bug or suggest a feature by [opening a new issue](https://github.com/anz-bank/sysl/issues/new/choose). Choose the issue template you want and follow the hints; it's that easy!

## How to publish a new release

Please follow the steps in the [releasing](releasing.md) documentaion.

## Any contributions you make will be under the Apache License 2.0

In short, when you submit code changes, your submissions are understood to be under the same [Apache License 2.0](https://github.com/anz-bank/sysl/blob/master/LICENSE) that covers the project. Feel free to contact the maintainers if that's a concern.
