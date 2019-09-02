package main

import (
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

type dataArgs struct {
	root     string
	title    string
	output   string
	project  string
	modules  string
	expected map[string]string
}

func TestGenerateDataDiagFail(t *testing.T) {
	t.Parallel()
	_, err := parse.NewParser().Parse("doesn't-exist.sysl", syslutil.NewChrootFs(afero.NewOsFs(), ""))
	require.Error(t, err)
}

func TestDoGenerateDataDiagrams(t *testing.T) {
	args := &dataArgs{
		root:    "./tests/",
		modules: "data.sysl",
		output:  "%(epname).png",
		project: "Project",
	}
	argsData := []string{"sysl", "data", "--root", args.root, "-o", args.output, "-j", args.project, args.modules}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")
	configureCmdlineForDatagen(syslCmd)
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "data")
}

func TestDoConstructDataDiagrams(t *testing.T) {
	args := &dataArgs{
		root:    "./tests/",
		modules: "data.sysl",
		output:  "%(epname).png",
		project: "Project",
		title:   "empdata",
		expected: map[string]string{
			"Relational-Model.png": "tests/relational-model-golden.puml",
			"Object-Model.png":     "tests/object-model-golden.puml",
		},
	}
	result, err := DoConstructDataDiagramsWithParams(args.root, "", args.title, args.output, "warn", args.project,
		args.modules, false)
	assert.Nil(t, err, "Generating the data diagrams failed")
	comparePUML(t, args.expected, result)
}

func DoConstructDataDiagramsWithParams(
	rootModel, filter, title, output, loglevel, project, modules string,
	isVerbose bool,
) (map[string]string, error) {
	plantuml := ""
	classFormat := "%(classname)"
	cmdContextParamDatagen := &CmdContextParamDatagen{
		root:        &rootModel,
		filter:      &filter,
		title:       &title,
		output:      &output,
		loglevel:    &loglevel,
		project:     &project,
		isVerbose:   &isVerbose,
		plantuml:    &plantuml,
		modules:     &modules,
		classFormat: &classFormat,
	}
	return GenerateDataModels(cmdContextParamDatagen)
}
