package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
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

func MakeFnCIDR(arg, n, l *Value) *Value {
	return MakeIntrinsic(FnCIDR,
		[]*Value{
			arg,
			n,
			l,
		},
	)
}

func MakeFnEquals(value1, value2 *Value) *Value {
	return MakeIntrinsic(FnEquals, []*Value{value1, value2})
}

func MakeFnGetAtt(name string, attr *Value) *Value {
	return MakeIntrinsic(FnGetAtt, []interface{}{name, attr})
}
func MakeFnGetAttString(name string, attr string) *Value {
	return MakeFnGetAtt(name, NewString(attr))
}

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

func MakeFnSubString(arg string) *Value { return MakeIntrinsic(FnSub, arg) }

func MakeFnSelect(i, arg *Value) *Value {
	return MakeIntrinsic(FnSelect,
		[]interface{}{
			i,
			arg,
		},
	)
}

func MakeFnSplit(sep string, arg *Value) *Value {
	return MakeIntrinsic(FnSplit,
		[]interface{}{sep, arg},
	)
}

func MakeFnSplitString(sep string, arg string) *Value {
	return MakeFnSplit(sep, NewString(arg))
}

func MakeFnFindInMap(mapName *Value, topLevelKey *Value, secondLevelKey *Value) *Value {
	return MakeIntrinsic(
		FnFindInMap,
		[]interface{}{mapName, topLevelKey, secondLevelKey},
	)
}

// TODO MakeFnBase64
// TODO MakeFnGetAZs
// TODO MakeFnIf
// TODO MakeFnNot
// TODO MakeFnAnd
// TODO MakeFnOr

func (v Intrinsic) MarshalJSON() ([]byte, error) { return json.Marshal(&v.Value) }

func makeIntrinsic(typ string, arg interface{}) (*Value, error) {
	switch typ {
	case Ref:
		switch a := arg.(type) {
		case string:
			return MakeRef(a), nil
		default:
			return nil, fmt.Errorf("Cannot create Ref")
		}
	case FnGetAtt:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 2 {
				return nil, fmt.Errorf("Cannot parse FnGetAtt")
			}
			logicalName, ok := a[0].(string)
			if !ok {
				return nil, fmt.Errorf("Cannot parse FnGetAtt")
			}
			attr, err := NewValueFromPrimitive(a[1])
			if err != nil {
				return nil, err
			}
			return MakeFnGetAtt(logicalName, attr), nil
		case string:
			splits := strings.Split(a, ".")
			if len(splits) != 2 {
				return nil, fmt.Errorf("Cannot parse parameters for FnGetAtt")
			}
			return MakeFnGetAtt(splits[0], NewString(splits[1])), nil
		default:
			return nil, fmt.Errorf("Cannot parse FnGetAtt")
		}
	case FnImportValue:
		argV, err := NewValueFromPrimitive(arg)
		if err != nil {
			return nil, err
		}
		return MakeIntrinsic(typ, argV), nil
	case FnCIDR:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 3 {
				return nil, fmt.Errorf("Cannot parse FnGetAtt")
			}
			var vs []*Value
			for _, v := range a {
				vv, err := NewValueFromPrimitive(v)
				if err != nil {
					return nil, err
				}
				vs = append(vs, vv)
			}
			return MakeFnCIDR(vs[0], vs[1], vs[2]), nil
		}
		return nil, fmt.Errorf("Expected array as argument for Fn::CIDR")
	case FnJoin:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 2 {
				return nil, fmt.Errorf("Couldn't parse Fn::Join")
			}
			del, ok := a[0].(string)
			if !ok {
				return nil, fmt.Errorf("Couldn't parse delimiter")
			}
			rawVs, ok := a[1].([]interface{})
			if !ok {
				return nil, fmt.Errorf("Couldn't parse delimiter")
			}
			var vs []*Value
			for _, v := range rawVs {
				vv, err := NewValueFromPrimitive(v)
				if err != nil {
					return nil, err
				}
				vs = append(vs, vv)
			}
			return MakeFnJoin(del, vs), nil
		}
		return nil, fmt.Errorf("Expected array as argument for Fn::CIDR")
	case FnEquals:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 2 {
				return nil, fmt.Errorf("Cannot parse FnEquals")
			}
			var vs []*Value
			for _, v := range a {
				vv, err := NewValueFromPrimitive(v)
				if err != nil {
					return nil, err
				}
				vs = append(vs, vv)
			}
			return MakeFnEquals(vs[0], vs[1]), nil
		}
		return nil, fmt.Errorf("Expected array as argument for Fn::CIDR")
	case FnSelect:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 2 {
				return nil, fmt.Errorf("Cannot parse Fn::Select")
			}
			var vs []*Value
			for _, v := range a {
				vv, err := NewValueFromPrimitive(v)
				if err != nil {
					return nil, err
				}
				vs = append(vs, vv)
			}
			return MakeFnSelect(vs[0], vs[1]), nil
		}
		return nil, fmt.Errorf("Expected array as argument for Fn::Select")
	case FnSub:
		return MakeIntrinsic(typ, arg), nil
	case FnSplit:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 2 {
				return nil, fmt.Errorf("Couldn't parse Fn::Split")
			}
			del, ok := a[0].(string)
			if !ok {
				return nil, fmt.Errorf("Couldn't parse delimiter")
			}
			src, err := NewValueFromPrimitive(a[1])
			if err != nil {
				return nil, err
			}
			return MakeFnSplit(del, src), nil
		}
		return nil, fmt.Errorf("Expected array as argument for Fn::CIDR")
	case FnFindInMap:
		switch a := arg.(type) {
		case []interface{}:
			if len(a) != 3 {
				return nil, fmt.Errorf("Cannot parse %s intrinsic arguments as list", FnFindInMap)
			}
			var vs []*Value
			for _, v := range a {
				vv, err := NewValueFromPrimitive(v)
				if err != nil {
					return nil, err
				}
				vs = append(vs, vv)
			}
			return MakeFnFindInMap(vs[0], vs[1], vs[2]), nil
		}
		return nil, fmt.Errorf("Expected array as argument for %s", FnFindInMap)
	// "Unsupported" intrinsics
	case FnBase64, FnGetAZs, FnAnd, FnIf, FnNot, FnOr:
		return MakeIntrinsic(typ, arg), nil
	default:
		return nil, nil
	}
}

// IntrinsicFromMap attempts to parse a map to an intrinsic function
func IntrinsicFromMap(raw map[string]interface{}) (*Value, error) {
	if len(raw) != 1 {
		return nil, fmt.Errorf("map must contain one key")
	}
	for k, v := range raw {
		intr, err := makeIntrinsic(k, v)
		if err != nil {
			return nil, err
		}
		if intr != nil {
			return intr, nil
		}
		return nil, fmt.Errorf("map didn't contain intrinsic function key")
	}
	panic("Unreachable!")
}
