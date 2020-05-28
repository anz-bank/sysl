package datamodeldiagram

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
	"github.com/spf13/afero"
)

func TestGenerateMermaidDataModelDiagram(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/petshop/petshop.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	assert.NoError(t, err)
	mapper := syslwrapper.MakeAppMapper(m)
	mapper.IndexTypes()
	mapper.ConvertTypes()
	r, err := GenerateFullDataDiagram(m)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}

func TestGenerateMermaidDataModelDiagramWithAppAndType(t *testing.T) {
	m, err := parse.NewParser().Parse("demo/petshop/petshop.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), mermaid.ProjectDir))
	assert.NoError(t, err)
	r, err := GenerateDataDiagramWithAppAndType(m, "PetShopModel", "Pet")
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.NoError(t, err)
}
