---
id: import
title: Import
keywords:
  - language
---

An Import merges the contents of another Sysl [Module](./module.md) into the current one.

Once a Module is imported, it doesn't matter where it came from. It can be local or remote, with any directory structure or filename. Only the elements in the Module are included.

## Example

Given some `server.sysl`:

```
Server:
  Login: ...
  Register: ...
```

... and some `client.sysl` with an `import`:

```
import server

Client:
  Login:
    Server <- Login
```

The resulting model in `client.sysl` will include both the `Server` and `Client` apps.

TODO:

- `root` CLI arg
- Relative, aboslute and external file imports
- Non-Sysl imports
- GOP

## See also

- [Module](./module.md)
