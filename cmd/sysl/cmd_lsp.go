package main

import (
	"context"
	"os"

	"github.com/anz-bank/sysl/internal/jsonrpc2"
	"github.com/anz-bank/sysl/internal/lsp/cache"
	"github.com/anz-bank/sysl/internal/lsp/lsprpc"
	"github.com/anz-bank/sysl/pkg/cmdutils"

	"gopkg.in/alecthomas/kingpin.v2"
)

type lspCmd struct{}

func (p *lspCmd) Name() string       { return "lsp" }
func (p *lspCmd) MaxSyslModule() int { return 0 }

func (p *lspCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Start the sysl language server")
	return cmd
}

// Serve struct internal/lsp/serve.go
func (p *lspCmd) Execute(args cmdutils.ExecuteArgs) error {
	// start the lsp debug server here
	ctx := context.Background()
	/*
		di := debug.GetInstance(ctx)
		if di != nil {
			closeLog, err := di.SetLogFile(s.Logfile)
			if err != nil {
				return err
			}
			defer closeLog()
			di.ServerAddress = s.Address
			di.DebugAddress = s.Debug
			di.Serve(ctx)
			di.MonitorMemory(ctx)
		}
	*/

	var ss jsonrpc2.StreamServer = lsprpc.NewStreamServer(cache.New(ctx, nil), true)
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	return ss.ServeStream(ctx, stream)
}
