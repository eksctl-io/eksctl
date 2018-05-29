package lol

import (
	"fmt"
)

var w = &Writer{Output: stdout, ColorMode: ColorMode256}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, a...)
}
