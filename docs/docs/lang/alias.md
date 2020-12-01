---
id: alias
title: Alias
keywords:
  - language
---

An alias is a redefinition of some [Type](./type.md) with a different name. This is commonly used to simplify an existing type for a different context.

## Example

```
Bank:
    !alias Accounts:
        sequence of Account
```

## See also

- [Type](./type.md)
