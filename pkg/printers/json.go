package printers

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/kris-nova/logger"

	"k8s.io/apimachinery/pkg/runtime"
	cliruntime "k8s.io/cli-runtime/pkg/genericclioptions/printers"
)

// JSONPrinter is a printer that outputs an object formatted
// as JSON
type JSONPrinter struct {
	runtimePrinter cliruntime.JSONPrinter
}

// NewJSONPrinter creates a new JsonPrinter with defaults for
// level of indentation.
func NewJSONPrinter() OutputPrinter {
	return &JSONPrinter{cliruntime.JSONPrinter{}}
}

// PrintObj will print the passed object formatted as JSON to
// the supplied writer.
func (j *JSONPrinter) PrintObj(obj interface{}, writer io.Writer) error {
	if obj, ok := obj.(runtime.Object); ok {
		if err := j.runtimePrinter.PrintObj(obj, writer); err == nil {
			// if an error occurred, we may still be able to serialise using json package directly
			return nil
		}
	}

	b, err := json.MarshalIndent(obj, "", "    ")
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
func (j *JSONPrinter) PrintObjWithKind(kind string, obj interface{}, writer io.Writer) error {
	return j.PrintObj(obj, writer)
}

// LogObj will print the passed object formatted as JSON to
// the logger.
func (j *JSONPrinter) LogObj(log logger.Logger, msgFmt string, obj interface{}) error {
	b := &bytes.Buffer{}
	if err := j.PrintObj(obj, b); err != nil {
		return err
	}

	log(msgFmt, strings.ReplaceAll(b.String(), "%", "%%"))

	return nil
}
