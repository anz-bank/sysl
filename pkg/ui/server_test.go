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
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode, "expected status code to be 200 but got %d", result.StatusCode)
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

	if _, err := syslUI.GenerateServer(); err == nil {
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
	if err = server.Setup(); err != nil {
		t.Errorf("Error running Setup on server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/rest/spec/example", nil)
	w := httptest.NewRecorder()
	server.handleAPIDoc(w, req)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode, "Expected return status code of 200")
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
	if err = server.Setup(); err != nil {
		t.Errorf("Error running Setup on server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/services.json", nil)
	w := httptest.NewRecorder()
	server.handleJSONServices(w, req)
	result := w.Result()
	defer result.Body.Close()
	jsonResponse, err := simplejson.NewFromReader(result.Body)
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
	if err := server.Setup(); err != nil {
		t.Errorf("Error running Setup on server %s", err)
	}
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/services.json", nil)
	w := httptest.NewRecorder()
	defer w.Result().Body.Close()
	server.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode, "expected status code to be 200 but got %d", result.StatusCode)
}
