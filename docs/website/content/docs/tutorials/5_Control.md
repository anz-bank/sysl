---
title: "Control"
date: 2019
weight: 50
draft: true
bref: ""
toc: false
---

Sysl allows you to express high level of detail about your design. You can specify decisions, processing loops etc.

##### If, else
You can express an endpoint's critical decisions using IF/ELSE statement:
```
Server:
  HandleFormSubmit:
    validate input
    IF session exists:
      use existing session
    Else:
      create new session
    process input
```
See [/assets/if-else.sysl](/assets/if-else.sysl) for complete example.

`IF` and `ELSE` keywords are case-insensitive. Here is how sysl will render these statements:

![](/assets/if-else-Seq.png)

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
