---
title: "Iteration"
date: 2019
weight: 50
sequence: true
---

##### For, Loop, Until, While
Express processing loop using FOR:
```
Server:
  HandleFormSubmit:
    validate input
    For each element in input:
      process element
```
See [/assets/for-loop.sysl](/assets/for-loop.sysl) for complete example.

`FOR` keyword is case insensitive. Here is how sysl will render these statements:

![](/assets/for-loop-Seq.png)

You can use `Loop`, `While`, `Until`, `Loop-N` as well (all case-insensitive).
