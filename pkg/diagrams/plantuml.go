package diagrams

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/spf13/afero"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	PlantUMLEnvVar  = "SYSL_PLANTUML"
	PlantUMLDefault = "http://localhost:8080/plantuml"
	testDir         = "../../tests/"
)

type Plantumlmixin struct {
	value string
}

func (p *Plantumlmixin) AddFlag(cmd *kingpin.CmdClause) {
	cmd.Flag("plantuml",
		"base url of plantuml server (default: "+PlantUMLEnvVar+" or "+
			PlantUMLDefault+" see "+
			"http://plantuml.com/server.html#install for more info)",
	).Short('p').StringVar(&p.value)
}

func (p *Plantumlmixin) Value() string {
	if p.value == "" {
		p.value = os.Getenv(PlantUMLEnvVar)
		if p.value == "" {
			p.value = PlantUMLDefault
		}
	}
	return p.value
}

func (p *Plantumlmixin) GenerateFromMap(m map[string]string, fs afero.Fs) error {
	for k, v := range m {
		if err := OutputPlantuml(k, p.Value(), v, fs); err != nil {
			return err
		}
	}
	return nil
}

func ComparePUML(t *testing.T, expected, actual map[string]string) {
	for name, goldenFile := range expected {
		golden, err := ioutil.ReadFile(goldenFile)
		assert.Nil(t, err)
		if string(golden) != actual[name] {
			err := ioutil.WriteFile(filepath.Join(testDir, name+".puml"), []byte(actual[name]), 0777)
			assert.Nil(t, err)
		}
		golden = syslutil.HandleCRLF(golden)
		assert.Equal(t, string(golden), actual[name])
	}

	// Then
	assert.Equal(t, len(expected), len(actual))
}
