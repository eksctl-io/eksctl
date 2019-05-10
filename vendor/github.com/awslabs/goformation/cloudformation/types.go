package cloudformation

import (
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

	Ref           = "Ref"
	FnBase64      = "Fn::Base64"
	FnCIDR        = "Fn::Cidr"
	FnAnd         = "Fn::And"
	FnEquals      = "Fn::Equals"
	FnIf          = "Fn::If"
	FnNot         = "Fn::Not"
	FnOr          = "Fn::Or"
	FnFindInMap   = "Fn::FindInMap"
	FnGetAtt      = "Fn::GetAtt"
	FnGetAZs      = "Fn::GetAZs"
	FnImportValue = "Fn::ImportValue"
	FnJoin        = "Fn::Join"
	FnSub         = "Fn::Sub"
	FnSelect      = "Fn::Select"
	FnSplit       = "Fn::Split"
)

var (
	RefAccountID        = MakeRef(AccountID)
	RefNotificationARNs = MakeRef(NotificationARNs)
	RefNoValue          = MakeRef(NoValue)
	RefPartition        = MakeRef(Partition)
	RefRegion           = MakeRef(Region)
	RefStackID          = MakeRef(StackID)
	RefStackName        = MakeRef(StackName)
)

type Value struct {
	value json.Marshaler
}

func NewValue(v json.Marshaler) *Value { return &Value{value: v} }

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
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch raw.(type) {
	case string:
		v.value = String(raw.(string))
		return nil
	case float64:
		v.value = Double(raw.(float64))
	case bool:
		v.value = Boolean(raw.(bool))
	case map[string]interface{}:
		v.value = AnythingMap(raw.(map[string]interface{}))
	case []interface{}:
		v.value = AnythingSlice(raw.([]interface{}))
	default:
		return fmt.Errorf("cannot handle type %s", reflect.ValueOf(raw).Kind())
	}

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

type StringSlice []string

func NewStringSlice(v ...string) *Value { return NewValue(StringSlice(v)) }

func (v StringSlice) MarshalJSON() ([]byte, error) {
	x := []string(v)
	return json.Marshal(&x)
}

type AnythingMap map[string]interface{}

func (v AnythingMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
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

type AnythingSlice []interface{}

func (v AnythingSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

func (v AnythingSlice) Convert(obj interface{}) error {
	x := []interface{}{}
	x = v
	data, err := json.Marshal(x)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

type Long int64

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

type Intrinsic struct {
	Value map[string]interface{}
}

func MakeIntrinsic(k string, v interface{}) *Value {
	return NewValue(
		Intrinsic{
			map[string]interface{}{
				k: v,
			},
		},
	)
}

func MakeRef(r string) *Value { return MakeIntrinsic(Ref, r) }

// TODO MakeFnBase64
// TODO MakeFnCIDR
// TODO MakeFnAnd
// TODO MakeFnEquals
// TODO MakeFnIf
// TODO MakeFnNot
// TODO MakeFnOr
// TODO MakeFnFindInMap

func MakeFnGetAtt(arg *Value) *Value       { return MakeIntrinsic(FnGetAtt, arg) }
func MakeFnGetAttString(arg string) *Value { return MakeFnGetAtt(NewString(arg)) }

// TODO MakeFnGetAZs

func MakeFnImportValue(arg *Value) *Value       { return MakeIntrinsic(FnImportValue, arg) }
func MakeFnImportValueString(arg string) *Value { return MakeFnImportValue(NewString(arg)) }

func MakeFnJoin(sep string, args []*Value) *Value {
	return MakeIntrinsic(FnJoin,
		[]interface{}{
			sep,
			args,
		},
	)
}

func MakeFnSub(arg *Value) *Value       { return MakeIntrinsic(FnSub, arg) }
func MakeFnSubString(arg string) *Value { return MakeFnSub(NewString(arg)) }

func MakeFnSelect(i int, arg *Value) *Value {
	return MakeIntrinsic(FnSelect,
		[]interface{}{
			i,
			arg,
		},
	)
}

func MakeFnSplit(sep string, arg *Value) *Value {
	return MakeIntrinsic(FnSplit,
		[]interface{}{
			sep,
			arg,
		},
	)
}

func MakeFnSplitString(sep string, arg string) *Value {
	return MakeFnSplit(sep, NewString(arg))
}

func MakeFnCIDR(arg *Value, n, len int) *Value {
	return MakeIntrinsic(FnCIDR,
		[]interface{}{
			arg,
			n,
			len,
		},
	)
}

func (v Intrinsic) MarshalJSON() ([]byte, error) { return json.Marshal(&v.Value) }
