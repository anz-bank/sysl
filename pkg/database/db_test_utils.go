package database

import (
	"os"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/stretchr/testify/assert"
)

func CompareSQL(t *testing.T, expected map[string]string, actual []ScriptOutput) {
	for _, entry := range actual {
		name := entry.filename
		goldenFile := expected[name]
		golden, err := os.ReadFile(goldenFile)
		assert.Nil(t, err)
		golden = syslutil.HandleCRLF(golden)
		if string(golden) != entry.content {
			err := os.WriteFile(name, []byte(entry.content), 0600)
			assert.Nil(t, err)
		}
		golden = syslutil.HandleCRLF(golden)
		assert.Equal(t, string(golden), entry.content)
	}
	assert.Equal(t, len(expected), len(actual))
}

func CompareContent(t *testing.T, goldenFile, output string) {
	golden, err := os.ReadFile(goldenFile)
	assert.Nil(t, err)
	golden = syslutil.HandleCRLF(golden)
	assert.Equal(t, string(golden), output)
}
