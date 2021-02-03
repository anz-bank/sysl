---
id: tag
title: Tag
keywords:
  - language
---

A Tag is an [Annotation](./annotation.md) without a value; its mere presence indicates some characteristic of the Tagged element.

:::info
This is an extension point of the Sysl language, allowing the capture of arbitrary information, especially concepts specific to a domain or organization.
:::

## Syntax

A Tag is a name preceded by a `~` in square brackets following the Tagged element. For example, an [Application](./application.md) with the `db` tag would be written as:

```sysl
App [~db]:
    ...
```

## Standard Tags

Some Tags have a standard meaning in all Sysl models. They should be used for that purpose, and not for any other. They include:

- `db`: The [Application](./application.md) is (or sometimes has) a database.
- `human`: The Application represents a human user.
- `external`: The Application sits outside the organization. This will, for example, force it to the far right of a sequence diagram.
- `pk`: Primary Key in a [Table](./table.md).
- `unique`: Unique Key in a Table.
- `hex`: The [Field](./field.md)'s value should be interpreted as a hexadecimal string.
- `body`: The [Endpoint](./endpoint.md) Parameter is expected in the body of the HTTP request.
- `header`: The Endpoint Parameter is expected in the header of the HTTP request.
- `cookie`: The Endpoint Parameter passed in the Cookie header, such as `Cookie: debug=0; csrftoken=BUSe35dohU3O1MZvDCU`

## See also

- [Application](./application.md): parent element
- [Annotation](./annotation.md): Tag with a value
