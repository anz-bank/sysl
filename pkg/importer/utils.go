package importer

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"
)

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

func isBuiltInType(item Type) bool {
	switch item.(type) {
	case *SyslBuiltIn, *ImportedBuiltInAlias:
		return true
	}
	return false
}

func isUnionType(item Type) bool {
	if _, ok := item.(*Union); ok {
		return true
	}
	return false
}

func isExternalAlias(item Type) bool {
	switch item.(type) {
	case *ExternalAlias, *Array, *Enum, *Alias:
		return true
	}
	return false
}

func getSyslTypeName(item Type) string {
	if item == nil {
		return ""
	}

	switch t := item.(type) {
	case *Array:
		return "sequence of " + getSyslTypeName(t.Items)
	case *Enum, *Alias, *Union:
		return item.Name()
	case *ImportedBuiltInAlias:
		return getSyslTypeName(t.Target)
	}
	if isExternalAlias(item) {
		return "EXTERNAL_" + item.Name()
	}
	if !isBuiltInType(item) {
		lower := strings.ToLower(item.Name())
		for _, bi := range syslutil.BuiltInTypes {
			if strings.HasPrefix(lower, bi) {
				return "_" + item.Name()
			}
		}
	}
	return item.Name()
}

func spaceSeparate(items ...string) string {
	var t []string
	for _, i := range items {
		if i != "" {
			t = append(t, i)
		}
	}
	return strings.Join(t, " ")
}

// escapeUnsafeSyslChars escapes special characters
// returns a string with URL encoded replacements
func escapeUnsafeSyslChars(name string) string {
	// url.PathEscape does not escape '. :'
	charsToReplace := map[string]string{
		".": `%2E`,
		":": `%3A`,
		"+": `%2B`,
		"$": `%24`,
		"&": `%26`,
	}
	name = url.PathEscape(name)
	for realChar, hex := range charsToReplace {
		name = strings.ReplaceAll(name, realChar, hex)
	}
	return name
}

func getSyslSafeName(name string) string {
	name = escapeUnsafeSyslChars(name)

	// Taken from `Name` in the .g4
	validSyslNameStart := regexp.MustCompile("^(%[0-9a-fA-F][0-9a-fA-F])*[a-zA-Z_]")

	if !validSyslNameStart.MatchString(name) {
		name = "_" + name
	}

	return name
}

func getSyslSafeURI(endpoint string) string {
	endpoint = escapeUnsafeSyslChars(endpoint)

	charsToKeep := map[string]string{
		`%2F`: "/",
		`%7B`: "{",
		`%7D`: "}",
		`%3D`: "=",
		`%3F`: "?",
		`%26`: "&",
	}
	for hex, realChar := range charsToKeep {
		endpoint = strings.ReplaceAll(endpoint, hex, realChar)
	}

	return endpoint
}

func cleanEndpointPath(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"{", "_",
		"}", "_",
		"-", "_").Replace(path)
}

func cleanMediaType(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"+", "_",
		"-", "_",
		".", "_",
		"*", "_").Replace(path)
}

func convertToSyslSafe(name string) string {
	if !strings.ContainsAny(name, "- ") {
		return name
	}

	syslSafe := strings.Builder{}
	toUppercase := false
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '-':
			toUppercase = true
		case ' ':
			continue
		default:
			if toUppercase {
				syslSafe.WriteString(strings.ToUpper(string(name[i])))
				toUppercase = false
			} else {
				syslSafe.WriteByte(name[i])
			}
		}
	}
	return syslSafe.String()
}
