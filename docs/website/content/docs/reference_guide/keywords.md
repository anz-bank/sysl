---
title: "Keywords"
description: "Sysl keyword dictionary"
date: 2018-11-26
weight: 5
draft: false
toc: true
sequence: true
---
## Type
`!type`
The type keyword is used to define a type. 
In the following example we define a `Post` type made up of multiple attributes.
```
  !type Post:
    userId <: int
    id <: int
    title <: string
    body <: string
```

## Alias
`!alias`
Alias' can be used to simplify a type;

```
  !alias Posts:
    sequence of Post
```
## View
`!view`
Views are sysl's functions; we can use them in the transformation language, see [docs/transformation.html]for more info

## Union
`!union`
Unions are a union type; 
`!union string, int32`
can either be a string, int32, but not both.

## Table
`!table` Add more here

## Wrap
`!wrap` Add more here



