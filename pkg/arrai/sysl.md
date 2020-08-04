# Sysl

`sysl.arrai` contains helper functions to make working with Sysl models more convenient.

## polish

`polish` takes an imported Sysl model and returns a copy with a cleaner internal structure, making it easier to work with. It applies the following tweaks:

- Parse `return` statement payloads into `(status, type)` tuples.

## normalize

`normalize` takes an imported Sysl model and returns an equivalent relational model. The relational model is a tuple structured like a database, with tables for each Sysl concept (e.g. `apps`, `endpoints`, `types`, `calls`, etc.) that can be naturallys joined with the `<&>` operator.

The relational model is a powerful option for transformations that resemble queries with joins and filters (i.e. `where` clauses). See the diagram scripts for examples.
