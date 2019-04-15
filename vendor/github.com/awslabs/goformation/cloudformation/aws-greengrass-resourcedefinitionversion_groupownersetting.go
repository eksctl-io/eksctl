package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_GroupOwnerSetting AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.GroupOwnerSetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-groupownersetting.html
type AWSGreengrassResourceDefinitionVersion_GroupOwnerSetting struct {

	// AutoAddGroupOwner AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-groupownersetting.html#cfn-greengrass-resourcedefinitionversion-groupownersetting-autoaddgroupowner
	AutoAddGroupOwner *Value `json:"AutoAddGroupOwner,omitempty"`

	// GroupOwner AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-groupownersetting.html#cfn-greengrass-resourcedefinitionversion-groupownersetting-groupowner
	GroupOwner *Value `json:"GroupOwner,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_GroupOwnerSetting) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.GroupOwnerSetting"
}

func (r *AWSGreengrassResourceDefinitionVersion_GroupOwnerSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
