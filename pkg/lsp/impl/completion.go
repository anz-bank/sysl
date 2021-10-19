package impl

import (
	"context"

	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
)

func (s *Server) completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return nil, nil
}
