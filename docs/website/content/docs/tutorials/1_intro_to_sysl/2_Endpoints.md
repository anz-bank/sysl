---
title: "Endpoints"
date: 2019
weight: 2
sequence: true
toc: true
---

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

One 
-----------

Two
-----------