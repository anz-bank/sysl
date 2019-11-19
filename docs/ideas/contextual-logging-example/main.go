package main

import (
	"context"

	logger "github.com/anz-bank/sysl/docs/ideas/contextual-logging-example/syslLogger"
)

func main() {
	ctx := context.Background()
	ctx = setupLogger(ctx)

	ctx = logger.AddFields(ctx, map[string]interface{}{
		"id": "12345",
	})

	/**
	 * Expected to log
	 * 2006-01-02T15:04:05Z07:00 id=12345 DEBUG This is a debug with context specific field (main.go:18)
	 */
	logger.Debug(ctx, "This is a debug with context specific field")

	/**
	 * Expected to log
	 * 2006-01-02T15:04:05Z07:00 id=12345 log=hello DEBUG This is a debug with log specific field (main.go:26)
	 */
	logger.LogFields(ctx, map[string]interface{}{
		"log": "hello",
	}).Debug("This is a debug with log specific field")

	/**
	 * Expected to log
	 * 2006-01-02T15:04:05Z07:00 id=12345 DEBUG This is a debug with context specific field again (main.go:34)
	 */
	logger.Debug(ctx, "This is a debug with context specific field again")
}

func setupLogger(ctx context.Context) context.Context {
	log := logger.NewLogger()
	return logger.AddLogger(ctx, log)
}
