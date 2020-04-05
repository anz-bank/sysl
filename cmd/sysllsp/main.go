package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anz-bank/sysl/pkg/lspimpl/lspframework/jsonrpc2"
	"github.com/anz-bank/sysl/pkg/lspimpl/lsprpc"
)

func main() {
	// start the lsp debug server here
	ctx := context.Background()

	// TODO: reuse this to handle sysl code
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

	ss := lsprpc.NewStreamServer(true)
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	err := ss.ServeStream(ctx, stream)

	if err != nil {
		fmt.Println(err)
	}
}
