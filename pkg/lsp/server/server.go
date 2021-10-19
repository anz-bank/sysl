// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lsp/server implements a jsonrpc2.StreamServer that may be used to
// serve the LSP on a jsonrpc2 channel.
package server

import (
	"context"

	"github.com/anz-bank/sysl/pkg/lsp/framework/jsonrpc2"
	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
	"github.com/anz-bank/sysl/pkg/lsp/impl"
	"github.com/pkg/errors"
)

// The StreamServer type is a jsonrpc2.StreamServer that handles incoming
// streams as a new LSP session.
type StreamServer struct {
	// serverForTest may be set to a test fake for testing.
	serverForTest protocol.Server
}

// NewStreamServer creates a StreamServer
func NewStreamServer() *StreamServer {
	return &StreamServer{}
}

// ServeStream implements the jsonrpc2.StreamServer interface, by handling
// incoming streams using a new lsp server.
func (s *StreamServer) ServeStream(ctx context.Context, stream jsonrpc2.Stream) error {
	conn := jsonrpc2.NewConn(stream)
	client := protocol.ClientDispatcher(conn)
	server := s.serverForTest
	if server == nil {
		server = impl.NewServer(client)
	}
	// Clients may or may not send a shutdown message. Make sure the server is
	// shut down.
	// TODO(rFindley): this shutdown should perhaps be on a disconnected context.
	defer server.Shutdown(ctx)
	conn.Go(ctx, protocol.Handlers(protocol.ServerHandler(server, jsonrpc2.MethodNotFound)))

	select {
	case <-conn.Done():
		client.Close()
	}

	var err error
	if conn.Err() != nil {
		err = errors.Errorf("remote disconnected: %v", conn.Err())
	}
	return err
}
