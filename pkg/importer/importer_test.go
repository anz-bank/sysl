package importer

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/require"
)

var (
	update = flag.Bool("update", false, "Update golden test files")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

type testConfig struct {
	name          string
	testDir       string
	testExtension string
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
		if strings.EqualFold(ext, cfg.testExtension) {
			filename := strings.TrimSuffix(f.Name(), ext)
			t.Run(fmt.Sprintf("%s-%s", cfg.name, filename), func(t *testing.T) {
				t.Parallel()
				syslFile := filepath.Join(cfg.testDir, filename+".sysl")
				fileToImport := syslutil.MustAbsolute(t, filepath.Join(cfg.testDir, filename+cfg.testExtension))
				input, err := ioutil.ReadFile(fileToImport)
				require.NoError(t, err)
				absFilePath, err := filepath.Abs(filepath.Join(cfg.testDir, filename+cfg.testExtension))
				require.NoError(t, err)
				imp, err := Factory(absFilePath, input, logger)
				require.NoError(t, err)
				imp.WithAppName("testapp").WithPackage("package_foo")
				result, err := imp.Load(string(input))
				require.NoError(t, err)
				if *update {
					err = ioutil.WriteFile(syslFile, []byte(result), 0644)
					if err != nil {
						t.Error(err)
					}
				}
				expected, err := ioutil.ReadFile(syslFile)
				require.NoError(t, err)
				expected = syslutil.HandleCRLF(expected)
				require.NoError(t, err)
				require.Equal(t, string(expected), result)
			})
		}
	}
}

func TestLoadSwaggerJSONFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadSwaggerFromTestFiles",
		testDir:       "tests-swagger",
		testExtension: ".json",
	})
}

func TestLoadSwaggerFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadSwaggerFromTestFiles",
		testDir:       "tests-swagger",
		testExtension: ".yaml",
	})
}

func TestLoadOpenAPIFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadOpenAPIFromTestFiles",
		testDir:       "tests-openapi",
		testExtension: ".yaml",
	})
}

func TestLoadXSDFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadXSDFromTestFiles",
		testDir:       "tests-xsd",
		testExtension: ".xsd",
	})
}

/*
func TestLoadGrammarFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadGrammarFromTestFiles",
		testDir:       "tests-grammar",
		testExtension: "g",
		mode:          ModeGrammar,
		fn:            LoadGrammar,
	})
}
*/
