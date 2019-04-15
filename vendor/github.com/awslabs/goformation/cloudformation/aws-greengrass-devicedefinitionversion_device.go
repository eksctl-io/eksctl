package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassDeviceDefinitionVersion_Device AWS CloudFormation Resource (AWS::Greengrass::DeviceDefinitionVersion.Device)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinitionversion-device.html
type AWSGreengrassDeviceDefinitionVersion_Device struct {

	// CertificateArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinitionversion-device.html#cfn-greengrass-devicedefinitionversion-device-certificatearn
	CertificateArn *Value `json:"CertificateArn,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinitionversion-device.html#cfn-greengrass-devicedefinitionversion-device-id
	Id *Value `json:"Id,omitempty"`

	// SyncShadow AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinitionversion-device.html#cfn-greengrass-devicedefinitionversion-device-syncshadow
	SyncShadow *Value `json:"SyncShadow,omitempty"`

	// ThingArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinitionversion-device.html#cfn-greengrass-devicedefinitionversion-device-thingarn
	ThingArn *Value `json:"ThingArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassDeviceDefinitionVersion_Device) AWSCloudFormationType() string {
	return "AWS::Greengrass::DeviceDefinitionVersion.Device"
}

func (r *AWSGreengrassDeviceDefinitionVersion_Device) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
