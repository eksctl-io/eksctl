package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingAutoScalingGroup_LaunchTemplate AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.LaunchTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-launchtemplate.html
type AWSAutoScalingAutoScalingGroup_LaunchTemplate struct {

	// LaunchTemplateSpecification AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-launchtemplate.html#cfn-as-group-launchtemplate
	LaunchTemplateSpecification *AWSAutoScalingAutoScalingGroup_LaunchTemplateSpecification `json:"LaunchTemplateSpecification,omitempty"`

	// Overrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-launchtemplate.html#cfn-as-mixedinstancespolicy-overrides
	Overrides []AWSAutoScalingAutoScalingGroup_LaunchTemplateOverrides `json:"Overrides,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingAutoScalingGroup_LaunchTemplate) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.LaunchTemplate"
}

func (r *AWSAutoScalingAutoScalingGroup_LaunchTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
