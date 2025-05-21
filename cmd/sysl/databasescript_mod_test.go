package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/database"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

type scriptModArgs struct {
	title     string
	outputDir string
	appNames  string
	orgSource string
	newSource string
	expected  map[string]string
}

func TestDoGenerateDataScriptMod(t *testing.T) {
	t.Parallel()

	args := &scriptModArgs{
		orgSource: database.DBTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
		newSource: database.DBTestDir + "db_scripts/dataForSqlScriptModifiedTwo.sysl",
		outputDir: database.DBTestDir + "db_scripts/",
		appNames:  "RelModel,RelModelNew",
		title:     "Petstore Schema",
	}
	argsData := []string{"sysl", "generatedbscriptsdelta", "-t", args.title, "-o", args.outputDir,
		"-a", args.appNames, args.orgSource, args.newSource}
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := cmdRunner{}
	assert.NoError(t, r.Configure(syslCmd))
	selectedCommand, err := syslCmd.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl data")
	assert.Equal(t, selectedCommand, "generate-db-scripts-delta")
}

func TestModDBScriptValidSyslFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")

	main2([]string{"sysl", "generatedbscriptsdelta", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl"),
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptModifiedTwoApps.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/RelModel.sql")
}

func TestModDBScriptInValidOrgSyslFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatedbscriptsdelta", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DBTestDir, "db_scripts/invalid.sysl"),
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptModified.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 2, err)
}

func TestModDBScriptOneModule(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatedbscriptsdelta", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 1, err)
}

func TestModDBScriptThreeModule(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatedbscriptsdelta", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl"),
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptModified.sysl"),
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptModified.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 1, err)
}

func TestModDBScriptInValidNewSyslFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatedbscriptsdelta", "-t", "PetStore", "-o", "", "-a", "RelModel",
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl"),
		filepath.Join(database.DBTestDir, "db_scripts/invalid.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 2, err)
}

func TestModDBScriptNewAppSyslFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")

	err := main2([]string{"sysl", "generatedbscriptsdelta", "-t", "PetStore", "-o", "", "-a", "RelModel,RelModelNew",
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptOrg.sysl"),
		filepath.Join(database.DBTestDir, "db_scripts/dataForSqlScriptModifiedTwoApps.sysl")},
		fs, logger, os.Stdin, io.Discard, main3)
	assert.Equal(t, 0, err)
	syslutil.AssertFsHasExactly(t, memFs, "/RelModel.sql", "/RelModelNew.sql")
}

func TestDoConstructDatabaseScriptMod(t *testing.T) {
	t.Parallel()

	args := &scriptModArgs{
		orgSource: database.DBTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
		newSource: database.DBTestDir + "db_scripts/dataForSqlScriptModified.sysl",
		outputDir: database.DBTestDir + "db_scripts/",
		appNames:  "RelModel",
		title:     "Petstore Schema",
		expected: map[string]string{
			filepath.Join(database.DBTestDir, "db_scripts/RelModel.sql"): filepath.Join(database.DBTestDir,
				"db_scripts/postgres-modify-script-golden.sql"),
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.title,
		args.outputDir, args.appNames, args.orgSource, args.newSource)
	assert.Nil(t, err, "Generating the sql script failed")
	database.CompareSQL(t, args.expected, result)
}

func TestDoConstructDatabaseScriptModTwoApps(t *testing.T) {
	t.Parallel()

	args := &scriptModArgs{
		orgSource: database.DBTestDir + "db_scripts/dataForSqlScriptOrg.sysl",
		newSource: database.DBTestDir + "db_scripts/dataForSqlScriptModifiedTwoApps.sysl",
		outputDir: database.DBTestDir + "db_scripts/",
		appNames:  "RelModel,RelModelNew",
		title:     "Petstore Schema",
		expected: map[string]string{
			filepath.Join(database.DBTestDir, "db_scripts/RelModel.sql"): filepath.Join(database.DBTestDir,
				"/db_scripts/postgres-modify-script-golden.sql"),
			filepath.Join(database.DBTestDir, "db_scripts/RelModelNew.sql"): filepath.Join(database.DBTestDir,
				"/db_scripts/postgres-modify-script-golden_second_app.sql"),
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.title, args.outputDir,
		args.appNames, args.orgSource, args.newSource)
	assert.Nil(t, err, "Generating the sql script failed")
	database.CompareSQL(t, args.expected, result)
}

func DoConstructModDatabaseScriptWithParams(
	title, outputDir, appNames, orgSource, newSource string,
) ([]database.ScriptOutput, error) {
	cmdDatabaseScriptMod := &cmdutils.CmdDatabaseScriptParams{
		Title:     title,
		OutputDir: outputDir,
		AppNames:  appNames,
	}

	logger, _ := test.NewNullLogger()
	modelOld, _, err := loader.LoadSyslModule("", orgSource, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	modelNew, _, err := loader.LoadSyslModule("", newSource, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateModDatabaseScripts(cmdDatabaseScriptMod, modelOld, modelNew, logger)
}
