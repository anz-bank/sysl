---
id: source-context
title: Source Context
keywords:
  - language
---

Source Context is a bundle of metadata attached to each Sysl element that describes where in the source that element was described.

Both the start and end locations of each element are recorded using zero-based line and column indexes for each. (Note that in compiled output, actual zero values are optimised away and may appear to be missing.)

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
        "part": ["Foo", "Bar"]
      },
      "sourceContexts": [
        {
          "file": "example1.sysl",
          "start": {
            "line": 0,
            "col": 0
          },
          "end": {
            "line": 2,
            "col": 10
          }
        },
        {
          "file": "example1.sysl",
          "start": {
            "line": 4,
            "col": 0
          },
          "end": {
            "line": 6,
            "col": 10
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
        "part": ["Foo", "Bar"]
      },
      "sourceContexts": [
        {
          "file": "example2.sysl",
          "start": {
            "line": 0,
            "col": 0
          },
          "end": {
            "line": 4,
            "col": 10
          }
        }
      ]
    }
  }
}
```

The `sourceContext` field is deprecated and contains only a single source instance. This value is presently retained in order to support legacy implementations.
