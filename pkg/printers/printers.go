package printers

import (
	"fmt"
	"io"

	"github.com/kris-nova/logger"
)

// Type is the type representing all supported printer types.
type Type = string

const (
	// YAMLType represents a printer of YAML type.
	YAMLType = Type("yaml")
	// JSONType represents a printer of JSON type.
	JSONType = Type("json")
	// TableType represents a printer of Table type.
	TableType = Type("table")
)

// OutputPrinter is the interface that printer must implement. This allows
// new printers to be added in the future.
type OutputPrinter interface {
	PrintObjWithKind(kind string, obj interface{}, writer io.Writer) error
	PrintObj(obj interface{}, writer io.Writer) error
	LogObj(log logger.Logger, msgFmt string, obj interface{}) error
}

// NewPrinter creates a new printer based in the printer type requested.
func NewPrinter(printerType Type) (OutputPrinter, error) {
	var printer OutputPrinter

	switch printerType {
	case YAMLType:
		printer = NewYAMLPrinter()
	case JSONType:
		printer = NewJSONPrinter()
	case TableType:
		printer = NewTablePrinter()
	default:
		return nil, errInvalidPrinterType(printerType)
	}

	return printer, nil
}

func errInvalidPrinterType(printerType Type) error {
	return fmt.Errorf("unknown output printer type: expected {%q,%q,%q} but got %q", YAMLType, JSONType, TableType, printerType)
}
