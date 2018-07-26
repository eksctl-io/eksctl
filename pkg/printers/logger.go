package printers

import (
	"io"
)

type loggerFunc func(format string, a ...interface{})

// LoggerPrinter is a printer that outputs an object to
// the log output
type LoggerPrinter struct {
	loggerFn loggerFunc
}

// NewLoggerPrinter creates a new LoggerPrinter
func NewLoggerPrinter(log loggerFunc) *LoggerPrinter {
	return &LoggerPrinter{loggerFn: log}
}

// PrintObj will print the passed object formatted as object
// using the logger function.
func (l *LoggerPrinter) PrintObj(obj interface{}, writer io.Writer) error {
	l.loggerFn("%#v", obj)

	return nil
}
