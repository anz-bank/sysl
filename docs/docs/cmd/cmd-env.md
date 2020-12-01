---
id: cmd-env
title: Environment Variables
sidebar_label: env
keywords:
  - command
---

## Summary

The `sysl env` command prints out the values of environment variables read by sysl to aid in diagnosing and reproducing reported issues.

## Usage

```
sysl env
```

## Output

### SYSL_MODULES

`SYSL_MODULES` is a flag to indicate whether Sysl modules are enabled.

### SYSL_PLANTUML

`SYSL_PLANTUML` is used to configure the address of the PlantUML server used for diagram generation
