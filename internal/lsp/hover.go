// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"github.com/anz-bank/sysl/internal/lsp/protocol"
	"github.com/anz-bank/sysl/internal/lsp/source"
)

func (s *Server) hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	s.client.LogMessage(ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "hover"})

	_, _, ok, err := s.beginFileRequest(params.TextDocument.URI, source.UnknownKind)
	if !ok {
		return nil, err
	}
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.PlainText,
			Value: "lmao",
		},
		Range: protocol.Range{
			Start: params.Position,
			End:   params.Position,
		},
	}, nil
}
