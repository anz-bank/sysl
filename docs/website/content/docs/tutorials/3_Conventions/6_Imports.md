---
title: "Imports"
date: 2019
weight: 50
sequence: true
---

To keep things modular, sysl allows you to import definitions created in other `.sysl` files.

E.g. `server.sysl`
```
Server:
  Login: ...
  Register: ...
```

and you use `import` in `client.sysl`

```js
import server

Client:
  Login:
    Server <- Login
```

Above code assumes, server and client files are in the same directory. If they are in different directories, you must have atleast a common root directory and `import /path/from/root`.

All sysl commands accept `--root` argument. Run `sysl -h` or `reljam -h` for more details.
