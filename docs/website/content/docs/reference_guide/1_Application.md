---
title: "First System In Sysl"
date: 2019
weight: 1
sequence: true
toc: true
---
Sysl allows you to specify Application behaviour and Data Models that are shared between your applications. This is useful in many use cases, especially in software projects with many moving parts which need their documentation kept up to date.

To explain these concepts, we will design an application called `MobileApp` which interacts with another application called `Server`.

Applications
====

An __application__ is an independent entity that provides services via its various __endpoints__.

Here is how an application is defined in sysl.
```
MobileApp:
  ...

Server:
  ...
```
`MobileApp` and `Server` are user-defined Applications that do not have any endpoints yet. We will design this app as we move along.

Notes about sysl syntax:

  * `:` and `...` have special meaning. `:` followed by an `indent` is used to create a parent-child relationship.
    * All lines after `:` should be indented. The only exception to this rule is when you want to use the shortcut `...`.
    * The `...` (aka shortcut) means that we don't have enough details yet to describe how this endpoint behaves. Sysl allows you to take an iterative approach in documenting the behaviour. You add more as you know more.

Endpoints
===

Endpoints are the services that an application offers. Let's add endpoints to our `MobileApp`.
```
MobileApp:
  Login: ...
  Search: ...
  Order: ...
```
Now, our `MobileApp` has three `endpoints`: `Login`, `Search` and `Orders`.

Notes about sysl syntax:

 * Again, `...` is used to show we don't have enough details yet about each endpoint.
 * All endpoints are indented. Use a `tab` or `spaces` to indent.
 * These endpoints can also be REST api's. See section on [Rest](#rest) below on how to define rest api endpoints.

Each endpoint should have statements that describe its behaviour. Before that lets took at data types and how it can used in sysl.

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

Data Types
==== 

Sysl supports following data types out of the box.
```
double 
int64 
float64 
string 
bool 
date.Date 
time.Time 
```
We can Also define our own datatypes using the `!type` keyword within an application.

```
!type response:
  data <: string
  type <: int
```

Now, we have two apps `MobileApp` and `Server`, but they do not interact with each other. Time to add some statements.


#### Return response
An endpoint can return response to the caller. Everything after `return` keyword till the end-of-line is considered response payload. You can have:
  * string - a description of what is returned, or
  * Sysl type - formal type to return to the caller

```
MobileApp:
  Login: 
    return string
  Search: 
    return string
  Order: 
    return sequence of string
```

