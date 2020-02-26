package syslutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCMDFlags(t *testing.T) {
	t.Parallel()
	flags, err := ReadCMDFlags("tests/config.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(flags))
	assert.Equal(t, "--grammar=go.gen.g", flags[0])
	assert.Equal(t, "--transform=go.gen.sysl", flags[1])
	assert.Equal(t, "model.sysl", flags[2])
}

func TestPopulateCMDFlagsFromFile(t *testing.T) {
	t.Parallel()
	cmdArgs := []string{"sysl", "codegen", "@tests/config.txt"}
	flags, err := PopulateCMDFlagsFromFile(cmdArgs)
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(flags))
	assert.Equal(t, "sysl", flags[0])
	assert.Equal(t, "codegen", flags[1])
	assert.Equal(t, "--grammar=go.gen.g", flags[2])
	assert.Equal(t, "--transform=go.gen.sysl", flags[3])
	assert.Equal(t, "model.sysl", flags[4])
}
