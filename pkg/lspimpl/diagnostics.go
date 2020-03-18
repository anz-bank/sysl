// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lspimpl

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/anz-bank/sysl/internal/jsonrpc2"
	parser "github.com/anz-bank/sysl/pkg/grammar"

	"github.com/anz-bank/sysl/internal/lsp/protocol"
)

func (s *Server) diagnoseRaw(input string, listener *SyslErrorListener) parser.ISysl_fileContext {
	var chars = antlr.NewInputStream(input)
	var lexer = parser.NewSyslLexer(chars)
	var tokens = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	var parser = parser.NewSyslParser(tokens)

	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	parser.RemoveErrorListeners()
	parser.BuildParseTrees = true
	parser.AddErrorListener(listener)

	// parse file here
	return parser.Sysl_file()
}

func (s *Server) provideSyntaxDiagnostics(ctx context.Context, uri protocol.DocumentURI, version float64) error {
	errListener := NewSyslErrorListener(ctx, s.client)

	spanUri := uri.SpanURI()
	if !spanUri.IsFile() {
		return nil
	}

	// TODO: in the future support didChange (syntax checking)
	//	text, err := s.changedText(ctx, uri, params.ContentChanges)
	//	if err != nil {
	//		return err
	//	}

	fh, err := os.Open(spanUri.Filename())
	if err != nil {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "file not found (%v)", err)
	}
	defer fh.Close()

	content, err := ioutil.ReadAll(fh)
	if err != nil {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "file not found (%v)", err)
	}

	// TODO: figure out a way to run this in a goroutine (not sure why it doesn't work)
	s.diagnoseRaw(string(content), errListener)
	// publish the errors
	err = s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		Diagnostics: errListener.Errors,
		URI:         uri,
		Version:     version,
	})
	if err != nil {
		//fmt.Println("publishReports failed")
		return err
	}

	return nil
}
