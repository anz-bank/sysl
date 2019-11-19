# Contextual Logging Design Proposal

## Problem

The current logger library, Logrus, makes the use of context quite difficult. In Logrus, the use of context makes use of an intermediate data structure that processes the context data before logging. To achieve this, it is necessary to set up how context is processed through the use of Logrus interface method `Format`. This library aims to make using context with logging easier by standardising log formats and enforces the use of context as part of the effort to use context more in sysl (see issue [#397](https://github.com/anz-bank/sysl/issues/397))

## Idea

This logging library will have the following features

### Logger API

The idea is to have [logrus](https://github.com/sirupsen/logrus)-like API while using logrus as the underlying logger initially. Instead of acting as a wrapper of logrus, a standardised interface will be created which an underlying logger can implement. This makes it easy in the future should we decide to remove logrus. The API will also enforce the use of context.

### Standardised Formatting

Formatting is to be standardised. General format would look like

```txt
[Time in RFC3339 format] Fields Level Message (Caller)
```

For example

```txt
2006-01-02T15:04:05Z07:00 caller="file" root="." WARN This log was called (file.go:5)
```

Fields can be customised according to needs through the use of API.

## Design

The main API will be similar to Logrus and consist of these functions

```go
package syslLogger

func AddLogger(ctx context.Context, logger Logger) context.Context {...}
func NewLogger() Logger {...} // One for each logger type

/**
 * The fields setup in AddField and AddFields are for context specific fields
 * Fields will be logged alphabetically
 */
func AddField(key string, val interface{}) context.Context {...}
func AddFields(fields map[string]interface{}) context.Context {...}

/**
 * LogFields and LogField are meant to be used for a specific log. It does not return context with
 * a new logger, fields here will be logged in the AdditionalField format
 */
func LogFields(ctx context.Context, fields map[string]interface{}) Logger {...}
func LogField(ctx context.Context, key string, value interface{}) Logger {...}

// Logger functionalities, Logger API will be similar to this
func Debug(ctx context.Context, args... interface{}) {...}
func Debugf(ctx context.Context, format string, args... interface{}) {...}
func Error(ctx context.Context, args... interface{}) {...}
func Errorf(ctx context.Context, format string, args... interface{}) {...}
func Exit(ctx context.Context, code int) {...}
func Fatal(ctx context.Context, args... interface{}) {...}
func Fatalf(ctx context.Context, format string, args... interface{}) {...}
func Info(ctx context.Context, args... interface{}) {...}
func Infof(ctx context.Context, format string, args... interface{}) {...}
func Log(ctx context.Context, level Level, args... interface{}) {...}
func Logf(ctx context.Context, level Level, format string, args... interface{}) {...}
func Panic(ctx context.Context, args... interface{}) {...}
func Panicf(ctx context.Context, format string, args... interface{}) {...}
func Print(ctx context.Context, args... interface{}) {...}
func Printf(ctx context.Context, format string, args... interface{}) {...}
func Trace(ctx context.Context, args... interface{}) {...}
func Tracef(ctx context.Context, format string, args... interface{}) {...}
func Warn(ctx context.Context, args... interface{}) {...}
func Warnf(ctx context.Context, format string, args... interface{}) {...}
```

These functions will extract logger from the context. To keep logger modular, a custom interface `Logger` will be made which can be implemented for any custom loggers. The logger interface itself will consist of trimmed logrus API similar to the above functions without context for the argument. Additional API can be added in the future should the need arises but so far, only these functions are used in the sysl engine. Initially Logrus' `StandardLogger` and `NullLogger` will be used as the underlying logger as they are used in the sysl engine frequently.
