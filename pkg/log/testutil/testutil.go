package testutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type LogCases struct {
	Name, Format string
	Arguments    []interface{}
}

type SingleField struct {
	Name, Key string
	Val       interface{}
}

type MultipleFields struct {
	Name   string
	Fields map[string]interface{}
}

func GenerateSingleFieldCases() []SingleField {
	return []SingleField{
		{
			Name: "String Value",
			Key:  "random",
			Val:  "Value",
		},
		{
			Name: "Number Value",
			Key:  "int",
			Val:  3,
		},
		{
			Name: "Byte Value",
			Key:  "byte",
			Val:  'q',
		},
		{
			Name: "Empty Key",
			Key:  "",
			Val:  "Empty",
		},
		{
			Name: "Empty Value",
			Key:  "Empty",
			Val:  "",
		},
		{
			Name: "Nil Value",
			Key:  "nil",
			Val:  nil,
		},
	}
}

func GenerateMultipleFieldsCases() []MultipleFields {
	return []MultipleFields{
		{
			Name: "Multiple types of Values",
			Fields: map[string]interface{}{
				"byte":   '1',
				"int":    123,
				"string": "this is an unnecessarily long sentence",
			},
		},
		{
			Name: "Empty Key",
			Fields: map[string]interface{}{
				"": "stuff",
			},
		},
		{
			Name: "Nil Value",
			Fields: map[string]interface{}{
				"Nil": nil,
			},
		},
		{
			Name:   "Nil Fields",
			Fields: nil,
		},
	}
}

func GenerateLogCases() []LogCases {
	return []LogCases{
		{
			Name:      "Single message",
			Arguments: []interface{}{"Test"},
			Format:    "%s",
		},
		{
			Name:      "Multiple messages",
			Arguments: []interface{}{"This", "is", "a", "test"},
			Format:    "%s %s %s %s",
		},
		{
			Name:      "Empty message",
			Arguments: []interface{}{""},
			Format:    "",
		},
		{
			Name:      "Multiple Types",
			Arguments: []interface{}{"test", 1, 'k'},
			Format:    "%s %v %v",
		},
		{
			Name:      "Empty Argument",
			Arguments: []interface{}{},
			Format:    "",
		},
	}
}

// Adapted from https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
func RedirectOutput(t *testing.T, print func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = w

	print()

	outC := make(chan string)
	go func(tt *testing.T) {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		require.NoError(tt, err)
		outC <- buf.String()
	}(t)

	w.Close()
	os.Stderr = old
	return <-outC
}

func OutputFormattedFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, len(fields))
	index := 0
	for fieldKey := range fields {
		keys[index] = fieldKey
		index++
	}

	sort.Strings(keys)

	output := strings.Builder{}
	output.WriteString(fmt.Sprintf("%s=%v", keys[0], fields[keys[0]]))

	if len(fields) > 1 {
		for _, keyField := range keys[1:] {
			output.WriteString(fmt.Sprintf(" %s=%v", keyField, fields[keyField]))
		}
	}

	return output.String()
}

func GetSortedKeys(fields map[string]interface{}) []string {
	keys := make([]string, len(fields))
	index := 0
	for key := range fields {
		keys[index] = key
		index++
	}
	sort.Strings(keys)
	return keys
}
