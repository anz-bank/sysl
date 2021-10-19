package impl

import (
	"context"

	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
)

func (s *Server) codeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, nil
}
