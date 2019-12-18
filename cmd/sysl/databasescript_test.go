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

type scriptArgs struct {
	root     string
	title    string
	output   string
	project  string
	source   string
	expected map[string]string
}

func TestGenerateDataScriptFail(t *testing.T) {
	t.Parallel()
	_, err := parse.NewParser().Parse("doesn't-exist.sysl", syslutil.NewChrootFs(afero.NewOsFs(), ""))
	require.Error(t, err)
}

func TestDoGenerateDataScript(t *testing.T) {
	args := &scriptArgs{
		root:    testDir + "db_scripts/",
		source:  "dataForSqlScriptOrg.sysl",
		output:  "%(epname).sql",
		project: "Project",
		title:   "Petstore Schema",
	}
	argsData := []string{"sysl", "generatescript", "-t", args.title, "-o", args.output,
		"-r", args.root, "-s", args.source, "-j", args.project}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "generate-script")
}

func TestDoConstructDatabaseScript(t *testing.T) {
	args := &scriptArgs{
		root:    testDir + "db_scripts/",
		source:  "dataForSqlScriptOrg.sysl",
		output:  "%(epname).sql",
		project: "Project",
		title:   "Petstore Schema",
		expected: map[string]string{
			"Relational-Model.sql": filepath.Join(testDir, "db_scripts/postgres-create-script-golden.sql"),
		},
	}
	result, err := DoConstructDatabaseScriptWithParams(args.root, "", args.title, args.output, args.project,
		args.source)
	assert.Nil(t, err, "Generating the sql script failed")
	compareSQL(t, args.expected, result)
}

func DoConstructDatabaseScriptWithParams(
	rootModel, filter, title, output, project, modules string,
) (map[string]string, error) {
	cmdDatabaseScript := &CmdDatabaseScript{
		title:   title,
		output:  output,
		root:    rootModel,
		source:  modules,
		project: project,
	}

	logger, _ := test.NewNullLogger()
	return GenerateDatabaseScripts(cmdDatabaseScript, logger)
}
