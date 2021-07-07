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
	t.Helper()

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
	t.Helper()

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

func TestLoadOpenAPI2JSONFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadOpenAPI2JSONFromTestFiles",
		testDir:       "tests/openapi2",
		testExtension: ".json",
	})
}

func TestLoadOpenAPI2FromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadOpenAPI2FromTestFiles",
		testDir:       "tests/openapi2",
		testExtension: ".yaml",
	})
}

func TestLoadOpenAPI3FromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadOpenAPI3FromTestFiles",
		testDir:       "tests/openapi3",
		testExtension: ".yaml",
	})
}

func TestLoadXSDFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadXSDFromTestFiles",
		testDir:       "tests/xsd",
		testExtension: ".xsd",
	})
}

func TestLoadSpannerFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadSpannerFromTestFiles",
		testDir:       "sql/tests/spanner",
		testExtension: ".sql",
		format:        "spannerSQL",
	})
}

func TestLoadSpannerDirFromTestDir(t *testing.T) {
	runImportDirEqualityTests(t, testConfig{
		name:          "TestLoadSpannerDirFromTestDir",
		testDir:       "sql/tests/spanner/migrations",
		testExtension: "",
		format:        "spannerSQLdir",
	})
}

func TestLoadPostgresqlFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadPostgresqlFromTestFiles",
		testDir:       "sql/tests/postgresql",
		testExtension: ".sql",
		format:        "postgres",
	})
}

func TestLoadPostgresqlDirFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadPostgresqlDirFromTestFiles",
		testDir:       "sql/tests/postgresql/migrations",
		testExtension: "",
		format:        "postgresDir",
	})
}

func TestLoadMySQLFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadMySQLFromTestFiles",
		testDir:       "sql/tests/mysql",
		testExtension: ".sql",
		format:        "mysql",
	})
}

func TestLoadMySQLDirFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadMySQLDirFromTestFiles",
		testDir:       "sql/tests/mysql/migrations",
		testExtension: "",
		format:        "mysqlDir",
	})
}

func TestLoadBigQueryFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadBigQueryFromTestFiles",
		testDir:       "sql/tests/bigquery",
		testExtension: ".sql",
		format:        "bigquery",
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
