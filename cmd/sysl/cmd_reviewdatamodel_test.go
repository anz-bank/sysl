package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func TestDoGenerateDataDiagramsWithDataModelViewCmd(t *testing.T) {
	args := &dataArgs{
		modules: "datamodel/reviewdatamodelcmd.sysl",
		output:  "%(appname).png",
	}
	argsData := []string{"sysl", "reviewdata", "-o", args.output, args.modules}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "reviewdatamodel")
}

/*
func TestDoConstructDataDiagramsWithDataModelViewCmd(t *testing.T) {
	args := &dataArgs{
		root:    testDir,
		modules: "datamodel/reviewdatamodelcmd.sysl",
		output:  "%(appname).png",
		title:   "testdata",
		expected: map[string]string{
			"Test.png": filepath.Join(testDir, "datamodel/review-data-model-cmd.puml"),
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
	comparePUML(t, args.expected, result)
}
*/
