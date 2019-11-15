package importer

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLoadGrammar(t *testing.T) {
	t.Parallel()
	const testGrammarImpDir = "tests-grammar"
	files, err := ioutil.ReadDir(testGrammarImpDir)
	require.NoError(t, err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		parts := strings.Split(file.Name(), ".")
		if len(parts) > 2 && strings.EqualFold(parts[1], "gen") && strings.EqualFold(parts[2], "g") {
			t.Run(parts[0], func(t *testing.T) {
				t.Parallel()
				outArgs := OutputData{
					AppName:   strings.Title(strings.ToLower(parts[0])),
					StartRule: parts[0] + "File",
				}
				data, err := ioutil.ReadFile("../importer/tests-grammar/" + parts[0] + ".gen.g")
				if err != nil {
					t.Errorf("Error reading grammar %s", parts[0]+".gen.g")
				}
				var outSysl string
				outSysl, err = LoadGrammar(outArgs, string(data), logrus.StandardLogger())
				if err != nil {
					t.Errorf("Error generating sysl %s", parts[0]+".gen.g")
				}
				syslGenBytes, err := ioutil.ReadFile("../importer/tests-grammar/" + parts[0] + ".gen.sysl")
				require.NoError(t, err)
				require.Equal(t, string(syslGenBytes), outSysl)
			})
		}
	}
}
