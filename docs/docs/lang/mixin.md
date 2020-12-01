---
id: mixin
title: Mixin
keywords:
  - language
---

A Mixin merges the content of another [Application](./application.md) into the parent Application.

## Example

In the following example, the definition of a `BankTeller` contains a Mixin of the `ATM`. This includes all of the children of `ATM` in `BankTeller`, implying (roughly) that a `BankTeller` can do everything that an ATM can do.

```sysl
ATM:
    Withdraw:
        ...

BankTeller:
    -|> ATM

    Deposit:
        ...
```

## See also

- [Application](./application.md): parent element
