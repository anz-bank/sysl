package importer

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/require"
)

func TestLoadSwaggerFromTestFiles(t *testing.T) {
	const fileTestDir = "tests-swagger"
	files, err := ioutil.ReadDir(fileTestDir)
	require.NoError(t, err)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		logger, _ := test.NewNullLogger()

		parts := strings.Split(f.Name(), ".")
		if strings.EqualFold(parts[1], "yaml") {
			filename := strings.Join(parts[:len(parts)-1], ".")
			t.Run(fmt.Sprintf("TestLoadSwaggerFromTestFiles - %s", filename), func(t *testing.T) {

				input, err := ioutil.ReadFile(path.Join(fileTestDir, filename+".yaml"))
				require.NoError(t, err)
				expected, err := ioutil.ReadFile(path.Join(fileTestDir, filename+".sysl"))
				require.NoError(t, err)

				result, err := LoadSwaggerText(OutputData{
					AppName: "testapp",
					Package: "package_foo",
				}, string(input), logger)
				require.NoError(t, err)
				require.Equal(t, string(expected), result)

			})
		}
	}

}
