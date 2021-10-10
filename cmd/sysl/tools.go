// +build tools

package main

import (
	// Prevent 'go mod tidy' from removing below packages, otherwise 'make test' fails.
	_ "github.com/arr-ai/proto"
	_ "github.com/chzyer/readline"
	_ "github.com/gorilla/websocket"
	_ "github.com/rjeczalik/notify"
)
