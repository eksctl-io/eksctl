package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassDeviceDefinition_Device AWS CloudFormation Resource (AWS::Greengrass::DeviceDefinition.Device)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-device.html
type AWSGreengrassDeviceDefinition_Device struct {

	// CertificateArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-device.html#cfn-greengrass-devicedefinition-device-certificatearn
	CertificateArn *Value `json:"CertificateArn,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-device.html#cfn-greengrass-devicedefinition-device-id
	Id *Value `json:"Id,omitempty"`

	// SyncShadow AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-device.html#cfn-greengrass-devicedefinition-device-syncshadow
	SyncShadow *Value `json:"SyncShadow,omitempty"`

	// ThingArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-device.html#cfn-greengrass-devicedefinition-device-thingarn
	ThingArn *Value `json:"ThingArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassDeviceDefinition_Device) AWSCloudFormationType() string {
	return "AWS::Greengrass::DeviceDefinition.Device"
}

func (r *AWSGreengrassDeviceDefinition_Device) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
