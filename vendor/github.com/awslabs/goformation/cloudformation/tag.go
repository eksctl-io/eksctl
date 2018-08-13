package cloudformation

import (
	"encoding/json"
)

// Tag AWS CloudFormation Resource (Tag)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-tag.html
type Tag struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-tag.html#cfn-dms-endpoint-tag-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-tag.html#cfn-dms-endpoint-tag-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Tag) AWSCloudFormationType() string {
	return "Tag"
}

func (r *Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
