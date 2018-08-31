package logger

import "github.com/dlespiau/kube-test-harness/testing"

// LogLevel defines how verbose the Logger is.
type LogLevel int

const (
	// Debug will display all logs.
	Debug LogLevel = 1
	// Info will display only informational logs.
	Info LogLevel = 2
)

// Logger can output logs when running tests.
type Logger interface {
	ForTest(t testing.T) Logger
	SetLevel(level LogLevel)
	GetLevel() LogLevel
	Log(level LogLevel, msg string)
	Logf(level LogLevel, fmt string, args ...interface{})
}
