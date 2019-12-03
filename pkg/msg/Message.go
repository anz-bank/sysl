package msg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// Msg holds message id and params to generate message string
type Msg struct {
	MessageID   int      `json:"message_id"`
	MessageData []string `json:"message_data"`
}

// LogMessage holds message string and type of message
type LogMessage struct {
	Message string `json:"message"`
	MsgType int    `json:"msg_type"`
}

const (
	ERROR = 100
	WARN  = 200
	INFO  = 300
	TITLE = 900
	UNDEF = 500

	ErrValidationFailed        = 400
	ErrEntryPointUndefined     = 401
	ErrInvalidEntryPointReturn = 402
	ErrUndefinedView           = 403
	ErrInvalidReturn           = 404
	ErrMissingReqField         = 405
	ErrExcessAttr              = 406
	ErrInvalidOption           = 407
	ErrInvalidUnary            = 408
	ErrRedefined               = 409

	WarnValidatedWithWarn = 300
	WarnValidationSkipped = 301

	InfoValidatedSuccessfully = 200
	TitleViewName             = 999
)

// nolint:gochecknoglobals
var (
	Messages = map[int]string{
		ErrValidationFailed:        "Validation failed",
		ErrEntryPointUndefined:     "Entry point view: (%s) is undefined",
		ErrInvalidEntryPointReturn: "Return type of entry point view: (%s) should be (%s)",
		ErrUndefinedView:           "View (%s) is undefined",
		ErrInvalidReturn:           "In view (%s), return type is invalid. (%s)",
		ErrMissingReqField:         "Missing required field (%s) in View (%s) with return Type (%s)",
		ErrExcessAttr:              "Excess Attribute (%s), defined in View (%s), having return Type (%s)",
		ErrInvalidOption:           "In View %s, (%s) does not match any of the options for return Type (%s)",
		ErrInvalidUnary:            "In view (%s), unary operator used with invalid type: (%s)",
		ErrRedefined:               "In view (%s), (%s) is already defined",

		WarnValidatedWithWarn: "Validated with warnings",
		WarnValidationSkipped: "Validation skipped. Reason: %s",

		InfoValidatedSuccessfully: "Validation success",

		TitleViewName: "Error in %s",
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
		ErrRedefined:               ERROR,

		WarnValidatedWithWarn: WARN,
		WarnValidationSkipped: WARN,

		InfoValidatedSuccessfully: INFO,
		TitleViewName:             TITLE,
	}
)

func NewMsg(messageID int, messageData []string) *Msg {
	return &Msg{MessageID: messageID, MessageData: messageData}
}

func (m Msg) String() string {
	message, exists := Messages[m.MessageID]

	if !exists {
		return fmt.Sprintf("[validationError.String] invalid message ID: %v", m.MessageID)
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
	if messageType, exists := MessageType[m.MessageID]; exists {
		return LogMessage{Message: m.String(), MsgType: messageType}
	}
	return LogMessage{Message: m.String(), MsgType: UNDEF}
}

// LogMsg prints the message using logrus
func (m Msg) LogMsg() {
	formattedMsg := "[Validator]: " + m.String()
	switch MessageType[m.MessageID] {
	case ERROR:
		logrus.Error(formattedMsg)
	case WARN:
		logrus.Warn(formattedMsg)
	case INFO:
		logrus.Info(formattedMsg)
	case TITLE:
		logrus.Println(generateTitle(strings.Split(m.String(), ":")[0]))
	default:
		logrus.Println(formattedMsg)
	}
}

func generateTitle(s string) string {
	str := s
	length := len(s)

	if dash := (100 - (length)) / 2; dash > 0 {
		dashStr := ""
		for i := 0; i < dash-1; i++ {
			dashStr += "-"
		}
		str = fmt.Sprintf("%s %s %s", dashStr, s, dashStr)
	}

	if length < 100 && length%2 == 1 {
		str += "-"
	}

	return str
}
