---
id: identifiers
title: Identifiers
keywords:
  - language
---

An Identifier is the name given to something, such as an Application, Endpoint, Type, Field, etc.

## Whitespace

Sysl supports whitespace within identifiers, however its use is strongly discouraged.

As an alternative, consider adhering to standard naming conventions and use the "long name" of an Application:

```
GCP "Google Cloud Platform":
  ...
```

## Reserved words

Identifiers cannot be named any of the following reserved words:

- `alt`
- `any`
- `as`
- `bool`
- `bytes`
- `date`
- `datetime`
- `decimal`
- `each`
- `else`
- `float`
- `float32`
- `float64`
- `for`
- `if`
- `int`
- `int32`
- `int64`
- `loop`
- `string`
- `until`
- `while`

## Special Characters

If special characters such as `:` or `.` are needed in the name of a Type or Endpoint, they can be expressed using their URL encoded equivalent instead. For reference, see https://www.urlencoder.org/.

## Multiple Declarations

Sysl allows an element with the same Identifier to be defined multiple times. If multiple definitions of the same thing occur in the same context (whether inline or via [Import](./import.md)), the details of the definitions will be merged.

:::caution
Sysl does not have redefinition errors, but the results of the merging may not always be what you expect or want. As such, Sysl logs a message for each element that is redefined and merged.
:::

For example, the follow Sysl model would be treated as a single `UserService` Application containing two Endpoints:

```
UserService:
  Login: ...

UserService:
  Register: ...
```

This capability is particularly useful for splitting the definition of a very large model across multiple files, and importing them into a unified model.

## See also

- https://www.urlencoder.org/ for encoding and decoding special characters.
