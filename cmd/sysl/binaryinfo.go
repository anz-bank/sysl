package main

// Binary info variables are dynamically injected via the `-ldflags` flag with `go build`
// Version   - Binary version
// GitCommit - Commit SHA of the source code
// BuildDate - Binary build date
// GoVersion - Binary build GoVersion
// BuildOS   - Operating System used to build binary
//nolint:gochecknoglobals
var (
	Version   = "unspecified"
	GitCommit = "unspecified"
	BuildDate = "unspecified"
	GoVersion = "unspecified"
	BuildOS   = "unspecified"
)
