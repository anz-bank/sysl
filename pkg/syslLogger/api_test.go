package sysllogger

import (
	"context"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslLogger/loggers"
	"github.com/anz-bank/sysl/pkg/syslLogger/testutil"
	"github.com/stretchr/testify/require"
)

func TestAddLogger(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := loggers.NewMockLogger()
	logger.On("Copy").Return(loggers.NewMockLogger()).Once()
	newCtx := AddLogger(ctx, logger)
	newLogger := newCtx.Value(loggerKey).(*loggers.MockLogger)

	require.True(t, logger.AssertExpectations(t))
	require.True(t, &ctx != &newCtx)
	require.True(t, &logger != &newLogger)
}

func TestAddField(t *testing.T) {
	cases := testutil.GenerateSingleFieldCases()
	for _, c := range cases {
		t.Run(c.Name, func(sc testutil.SingleField) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := loggers.NewMockLogger()
				logger.On("WithField", sc.Key, sc.Val).Return(logger)
				logger.On("Copy").Return(loggers.NewMockLogger())
				ctx := context.WithValue(context.Background(), loggerKey, logger)

				newCtx := AddField(ctx, sc.Key, sc.Val)
				newLogger := newCtx.Value(loggerKey).(*loggers.MockLogger)

				require.True(tt, logger.AssertExpectations(tt))
				require.True(tt, &ctx != &newCtx)
				require.True(tt, &logger != &newLogger)
			}
		}(c))
	}
}

func TestLogField(t *testing.T) {
	cases := testutil.GenerateSingleFieldCases()
	for _, c := range cases {
		t.Run(c.Name, func(sc testutil.SingleField) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := loggers.NewMockLogger()

				// this is done so that the mock logger is of the
				// same reference, for testing purposes
				logger.On("Copy").Return(logger)
				logger.On("WithField", sc.Key, sc.Val).Return(loggers.NewMockLogger())
				ctx := context.WithValue(context.Background(), loggerKey, logger)

				newLogger := LogField(ctx, sc.Key, sc.Val).(*loggers.MockLogger)

				require.True(tt, logger.AssertExpectations(tt))
				require.True(tt, &logger != &newLogger)
			}
		}(c))
	}
}

func TestAddFields(t *testing.T) {
	cases := testutil.GenerateMultipleFieldCases()
	for _, c := range cases {
		t.Run(c.Name, func(mc testutil.MultipleFields) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := loggers.NewMockLogger()
				logger.On("WithFields", mc.Fields).Return(logger)
				logger.On("Copy").Return(loggers.NewMockLogger())
				ctx := context.WithValue(context.Background(), loggerKey, logger)

				newCtx := AddFields(ctx, mc.Fields)
				newLogger := newCtx.Value(loggerKey).(*loggers.MockLogger)

				require.True(tt, logger.AssertExpectations(tt))
				require.True(tt, &ctx != &newCtx)
				require.True(tt, &logger != &newLogger)
			}
		}(c))
	}
}

func TestLogFields(t *testing.T) {
	cases := testutil.GenerateMultipleFieldCases()
	for _, c := range cases {
		t.Run(c.Name, func(mc testutil.MultipleFields) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := loggers.NewMockLogger()

				// this is done so that the mock logger is of the
				// same reference, for testing purposes
				logger.On("Copy").Return(logger)
				logger.On("WithFields", mc.Fields).Return(loggers.NewMockLogger())
				ctx := context.WithValue(context.Background(), loggerKey, logger)

				newLogger := LogFields(ctx, mc.Fields).(*loggers.MockLogger)

				require.True(tt, logger.AssertExpectations(tt))
				require.True(tt, &logger != &newLogger)
			}
		}(c))
	}
}

func testAllCases(t *testing.T, name string, test func(testutil.LogCases, *testing.T) func(*testing.T)) {
	for _, c := range testutil.GenerateLogCases() {
		t.Run(name+" "+c.Name, test(c, t))
	}
}

func TestPrint(t *testing.T) {
	testAllCases(t, "TestPrint", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Print", Print)
			testFormattedLog(tt, logCase, "Printf", Printf)
		}
	})
}

func TestWarn(t *testing.T) {
	testAllCases(t, "TestWarn", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Warn", Warn)
			testFormattedLog(tt, logCase, "Warnf", Warnf)
		}
	})
}

func TestPanic(t *testing.T) {
	testAllCases(t, "TestPanic", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Panic", Panic)
			testFormattedLog(tt, logCase, "Panicf", Panicf)
		}
	})
}

func TestTrace(t *testing.T) {
	testAllCases(t, "TestTrace", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Trace", Trace)
			testFormattedLog(tt, logCase, "Tracef", Tracef)
		}
	})
}

func TestError(t *testing.T) {
	testAllCases(t, "TestError", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Error", Error)
			testFormattedLog(tt, logCase, "Errorf", Errorf)
		}
	})
}

func TestDebug(t *testing.T) {
	testAllCases(t, "TestDebug", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Debug", Debug)
			testFormattedLog(tt, logCase, "Debugf", Debugf)
		}
	})
}

func TestFatal(t *testing.T) {
	testAllCases(t, "TestFatal", func(logCase testutil.LogCases, t *testing.T) func(*testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			testNonFormatLog(tt, logCase, "Fatal", Fatal)
			testFormattedLog(tt, logCase, "Fatalf", Fatalf)
		}
	})
}

func TestExit(t *testing.T) {
	logger := loggers.NewMockLogger()
	logger.On("Exit", 1)
	ctx := context.WithValue(context.Background(), loggerKey, logger)

	Exit(ctx, 1)

	require.True(t, logger.AssertExpectations(t))
}

func testNonFormatLog(t *testing.T, logCase testutil.LogCases, name string, log func(context.Context, ...interface{})) {
	logger := loggers.NewMockLogger()
	logger.On(name, logCase.Arguments...).Once()
	ctx := context.Background()
	ctx = context.WithValue(ctx, loggerKey, logger)

	log(ctx, logCase.Arguments...)

	require.True(t, logger.AssertExpectations(t))
}

func testFormattedLog(
	t *testing.T,
	logCase testutil.LogCases,
	name string,
	log func(context.Context, string, ...interface{}),
) {
	logger := loggers.NewMockLogger()
	arguments := append([]interface{}{logCase.Format}, logCase.Arguments...)
	logger.On(name, arguments...).Once()
	ctx := context.Background()
	ctx = context.WithValue(ctx, loggerKey, logger)

	log(ctx, logCase.Format, logCase.Arguments...)

	require.True(t, logger.AssertExpectations(t))
}
