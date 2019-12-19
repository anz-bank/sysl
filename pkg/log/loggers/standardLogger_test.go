package loggers

import (
	"sort"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/log/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMessage  = "This is a test message"
	simpleFormat = "%s"
)

// to test fields output for all log
var testField = testutil.GenerateMultipleFieldsCases()[0].Fields

func TestCopyStandardLogger(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger().PutFields(map[string]interface{}{
		"numberVal": 1,
		"byteVal":   'k',
		"stringVal": "this is a sentence",
	}).(*standardLogger)
	copiedLogger := logger.Copy().(*standardLogger)
	assert.True(t, logger.internal != copiedLogger.internal)
	assert.Equal(t, logger.fields, copiedLogger.fields)
	assert.Equal(t, logger.sortedFields, copiedLogger.sortedFields)
	assert.True(t, sort.StringsAreSorted(logger.sortedFields))
	assert.True(t, logger != copiedLogger)
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

	expectedOutput := strings.Join([]string{outputtedFields, strings.ToUpper(level.String()), testMessage}, " ")
	actualOutput := testutil.RedirectOutput(t, logFunc)

	// uses Contains to avoid checking timestamps
	assert.Contains(t, actualOutput, expectedOutput)
}

func TestNewStandardLogger(t *testing.T) {
	t.Parallel()

	logger := NewStandardLogger()

	require.NotNil(t, logger)
	assert.IsType(t, logger, &standardLogger{})
}

func TestGetFormattedFieldEmptyFields(t *testing.T) {
	t.Parallel()

	require.Equal(t, getNewStandardLogger().getFormattedField(), "")
}

func TestGetFormattedFieldWithFields(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger().PutFields(map[string]interface{}{
		"numberVal": 1,
		"byteVal":   byte('k'),
		"stringVal": "this is a sentence",
	}).(*standardLogger)

	expected := "byteVal=107 numberVal=1 stringVal=this is a sentence"
	assert.Equal(t, expected, logger.getFormattedField())
}

func TestInsertFieldsKeyEmpty(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	logger.insertFieldsKey()
	assert.Equal(t, 0, len(logger.sortedFields))
}

func TestInsertFieldsKey(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	fields := []string{"some", "random", "fields"}
	logger.insertFieldsKey(fields...)

	sort.Strings(fields)
	assert.Equal(t, fields, logger.sortedFields)
}

func TestInsertFieldsKeyAddMoreFields(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	fields1 := []string{"some", "random", "fields"}
	fields2 := []string{"even", "more", "stuff"}

	logger.insertFieldsKey(fields1...)
	logger.insertFieldsKey(fields2...)

	combined := append(fields1, fields2...)
	sort.Strings(combined)
	assert.Equal(t, combined, logger.sortedFields)
}

func TestSetInfo(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		if c.Fields == nil {
			continue
		}

		t.Run("TestSetInfo"+c.Name, func(mc testutil.MultipleFields) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := getNewStandardLogger().PutFields(mc.Fields).(*standardLogger)
				entry := logger.setInfo()
				expected := testutil.OutputFormattedFields(mc.Fields)

				assert.Equal(tt, expected, entry.Data[keyFields])
			}
		}(c))
	}
}

func TestWithFields(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		t.Run("TestWithFields"+c.Name,
			func(mc testutil.MultipleFields) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()

					if mc.Fields == nil {
						require.Panics(tt, func() {
							logger.PutFields(mc.Fields)
						})
						return
					}

					logger.PutFields(mc.Fields)
					expectedKeys := testutil.GetSortedKeys(mc.Fields)
					assert.Equal(tt, expectedKeys, logger.sortedFields)
					assert.Equal(tt, mc.Fields, logger.fields)
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
					logger.PutField(sc.Key, sc.Val)
					value, exists := logger.fields[sc.Key]

					require.True(tt, exists)
					assert.Equal(tt, sc.Val, value)
				}
			}(c))
	}
}

func TestWithFieldWithAddingMoreValues(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
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
						logger.PutField(k, v)
					}

					expectedKeys := testutil.GetSortedKeys(mc.Fields)
					assert.Equal(tt, expectedKeys, logger.sortedFields)
					assert.Equal(tt, mc.Fields, logger.fields)
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

	logger.PutField(key, oldVal)
	assertFieldExists(t, logger, map[string]interface{}{key: oldVal})

	logger.PutField(key, newVal)
	assertFieldExists(t, logger, map[string]interface{}{key: newVal})
	assert.Equal(t, []string{key}, logger.sortedFields)
}

func TestWithFieldsReplaceValues(t *testing.T) {
	t.Parallel()

	field := map[string]interface{}{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	logger := getNewStandardLogger().PutFields(field).(*standardLogger)

	assertFieldExists(t, logger, field)

	for k := range field {
		field[k] = "replaced"
	}
	logger.PutFields(field)

	assertFieldExists(t, logger, field)
	assert.Equal(t, testutil.GetSortedKeys(field), logger.sortedFields)
}

func assertFieldExists(t *testing.T, logger *standardLogger, fieldToCheck map[string]interface{}) {
	for key, expectedVal := range fieldToCheck {
		curVal, exists := logger.fields[key]
		require.True(t, exists)
		assert.Equal(t, expectedVal, curVal)
	}
}

func getNewStandardLogger() *standardLogger {
	return NewStandardLogger().(*standardLogger)
}

func getStandardLoggerWithFields() *standardLogger {
	logger := NewStandardLogger().PutFields(testField)
	return logger.(*standardLogger)
}
