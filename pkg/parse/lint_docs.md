# Linter Design

Linting happens during parsing of the abstract syntax tree. The linter itself
is a part of the `TreeShapeListener` struct in the `linter` field. The `linter`
field's job is to collect data that is required to do its linting operation.
The linting operation can happen during walking the abstract syntax tree or
after that process. This depends on what the linter is detecting.

Table of Contents:
- [Linter Design](#linter-design)
  - [Linter Format](#linter-format)
  - [Return Statement Lint](#return-statement-lint)
  - [Case-sensitive Applications Redefinition](#case-sensitive-applications-redefinition)
  - [Call Statement Linting](#call-statement-linting)
  - [Additional Notes](#additional-notes)

## Linter Format
The format of the linter message is the following:
```
lint path/to/file:lineNumber:colNumber: linter message
```

But there are linter messages that involve multiple locations.
For that messages, the format is the following:
```
lint: linter message:
path/to/file1:lineNumber:colNumber
path/to/file2:lineNumber:colNumber
path/to/file3:lineNumber:colNumber
...
```

## Return Statement Lint

For this feature, it is easier to lint during parsing. A simple regex match in
the `EnterRet_stmt` function is done to ensure it uses the `ok <: some_type`
format.

## Case-sensitive Applications Redefinition

For this feature, it is not possible to lint during parsing as it needs to walk
through all the modules used. This is where the `linter` field comes in. During
the walk operation, the `linter` field will collect all the required data. The
following is the definition for `linterRecords` which is the type for the
`linter` field.

```go
type linterRecords struct {
    apps  map[string]*graph
    calls *graph
}

type graph map[string]*graphData

type graphData struct {
    locations map[string]bool
    rec       *graph
}
```

For this feature, the linting mainly uses the `linterRecords.apps` field. The
data structure is a `map[string]*graph`, the key is an application name that is
lowercased and its value is a map of applications whose names correspond to it.

When the walker reaches `EnterApp_decl`, `*TreeShapeListener.recordApp` is
called. It checks if the lowercased application name is ever recorded. If it is,
it will record it as another entry in the corresponding graph. If it is not
recorded yet, it will create a new entry in `linterRecords.apps` and assign
a new graph for it. The implementation of the linting can be found in the
function `*TreeShapeListener.lintAppDefs`

## Call Statement Linting

This feature also requires the linting to be done after the walk operations. For
this linting, it needs to record two types of data. The defined endpoints for
each application and all the endpoint calls in all the imported modules. It uses
the same `linterRecords` struct as above to do this.

This is where the `graph` data structure comes. It is a rudimentary graph, the
structure for the `apps` and `calls` field is the same. The root (key) is an
application name, the corresponding value `*graphData` contains locations of the
application definitions stored in the `locations` field and the endpoints in the
`rec` field. The endpoints `rec` field is another `graph` whose keys are the
name of the endpoints and the corresponding `*graphData` contains locations and
another graph which contains REST methods (GET, POST, etc). However, only REST
endpoints contain the REST methods graph, if the `rec` field of endpoints is nil,
the endpoint is a simple endpoint.

The `locations` field is defined as a `map[string]bool` because the same data can
be defined in multiple locations (e.g. endpoint calls `App <- Endpoint`).

The `apps` field in `linterRecords` contains all the defined endpoints for each
applications in the sysl modules. On the other hand, the `calls` field records
all the call statements in the sysl modules. For optimisation purposes, the
location data for each call statement is only stored at the leaf of the graph.
So for REST endpoints, they are stored at the method graph and for simple
endpoints, they are stored in the endpoint graph.

The linting implementation for this feature can be found in the function
`*TreeShapeListener.lintEndpoint`. The implementation just checks that all the
calls in the `calls` field exist in the `apps` field. It checks that
applications, endpoints, and methods (for REST endpoints) are defined.

## Additional Notes

The `recordApp`, `recordEndpoint`, and `recordMethod` functions
of the struct `*graph` return errors when a key already exists
(or doesn't exist) in the graph. The errors are meant to just
catch implementation mistakes and it should not happen, unless
there's a bug.
