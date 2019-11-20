package loggers

import (
	"sort"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslLogger/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	testMessage  = "This is a test message"
	simpleFormat = "%s"
)

// to test fields output for all log
var testField = testutil.GenerateMultipleFieldCases()[0].Fields

func TestCopyStandardLogger(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	logger.WithFields(map[string]interface{}{
		"numberVal": 1,
		"byteVal":   'k',
		"stringVal": "this is a sentence",
	})
	copiedLogger := logger.Copy().(*standardLogger)

	require.True(t, &logger.internal != &copiedLogger.internal)
	require.Equal(t, logger.fields, copiedLogger.fields)
	require.Equal(t, logger.sortedFields, copiedLogger.sortedFields)
	require.True(t, sort.StringsAreSorted(logger.sortedFields))
	require.True(t, &logger != &copiedLogger)
}

func TestDebug(t *testing.T) {
	testLogOutput(t, logrus.DebugLevel, nil, func() {
		NewStandardLogger().Debug(testMessage)
	})

	testLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debug(testMessage)
	})
}

func TestPrint(t *testing.T) {
	testLogOutput(t, logrus.InfoLevel, nil, func() {
		NewStandardLogger().Print(testMessage)
	})

	testLogOutput(t, logrus.InfoLevel, testField, func() {
		getStandardLoggerWithFields().Print(testMessage)
	})
}

func TestWarn(t *testing.T) {
	testLogOutput(t, logrus.WarnLevel, nil, func() {
		NewStandardLogger().Warn(testMessage)
	})

	testLogOutput(t, logrus.WarnLevel, testField, func() {
		getStandardLoggerWithFields().Warn(testMessage)
	})
}

func TestTrace(t *testing.T) {
	testLogOutput(t, logrus.TraceLevel, nil, func() {
		NewStandardLogger().Trace(testMessage)
	})

	testLogOutput(t, logrus.TraceLevel, testField, func() {
		getStandardLoggerWithFields().Trace(testMessage)
	})
}

func TestPanic(t *testing.T) {
	testLogOutput(t, logrus.PanicLevel, nil, func() {
		require.Panics(t, func() {
			NewStandardLogger().Panic(testMessage)
		})
	})

	testLogOutput(t, logrus.PanicLevel, nil, func() {
		require.Panics(t, func() {
			getStandardLoggerWithFields().Panic(testMessage)
		})
	})
}

func TestError(t *testing.T) {
	testLogOutput(t, logrus.ErrorLevel, nil, func() {
		NewStandardLogger().Error(testMessage)
	})

	testLogOutput(t, logrus.ErrorLevel, testField, func() {
		getStandardLoggerWithFields().Error(testMessage)
	})
}

