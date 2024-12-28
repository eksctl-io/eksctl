package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	AccountID        = "AWS::AccountId"
	NotificationARNs = "AWS::NotificationARNs"
	NoValue          = "AWS::NoValue"
	Partition        = "AWS::Partition"
	Region           = "AWS::Region"
	StackID          = "AWS::StackId"
	StackName        = "AWS::StackName"
)

type Value struct {
	value json.Marshaler
}

func NewValue(v json.Marshaler) *Value { return &Value{value: v} }

func NewValueFromPrimitive(raw interface{}) (*Value, error) {
	switch p := raw.(type) {
	case string:
		return NewString(p), nil
	case float64:
		return NewDouble(p), nil
	case json.Number:
		i, err := p.Int64()
		if err == nil {
			if i <= int64(^uint(0) >> 1) {
				return NewInteger(int(i)), nil
			}
			return NewLong(i), nil
		}
		f, err := p.Float64()
		if err == nil {
			return NewDouble(f), nil
		}
		return NewString(p.String()), nil
	case int:
		return NewInteger(p), nil
	case int64:
		return NewLong(p), nil
	case bool:
		return NewBoolean(p), nil
	case map[string]interface{}:
		asIntr, err := IntrinsicFromMap(p)
		if err == nil {
			return asIntr, nil
		}
		return NewValue(AnythingMap(p)), nil
	case []interface{}:
		var vs []*Value
		for _, s := range p {
			v, err := NewValueFromPrimitive(s)
			if err != nil {
				return nil, err
			}
			vs = append(vs, v)
		}
		return NewSlice(vs...), nil
	default:
		return nil, fmt.Errorf("cannot handle type %s", reflect.ValueOf(p).Kind())
	}

}

func (v *Value) Raw() interface{} { return v.value }

func (v *Value) String() string {
	if v == nil {
		return "<nil>"
	}
	if x, ok := v.Raw().(string); ok {
		return x
	}
	return fmt.Sprintf("%v", v.value)
}

func (v *Value) UnmarshalJSON(b []byte) error {
	var raw interface{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
		return err
	}

	newV, err := NewValueFromPrimitive(raw)
	if err != nil {
		return err
	}
	*v = *newV
	return nil
}

func (v Value) MarshalJSON() ([]byte, error) {
	return v.value.MarshalJSON()
}

type String string

func NewString(v string) *Value { return NewValue(String(v)) }

func (v String) MarshalJSON() ([]byte, error) {
	x := string(v)
	return json.Marshal(&x)
}

type Slice []*Value

// NewSlice returns nil if it would otherwise be empty
func NewSlice(v ...*Value) *Value {
	if len(v) == 0 {
		return nil
	}
	return NewValue(Slice(v))
}

func (v Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal([]*Value(v))
}

func NewStringSlice(ss ...string) *Value {
	var vs []*Value
	for _, s := range ss {
		vs = append(vs, NewString(s))
	}
	return NewSlice(vs...)
}

type AnythingMap map[string]interface{}

func (v AnythingMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(v))
}

func (v AnythingMap) Convert(obj interface{}) error {
	x := map[string]interface{}{}
	x = v
	data, err := json.Marshal(x)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

type Long int64

func NewLong(v int64) *Value { return NewValue(Long(v)) }

func (v Long) MarshalJSON() ([]byte, error) {
	x := int64(v)
	return json.Marshal(&x)
}

type Integer int

func NewInteger(v int) *Value { return NewValue(Integer(v)) }

func (v Integer) MarshalJSON() ([]byte, error) {
	x := int(v)
	return json.Marshal(&x)
}

type Double float64

func NewDouble(v float64) *Value { return NewValue(Double(v)) }

func (v Double) MarshalJSON() ([]byte, error) {
	x := float64(v)
	return json.Marshal(&x)
}

type Boolean bool

var (
	trueValue  = NewBoolean(true)
	falseValue = NewBoolean(false)
)

func True() *Value  { return trueValue }
func False() *Value { return falseValue }

func NewBoolean(v bool) *Value { return NewValue(Boolean(v)) }

func (v Boolean) MarshalJSON() ([]byte, error) {
	x := bool(v)
	return json.Marshal(&x)
}
