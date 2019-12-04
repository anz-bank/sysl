---
title: "Conventions"
date: 2019
weight: 50
sequence: true
---

#### Multiple Declarations
Sysl allows you to define an application in multiple places. There is no redefinition error in sysl.

```
UserService:
  Login: ...

UserService:
  Register: ...
```

Result will be as-if it was declared like so:

```
UserService:
  Login: ...
  Register: ...
```

### Projects
Most of the changes to your system will be done as part of a well defined `project` or a `software release`.

`TODO: Elaborate`