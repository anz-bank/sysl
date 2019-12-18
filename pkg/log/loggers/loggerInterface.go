package loggers

// Logger is the underlying logger that is to be added to a context
type Logger interface {
	// Debug logs the message at the Debug level
	Debug(args ...interface{})
	// Debugf logs the message at the Debug level
	Debugf(format string, args ...interface{})
	// Error logs the message at the Error level
	Error(args ...interface{})
	// Errorf logs the message at the Error level
	Errorf(format string, args ...interface{})
	// Exit exits the program with specified code
	Exit(code int)
	// Fatal logs the message at the Fatal level and exits the program with code 1
	Fatal(args ...interface{})
	// Fatalf logs the message at the Fatal level and exits the program with code 1
	Fatalf(format string, args ...interface{})
	// Panic logs the message at the Panic level
	Panic(args ...interface{})
	// Panicf logs the message at the Panic level
	Panicf(format string, args ...interface{})
	// Print logs the message at the Print level
	Print(args ...interface{})
	// Printf logs the message at the Print level
	Printf(format string, args ...interface{})
	// Trace logs the message at the Trace level
	Trace(args ...interface{})
	// Tracef logs the message at the Trace level
	Tracef(format string, args ...interface{})
	// Warn logs the message at the Warn level
	Warn(args ...interface{})
	// Warnf logs the message at the Warn level
	Warnf(format string, args ...interface{})
	// PutField returns the Logger with the new field added
	PutField(key string, val interface{}) Logger
	// PutFields returns the Logger with the new fields added
	PutFields(fields map[string]interface{}) Logger
	// Copy returns a logger whose data is copied from the caller
	Copy() Logger
}
