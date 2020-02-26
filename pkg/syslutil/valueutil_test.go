package syslutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetVal(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "new", GetNonEmpty("", "new"))
	assert.Equal(t, "init", GetNonEmpty("init", "new"))
	assert.Equal(t, "init", GetNonEmpty("init", ""))
}
