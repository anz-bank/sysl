package integrationdiagram

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

func TestBadInputsToGenerateMermaidIntegrationDiagram(t *testing.T) {
	appName := "wrongName"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
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
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "WebFrontend")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram2(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "WebFrontend")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram3(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/sysl-ints.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "IntegratedSystem")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram4(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateIntegrationDiagram(m, "DragonEater")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagramWithoutApp(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateFullIntegrationDiagram(m)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}
func TestGenerateMermaidIntegrationDiagraMultipleApps(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	apps := []string{"TheWorld", "DragonEater"}
	r, err := GenerateMultipleAppIntegrationDiagram(m, apps)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagraMultipleAppsWithUnrelated(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/sizzle/sizzle.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	assert.NoError(t, err)
	apps := []string{"Visa", "MasterCard", "ShouldNotDisplay"}
	r, err := GenerateMultipleAppIntegrationDiagram(m, apps)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(r, "Visa"))
	assert.True(t, strings.Contains(r, "MasterCard"))
}

func TestGenerateMermaidIntegrationDiagraMultipleAppsDependantApps(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/sizzle/sizzle.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	assert.NoError(t, err)
	apps := []string{"MegaDatabase"}
	r, err := GenerateMultipleAppIntegrationDiagram(m, apps)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(r, "MegaDatabase"))
	assert.True(t, strings.Contains(r, ` IdentityServer["IdentityServer"] --> MegaDatabase["MegaDatabase"]`))
}

func TestGenerateMermaidIntegrationDiagraMultipleAppsDependantAppsCallInConditional(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/sizzle/sizzle.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	assert.NoError(t, err)
	apps := []string{"Visa"}
	r, err := GenerateMultipleAppIntegrationDiagram(m, apps)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(r, "Visa"))
	assert.True(t, strings.Contains(r, `PaymentServer["PaymentServer"] --> Visa["Visa"]`))
}

func TestPrintClassStatementWithSpaces(t *testing.T) {
	r := printClassStatement("my application")
	expected := "    my_application\n"
	assert.Equal(t, expected, r)
}
