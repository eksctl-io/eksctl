package printers

import (
	"io"

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
	return t.table.Render(obj, writer, t.columnames...)
}

func (t *TablePrinter) AddColumn(name string, getter interface{}) {
	t.columnames = append(t.columnames, name)
	t.table.AddColumn(name, getter)
}
