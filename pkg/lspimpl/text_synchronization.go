// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lspimpl

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/sysl/pkg/lspimpl/lspframework/lsp/protocol"
)

type SyslErrorListener struct {
	antlr.ErrorListener

	Errors []protocol.Diagnostic
	ctx    context.Context
	client protocol.Client
}

func NewSyslErrorListener(ctx context.Context, client protocol.Client) *SyslErrorListener {
	return &SyslErrorListener{Errors: []protocol.Diagnostic{}, ctx: ctx, client: client}
}

func (l *SyslErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.client.LogMessage(l.ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "SyntaxError"})

	// assume an offendingSymbol is an antlr.Token type
	token := offendingSymbol.(antlr.Token)

	// continue
	l.Errors = append(l.Errors, protocol.Diagnostic{
		Message: msg,
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      float64(line - 1),
				Character: float64(column),
			},
			End: protocol.Position{
				Line:      float64(line - 1),
				Character: float64(column + len(token.GetText())),
			},
		},
		Severity: protocol.SeverityError,
		Source:   "sysl-lsp",
		Tags:     []protocol.DiagnosticTag{},
	})
}

func (l *SyslErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	l.client.LogMessage(l.ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "ReportAmbiguity"})
}

func (l *SyslErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	l.client.LogMessage(l.ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "ReportAttemptingFullContext"})
}

func (l *SyslErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	l.client.LogMessage(l.ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "ReportContextSensitivity"})
}

func (s *Server) didOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.client.LogMessage(ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "didOpen"})
	return s.didModifyFiles(ctx, params.TextDocument.URI, params.TextDocument.Version)
}

func (s *Server) didChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	s.client.LogMessage(ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "didChange"})
	return nil
}

// only show syntax diagnostics on save
func (s *Server) didSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	s.client.LogMessage(ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "didSave"})
	return s.didModifyFiles(ctx, params.TextDocument.URI, params.TextDocument.Version)
}

func (s *Server) didClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.client.LogMessage(ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "didClose"})

	// If a file has been closed and is not on disk, clear its diagnostics.
	uri := params.TextDocument.URI.SpanURI()
	fh, err := os.Open(uri.Filename())
	if err != nil {
		return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
			URI:         protocol.URIFromSpanURI(uri),
			Diagnostics: []protocol.Diagnostic{},
		})
	}
	defer fh.Close()

	if _, err := ioutil.ReadAll(fh); err != nil {
		return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
			URI:         protocol.URIFromSpanURI(uri),
			Diagnostics: []protocol.Diagnostic{},
		})
	}
	return nil
}

func (s *Server) didModifyFiles(ctx context.Context, uri protocol.DocumentURI, version float64) error {
	s.client.LogMessage(ctx, &protocol.LogMessageParams{Type: protocol.Log, Message: "didModifyFiles"})

	return s.provideSyntaxDiagnostics(ctx, uri, version)
}
