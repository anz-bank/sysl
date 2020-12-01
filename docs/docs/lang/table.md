---
id: table
title: Table
keywords:
  - language
---

A Table is a kind of [Type](./type.md) that represents a physical table in a database. Each column in the Table is represented by a [Field](./field.md).

A Table is defined with the `!table` keyword (in contrast to a regular [Type](./type.md) which uses `!type`) followed by the Table's name. Nested in the Table definition are a list of Field definitions.

## Keys

One _key_ difference between a Table and an ordinary Type is that a Table has more a formal notion of Keys (like in a database).

- The set of Fields that comprise the Table's Primary Key each have the [Tag](./tag.md) `~pk`.
- A Field whose Type is a reference to another other Field is implicitly considered a Foreign Key. The expectation is that its value will be equal to an existing value of the referenced Field.
- A Field whose values must be unique (i.e. a Unique Key) should have the Tag `~unique`. If a set of Fields must together have a unique set of values... TODO: Annotation on the `!table`?

## Indexes

TODO: `indexes=[...]` Annotation on the `!table` declaration.

## Cascading

TODO: More Annotations.

## See also

- [Application](./application.md): parent element
- [Type](./type.md)
- [Field](./field.md)
