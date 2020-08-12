# Blueprint Make

Generates a Makefile to execute the tests of downstream repos.

Blueprint files output an arr.ai tuple of the form:

```arrai
(
    name: 'name',
    repo: 'https://github.com/org/repo',
    path: 'optional/path/within/repo',
    tests: {
        'name': 'command',
    },
    # Blueprints that depend on this one.
    children: {
        (...),
        //{./path/within/repo/blueprint.arrai},
        //{github.com/org/other-repo/blueprint.arrai},
    }
)
```

The Blueprint tool performs ecosystem testing as follows:

- Locate a `blueprint.arrai` file in the current working directory
- Invoke `blueprint_make.arrai` to generate a `Makefile` to perform the testing
- Run `make` with the generated `Makefile` from the working directory of the `blueprint.arrai` file
- `make` will then clone the downstream repos into `.blueprint/children` and execute their test commands from the working directories of their `blueprint.arrai` files

```bash
arrai run blueprint_make.arrai
all: ...
...
```
