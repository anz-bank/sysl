# Forked code from golang.org/x/tools/internal/

## Original gopls intentions

From what I can see, the gopls implementation was originally written to support multiple workspaces connecting to a remote gopls instance as well as connecting locally. This is probably so multiple editors with a LSP client can connect to a single process. (NOT EXACTLY SURE IF THIS IS TRUE)

The sysl lsp server will not contain this feature and each client will create a new instance of the sysl lsp server (at least for now)

## Forked on 2020/03/12

* This folder mostly contains forked lsp code from the gopls implementation.
* Some code is modified. Code which is modified is prepended with a comment `// MODIFIED: SYSL_LSP`
* A general overview of how this lsp extension works:
  * The LSP interface and data structures are transpiled from TypeScript to Golang using the typescript repository: https://github.com/microsoft/vscode-languageserver-node.
    * This is located in the internal/protocol/ folder and more information is contained there.
  * The rest of the code is wrapped/frameworking code to originally provide gopls with debugging/golang diagnostics
    * The only partly stripped/modified code is the debugging server for lsprpc (which currently is not working, but may be used for future debugging)
  * There is a helper/helper.go binary for server.go generation --> more information located in the internal/lsp/helper/ folder
  * Use the helper binary to generate the concrete server implementation (maps the function to the lsp protocol call)
    * e.g. "textDocument/didOpen" --> mapped to `func (s *Server) didOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error` in internal/lspimpl/text_synchronization.go
    * Output file resulted in internal/lspimpl/server_gen.go

## Folder structure

* internal/jsonrpc2: contains a minimal implementation of the jsonrpc2 spec
* internal/span: contains some code for OS independent file names and positions and ranges in files
* internal/lsp: lsp frameworking code
* internal/telemetry: seems to be telemetry code for prometheous or ocagent (used by the autogenerator)
* internal/xcontext: extra functionality for the inbuilt context package (used by the autogenerator)

* pkg/lspimpl (non gopls code): this package is created to separate lsp framework code from the sysl diagnostic code

## Updating instructions

1. Update the ts protocol files in go. Instructions in internal/lsp/protocol/typescript/README.md
2. Compile helper binary in internal/lsp/helper/ [cd internal/lsp/helper && go build .]
3. Generate server_gen.go in pkg/lspimpl/ [cd ../pkg/lspimpl && go generate server.go]
4. Compile sysllsp [make buildlsp] in project root
