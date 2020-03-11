package catalog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var catalogFields = `
team,
team.slack,
owner.name,
owner.email,
file.version,
release.version,
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
	module, err := parse.NewParser().Parse("rest_catalog.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), "../../tests/"))
	if err != nil {
		t.Errorf("Error parsing test modules %s", err)
	}

	modules := []*sysl.Module{module}

	syslCatalog := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(catalogFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	server, err := syslCatalog.GenerateServer()
	if err != nil {
		t.Errorf("Error generating server %s", err)
	}

	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/rest/spec/example", nil)
	w := httptest.NewRecorder()
	server.handleRestSpec(w, req)
	//nolint:bodyclose
	if w.Result().StatusCode != 200 {
		t.Errorf("Not returning 200")
	}
}

func TestGenerateServerHandlesEmptyArray(t *testing.T) {
	modules := []*sysl.Module{}

	syslCatalog := SyslUI{
		Host:    "localhost:8080",
		Fields:  strings.Split(catalogFields, ","),
		Fs:      afero.NewOsFs(),
		Log:     logrus.New(),
		Modules: modules,
	}

	_, err := syslCatalog.GenerateServer()

	if err == nil {
		t.Error("Empty input array not caught")
	}
}
