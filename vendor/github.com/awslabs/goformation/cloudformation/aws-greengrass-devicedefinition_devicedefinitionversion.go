package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassDeviceDefinition_DeviceDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::DeviceDefinition.DeviceDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-devicedefinitionversion.html
type AWSGreengrassDeviceDefinition_DeviceDefinitionVersion struct {

	// Devices AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-devicedefinition-devicedefinitionversion.html#cfn-greengrass-devicedefinition-devicedefinitionversion-devices
	Devices []AWSGreengrassDeviceDefinition_Device `json:"Devices,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassDeviceDefinition_DeviceDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::DeviceDefinition.DeviceDefinitionVersion"
}

func (r *AWSGreengrassDeviceDefinition_DeviceDefinitionVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
