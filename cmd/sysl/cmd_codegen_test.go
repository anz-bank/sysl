package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestLoadFlags(t *testing.T) {
	var cmd codegenCmd
	var err error
	cmd = codegenCmd{}
	cmd.Config = "../../pkg/config/tests/config.yml"
	err = cmd.loadFlags()
	assert.Equal(t, nil, err)
	assert.Equal(t, "go.gen.g", cmd.Grammar)
	assert.Equal(t, "go.gen.sysl", cmd.Transform)
	assert.Equal(t, "depPath", cmd.DepPath)
	assert.Equal(t, "basePath", cmd.BasePath)
	assert.Equal(t, "appName", cmd.AppName)

	cmd = codegenCmd{}
	cmd.Grammar = "grammar"
	cmd.Transform = "transform"
	err = cmd.loadFlags()
	assert.Equal(t, nil, err)
	assert.Equal(t, "grammar", cmd.Grammar)
	assert.Equal(t, "transform", cmd.Transform)
	assert.Equal(t, "", cmd.DepPath)
	assert.Equal(t, "", cmd.BasePath)
	assert.Equal(t, "", cmd.AppName)

	cmd = codegenCmd{}
	assert.Error(t, cmd.loadFlags())

	cmd = codegenCmd{}
	cmd.Grammar = "grammar"
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
