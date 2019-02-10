package main

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func assertNoLog(t *testing.T, level logrus.Level, f func()) {
	var buf bytes.Buffer
	std := logrus.StandardLogger()

	savedOut := std.Out
	std.SetOutput(&buf)
	defer std.SetOutput(savedOut)

	savedLevel := std.Level
	std.SetLevel(level)
	defer std.SetLevel(savedLevel)

	f()

	assert.Empty(t, buf.String())
}

func assertLog(t *testing.T, msg string, level logrus.Level, f func()) {
	var buf bytes.Buffer
	std := logrus.StandardLogger()

	savedOut := std.Out
	std.SetOutput(&buf)
	defer std.SetOutput(savedOut)

	savedLevel := std.Level
	std.SetLevel(level)
	defer std.SetLevel(savedLevel)

	f()

	timeRE := `time="\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(?:Z|[-+]\d\d:\d\d)?"`
	expected := fmt.Sprintf(` level=%s msg="%s"\n`, level, regexp.QuoteMeta(msg))
	assert.Regexp(t, timeRE+expected, buf.String())
}

type dummyRecognizer struct {
	*antlr.BaseRecognizer
}

func (r *dummyRecognizer) GetATN() *antlr.ATN {
	panic("Not implemented")
}

func TestSyslParserErrorListenerSyntaxError(t *testing.T) {
	listener := &SyslParserErrorListener{}
	recognizer := &dummyRecognizer{antlr.NewBaseRecognizer()}
	recognizer.SymbolicNames = []string{"some_token", "some_other_token"}
	source := &antlr.TokenSourceCharStreamPair{}
	offendingSymbol := antlr.NewCommonToken(source, 1, 0, 0, 0)
	assertNoLog(t, logrus.WarnLevel, func() {
		listener.SyntaxError(recognizer, offendingSymbol, 1, 1, "some error", nil)
	})
	assertLog(t, `SyntaxError: Token: some_other_token\n`, logrus.InfoLevel, func() {
		listener.SyntaxError(recognizer, offendingSymbol, 1, 1, "some error", nil)
	})
}

func TestSyslParserErrorListenerReportAttemptingFullContext(t *testing.T) {
	listener := &SyslParserErrorListener{}
	assertNoLog(t, logrus.WarnLevel, func() {
		listener.ReportAttemptingFullContext(nil, nil, 42, 43, nil, nil)
	})
	assertLog(t, `ReportAttemptingFullContext: 42 43\n`, logrus.InfoLevel, func() {
		listener.ReportAttemptingFullContext(nil, nil, 42, 43, nil, nil)
	})
}

func TestSyslParserErrorListenerReportAmbiguity(t *testing.T) {
	listener := &SyslParserErrorListener{}
	assertNoLog(t, logrus.WarnLevel, func() {
		listener.ReportAmbiguity(nil, nil, 42, 43, false, nil, nil)
	})
	assertLog(t, `ReportAmbiguity: 42 43\n`, logrus.InfoLevel, func() {
		listener.ReportAmbiguity(nil, nil, 42, 43, false, nil, nil)
	})
}

func TestSyslParserErrorListenerReportContextSensitivity(t *testing.T) {
	listener := &SyslParserErrorListener{}
	assertNoLog(t, logrus.WarnLevel, func() {
		listener.ReportContextSensitivity(nil, nil, 42, 43, 0, nil)
	})
	assertLog(t, `ReportContextSensitivity: 42 43\n`, logrus.InfoLevel, func() {
		listener.ReportContextSensitivity(nil, nil, 42, 43, 0, nil)
	})
}
