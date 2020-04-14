---
title: "Import"
date: 2018-02-28T10:11:24+11:00
description: "Sysl can import files of various formats."
weight: 10
bref: "Sysl can import files of various formats."
topic: "Import"
layout: "single"
toc: true
---

Sysl can import files of various formats. Currently supported formats include

- OpenAPI 2.0
- OpenAPI 3.0
- Swagger 2.0

There are two ways of using specifications defined in other formats.

## Import the specification directly into Sysl

This lets you convert a file into sysl, and work from there. This is suitable for when you want to migrate to Sysl as the source of truth for your system specifications. When used with code generation, it offers higher levels of customization. For more details refer to [import](/docs/commands/import)

## Reference the specification file in a Sysl file

This lets you use Sysl whilst still working with specs defined in other languages, often supplied by vendors and other teams.

For more details on how to import external sysl files, refer to [non-sysl-file](/docs/language/#non-sysl-file)
