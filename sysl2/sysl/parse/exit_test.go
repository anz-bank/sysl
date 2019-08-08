package parse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	format := "Exiting: %s"
	param := "Oopsies!"
	message := fmt.Sprintf(format, param)
	code := 42
	e := Exitf(code, format, param)
	assert.Error(t, e)
	assert.Equal(t, message, e.Error())
	assert.Equal(t, 42, e.Code)
}
