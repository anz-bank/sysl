# Sysl LSP Language Server Implementation

(forked from https://github.com/golang/tools/tree/378b9e1d59e2352276eff57e8153a4ff4053d8a7/internal/lsp)

`gopls` is a sophisticated [LSP](https://microsoft.github.io/language-server-protocol/) implementation providing language features for Golang. It's a large project that contains a combination of business logic, plumbing and boilerplate. Unfortunately it's all `internal` and can't be used as a library.

This Sysl LSP language server implementation copies the boilerplate from `gopls`, removes most of the plumbing for unnecessary features (e.g. remote servers, telemetry), and focuses on the much simpler Sysl language business logic.

## Usage

For simple distribution, the Sysl LSP language server is built into the Sysl binary. It can be invoked via stdio with the command:

```
sysl lsp
```

## Design

* The LSP interface and data structures are [transpiled from TypeScript](framework/lsp/protocol) to Golang using the TypeScript repository https://github.com/microsoft/vscode-languageserver-node.
* [`server.go`](impl/server.go) declares the main `Server` struct, to which all LSP handler functions are added in [`server_gen.go`](impl/server_gen.go).
* [lsp/helper/helper.go](framework/lsp/helper/helper.go) generates the stub implementations in [`server_gen.go`](impl/server_gen.go) based on the existence of appropriately-named methods on `Server` in this package. Run `make generate` in the root directory to regenerate it with any changes in available methods.
* The business logic for the various language features is divided between files in [`impl`](impl) (other than `server.go` and `server_gen.go`).

## Folder structure

* `framework/`: The LSP boilerplate and plumbing copied from `gopls`
    * `fakenet/`: implementation of the `net.Conn` interface over a reader and writer (used for stdio)
    * `jsonrpc2/`: minimal implementation of the JSON-RPC 2.0 spec
    * `lsp/`: LSP server boilerplate derived from the spec
    * `xcontext/`: extra functionality for the in-built `context` package
* `impl/`: business logic computing and providing Sysl language feature content to client
* `server/`: Connects the server in `impl` to a stream
