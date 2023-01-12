package diagrams

import (
	"os"
	"path/filepath"
	"strings"
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
		golden, err := os.ReadFile(goldenFile)
		assert.Nil(t, err)
		if string(golden) != actual[name] {
			err := os.WriteFile(filepath.Join(testDir, name+".puml"), []byte(actual[name]), 0600)
			assert.Nil(t, err)
		}
		a := strings.TrimSuffix(string(syslutil.HandleCRLF(golden)), "\n")
		b := strings.TrimSuffix(actual[name], "\n")
		assert.Equal(t, a, b, "PlantUML not equal for %s", name)
	}

	// Then
	assert.Len(t, actual, len(expected))
}
