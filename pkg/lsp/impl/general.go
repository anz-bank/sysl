package impl

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/anz-bank/sysl/pkg/lsp/framework/jsonrpc2"
	"github.com/anz-bank/sysl/pkg/lsp/framework/lsp/protocol"
)

func (s *Server) initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	s.stateMu.Lock()
	state := s.state
	s.stateMu.Unlock()
	if state >= serverInitializing {
		return nil, errors.Wrapf(jsonrpc2.ErrInvalidRequest, "initialize called while server in %v state", s.state)
	}
	s.stateMu.Lock()
	s.state = serverInitializing
	s.stateMu.Unlock()

	//options := s.session.Options()
	//defer func() { s.session.SetOptions(options) }()

	// TODO
	// client initialization options
	//fmt.Println(params.InitializationOptions)

	// client capabilities
	//fmt.Println(params.Capabilities)

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			DefinitionProvider: true,
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				Change:    protocol.Full,
				OpenClose: true,
				Save: protocol.SaveOptions{
					IncludeText: false,
				},
			},
		},
	}, nil
}

func (s *Server) initialized(ctx context.Context, params *protocol.InitializedParams) error {
	s.stateMu.Lock()
	s.state = serverInitialized
	s.stateMu.Unlock()

	//options := s.session.Options()
	//defer func() { s.session.SetOptions(options) }()

	/*
		var registrations []protocol.Registration
		if options.ConfigurationSupported && options.DynamicConfigurationSupported {
			registrations = append(registrations,
				protocol.Registration{
					ID:     "workspace/didChangeConfiguration",
					Method: "workspace/didChangeConfiguration",
				},
				protocol.Registration{
					ID:     "workspace/didChangeWorkspaceFolders",
					Method: "workspace/didChangeWorkspaceFolders",
				},
			)
		}

		if options.DynamicWatchedFilesSupported {
			registrations = append(registrations, protocol.Registration{
				ID:     "workspace/didChangeWatchedFiles",
				Method: "workspace/didChangeWatchedFiles",
				RegisterOptions: protocol.DidChangeWatchedFilesRegistrationOptions{
					Watchers: []protocol.FileSystemWatcher{{
	*/
	//					GlobPattern: "**/*.sysl",
	/*
						Kind:        float64(protocol.WatchChange + protocol.WatchDelete + protocol.WatchCreate),
					}},
				},
			})
		}

		if len(registrations) > 0 {
			s.client.RegisterCapability(ctx, &protocol.RegistrationParams{
				Registrations: registrations,
			})
		}

		buf := &bytes.Buffer{}
		debug.PrintVersionInfo(buf, true, debug.PlainText)
		log.Print(ctx, buf.String())

		s.addFolders(ctx, s.pendingFolders)
		s.pendingFolders = nil

	*/
	return nil
}

/*
func (s *Server) addFolders(ctx context.Context, folders []protocol.WorkspaceFolder) {
	originalViews := len(s.session.Views())
	viewErrors := make(map[span.URI]error)

	for _, folder := range folders {
		uri := span.URIFromURI(folder.URI)
		_, snapshot, err := s.addView(ctx, folder.Name, uri)
		if err != nil {
			viewErrors[uri] = err
			continue
		}
		go s.diagnoseDetached(snapshot)
	}
	if len(viewErrors) > 0 {
		errMsg := fmt.Sprintf("Error loading workspace folders (expected %v, got %v)\n", len(folders), len(s.session.Views())-originalViews)
		for uri, err := range viewErrors {
			errMsg += fmt.Sprintf("failed to load view for %s: %v\n", uri, err)
		}
		s.client.ShowMessage(ctx, &protocol.ShowMessageParams{
			Type:    protocol.Error,
			Message: errMsg,
		})
	}
}
*/

/*
func (s *Server) fetchConfig(ctx context.Context, name string, folder span.URI, o *source.Options) error {
	if !s.session.Options().ConfigurationSupported {
		return nil
	}
	v := protocol.ParamConfiguration{
		ConfigurationParams: protocol.ConfigurationParams{
			Items: []protocol.ConfigurationItem{{
				ScopeURI: string(folder),
				Section:  "gopls",
			}, {
				ScopeURI: string(folder),
				Section:  fmt.Sprintf("gopls-%s", name),
			}},
		},
	}
	configs, err := s.client.Configuration(ctx, &v)
	if err != nil {
		return err
	}
	for _, config := range configs {
		results := source.SetOptions(o, config)
		for _, result := range results {
			if result.Error != nil {
				s.client.ShowMessage(ctx, &protocol.ShowMessageParams{
					Type:    protocol.Error,
					Message: result.Error.Error(),
				})
			}
			switch result.State {
			case source.OptionUnexpected:
				s.client.ShowMessage(ctx, &protocol.ShowMessageParams{
					Type:    protocol.Error,
					Message: fmt.Sprintf("unexpected config %s", result.Name),
				})
			case source.OptionDeprecated:
				msg := fmt.Sprintf("config %s is deprecated", result.Name)
				if result.Replacement != "" {
					msg = fmt.Sprintf("%s, use %s instead", msg, result.Replacement)
				}
				s.client.ShowMessage(ctx, &protocol.ShowMessageParams{
					Type:    protocol.Warning,
					Message: msg,
				})
			}
		}
	}
	return nil
}
*/

// beginFileRequest checks preconditions for a file-oriented request and routes
// it to a snapshot.
// We don't want to return errors for benign conditions like wrong file type,
// so callers should do if !ok { return err } rather than if err != nil.
/*
func (s *Server) beginFileRequest(pURI protocol.DocumentURI, expectKind source.FileKind) (source.Snapshot, source.FileHandle, bool, error) {
	uri := pURI.SpanURI()
	if !uri.IsFile() {
		// Not a file URI. Stop processing the request, but don't return an error.
		return nil, nil, false, nil
	}
	view, err := s.session.ViewOf(uri)
	if err != nil {
		return nil, nil, false, err
	}
	snapshot := view.Snapshot()
	fh, err := snapshot.GetFile(uri)
	if err != nil {
		return nil, nil, false, err
	}
	if expectKind != source.UnknownKind && fh.Identity().Kind != expectKind {
		// Wrong kind of file. Nothing to do.
		return nil, nil, false, nil
	}
	return snapshot, fh, true, nil
}
*/

func (s *Server) shutdown(ctx context.Context) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state < serverInitialized {
		logrus.WithContext(ctx).Error("server shutdown without initialization")
	}
	if s.state != serverShutDown {
		// drop all the active views
		// s.session.Shutdown(ctx)
		s.state = serverShutDown
	}
	return nil
}

// ServerExitFunc is used to exit when requested by the client. It is mutable
// for testing purposes.
var ServerExitFunc = os.Exit

func (s *Server) exit(ctx context.Context) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state != serverShutDown {
		ServerExitFunc(1)
	}
	ServerExitFunc(0)
	return nil
}

func setBool(b *bool, m map[string]interface{}, name string) {
	if v, ok := m[name].(bool); ok {
		*b = v
	}
}

func setNotBool(b *bool, m map[string]interface{}, name string) {
	if v, ok := m[name].(bool); ok {
		*b = !v
	}
}
