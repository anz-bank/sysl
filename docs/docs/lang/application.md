---
id: application
title: Application
keywords:
  - language
---

Applications are the top-level elements of the Sysl model.

## Name

Each Application has a name, but it's more than just a simple string. The Application's full name is actually an array of parts separated by `::` in Sysl source. This full name must be used when referring to the Application in other contexts (e.g. [Type](./type.md) references and Call [Statement](./statement.md)s).

The parts of the name preceding the final part are commonly known as the Namespace. This is analogous to the package of a fully-qualified Java class. The final part may be used by itself in some outputs generated from the model, but fundamentally it is no more special than the other parts.

Applications may also have a "long name" that, like a display name, can be more descriptive, and will be used to represent the Application in various outputs. The long name follows the Application name in double quotes, like `App "Long Name"`.

For guidance on choosing names, see [Identifiers](./identifiers.md) and [Best Practices](../best-practices/intro.md).

### Examples

```sysl
# Application with simple name.
App:
    ...

# Application with a namespace.
Namespace :: App:
    ...

# Application with a multi-part namespace, common for large systems.
Organization :: Team :: System :: App:
    ...

# Application with a more descriptive "long name".
GCP "Google Cloud Platform":
    ...
```

## See also

Children of Application:

- [Endpoint](./endpoint.md)
- [Type](./type.md)
- [Table](./table.md)
- [Enum](./enum.md)
- [Union](./union.md)
- [Alias](./alias.md)
- [Annotation](./annotation.md)
- [Tag](./tag.md)
- [View](./view.md)
- [Mixin](./mixin.md)

* [Project](./project.md): a special interpretation of an Application
