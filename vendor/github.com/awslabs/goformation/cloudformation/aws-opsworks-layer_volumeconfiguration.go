package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksLayer_VolumeConfiguration AWS CloudFormation Resource (AWS::OpsWorks::Layer.VolumeConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html
type AWSOpsWorksLayer_VolumeConfiguration struct {

	// Iops AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html#cfn-opsworks-layer-volconfig-iops
	Iops *Value `json:"Iops,omitempty"`

	// MountPoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html#cfn-opsworks-layer-volconfig-mountpoint
	MountPoint *Value `json:"MountPoint,omitempty"`

	// NumberOfDisks AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html#cfn-opsworks-layer-volconfig-numberofdisks
	NumberOfDisks *Value `json:"NumberOfDisks,omitempty"`

	// RaidLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html#cfn-opsworks-layer-volconfig-raidlevel
	RaidLevel *Value `json:"RaidLevel,omitempty"`

	// Size AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html#cfn-opsworks-layer-volconfig-size
	Size *Value `json:"Size,omitempty"`

	// VolumeType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-volumeconfiguration.html#cfn-opsworks-layer-volconfig-volumetype
	VolumeType *Value `json:"VolumeType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksLayer_VolumeConfiguration) AWSCloudFormationType() string {
	return "AWS::OpsWorks::Layer.VolumeConfiguration"
}

func (r *AWSOpsWorksLayer_VolumeConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
