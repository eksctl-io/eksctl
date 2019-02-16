package goptions

import (
	"fmt"
	"reflect"
	"strings"
)

// Flag represents a single flag of a FlagSet.
type Flag struct {
	Short        string
	Long         string
	MutexGroups  []string
	Description  string
	Obligatory   bool
	WasSpecified bool
	value        reflect.Value
	optionMeta   map[string]interface{}
	DefaultValue interface{}
}

// Return the name of the flag preceding the right amount of dashes.
// The long name is preferred. If no name has been specified, "<unspecified>"
// will be returned.
func (f *Flag) Name() string {
	if len(f.Long) > 0 {
		return "--" + f.Long
	}
	if len(f.Short) > 0 {
		return "-" + f.Short
	}
	return "<unspecified>"
}

// NeedsExtraValue returns true if the flag expects a separate value.
func (f *Flag) NeedsExtraValue() bool {
	// Explicit over implicit
	if f.value.Type() == reflect.TypeOf(new([]bool)).Elem() ||
		f.value.Type() == reflect.TypeOf(new(bool)).Elem() {
		return false
	}
	if _, ok := f.value.Interface().(Help); ok {
		return false
	}
	return true
}

// IsMulti returns true if the flag can be specified multiple times.
func (f *Flag) IsMulti() bool {
	if f.value.Kind() == reflect.Slice {
		return true
	}
	return false
}

func isShort(arg string) bool {
	return strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) >= 2
}

func isLong(arg string) bool {
	return strings.HasPrefix(arg, "--") && len(arg) >= 3
}

func (f *Flag) Handles(arg string) bool {
	return (isShort(arg) && arg[1:2] == f.Short) ||
		(isLong(arg) && arg[2:] == f.Long)

}

func (f *Flag) Parse(args []string) ([]string, error) {
	param, value := args[0], ""
	if f.NeedsExtraValue() &&
		(len(args) < 2 || (isShort(param) && len(param) > 2)) {
		return args, fmt.Errorf("Flag %s needs an argument", f.Name())
	}
	if f.WasSpecified && !f.IsMulti() {
		return args, fmt.Errorf("Flag %s can only be specified once", f.Name())
	}
	if isShort(param) && len(param) > 2 {
		// Short flag cluster
		args[0] = "-" + param[2:]
	} else if f.NeedsExtraValue() {
		value = args[1]
		args = args[2:]
	} else {
		args = args[1:]
	}
	f.WasSpecified = true
	return args, f.setValue(value)
}
