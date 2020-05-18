package mermaid

import "strings"

//cleanString replaces certain characters in the string suitable for mermaid
func CleanString(temp string) string {
	temp = strings.ReplaceAll(temp, " ", "")
	temp = strings.ReplaceAll(temp, "{", "_")
	temp = strings.ReplaceAll(temp, "}", "_")
	temp = strings.ReplaceAll(temp, "[", "_")
	temp = strings.ReplaceAll(temp, "]", "_")
	temp = strings.ReplaceAll(temp, "\"", "")
	temp = strings.ReplaceAll(temp, "~", "")
	temp = strings.ReplaceAll(temp, ":", "_")
	return temp
}
