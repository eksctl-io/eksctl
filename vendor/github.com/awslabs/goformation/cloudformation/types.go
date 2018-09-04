package cloudformation

import (
	"encoding/json"
)

type StringIntrinsic struct {
	stringValue    *string
	intrinsicValue *intrinsic
}

type intrinsic = map[string]interface{}

func NewString(value string) *StringIntrinsic {
	return &StringIntrinsic{
		stringValue: &value,
	}
}

func NewStringIntrinsic(key string, value interface{}) *StringIntrinsic {
	return &StringIntrinsic{
		intrinsicValue: &intrinsic{
			key: value,
		},
	}
}

func NewStringRef(value string) *StringIntrinsic {
	return NewStringIntrinsic("Ref", value)
}

func (r StringIntrinsic) value() interface{} {
	if r.stringValue != nil {
		return r.stringValue
	}
	if r.intrinsicValue != nil {
		return r.intrinsicValue
	}
	return nil
}

func (r *StringIntrinsic) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value())
}

func (r *StringIntrinsic) UnmarshalJSON(b []byte) error {
	var typecheck interface{}
	if err := json.Unmarshal(b, &typecheck); err != nil {
		return err
	}
	switch val := typecheck.(type) {
	case string:
		r.stringValue = &val
	case map[string]interface{}:
		r.intrinsicValue = &val
	}
	return nil
}
