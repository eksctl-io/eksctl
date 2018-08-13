package cloudformation

import (
	"encoding/json"
)

// AWSEC2Instance_Volume AWS CloudFormation Resource (AWS::EC2::Instance.Volume)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-mount-point.html
type AWSEC2Instance_Volume struct {

	// Device AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-mount-point.html#cfn-ec2-mountpoint-device
	Device *Value `json:"Device,omitempty"`

	// VolumeId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-mount-point.html#cfn-ec2-mountpoint-volumeid
	VolumeId *Value `json:"VolumeId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2Instance_Volume) AWSCloudFormationType() string {
	return "AWS::EC2::Instance.Volume"
}

func (r *AWSEC2Instance_Volume) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
