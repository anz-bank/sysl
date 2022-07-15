---
id: pbmodel
title: Protobuf Model
---

## Why Protobuf?

Protocol buffers (protobufs) are widely used to define data schemas. Given a protobuf schema, protobuf libraries (available in many languages) can be used to encode data into protobuf messages, serialize data for transfer over a network, and decode data back to messages.

Protobuf schemas are forwards and backwards compatible, and message serialization is very efficient. All of this makes protobufs suitable for encoding a Sysl specification into a machine-readable form.

This is what the `sysl protobuf` (aka `sysl pb`) command does: parse the source, populate a protobuf message with [sysl.proto](https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto), then serialize that message accordin to the `--mode` flag.

## What's the model

[sysl.proto](https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto) is the protobuf schema, and the `pbModel` module is the equivalent model in the TypeScript library. It's not very pleasant to work with (that's what [`model`](./model.md) is for), but it's the first step to get Sysl protobuf message data into TypeScript.

Indeed, the `pbModel` is designed specifically to deserialize a JSON-encoded Sysl protobuf message.

### Module

The model describes a `Module`. If multiple files were included in the Sysl spec via `import` statements, the `Module` is the result of resolving and merging their contents.

The `Module` contains primarily:
- `imports`: any `import` statements that were included
- `apps`: map of `Application`s indexed by name

`Application` is the basic unit of a Sysl model, containing primarily `types` for modelling data and `endpoints` for modelling APIs, each maps of objects indexed by name.

### Types

At the `Application` level, a `Type` represents either a tuple (`!type`), a relation (`!table`), an enum (`!enum`), union/one-of (`!union`) or an alias (`!alias`). The kind of a `Type` is determined based on which of its properties is populated (a `oneof` in the protobuf schema). Most types will be tuples or relations, which have respectively a `tuple.attrDefs` or `relation.attrDefs` property listing their fields (aka columns).

Each field is modelled with the same `Type` class, since it has all the same properties, except it *has a* type rather than *is a* type. The kind of a field is typically `primitive` (e.g. `int`, `string`, `datetime`) or `type_ref` (reference to another `Type`). Fields representing collections will have `sequence` or `set` populated, wrapping another `Type`.

### Endpoints

An `Endpoint` contains a list of `Statement`s, each of a certain kind, and some of which may have child `Statement`s of their own. As with a `Type`, a `Statement`'s kind is determined by which of its properties is populated:

- `action`: a simple string (pseudocode, like `validate the request`)
- `call`: a call to another endpoint (aka an integration)
- `ret` (return): indicates the `Type` returned by the endpoint in a particular case (e.g. `ok`, `error`, `200` or `404`), like `return ok <: MyResponse`
- `cond`, `loop`, `loop_n`, `foreach`, `alt`, `group`: a branch, loop, or other block with nested `Statement`s

An `Endpoint` may also have a list of `Param`s (parameters) in its `param` property which (like fields) have a name and `Type`. If it's a REST `Endpoint`, it will have `rest_params` as well, describing the HTTP `method`, `path`, `url_param`s and `query_param`s.

### Metadata

Finally, almost all elements (apps, types, fields, endpoints, statements) have additional metadata:

`attrs` is a map of `Attribute`s, indexed by name. The name and `Attribute` value is an arbitrary key-value pair the adds metadata to the element. The value of the `Attribute` is either `s` (string), `i` (integer), `n` (double), or `a` (array with more `Attribute`s in its `elt` property) depending on the kind of value stored.

:::info
Note: we often refer to `attrs` as "annotations" (e.g. `@foo = "bar"`), except for the special `patterns` attr, the values of which are `tags` (e.g. `[~foo]`).
:::

`source_contexts`: a list of locations in the source code where the element was defined. Since the same element can be redefined in multiple files (e.g. with different sets of `attrs`) and merged together, there may be multiple locations. They are listed in the order they were seen by the parser.

### More

These elements make up the basic structure of the protobuf model, but there more detail and nuance if you need it (e.g. optional fields, type constraints, foreign keys, mixins, pub/sub topics, call statement arguments). You can learn more about these in the [Language Reference](../lang/intro.md).
