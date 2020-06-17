package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/require"
)

type testConfig struct {
	name          string
	testDir       string
	testExtension string
	mode          string
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
		if strings.EqualFold(ext, cfg.testExtension) {
			filename := strings.TrimSuffix(f.Name(), ext)
			t.Run(fmt.Sprintf("%s - %s", cfg.name, filename), func(t *testing.T) {
				t.Parallel()
				swaggerFile := syslutil.MustAbsolute(t, filepath.Join(cfg.testDir, filename+cfg.testExtension))
				input, err := ioutil.ReadFile(swaggerFile)
				require.NoError(t, err)
				syslFile := filepath.Join(cfg.testDir, filename+".sysl")
				expected, err := ioutil.ReadFile(syslFile)
				require.NoError(t, err)
				expected = syslutil.HandleCRLF(expected)

				outputData := OutputData{
					AppName:     "testapp",
					Package:     "package_foo",
					SwaggerRoot: filepath.Dir(swaggerFile),
					Mode:        cfg.mode,
				}
				result, err := cfg.fn(outputData, string(input), logger)
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
		testExtension: ".yaml",
		mode:          ModeSwagger,
		fn:            LoadSwaggerText,
	})
}

func TestLoadOpenApiFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadOpenAPIFromTestFiles",
		testDir:       "tests-openapi",
		testExtension: ".yaml",
		mode:          ModeOpenAPI,
		fn:            LoadOpenAPIText,
	})
}

func TestLoadXSDFromTestFiles(t *testing.T) {
	runImportEqualityTests(t, testConfig{
		name:          "TestLoadXSDFromTestFiles",
		testDir:       "tests-xsd",
		testExtension: ".xsd",
		mode:          ModeXSD,
		fn:            LoadXSDText,
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
