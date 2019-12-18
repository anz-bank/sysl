package log

import (
	"context"
	"testing"

	"github.com/anz-bank/sysl/pkg/log/loggers"
	"github.com/anz-bank/sysl/pkg/log/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLogger(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := loggers.NewMockLogger()
	logger.On("Copy").Return(loggers.NewMockLogger()).Once()
	newCtx := WithLogger(ctx, logger)
	newLogger := newCtx.Value(loggerKey).(*loggers.MockLogger)

	require.True(t, logger.AssertExpectations(t))
	assert.True(t, &ctx != &newCtx)
	assert.True(t, &logger != &newLogger)
}

func TestWithField(t *testing.T) {
	testSingleFieldCases(t, "TestWithField",
		func(tt *testing.T, ctx context.Context, sc testutil.SingleField) {
			logger := ctx.Value(loggerKey).(*loggers.MockLogger)
			newCtx := WithField(ctx, sc.Key, sc.Val)
			newLogger := newCtx.Value(loggerKey).(*loggers.MockLogger)

			require.True(tt, logger.AssertExpectations(tt))
			assert.True(tt, &ctx != &newCtx)
			assert.True(tt, &logger != &newLogger)
		})
}

func TestField(t *testing.T) {
	testSingleFieldCases(t, "TestField",
		func(tt *testing.T, ctx context.Context, sc testutil.SingleField) {
			logger := ctx.Value(loggerKey).(*loggers.MockLogger)
			newLogger := Field(ctx, sc.Key, sc.Val).(*loggers.MockLogger)

			require.True(tt, logger.AssertExpectations(tt))
			assert.True(tt, &logger != &newLogger)
		})
}

func TestWithFields(t *testing.T) {
	testMultipleFieldsCases(t, "TestWithFields",
		func(tt *testing.T, ctx context.Context, mc testutil.MultipleFields) {
			logger := ctx.Value(loggerKey).(*loggers.MockLogger)
			newCtx := WithFields(ctx, mc.Fields)
			newLogger := newCtx.Value(loggerKey).(*loggers.MockLogger)

			require.True(tt, logger.AssertExpectations(tt))
			assert.True(tt, &ctx != &newCtx)
			assert.True(tt, &logger != &newLogger)
		})
}

func TestFields(t *testing.T) {
	testMultipleFieldsCases(t, "TestFields",
		func(tt *testing.T, ctx context.Context, mc testutil.MultipleFields) {
			logger := ctx.Value(loggerKey).(*loggers.MockLogger)
			newLogger := Fields(ctx, mc.Fields).(*loggers.MockLogger)

			require.True(tt, logger.AssertExpectations(tt))
			assert.True(tt, &logger != &newLogger)
		})
}

func testAllCases(t *testing.T, name string, test func(testutil.LogCases, *testing.T)) {
	for _, c := range testutil.GenerateLogCases() {
		c := c
		t.Run(name+" "+c.Name, func(tt *testing.T) { test(c, tt) })
	}
}

func TestPrint(t *testing.T) {
	testAllCases(t, "TestPrint", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Print", Print)
		testFormattedLog(tt, logCase, "Printf", Printf)
	})
}

func TestWarn(t *testing.T) {
	testAllCases(t, "TestWarn", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Warn", Warn)
		testFormattedLog(tt, logCase, "Warnf", Warnf)
	})
}

func TestPanic(t *testing.T) {
	testAllCases(t, "TestPanic", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Panic", Panic)
		testFormattedLog(tt, logCase, "Panicf", Panicf)
	})
}

func TestTrace(t *testing.T) {
	testAllCases(t, "TestTrace", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Trace", Trace)
		testFormattedLog(tt, logCase, "Tracef", Tracef)
	})
}

func TestError(t *testing.T) {
	testAllCases(t, "TestError", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Error", Error)
		testFormattedLog(tt, logCase, "Errorf", Errorf)
	})
}

func TestDebug(t *testing.T) {
	testAllCases(t, "TestDebug", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Debug", Debug)
		testFormattedLog(tt, logCase, "Debugf", Debugf)
	})
}

func TestFatal(t *testing.T) {
	testAllCases(t, "TestFatal", func(logCase testutil.LogCases, tt *testing.T) {
		tt.Parallel()

		testNonFormatLog(tt, logCase, "Fatal", Fatal)
		testFormattedLog(tt, logCase, "Fatalf", Fatalf)
	})
}

func TestExit(t *testing.T) {
	logger := loggers.NewMockLogger()
	logger.On("Exit", 1)
	ctx := context.WithValue(context.Background(), loggerKey, logger)

	Exit(ctx, 1)

	require.True(t, logger.AssertExpectations(t))
}

func testNonFormatLog(
	t *testing.T,
	logCase testutil.LogCases,
	name string,
	logFunc func(context.Context, ...interface{}),
) {
	logger := loggers.NewMockLogger()
	logger.On(name, logCase.Arguments...).Once()
	ctx := context.Background()
	ctx = context.WithValue(ctx, loggerKey, logger)

	logFunc(ctx, logCase.Arguments...)

	require.True(t, logger.AssertExpectations(t))
}

func testFormattedLog(
	t *testing.T,
	logCase testutil.LogCases,
	name string,
	logFunc func(context.Context, string, ...interface{}),
) {
	logger := loggers.NewMockLogger()
	arguments := append([]interface{}{logCase.Format}, logCase.Arguments...)
	logger.On(name, arguments...).Once()
	ctx := context.Background()
	ctx = context.WithValue(ctx, loggerKey, logger)

	logFunc(ctx, logCase.Format, logCase.Arguments...)

	require.True(t, logger.AssertExpectations(t))
}

func testSingleFieldCases(
	t *testing.T,
	name string,
	testFunc func(tt *testing.T, ctx context.Context, sc testutil.SingleField),
) {
	cases := testutil.GenerateSingleFieldCases()
	for _, c := range cases {
		c := c
		t.Run(name+" "+c.Name, func(tt *testing.T) {
			tt.Parallel()

			logger := loggers.NewMockLogger()
			// this is done so that the mock logger is of the
			// same reference, for testing purposes
			logger.On("Copy").Return(logger)
			logger.On("PutField", c.Key, c.Val).Return(loggers.NewMockLogger())
			testFunc(tt, context.WithValue(context.Background(), loggerKey, logger), c)
		})
	}
}

func testMultipleFieldsCases(
	t *testing.T,
	name string,
	testFunc func(*testing.T, context.Context, testutil.MultipleFields),
) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		c := c
		t.Run(name+" "+c.Name, func(tt *testing.T) {
			tt.Parallel()

			logger := loggers.NewMockLogger()
			// this is done so that the mock logger is of the
			// same reference, for testing purposes
			logger.On("Copy").Return(logger)
			logger.On("PutFields", c.Fields).Return(loggers.NewMockLogger())
			testFunc(tt, context.WithValue(context.Background(), loggerKey, logger), c)
		})
	}
}
