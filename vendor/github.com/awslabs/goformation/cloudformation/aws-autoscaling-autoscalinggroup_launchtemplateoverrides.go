package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingAutoScalingGroup_LaunchTemplateOverrides AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.LaunchTemplateOverrides)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-launchtemplateoverrides.html
type AWSAutoScalingAutoScalingGroup_LaunchTemplateOverrides struct {

	// InstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-launchtemplateoverrides.html#cfn-autoscaling-autoscalinggroup-launchtemplateoverrides-instancetype
	InstanceType *Value `json:"InstanceType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingAutoScalingGroup_LaunchTemplateOverrides) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.LaunchTemplateOverrides"
}

func (r *AWSAutoScalingAutoScalingGroup_LaunchTemplateOverrides) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
