# Sysl

`sysl.arrai` contains helper functions to make working with Sysl models more convenient.

Examples in this file assume you've imported `sysl.arrai` using the following syntax:
```arrai
let sysl = //{github.com/anz-bank/sysl/pkg/arrai/sysl};
```

## polish()

`polish` takes an imported Sysl model and returns a copy with a cleaner internal structure, making it easier to work with. It applies the following tweaks:

- Parse `return` statement payloads into `(status, type)` tuples.

## ~~normalize()~~ *[Deprecated]*

`normalize` takes an imported Sysl model and returns an equivalent relational model. The relational model is a tuple structured like a database, with tables for each Sysl concept (e.g. `apps`, `endpoints`, `types`, `calls`, etc.) that can be naturally joined with the `<&>` operator.

The relational model is a powerful option for transformations that resemble queries with joins and filters (i.e. `where` clauses). See the diagram scripts for examples.

## load()

Loads a sysl model that is stored in a protobuf file. To generate such a file for `myModel.sysl`, run the Sysl tool with the following arguments:
```sh
sysl protobuf --output=myModel.pb --mode=pb myModel.sysl
```

You can then to load the model into arr.ai, use the following:
```arrai
let model = sysl.load('myModel.pb');
```

## loadBytes()

Loads a sysl model protobuf-encoded bytes (e.g. from `//os.stdin`, `//[//encoding.bytes]{./path}` or `//os.file()`). For example, you can generate those bytes from `myModel.sysl` and pipe them into an arr.ai script `myScript.arrai` using the following command:
```sh
sysl protobuf --mode=pb myModel.sysl | arrai r myScript.arrai
```

**myScript.arrai** could then contain the following:
```arrai
let model = sysl.loadBytes(//os.stdin);
```

## normalizeNew()

After a model was loaded with `load()` or `loadBytes()`, it can then be transformed to the relational model using `normalizeNew()` as such:

```arrai
var relModel = sysl.newNormalize(model);
```
