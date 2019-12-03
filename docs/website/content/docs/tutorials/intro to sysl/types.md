this
- Types are either `primitive` or composite types
- Primitives are:
    - decimals
    - int
    - float
    - string
    - bool
    - date
    - datetime


Custom types:
```
!type Post:
    userId <: int
    id <: int
    title <: string
    body <: string

!alias Posts: sequence of Post

!alias Foo: Post

```
Struct embedding:
```
!type ResourceNotFoundError:
    status <: string
```

Required and optional fields:
```
!type Account:
    status <: string
    productCode <: string
    name <: string
    limits <: sequence of Limit?
```

Package declerations:
```
AppName [package="app.name.com"]:
```