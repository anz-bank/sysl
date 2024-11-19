package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

type scriptArgs struct {
	title     string
	outputDir string
	source    string
	appNames  string
	expected  map[string]string
}

func TestDoGenerateDataScript(t *testing.T) {
	t.Parallel()

	args := &scriptArgs{
		source:    database.DBTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
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
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl"),
		"-a", "RelModel"},
		fs, logger, os.Stdin, io.Discard, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/RelModel.sql")
}

func TestCreateDBScriptInValidSyslFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatescript", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DBTestDir, "db_scripts/invalid.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 1, err)
}

func TestCreateDBScriptNoAppSyslFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatescript", "-t", "PetStore", "-o", "", "-a", "Proj123",
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 1, err)
}

func TestDoConstructDatabaseScript(t *testing.T) {
	t.Parallel()

	args := &scriptArgs{
		source:    database.DBTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
		outputDir: database.DBTestDir + "db_scripts/",
		title:     "Petstore Schema",
		appNames:  "RelModel",
		expected: map[string]string{
			filepath.Join(database.DBTestDir, "db_scripts/RelModel.sql"): filepath.Join(database.DBTestDir,
				"db_scripts/postgres-create-script-golden.sql"),
		},
	}
	result, err := DoConstructDatabaseScriptWithParams(args.title, args.outputDir,
		args.appNames, args.source)
	assert.Nil(t, err, "Generating the sql script failed")
	database.CompareSQL(t, args.expected, result)
}

func TestDoConstructDatabaseScriptInvalidFile(t *testing.T) {
	t.Parallel()

	args := &scriptArgs{
		source:    database.DBTestDir + "db_scripts/invalid.sysl",
		outputDir: database.DBTestDir + "db_scripts/",
		title:     "Petstore Schema",
		appNames:  "RelModel",
	}
	_, err := DoConstructDatabaseScriptWithParams(args.title, args.outputDir,
		args.appNames, args.source)
	actualErr, isParseExit := err.(syslutil.Exit)
	assert.True(t, isParseExit)
	assert.Equal(t, 2, actualErr.Code)
	assert.True(t, strings.Contains(actualErr.Error(), "invalid.sysl has syntax errors\n"))
}

func DoConstructDatabaseScriptWithParams(
	title, output, appNames, source string,
) ([]database.ScriptOutput, error) {
	cmdDatabaseScript := &cmdutils.CmdDatabaseScriptParams{
		Title:     title,
		OutputDir: output,
		AppNames:  appNames,
	}

	logger, _ := test.NewNullLogger()
	mod, _, err := loader.LoadSyslModule("", source, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateDatabaseScripts(cmdDatabaseScript, mod, logger)
}
