package catalog

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

const testDir = "../../tests/"

var filenames = []string{"grpc_catalog.sysl", "rest_catalog.sysl"}
var testAppNames = []string{"Greeter", "example"}
var testAppAttrKeys = []string{"owner", "owner"}
var testAppAttrVals = []string{"foo", "josh"}

func TestImport(t *testing.T) {
	var apps []*sysl.Application

	for _, filename := range filenames {
		appImports, err := importAppsFromFile(filename)
		if err != nil {
			t.Error("Failed to import webservice file")
		}
		apps = append(apps, appImports...)
	}

	for index, app := range apps {
		service, err := BuildWebService(app)
		if err != nil {
			t.Error("Failed to import webservice file")
		}
		if service.Name != testAppNames[index] {
			t.Errorf("Incorrect AppName, expected:%s got:%s", testAppNames[index], service.Name)
		}
		if service.Attributes[testAppAttrKeys[index]] != testAppAttrVals[index] {
			t.Errorf("Incorrect Attribute, expected:%s got:%s",
				testAppAttrVals[index],
				service.Attributes[testAppAttrKeys[index]])
		}
	}
}

func importAppsFromFile(filename string) ([]*sysl.Application, error) {
	var apps []*sysl.Application
	module, err := parse.NewParser().Parse(filename,
		syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	if err != nil {
		return nil, err
	}

	for _, a := range module.GetApps() {
		apps = append(apps, a)
	}
	return apps, nil
}
