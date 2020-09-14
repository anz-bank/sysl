# Data Flow Analysis Tool Design

## Objective

Create a tool to track a type or field in a sysl spec how it is being sent
throughout all the related endpoints.

## Background

It is important to be able to track where data is going throughout all the
systems as part of data flow analysis.

## Summary

Have not decided on one yet

## Detailed Design

Have not decided on one yet

## Previous Solutions

### First Solution
<a name="first-solution"></a>

The first solution was to create a graph out of the protobuf of a sysl
specification. The script `graph.arrai` takes the encoded json sysl protobuf
and creates a data structure that represents the parsed sysl module. The
`schema.arrai` represents what the output of `graph.arrai` looks like.

After creating a graph, the script with then try to traverse through the graph.
The script `sources.arrai` finds endpoints that returns a target type. If the
target type is a type, then it searches for any endpoints that return that type.
But if the target type is a field, it will look for any endpoints that return
that field and the parent type of that field.

After finding all the sources, it starts to traverse through all the sources.
It does breadth-first search on the sources and create a set of edges. So it
uses the list of sources and find the next sources recursively until it does not
find anymore sources. An edge between two nodes is collected if one of the node
calls the other node and the other node returns a target type or a transformed
version of that target type.

In this iteration, the solution handles transformation tracking by going
allowing a special attribute `map_of` in the type or alias definition of the
sysl specification. For example:

```sysl
App:
    !alias FullName:
        @map_of = ["App.Name.FirstName", "App.Name.LastName"]
        string

    !type Name:
        FirstName <: string
        LastName <: string

    GetName:
        return ok <: App.Name

    GetFullName:
        . <- GetName
        return ok <: App.FullName
```

The above example shows that the type `App.FullName` is a result of
transformation from the fields `App.Name.FirstName` and `App.Name.LastName`.

## Alternative Solutions

Several alternative solutions are being considered in creating the tool.

### Simpler Graph

The graph of [Solution 1](##first-solution) is easy to query through but its
nested form can create problem as arrai is better at handling flatter data
structure. In this solution, the result of `graph.arrai` is changed to be
a set of edges. Each edge is tuple with the following format:

```arrai
(
    from   : (app: appNameInString, endpoint: endpointNameInString),
    to     : (app: appNameInString, endpoint: endpointNameInString),
    returns: {typesInString},
    map_of : {typesInString}
)
```

The `returns` field stores all the types of the endpoint listed in `from`
returns. The `map_of` attribute is the union of all the `map_of` types for
each return types of that endpoint.

#### Why it didn't work

This was designed with the `map_of` attribute in mind. Since the types name in
the `returns` and `map_of` fields will be consistent throughout all endpoints,
the edges will be easy to traverse through. But the `map_of` attribute is very
limiting as it is a special attribute in the type definition of the sysl
specification. This is very limiting as each endpoint can transform types
differently. This makes the `map_of` attribute not the most ideal solution.

### Even Simpler Graph
<a name="even-simpler-graph"></a>

The `map_of` attribute is being replaced by the `dataflow` attribute. The
`dataflow` attribute is similar to the `map_of` attribute but it is being
placed in the endpoint instead of type definition.

The `dataflow` attribute can be used like this:

```sysl
App:
    !type Type1: ...
    !type Type2: ...
    !type Type3: ...

    Endpoint:
        @dataflow = [[["App.Type1", "App.Type2"], "App.Type3"]]
        return ok <: App.Type3
```

The above sysl specification states that `App.Type1` and `App.Type2` are
transformed into `App.Type3` in the endpoint `App.Endpoint`.

To create an even flatter graph, this solution change the result of
`graph.arrai` to be a set of edges which is in the form of tuples in the
following format:

```arrai
(
    caller: (
        app     : appNameInString,
        endpoint: endpointNameInString,
        type    : typeInString
    ),
    callee: (
        app     : appNameInString,
        endpoint: endpointNameInString,
        type    : typeInString
    )
)
```

This format should work and it is actually the currently used data structure.
The set of edges is flat but it does create a lot of edges as each endpoint can
return multiple types.

#### Problems

The `dataflow` attribute in sysl specification does not limit the way user can
show transformation. But this creates a problem in the implementation as this
means transformation is contextual, it is no longer constant for every type and
it depends on each endpoint itself.

Creating the set of edges should not be too problematic but traversing the
edges will be complicated. It has to do depth-first search in tracking how the
target type is being transferred throughout the systems. This will be a very
massive task for the script, especially on large sysl specifications. This can
be considered the brute force method. Optimization was considered through
memoization but other ways are still being considered.

An optimization was considered during the edges collection by only collecting
edges that return a target type or the transformed version of it. But this seems
to complicate the task due to deep transformation. For example:

```sysl
App:
    !type Type1: ...
    !type Type2: ...
    !type Type3: ...
    !type Type4: ...
    !type Type5: ...

    Endpoint1:
        return ok <: Type1

    Endpoint2:
        @dataflow = [[["App.Type1"], "App.Type2"]]
        . <- App.Endpoint1
        return ok <: App.Type

    Endpoint3:
        @dataflow = [[["App.Type2"], "App.Type3"]]
        . <- App.Endpoint2
        return ok <: App.Type2

    Endpoint4:
        @dataflow = [[["App.Type3"], "App.Type4"]]
        . <- App.Endpoint3
        return ok <: App.Type4

    Endpoint5:
        @dataflow = [[["App.Type4"], "App.Type5"]]
        . <- App.Endpoint4
        return ok <: App.Type5
```

Let's say the script is tracking `App.Type1`, during the edges creation in
`App.Endpoint5`, it will check for `App.Type1` and it will not find anything. It
will then have to traverse to `App.Endpoint4` due to the call statement and
repeat until it find `App.Endpoint2`. As types can be transformed many times,
this will create a deep traversal. So the edges creation by filtering might not
be the most ideal. There might be optimizations through memoization but I have
not thought of any way yet.

### Simpler Dataflow Attribute

The `dataflow` attribute can get verbose in a complicated endpoint so a
different form is being considered. Currently, the return statement grammar is
very relax and will allow most things. The form considered here is the
following:

```sysl
App:
    !type Type1: ...
    !type Type2: ...
    !type Type3: ...

    Endpoint:
        return ok <: SomeType [~App.Type1, ~App.Type2, ~App.Type3]
```

The above sysl shows that the type `SomeType` in `App.Endpoint` is a transform
of `App.Type1`, `App.Type2`, and `App.Type3`.

#### Problem

The problems are mostly the same with the previous solution but also, this will
require grammar and sysl compiler changes.

### Simpler (Naive) Implementation

Another simpler implementation was considered. This solution does not require
any special attributes in the sysl specification. This solution however does
not track transformation as accurate as the previous solutions. Basically, the
script will does the same thing by creating a set of edges just like in
[Even Simpler Graph](#even-simpler-graph) but it makes the assumption that if
any endpoints that return anything, the return type will be considered the
transformed version of whatever types that endpoint has at any time. For
example:

```sysl
App:
    !type Type1: ...
    !type Type2: ...

    Endpoint1:
        return ok <: App.Type1

    Endpoint2:
        . <- App.Endpoint1
        return ok <: App.Type2
```

It assumes that, in `App.Endpoint2`, the return type is a transformed version
of `App.Type1`. This might not necessarily be true as the implementation might
use `App.Type1` differently and not transform it at all.

#### Problem

The problem here is that this solution might create a lot of false positive and
therefore not the most ideal solution.
