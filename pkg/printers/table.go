package printers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/kops/util/pkg/tables"
)

// TablePrinter is a printer that outputs an object formatted
// as table
type TablePrinter struct {
	table      *tables.Table
	columnames []string
}

// NewTablePrinter creates a new TablePrinter with defaults.
func NewTablePrinter() OutputPrinter {
	return &TablePrinter{table: &tables.Table{}}
}

// PrintObj will print the passed object formatted as textual
// table to the supplied writer.
func (t *TablePrinter) PrintObj(obj interface{}, writer io.Writer) error {
	return t.PrintObjWithKind("objects", obj, writer)
}

// PrintObjWithKind will print the passed object formatted as textual
// table to the supplied writer.
func (t *TablePrinter) PrintObjWithKind(kind string, obj interface{}, writer io.Writer) error {
	itemsValue := reflect.ValueOf(obj)
	if itemsValue.Kind() != reflect.Slice {
		return errors.Errorf("table printer expects a slice but the kind was %v", itemsValue.Kind())
	}

	if itemsValue.Len() == 0 {
		w := bufio.NewWriter(writer)
		if _, err := w.WriteString(fmt.Sprintf("No %s found\n", strings.ToLower(kind))); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
		return nil
	}

	return t.table.Render(obj, writer, t.columnames...)
}

// LogObj will print the passed object formatted as a table to
// the logger.
func (t *TablePrinter) LogObj(log logger.Logger, prefixFmt string, obj interface{}) error {
	b := &bytes.Buffer{}
	if err := t.PrintObj(obj, b); err != nil {
		return err
	}

	log(prefixFmt+"%s", b.String())

	return nil
}

// AddColumn adds a column to the table that will be printed
func (t *TablePrinter) AddColumn(name string, getter interface{}) {
	t.columnames = append(t.columnames, name)
	t.table.AddColumn(name, getter)
}
