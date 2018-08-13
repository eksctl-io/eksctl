package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_EbsBlockDeviceConfig AWS CloudFormation Resource (AWS::EMR::Cluster.EbsBlockDeviceConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-ebsblockdeviceconfig.html
type AWSEMRCluster_EbsBlockDeviceConfig struct {

	// VolumeSpecification AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-ebsblockdeviceconfig.html#cfn-elasticmapreduce-cluster-ebsblockdeviceconfig-volumespecification
	VolumeSpecification *AWSEMRCluster_VolumeSpecification `json:"VolumeSpecification,omitempty"`

	// VolumesPerInstance AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-ebsblockdeviceconfig.html#cfn-elasticmapreduce-cluster-ebsblockdeviceconfig-volumesperinstance
	VolumesPerInstance *Value `json:"VolumesPerInstance,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_EbsBlockDeviceConfig) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.EbsBlockDeviceConfig"
}

func (r *AWSEMRCluster_EbsBlockDeviceConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
