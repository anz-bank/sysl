package sysllogger

import (
	"context"

	"github.com/anz-bank/sysl/pkg/syslLogger/loggers"
)

type loggerContextKey int

const loggerKey = loggerContextKey(1)

func AddLogger(ctx context.Context, logger loggers.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger.Copy())
}

/**
 * The fields setup in AddField and AddFields are for context specific fields
 * Fields will be logged alphabetically
 */
func AddField(ctx context.Context, key string, val interface{}) context.Context {
	newLogger := getLogger(ctx)
	newLogger.WithField(key, val)
	return AddLogger(ctx, newLogger)
}

func AddFields(ctx context.Context, fields map[string]interface{}) context.Context {
	newLogger := getLogger(ctx)
	newLogger.WithFields(fields)
	return AddLogger(ctx, newLogger)
}

func LogField(ctx context.Context, key string, value interface{}) loggers.Logger {
	return getLogger(ctx).Copy().WithField(key, value)
}

func LogFields(ctx context.Context, fields map[string]interface{}) loggers.Logger {
	newLogger := getLogger(ctx).Copy()
	newLogger.WithFields(fields)
	return newLogger
}

func Debug(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Debugf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Errorf(format, args...)
}

func Exit(ctx context.Context, code int) {
	getLogger(ctx).Exit(code)
}

func Fatal(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Fatalf(format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Panic(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Panicf(format, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Print(args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Printf(format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Trace(args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Tracef(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	getLogger(ctx).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Warnf(format, args...)
}
