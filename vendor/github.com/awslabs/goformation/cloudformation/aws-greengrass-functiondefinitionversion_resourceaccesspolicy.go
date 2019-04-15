package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_ResourceAccessPolicy AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.ResourceAccessPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-resourceaccesspolicy.html
type AWSGreengrassFunctionDefinitionVersion_ResourceAccessPolicy struct {

	// Permission AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-resourceaccesspolicy.html#cfn-greengrass-functiondefinitionversion-resourceaccesspolicy-permission
	Permission *Value `json:"Permission,omitempty"`

	// ResourceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-resourceaccesspolicy.html#cfn-greengrass-functiondefinitionversion-resourceaccesspolicy-resourceid
	ResourceId *Value `json:"ResourceId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_ResourceAccessPolicy) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.ResourceAccessPolicy"
}

func (r *AWSGreengrassFunctionDefinitionVersion_ResourceAccessPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
