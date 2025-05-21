package main

import (
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/stretchr/testify/assert"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"
)

func TestDoGenerateDataDiagramsWithProjectMannerModuleCMD(t *testing.T) {
	t.Parallel()

	args := &datamodeldiagram.DataArgs{
		Modules: "data.sysl",
		Output:  "%(epname).png",
		Project: "Project",
	}
	argsData := []string{"sysl", "data", "-o", args.Output, "-j", args.Project, args.Modules}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "datamodel")
}

func TestDoGenerateDataDiagramsWithPureModuleCMD(t *testing.T) {
	t.Parallel()

	args := &datamodeldiagram.DataArgs{
		Modules: "reviewdatamodelcmd.sysl",
		Output:  "%(epname).png",
	}
	argsData := []string{"sysl", "data", "-d", "-o", args.Output, args.Modules}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "datamodel")
}
