package ui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/bitly/go-simplejson"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var uiFields = `
team.slack,
description,
deploy.env1.url,
deploy.sit1.url,
deploy.sit2.url,
deploy.qa.url,
deploy.prod.url,
repo.url,
docs.url,
type`

func TestGenerateServer(t *testing.T) {
	module, err := parse.NewParser().Parse("ui_rest.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), "../../tests/"))
	if err != nil {
		t.Errorf("Error parsing test modules %s", err)
	}

	modules := []*sysl.Module{module}

	syslUI := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(uiFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	server, err := syslUI.GenerateServer()
	if err != nil {
		t.Errorf("Error generating server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/rest/spec/example", nil)
	w := httptest.NewRecorder()
	server.handleRestSpec(w, req)
	if w.Result().StatusCode != 200 {
		t.Errorf("Not returning 200")
	}
}

func TestGenerateServerHandlesEmptyArray(t *testing.T) {
	modules := []*sysl.Module{}

	syslUI := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(uiFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	_, err := syslUI.GenerateServer()

	if err == nil {
		t.Error("Empty input array not caught")
	}
}

func TestServerSetupRuns(t *testing.T) {
	module, err := parse.NewParser().Parse("ui_rest.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), "../../tests/"))
	if err != nil {
		t.Errorf("Error parsing test modules %s", err)
	}

	modules := []*sysl.Module{module}

	syslUI := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(uiFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	server, err := syslUI.GenerateServer()

	if err != nil {
		t.Errorf("Error generating server %s", err)
	}

	err = server.Setup()
	if err != nil {
		t.Errorf("Error running Setup on server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/rest/spec/example", nil)
	w := httptest.NewRecorder()
	server.handleAPIDoc(w, req)
	assert.Equal(t, 200, w.Result().StatusCode, "Expected return status code of 200")
}

func TestHandleJSONServices(t *testing.T) {
	module, err := parse.NewParser().Parse("ui_rest.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), "../../tests/"))
	if err != nil {
		t.Errorf("Error parsing test modules %s", err)
	}

	modules := []*sysl.Module{module}

	syslUI := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(uiFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	server, err := syslUI.GenerateServer()

	if err != nil {
		t.Errorf("Error generating server %s", err)
	}

	err = server.Setup()
	if err != nil {
		t.Errorf("Error running Setup on server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/services.json", nil)
	w := httptest.NewRecorder()
	server.handleJSONServices(w, req)
	jsonResponse, err := simplejson.NewFromReader(w.Result().Body) //nolint:bodycloser
	if err != nil {
		t.Errorf("Error parsing JSON")
	}
	name := jsonResponse.GetIndex(0).Get("Name").MustString()
	assert.Equal(t, "example", name, "expected example but got %s", name)
}

func TestServeHTTP(t *testing.T) {
	module, err := parse.NewParser().Parse("ui_rest.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), "../../tests/"))
	if err != nil {
		t.Errorf("Error parsing test modules %s", err)
	}

	modules := []*sysl.Module{module}

	syslUI := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(uiFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	server, err := syslUI.GenerateServer()

	if err != nil {
		t.Errorf("Error generating server %s", err)
	}

	err = server.Setup()
	if err != nil {
		t.Errorf("Error running Setup on server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/services.json", nil)
	w := httptest.NewRecorder()
	defer w.Result().Body.Close()
	server.ServeHTTP(w, req)
	assert.Equal(t, w.Result().StatusCode, 200)
}
