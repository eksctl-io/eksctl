package api

import (
	"strings"
	"fmt"
	"errors"
)

func (f NodeLabels) String() string {
	s := ""
	for k, v := range f {
		if s != "" {
			s += ","
		}
		s += k + "=" + v
	}
	return s
}

// Set parses the value and add the resulting key-value pairs as node labels.
// Each value is expected to include one or more key-value pairs, formatted as `KEY_1=VAL_1,KEY_2=VAL_2,...`.
// That is, the key and the value must be delimited by an equal sign, while each pair must be delimited by a comma.
func (f *NodeLabels) Set(value string) error {
	if *f == nil {
		*f = map[string]string{}
	}
	kvs := strings.Split(value, ",")
	for i := range kvs {
		kv := strings.Split(kvs[i], "=")
		if len(kv) != 2 {
			return fmt.Errorf("node label must be formatted K=V, but it was: %s", kvs[i])
		}
		if kv[0] == "" {
			return errors.New("key must not be omitted in --node-labels")
		}
		(*f)[kv[0]] = kv[1]
	}
	return nil
}

// Type returns the name of the type as a string
func (f *NodeLabels) Type() string {
	return "NodeLabels"
}
