package database

import (
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDatabaseScriptCreate(t *testing.T) {
	goldenFileName := "db_scripts/postgres-create-script-golden.sql"
	modelParser := parse.NewParser()
	mod, _, err := parse.LoadAndGetDefaultApp("database/db_scripts/dataForSqlScriptOrg.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
	assert.Nil(t, err)
	types := mod.GetApps()[testAppName].GetTypes()
	v := MakeDatabaseScriptView(testTitle, logrus.StandardLogger())
	outputStr := v.GenerateDatabaseScriptCreate(types, testDBType, testAppName)
	CompareContent(t, goldenFileName, outputStr)
}

func TestGenerateDatabaseScriptModify(t *testing.T) {
	goldenFileName := "db_scripts/postgres-modify-script-golden.sql"
	modelParser := parse.NewParser()
	modOld, _, err := parse.LoadAndGetDefaultApp("database/db_scripts/dataForSqlScriptOrg.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
	assert.Nil(t, err)
	modNew, _, err := parse.LoadAndGetDefaultApp("database/db_scripts/dataForSqlScriptModified.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
	assert.Nil(t, err)
	appsOld := modOld.GetApps()
	appsNew := modNew.GetApps()
	appNames := strings.Split(testAppName, Delimiter)
	v := MakeDatabaseScriptView(testTitle, logrus.StandardLogger())
	outputStr := v.ProcessModSysls(appsOld, appsNew, appNames, "", testDBType)
	CompareContent(t, goldenFileName, outputStr[0].content)
}

func TestGenerateDatabaseScriptModifyTwoApps(t *testing.T) {
	expected := map[string]string{
		"RelModel.sql":    "db_scripts/postgres-modify-script-golden.sql",
		"RelModelNew.sql": "db_scripts/postgres-modify-script-golden_second_app.sql",
	}
	modelParser := parse.NewParser()
	modOld, _, err := parse.LoadAndGetDefaultApp("database/db_scripts/dataForSqlScriptOrg.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
	assert.Nil(t, err)
	modNew, _, err := parse.LoadAndGetDefaultApp("database/db_scripts/dataForSqlScriptModifiedTwoApps.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
	assert.Nil(t, err)
	appsOld := modOld.GetApps()
	appsNew := modNew.GetApps()
	appNames := strings.Split(testTwoAppNames, Delimiter)
	v := MakeDatabaseScriptView(testTitle, logrus.StandardLogger())
	output := v.ProcessModSysls(appsOld, appsNew, appNames, "", testDBType)
	CompareSQL(t, expected, output)
}
