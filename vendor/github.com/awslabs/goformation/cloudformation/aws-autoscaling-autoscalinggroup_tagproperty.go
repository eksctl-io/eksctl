package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingAutoScalingGroup_TagProperty AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.TagProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-tags.html
type AWSAutoScalingAutoScalingGroup_TagProperty struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-tags.html#cfn-as-tags-Key
	Key *Value `json:"Key,omitempty"`

	// PropagateAtLaunch AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-tags.html#cfn-as-tags-PropagateAtLaunch
	PropagateAtLaunch *Value `json:"PropagateAtLaunch,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-tags.html#cfn-as-tags-Value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingAutoScalingGroup_TagProperty) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.TagProperty"
}

func (r *AWSAutoScalingAutoScalingGroup_TagProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
