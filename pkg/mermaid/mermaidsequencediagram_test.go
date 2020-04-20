package mermaid

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

func TestBadInputsToGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appname := "wrongname"
	epname := "wrongep"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.Empty(t, r)
	assert.Error(t, err)
}

func TestGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appname := "WebFrontend"
	epname := "RequestProfile"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestGenerateMermaidSequenceDiagram2(t *testing.T) {
	t.Parallel()
	appname := "WebFrontend"
	epname := "RequestProfile"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestGenerateMermaidSequenceDiagram3(t *testing.T) {
	t.Parallel()
	appname := "MobileApp"
	epname := "Login"
	m, err := parse.NewParser().Parse("demo/simple/sysl-app-hyperlink.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestGenerateMermaidSequenceDiagramWithIfElseLoopActionAndGroupStatements(t *testing.T) {
	t.Parallel()
	appname := "DragonEater"
	epname := "EatDragon"
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}
