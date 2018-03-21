---
title: "Output formats"
description: "Explore various output formats from sequence diagrams to Java code."
date: 2018-02-28T14:05:31+11:00
weight: 60
draft: false
bref: "Explore various output formats from diagrams to code"
toc: false
---

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
| view    | Java implementation of Sysl model transformations|
| xsd     | XSD representation of Sysl model |
| swagger | Swagger representation of REST APIs and models |
| spring-rest-service | Java Spring REST API implementation |
