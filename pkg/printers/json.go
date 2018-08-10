package printers

import (
	"encoding/json"
	"io"
)

const (
	DefaultPrefix     = ""
	DefaultIdentation = "  "
)

// JSONPrinter is a printer that outputs an object formatted
// as JSON
type JSONPrinter struct {
	prefix string
	indent string
}

// NewJSONPrinter creates a new JsonPrinter with defaults for
// level of identation.
func NewJSONPrinter() OutputPrinter {
	return &JSONPrinter{prefix: DefaultPrefix, indent: DefaultIdentation}
}

// PrintObj will print the passed object formatted as JSON to
// the supplied writer.
func (j *JSONPrinter) PrintObj(kind string, obj interface{}, writer io.Writer) error {
	b, err := json.MarshalIndent(obj, j.prefix, j.indent)
	if err != nil {
		return err
	}
	writer.Write(b)

	return nil
}
