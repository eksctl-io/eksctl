package logger

import "github.com/dlespiau/kube-test-harness/testing"

type baseLogger struct {
	level LogLevel
	t     testing.T
}

// SetLevel implements Logger.
func (l *baseLogger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel implements Logger.
func (l *baseLogger) GetLevel() LogLevel {
	return l.level
}

func (l *baseLogger) shouldLog(level LogLevel) bool {
	return level >= l.level
}
