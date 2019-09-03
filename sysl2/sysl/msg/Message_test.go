package msg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestValidationMsgStringer(t *testing.T) {
	t.Parallel()

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

func TestValidationMsgGenLogMessage(t *testing.T) {
	t.Parallel()

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
			input:    Msg{MessageID: 900},
			expected: LogMessage{Message: "[validationError.String] invalid message ID: 900", MsgType: UNDEF}},
		"Lacks args": {
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

func TestValidationMsgLogMsg(t *testing.T) {
	t.Parallel()

	defer logrus.SetOutput(logrus.StandardLogger().Out)

	cases := map[string]struct {
		input    Msg
		expected string
	}{
		"Log error message": {
			input:    Msg{MessageID: ErrValidationFailed},
			expected: "level=error msg=\"[Validator]: Validation failed\"\n"},
		"Log warning message": {
			input:    Msg{MessageID: WarnValidatedWithWarn},
			expected: "level=warning msg=\"[Validator]: Validated with warnings\"\n"},
		"Log info message": {
			input:    Msg{MessageID: InfoValidatedSuccessfully},
			expected: "level=info msg=\"[Validator]: Validation success\"\n"},
		"Log title": {
			input: Msg{MessageID: TitleViewName, MessageData: []string{"Title"}},
			expected: "level=info msg=\"------------------------------------------ Error in Title " +
				"------------------------------------------\"\n"},
		"Log unhandled message": {
			input:    Msg{MessageID: 900},
			expected: "level=info msg=\"[Validator]: [validationError.String] invalid message ID: 900\"\n"},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			logrus.SetOutput(&buf)
			input.LogMsg()
			assert.True(t, strings.Contains(buf.String(), expected))
		})
	}
}

func TestNewMsg(t *testing.T) {
	t.Parallel()

	msg := NewMsg(InfoValidatedSuccessfully, []string{})
	assert.Equal(t, InfoValidatedSuccessfully, msg.MessageID)
	assert.Equal(t, []string{}, msg.MessageData)
}

func TestGenerateTitle(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input    string
		expected string
	}{
		"Length 95": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: "- aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa --"},
		"Length 96": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: "- aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa -"},
		"Length 97": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: " aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa -"},
		"Length 98": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: " aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa "},
		"Length 99": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-"},
		"Length 100": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		"Length 101": {
			input:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			title := generateTitle(input)
			assert.Equal(t, expected, title)
		})
	}
}
