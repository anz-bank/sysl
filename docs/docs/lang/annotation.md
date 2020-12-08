---
id: annotation
title: Annotation
keywords:
  - language
---

An Annotation associates a key-value pair of metadata with an element of the Sysl model. The key is always a string, and the value can be a string or an array of strings (one- or two-dimensional).

:::info
This is an extension point of the Sysl language, allowing the capture of arbitrary information, especially concepts specific to a domain or organization.
:::

## Syntax

An Annotation appears in square brackets following an element, with an `=` separating the key and value. If the value is a string, it is wrapped in double quotes; an array is wrapped in square brackets. For example:

```sysl
# An Application with a string Annotation.
App [package="io.sysl.example"]:
    ...

# A Table Annotated with an array of string arrays.
App:
    !table [indexes=[["name:foo", "key_parts:col1,col2"], ["name:bar", "key_parts:col3,col4"]]:
        ...
```

Annotations can also be nested as children when prefixed with `@`:

```sysl
App:
    @package = "io.sysl.example"
    @description =:
        | Summary of the application.
        |
        | A long string can be split over multiple lines, with
        | two newlines to separate paragraphs.
    @some_array = ["foo", "bar"]
    @some_2d_array = [["foo"], "bar"]
```

See [formatting best practices](../best-practices/formatting.md) for guidance on how to choose the right style.

## Standard Annotations

Some Annotations have a standard meaning in all Sysl models. They should be used for that purpose, and not for any other. They include:

- `description`: Describes an element.
- `physical_name`: Records the name of a table or field in an actual, physical implementation (as opposed to the logical name). You are encouraged to use the physical name as the regular name, but sometimes this is impossible.

## See also

- [Tag](./tag.md): Annotation without a value
