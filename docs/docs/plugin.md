---
id: plugin
title: Plugin
sidebar_label: Plugin
---

:::caution
WIP

**TODO:**

- Update and polish content.
- Move referenced assets to a permanent directory on GitHub and update links.

:::

---

## Extending Sysl

In order to easily reuse and extend Sysl across systems, the Sysl compiler translates Sysl files into an intermediate representation expressed as protocol buffer messages. These protobuf messages can be consumed in your favorite programming language and transformed to your desired output. In this way you are creating your own Sysl exporter.

Using the protoc compiler you can translate the definition file of the intermediate representation pkg/proto/sysl.proto into your preferred programming language in a one-off step or on every build. You can then easily consume Sysl models in your programming language of choice in a type-safe way without having to write a ton of mapping boilerplate. With that you can create your own tailored output diagrams, source code, views, integrations or other desired outputs.

## Status

Sysl is currently targeted at early adopters. The current focus is to improve documentation and usability, especially error messages and warnings.
