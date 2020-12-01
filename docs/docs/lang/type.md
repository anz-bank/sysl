---
id: type
title: Type
keywords:
  - language
---

A Type is a data schema. It defines the characteristics of some data in the system.

There are various _kinds_ of Type that are specialized to model different kinds of data. The most common meaning of Type in Sysl is the `!type` declaration, which is in fact the specification of a `Tuple`, a kind of Type with multiple [Field](./field.md)s. However each Field also has its own Type, so the term is overloaded.

## Kinds of Type

```
any         as          bool        bytes       date
datetime    decimal     else        float       float32
float64     if          int         int32       int64       string
```

## See also

- [Application](./application.md): parent element

Kinds of Type:

- [Primitives](./primitives.md)
- [Table](./table.md)
- [Enum](./enum.md)
- [Union](./union.md)
- [Alias](./alias.md)
