package logger

import (
	"context"
	"errors"
)

func AddLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func Debug(ctx context.Context, args ...interface{}) {
	log := ctx.Value(key).(Logger)
	if log == nil {
		panic(errors.New("Logger does not exist in context"))
	}
	log.Debug(args...)
}

func AddFields(ctx context.Context, fields map[string]interface{}) context.Context {
	log := ctx.Value(key).(Logger)
	if log == nil {
		panic(errors.New("Logger does not exist in context"))
	}
	// add fields to log here, maybe a new API for Logger
	return context.WithValue(ctx, key, log)
}

func LogFields(ctx context.Context, fields map[string]interface{}) Logger {
	log := ctx.Value(key).(Logger)
	if log == nil {
		panic(errors.New("Logger does not exist in context"))
	}
	// add fields to log here, maybe a new API for Logger
	return log
}
