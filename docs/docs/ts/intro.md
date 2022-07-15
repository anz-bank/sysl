---
id: intro
title: Introduction
---

# TypeScript for Sysl

The NPM package `@anz-bank/sysl` is a library designed for working with Sysl models in TypeScript. It provides a canonical set of types to describe common concepts and patterns to minimise the volume of decisions and code required to build Sysl tools.

## What is a Sysl model?

Sysl is a language for describing software systems. Once a software system is described, the `sysl` CLI tool can compile it into a portable Protobuf message that can be read and used in any language.

The data that goes into the Protobuf message is the **model** that was captured in the Sysl language. It contains all the same information as the source file, including all imports, but is easy for programs to read and transform.

Any program that generates something from Sysl should take a Sysl model as input (rather than the raw text of a Sysl source file). And any tool that generates Sysl should first produce a Sysl model, then serialize it to a string.

## How do I use it?

Sysl is very flexible, and Sysl models can be used for many different applications - there is no one-size-fits-all way to work with them.

Instead, multiple models are made available so you can use the one that's most appropriate for your use case.

The main models will be covered on the following pages. They are:

- [Protobuf model](./pbmodel.md)
- [TypeScript model](./model.md)
