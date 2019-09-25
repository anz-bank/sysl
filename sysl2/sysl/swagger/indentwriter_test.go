package swagger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIndentWriter(t *testing.T) {

	buf := bytes.Buffer{}
	i := NewIndentWriter(" ", &buf)

	i.Push()
	assert.NoError(t, i.Write())
	i.Pop()

	assert.Equal(t, " ", buf.String())
}
