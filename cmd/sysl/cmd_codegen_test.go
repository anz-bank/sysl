package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestLoadFlags(t *testing.T) {
	var cmd codegenCmd
	var err error
	cmd = codegenCmd{config: "../../pkg/config/tests/config.yml"}
	err = cmd.loadFlags()
	assert.Equal(t, nil, err)
	assert.Equal(t, "go.gen.g", cmd.grammar)
	assert.Equal(t, "go.gen.sysl", cmd.transform)
	assert.Equal(t, "depPath", cmd.depPath)
	assert.Equal(t, "basePath", cmd.basePath)
	assert.Equal(t, "appName", cmd.appName)

	cmd = codegenCmd{}
	cmd.grammar = "grammar"
	cmd.transform = "transform"
	err = cmd.loadFlags()
	assert.Equal(t, nil, err)
	assert.Equal(t, "grammar", cmd.grammar)
	assert.Equal(t, "transform", cmd.transform)
	assert.Equal(t, "", cmd.depPath)
	assert.Equal(t, "", cmd.basePath)
	assert.Equal(t, "", cmd.appName)

	cmd = codegenCmd{}
	assert.Error(t, cmd.loadFlags())

	cmd = codegenCmd{}
	cmd.grammar = "grammar"
	assert.Error(t, cmd.loadFlags())
}

func TestCodegenIndividualFlagsCMD(t *testing.T) {
	argsData := []string{"sysl", "codegen", "--grammar=grammar", "--transform=transform", "hello"}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "codegen")
}

func TestCodegenConfigCMD(t *testing.T) {
	argsData := []string{"sysl", "codegen", "--config=config.yml", "hello"}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "codegen")
}
