---
id: endpoint
title: Endpoint
keywords:
  - language
---

An Endpoint represents a behavior of an [Application](./application.md) that can be invoked. The behavior can be described with a sequence of [Statement](./statement.md)s.

Sysl supports two kinds of Endpoint:

- RPC: A function that can be invoked remotely. Parameters are specified in the function signature.
- REST: A function mapped to a URL with a specific path and HTTP method. Parameters can be specified in three ways:
  - Path parameters: `/foo/bar/{param <: int}`. Types must be primitive.
  - Query parameters: `GET ?foo=int&bar={TypeName}`. Types may be primitive (bare) or references (wrapped in curly braces).
  - Payload parameters: `POST (foo <: string [~body], head <: TypeName [~header], bar <: int [~cookie])`. Can contain header, cookie and body content, identified by the corresponding [Tag](./tag.md)s.

## Syntax

At a minimum, an RPC endpoint is represented by its name and the placeholder Statement (`...`):

```sysl
# An Application with a "Withdraw" RPC Endpoint.
Bank:
    Withdraw:
        ...
```

The Endpoint can be followed by parentheses containing a specification of the input Parameters. Each Parameter is simply a name and a [Type](./type.md), much like a [Field](./field.md).

```sysl
Bank:
    # An RPC Endpoint that takes two inputs.
    Withdraw(accountNumber <: string, amount <: int):
        ...

    # An RPC Endpoint that explicitly takes no inputs.
    Logout():
        ...
```

REST Endpoint specifications are similar, but are described by their path with a nested HTTP method. Paths can be nested to The HTTP method can be one of `GET`, `PUT`, `POST`, `DELETE` or `PATCH`.

```sysl
Bank:
    # Get details of an account.
    /account/{accountNumber <: int}:
        # Query params can be specified on the method.
        GET ?from=date&to=date:
            ...

        # This path builds on the parent path above, withdrawing from that account.
        /withdraw:
            POST (Transaction):
                ...

    # The schema of the POST request body.
    !type Transaction:
        amount <: int
```

## See also

- [Application](./application.md): parent of an Endpoint
- [Statement](./statement.md): children of an Endpoint
- [Field](./field.md)
