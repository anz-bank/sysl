package main

import (
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
)

func getAppName(appname *sysl.AppName) string {
	return strings.Join(appname.Part, " :: ")
}

func getApp(appName *sysl.AppName, mod *sysl.Module) *sysl.Application {
	return mod.Apps[getAppName(appName)]
}

func hasAbstractPattern(attrs map[string]*sysl.Attribute) bool {
	patterns, has := attrs["patterns"]
	if has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				if y.GetS() == "abstract" {
					return true
				}
			}
		}
	}
	return false
}

func isSameApp(a *sysl.AppName, b *sysl.AppName) bool {
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

func isSameCall(a *sysl.Call, b *sysl.Call) bool {
	return isSameApp(a.Target, b.Target) && a.Endpoint == b.Endpoint
}
