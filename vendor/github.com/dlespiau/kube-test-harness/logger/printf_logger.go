package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/dlespiau/kube-test-harness/testing"
)

// PrintfLogger is a logger using testing.T.Log for its output.
type PrintfLogger struct {
	baseLogger
}

var _ Logger = &PrintfLogger{}

// ForTest implements Logger.
func (l *PrintfLogger) ForTest(t testing.T) Logger {
	return l
}

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func decorate(s string) string {
	_, file, line, ok := runtime.Caller(3) // decorate + log + public function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}

// Log implements Logger.
func (l *PrintfLogger) Log(level LogLevel, msg string) {
	if !l.shouldLog(level) {
		return
	}
	fmt.Println(msg)
}

// Logf implements Logger.
func (l *PrintfLogger) Logf(level LogLevel, format string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}
	fmt.Print(decorate(fmt.Sprintf(format, args...)))
}
