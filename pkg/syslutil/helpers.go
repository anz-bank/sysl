package syslutil

import (
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

// LogLevels ...
// nolint:gochecknoglobals
var LogLevels = map[string]logrus.Level{
	"":      logrus.ErrorLevel,
	"off":   logrus.ErrorLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"trace": logrus.TraceLevel,
}

func CleanAppName(name string) string {
	parts := strings.Split(name, "::")
	for i := range parts {
		parts[i] = strings.Trim(parts[i], " \t")
	}
	return GetAppName(&sysl.AppName{Part: parts})
}

func GetAppName(appname *sysl.AppName) string {
	return strings.Join(appname.Part, " :: ")
}

func GetApp(appName *sysl.AppName, mod *sysl.Module) *sysl.Application {
	return mod.Apps[GetAppName(appName)]
}

func HasPattern(attrs map[string]*sysl.Attribute, pattern string) bool {
	patterns, has := attrs["patterns"]
	if has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				if y.GetS() == pattern {
					return true
				}
			}
		}
	}
	return false
}

func IsSameApp(a *sysl.AppName, b *sysl.AppName) bool {
	if len(a.Part) != len(b.Part) {
		return false
	}
	for i := range a.Part {
		if a.Part[i] != b.Part[i] {
			return false
		}
	}
	return true
}

func IsSameCall(a *sysl.Call, b *sysl.Call) bool {
	return IsSameApp(a.Target, b.Target) && a.Endpoint == b.Endpoint
}
