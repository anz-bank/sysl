package main

import (
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
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
		root:    testDir,
		modules: "data.sysl",
		output:  "%(epname).png",
		project: "Project",
		title:   "empdata",
		expected: map[string]string{
			"Relational-Model.png": filepath.Join(testDir, "relational-model-golden.puml"),
			"Object-Model.png":     filepath.Join(testDir, "object-model-golden.puml"),
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
	mod, _, err := LoadSyslModule(rootModel, modules, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateDataModels(cmdContextParamDatagen, mod, logger)
}
