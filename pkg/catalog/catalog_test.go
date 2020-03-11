package catalog

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

var catalogTestFiles = []string{"grpc_catalog.sysl", "rest_catalog.sysl"}
var catalogTestAppNames = []string{"Greeter", "example"}
var catalogTestPaths = []string{"/grpc/Greeter/", "/rest/example"}

func TestMakeAPIDocBuilder(t *testing.T) {
	var apps []*sysl.Application
	for _, filename := range catalogTestFiles {
		appImports, err := importWebserviceFromFile(filename)
		if err != nil {
			t.Error("Failed to import webservice file")
		}
		apps = append(apps, appImports...)
	}

	for index, app := range apps {
		builder := MakeAPIDocBuilder(app, logrus.New(), false)
		newDoc, err := builder.BuildAPIDoc()

		if err != nil {
			t.Error("Failed to build API Doc")
		}

		if newDoc.name != catalogTestAppNames[index] {
			t.Errorf("Incorrect AppName, expected:%s got:%s", catalogTestAppNames[index], newDoc.name)
		}
	}
}

func TestBuildAPIDoc(t *testing.T) {
	var apps []*sysl.Application
	for _, filename := range catalogTestFiles {
		appImports, err := importWebserviceFromFile(filename)
		if err != nil {
			t.Error("Failed to import webservice file")
		}
		apps = append(apps, appImports...)
	}

	for index, app := range apps {
		builder := MakeAPIDocBuilder(app, &logrus.Logger{}, false)
		newDoc, err := builder.BuildAPIDoc()

		if err != nil {
			t.Error("Failed to build API Doc")
		}

		if newDoc.path != catalogTestPaths[index] {
			t.Errorf("Incorrect path, expected:%s got:%s", catalogTestPaths[index], newDoc.path)
		}
	}
}
