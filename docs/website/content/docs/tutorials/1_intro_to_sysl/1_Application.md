---
title: "Applications"
date: 2019
weight: 1
sequence: true
toc: true
---
## Concepts
Sysl allows you to specify [Application](#application) behaviour and [Data Model](#data) that is shared between your applications. Another related concept is of software [Projects](#projects) where you can document what changes happened in each project or a release.

To explain these concepts, we will design an application called `MobileApp` which interacts with another application called `Server`.

### Application
An application is an independent entity that provides services via its various `endpoints`.

Here is how an application is defined in sysl.
```
MobileApp:
  ...
```
`MobileApp` is a user-defined Application that does not have any endpoints yet. We will design this app as we move along.

Notes about sysl syntax:

  * `:` and `...` have special meaning. `:` followed by an `indent` is used to create a parent-child relationship.
    * All lines after `:` should be indented. The only exception to this rule is when you want to use the shortcut `...`.
    * The `...` (aka shortcut) means that we don't have enough details yet to describe how this endpoint behaves. Sysl allows you to take an iterative approach in documenting the behaviour. You add more as you know more.

The next bit is to add endpoints to this app.
