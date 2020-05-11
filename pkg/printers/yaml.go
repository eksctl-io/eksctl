package printers

import (
	"bytes"
	"io"
	"strings"

	"github.com/kris-nova/logger"

	"k8s.io/apimachinery/pkg/runtime"
	cliruntime "k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/yaml"
)

// YAMLPrinter is a printer that outputs an object formatted
// as YAML
type YAMLPrinter struct {
	runtimePrinter cliruntime.YAMLPrinter
}

// NewYAMLPrinter creates a new YAMLPrinter
func NewYAMLPrinter() OutputPrinter {
	return &YAMLPrinter{cliruntime.YAMLPrinter{}}
}

// PrintObj will print the passed object formatted as YAML to
// the supplied writer.
func (y *YAMLPrinter) PrintObj(obj interface{}, writer io.Writer) error {
	if obj, ok := obj.(runtime.Object); ok {
		if err := y.runtimePrinter.PrintObj(obj, writer); err == nil {
			// if an error occurred, we may still be able to serialise using yaml package directly
			return nil
		}
	}

	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := writer.Write(b); err != nil {
		return err
	}

	return nil
}

// PrintObjWithKind will print the passed object formatted as YAML to
// the supplied writer. This printer ignores kind argument.
func (y *YAMLPrinter) PrintObjWithKind(kind string, obj interface{}, writer io.Writer) error {
	return y.PrintObj(obj, writer)
}

// LogObj will print the passed object formatted as YAML to
// the logger.
func (y *YAMLPrinter) LogObj(log logger.Logger, msgFmt string, obj interface{}) error {
	b := &bytes.Buffer{}
	if err := y.PrintObj(obj, b); err != nil {
		return err
	}

	log(msgFmt, strings.ReplaceAll(b.String(), "%", "%%"))

	return nil
}
