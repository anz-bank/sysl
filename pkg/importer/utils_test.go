package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSyslSafeName(t *testing.T) {
	input := "Something:Here"
	expected := "Something%3AHere"
	safeName := getSyslSafeName(input)
	assert.Equal(t, expected, safeName)
}
