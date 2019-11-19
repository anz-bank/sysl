package logger

import (
	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

var key = loggerKey{}

// Partial logger interface
type Logger interface {
	Debug(args... interface{})
}

type StandardLogger struct {
	Logger
	internal *logrus.Logger
}

func NewLogger() Logger {
	return StandardLogger{internal: logrus.New()}
}

func (s StandardLogger) Debug(args... interface{}) {
	s.internal.Debug(args...)
}
