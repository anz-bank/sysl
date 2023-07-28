---
title: TypeScript Wrapper
---

The Sysl CLI is a convenient tool when you're working in a terminal or a Bash script, but for more involved use cases you will want to use Sysl's features programmatically.

The Sysl SDK TypeScript library is intended to cover the same space as the Sysl CLI, but it will take a while. Where there are gaps, you can call the Sysl CLI from TypeScript using the `spawnSysl` function.

There is also a TypeScript-based CLI that wraps the existing Sysl CLI in `ts/src/cli`. This CLI mostly just passes your input through to `sysl`, but it also enhances some of the functionality of the CLI commands using functionality of the Sysl SDK (which is evolving much faster than the Sysl binary). The functions for which the behaviour is altered are documented below.

# Import

**Old:** `sysl import` overwrites any existing content when writing output to disk, clobbering any changes the user may have made (hence the `DO NOT EDIT` header on imported Sysl specs).

**New:** When importing to an existing `.sysl` file, some of the existing content of that file is loaded and merged with the new content before that new content is written to disk. Specifically retained information includes:
- [statements](../lang/statement.md) in [endpoints](../lang/endpoint.md), since these are never modelled in the sources that Sysl imports
- [import](../lang/import.md) statements (deduplicated), since they won't produce conflicts and are necessary to support call [statements](../lang/statement.md).

To use this enhanced `import` command, [install the Sysl binary as usual](../installation.md), and run:

```
npx sysl import [options]
```
