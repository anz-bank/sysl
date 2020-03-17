package ui

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

var uiTestFiles = []string{"ui_grpc.sysl", "ui_rest.sysl"}
var uiTestAppNames = []string{"Greeter", "example"}
var uiTestPaths = []string{"/grpc/Greeter/", "/rest/example"}

func TestMakeAPIDocBuilder(t *testing.T) {
	var apps []*sysl.Application
	for _, filename := range uiTestFiles {
		appImports, err := importAppsFromFile(filename)
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

		if newDoc.name != uiTestAppNames[index] {
			t.Errorf("Incorrect AppName, expected:%s got:%s", uiTestAppNames[index], newDoc.name)
		}
	}
}

func TestBuildAPIDoc(t *testing.T) {
	var apps []*sysl.Application
	for _, filename := range uiTestFiles {
		appImports, err := importAppsFromFile(filename)
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
		if newDoc.path != uiTestPaths[index] {
			t.Errorf("Incorrect path, expected:%s got:%s", uiTestPaths[index], newDoc.path)
		}
	}
}

func TestBuildAPIDocWithInaccessibleGRPC(t *testing.T) {
	app, err := importAppsFromFile(uiTestFiles[0])
	if err != nil {
		t.Error("Failed to import webservice file")
	}

	builder := MakeAPIDocBuilder(app[0], &logrus.Logger{}, true)
	_, err = builder.BuildAPIDoc()
	if err == nil {
		t.Errorf("Error not raised for nonexistent grpc URL")
	}
}
