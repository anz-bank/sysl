package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestDoGenerateDataDiagramsWithDataModelViewCmd(t *testing.T) {
	args := &dataArgs{
		modules: "dataviewcmd.sysl",
		output:  "%(epname).png",
	}
	argsData := []string{"sysl", "dataview", "-o", args.output, args.modules}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "datamodelview")
}

func TestDoConstructDataDiagramsWithDataModelViewCmd(t *testing.T) {
	args := &dataArgs{
		root:    testDir,
		modules: "dataviewcmd.sysl",
		output:  "%(epname).png",
		title:   "testdata",
		expected: map[string]string{
			"Default.png": filepath.Join(testDir, "data-view-cmd.puml"),
		},
	}

	var result map[string]string
	logger, _ := test.NewNullLogger()
	mod, _, err := LoadSyslModule(args.root, args.modules, afero.NewOsFs(), logger)
	if err != nil {
		result = nil
	} else {
		result, err = GenerateDataModelsView(&CmdContextParamDatagen{
			title:  args.title,
			output: args.output,
		}, mod, logger)
	}

	assert.Nil(t, err, "Generating the data diagrams failed")
	fmt.Printf("%s", result["Default.png"])
	comparePUML(t, args.expected, result)
}
