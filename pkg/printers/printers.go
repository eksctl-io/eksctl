package printers

import (
	"fmt"
	"io"
)

type OutputPrinter interface {
	PrintObj(interface{}, io.Writer) error
}

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
