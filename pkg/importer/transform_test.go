package importer

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/anz-bank/sysl/pkg/arrai/transform"
	"github.com/arr-ai/arrai/pkg/test"
	"github.com/arr-ai/arrai/rel"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestImporters tests each importer by running its transform on a file that matches the input.* pattern, then passing
// the resulting Sysl to a test function defined in the file tests.arrai which accepts a Sysl string and returns
// a standard test expression (leaves are true/false to indicate pass/fail). This expression is passed into the standard
// arrai testing mechanism.
func TestImporters(t *testing.T) {
	t.Parallel()

	dirEntries, err := bundles.BundlesFs.ReadDir("importers")
	require.NoError(t, err)
	require.NotEmpty(t, dirEntries)
	testFiles := make([]test.File, 0, len(dirEntries))

	// Run tests on each importer
	for _, importerDir := range dirEntries {
		// Read the importer's directory content
		require.Equal(t, os.ModeDir, importerDir.Type())
		dirPath := path.Join("importers", importerDir.Name())
		files, err := bundles.BundlesFs.ReadDir(dirPath)
		require.NoError(t, err)

		// Find and read the input file. It can be of any filetype, so look for something that starts with 'input.'
		var inputFile string
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "input.") {
				inputFileBytes, err := bundles.BundlesFs.ReadFile(path.Join(dirPath, file.Name()))
				require.NoError(t, err)
				inputFile = string(inputFileBytes)
				break
			}
		}
		require.NotEmpty(t, inputFile, "Failed to find input file for importer, "+
			"no files matching 'input.*' pattern in embedded directory '%s'", dirPath)

		imp, err := MakeTransformImporter(logrus.StandardLogger(), importerDir.Name()).
			Configure(&ImporterArg{
				AppName: "TestNamespace::TestApp",
			})
		require.NoError(t, err)
		// Run the importer and get the resulting Sysl
		sysl, err := imp.Load(inputFile)
		require.NoError(t, err)

		// Read the test script
		testsPath := path.Join(dirPath, "tests.arrai")
		testScriptBytes, err := bundles.BundlesFs.ReadFile(testsPath)
		require.NoError(t, err)

		// Run the test script with the imported Sysl as the argument
		testFile, err := transform.RunTests(testScriptBytes, importerDir.Name(), rel.NewString([]rune(sysl)))
		require.NoError(t, err)
		testFiles = append(testFiles, testFile)
	}

	// Run each arr.ai test as a native go subtest with the correct hierarchical position
	for _, testFile := range testFiles {
		testFile := testFile
		t.Run(testFile.Path, func(t *testing.T) {
			for _, result := range testFile.Results {
				result := result
				t.Run(result.Name, func(t *testing.T) {
					assert.GreaterOrEqual(t, result.Outcome, test.Ignored, result.Message)
				})
			}
		})
	}
}
