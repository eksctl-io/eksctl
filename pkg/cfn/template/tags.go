package template

import (
	"fmt"
	"reflect"
)

// Tag is a CloudFormation tag
type Tag struct {
	Key   interface{}
	Value interface{}

	PropagateAtLaunch string
}

// maybeSetNameTag adds a Name tag to any resource that supports tags
// it calls makeAutoNameTag to format the tag value
func maybeSetNameTag(name string, resource interface{}) {
	e := reflect.ValueOf(resource).Elem()
	if e.Kind() == reflect.Struct {
		f := e.FieldByName("Tags")
		if f.IsValid() && f.CanSet() {
			tag := reflect.ValueOf(makeAutoNameTag(name))
			if f.Type() == reflect.ValueOf([]Tag{}).Type() {
				f.Set(reflect.Append(f, tag))
			}
		}
	}
}

// makeAutoNameTag create a new Name tag in the following format:
// {Key: "Name", Value: !Sub "${AWS::StackName}/<logicalResourceName>"}
func makeAutoNameTag(suffix string) Tag {
	return Tag{
		Key:   NewString("Name"),
		Value: MakeFnSubString(fmt.Sprintf("${%s}/%s", StackName, suffix)),
	}
}
