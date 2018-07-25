package printers

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	DefaultTablePadding = 1
)

var (
	timeType reflect.Type
)

func init() {
	timeType = reflect.TypeOf(time.Time{})
}

// TablePrinter is a printer that outputs an object formatted
// as table
type TablePrinter struct {
	padding   int
	padchar   byte
	tabwriter *tabwriter.Writer
}

// NewTablePrinter creates a new TablePrinter with defaults.
func NewTablePrinter() *TablePrinter {
	return &TablePrinter{padding: DefaultTablePadding, padchar: ' '}
}

// PrintObj will print the passed object formatted as textual
// table to the supplied writer.
func (t *TablePrinter) PrintObj(obj interface{}, writer io.Writer) error {
	in := reflect.ValueOf(obj).Elem()
	if in.Kind() != reflect.Struct {
		return fmt.Errorf("PrintObj can only print structs. The kind of object passed in was: %s", in.Kind())
	}

	t.tabwriter = tabwriter.NewWriter(writer, 0, 0, t.padding, t.padchar, tabwriter.Debug)

	if err := t.processStruct(in, 0); err != nil {
		return err
	}

	return t.tabwriter.Flush()
}

func (t *TablePrinter) processStruct(in reflect.Value, level int) error {
	for i := 0; i < in.NumField(); i++ {
		f := in.Field(i)
		name := in.Type().Field(i).Name

		// Make sure we are only dealing with exported fields
		if !f.CanInterface() {
			continue
		}

		t.processField(name, f, level)
	}

	return nil
}

func (t *TablePrinter) processField(name string, in reflect.Value, level int) {
	if in.Kind() == reflect.Ptr && in.IsNil() {
		t.printField(level, name, "<nil>")
		return
	}

	val := reflect.ValueOf(in.Interface())

	//fmt.Printf("Name: %s, Kind: %s, Level: %d\n", name, val.Kind(), level)
	switch val.Kind() {
	case reflect.Interface:
		t.processField(name, val.Elem(), level)
	case reflect.Ptr:
		t.processField(name, val.Elem(), level)
	case reflect.Struct:
		if val.Type() == timeType {
			t.processTime(name, val, level)
		} else {
			t.printField(level, name, "")
			newLevel := level + 1
			t.processStruct(val, newLevel)
		}
	case reflect.Slice, reflect.Array:
		t.processSlice(name, val, level)
	default:
		t.printField(level, name, val.Interface())
	}
}

func (t *TablePrinter) processSlice(name string, in reflect.Value, level int) {
	numItems := in.Len()
	for i := 0; i < numItems; i++ {
		t.processField(name, in.Index(i), level)
	}
}

func (t *TablePrinter) processTime(name string, in reflect.Value, level int) {
	tt := in.Interface().(time.Time)
	formatted := tt.Format(time.RFC3339Nano)
	t.printField(level, name, formatted)
}

func (t *TablePrinter) printField(indent int, fieldName string, filedValue interface{}) {
	prefix := strings.Repeat("\t", indent)
	fmt.Fprintf(t.tabwriter, "%s%s\t%v\n", prefix, fieldName, filedValue)
}
