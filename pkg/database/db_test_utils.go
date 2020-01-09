package database

import (
	"io/ioutil"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/stretchr/testify/assert"
)

func CompareSQL(t *testing.T, expected map[string]string, actual []ScriptOutput) {
	for _, entry := range actual {
		name := entry.filename
		goldenFile := expected[name]
		golden, err := ioutil.ReadFile(goldenFile)
		assert.Nil(t, err)
		if string(golden) != entry.content {
			err := ioutil.WriteFile(name, []byte(entry.content), 0777)
			assert.Nil(t, err)
		}
		golden = syslutil.HandleCRLF(golden)
		assert.Equal(t, string(golden), entry.content)
	}
	assert.Equal(t, len(expected), len(actual))
}

func CompareContent(t *testing.T, goldenFile, output string) {
	golden, err := ioutil.ReadFile(goldenFile)
	assert.Nil(t, err)
	assert.Equal(t, string(golden), output)
}
