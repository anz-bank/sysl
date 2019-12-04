---
title: "Data"
date: 2019
weight: 3
sequence: true
---

You will have various kinds of data passing through your systems. Sysl allows you to express ownership, information classification and other attributes of your data in one place.

Continuing with the previous example, let's define a `Server` that expects `LoginData` for the `Login` Flow.
```
Server:
  Login (request <: Server.LoginData): ...

  !type LoginData:
    username <: string
    password <: string
```
In the above example, we have defined another application called `Server` that has an endpoint called `Login`. It also defines a new data type called `LoginData` that it expects callers to pass in the login call.

Notes about sysl syntax:
  * `<:` is used to define the arguments to `Login` endpoint.
  * `!type` is used to define a new data type `LoginData`.
    * Note the indent to create fields `username` and `password` of type `string`.
    * See [Data Types](#data-types) to see what all the supported data types.
  * Data types (like `LoginData`) belong to the app under which it is defined.
  * Refer to the newly defined type by its fully qualified name. e.g. `Server.LoginData`.

#### Data Types
Sysl supports following data types out of the box.
  * int, int64, int32
  * float, decimal
  * string
  * bool
  * datetime, date
  * any
  * xml


Now, we have two apps `MobileApp` and `Server`, but they do not interact with each other. Time to add some statements.


#### Return response
An endpoint can return response to the caller. Everything after `return` keyword till the end-of-line is considered response payload. You can have:
  * string - a description of what is returned, or
  * Sysl type - formal type to return to the caller

Sequence diagram will render the response accordingly. In the previous example, `data` is a generic description of what `DB <- Query` returns. `Server.Response` is the Sysl type that is returned by `Login` endpoint.
