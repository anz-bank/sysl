package log

import (
	"context"

	"github.com/anz-bank/sysl/pkg/log/loggers"
)

type loggerContextKey int

const loggerKey loggerContextKey = iota

// WithLogger adds a copy of the logger to the context
func WithLogger(ctx context.Context, logger loggers.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger.Copy())
}

// The fields setup in WithField and WithFields are for context specific fields
// Fields will be logged alphabetically

// WithField adds a single field in the scope of the context
func WithField(ctx context.Context, key string, val interface{}) context.Context {
	newLogger := getLogger(ctx).Copy()
	newLogger.PutField(key, val)
	return context.WithValue(ctx, loggerKey, newLogger)
}

// WithFields adds multiple fields in the scope of the context
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	newLogger := getLogger(ctx).Copy()
	newLogger.PutFields(fields)
	return context.WithValue(ctx, loggerKey, newLogger)
}

// Field adds a single field for log specific fields. Does not add to the context.
func Field(ctx context.Context, key string, value interface{}) loggers.Logger {
	return getLogger(ctx).Copy().PutField(key, value)
}

// Fields adds multiple fields for log specific fields. Does not add to the context.
func Fields(ctx context.Context, fields map[string]interface{}) loggers.Logger {
	newLogger := getLogger(ctx).Copy()
	newLogger.PutFields(fields)
	return newLogger
}

// Debug logs the message at the Debug level
func Debug(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Debug(args...)
}

// Debugf logs the message at the Debug level
func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Debugf(format, args...)
}

// Error logs the message at the Error level
func Error(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Error(args...)
}

// Errorf logs the message at the Error level
func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Errorf(format, args...)
}

// Exit exits the program with specified code
func Exit(ctx context.Context, code int) {
	getLogger(ctx).Exit(code)
}

// Fatal logs the message at the Fatal level and exits the program with code 1
func Fatal(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Fatal(args...)
}

// Fatalf logs the message at the Fatal level and exits the program with code 1
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Fatalf(format, args...)
}

// Panic logs the message at the Panic level
func Panic(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Panic(args...)
}

// Panicf logs the message at the Panic level
func Panicf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Panicf(format, args...)
}

// Print logs the message at the Print level
func Print(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Print(args...)
}

// Printf logs the message at the Print level
func Printf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Printf(format, args...)
}

// Trace logs the message at the Trace level
func Trace(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Trace(args...)
}

// Tracef logs the message at the Trace level
func Tracef(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Tracef(format, args...)
}

// Warn logs the message at the Warn level
func Warn(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Warn(args...)
}

// Warnf logs the message at the Warn level
func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Warnf(format, args...)
}
