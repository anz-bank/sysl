---
id: primitives
title: Primitives
keywords:
  - language
---

Primitive types are simple [Type](./type.md)s that are built into Sysl. The [`!type`](./type.md) declaration can be used to create user-defined Types, composed of Primitives and other user-defined Types.

The basic, fairly self-explanatory Primitives are as follows:

- `any`
- `int`
- `float`
- `decimal`
- `string`
- `bytes`
- `date`
- `datetime`

There are also some more specialized, constrained versions of these Primitives:

- `int32`: int with bit width 32
- `int64`: int with bit width 64
- `float32`: float with bit width 32
- `float64`: float with bit width 64
- `decimal(p.s)`: decimal with precision `p` and scale `s` (e.g. `decimal(5.2)`)
- `string(max)`: string with maximum length (e.g. `string(100)`)
- `string(min..max)`: string with minimum and maximum lengths (e.g. `string(10..12)`)
- `xml`: an XML string

## See also

- [Type](./type.md)
- [Field](./field.md)
