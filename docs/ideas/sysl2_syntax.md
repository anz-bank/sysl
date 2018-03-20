sysl2 Syntax
============
The following is a collection of unrelated, unevaluated ideas on how to improve the sysl Syntax:

* Indent vs `{}`
* Whitespace free identifiers
* `.` instead `<-` (for Method invocation)
* Rethink `::` namespaces
* Eliminate `<:` e.g. `Counter <: int` -> `Counter int`
* Revisit transformation language change

```
      X -> (:
        x = .a
        y = 2 + .b
     )

     X -> {
     	x: .a
     	b: 2 + .b
     }
```
* `!table` `!type` `!wrap` `!view` get rid of bang `!`
* "Scoping"  - reference syntax needs to be rationalized
* Introduce Top-level, "intent" keywords, e.g. `model`, `app`, `project` (?), etc.
* Get rid of `project` structure for diagrams.
* Single word keywords instead of `set of` use `setOf` or `set()`
* Meta data on types in brackets `[~pk, ~autoinc]` and part of typedef: `int?` (`?` means optional)
* Inconsistent naming in `reljam` commandline vs `*.sysl` file: `!wrap` vs `facade`
* Find a better term for `!type` maybe `!object` (invokes association with JSON)
