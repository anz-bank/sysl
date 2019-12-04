---
title: "Applications"
date: 2019
weight: 1
sequence: true
toc: true
---
Sysl allows you to specify Application behaviour and Data Models that are shared between your applications. This is useful in many use cases, expecially in software Projects where documentation of what changes occur in each release.

To explain these concepts, we will design an application called `MobileApp` which interacts with another application called `Server`.

### Applications

An __application__ is an independent entity that provides services via its various __endpoints__.

Here is how an application is defined in sysl.
```
MobileApp:
  ...

Server:
  ...
```
`MobileApp` and `Server are user-defined Applications that does not have any endpoints yet. We will design this app as we move along.

Notes about sysl syntax:

  * `:` and `...` have special meaning. `:` followed by an `indent` is used to create a parent-child relationship.
    * All lines after `:` should be indented. The only exception to this rule is when you want to use the shortcut `...`.
    * The `...` (aka shortcut) means that we don't have enough details yet to describe how this endpoint behaves. Sysl allows you to take an iterative approach in documenting the behaviour. You add more as you know more.

This is the most simple 
