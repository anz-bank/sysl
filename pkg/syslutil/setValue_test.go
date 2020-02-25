package syslutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetVal(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "new", ResetVal("", "new"))
	assert.Equal(t, "init", ResetVal("init", "new"))
	assert.Equal(t, "init", ResetVal("init", ""))
}
