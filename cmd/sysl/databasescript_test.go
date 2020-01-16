package main

import (
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

type scriptArgs struct {
	title     string
	outputDir string
	source    string
	appNames  string
	expected  map[string]string
}

func TestDoGenerateDataScript(t *testing.T) {
	args := &scriptArgs{
		source:    database.DbTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
		outputDir: "",
		title:     "Petstore Schema",
		appNames:  "RelModel",
	}
	argsData := []string{"sysl", "generatedbscripts", "-t", args.title, "-o", args.outputDir,
		"-a", args.appNames, args.source}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "generate-db-scripts")
}

func TestCreateDBScriptValidSyslFile(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")

	main2([]string{"sysl", "generatedbscripts", "-t", "PetStore",
		"-o", "",
		filepath.Join(database.DbTestDir, "db_scripts/dataForSqlScriptOrg.sysl"),
		"-a", "RelModel"},
		fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/RelModel.sql")
}

func TestCreateDBScriptInValidSyslFile(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatescript", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DbTestDir, "db_scripts/invalid.sysl")},
		fs, logger, main3)
	assert.Equal(t, 1, err)
}

func TestCreateDBScriptNoAppSyslFile(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatescript", "-t", "PetStore", "-o", "", "-a", "Proj123",
		filepath.Join(database.DbTestDir, "db_scripts/dataForSqlScriptOrg.sysl")},
		fs, logger, main3)
	assert.Equal(t, 1, err)
}

func TestDoConstructDatabaseScript(t *testing.T) {
	args := &scriptArgs{
		source:    database.DbTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
		outputDir: database.DbTestDir + "db_scripts/",
		title:     "Petstore Schema",
		appNames:  "RelModel",
		expected: map[string]string{
			filepath.Join(database.DbTestDir, "db_scripts/RelModel.sql"): filepath.Join(database.DbTestDir,
				"db_scripts/postgres-create-script-golden.sql"),
		},
	}
	result, err := DoConstructDatabaseScriptWithParams("", args.title, args.outputDir,
		args.appNames, args.source)
	assert.Nil(t, err, "Generating the sql script failed")
	database.CompareSQL(t, args.expected, result)
}

func TestDoConstructDatabaseScriptInvalidFile(t *testing.T) {
	args := &scriptArgs{
		source:    database.DbTestDir + "db_scripts/invalid.sysl",
		outputDir: database.DbTestDir + "db_scripts/",
		title:     "Petstore Schema",
		appNames:  "RelModel",
	}
	_, err := DoConstructDatabaseScriptWithParams("", args.title, args.outputDir,
		args.appNames, args.source)
	expectedError := parse.Exitf(2, "invalid.sysl has syntax errors\n")
	assert.Equal(t, expectedError, err)
}

func DoConstructDatabaseScriptWithParams(
	filter, title, output, appNames, source string,
) ([]database.ScriptOutput, error) {
	cmdDatabaseScript := &CmdDatabaseScriptParams{
		title:     title,
		outputDir: output,
		appNames:  appNames,
	}

	logger, _ := test.NewNullLogger()
	mod, _, err := LoadSyslModule("", source, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateDatabaseScripts(cmdDatabaseScript, mod, logger)
}
