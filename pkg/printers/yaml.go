package printers

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

// YAMLPrinter is a printer that outputs an object formatted
// as YAML
type YAMLPrinter struct{}

// NewYAMLPrinter creates a new YAMLPrinter
func NewYAMLPrinter() *YAMLPrinter {
	return &YAMLPrinter{}
}

// PrintObj will print the passed object formatted as YAML to
// the supplied writer.
func (j *YAMLPrinter) PrintObj(obj interface{}, writer io.Writer) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	writer.Write(b)

	return nil
}
