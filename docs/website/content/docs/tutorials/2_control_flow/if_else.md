---
title: "If-Else"
date: 2019
weight: 50
sequence: true
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