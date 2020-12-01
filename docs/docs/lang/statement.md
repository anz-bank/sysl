---
id: statement
title: Statement
keywords:
  - language
---

A Statement describes some behavior of an [Endpoint](./endpoint.md).

## Syntax

Statements are nested under an Endpoint or another parent Statement. There are quite a few kinds of Statement to capture different kinds of behavior:

- `Action`: Just some plain text describing something, like `validate input`.
- `Call`: A call to an Endpoint of an [Application](./application.md) (`.` for same, or the name of another).
  - Used to [generate diagrams](../gen-diagram.md) with interactions between Applications (e.g. integration and sequence diagrams).
- `Cond`: A conditional branch (`if something`, `else if something`, `else`).
  - The `if` and `else` keywords are case-insensitive.
- `Loop`:
- `LoopN`:
- `Foreach`:
- `Alt`:
- `Group`: A collection of child Statements.
- `Return`: The return of a status and value from the Endpoint (`return status <: Type`).

## Examples

```sysl
Statements:
    IfStmt:
        if predicate1:
            # more statements
            return ok <: string
        else if predicate2:
            # more statements
            . <- IfStmt
        else:
            # more statements
            ...

    Loops:
        alt predicate:
            # more statements
            ...
        until predicate:
            # more statements
            ...
        for each predicate:
            # more statements
            ...
        for predicate:
            # more statements
            ...
        loop predicate:
            # more statements
            ...
        while predicate:
            # more statements
            ...

    Returns:
        return ok <: string
        return ok <: Types.Type
        return error <: Types.Type

    Calls:
        # self call
        . <- Returns

        # Rest endpoint call
        RestEndpoint <- GET /param

    OneOfStatements:
        one of:
            case1:
                # more statements
                return ok <: string
            case number 2:
                # more statements
                return ok <: int
            "case 3":
                # more statements
                return ok <: Types.Type
            :
                return error <: string

    GroupStatements:
        grouped:
            # more statements
            . <- GroupStatements

    Annotations:
        @annotation1 = "you can do string annotation like this"
        @annotation2 = ["or", "in", "an", "array"]
        @annotation3 =:
            | you can also do
            | multiline annotations
            | like this

    AnnotatedStatements:
        . <- Miscellanous [annotation="annotation can be added for any statement as a string"]
        return ok <: string [annotation=["as", "an", "array"]]
        "statement" [annotation=[["or", "as", "an"], ["array", "of", "arrays"]]]

    Miscellanous:
        | you can add comments like this
        "string statements"
        SimpleEndpoint -> SimpleEp # TODO: what's this for?
```

## See also

-
