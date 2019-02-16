package goptions

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"
)

type valueParser func(f *Flag, val string) (reflect.Value, error)

var (
	parserMap = map[reflect.Type]valueParser{
		reflect.TypeOf(new(bool)).Elem():          boolValueParser,
		reflect.TypeOf(new(string)).Elem():        stringValueParser,
		reflect.TypeOf(new(float64)).Elem():       float64ValueParser,
		reflect.TypeOf(new(float32)).Elem():       float32ValueParser,
		reflect.TypeOf(new(int)).Elem():           intValueParser,
		reflect.TypeOf(new(int64)).Elem():         int64ValueParser,
		reflect.TypeOf(new(int32)).Elem():         int32ValueParser,
		reflect.TypeOf(new(Help)).Elem():          helpValueParser,
		reflect.TypeOf(new(*os.File)).Elem():      fileValueParser,
		reflect.TypeOf(new(*net.TCPAddr)).Elem():  tcpAddrValueParser,
		reflect.TypeOf(new(*url.URL)).Elem():      urlValueParser,
		reflect.TypeOf(new(time.Duration)).Elem(): durationValueParser,
	}
)

func parseMarshalValue(value reflect.Value, s string) error {
	newval := reflect.New(value.Type()).Elem()
	if newval.Kind() == reflect.Ptr {
		newptrval := reflect.New(value.Type().Elem())
		newval.Set(newptrval)
	}
	err := newval.Interface().(Marshaler).MarshalGoption(s)
	value.Set(newval)
	return err
}

func (f *Flag) setValue(s string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	if f.value.Type().Implements(reflect.TypeOf(new(Marshaler)).Elem()) {
		return parseMarshalValue(f.value, s)
	}
	vtype := f.value.Type()
	newval := reflect.New(vtype).Elem()
	if f.value.Kind() == reflect.Slice {
		vtype = f.value.Type().Elem()
		if vtype.Implements(reflect.TypeOf(new(Marshaler)).Elem()) {
			newval = reflect.New(vtype).Elem()
			err := parseMarshalValue(newval, s)
			f.value.Set(reflect.Append(f.value, newval))
			return err
		}
	}
	if parser, ok := parserMap[vtype]; ok {
		val, err := parser(f, s)
		if err != nil {
			return err
		}
		if f.value.Kind() == reflect.Slice {
			f.value.Set(reflect.Append(f.value, val))
		} else {
			f.value.Set(val)
		}
		return nil
	} else {
		return fmt.Errorf("Unsupported flag type: %s", f.value.Type())
	}
	panic("Invalid execution path")
}

func boolValueParser(f *Flag, val string) (reflect.Value, error) {
	return reflect.ValueOf(true), nil
}

func stringValueParser(f *Flag, val string) (reflect.Value, error) {
	return reflect.ValueOf(val), nil
}

func float64ValueParser(f *Flag, val string) (reflect.Value, error) {
	floatval, err := strconv.ParseFloat(val, 64)
	return reflect.ValueOf(float64(floatval)), err
}

func float32ValueParser(f *Flag, val string) (reflect.Value, error) {
	floatval, err := strconv.ParseFloat(val, 32)
	return reflect.ValueOf(float32(floatval)), err
}

func int64ValueParser(f *Flag, val string) (reflect.Value, error) {
	intval, err := strconv.ParseInt(val, 10, 64)
	return reflect.ValueOf(int64(intval)), err
}

func int32ValueParser(f *Flag, val string) (reflect.Value, error) {
	intval, err := strconv.ParseInt(val, 10, 32)
	return reflect.ValueOf(int32(intval)), err
}

func intValueParser(f *Flag, val string) (reflect.Value, error) {
	intval, err := strconv.ParseInt(val, 10, 64)
	return reflect.ValueOf(int(intval)), err
}

func fileValueParser(f *Flag, val string) (reflect.Value, error) {
	mode := 0
	if v, ok := f.optionMeta["file_mode"]; ok {
		mode = v.(int)
	}
	if val == "-" {
		if mode&1 == os.O_RDONLY {
			return reflect.ValueOf(os.Stdin), nil
		} else if mode&1 == os.O_WRONLY {
			return reflect.ValueOf(os.Stdout), nil
		}
	} else {
		perm := uint32(0644)
		if v, ok := f.optionMeta["file_perm"].(uint32); ok {
			perm = v
		}
		f, e := os.OpenFile(val, mode, os.FileMode(perm))
		return reflect.ValueOf(f), e
	}
	panic("Invalid execution path")
}

func tcpAddrValueParser(f *Flag, val string) (reflect.Value, error) {
	addr, err := net.ResolveTCPAddr("tcp", val)
	return reflect.ValueOf(addr), err
}

func urlValueParser(f *Flag, val string) (reflect.Value, error) {
	url, err := url.Parse(val)
	return reflect.ValueOf(url), err
}

func durationValueParser(f *Flag, val string) (reflect.Value, error) {
	d, err := time.ParseDuration(val)
	return reflect.ValueOf(d), err
}

func helpValueParser(f *Flag, val string) (reflect.Value, error) {
	return reflect.Value{}, ErrHelpRequest
}
