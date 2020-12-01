---
id: formatting
title: Formatting
keywords:
  - standard
  - format
  - style
---

:::info
A standard Sysl source formatting tool is on our roadmap, which should address some of these points automatically. However it is unlikely that there will ever be exactly one right way to format everything.
:::

### Tags and Annotations

An annotation can be added to the same line as its subject, or (following a `:`) the following line with an indent. Decide as follows:

- **Tags**: same line
  ```sysl
  App [~db]:
      !table Foo:
          bar <: string [~pk]
  ```
- **Short, technical annotations** (relating directly to technical aspects of their subject): same line
  ```sysl
  App:
      !table Foo:
          bar <: string? [default="baz"]
  ```
- **Everything else**: separate line
  ```sysl
  App:
      @package = "io.sysl"
      !table Foo:
          bar <: string?
              @displayName = "Bar"
  ```
- Multiline annotations (where line breaks are expected, even if not yet present): separate line, multiline format
  ```sysl
  App:
      @description =:
          | A short summary of the application.
          |
          | A longer description that can span across multiple
          | lines.
      !table Foo:
          bar <: string
              @description =:
                  | Contains the bar.
  ```
