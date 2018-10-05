package printers

import (
	"io"

	"github.com/ghodss/yaml"
)

// YAMLPrinter is a printer that outputs an object formatted
// as YAML
type YAMLPrinter struct{}

// NewYAMLPrinter creates a new YAMLPrinter
func NewYAMLPrinter() OutputPrinter {
	return &YAMLPrinter{}
}

// PrintObj will print the passed object formatted as YAML to
// the supplied writer.
func (j *YAMLPrinter) PrintObj(kind string, obj interface{}, writer io.Writer) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := writer.Write(b); err != nil {
		return err
	}

	return nil
}
