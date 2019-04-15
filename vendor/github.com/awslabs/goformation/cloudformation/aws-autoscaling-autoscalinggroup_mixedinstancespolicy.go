package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingAutoScalingGroup_MixedInstancesPolicy AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.MixedInstancesPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-group-mixedinstancespolicy.html
type AWSAutoScalingAutoScalingGroup_MixedInstancesPolicy struct {

	// InstancesDistribution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-group-mixedinstancespolicy.html#cfn-as-mixedinstancespolicy-instancesdistribution
	InstancesDistribution *AWSAutoScalingAutoScalingGroup_InstancesDistribution `json:"InstancesDistribution,omitempty"`

	// LaunchTemplate AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-group-mixedinstancespolicy.html#cfn-as-mixedinstancespolicy-launchtemplate
	LaunchTemplate *AWSAutoScalingAutoScalingGroup_LaunchTemplate `json:"LaunchTemplate,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingAutoScalingGroup_MixedInstancesPolicy) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.MixedInstancesPolicy"
}

func (r *AWSAutoScalingAutoScalingGroup_MixedInstancesPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
