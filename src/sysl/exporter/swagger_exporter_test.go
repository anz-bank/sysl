package exporter

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/src/sysl/parse"
	"github.com/anz-bank/sysl/src/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestExportAll(t *testing.T) {
	t.Parallel()
	modelParser := parse.NewParser()
	const syslTestDir = "test-data"
	files, err := ioutil.ReadDir(syslTestDir)
	require.NoError(t, err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		parts := strings.Split(file.Name(), ".")
		if strings.EqualFold(parts[1], "sysl") {
			t.Run(parts[0], func(t *testing.T) {
				t.Parallel()
				mod, _, err := parse.LoadAndGetDefaultApp("exporter/test-data/"+parts[0]+`.sysl`,
					syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
				require.NoError(t, err)
				if err != nil {
					t.Errorf("Error reading sysl %s", parts[0]+`.sysl`)
				}
				swaggerExporter := MakeSwaggerExporter(mod.GetApps()["testapp"], logrus.StandardLogger())
				require.NoError(t, swaggerExporter.GenerateSwagger())
				out, err := swaggerExporter.SerializeOutput("yaml")
				require.NoError(t, err)
				yamlFileBytes, err := ioutil.ReadFile("../exporter/test-data/" + parts[0] + `.yaml`)
				require.NoError(t, err)

				yamlFileBytes = syslutil.HandleCRLF(yamlFileBytes)
				require.Equal(t, string(yamlFileBytes), string(out))
			})
		}
	}
}
