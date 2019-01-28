package plugins

import (
	"context"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Plugin name constant
const (
	CodeGeneratorName = "CodeGenerator"
)

// SingleFileGenerateCodeResponse returns a GenerateCodeResponse containing a
// single file.
func SingleFileGenerateCodeResponse(filename, contents string) *GenerateCodeResponse {
	return &GenerateCodeResponse{
		OutputFile: []*GenerateCodeResponse_OutputFile{
			{
				Filename: filename,
				Contents: contents,
			},
		},
	}
}

// Handshake is shared between sysl and its plugins.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "SYSL_PLUGIN",
	MagicCookieValue: "2d76d777-4fde-4659-b7f9-211bd4e2cfe2",
}

// PluginMap is used by the host to define the plugins it supports.
var PluginMap = map[string]plugin.Plugin{
	CodeGeneratorName: &CodeGeneratorPlugin{},
}

// ServeCodeGenerator serves a code generator gRPC service as a plugin.
func ServeCodeGenerator(server CodeGeneratorServer) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		GRPCServer:      plugin.DefaultGRPCServer,
		Plugins: map[string]plugin.Plugin{
			CodeGeneratorName: &CodeGeneratorPlugin{
				Impl: server,
			},
		},
	})
}

// CodeGeneratorPlugin provides the CodeGenerator proto service as a plugin.
type CodeGeneratorPlugin struct {
	plugin.Plugin
	Impl CodeGeneratorServer
}

// GRPCServer creates a gRPC server.
func (p *CodeGeneratorPlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	RegisterCodeGeneratorServer(s, p.Impl)
	return nil
}

// GRPCClient creates a gRPC client.
func (p *CodeGeneratorPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return NewCodeGeneratorClient(c), nil
}
