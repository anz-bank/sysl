---
id: model
title: TypeScript Model
---

TypeScript (and JavaScript) is broadly popular and works well in any environment (server, browser, CLI). Any schemas or logic you write can be packaged and reused on any platform (with the exception of native APIs file filesystem access), so it's a good default choice for any cross-platform system.

Compared to protocol buffers (protobufs), which are used for low-level schema specification and data transfer (see [Protobuf Model](./pbmodel.md)), TypeScript is a general-purpose programming language and can provide much more functionality on top of the basic schema. The classes in the `model` module provide this functionality.

## What is the TypeScript model?

This model has a similar structure to the protobuf model, i.e. `Application`s containing `Type`s, `Endpoint`s, etc. However its usability has been improved in a few ways:

- Hide concepts that are deprecated or not commonly used (e.g. `!view`, `list`, `map`)
- Rename concepts to conform to common usage (e.g. `annos` and `tags` instead of `attrs`)
- Expose constructors that take named arguments with sensible defaults
- Provide methods to extract data for common use cases (e.g. filtering elements by the files they are defined in)
- Each element has a `toSysl()` function to serialize to Sysl source

## Getting Started

First you'll need to install the `@anz-bank/sysl` package from NPM:

```sh
yarn add @anz-bank/sysl
```

You can import it like so:

```ts
import * as sysl from '@anz-bank/sysl';
// or to be more specific
import * as model from '@anz-bank/sysl/dist/model';
// or
import { Model, Application } from '@anz-bank/sysl/dist/model';
```

The `Model` class represents a whole Sysl model. There are two ways to get an instance of `Model`:

1. **Construct one from scratch**: this is appropriate if you are generating Sysl yourself (e.g. in a custom importer).
2. **Generate one from a Sysl source or `pbModel`**: this is the typical case if you're transforming an existing Sysl specification.

### Generating from `pbModel`

Given a Sysl specification, there is a standard process of loading it into a `pbModel`, [described previously](./pbmodel.md).

To load Sysl source text into a `Model`, you can call `Model.fromText(...)`.

Behind the scenes, this first loads the data into the `pbModel`. Each element in the `pbModel` has a `.toModel()` method, including the `Model` class itself, which returns an equivalent object in the TypeScript model.


### Constructing from scratch

Creating an element in the TypeScript model is as simple as invoking a constructor or factory method. For example:

```ts
const model = new Model();
```

If you already have model contents (e.g. a list of `Application`s) you can pass them into the constructor:

```ts
const app = new Application({name: ["Foo", "Bar"]});
const model = new Model({apps: [app]});
```

If not, you can construct them separately and add them to the `Model` after the fact:

```ts
const model = new Model();
const app = new Application({name: ["Foo", "Bar"]});
model.apps.push(app);
```

## Usage

Once you have a Sysl model, there's a lot you can do with. Here are a few common examples:

### Serialize

The simplest use case is to generate human-readable Sysl source code containing all the information in your model (e.g. to version control in Git).

Serializing a whole merged model into a single file is as simple as calling `model.toSysl()`. This returns a string that you can write to stdout or the file system.

:::info
Serializing a merged model back to multiple files is possible (based on source context in `locations`), but not yet implemented by the library.
:::

### Transform

Often you'll want to produce some artefact from your Sysl model, such as rendering a diagram, exporting to another format, or generating code or config. This will involve traversing the Sysl model and building up some other model to be serialized.

:::info
The library does not yet provide helper functions for walking the Sysl model, but they are planned.
:::

### Validate

Validating the content of a Sysl model can be thought of as a transform from a model to a list of problems with the content. If the list is empty, the model is valid.

:::info
The library does not yet provide helper functions for validating a Sysl model, but they are planned.
:::
