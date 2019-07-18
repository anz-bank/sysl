package main

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

type Msg struct {
	MessageID   int      `json:"message_id"`
	MessageData []string `json:"message_data"`
}

type LogMessage struct {
	Message string `json:"message"`
	MsgType int    `json:"msg_type"`
}

const (
	ERROR = 100
	WARN  = 200
	INFO  = 300
	UNDEF = 400

	ErrValidationFailed        = 400
	ErrEntryPointUndefined     = 401
	ErrInvalidEntryPointReturn = 402
	ErrUndefinedView           = 403
	ErrInvalidReturn           = 404
	ErrMissingReqField         = 405
	ErrExcessAttr              = 406
	ErrInvalidOption           = 407
	ErrInvalidUnary            = 408

	WarnValidatedWithWarn = 300

	InfoValidatedSuccessfully = 200
)

//nolint:gochecknoglobals
var (
	Messages = map[int]string{
		ErrValidationFailed:        "Validation failed",
		ErrEntryPointUndefined:     "Entry point view: (%s) is undefined",
		ErrInvalidEntryPointReturn: "Return type of entry point view: (%s) should be (%s)",
		ErrUndefinedView:           "View (%s) is undefined",
		ErrInvalidReturn:           "In view (%s), return type should be (%s)",
		ErrMissingReqField:         "Missing required field (%s) in View (%s) with return Type (%s)",
		ErrExcessAttr:              "Excess Attribute (%s), defined in View (%s), having return Type (%s)",
		ErrInvalidOption:           "In View %s, (%s) does not match any of the options for return Type (%s)",
		ErrInvalidUnary:            "In view (%s), unary operator used with invalid type: (%s)",

		WarnValidatedWithWarn: "Validated with warnings",

		InfoValidatedSuccessfully: "Validation success",
	}

	MessageType = map[int]int{
		ErrValidationFailed:        ERROR,
		ErrEntryPointUndefined:     ERROR,
		ErrInvalidEntryPointReturn: ERROR,
		ErrUndefinedView:           ERROR,
		ErrInvalidReturn:           ERROR,
		ErrMissingReqField:         ERROR,
		ErrExcessAttr:              ERROR,
		ErrInvalidOption:           ERROR,
		ErrInvalidUnary:            ERROR,

		WarnValidatedWithWarn: WARN,

		InfoValidatedSuccessfully: INFO,
	}
)

func NewMsg(messageID int, messageData []string) *Msg {
	return &Msg{MessageID: messageID, MessageData: messageData}
}

func (m Msg) String() string {
	message, exists := Messages[m.MessageID]

	if !exists {
		fmt.Println("[validationError.String] invalid message ID")
		return ""
	}

	argPattern := regexp.MustCompile("%[sv]")
	matches := argPattern.FindAllStringIndex(message, -1)
	matchCount := len(matches)
	argCount := len(m.MessageData)

	switch {
	case matchCount > argCount:
		fmt.Printf("[validationError.String] Insufficient args provided. expected %d, provided %d\n", matchCount, argCount)
		return ""
	case matchCount < argCount:
		fmt.Printf("[validationError.String] Too many args provided. expected %d, provided %d", matchCount, argCount)
		return ""
	}

	args := make([]interface{}, argCount)
	for i, arg := range m.MessageData {
		args[i] = arg
	}
	return fmt.Sprintf(message, args...)
}

func (m Msg) genLogMessage() LogMessage {
	messageType, exists := MessageType[m.MessageID]

	if exists {
		return LogMessage{Message: m.String(), MsgType: messageType}
	}
	return LogMessage{Message: m.String(), MsgType: UNDEF}
}

func (m Msg) logMsg() {
	formattedMsg := "[Validator]: " + m.String()
	switch MessageType[m.MessageID] {
	case ERROR:
		logrus.Error(formattedMsg)
	case WARN:
		logrus.Warn(formattedMsg)
	case INFO:
		logrus.Info(formattedMsg)
	default:
		fmt.Println("[validationError.logMsg] Unhandled message", formattedMsg)
	}
}
