package exporter

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/roothandler"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
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
				root := ".."
				rootHandler, err := roothandler.NewRootHandler(root, "/exporter/test-data/"+parts[0]+`.sysl`,
					afero.NewOsFs(), logrus.StandardLogger())
				require.NoError(t, err)
				mod, _, err := parse.LoadAndGetDefaultApp(rootHandler,
					syslutil.NewChrootFs(afero.NewOsFs(), rootHandler.Root()), modelParser)
				require.NoError(t, err)
				if err != nil {
					t.Errorf("Error reading sysl %s", parts[0]+`.sysl`)
				}
				swaggerExporter := MakeSwaggerExporter(mod.GetApps()["testapp"], logrus.StandardLogger())
				require.NoError(t, swaggerExporter.GenerateSwagger())
				out, err := swaggerExporter.SerializeToYaml()
				require.NoError(t, err)
				yamlFileBytes, err := ioutil.ReadFile("../exporter/test-data/" + parts[0] + `.yaml`)
				require.NoError(t, err)
				require.Equal(t, string(yamlFileBytes), string(out))
			})
		}
	}
}
