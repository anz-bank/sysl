package importer

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus/hooks/test"
)

func TestLoadXsdFromTestFiles(t *testing.T) {
	const fileTestDir = "tests-xsd"
	files, err := ioutil.ReadDir(fileTestDir)
	require.NoError(t, err)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		logger, _ := test.NewNullLogger()

		parts := strings.Split(f.Name(), ".")
		if strings.EqualFold(parts[1], "xsd") {
			filename := strings.Join(parts[:len(parts)-1], ".")
			t.Run(fmt.Sprintf("TestLoadXsdFromTestFiles - %s", filename), func(t *testing.T) {
				input, err := ioutil.ReadFile(path.Join(fileTestDir, filename+".xsd"))
				require.NoError(t, err)
				expected, err := ioutil.ReadFile(path.Join(fileTestDir, filename+".sysl"))
				require.NoError(t, err)

				result, err := LoadXSDText(OutputData{
					AppName: "testapp",
					Package: "package_foo",
				}, string(input), "", logger)
				require.NoError(t, err)
				require.Equal(t, string(expected), result)
			})
		}
	}
}
