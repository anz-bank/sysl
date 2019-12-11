---
title: "Conventions"
date: 2019
weight: 3
sequence: true
toc: true
---

## Multiple Declarations

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

## Projects
Most of the changes to your system will be done as part of a well defined `project` or a `software release`.

`TODO: Elaborate`

## Imports


To keep things modular, sysl allows you to import definitions created in other `.sysl` files.

E.g. `server.sysl`
```
Server:
  Login: ...
  Register: ...
```

and you use `import` in `client.sysl`

```
import server

Client:
  Login:
    Server <- Login
```

Above code assumes, server and client files are in the same directory. If they are in different directories, you must have atleast a common root directory and `import /path/from/root`.

All sysl commands accept `--root` argument. Run `sysl -h` or `reljam -h` for more details.

