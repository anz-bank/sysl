package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationErrorMessagesStringer(t *testing.T) {
	cases := map[string]struct {
		input    Msg
		expected string
	}{
		"Success no args": {input: Msg{MessageID: InfoValidatedSuccessfully}, expected: "Validation success"},
		"Success with args": {input: Msg{MessageID: ErrEntryPointUndefined, MessageData: []string{"Foo"}},
			expected: "Entry point view: (Foo) is undefined"},
		"Invalid message ID": {input: Msg{MessageID: 999}, expected: ""},
		"Lack args":          {input: Msg{MessageID: ErrEntryPointUndefined, MessageData: []string{}}, expected: ""},
		"Too many args": {input: Msg{MessageID: ErrEntryPointUndefined, MessageData: []string{"Foo", "Bar"}},
			expected: ""},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, expected, input.String(), "Unexpected result")
		})
	}
}

func TestValidationErrorMessagesGenLogMessage(t *testing.T) {
	cases := map[string]struct {
		input    Msg
		expected LogMessage
	}{
		"Success no args": {
			input:    Msg{MessageID: InfoValidatedSuccessfully},
			expected: LogMessage{Message: "Validation success", MsgType: INFO}},
		"Success with args": {
			input:    Msg{MessageID: ErrEntryPointUndefined, MessageData: []string{"Foo"}},
			expected: LogMessage{Message: "Entry point view: (Foo) is undefined", MsgType: ERROR}},
		"Invalid message ID": {
			input:    Msg{MessageID: 999},
			expected: LogMessage{Message: "", MsgType: UNDEF}},
		"Lack args": {
			input:    Msg{MessageID: ErrEntryPointUndefined, MessageData: []string{}},
			expected: LogMessage{Message: "", MsgType: ERROR}},
		"Too many args": {
			input:    Msg{MessageID: ErrEntryPointUndefined, MessageData: []string{"Foo", "Bar"}},
			expected: LogMessage{Message: "", MsgType: ERROR}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, expected, input.genLogMessage(), "Unexpected result")
		})
	}
}

func TestValidationErrorMessagesLogMsg(t *testing.T) {
	cases := map[string]struct {
		input Msg
	}{
		"Log error message":     {input: Msg{MessageID: ErrValidationFailed}},
		"Log warning message":   {input: Msg{MessageID: WarnValidatedWithWarn}},
		"Log info message":      {input: Msg{MessageID: InfoValidatedSuccessfully}},
		"Log unhandled message": {input: Msg{MessageID: 999}},
	}

	for name, test := range cases {
		input := test.input
		t.Run(name, func(t *testing.T) {
			input.logMsg()
		})
	}
}
