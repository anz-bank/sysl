---
id: field
title: Field
keywords:
  - language
---

A Field is an attribute of a [Type](./type.md), be it a `!type` declaration (aka a `Tuple`) or a `!table` (aka a `Relation`).

## Syntax

A Field is nested inside a Type definition. It is comprised of a name [Identifier](./identifiers.md), the "element of" operator (`<:`), and a Type definition. The Type is usually a [Primitive](./primitives.md) or a reference to a Type, but it can also be an inline definition of a new Type.

If the Type definition ends with a question mark (`?`), the Field's Type is nullable (by default all Types are presumed to be non-nullable).

## Example

A Type with some Primitive Fields:

```sysl
App:
    !type Type:
        foo <: int
        bar <: string?
```

## See also

- [Type](./type.md)
- [Table](./table.md)
