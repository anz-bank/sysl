package loggers

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const keyFields = "_fields"

type standardLogger struct {
	fields       map[string]interface{}
	sortedFields []string
	internal     *logrus.Logger
}

type standardFormat struct{}

func (sf *standardFormat) Format(entry *logrus.Entry) ([]byte, error) {
	message := strings.Builder{}
	message.WriteString(entry.Time.Format(time.RFC3339Nano))
	message.WriteByte(' ')

	if len(entry.Data) > 0 {
		message.WriteString(entry.Data[keyFields].(string))
		message.WriteByte(' ')
	}

	message.WriteString(parseLevel(entry.Level))

	if entry.Message != "" {
		message.WriteByte(' ')
		message.WriteString(entry.Message)
	}

	// TODO: add codelinker's message here
	message.WriteByte('\n')
	return []byte(message.String()), nil
}

func parseLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.TraceLevel:
		return "TRACE"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.InfoLevel:
		// Print logs at the InfoLevel
		return "PRINT"
	case logrus.PanicLevel:
		return "PANIC"
	case logrus.FatalLevel:
		return "FATAL"
	default:
		panic("Logging: Unsupported log level")
	}
}

// NewStandardLogger returns a logger with logrus standard logger as the internal logger
func NewStandardLogger() Logger {
	logger := setupStandardLogger()

	return &standardLogger{
		fields:       make(map[string]interface{}),
		sortedFields: []string{},
		internal:     logger,
	}
}

func (sl *standardLogger) Debug(args ...interface{}) {
	sl.setInfo().Debug(args...)
}

func (sl *standardLogger) Debugf(format string, args ...interface{}) {
	sl.setInfo().Debugf(format, args...)
}

func (sl *standardLogger) Error(args ...interface{}) {
	sl.setInfo().Error(args...)
}

func (sl *standardLogger) Errorf(format string, args ...interface{}) {
	sl.setInfo().Errorf(format, args...)
}

func (sl *standardLogger) Exit(code int) {
	sl.internal.Exit(code)
}

func (sl *standardLogger) Fatal(args ...interface{}) {
	sl.setInfo().Fatal(args...)
}

func (sl *standardLogger) Fatalf(format string, args ...interface{}) {
	sl.setInfo().Fatalf(format, args...)
}

func (sl *standardLogger) Panic(args ...interface{}) {
	sl.setInfo().Panic(args...)
}

func (sl *standardLogger) Panicf(format string, args ...interface{}) {
	sl.setInfo().Panicf(format, args...)
}

func (sl *standardLogger) Print(args ...interface{}) {
	sl.setInfo().Print(args...)
}

func (sl *standardLogger) Printf(format string, args ...interface{}) {
	sl.setInfo().Printf(format, args...)
}

func (sl *standardLogger) Trace(args ...interface{}) {
	sl.setInfo().Trace(args...)
}

func (sl *standardLogger) Tracef(format string, args ...interface{}) {
	sl.setInfo().Tracef(format, args...)
}

func (sl *standardLogger) Warn(args ...interface{}) {
	sl.setInfo().Warn(args...)
}

func (sl *standardLogger) Warnf(format string, args ...interface{}) {
	sl.setInfo().Warnf(format, args...)
}

func (sl *standardLogger) WithField(key string, val interface{}) Logger {
	if _, exists := sl.fields[key]; !exists {
		sl.insertFieldsKey(key)
	}
	sl.fields[key] = val
	return sl
}

func (sl *standardLogger) WithFields(fields map[string]interface{}) {
	if fields == nil {
		panic("fields can not be empty")
	}

	keys := make([]string, len(fields))
	index := 0
	for key, val := range fields {
		if _, exists := sl.fields[key]; !exists {
			keys[index] = key
			index++
		}
		sl.fields[key] = val
	}
	sl.insertFieldsKey(keys[:index]...)
}

func (sl *standardLogger) insertFieldsKey(fields ...string) {
	newFields := append(sl.sortedFields, fields...)
	sort.Strings(newFields)
	sl.sortedFields = newFields
}

func (sl *standardLogger) setInfo() *logrus.Entry {
	// TODO: set linker here
	return sl.internal.WithFields(logrus.Fields{
		keyFields: sl.getFormattedField(),
	})
}

func (sl *standardLogger) getFormattedField() string {
	if len(sl.fields) == 0 {
		return ""
	}

	fields := strings.Builder{}
	fields.WriteString(fmt.Sprintf("%s=%v", sl.sortedFields[0], sl.fields[sl.sortedFields[0]]))
	if len(sl.fields) > 1 {
		for _, field := range sl.sortedFields[1:] {
			fields.WriteString(fmt.Sprintf(" %s=%v", field, sl.fields[field]))
		}
	}
	return fields.String()
}

func (sl *standardLogger) Copy() Logger {
	fieldsCopy := make(map[string]interface{})
	for key, val := range sl.fields {
		fieldsCopy[key] = val
	}
	sortedFields := make([]string, len(sl.fields))
	copy(sortedFields, sl.sortedFields)
	logger := setupStandardLogger()

	return &standardLogger{
		fields:       fieldsCopy,
		internal:     logger,
		sortedFields: sortedFields,
	}
}

func setupStandardLogger() *logrus.Logger {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&standardFormat{})

	// makes sure that it always logs every level
	logger.SetLevel(logrus.TraceLevel)

	// explicitly set it to os.Stderr
	logger.SetOutput(os.Stderr)

	return logger
}
