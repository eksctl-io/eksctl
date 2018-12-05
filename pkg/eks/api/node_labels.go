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

func (f *NodeLabels) Type() string {
	return "NodeLabels"
}
