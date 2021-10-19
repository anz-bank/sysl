package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anz-bank/sysl/pkg/lsp/framework/fakenet"
	"github.com/anz-bank/sysl/pkg/lsp/framework/jsonrpc2"
	"github.com/anz-bank/sysl/pkg/lsp/server"
)

func main() {
	ctx := context.Background()

	ss := server.NewStreamServer()
	stream := jsonrpc2.NewHeaderStream(fakenet.NewConn("stdio", os.Stdin, os.Stdout))
	err := ss.ServeStream(ctx, stream)

	if err != nil {
		fmt.Println(err)
	}
}
