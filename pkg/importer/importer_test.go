package importer

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl/syslutil"
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
	files, err := ioutil.ReadDir(cfg.testDir)
	require.NoError(t, err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		logger, _ := test.NewNullLogger()
		ext := filepath.Ext(f.Name())
		if strings.EqualFold(ext[1:], cfg.testExtension) {
			filename := strings.TrimSuffix(f.Name(), ext)
			t.Run(fmt.Sprintf("%s - %s", cfg.name, filename), func(t *testing.T) {
				t.Parallel()
				input, err := ioutil.ReadFile(path.Join(cfg.testDir, filename+"."+cfg.testExtension))
				require.NoError(t, err)
				syslFile := path.Join(cfg.testDir, filename+".sysl")
				expected, err := ioutil.ReadFile(syslFile)
				require.NoError(t, err)
				expected = syslutil.HandleCRLF(expected)

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

func TestLoadGrammarFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadGrammarFromTestFiles",
		testDir:       "tests-grammar",
		testExtension: "g",
		fn:            LoadGrammar,
	})
}
