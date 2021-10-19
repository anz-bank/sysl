package main

import (
	"context"
	"os"

	"github.com/anz-bank/sysl/pkg/lsp/framework/fakenet"
	"github.com/anz-bank/sysl/pkg/lsp/framework/jsonrpc2"
	"github.com/anz-bank/sysl/pkg/lsp/server"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"gopkg.in/alecthomas/kingpin.v2"
)

type lspCmd struct {
}

func (p *lspCmd) Name() string       { return "lsp" }
func (p *lspCmd) MaxSyslModule() int { return 0 }

func (p *lspCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Handle an LSP request as a language server")
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *lspCmd) Execute(args cmdutils.ExecuteArgs) error {
	ctx := context.Background()

	ss := server.NewStreamServer()
	stream := jsonrpc2.NewHeaderStream(fakenet.NewConn("stdio", os.Stdin, os.Stdout))
	return ss.ServeStream(ctx, stream)
}
