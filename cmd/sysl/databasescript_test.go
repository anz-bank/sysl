package main

import (
	"path/filepath"
	"testing"

	db "github.com/anz-bank/sysl/pkg/database"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

type scriptArgs struct {
	inputDir  string
	title     string
	outputDir string
	source    string
	appNames  string
	expected  map[string]string
}

func TestDoGenerateDataScript(t *testing.T) {
	args := &scriptArgs{
		inputDir:  testDir + "db_scripts/",
		source:    "dataForSqlScriptOrg.sysl",
		outputDir: testDir + "db_scripts/",
		title:     "Petstore Schema",
		appNames:  "RelModel",
	}
	argsData := []string{"sysl", "generatedbscripts", "-t", args.title, "-o", args.outputDir,
		"-i", args.inputDir, "-s", args.source, "-a", args.appNames}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "generate-db-scripts")
}

func TestDoConstructDatabaseScript(t *testing.T) {
	args := &scriptArgs{
		inputDir:  testDir + "db_scripts/",
		source:    "dataForSqlScriptOrg.sysl",
		outputDir: testDir + "db_scripts/",
		title:     "Petstore Schema",
		appNames:  "RelModel",
		expected: map[string]string{
			filepath.Join(testDir, "db_scripts/RelModel.sql"): filepath.Join(testDir,
				"db_scripts/postgres-create-script-golden.sql"),
		},
	}
	result, err := DoConstructDatabaseScriptWithParams(args.inputDir, "", args.title, args.outputDir,
		args.appNames, args.source)
	assert.Nil(t, err, "Generating the sql script failed")
	db.CompareSQL(t, args.expected, result)
}

func DoConstructDatabaseScriptWithParams(
	inputDir, filter, title, output, appNames, source string,
) ([]db.ScriptOutput, error) {
	cmdDatabaseScript := &CmdDatabaseScript{
		title:     title,
		outputDir: output,
		inputDir:  inputDir,
		source:    source,
		appNames:  appNames,
	}

	logger, _ := test.NewNullLogger()
	return GenerateDatabaseScripts(cmdDatabaseScript, logger)
}
