package impl

import (
	"context"
	"errors"

	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
)

func (s *Server) executeCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	switch params.Command {
	case "tidy":
		return nil, errors.New("tidy error")
	case "upgrade.dependency":
		return nil, errors.New("upgrade dependency error")
	}
	return nil, nil
}
