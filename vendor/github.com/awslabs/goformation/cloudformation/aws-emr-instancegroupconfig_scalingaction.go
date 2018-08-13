package cloudformation

import (
	"encoding/json"
)

// AWSEMRInstanceGroupConfig_ScalingAction AWS CloudFormation Resource (AWS::EMR::InstanceGroupConfig.ScalingAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancegroupconfig-scalingaction.html
type AWSEMRInstanceGroupConfig_ScalingAction struct {

	// Market AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancegroupconfig-scalingaction.html#cfn-elasticmapreduce-instancegroupconfig-scalingaction-market
	Market *Value `json:"Market,omitempty"`

	// SimpleScalingPolicyConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancegroupconfig-scalingaction.html#cfn-elasticmapreduce-instancegroupconfig-scalingaction-simplescalingpolicyconfiguration
	SimpleScalingPolicyConfiguration *AWSEMRInstanceGroupConfig_SimpleScalingPolicyConfiguration `json:"SimpleScalingPolicyConfiguration,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRInstanceGroupConfig_ScalingAction) AWSCloudFormationType() string {
	return "AWS::EMR::InstanceGroupConfig.ScalingAction"
}

func (r *AWSEMRInstanceGroupConfig_ScalingAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
