package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCodeGenFlags(t *testing.T) {
	t.Parallel()
	config, err := ReadCodeGenFlags("tests/config.yml")
	assert.Equal(t, nil, err)
	assert.Equal(t, "go.gen.g", config.Grammar)
	assert.Equal(t, "go.gen.sysl", config.Transform)
	assert.Equal(t, "depPath", config.DepPath)
	assert.Equal(t, "basePath", config.BasePath)
	assert.Equal(t, "appName", config.AppName)
}
