// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"context"

	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
)

func (s *Server) documentLink(ctx context.Context, params *protocol.DocumentLinkParams) (links []protocol.DocumentLink, err error) {
	return nil, nil
}
