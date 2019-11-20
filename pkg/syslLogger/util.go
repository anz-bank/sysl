package sysllogger

import (
	"context"

	"github.com/anz-bank/sysl/pkg/syslLogger/loggers"
)

func getLogger(ctx context.Context) loggers.Logger {
	log, exists := ctx.Value(loggerKey).(loggers.Logger)
	if !exists {
		panic("Logger does not exist in context")
	}
	return log
}
