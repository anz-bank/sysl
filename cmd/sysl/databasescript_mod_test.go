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

type scriptModArgs struct {
	root      string
	title     string
	output    string
	project   string
	sourceOld string
	sourceNew string
	expected  map[string]string
}

func TestGenerateDataScriptModFail(t *testing.T) {
	t.Parallel()
	_, err := parse.NewParser().Parse("doesn't-exist.sysl", syslutil.NewChrootFs(afero.NewOsFs(), ""))
	require.Error(t, err)
}

func TestDoGenerateDataScriptMod(t *testing.T) {
	args := &scriptModArgs{
		root:      testDir + "db_scripts/",
		sourceOld: "dataForSqlScriptOrg.sysl",
		sourceNew: "dataForSqlScriptModified.sysl",
		output:    "%(epname).sql",
		project:   "Project",
		title:     "Petstore Schema",
	}
	argsData := []string{"sysl", "generatescriptdelta", "-t", args.title, "-o", args.output,
		"-r", args.root, "-s", args.sourceOld, "-n", args.sourceNew, "-j", args.project}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "generate-script-delta")
}
func TestDoConstructDatabaseScriptMod(t *testing.T) {
	args := &scriptModArgs{
		root:      testDir + "db_scripts/",
		sourceOld: "dataForSqlScriptOrg.sysl",
		sourceNew: "dataForSqlScriptModified.sysl",
		output:    "%(epname).sql",
		project:   "Project",
		title:     "Petstore Schema",
		expected: map[string]string{
			"Relational-Model.sql": filepath.Join(testDir, "db_scripts/postgres-modify-script-golden.sql"),
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.root, "", args.title, args.output, args.project,
		args.sourceOld, args.sourceNew)
	assert.Nil(t, err, "Generating the sql script failed")
	compareSQL(t, args.expected, result)
}

func TestDoConstructDatabaseScriptModTwoApps(t *testing.T) {
	args := &scriptModArgs{
		root:      testDir + "db_scripts/",
		sourceOld: "dataForSqlScriptOrg.sysl",
		sourceNew: "dataForSqlScriptModifiedTwoApps.sysl",
		output:    "%(epname).sql",
		project:   "Project",
		title:     "Petstore Schema",
		expected: map[string]string{
			"Relational-Model.sql": filepath.Join(testDir, "/db_scripts/postgres-modify-script-golden_two_apps.sql"),
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.root, "", args.title, args.output, args.project,
		args.sourceOld, args.sourceNew)
	assert.Nil(t, err, "Generating the sql script failed")
	compareSQL(t, args.expected, result)
}

func DoConstructModDatabaseScriptWithParams(
	rootModel, filter, title, output, project, modulesOld, modulesNew string,
) (map[string]string, error) {
	cmdDatabaseScriptMod := &CmdDatabaseScriptMod{
		title:     title,
		root:      rootModel,
		orgSource: modulesOld,
		newSource: modulesNew,
		output:    output,
		project:   project,
	}

	logger, _ := test.NewNullLogger()
	return GenerateModDatabaseScripts(cmdDatabaseScriptMod, logger)
}
