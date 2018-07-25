package printers

import (
	"io"
)

type OutputPrinter interface {
	PrintObj(interface{}, io.Writer) error
}
