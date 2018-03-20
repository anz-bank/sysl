---
title: "Output formats"
description: "Explore various output formats from sequence diagrams to Java code."
date: 2018-02-28T14:05:31+11:00
weight: 60
draft: false
bref: "Explore various output formats from diagrams to code"
toc: false
---

Sysl consists of two executables: `sysl` and `reljam`. The **Sy**stem **S**pecification **L**anguage `sysl` is mainly concerned with diagram creation whereas the **Rel**ational **Ja**va **M**odels creator `reljam` generates different types of source code output.

Sysl outputs
------------
| Command | Description |
|---------|-------------|
| data    | Data model diagrams |
| ints    | Integration Diagrams |
| sd      | Sequence Diagrams |
| pb      | Binary Protocol Buffer files of the Sysl definitions (plugins)    |
| textpb  | Text based Protocol Buffer files of the Sysl definitions (plugins, debugging) |


Reljam outputs
--------------
| Command | Description |
|---------|-------------|
| model   | Java model implementation (in memory) |
| facade  | Java facade implementation (restricted access to creating and populating models) |
| view    | Java implementation of sysl model transformatios|
| xsd     | XSD represtation of sysl model |
| swagger | Swagger representation of REST APIs and models |
| spring-rest-service | Java Spring REST API implementation |

Sysl diagram samples
-------------------
[[TODO]]
