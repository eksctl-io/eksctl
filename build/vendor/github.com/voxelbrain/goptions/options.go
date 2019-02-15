package goptions

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type optionFunc func(f *Flag, option, value string) error
type optionMap map[string]optionFunc

var (
	typeOptionMap = map[reflect.Type]optionMap{
		// Global options
		nil: optionMap{
			"description": description,
			"obligatory":  obligatory,
			"mutexgroup":  mutexgroup,
		},
		reflect.TypeOf(new(*os.File)).Elem(): optionMap{
			"create": initOptionMeta(file_create, "file_mode", 0),
			"append": initOptionMeta(file_append, "file_mode", 0),
			"rdonly": initOptionMeta(file_rdonly, "file_mode", 0),
			"wronly": initOptionMeta(file_wronly, "file_mode", 0),
			"rdwr":   initOptionMeta(file_rdwr, "file_mode", 0),
			"excl":   initOptionMeta(file_excl, "file_mode", 0),
			"sync":   initOptionMeta(file_sync, "file_mode", 0),
			"trunc":  initOptionMeta(file_trunc, "file_mode", 0),
			"perm":   file_perm,
		},
	}
)

// Wraps another optionFunc and inits optionMeta[field] with value if it does
// not have one already.
func initOptionMeta(fn optionFunc, field string, init_value interface{}) optionFunc {
	return func(f *Flag, option, value string) error {
		if _, ok := f.optionMeta[field]; !ok {
			f.optionMeta[field] = init_value
		}
		return fn(f, option, value)
	}
}

func description(f *Flag, option, value string) error {
	f.Description = strings.Replace(value, `\`, ``, -1)
	return nil
}

func obligatory(f *Flag, option, value string) error {
	f.Obligatory = true
	return nil
}

func mutexgroup(f *Flag, option, value string) error {
	if len(value) <= 0 {
		return fmt.Errorf("Mutexgroup option needs a value")
	}
	for _, group := range strings.Split(value, ",") {
		f.MutexGroups = append(f.MutexGroups, group)
	}
	return nil
}

func file_create(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_CREATE
	return nil
}

func file_append(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_APPEND
	return nil
}

func file_rdonly(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_RDONLY
	return nil
}

func file_wronly(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_WRONLY
	return nil
}

func file_rdwr(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_RDWR
	return nil
}

func file_excl(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_EXCL
	return nil
}

func file_sync(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_SYNC
	return nil
}

func file_trunc(f *Flag, option, value string) error {
	f.optionMeta["file_mode"] = f.optionMeta["file_mode"].(int) | os.O_TRUNC
	return nil
}

func file_perm(f *Flag, option, value string) error {
	perm, err := strconv.ParseInt(value, 8, 32)
	if err != nil {
		return err
	}
	f.optionMeta["file_perm"] = uint32(perm)
	return nil
}

func optionMapForType(t reflect.Type) optionMap {
	g := typeOptionMap[nil]
	m, _ := typeOptionMap[t]
	r := make(optionMap)
	for k, v := range g {
		r[k] = v
	}
	for k, v := range m {
		r[k] = v
	}
	return r
}
