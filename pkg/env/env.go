package env

import "os"

// 	SYSL_PLANTUML
// 		URL of PlantUML server. Sysl depends upon
// 		[PlantUML](http://plantuml.com/) for diagram generation.
// 	SYSL_MODULES
// 		Whether the sysl modules is enabled.
// 		Enable by default, set to "off" to disable sysl modules.
// 	SYSL_TOKENS
// 		Tokens to use for git/github credentials + domains to use them on
// 		eg: SYSL_TOKENS=github.com:1234,gitlab.com:567
// 	SYSL_CACHE
// 		Cache location in current directory, defaults to "sysl-modules" if SYSL_MODULES is enabled
// 	SYSL_PROXY
// 		Proxy service to use, won't use SYSL_PROXY if not set

type VAR int

/* Use iotas instead of maps because this is thread safe */
//nolint:stylecheck,golint
const (
	SYSL_MODULES VAR = iota
	SYSL_PLANTUML
	SYSL_TOKENS
	SYSL_CACHE
	SYSL_PROXY
)

var VARS = []VAR{SYSL_MODULES, SYSL_PLANTUML, SYSL_TOKENS, SYSL_CACHE, SYSL_PROXY}

func (e VAR) Default() string {
	return [...]string{
		"on",                            //SYSL_MODULES
		"https://plantuml.com/plantuml", //SYSL_PLANTUML
		"",                              //SYSL_TOKENS
		"sysl_modules",                  //SYSL_CACHE
		"",                              //SYSL_TOKENS
	}[e]
}

func (e VAR) String() string {
	return [...]string{"SYSL_MODULES", "SYSL_PLANTUML", "SYSL_TOKENS", "SYSL_CACHE", "SYSL_PROXY"}[e]
}

func (e VAR) Value() string {
	if e := os.Getenv(e.String()); e != "" {
		return e
	}
	return e.Default()
}
