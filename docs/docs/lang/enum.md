---
id: enum
title: Enum
keywords:
  - language
---

An enum (enumeration) is a kind of [Type](./type.md) that can have one of a finite set of values. It can be referenced by name just like any other kind of Type.

## Syntax

`!enum` followed by an [Identifier](./identifiers.md) declares an Enum. It has a sequence of child `name: value` pairs corresponding to each possible value of the Enum.

:::caution
The syntax for enumerations will likely change from `name: value` to `name = value` in future. Limitations in the current parser prevent parsing of the second form.
:::

## Example

```
Server:
  Login (request <: Server.LoginData):
    return Server.Code

  !enum Code:
    success: 1
    invalid: 2
    tooManyAttempts: 3

  !type LoginData: ...
```

## See also

- [Type](./type.md)
