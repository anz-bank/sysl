package utils

import (
	"reflect"
	"sort"
	"unicode"
)

// Contains returns whether a string is in a string list
func Contains(needle string, haystack []string) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

// OrderedKeys takes a map and returns the ordered map keys
func OrderedKeys(mapObj interface{}) []string {
	var typeNames []string
	for _, k := range reflect.ValueOf(mapObj).MapKeys() {
		typeNames = append(typeNames, k.String())
	}
	sort.Strings(typeNames)
	return typeNames
}

// ToCamel camelcases the given string.
func ToCamel(in string) string {
	runes := []rune(in)
	var out []rune

	for i, r := range runes {
		if r == '_' {
			continue
		}
		if i == 0 || runes[i-1] == '_' {
			out = append(out, unicode.ToUpper(r))
			continue
		}
		out = append(out, r)
	}

	return string(out)
}
