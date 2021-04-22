---
id: source-context
title: Source Context
keywords:
  - language
---

Source Context is a bundle of metadata attached to each Sysl element that describes where in the source that element was described.

It's possible for elements within a Sysl specification to be specified more than once. When this occurs, Sysl merges the elements together. For example, the following specifications represent identical systems:

#### Example 1

```
Foo :: Bar:
    Endpoint1:
        ...

Foo :: Bar:
    Endpoint2:
        ...

```

#### Example 2

```
Foo :: Bar:
    Endpoint1:
        ...
    Endpoint2:
        ...

```

The `sourceContexts` field contains all locations where a given element can be found within source files. For example, using the examples above, the `Foo :: Bar` application would have the following `sourceContexts` values:

#### Example 1: Source Contexts

```json
{
 "apps": {
  "Foo :: Bar": {
   "name": {
    "part": [
     "Foo",
     "Bar"
    ]
   },
   "sourceContexts": [
    {
     "file": "example1.sysl",
     "start": {
      "line": 1,
      "col": 1
     },
     "end": {
      "line": 1,
      "col": 7
     }
    },
    {
     "file": "example1.sysl",
     "start": {
      "line": 5,
      "col": 1
     },
     "end": {
      "line": 5,
      "col": 7
     }
    }
   ]
  }
 }
}
```

#### Example 2: Source Contexts

```json
{
 "apps": {
  "Foo :: Bar": {
   "name": {
    "part": [
     "Foo",
     "Bar"
    ]
   },
   "sourceContexts": [
    {
     "file": "example2.sysl",
     "start": {
      "line": 1,
      "col": 1
     },
     "end": {
      "line": 1,
      "col": 7
     }
    }
   ]
  }
 }
}
```

The `sourceContext` field is deprecated and contains only a single source instance. This value is presently retained in order to support legacy implementations.