func TestErrorf(t *testing.T) {
	testLogOutput(t, logrus.ErrorLevel, nil, func() {
		NewStandardLogger().Errorf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.ErrorLevel, testField, func() {
		getStandardLoggerWithFields().Errorf(simpleFormat, testMessage)
	})
}

func TestDebugf(t *testing.T) {
	testLogOutput(t, logrus.DebugLevel, nil, func() {
		NewStandardLogger().Debugf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debugf(simpleFormat, testMessage)
	})
}

func TestPrintf(t *testing.T) {
	testLogOutput(t, logrus.InfoLevel, nil, func() {
		NewStandardLogger().Printf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.InfoLevel, testField, func() {
		getStandardLoggerWithFields().Printf(simpleFormat, testMessage)
	})
}

func TestWarnf(t *testing.T) {
	testLogOutput(t, logrus.WarnLevel, nil, func() {
		NewStandardLogger().Warnf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.WarnLevel, testField, func() {
		getStandardLoggerWithFields().Warnf(simpleFormat, testMessage)
	})
}

func TestTracef(t *testing.T) {
	testLogOutput(t, logrus.TraceLevel, nil, func() {
		NewStandardLogger().Tracef(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.TraceLevel, testField, func() {
		getStandardLoggerWithFields().Tracef(simpleFormat, testMessage)
	})
}

func TestPanicf(t *testing.T) {
	testLogOutput(t, logrus.PanicLevel, nil, func() {
		require.Panics(t, func() {
			NewStandardLogger().Panicf(simpleFormat, testMessage)
		})
	})

	testLogOutput(t, logrus.PanicLevel, testField, func() {
		require.Panics(t, func() {
			getStandardLoggerWithFields().Panicf(simpleFormat, testMessage)
		})
	})
}

func testLogOutput(t *testing.T, level logrus.Level, fields map[string]interface{}, logFunc func()) {
	outputtedFields := ""
	if fields != nil {
		outputtedFields = testutil.OutputFormattedFields(fields)
	}
	parsedLevel := parseLevel(level)

	expectedOutput := strings.Join([]string{outputtedFields, parsedLevel, testMessage}, " ")
	actualOutput := testutil.RedirectOutput(t, logFunc)

	// uses Contains to avoid checking timestamps
	require.Contains(t, actualOutput, expectedOutput)
}

func TestNewStandardLogger(t *testing.T) {
	t.Parallel()

	logger := NewStandardLogger()

	require.NotNil(t, logger)
	require.IsType(t, logger, &standardLogger{})
}

func TestGetFormattedFieldEmptyFields(t *testing.T) {
	t.Parallel()

	require.Equal(t, getNewStandardLogger().getFormattedField(), "")
}

func TestGetFormattedFieldWithFields(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	logger.WithFields(map[string]interface{}{
		"numberVal": 1,
		"byteVal":   byte('k'),
		"stringVal": "this is a sentence",
	})

	expected := "byteVal=107 numberVal=1 stringVal=this is a sentence"
	require.Equal(t, expected, logger.getFormattedField())
}

func TestInsertFieldsKeyEmpty(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	logger.insertFieldsKey()
	require.Equal(t, 0, len(logger.sortedFields))
}

func TestInsertFieldsKey(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	fields := []string{"some", "random", "fields"}
	logger.insertFieldsKey(fields...)

	require.True(t, sort.StringsAreSorted(logger.sortedFields))

	sort.Strings(fields)
	require.Equal(t, fields, logger.sortedFields)
}

func TestInsertFieldsKeyAddMoreFields(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	fields1 := []string{"some", "random", "fields"}
	fields2 := []string{"even", "more", "stuff"}

	logger.insertFieldsKey(fields1...)
	logger.insertFieldsKey(fields2...)
	require.True(t, sort.StringsAreSorted(logger.sortedFields))

	combined := append(fields1, fields2...)
	sort.Strings(combined)
	require.Equal(t, combined, logger.sortedFields)
}

func TestSetInfo(t *testing.T) {
	cases := testutil.GenerateMultipleFieldCases()
	for _, c := range cases {
		if c.Fields == nil {
			continue
		}

		t.Run("TestSetInfo"+c.Name, func(mc testutil.MultipleFields) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := getNewStandardLogger()
				logger.WithFields(mc.Fields)
				entry := logger.setInfo()
				expected := testutil.OutputFormattedFields(mc.Fields)

				require.Equal(tt, expected, entry.Data[keyFields])
			}
		}(c))
	}
}

func TestWithFields(t *testing.T) {
	cases := testutil.GenerateMultipleFieldCases()
	for _, c := range cases {
		t.Run("TestWithFields"+c.Name,
			func(mc testutil.MultipleFields) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()

					if mc.Fields == nil {
						require.Panics(tt, func() {
							logger.WithFields(mc.Fields)
						})
						return
					}

					logger.WithFields(mc.Fields)
					require.True(tt, sort.StringsAreSorted(logger.sortedFields))

					expectedKeys := testutil.GetSortedKeys(mc.Fields)
					require.Equal(tt, expectedKeys, logger.sortedFields)
					require.Equal(tt, mc.Fields, logger.fields)
				}
			}(c))
	}
}

func TestWithField(t *testing.T) {
	cases := testutil.GenerateSingleFieldCases()
	for _, c := range cases {
		t.Run("TestWithField"+c.Name,
			func(sc testutil.SingleField) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()
					logger.WithField(sc.Key, sc.Val)
					value, exists := logger.fields[sc.Key]

					require.True(tt, exists)
					require.Equal(tt, sc.Val, value)
				}
			}(c))
	}
}

func TestWithFieldWithAddingMoreValues(t *testing.T) {
	cases := testutil.GenerateMultipleFieldCases()
	for _, c := range cases {
		if c.Fields == nil {
			continue
		}

		t.Run("TestWithFieldWithAddingMoreValues"+c.Name,
			func(mc testutil.MultipleFields) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()

					for k, v := range mc.Fields {
						logger.WithField(k, v)
					}

					require.True(tt, sort.StringsAreSorted(logger.sortedFields))
					expectedKeys := testutil.GetSortedKeys(mc.Fields)
					require.Equal(tt, expectedKeys, logger.sortedFields)
					require.Equal(tt, mc.Fields, logger.fields)
				}
			}(c))
	}
}

func TestWithFieldReplaceValues(t *testing.T) {
	t.Parallel()

	key := "random"
	oldVal := 1
	newVal := 2

	logger := getNewStandardLogger()

	logger.WithField(key, oldVal)
	assertFieldExists(t, logger, map[string]interface{}{key: oldVal})

	logger.WithField(key, newVal)
	assertFieldExists(t, logger, map[string]interface{}{key: newVal})
	require.Equal(t, []string{key}, logger.sortedFields)
}

func TestWithFieldsReplaceValues(t *testing.T) {
	t.Parallel()

	field := map[string]interface{}{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	logger := getNewStandardLogger()
	logger.WithFields(field)

	assertFieldExists(t, logger, field)

	for k := range field {
		field[k] = "replaced"
	}
	logger.WithFields(field)

	assertFieldExists(t, logger, field)
	require.Equal(t, testutil.GetSortedKeys(field), logger.sortedFields)
}

func assertFieldExists(t *testing.T, logger *standardLogger, fieldToCheck map[string]interface{}) {
	for key, expectedVal := range fieldToCheck {
		curVal, exists := logger.fields[key]
		require.True(t, exists)
		require.Equal(t, expectedVal, curVal)
	}
}

func getNewStandardLogger() *standardLogger {
	return NewStandardLogger().(*standardLogger)
}

func getStandardLoggerWithFields() *standardLogger {
	logger := NewStandardLogger()
	logger.WithFields(testField)
	return logger.(*standardLogger)
}
