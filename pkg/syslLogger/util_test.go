package sysllogger

import (
	"context"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslLogger/loggers"
	"github.com/stretchr/testify/require"
)

func TestGetLogger(t *testing.T) {
	t.Parallel()

	t.Run("Context has no logger", testEmptyContext)
	t.Run("Context has a logger", testContextWithLogger)
}

func testEmptyContext(t *testing.T) {
	require.Panics(t, func() {
		getLogger(context.Background())
	})
}

func testContextWithLogger(t *testing.T) {
	ctx := context.Background()
	ctx = AddLogger(ctx, loggers.NewStandardLogger())

	logger := getLogger(ctx)
	require.NotNil(t, logger)

	fromContext := ctx.Value(loggerKey).(loggers.Logger)
	require.True(t, &logger != &fromContext)
}
