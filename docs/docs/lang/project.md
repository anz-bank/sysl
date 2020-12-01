---
id: project
title: Project
keywords:
  - language
---

When modelling a large-scale system, generated outputs from the entire model will often produce far too much content. Some outputs can be generated for a single element of the model (e.g. a service or database), but the most interesting artifacts describe the overlaps and interactions between parts of the system.

In Sysl, a Project is a _projection_ of a subset of a larger model. It is a means of selecting a few elements of interest from the system, and generate outputs just for them.

A Project is not a fundamental concept in Sysl; it is merely an interpretation of an [Application](./application.md). A Project is an Application with a sequence of names as children (instead of Endpoints), and each child having a sequence of Application names (instead of Statements). Using a Project as an input is an instruction to a tool to include the Applications named in the Statements.

TODO: More detail.

## See also

- [Application](./application.md)
- [Diagram Generation](../gen-diagram.md)
