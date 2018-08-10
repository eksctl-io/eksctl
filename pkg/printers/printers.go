package printers

import (
	"fmt"
	"io"
)

// OutputPrinter is the interface that printer must implement. This allows
// new printers to be added in the future.
type OutputPrinter interface {
	PrintObj(kind string, obj interface{}, writer io.Writer) error
}

// NewPrinter creates a new printer based in the printer type requested
// as a string.
func NewPrinter(printerType string) (OutputPrinter, error) {
	var printer OutputPrinter

	switch printerType {
	case "yaml":
		printer = NewYAMLPrinter()
	case "json":
		printer = NewJSONPrinter()
	case "table":
		printer = NewTablePrinter()
	default:
		return nil, fmt.Errorf("unknown output printer type: %s", printerType)
	}

	return printer, nil
}
