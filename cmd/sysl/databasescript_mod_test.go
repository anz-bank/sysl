package main

import (
	"github.com/sirupsen/logrus/hooks/test"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type scriptModArgs struct {
	root       string
	title      string
	output     string
	project    string
	modulesOld string
	modulesNew string
	expected   map[string]string
}

func TestDoConstructDatabaseScriptMod(t *testing.T) {
	args := &scriptModArgs{
		root:       "./tests/",
		modulesOld: "dataForSqlScriptOrg.sysl",
		modulesNew: "dataForSqlScriptModified.sysl",
		output:     "%(epname).sql",
		project:    "Project",
		title:      "Petstore Schema",
		expected: map[string]string{
			"Relational-Model.sql": "tests/postgres-modify-script-golden.sql",
		},
	}
	result, err := DoConstructModDatabaseScriptWithParams(args.root, "", args.title, args.output, args.project,
		args.modulesOld, args.modulesNew)
	assert.Nil(t, err, "Generating the sql script failed")
	compareSQL(t, args.expected, result)
}

func DoConstructModDatabaseScriptWithParams(
	rootModel, filter, title, output, project, modulesOld, modulesNew string,
) (map[string]string, error) {
	classFormat := databaseScriptHeader + "(classname)"
	cmdContextParamDatagen := &CmdContextParamDatagen{
		filter:      filter,
		title:       title,
		output:      output,
		project:     project,
		classFormat: classFormat,
	}

	logger, _ := test.NewNullLogger()
	modOld, _, err1 := LoadSyslModule(rootModel, modulesOld, afero.NewOsFs(), logger)
	modNew, _, err2 := LoadSyslModule(rootModel, modulesNew, afero.NewOsFs(), logger)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return GenerateModDatabaseScripts(cmdContextParamDatagen, modOld, modNew, logger)
}
