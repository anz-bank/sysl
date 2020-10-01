package endpointanalysisdiagram

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

func TestGenerateMermaidIntegrationDiagram(t *testing.T) {
	m, err := parse.NewParser().ParseFromFs("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateEndpointAnalysisDiagram(m)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram1(t *testing.T) {
	m, err := parse.NewParser().ParseFromFs("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateEndpointAnalysisDiagram(m)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram2(t *testing.T) {
	m, err := parse.NewParser().ParseFromFs("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateEndpointAnalysisDiagram(m)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram3(t *testing.T) {
	m, err := parse.NewParser().ParseFromFs("demo/simple/sysl-app-hyperlink.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateEndpointAnalysisDiagram(m)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidIntegrationDiagram4(t *testing.T) {
	m, err := parse.NewParser().ParseFromFs("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	if err != nil {
		t.Error(err)
	}
	apps := []string{"TheWorld", "DragonEater"}
	r, err := GenerateMultipleAppEndpointAnalysisDiagram(m, apps)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidEPAWithNamespace(t *testing.T) {
	t.Parallel()
	apps := []string{"Org :: Team :: TestApp", "Org :: Team :: TestApp2"}
	syslInput := `
Org :: Team :: TestApp:
	myEndPoint1:
		Org :: Team :: TestApp2 <- myEndPoint2
		return ok <: string

Org :: Team :: TestApp2:
	myEndPoint2:
		return ok <: string
	`
	expected := `Org__Team__TestApp2-myEndPoint2 --> myEndPoint2`
	m, err := parse.NewParser().ParseString(syslInput)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	r, err := GenerateMultipleAppEndpointAnalysisDiagram(m, apps)
	assert.NotEmpty(t, r)
	assert.NoError(t, err)
	assert.Contains(t, r, expected)
}
