package template

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Commonly-used constants
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

// Commonly-used references
var (
	RefAccountID        = MakeRef(AccountID)
	RefNotificationARNs = MakeRef(NotificationARNs)
	RefNoValue          = MakeRef(NoValue)
	RefPartition        = MakeRef(Partition)
	RefRegion           = MakeRef(Region)
	RefStackID          = MakeRef(StackID)
	RefStackName        = MakeRef(StackName)
)

// Value represents any JSON value
type Value struct {
	value json.Marshaler
}

// NewValue create a value
func NewValue(v json.Marshaler) *Value { return &Value{value: v} }

// Raw returns the value as an empty interface
func (v *Value) Raw() interface{} { return v.value }

// String representation of the value
func (v *Value) String() string {
	if v == nil {
		return "<nil>"
	}
	if x, ok := v.Raw().(string); ok {
		return x
	}
	return fmt.Sprintf("%v", v.value)
}

// UnmarshalJSON parses JSON data into a value
func (v *Value) UnmarshalJSON(b []byte) error {
	var raw interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch r := raw.(type) {
	case string:
		v.value = String(r)
		return nil
	case float64:
		v.value = Double(r)
	case bool:
		v.value = Boolean(r)
	case map[string]interface{}:
		v.value = AnythingMap(r)
	case []interface{}:
		v.value = AnythingSlice(r)
	default:
		return fmt.Errorf("cannot handle type %s", reflect.ValueOf(raw).Kind())
	}

	return nil
}

// MarshalJSON serialises the value as JSON data
func (v Value) MarshalJSON() ([]byte, error) {
	return v.value.MarshalJSON()
}

// String is a extension of string type
type String string

// NewString creates a String from v
func NewString(v string) *Value { return NewValue(String(v)) }

// MarshalJSON serialises the value as JSON data
func (v String) MarshalJSON() ([]byte, error) {
	x := string(v)
	return json.Marshal(&x)
}

// StringSlice represents a slice of strings
type StringSlice []string

// NewStringSlice constructs a representation of slice of strings from the given arguments
func NewStringSlice(v ...string) *Value { return NewValue(StringSlice(v)) }

// MarshalJSON serialises the value as JSON data
func (v StringSlice) MarshalJSON() ([]byte, error) {
	x := []string(v)
	return json.Marshal(&x)
}

// AnythingMap is an extension of map[string]interface{} type
type AnythingMap MapOfInterfaces

// MarshalJSON serialises the value as JSON data
func (v AnythingMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(MapOfInterfaces(v))
}

// Convert will serialise the receiver as JSON, and deserialise it into obj
func (v AnythingMap) Convert(obj interface{}) error {
	data, err := json.Marshal(MapOfInterfaces(v))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

// AnythingSlice is an extension of []interface{} type
type AnythingSlice SliceOfInterfaces

// MarshalJSON serialises the value as JSON data
func (v AnythingSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

// Convert will serialise the receiver as JSON, and deserialise it into obj
func (v AnythingSlice) Convert(obj interface{}) error {
	data, err := json.Marshal(SliceOfInterfaces(v))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

// Integer is an extension of int type
type Integer int

// NewInteger creates an Integer from v
func NewInteger(v int) *Value { return NewValue(Integer(v)) }

// MarshalJSON serialises the value as JSON data
func (v Integer) MarshalJSON() ([]byte, error) {
	x := int(v)
	return json.Marshal(&x)
}

// Double is an extension of float64 type
type Double float64

// NewDouble creates a Double from v
func NewDouble(v float64) *Value { return NewValue(Double(v)) }

// MarshalJSON serialises the value as JSON data
func (v Double) MarshalJSON() ([]byte, error) {
	x := float64(v)
	return json.Marshal(&x)
}

// Boolean is an extension of bool type
type Boolean bool

var (
	trueValue  = NewBoolean(true)
	falseValue = NewBoolean(false)
)

// True will return a new true value
func True() *Value { return trueValue }

// False will return a new false value
func False() *Value { return falseValue }

// NewBoolean creates a Boolean from v
func NewBoolean(v bool) *Value { return NewValue(Boolean(v)) }

// MarshalJSON serialises the value as JSON data
func (v Boolean) MarshalJSON() ([]byte, error) {
	x := bool(v)
	return json.Marshal(&x)
}

// Intrinsic represents an intrinsics
type Intrinsic struct {
	Value MapOfInterfaces
}

// MarshalJSON serialises the value as JSON data
func (v Intrinsic) MarshalJSON() ([]byte, error) { return json.Marshal(&v.Value) }

// MakeIntrinsic constructs an intrinsic
func MakeIntrinsic(k string, v interface{}) *Value {
	return NewValue(Intrinsic{Value: MapOfInterfaces{k: v}})
}

// MakeRef constructs a Ref intrinsic
func MakeRef(r string) *Value { return MakeIntrinsic(Ref, r) }

// MakeFnGetAtt constructs an Fn::GetAtt intrinsic
func MakeFnGetAtt(arg *Value) *Value { return MakeIntrinsic(FnGetAtt, arg) }

// MakeFnGetAttString constructs an Fn::GetAtt intrinsic for a given string
func MakeFnGetAttString(arg string) *Value { return MakeFnGetAtt(NewString(arg)) }

// MakeFnImportValue constructs an Fn::ImportValue intrinsic
func MakeFnImportValue(arg *Value) *Value { return MakeIntrinsic(FnImportValue, arg) }

// MakeFnImportValueString constructs an Fn::ImportValue intrinsic for a given string
func MakeFnImportValueString(arg string) *Value { return MakeFnImportValue(NewString(arg)) }

// MakeFnJoin constructs an Fn::Join intrinsic
func MakeFnJoin(sep string, args []*Value) *Value {
	return MakeIntrinsic(FnJoin, SliceOfInterfaces{sep, args})
}

// MakeFnSplit constructs an Fn::Split intrinsic
func MakeFnSplit(sep string, arg *Value) *Value {
	return MakeIntrinsic(FnSplit, SliceOfInterfaces{sep, arg})
}

// MakeFnSplitString constructs an Fn::Split intrinsic for a given string
func MakeFnSplitString(sep string, arg string) *Value {
	return MakeFnSplit(sep, NewString(arg))
}

// MakeFnSub constructs an Fn::Sub intrinsic
func MakeFnSub(arg *Value) *Value { return MakeIntrinsic(FnSub, arg) }

// MakeFnSubString constructs an Fn::Sub intrinsic for a given string
func MakeFnSubString(arg string) *Value { return MakeFnSub(NewString(arg)) }

// MakeFnSelect constructs an Fn::Select intrinsic
func MakeFnSelect(i int, arg *Value) *Value {
	return MakeIntrinsic(FnSelect, SliceOfInterfaces{i, arg})
}

// MakeFnCIDR constructs an Fn::Cidr intrinsic
func MakeFnCIDR(arg *Value, n, len int) *Value {
	return MakeIntrinsic(FnCIDR, SliceOfInterfaces{arg, n, len})
}

// MakeName is syntactic sugar for {"Fn::Sub": "${AWS::Stack}-<name>"}
func MakeName(suffix string) *Value {
	return MakeFnSubString(fmt.Sprintf("${%s}-%s", StackName, suffix))
}

// MakeSlice makes a slice from a list of arguments
func MakeSlice(i ...*Value) []*Value {
	return i
}

// MakeStringSlice makes a slice from a list of string arguments
func MakeStringSlice(s ...string) []*Value {
	slice := []*Value{}
	for _, i := range s {
		slice = append(slice, NewString(i))
	}
	return slice
}
