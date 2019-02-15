package goptions

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	_LONG_FLAG_REGEXP     = `--[[:word:]-]+`
	_SHORT_FLAG_REGEXP    = `-[[:alnum:]]`
	_QUOTED_STRING_REGEXP = `'((?:\\'|[^\\'])+)'`
	_OPTION_REGEXP        = `([[:word:]-]+)(?:=` + _QUOTED_STRING_REGEXP + `)?`
)

var (
	optionRegexp = regexp.MustCompile(`^(` + strings.Join([]string{_SHORT_FLAG_REGEXP, _LONG_FLAG_REGEXP, _OPTION_REGEXP}, "|") + `)(?:,|$)`)
)

func parseStructField(fieldValue reflect.Value, tag string) (*Flag, error) {
	f := &Flag{
		value:        fieldValue,
		DefaultValue: fieldValue.Interface(),
		optionMeta:   make(map[string]interface{}),
	}
	for {
		tag = strings.TrimSpace(tag)
		if len(tag) == 0 {
			break
		}
		idx := optionRegexp.FindStringSubmatchIndex(tag)
		if idx == nil {
			return nil, fmt.Errorf("Could not find a valid flag definition at the beginning of \"%s\"", tag)
		}
		option := tag[idx[2]:idx[3]]

		if strings.HasPrefix(option, "--") {
			if f.Long != "" {
				return nil, fmt.Errorf("Multiple flags assigned to a member: %s", strings.Join([]string{"--" + f.Long, option}, ", "))
			}
			f.Long = option[2:]
		} else if strings.HasPrefix(option, "-") {
			if f.Short != "" {
				return nil, fmt.Errorf("Multiple flags assigned to a member: %s", strings.Join([]string{"-" + f.Short, option}, ", "))
			}
			f.Short = option[1:]
		} else {
			option := tag[idx[4]:idx[5]]
			value := ""
			if idx[6] != -1 {
				value = tag[idx[6]:idx[7]]
			}
			optionmap := optionMapForType(fieldValue.Type())
			opf, ok := optionmap[option]
			if !ok {
				return nil, fmt.Errorf("Unknown option %s", option)
			}
			err := opf(f, option, value)
			if err != nil {
				return nil, fmt.Errorf("Option %s invalid: %s", option, err)
			}
		}
		// Keep remainder
		tag = tag[idx[1]:]
	}
	return f, nil
}
