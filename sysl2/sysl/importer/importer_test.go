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

type testConfig struct {
	name          string
	testDir       string
	testExtension string
	fn            Func
}

func runImportEqualityTests(t *testing.T, cfg testConfig) {
	t.Parallel()
	files, err := ioutil.ReadDir(cfg.testDir)
	require.NoError(t, err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		logger, _ := test.NewNullLogger()

		parts := strings.Split(f.Name(), ".")
		if strings.EqualFold(parts[1], cfg.testExtension) {
			filename := strings.Join(parts[:len(parts)-1], ".")
			t.Run(fmt.Sprintf("%s - %s", cfg.name, filename), func(t *testing.T) {
				input, err := ioutil.ReadFile(path.Join(cfg.testDir, filename+"."+cfg.testExtension))
				require.NoError(t, err)
				expected, err := ioutil.ReadFile(path.Join(cfg.testDir, filename+".sysl"))
				require.NoError(t, err)

				result, err := cfg.fn(OutputData{AppName: "testapp", Package: "package_foo"}, string(input), logger)
				require.NoError(t, err)
				require.Equal(t, string(expected), result)
			})
		}
	}
}

func TestLoadSwaggerFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadSwaggerFromTestFiles",
		testDir:       "tests-swagger",
		testExtension: "yaml",
		fn:            LoadSwaggerText,
	})
}

func TestLoadXSDFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadXSDFromTestFiles",
		testDir:       "tests-xsd",
		testExtension: "xsd",
		fn:            LoadXSDText,
	})
}
