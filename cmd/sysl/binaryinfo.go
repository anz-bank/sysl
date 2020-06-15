package main

import (
	"fmt"
	"runtime"
)

// Binary info variables are dynamically injected via the `-ldflags` flag with `go build`
// Version   - Binary version
// GitFullCommit - Commit SHA of the source code
// BuildDate - Binary build date
// GoVersion - Binary build GoVersion
// BuildOS   - Operating System used to build binary
//nolint:gochecknoglobals
var (
	Version       = "unspecified"
	GitFullCommit = "unspecified"
	BuildDate     = "unspecified"
	GoVersion     = "unspecified"
	BuildOS       = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)
