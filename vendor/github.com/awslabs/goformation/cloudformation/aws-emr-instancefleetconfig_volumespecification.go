package cloudformation

import (
	"encoding/json"
)

// AWSEMRInstanceFleetConfig_VolumeSpecification AWS CloudFormation Resource (AWS::EMR::InstanceFleetConfig.VolumeSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-volumespecification.html
type AWSEMRInstanceFleetConfig_VolumeSpecification struct {

	// Iops AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-volumespecification.html#cfn-elasticmapreduce-instancefleetconfig-volumespecification-iops
	Iops *Value `json:"Iops,omitempty"`

	// SizeInGB AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-volumespecification.html#cfn-elasticmapreduce-instancefleetconfig-volumespecification-sizeingb
	SizeInGB *Value `json:"SizeInGB,omitempty"`

	// VolumeType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-volumespecification.html#cfn-elasticmapreduce-instancefleetconfig-volumespecification-volumetype
	VolumeType *Value `json:"VolumeType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRInstanceFleetConfig_VolumeSpecification) AWSCloudFormationType() string {
	return "AWS::EMR::InstanceFleetConfig.VolumeSpecification"
}

func (r *AWSEMRInstanceFleetConfig_VolumeSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
