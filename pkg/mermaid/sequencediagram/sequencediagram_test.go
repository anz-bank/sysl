package sequencediagram

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

func TestBadInputsToGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appName := "wrongName"
	epName := "wrongEp"
	m, err := parse.NewParser().ParseFs("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateSequenceDiagram(m, appName, epName)
	assert.NotNil(t, m)
	assert.Empty(t, r)
	assert.Error(t, err)
}

func TestGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appName := "WebFrontend"
	epName := "RequestProfile"
	m, err := parse.NewParser().ParseFs("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateSequenceDiagram(m, appName, epName)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidSequenceDiagram2(t *testing.T) {
	t.Parallel()
	appName := "WebFrontend"
	epName := "RequestProfile"
	m, err := parse.NewParser().ParseFs("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateSequenceDiagram(m, appName, epName)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidSequenceDiagram3(t *testing.T) {
	t.Parallel()
	appName := "MobileApp"
	epName := "Login"
	m, err := parse.NewParser().ParseFs("demo/simple/sysl-app-hyperlink.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateSequenceDiagram(m, appName, epName)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidSequenceDiagramWithIfElseLoopActionAndGroupStatements(t *testing.T) {
	t.Parallel()
	appName := "DragonEater"
	epName := "EatDragon"
	m, err := parse.NewParser().ParseFs("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projectDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateSequenceDiagram(m, appName, epName)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}
