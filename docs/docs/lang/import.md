---
id: import
title: Import
keywords:
  - language
---

Each Sysl file defines a [Module](./module.md). However a large system defined in a single file is hard to manage. The `import` keyword allows one file to include the contents of another file, merging the two modules together. This allows a large system to be split across many simple files.

An `import` statement merges the contents of another Sysl [Module](./module.md) into the current one.

For example, given some Sysl file `server.sysl`:

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

Once a Module is imported, it doesn't matter where it came from. It can be local or remote, with any directory structure or filename.

Imported files can even be non-Sysl files, if Sysl knows how to import them. For example, you can `import` a `.json` file containing an OpenAPI spec, and it will be automatically converted to Sysl.

If the imported file is a `.sysl` file, you can omit the `.sysl` extension from the import location.

### Local relative file

Imagine that you have `server.sysl`, `client.sysl` and `deps.sysl` files structured as below:

```
.
├── server.sysl
├── client.sysl
└── deps
    └── deps.sysl
```

To import both `server.sysl` and `deps.sysl` files in `client.sysl`:

```sysl title="client.sysl"

import server
import deps/deps

Client:
  Login:
    Server <- Login
  Dep:
  	Deps <- Dep
```

```sysl title="server.sysl"
Server:
  Login: ...
  Register: ...
```

```sysl title="deps.sysl"
Deps:
	Dep: ...
```

The resulting model in `client.sysl` will include three apps: `Server`, `Client` and `Deps`.

### Local absolute file

If the imported are in the same project but outside of current folder, you must have at least a common root directory and `import /path/from/root`.

All Sysl commands accept the `--root` option to specify this root directory from which paths are resolved. Run `sysl help` for more details.

```
<root-dir>
├── clients
│   └── client.sysl
└── servers
    └── server.sysl
```

```sysl title="<root-dir>/servers/server.sysl"
Server:
  Login: ...
  Register: ...
```

```sysl title="<root-dir>/clients/client.sysl"
import /servers/server

Client:
  Login:
    Server <- Login
```

### Remote file

Sysl supports importing Sysl files from any [Git](https://git-scm.com/) repository.

The location for a remote import starts with `//`, and omits the `blob/main/` path segment you may see in GitHub. For example, to import the [Sysl bank demo](https://github.com/anz-bank/sysl/blob/master/demo/bank/bank.sysl), you would `import //github.com/anz-bank/sysl/demo/bank/bank.sysl`

```sysl title="servers/server.sysl in repo github.com/foo/bar"
Server:
  Login: ...
  Register: ...
```

```sysl title="client.sysl in Git repo github.com/your/repo"
import //github.com/foo/bar/servers/server

Client:
  Login:
    Server <- Login
```

When importing remote files from private repos, Sysl supports multiple authentication methods.

**SSH Agent**

If you have an SSH agent to connect to GitHub, Sysl will use that to fetch remote files.

For GitHub repositories, [generate a new SSH key](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent#generating-a-new-ssh-key) and [add it to your GitHub account](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account) if you haven't yet. Also make sure it's [added to `ssh-agent`](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent#adding-your-ssh-key-to-the-ssh-agent).

**Tokens**

Git repo hosting services provide access to individual files via API if you provide an API token. If you set the `SYSL_TOKENS` environment variable, Sysl will use the token you specify for the domain of the remote import.
On GitHub, this token is a [personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token).
The `SYSL_TOKENS` value is a comma-separated list of `domain:token` pairs. For example:

```
export SYSL_TOKENS="github.com:<GITHUB-PAT>,gitlab.com:<GITLAB-TOKEN>"
```

**SSH keys**

Sysl can connect directly to GitHub using [the SSH private key you've registered with GitHub](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent). Tell it where to find your key with the `SYSL_SSH_KEYS` environment variable.
`SYSL_SSH_KEYS` is a comma-separated list of `domain:key-path:key-passphrase` triples. If your key doesn't have a passphrase, you can leave that part empty. For example:

```
export SYSL_SSH_KEYS="github.com:<github-key-path>:<github-key-passphrase>,gitlab.com:<gitlab-key-path>:<gitlab-key-passphrase>"
```

### Non-Sysl file

To import a non-Sysl file, include the file extension and use `as AppName` to give a name to the imported application.
For example, import a [Swagger](https://swagger.io/) file like so:

```sysl
import foreign_import_swagger.yaml as Foreign :: App
```

Refer to the [`sysl import`](https://sysl.io/docs/cmd/cmd-import/#usage) command for all supported non-Sysl types. (Note that not all of them can be imported with the `import` keyword.)

### Compiled Sysl file

In some cases, you may want to import a pre-compiled Sysl model from a `.pb` (or `.textpb`) file using an `import` statement in a Sysl source file. You can do this by simply importing the file (without the `as AppName` suffix, since all the apps within will already have names):

```sysl
import foreign_import_pb.pb
import foreign_import_textpb.textpb
```

This is generally discouraged as a regular modelling practice, but may be useful for debugging or experimenting with Sysl models. For example, in order to test some new Sysl use case, you might want to add a few new annotations to a large, precompiled `.pb` model without recompiling the whole thing from scratch.

## See also

- [Module](./module.md)
