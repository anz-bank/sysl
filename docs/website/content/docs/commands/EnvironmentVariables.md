---
title: "Environment Variables"
description: "Environment Variables"
date: 2018-02-28T14:05:40+11:00
weight: 2
draft: false
bref: "Environment Variables"
toc: true
---

Several commands require environment variables to be set before they are able to correctly work.

## SYSL_PLANTUML

`SYSL_PLANTUML` is a flag to indicate the PlantUML server address. By default, it is `http://www.plantuml.com/plantuml`.

For more details, refer to [Install](/docs/install/)

## SYSL_MODULES

`SYSL_MODULES` is a flag to indicate whether Sysl modules are enabled. By default, if this is not declared, Sysl modules are enabled.
To disable Sysl Modules, set the environment variable `SYSL_MODULES=off`.
