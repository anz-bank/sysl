package exporter

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var (
	update = flag.Bool("update", false, "Update golden test files")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestExportAll(t *testing.T) {
	t.Parallel()
	modelParser := parse.NewParser()
	const syslTestDir = "test-data/openapi2/"
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
				syslFileName := "exporter/" + syslTestDir + parts[0] + `.sysl`
				openAPIFileName := "../exporter/" + syslTestDir + parts[0] + `.yaml`
				mod, _, err := parse.LoadAndGetDefaultApp(syslFileName,
					syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
				require.NoError(t, err)
				if err != nil {
					t.Errorf("Error reading sysl %s", parts[0]+`.sysl`)
				}
				swaggerExporter := MakeSwaggerExporter(mod.GetApps()["testapp"], logrus.StandardLogger())
				require.NoError(t, swaggerExporter.GenerateSwagger())
				out, err := swaggerExporter.SerializeOutput("yaml")
				require.NoError(t, err)
				if *update {
					err = ioutil.WriteFile(openAPIFileName, out, 0644)
					if err != nil {
						t.Error(err)
					}
				}
				yamlFileBytes, err := ioutil.ReadFile(openAPIFileName)
				require.NoError(t, err)

				yamlFileBytes = syslutil.HandleCRLF(yamlFileBytes)
				require.Equal(t, string(yamlFileBytes), string(out))
			})
		}
	}
}
