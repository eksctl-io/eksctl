package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_LocalDeviceResourceData AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.LocalDeviceResourceData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localdeviceresourcedata.html
type AWSGreengrassResourceDefinitionVersion_LocalDeviceResourceData struct {

	// GroupOwnerSetting AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localdeviceresourcedata.html#cfn-greengrass-resourcedefinitionversion-localdeviceresourcedata-groupownersetting
	GroupOwnerSetting *AWSGreengrassResourceDefinitionVersion_GroupOwnerSetting `json:"GroupOwnerSetting,omitempty"`

	// SourcePath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localdeviceresourcedata.html#cfn-greengrass-resourcedefinitionversion-localdeviceresourcedata-sourcepath
	SourcePath *Value `json:"SourcePath,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_LocalDeviceResourceData) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.LocalDeviceResourceData"
}

func (r *AWSGreengrassResourceDefinitionVersion_LocalDeviceResourceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
