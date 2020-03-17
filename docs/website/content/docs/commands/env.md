---
title: "env"
description: "Print sysl environment information."
date: 2018-02-28T14:05:40+11:00
weight: 70
draft: false
bref: "Print sysl environment information."
toc: true
---

The `sysl env` command prints out the values of environment variables read by sysl to aid in diagnosing and reproducing reported issues.

## Usage

```bash
usage: sysl env
```

## Output

### SYSL_MODULES

`SYSL_MODULES` is a flag to indicate whether Sysl modules are enabled.

### SYSL_PLANTUML

`SYSL_PLANTUML` is used to configure the address of the PlantUML server used for diagram generation
