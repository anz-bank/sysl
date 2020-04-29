// alphabetical.go are the functions to get a sorted list of keys from sysls map types
// Unfortunately go doesn't allow generics and pointer types don't implement empty interfaces
// so map[string]interface{} doesn't cut it

package printer

import (
	"sort"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func alphabeticalAttributes(m map[string]*sysl.Attribute) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalApplications(m map[string]*sysl.Application) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalEndpoints(m map[string]*sysl.Endpoint) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalTypes(m map[string]*sysl.Type) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalInts(m map[string]int64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
