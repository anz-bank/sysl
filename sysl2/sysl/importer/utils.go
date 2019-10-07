package importer

import "strings"

func getDescription(d string) string {
	if d == "" {
		return "No description."
	}
	return d
}

func quote(s string) string {
	if s == "" {
		return ""
	}
	return `"` + s + `"`
}

func isExternalAlias(item Type) bool {
	switch item.(type) {
	case *Alias, *Array, *Enum:
		return true
	}
	return false
}

func getSyslTypeName(item Type) string {
	switch t := item.(type) {
	case *Array:
		return "sequence of " + getSyslTypeName(t.Items)
	case *Enum:
		return item.Name()
	}
	if isExternalAlias(item) {
		return "EXTERNAL_" + item.Name()
	}
	return item.Name()
}

func spaceSeperate(items ...string) string {
	var t []string
	for _, i := range items {
		if i != "" {
			t = append(t, i)
		}
	}
	return strings.Join(t, " ")
}
