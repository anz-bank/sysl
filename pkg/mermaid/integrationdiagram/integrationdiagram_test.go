package integrationdiagram

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

func TestBadInputsToGenerateMermaidIntegrationDiagram(t *testing.T) {
	appName := "wrongName"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, appName)
	assert.NotNil(t, m)
	assert.Empty(t, r)
	assert.Error(t, err)
}

func TestGenerateMermaidIntegrationDiagram(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "WebFrontend")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestGenerateMermaidIntegrationDiagram2(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "WebFrontend")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestGenerateMermaidIntegrationDiagram3(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/sysl-ints.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "IntegratedSystem")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestGenerateMermaidIntegrationDiagram4(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "DragonEater")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}
