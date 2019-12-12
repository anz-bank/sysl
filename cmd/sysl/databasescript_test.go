package main

import (
	"github.com/sirupsen/logrus/hooks/test"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type scriptArgs struct {
	root     string
	title    string
	output   string
	project  string
	modules  string
	expected map[string]string
}

func TestDoConstructDatabaseScript(t *testing.T) {
	args := &scriptArgs{
		root:    "./tests/",
		modules: "dataForSqlScriptOrg.sysl",
		output:  "%(epname).sql",
		project: "Project",
		title:   "Petstore Schema",
		expected: map[string]string{
			"Relational-Model.sql": "tests/postgres-create-script-golden.sql",
		},
	}
	result, err := DoConstructDatabaseScriptWithParams(args.root, "", args.title, args.output, args.project,
		args.modules)
	assert.Nil(t, err, "Generating the sql script failed")
	compareSQL(t, args.expected, result)
}

func DoConstructDatabaseScriptWithParams(
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
	return GenerateDatabaseScripts(cmdContextParamDatagen, mod, logger)
}
