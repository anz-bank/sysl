package loggers

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	Logger
	mock.Mock
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Exit(code int) {
	m.Called(code)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Panic(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Panicf(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Print(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Printf(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Trace(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Tracef(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	arguments := append([]interface{}{format}, args...)
	m.Called(arguments...)
}

func (m *MockLogger) WithField(key string, val interface{}) Logger {
	return m.Called(key, val).Get(0).(Logger)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) {
	m.Called(fields)
}

func (m *MockLogger) Copy() Logger {
	return m.Called().Get(0).(Logger)
}
