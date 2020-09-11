package datamodeldiagram

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDataDiagFail(t *testing.T) {
	t.Parallel()
	_, err := parse.NewParser().Parse("doesn't-exist.sysl", syslutil.NewChrootFs(afero.NewOsFs(), ""))
	require.Error(t, err)
}

func TestDoConstructDataDiagramsWithProjectMannerModule(t *testing.T) {
	args := &DataArgs{
		Root:    testDir,
		Modules: "data.sysl",
		Output:  "%(epname).png",
		Project: "Project",
		Title:   "empdata",
		Expected: map[string]string{
			"Relational-Model.png":      filepath.Join(testDir, "relational-model-golden.puml"),
			"Object-Model.png":          filepath.Join(testDir, "object-model-golden.puml"),
			"Primitive-Alias-Model.png": filepath.Join(testDir, "primitive-alias-model-golden.puml"),
		},
	}
	result, err := DoConstructDataDiagramsWithParams(args.Root, "", args.Title, args.Output, args.Project,
		args.Modules)
	assert.Nil(t, err, "Generating the data diagrams failed")
	diagrams.ComparePUML(t, args.Expected, result)
}

func DoConstructDataDiagramsWithParams(
	rootModel, filter, title, output, project, modules string,
) (map[string]string, error) {
	classFormat := "%(classname)"
	cmdContextParamDatagen := &cmdutils.CmdContextParamDatagen{
		Filter:      filter,
		Title:       title,
		Output:      output,
		Project:     project,
		ClassFormat: classFormat,
	}

	logger, _ := test.NewNullLogger()
	mod, _, err := loader.LoadSyslModule(rootModel, modules, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateDataModels(cmdContextParamDatagen, mod, logger)
}

func TestDoConstructDataDiagramsWithPureModule(t *testing.T) {
	args := &DataArgs{
		Root:    testDir,
		Modules: "reviewdatamodelcmd.sysl",
		Output:  "%(epname).png",
		Title:   "testdata",
		Expected: map[string]string{
			"Test.png": filepath.Join(testDir, "review-data-model-cmd.puml"),
		},
	}

	var result map[string]string
	logger, _ := test.NewNullLogger()
	mod, _, err := loader.LoadSyslModule(args.Root, args.Modules, afero.NewOsFs(), logger)
	if err != nil {
		result = nil
	} else {
		result, err = GenerateDataModels(&cmdutils.CmdContextParamDatagen{
			Title:  args.Title,
			Output: args.Output,
			Direct: true,
		}, mod, logger)
	}
	assert.Nil(t, err, "Generating the data diagrams failed")
	fmt.Printf("\n\n\n%s", result)

	diagrams.ComparePUML(t, args.Expected, result)
}
