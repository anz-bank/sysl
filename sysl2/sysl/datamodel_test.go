package main

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/roothandler"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
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
	root := ""
	rootHandler, err := roothandler.NewRootHandler(root, "doesn't-exist.sysl",
		afero.NewOsFs(), logrus.StandardLogger())
	assert.NoError(t, err)
	_, err = parse.NewParser().Parse(rootHandler, syslutil.NewChrootFs(afero.NewOsFs(), rootHandler.Root()))
	require.Error(t, err)
}

func TestDoGenerateDataDiagrams(t *testing.T) {
	args := &dataArgs{
		modules: "data.sysl",
		output:  "%(epname).png",
		project: "Project",
	}
	argsData := []string{"sysl", "data", "-o", args.output, "-j", args.project, args.modules}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "datamodel")
}

func TestDoConstructDataDiagrams(t *testing.T) {
	args := &dataArgs{
		root:    "./tests/",
		modules: "/data.sysl",
		output:  "%(epname).png",
		project: "Project",
		title:   "empdata",
		expected: map[string]string{
			"Relational-Model.png": "tests/relational-model-golden.puml",
			"Object-Model.png":     "tests/object-model-golden.puml",
		},
	}
	result, err := DoConstructDataDiagramsWithParams(args.root, "", args.title, args.output, args.project,
		args.modules)
	assert.Nil(t, err, "Generating the data diagrams failed")
	comparePUML(t, args.expected, result)
}

func DoConstructDataDiagramsWithParams(
	rootModel, filter, title, output, project, modules string,
) (map[string]string, error) {
	classFormat := "%(classname)"
	cmdContextParamDatagen := &CmdContextParamDatagen{
		filter:      filter,
		title:       title,
		output:      output,
		project:     project,
		classFormat: classFormat,
	}

	logger, _ := test.NewNullLogger()
	rootHandler, err := roothandler.NewRootHandler(rootModel, modules, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	mod, _, err := LoadSyslModule(rootHandler, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateDataModels(cmdContextParamDatagen, mod, logger)
}
