package main

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/alecthomas/kingpin/v2"
	"github.com/stretchr/testify/assert"
)

func TestDoGenerateIntegrations(t *testing.T) {
	t.Parallel()

	args := &integrationdiagram.IntsArg{
		Modules: "indirect_1.sysl",
		Output:  "%(epname).png",
		Project: "Project",
	}
	argsData := []string{"sysl", "ints", "-o", args.Output, "-j", args.Project, args.Modules}
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(sysl))
	selectedCommand, err := sysl.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl ints")
	assert.Equal(t, selectedCommand, "integrations")
}
