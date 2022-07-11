# Sysl

Herein lies a TypeScript library for working with Sysl models.

## What is a Sysl model?

Sysl is a language for describing software systems. Once a software system is described, the `sysl` CLI tool can compile it into a portable Protobuf message that can be read and used in any language.

The data that goes into the Protobuf message is the **model** that was captured in the Sysl language. It contains all the same information as the source file, including all imports, but is easy for programs to read and transform.

Any program that generates something from Sysl should take a Sysl model as input (rather than the raw text of a Sysl source file). And any tool that generates Sysl should first produce a Sysl model, then serialize it to a string.

## Working with Sysl in TypeScript

TypeScript (and JavaScript) is broadly popular and works well in any environment (server, browser, CLI). Any schemas or logic you write can be packaged and reused on any platform (with the exception of native APIs file filesystem access), so it's a good default choice for any cross-platform system.

In this TypeScript library, the `model` classes provide common functionality for working with Sysl models - both transforming existing models, and synthesizing new models. As such it's important to understand the design and APIs of these classes.

**TODO**
