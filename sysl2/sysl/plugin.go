package main

import (
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/plugins"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

func runPlugin(
	generate string,
	openOutfile func() (io.Writer, error),
	module *sysl.Module,
	debug bool,
) error {
	parts := strings.SplitN(generate, ":", 2)
	generator := parts[0]
	parameter := ""
	if len(parts) > 1 {
		parameter = parts[1]
	}

	var logger hclog.Logger
	if !debug {
		logger = hclog.NewNullLogger()
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  plugins.Handshake,
		Plugins:          plugins.PluginMap,
		Cmd:              exec.Command("sh", "-c", "sysl_"+generator),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	codegenRaw, err := rpcClient.Dispense(plugins.CodeGeneratorName)
	if err != nil {
		return err
	}

	codegen := codegenRaw.(plugins.CodeGeneratorClient)
	response, err := codegen.GenerateCode(
		context.Background(),
		&plugins.GenerateCodeRequest{
			GeneratorName: generator,
			Parameter:     parameter,
			Module:        module,
		},
	)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return errors.Errorf("Plugin error: %s", response.Error)
	}

	for _, file := range response.OutputFile {
		outfile, err := openOutfile()
		if err != nil {
			return err
		}
		if _, err := outfile.Write([]byte(file.Contents)); err != nil {
			return errors.Wrapf(err, "Writing generated file")
		}
	}

	return nil
}
