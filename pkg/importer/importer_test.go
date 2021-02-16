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
	format        string
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
				imp, err := Factory(absFilePath, false, cfg.format, input, logger)
				require.NoError(t, err)
				imp.WithAppName("TestApp").WithPackage("com.example.package")
				result, err := imp.LoadFile(absFilePath)
				require.NoError(t, err)
				if *update {
					err = ioutil.WriteFile(syslFile, []byte(result), 0600)
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

func runImportDirEqualityTests(t *testing.T, cfg testConfig) {
	logger, _ := test.NewNullLogger()
	syslFile := filepath.Join(cfg.testDir, filepath.Base(cfg.testDir)+".sysl")
	path := syslutil.MustAbsolute(t, cfg.testDir)
	imp, err := Factory(path, true, cfg.format, nil, logger)
	require.NoError(t, err)
	out, err := imp.WithAppName("TestApp").WithPackage("com.example.package").LoadFile(path)
	require.NoError(t, err)
	expected, err := ioutil.ReadFile(syslFile)
	require.NoError(t, err)
	expected = syslutil.HandleCRLF(expected)
	require.NoError(t, err)
	require.Equal(t, string(expected), out)
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

func TestLoadSpannerFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadSpannerFromTestFiles",
		testDir:       "spanner/tests",
		testExtension: ".sql",
		format:        "spannerSQL",
	})
}

func TestLoadSpannerDirFromTestDir(t *testing.T) {
	runImportDirEqualityTests(t, testConfig{
		name:          "TestLoadSpannerDirFromTestDir",
		testDir:       "spanner/tests/migrations",
		testExtension: "",
		format:        "spannerSQLdir",
	})
}

func TestLoadPostgresqlFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadPostgresqlFromTestDir",
		testDir:       "postgresql/tests",
		testExtension: ".sql",
		format:        "postgres",
	})
}

func TestLoadPostgresqlDirFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadPostgresqlDirFromTestDir",
		testDir:       "postgresql/tests/migrations",
		testExtension: "",
		format:        "postgresDir",
	})
}

func TestLoadMySQLFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadMySQLFromTestDir",
		testDir:       "mysql/tests",
		testExtension: ".sql",
		format:        "mysql",
	})
}

func TestLoadMySQLDirFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadMySQLDirFromTestDir",
		testDir:       "mysql/tests/migrations",
		testExtension: "",
		format:        "mysqlDir",
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

func TestLoadAvroFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadAvroFromTestFiles",
		testDir:       "avro/tests",
		testExtension: ".avsc",
	})
}
