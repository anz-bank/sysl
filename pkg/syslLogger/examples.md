# Examples of Usage

## Setup

```go
package main

import (
    "context"

    sysllogger "github.com/anz-bank/sysl/pkg/syslLogger"
    "github.com/anz-bank/sysl/pkg/syslLogger/loggers"
)

func main() {
    // User is expected to choose a logger and add it to the context using the library's API
    ctx := context.Background()

    // this is a logger based on the logrus standard logger
    logger := loggers.NewStandardLogger()

    // AddLogger returns a new context
    ctx = sysllogger.AddLogger(ctx, logger)
}
```

That's all in setup, now logging can be used by using the context.

## Usage

```go
import (
    sysllogger "github.com/anz-bank/sysl/pkg/syslLogger"
)

func stuffToLog(ctx context.Context) {
    // logging uses the context variable so it must be given to any function that requires it
    sysllogger.Debug(ctx, "Debug")
    sysllogger.Print(ctx, "Print")
    sysllogger.Trace(ctx, "Trace")
    sysllogger.Warn(ctx, "Warn")
    sysllogger.Error(ctx, "Error")
    sysllogger.Fatal(ctx, "Fatal")
    sysllogger.Panic(ctx, "Panic")

    /**
     * Expected to log
     * (time in RFC3339Nano Format) (Level) (Message)
     *
     * Example:
     * 2019-12-12T08:23:59.210878+11:00 PRINT Hello There
     *
     * Each API also has its Format counterpart (Debugf, Printf, Tracef, Warnf, Errorf, Fatalf, Panicf)
     */
}
```

Fields are also supported in the logging. There are two kinds of fields, context-level field and log-level field.

```go
import (
    sysllogger "github.com/anz-bank/sysl/pkg/syslLogger"
)

    /**
     * With fields, it is expected to log
     * (time in RFC3339Nano Format) (Fields) (Level) (Message)
     *
     * Fields will be logged ALPHABETICALLY. If the same key field is added to the context logger,
     * it will replace the existing value that corresponds to that key.
     *
     * Example:
     * 2019-12-12T08:23:59.210878+11:00 random=stuff very=random PRINT Hello There
     *
     * Each API also has its Format counterpart (Debugf, Printf, Tracef, Warnf, Errorf, Fatalf, Panicf)
     */

func logWithField(ctx context.Context) {
    // context-level field adds fields to the context and creates a new context
    ctx = sysllogger.AddField(ctx, "random", "stuff")
    ctx = sysllogger.AddFields(ctx, map[string]interface{}{
        "just": "stuff",
        "stuff": 1
    })

    /**
     * Any log at this point will have fields and to any function that uses the same context
     * just=stuff random=stuff stuff=1
     */

    contextLevelField(ctx)
    logLevelField(ctx)
}

func contextLevelField(ctx context.Context) {
    /**
     * This is expected to log something like
     *
     * 2019-12-12T08:23:59.210878+11:00 just=stuff random=stuff stuff=1 WARN Warn
     */
    sysllogger.Warn(ctx, "Warn")
}

func logLevelField(ctx context.Context) {
    /**
     * Log level field returns a logger that does not modify the existing context
     * The APIs for Log level fields are meant to directly interact with the Logger API without
     * modifying context. While it is possible to extract logger and use it directly, that is
     * NOT RECOMMENDED.
     */

    /**
     * Log level fields will add fields on top of the existing context level fields
     * This is expected to log something like
     *
     * 2019-12-12T08:23:59.210878+11:00 just=stuff more=random stuff random=stuff stuff=1 very=random WARN Warn
     */
    sysllogger.LogFields(ctx, map[string]interface{}{
        "more": "random stuff",
        "very": "random"
    }).Warn("Warn")

    /**
     * With LogField it is possible to chain LogField with more Field related API that Logger has.
     * It also replaces existing value but only in log level, not context level.
     * This is expected to log something like
     *
     * 2019-12-12T08:23:59.210878+11:00 epicly=random just=stuff more=stuff random=crap stuff=1 WARN Warn
     */

    sysllogger.LogField(ctx, "epicly", "random").WithField("more", "stuff").WithField("random", "crap").Warn("Warn")

    /**
     * As long as context logger is not modified, it will log again the context level field
     */
    contextLevelField(ctx)
}

```