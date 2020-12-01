---
id: comment
title: Comment
keywords:
  - language
---

A Comment is some text in a Sysl source file that is ignored by the parser. This behavior is the same as comments in other languages.

A Comment begins with a `#` and continues to end of the line.

## Example

```sysl
# Comments on the first line of a file describe the file.

# This is a comment on an Application.
# To span multiple lines, just start each with a #.
App:
    Endpoint: # Comments can be inline as well.
        ...
```
