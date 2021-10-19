package impl

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
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
	// assume an offendingSymbol is an antlr.Token type
	token := offendingSymbol.(antlr.Token)

	// continue
	l.Errors = append(l.Errors, protocol.Diagnostic{
		Message: msg,
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(line - 1),
				Character: uint32(column),
			},
			End: protocol.Position{
				Line:      uint32(line - 1),
				Character: uint32(column + len(token.GetText())),
			},
		},
		Severity: protocol.SeverityError,
		Source:   "sysl-lsp",
		Tags:     []protocol.DiagnosticTag{},
	})
}

func (l *SyslErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

func (l *SyslErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

func (l *SyslErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
}

func (s *Server) didOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	return s.didModifyFiles(ctx, params.TextDocument.URI, params.TextDocument.Version)
}

func (s *Server) didChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	return nil
}

// only show syntax diagnostics on save
func (s *Server) didSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	return s.didModifyFiles(ctx, params.TextDocument.URI, 0)
}

func (s *Server) didClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	// If a file has been closed and is not on disk, clear its diagnostics.
	uri := params.TextDocument.URI
	url, err := url.Parse(string(params.TextDocument.URI))
	if err != nil {
		return err
	}
	fh, err := os.Open(url.Path)
	if err != nil {
		return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: []protocol.Diagnostic{},
		})
	}
	defer fh.Close()

	if _, err := ioutil.ReadAll(fh); err != nil {
		return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: []protocol.Diagnostic{},
		})
	}
	return nil
}

func (s *Server) didModifyFiles(ctx context.Context, uri protocol.DocumentURI, version int32) error {
	return s.provideSyntaxDiagnostics(ctx, uri, version)
}
