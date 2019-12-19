package main

import (
	"path/filepath"
	"testing"

	db "github.com/anz-bank/sysl/pkg/database"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

type scriptModArgs struct {
	inputDir  string
	title     string
	outputDir string
	appNames  string
	orgSource string
	newSource string
	expected  map[string]string
}

func TestDoGenerateDataScriptMod(t *testing.T) {
	args := &scriptModArgs{
		inputDir:  testDir + "db_scripts/",
		orgSource: "dataForSqlScriptOrg.sysl",
		newSource: "dataForSqlScriptModifiedTwo.sysl",
		outputDir: testDir + "db_scripts/",
		appNames:  "RelModel,RelModelNew",
		title:     "Petstore Schema",
	}
	argsData := []string{"sysl", "generatedbscriptsdelta", "-t", args.title, "-o", args.outputDir,
		"-r", args.inputDir, "-s", args.newSource, "-n", args.orgSource, "-a", args.appNames}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "generate-db-scripts-delta")
}
func TestDoConstructDatabaseScriptMod(t *testing.T) {
	args := &scriptModArgs{
		inputDir:  testDir + "db_scripts/",
		orgSource: "dataForSqlScriptOrg.sysl",
		newSource: "dataForSqlScriptModified.sysl",
		outputDir: testDir + "db_scripts/",
		appNames:  "RelModel",
		title:     "Petstore Schema",
		expected: map[string]string{
			filepath.Join(testDir, "db_scripts/RelModel.sql"): filepath.Join(testDir,
				"db_scripts/postgres-modify-script-golden.sql"),
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.inputDir, "", args.title,
		args.outputDir, args.appNames, args.orgSource, args.newSource)
	assert.Nil(t, err, "Generating the sql script failed")
	db.CompareSQL(t, args.expected, result)
}

func TestDoConstructDatabaseScriptModTwoApps(t *testing.T) {
	args := &scriptModArgs{
		inputDir:  testDir + "db_scripts/",
		orgSource: "dataForSqlScriptOrg.sysl",
		newSource: "dataForSqlScriptModifiedTwoApps.sysl",
		outputDir: testDir + "db_scripts/",
		appNames:  "RelModel,RelModelNew",
		title:     "Petstore Schema",
		expected: map[string]string{
			filepath.Join(testDir, "db_scripts/RelModel.sql"): filepath.Join(testDir,
				"/db_scripts/postgres-modify-script-golden.sql"),
			filepath.Join(testDir, "db_scripts/RelModelNew.sql"): filepath.Join(testDir,
				"/db_scripts/postgres-modify-script-golden_second_app.sql"),
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.inputDir, "", args.title, args.outputDir,
		args.appNames, args.orgSource, args.newSource)
	assert.Nil(t, err, "Generating the sql script failed")
	db.CompareSQL(t, args.expected, result)
}

func DoConstructModDatabaseScriptWithParams(
	inputDir, filter, title, outputDir, appNames, orgSource, newSource string,
) ([]db.ScriptOutput, error) {
	cmdDatabaseScriptMod := &CmdDatabaseScriptMod{
		title:     title,
		inputDir:  inputDir,
		orgSource: orgSource,
		newSource: newSource,
		outputDir: outputDir,
		appNames:  appNames,
	}

	logger, _ := test.NewNullLogger()
	return GenerateModDatabaseScripts(cmdDatabaseScriptMod, logger)
}
