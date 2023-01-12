package impl

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	parser "github.com/anz-bank/sysl/pkg/grammar"
	"github.com/anz-bank/sysl/pkg/lsp/framework/jsonrpc2"

	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
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

func (s *Server) provideSyntaxDiagnostics(ctx context.Context, uri protocol.DocumentURI, version int32) error {
	fsPath := strings.TrimPrefix(string(uri), "file://")
	errListener := NewSyslErrorListener(ctx, s.client)

	//spanUri := uri.SpanURI()
	//if !spanUri.IsFile() {
	//	return nil
	//}

	// TODO: in the future support didChange (syntax checking)
	//	text, err := s.changedText(ctx, uri, params.ContentChanges)
	//	if err != nil {
	//		return err
	//	}

	fh, err := os.Open(fsPath)
	if err != nil {
		return errors.Wrapf(err, "file not found (%v)", jsonrpc2.ErrInternal)
	}
	defer fh.Close()

	content, err := io.ReadAll(fh)
	if err != nil {
		return errors.Wrapf(err, "file not found (%v)", jsonrpc2.ErrInternal)
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
		return errors.Wrapf(err, "PublishDiagnostics failed")
	}

	return nil
}
