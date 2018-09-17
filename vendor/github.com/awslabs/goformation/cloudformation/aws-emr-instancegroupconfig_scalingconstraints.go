package cloudformation

import (
	"encoding/json"
)

// AWSEMRInstanceGroupConfig_ScalingConstraints AWS CloudFormation Resource (AWS::EMR::InstanceGroupConfig.ScalingConstraints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancegroupconfig-scalingconstraints.html
type AWSEMRInstanceGroupConfig_ScalingConstraints struct {

	// MaxCapacity AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancegroupconfig-scalingconstraints.html#cfn-elasticmapreduce-instancegroupconfig-scalingconstraints-maxcapacity
	MaxCapacity *Value `json:"MaxCapacity,omitempty"`

	// MinCapacity AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancegroupconfig-scalingconstraints.html#cfn-elasticmapreduce-instancegroupconfig-scalingconstraints-mincapacity
	MinCapacity *Value `json:"MinCapacity,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRInstanceGroupConfig_ScalingConstraints) AWSCloudFormationType() string {
	return "AWS::EMR::InstanceGroupConfig.ScalingConstraints"
}

func (r *AWSEMRInstanceGroupConfig_ScalingConstraints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
