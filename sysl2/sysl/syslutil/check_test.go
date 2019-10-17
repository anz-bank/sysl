package syslutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertOK(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		Assert(true, "whatever")
	})
}

func TestAssertNotOK(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		Assert(false, "whatever")
	})
}

func TestPanicOnErrorWithNoError(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		PanicOnError(nil)
	})
}

func TestPanicOnErrorWithError(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		PanicOnError(fmt.Errorf("whatever"))
	})
}

func TestPanicOnErrorfWithNoError(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		PanicOnErrorf(nil, "whatever")
	})
}

func TestPanicfOnErrorWithError(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		PanicOnErrorf(fmt.Errorf("whatever"), "whatever")
	})
}
