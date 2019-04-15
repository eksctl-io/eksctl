package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_LocalVolumeResourceData AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.LocalVolumeResourceData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localvolumeresourcedata.html
type AWSGreengrassResourceDefinitionVersion_LocalVolumeResourceData struct {

	// DestinationPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localvolumeresourcedata.html#cfn-greengrass-resourcedefinitionversion-localvolumeresourcedata-destinationpath
	DestinationPath *Value `json:"DestinationPath,omitempty"`

	// GroupOwnerSetting AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localvolumeresourcedata.html#cfn-greengrass-resourcedefinitionversion-localvolumeresourcedata-groupownersetting
	GroupOwnerSetting *AWSGreengrassResourceDefinitionVersion_GroupOwnerSetting `json:"GroupOwnerSetting,omitempty"`

	// SourcePath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-localvolumeresourcedata.html#cfn-greengrass-resourcedefinitionversion-localvolumeresourcedata-sourcepath
	SourcePath *Value `json:"SourcePath,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_LocalVolumeResourceData) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.LocalVolumeResourceData"
}

func (r *AWSGreengrassResourceDefinitionVersion_LocalVolumeResourceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
