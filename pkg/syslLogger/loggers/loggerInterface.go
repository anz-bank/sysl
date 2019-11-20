package loggers

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Exit(code int)
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WithField(key string, val interface{}) Logger
	WithFields(fields map[string]interface{})
	Copy() Logger
}
