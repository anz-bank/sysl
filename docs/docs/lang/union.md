---
id: union
title: Union
keywords:
  - language
---

A Union represents a [Type](./type.md) whose value can be one of a set of other Types.

This is conceptually similar to an [Enum](./enum.md), in that it constrains the set of possible values of instances of the Union. However it is less restrictive than a single Type, because values can be of any of the Types in the Union.

## Example

The following example specifies that `User.id` can only be one of `string`, `int32` or `TypeUUID`:

```
  !type TypeUUID:
    id <: string

  !union UnionType:
    string
    int32
    TypeUUID

  !type User:
    id <: UnionType
```

## See also

- [Application](./application.md): parent element
- [Type](./type.md)
- [Enum](./enum.md)

External references:

- [Wikipedia page on union types](https://en.wikipedia.org/wiki/Union_type)
