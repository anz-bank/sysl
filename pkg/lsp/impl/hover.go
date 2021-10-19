package impl

import (
	"context"

	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
)

func (s *Server) hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	// TODO: reuse this to handle sysl code
	/*
		_, _, ok, err := s.beginFileRequest(params.TextDocument.URI, source.UnknownKind)
		if !ok {
			return nil, err
		}
	*/
	return nil, nil
}
